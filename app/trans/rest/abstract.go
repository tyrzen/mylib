package rest

import (
	"context"

	"github.com/delveper/mylib/app/ent"
)

type ReaderLogic interface {
	SignUp(context.Context, ent.Reader) error
	SignIn(context.Context, ent.Credentials) (*ent.TokenPair, error)
	SignOut(context.Context, ent.Token) error
	Refresh(context.Context, ent.Token) (*ent.TokenPair, error)
}

type BookLogic interface {
	Add(ent.Book) error
	GetAll(ent.DataFilter) []ent.Book
	GetByID(string) ent.Book
	AddToFavorites(ent.Book) error
	AddToWishList(ent.Book) error
	ExportLibrary() error
}
