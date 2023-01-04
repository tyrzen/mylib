package psql

import (
	"context"
	"database/sql"
	"time"

	"github.com/delveper/mylib/app/ent"
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

	if err != nil {
		return err
	}

	return nil
}
