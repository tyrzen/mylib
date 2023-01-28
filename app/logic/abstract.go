package logic

import (
	"context"
	"time"

	"github.com/delveper/mylib/app/ent"
)

type ReaderRepository interface {
	Add(context.Context, ent.Reader) error
	GetByID(context.Context, ent.Reader) (ent.Reader, error)
	GetByEmail(context.Context, ent.Reader) (ent.Reader, error)
}

type TokenRepository interface {
	Create(context.Context, ent.Token) error
	Find(context.Context, ent.Token) (ent.Token, error)
	Destroy(context.Context, ent.Token) error
}

type Tokenizer interface {
	Make(alg, key string, exp time.Duration, data any) (string, error)
	Parse(val, key string) (data any, err error)
}

type BookRepository interface {
	Add(ent.Book) error
	GetAll(ent.DataFilter) []ent.Book
	GetByID(string) ent.Book
	AddToFavorites(ent.Book) error
	AddToWishList(ent.Book) error
	ExportLibrary() error
}
