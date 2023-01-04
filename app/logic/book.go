package logic

import "github.com/delveper/mylib/app/ent"

type Book struct {
	repo   BookRepository
	logger ent.Logger
}

func NewBook(repo BookRepository, logger ent.Logger) Book {
	return Book{repo: repo, logger: logger}
}
