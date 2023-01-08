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

func (repo Reader) Create(ctx context.Context, reader ent.Reader) error {
	reader.CreatedAt = time.Now()

	_, err := repo.DB.ExecContext(ctx, queryCreateReader,
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
