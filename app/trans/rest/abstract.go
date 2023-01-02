package rest

import "github.com/delveper/mylib/app/ent"

type ReaderLogic interface {
	SignUp(ent.Reader) error
	SignOut(ent.Reader) error
	SignIn(ent.Reader) error
}

type BookLogic interface {
	Add(ent.Book) error
	GetAll(ent.DataFilter) []ent.Book
	GetByID(string) ent.Book
	AddToFavorites(ent.Book) error
	AddToWishList(ent.Book) error
	ExportLibrary() error
}
