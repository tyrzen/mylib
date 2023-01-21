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

// Add adds ent.Reader entity.
func (r Reader) Add(ctx context.Context, reader ent.Reader) error {
	reader.CreatedAt = time.Now()

	_, err := r.ExecContext(ctx, queryCreateReader,
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

// GetByEmailOrID will retrieve ent.Reader by given UserID, Email or both parameters.
func (r Reader) GetByEmailOrID(ctx context.Context, reader ent.Reader) (ent.Reader, error) {
	row := r.QueryRowContext(ctx, queryFindReader,
		reader.ID,    // $1
		reader.Email, // $2
	)

	err := row.Scan(
		&reader.ID,
		&reader.FirstName,
		&reader.LastName,
		&reader.Email,
		&reader.Password,
		&reader.CreatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return ent.Reader{}, fmt.Errorf("%w: %v", exc.ErrDeadline, err)
		case errors.Is(err, sql.ErrNoRows):
			return ent.Reader{}, fmt.Errorf("%w: %v", exc.ErrNoRecord, err)
		default:
			return ent.Reader{}, fmt.Errorf("%w: %v", exc.ErrUnexpected, err)
		}
	}

	return reader, nil
}
