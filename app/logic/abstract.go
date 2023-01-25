package logic

import (
	"context"

	"github.com/delveper/mylib/app/ent"
)

type ReaderRepository interface {
	Add(context.Context, ent.Reader) error
	GetByEmailOrID(context.Context, ent.Reader) (ent.Reader, error)
}

type TokenRepository interface {
	Create(context.Context, ent.Token) error
	Find(context.Context, ent.Token) (*ent.Token, error)
	Destroy(context.Context, ent.Token) error
}

type BookRepository interface {
	Add(ent.Book) error
	GetAll(ent.DataFilter) []ent.Book
	GetByID(string) ent.Book
	AddToFavorites(ent.Book) error
	AddToWishList(ent.Book) error
	ExportLibrary() error
}
