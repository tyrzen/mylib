package usecases

import (
	"context"

	"github.com/delveper/mylib/app/models"
)

type Book struct {
	repo BookRepository
}

func NewBook(repo BookRepository) Book {
	return Book{
		repo: repo,
	}
}

func (b Book) Import(ctx context.Context, book models.Book) error {
	return nil
}
