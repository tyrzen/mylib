package rest

import (
	"context"

	"github.com/delveper/mylib/app/ent"
)

type ReaderLogic interface {
	SignUp(context.Context, ent.Reader) error
	SignIn(context.Context, ent.Credentials) (ent.Reader, error)
	SignOut(context.Context, ent.Reader) error
}

type SessionLogic interface {
	Create(context.Context, ent.Token) error
	Find(context.Context, string) (ent.Token, error)
}

type BookLogic interface {
	Add(ent.Book) error
	GetAll(ent.DataFilter) []ent.Book
	GetByID(string) ent.Book
	AddToFavorites(ent.Book) error
	AddToWishList(ent.Book) error
	ExportLibrary() error
}
