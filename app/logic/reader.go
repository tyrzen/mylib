package logic

import (
	"context"

	"github.com/delveper/mylib/app/ent"
)

type Reader struct {
	repo ReaderRepository
}

func NewReader(repo ReaderRepository) Reader {
	return Reader{repo: repo}
}

func (r Reader) SignUp(ctx context.Context, reader ent.Reader) error {
	if err := r.repo.Create(ctx, reader); err != nil {
		return err
	}
	return nil
}

func (r Reader) SingIn(ctx context.Context, reader ent.Reader) (ent.Reader, error) {

	return ent.Reader{}, nil
}
