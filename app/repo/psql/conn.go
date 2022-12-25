package psql

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	DSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("REPO_HOST"),
		os.Getenv("REPO_PORT"),
		os.Getenv("REPO_USER"),
		os.Getenv("REPO_PASSWORD"),
		os.Getenv("REPO_NAME"),
	)

	db, err := sql.Open(os.Getenv("REPO_DIALECT"), DSN)
	if err != nil {
		return nil, fmt.Errorf("error onening connection to repo: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging conneciton: %w", err)
	}

	return db, nil
}
