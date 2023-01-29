package mig

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"os"

	"github.com/delveper/mylib/app/models"
	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var fs embed.FS

const root = "."

type Migration struct{ *embed.FS }

func New() Migration {
	return Migration{&fs}
}

func (m Migration) SetLogger(logger models.Logger) {
	goose.SetLogger(logger)
}

func (m Migration) Run(db *sql.DB) error {
	goose.SetBaseFS(m.FS)

	defer func() {
		goose.SetBaseFS(nil)
	}()

	d, ok := os.LookupEnv("DB_DIALECT")
	if !ok {
		return fmt.Errorf("there is no DB_DIALECT among environment variables")
	}

	if err := goose.SetDialect(d); err != nil {
		return fmt.Errorf("error setting dialect: %w", err)
	}

	direction, ok := os.LookupEnv("DB_MIGRATE")
	if !ok {
		return errors.New("there no DB_MIGRATE above environment variables")
	}

	if err := goose.Fix(root); err != nil {
		return fmt.Errorf("error during fixing migrations timestamps: %w", err)
	}

	switch direction {
	case "up":
		if err := goose.Up(db, root); err != nil {
			return fmt.Errorf("error running mingrations up: %w", err)
		}
	case "down":
		if err := goose.Down(db, root); err != nil {
			return fmt.Errorf("error run mingrations down: %w", err)
		}
	default:
		return fmt.Errorf("wrong direction argument: %s", direction)
	}

	return nil
}
