package psql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/delveper/mylib/app/ent"
	"github.com/delveper/mylib/app/exc"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
)

type Reader struct{ *sql.DB }

func NewReader(db *sql.DB) *Reader {
	return &Reader{db}
}

func (r Reader) Create(ctx context.Context, reader ent.Reader) error {
	reader.CreatedAt = time.Now()

	_, err := r.DB.ExecContext(ctx, queryCreateReader,
		reader.FirstName, // $1
		reader.LastName,  // $2
		reader.Email,     // $3
		reader.Password,  // $4
	)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("%w: %v", exc.ErrDeadline, err)
		}

		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) {
			switch pgxErr.ConstraintName {
			case "readers_email_key":
				return fmt.Errorf("%w: %v", exc.ErrDuplicateEmail, err)
			case "readers_pkey":
				return fmt.Errorf("%w: %v", exc.ErrDuplicateID, err)
			}
		}

		return fmt.Errorf("%w: %v", exc.ErrUnexpected, err)
	}

	return nil
}

func (r Reader) Authenticate(ctx context.Context, reader ent.Reader) error {
	var id, password string

	row := r.DB.QueryRowContext(ctx, queryAuthenticateReader, reader.Email)

	err := row.Scan(&id, &password)

	if err != nil {
		// switch {
		// case errors.Is(err, context.DeadlineExceeded):
		// 	return fmt.Errorf("%w: %v", exc.ErrDeadline, err)
		// case errors.Is(err, sql.ErrNoRows):
		// 	return fmt.Errorf("%w: %v", exc.ErrNoRecord, err)
		// default:
		// 	return fmt.Errorf("%w: %v", exc.ErrUnexpected, err)
		// }
		return fmt.Errorf("authentication error: %w", err)
	}

	reader = ent.Reader{ID: id, Password: password}

	return nil
}
