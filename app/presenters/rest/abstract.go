package rest

import (
	"context"

	"github.com/delveper/mylib/app/models"
)

type ReaderLogic interface {
	Auth(context.Context, models.AccessToken) error
	Refresh(context.Context, models.RefreshToken) (*models.TokenPair, error)
	SignUp(context.Context, models.Reader) error
	SignIn(context.Context, models.Credentials) (*models.TokenPair, error)
	SignOut(context.Context, models.AccessToken) error
}

type BookLogic interface {
	Import(context.Context, models.Book) error
	Fetch(context.Context, models.Book) (models.Book, error)
	FetchMany(context.Context, models.DataFilter) ([]models.Book, error)
	ExportToCSV(context.Context, models.DataFilter) ([]byte, error)
	AddToFavorites(context.Context, models.Reader, models.Book) error
	AddToWishlist(context.Context, models.Reader, models.Book) error
}
