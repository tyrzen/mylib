package psql

import (
	"context"
	"database/sql"
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
		reader.ID,        // $1
		reader.FirstName, // $2
		reader.LastName,  // $3
		reader.Email,     // $4
		reader.Password,  // $5
	)

	var pgxErr *pgconn.PgError
	if errors.As(err, &pgxErr) {
		switch pgxErr.ConstraintName {
		case "readers_email_key":
			return errors.Wrapf(exc.ErrDuplicateEmail, "reader with given email is already exists: %v", err)
		case "readers_pkey":
			return errors.Wrapf(exc.ErrDuplicateID, "reader with given ID is already exists: %v", err)
		default:
			return errors.Wrapf(exc.ErrUnexpected, "error querying: %v", err)
		}
	}

	return nil
}
