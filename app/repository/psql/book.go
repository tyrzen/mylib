package psql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/delveper/mylib/app/exceptions"
	"github.com/delveper/mylib/app/models"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
)

type Book struct{ *sql.DB }

func NewBook(db *sql.DB) *Book {
	return &Book{db}
}

func (b Book) Add(ctx context.Context, book models.Book) error {
	const SQL = `INSERT INTO books (id, author_id, title, genre, rate, size, created_at) 
					VALUES (GEN_RANDOM_UUID(), $1, $2, $3, $4, $5, NOW());`

	_, err := b.ExecContext(ctx, SQL,
		book.AuthorID, // $1
		book.Title,    // $2
		book.Genre,    // $3
		book.Rate,     // $4
		book.Size,     // $5
	)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("%w: %v", exceptions.ErrDeadline, err)
		}

		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) {
			switch pgxErr.ConstraintName {
			case "books_title_key":
				return fmt.Errorf("%w: %v", exceptions.ErrDuplicateTitle, err)
			case "books_pkey":
				return fmt.Errorf("%w: %v", exceptions.ErrDuplicateID, err)
			}
		}

		return fmt.Errorf("%w: %v", exceptions.ErrUnexpected, err)
	}

	return nil
}
