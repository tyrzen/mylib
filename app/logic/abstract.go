package logic

import (
	"context"

	"github.com/delveper/mylib/app/ent"
)

type ReaderRepository interface {
	Create(context.Context, ent.Reader) error
	Authenticate(context.Context, ent.Reader) error
}

type BookRepository interface {
	Add(ent.Book) error
	GetAll(ent.DataFilter) []ent.Book
	GetByID(string) ent.Book
	AddToFavorites(ent.Book) error
	AddToWishList(ent.Book) error
	ExportLibrary() error
}
