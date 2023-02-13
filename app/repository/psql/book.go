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
	const SQL = `INSERT INTO books (id, author_id, title, genre, rate, size, year) 
					VALUES (GEN_RANDOM_UUID(), $1, $2, $3, $4, $5, $6);`

	_, err := b.ExecContext(ctx, SQL,
		book.AuthorID, // $1
		book.Title,    // $2
		book.Genre,    // $3
		book.Rate,     // $4
		book.Size,     // $5
		book.Year,     // $6
	)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("%w: %w", exceptions.ErrDeadline, err)
		}

		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) {
			switch pgxErr.ConstraintName {
			case "books_title_key":
				return fmt.Errorf("%w: %w", exceptions.ErrDuplicateTitle, err)
			case "books_pkey":
				return fmt.Errorf("%w: %w", exceptions.ErrDuplicateID, err)
			}
		}

		return fmt.Errorf("%w: %w", exceptions.ErrUnexpected, err)
	}

	return nil
}

func (b Book) GetByID(ctx context.Context, book models.Book) (models.Book, error) {
	const SQL = `SELECT id, author_id, title, genre, rate, size, year
				 FROM books 
				 WHERE id=$1;`

	row := b.QueryRowContext(ctx, SQL, book.ID)

	err := row.Scan(
		&book.ID,
		&book.AuthorID,
		&book.Title,
		&book.Genre,
		&book.Rate,
		&book.Size,
		&book.Year,
	)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return models.Book{}, fmt.Errorf("%w: %w", exceptions.ErrDeadline, err)
		case errors.Is(err, sql.ErrNoRows):
			return models.Book{}, fmt.Errorf("%w: %w", exceptions.ErrRecordNotFound, err)
		default:
			return models.Book{}, fmt.Errorf("%w: %w", exceptions.ErrUnexpected, err)
		}
	}

	return book, nil
}

func (b Book) GetMany(ctx context.Context, filter models.DataFilter) ([]models.Book, error) {
	const SQL = `SELECT id, author_ID, title, genre, rate, size, year
				 FROM books
				 `

	query := SQL + "\n" + evalQuery(filter)

	rows, err := b.QueryContext(ctx, query)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return nil, fmt.Errorf("%w: %w", exceptions.ErrDeadline, err)
		case errors.Is(err, sql.ErrNoRows):
			return nil, fmt.Errorf("%w: %w", exceptions.ErrRecordNotFound, err)
		default:
			return nil, fmt.Errorf("%w: %w", exceptions.ErrUnexpected, err)
		}
	}

	defer rows.Close()

	var books []models.Book

	for rows.Next() {
		var book models.Book

		err := rows.Scan(
			&book.ID,
			&book.AuthorID,
			&book.Title,
			&book.Genre,
			&book.Rate,
			&book.Size,
			&book.Year,
		)

		if err != nil {
			switch {
			case errors.Is(err, context.DeadlineExceeded):
				return nil, fmt.Errorf("%w: %w", exceptions.ErrDeadline, err)
			case errors.Is(err, sql.ErrNoRows):
				return nil, fmt.Errorf("%w: %w", exceptions.ErrRecordNotFound, err)
			default:
				return nil, fmt.Errorf("%w: %w", exceptions.ErrUnexpected, err)
			}
		}

		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during iteration: %w", err)
	}

	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("error while closing connection: %w", err)
	}

	return books, nil
}

func (b Book) AddToFavorites(ctx context.Context, reader models.Reader, book models.Book) error {
	const SQL = `INSERT INTO favorites (reader_id, book_id, created_at) 
					VALUES ($1, $2, NOW());`

	_, err := b.ExecContext(ctx, SQL,
		reader.ID, // $1
		book.ID,   // $2
	)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("%w: %w", exceptions.ErrDeadline, err)
		}

		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) {
			switch pgxErr.ConstraintName {
			case "favorites_reader_id_book_id_key":
				return fmt.Errorf("%w: %w", exceptions.ErrRecordExists, err)
			case "favorites_book_id_fkey":
				return fmt.Errorf("%w: %w", exceptions.ErrBookNotFound, err)
			case "favorites_reader_id_fkey":
				return fmt.Errorf("%w: %w", exceptions.ErrReaderNotFound, err)
			default:
			}
		}

		return fmt.Errorf("%w: %w", exceptions.ErrUnexpected, err)
	}

	return nil
}

func (b Book) AddToWishlist(ctx context.Context, reader models.Reader, book models.Book) error {
	const SQL = `INSERT INTO wishlist (reader_id, book_id, created_at) 
					VALUES ($1, $2, NOW());`

	_, err := b.ExecContext(ctx, SQL,
		reader.ID, // $1
		book.ID,   // $2
	)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("%w: %w", exceptions.ErrDeadline, err)
		}

		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) {
			switch pgxErr.ConstraintName {
			case "wishlist_reader_id_book_id_key":
				return fmt.Errorf("%w: %w", exceptions.ErrRecordExists, err)
			case "wishlist_book_id_fkey":
				return fmt.Errorf("%w: %w", exceptions.ErrBookNotFound, err)
			case "wishlist_reader_id_fkey":
				return fmt.Errorf("%w: %w", exceptions.ErrReaderNotFound, err)
			default:
			}
		}

		return fmt.Errorf("%w: %w", exceptions.ErrUnexpected, err)
	}

	return nil
}
