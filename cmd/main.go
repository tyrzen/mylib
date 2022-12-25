package main

import (
	"log"

	"github.com/delveper/mylib/app/ent"
	repo "github.com/delveper/mylib/app/repo/psql"
	"github.com/delveper/mylib/cfg"
	"github.com/delveper/mylib/lgr"
	"github.com/delveper/mylib/mig"
)

func main() {
	if err := cfg.Load(); err != nil {
		log.Printf("Failed to load environment variables: %+v", err)

		return
	}

	var logger ent.Logger = lgr.New()
	defer func() {
		if err := logger.Flush(); err != nil {
			log.Printf("Failed flush logger: %+v", err)
		}
	}()
	logger.Infof("Logger set up with level: %s", logger.Level())

	conn, err := repo.Connect()
	if err != nil {
		logger.Errorf("Failed connecting to repo: %+v", err)

		return
	}

	defer func() {
		if err := conn.Close(); err != nil {
			logger.Warnf("Failed closing connection: %+v", err)
		}
	}()
	logger.Infof("Connection to repo was established.")

	if err := mig.Migrate(conn); err != nil {
		logger.Errorf("Failed run repo migrations: %+v", err)

		return
	}

	logger.Infof("Migrations went successfully.")
}
