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

	const SQL = `
				INSERT INTO readers (id, first_name, last_name, email, password, created_at)
					VALUES(gen_random_uuid(), $1, $2, $3, $4, NOW())--RETURNING *;
`
	_, err := r.ExecContext(ctx, SQL,
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

// GetByID retrieves ent.Reader by given UserID.
func (r Reader) GetByID(ctx context.Context, reader ent.Reader) (ent.Reader, error) {
	const SQL = `
				SELECT id, first_name, last_name, email, password, created_at 
				FROM readers 
				WHERE id=$1;
`
	row := r.QueryRowContext(ctx, SQL, reader.ID)

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

// GetByEmail retrieves ent.Reader email.
func (r Reader) GetByEmail(ctx context.Context, reader ent.Reader) (ent.Reader, error) {
	const SQL = `
				SELECT id, first_name, last_name, email, password, created_at 
				FROM readers 
				WHERE email=$1;
`
	row := r.QueryRowContext(ctx, SQL, reader.Email)

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
