package psql

import (
	"context"
	"database/sql"

	"github.com/delveper/mylib/app/models"
)

type Author struct{ *sql.DB }

func NewAuthor(db *sql.DB) *Author {
	return &Author{db}
}

func (a Author) GetByID(ctx context.Context, author models.Author) (models.Author, error) {
	return models.Author{}, nil
}
