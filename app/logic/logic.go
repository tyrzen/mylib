package logic

import "github.com/delveper/mylib/app/ent"

type ReaderService struct {
	repo   ReaderRepository
	logger ent.Logger
}

type BookService struct {
	repo   BookRepository
	logger ent.Logger
}
