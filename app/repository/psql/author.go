package psql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/delveper/mylib/app/exceptions"
	"github.com/delveper/mylib/app/models"
	"github.com/pkg/errors"
)

type Author struct{ *sql.DB }

func NewAuthor(db *sql.DB) *Author {
	return &Author{db}
}

func (a Author) GetByID(ctx context.Context, author models.Author) (models.Author, error) {
	const SQL = `SELECT id, first_name, last_name, created_at 
				 FROM authors
				 WHERE id=$1;`

	row := a.QueryRowContext(ctx, SQL, author.ID)

	err := row.Scan(
		&author.ID,
		&author.FirstName,
		&author.LastName,
		&author.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return models.Author{}, fmt.Errorf("%w: %v", exceptions.ErrDeadline, err)
		case errors.Is(err, sql.ErrNoRows):
			return models.Author{}, fmt.Errorf("%w: %v", exceptions.ErrRecordNotFound, err)
		default:
			return models.Author{}, fmt.Errorf("%w: %v", exceptions.ErrUnexpected, err)
		}
	}

	return author, nil
}
