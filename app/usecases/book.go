package usecases

import "github.com/delveper/mylib/app/models"

type Book struct {
	repo   BookRepository
	logger models.Logger
}

func NewBook(repo BookRepository, logger models.Logger) Book {
	return Book{
		repo:   repo,
		logger: logger,
	}
}
