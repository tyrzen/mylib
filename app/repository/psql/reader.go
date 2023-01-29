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

type Reader struct{ *sql.DB }

func NewReader(db *sql.DB) *Reader {
	return &Reader{db}
}

// Add adds models.Reader entity.
func (r Reader) Add(ctx context.Context, reader models.Reader) error {
	const SQL = `INSERT INTO readers (id, first_name, last_name, email, password, role, created_at)
					VALUES(GEN_RANDOM_UUID(), $1, $2, $3, $4, NOW());`

	_, err := r.ExecContext(ctx, SQL,
		reader.FirstName, // $1
		reader.LastName,  // $2
		reader.Email,     // $3
		reader.Role,      // $4
		reader.Password,  // $5
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("%w: %v", exceptions.ErrDeadline, err)
		}

		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) {
			switch pgxErr.ConstraintName {
			case "readers_email_key":
				return fmt.Errorf("%w: %v", exceptions.ErrDuplicateEmail, err)
			case "readers_pkey":
				return fmt.Errorf("%w: %v", exceptions.ErrDuplicateID, err)
			}
		}

		return fmt.Errorf("%w: %v", exceptions.ErrUnexpected, err)
	}

	return nil
}

// GetByID retrieves models.Reader by given UserID.
func (r Reader) GetByID(ctx context.Context, reader models.Reader) (models.Reader, error) {
	const SQL = `SELECT id, first_name, last_name, email, password, role, created_at 
				 FROM readers 
				 WHERE id=$1;`

	row := r.QueryRowContext(ctx, SQL, reader.ID)

	err := row.Scan(
		&reader.ID,
		&reader.FirstName,
		&reader.LastName,
		&reader.Email,
		&reader.Password,
		&reader.Role,
		&reader.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return models.Reader{}, fmt.Errorf("%w: %v", exceptions.ErrDeadline, err)
		case errors.Is(err, sql.ErrNoRows):
			return models.Reader{}, fmt.Errorf("%w: %v", exceptions.ErrNoRecord, err)
		default:
			return models.Reader{}, fmt.Errorf("%w: %v", exceptions.ErrUnexpected, err)
		}
	}

	return reader, nil
}

// GetByEmail retrieves models.Reader email.
func (r Reader) GetByEmail(ctx context.Context, reader models.Reader) (models.Reader, error) {
	const SQL = `SELECT id, first_name, last_name, email, password, role, created_at 
				 FROM readers 
				 WHERE email=$1;`

	row := r.QueryRowContext(ctx, SQL, reader.Email)

	err := row.Scan(
		&reader.ID,
		&reader.FirstName,
		&reader.LastName,
		&reader.Email,
		&reader.Password,
		&reader.Role,
		&reader.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return models.Reader{}, fmt.Errorf("%w: %v", exceptions.ErrDeadline, err)
		case errors.Is(err, sql.ErrNoRows):
			return models.Reader{}, fmt.Errorf("%w: %v", exceptions.ErrNoRecord, err)
		default:
			return models.Reader{}, fmt.Errorf("%w: %v", exceptions.ErrUnexpected, err)
		}
	}

	return reader, nil
}
