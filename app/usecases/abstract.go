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

type BookRepository interface {
	Add(models.Book) error
	GetAll(models.DataFilter) []models.Book
	GetByID(string) models.Book
	AddToFavorites(models.Book) error
	AddToWishList(models.Book) error
	ExportLibrary() error
}
