package psql

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib" // "github.com/jackc/pgx/v5/pgxpool"
)

func Connect() (*sql.DB, error) {
	DSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_EXPOSE_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSL_MODE"),
	)

	d := os.Getenv("DB_DIALECT") // pgx

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
