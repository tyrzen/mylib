package logic

import (
	"context"
	"fmt"

	"github.com/delveper/mylib/app/ent"
	"github.com/delveper/mylib/app/exc"
	"github.com/delveper/mylib/lib/hash"
)

type Reader struct {
	ReaderRepository
}

func NewReader(repo ReaderRepository) Reader {
	return Reader{repo}
}

func (r Reader) SignUp(ctx context.Context, reader ent.Reader) error {
	if err := r.Add(ctx, reader); err != nil {
		return err
	}

	return nil
}

func (r Reader) SignIn(ctx context.Context, creds ent.Credentials) (ent.Token, error) {
	reader, err := r.GetByEmailOrID(ctx, ent.Reader{Email: creds.Email})

	if err != nil {
		return ent.Token{}, fmt.Errorf("singin error: %w", err)
	}

	if err := hash.Check(creds.Password, reader.Password); err != nil {
		return ent.Token{}, exc.ErrInvalidCredits
	}

	return ent.Token{}, nil
}

func (r Reader) SignOut(ctx context.Context, reader ent.Reader) error {
	return nil
}
