package psql

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5" // "github.com/jackc/pgx/v5/pgxpool"
)

func Connect() (*sql.DB, error) {
	DSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("REPO_HOST"),
		os.Getenv("REPO_PORT"),
		os.Getenv("REPO_USER"),
		os.Getenv("REPO_PASSWORD"),
		os.Getenv("REPO_NAME"),
	)

	d := os.Getenv("REPO_DIALECT") // pgx

	db, err := sql.Open(d, DSN)
	if err != nil {
		return nil, fmt.Errorf("error onening connection to repo: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging connection: %w", err)
	}

	return db, nil
}
