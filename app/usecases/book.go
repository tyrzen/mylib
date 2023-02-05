package usecases

import (
	"context"
	"fmt"

	"github.com/delveper/mylib/app/models"
)

type Book struct {
	repo   BookRepository
	author AuthorRepository
}

func NewBook(repo BookRepository, author AuthorRepository) Book {
	return Book{
		repo:   repo,
		author: author,
	}
}

func (b Book) Import(ctx context.Context, book models.Book) error {
	author := models.Author{ID: book.AuthorID}
	if _, err := b.author.GetByID(ctx, author); err != nil {
		return fmt.Errorf("error getting author record: %w", err)
	}

	if err := b.repo.Add(ctx, book); err != nil {
		return fmt.Errorf("error adding book record: %w", err)
	}

	return nil
}

func (b Book) Fetch(ctx context.Context, book models.Book) (models.Book, error) {
	book, err := b.repo.GetByID(ctx, book)
	if err != nil {
		return models.Book{}, fmt.Errorf("error fetching book record: %w", err)
	}

	return book, nil
}

func (b Book) FetchMany(ctx context.Context, filter models.DataFilter) ([]models.Book, error) {
	books, err := b.repo.GetMany(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error fetching book records: %w", err)
	}

	return books, nil
}
