package usecases

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"

	"github.com/delveper/mylib/app/exceptions"
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

func (b Book) ExportToCSV(ctx context.Context, filter models.DataFilter) ([]byte, error) {
	books, err := b.repo.GetMany(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error fetching book records: %w", err)
	}

	if len(books) == 0 {
		return nil, fmt.Errorf("no content to show: %w", exceptions.ErrNoContent)
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	header := []string{"ID", "AuthorID", "Title", "Genre", "Rate", "Size", "Year"}
	if err := writer.Write(header); err != nil {
		return nil, fmt.Errorf("error writing header to csv: %w", err)
	}

	for i, book := range books {
		row := []string{
			book.ID,
			book.AuthorID,
			book.Title,
			book.Genre,
			fmt.Sprintf("%d", book.Rate),
			fmt.Sprintf("%d", book.Size),
			fmt.Sprintf("%d", book.Year),
		}
		if err := writer.Write(row); err != nil {
			return nil, fmt.Errorf("error writing %d row to csv: %w", i+1, err)
		}
	}

	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("error flushing csv writer: %w", err)
	}

	return buf.Bytes(), nil
}

func (b Book) AddToFavorites(ctx context.Context, reader models.Reader, book models.Book) error {
	if err := b.repo.AddToFavorites(ctx, reader, book); err != nil {
		return fmt.Errorf("error adding book to favorites: %w", err)
	}

	return nil
}

func (b Book) AddToWishlist(ctx context.Context, reader models.Reader, book models.Book) error {
	if err := b.repo.AddToWishlist(ctx, reader, book); err != nil {
		return fmt.Errorf("error adding book to wishlist: %w", err)
	}

	return nil
}
