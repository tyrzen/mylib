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
	Add(models.Book) error
	GetAll(models.DataFilter) []models.Book
	GetByID(string) models.Book
	AddToFavorites(models.Book) error
	AddToWishList(models.Book) error
	ExportLibrary() error
}
