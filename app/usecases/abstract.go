package usecases

import (
	"context"

	"github.com/delveper/mylib/app/models"
)

type ReaderRepository interface {
	Add(context.Context, models.Reader) error
	GetByID(context.Context, models.Reader) (models.Reader, error)
	GetByEmail(context.Context, models.Reader) (models.Reader, error)
}

type TokenRepository interface {
	Create(context.Context, models.Token) error
	Find(context.Context, models.Token) (models.Token, error)
	Destroy(context.Context, models.Token) error
}

type AuthorRepository interface {
	GetByID(context.Context, models.Author) (models.Author, error)
}

type BookRepository interface {
	Add(context.Context, models.Book) error
	GetByID(context.Context, models.Book) (models.Book, error)
	GetMany(context.Context, models.DataFilter) ([]models.Book, error)
}
