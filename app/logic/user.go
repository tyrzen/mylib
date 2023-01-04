package logic

import (
	"context"
	"errors"
	"fmt"

	"github.com/delveper/mylib/app/ent"
	"github.com/delveper/mylib/app/exc"
)

type Reader struct {
	repo ReaderRepository
}

func NewReader(repo ReaderRepository) Reader {
	return Reader{repo: repo}
}

func (r Reader) SignUp(ctx context.Context, reader ent.Reader) error {
	switch err := r.repo.Create(ctx, reader); {
	case errors.Is(err, exc.ErrDuplicateEmail):
		return fmt.Errorf("reader with given email is already exists: %w", err)
	case err != nil:
		return fmt.Errorf("error signup reader: %w", err)
	default:
		return nil
	}
}
