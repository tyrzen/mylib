package logic

import (
	"context"
	"fmt"
	"os"

	"github.com/delveper/mylib/app/ent"
	"github.com/delveper/mylib/app/exc"
	"github.com/delveper/mylib/lib/hash"
	"github.com/delveper/mylib/lib/jwt"
)

type Reader struct {
	repo ReaderRepository
	sess TokenRepository
}

func NewReader(repo ReaderRepository, sess TokenRepository) Reader {
	return Reader{
		repo: repo,
		sess: sess,
	}
}

func (r Reader) SignUp(ctx context.Context, reader ent.Reader) error {
	if err := r.repo.Add(ctx, reader); err != nil {
		return err
	}

	return nil
}

func (r Reader) SignIn(ctx context.Context, creds ent.Credentials) (*ent.TokenPair, error) {
	reader, err := r.repo.GetByEmailOrID(ctx, ent.Reader{Email: creds.Email})
	if err != nil {
		return nil, fmt.Errorf("errror fetching reader: %w", err)
	}

	if err := hash.Verify(creds.Password, reader.Password); err != nil {
		return nil, exc.ErrInvalidCredits
	}

	refreshToken, err := NewRefreshToken(reader.ID)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, exc.ErrTokenCreating)
	}

	if err = r.sess.Create(ctx, *refreshToken); err != nil {
		return nil, fmt.Errorf("%v: %w", err, exc.ErrTokenCreating)
	}

	accessToken, err := NewAccessToken(refreshToken.ID, reader.Role)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, exc.ErrTokenCreating)
	}

	if err = r.sess.Create(ctx, *accessToken); err != nil {
		return nil, fmt.Errorf("%v: %w", err, exc.ErrTokenCreating)
	}

	return &ent.TokenPair{
		Access:  *accessToken,
		Refresh: *refreshToken,
	}, nil
}

func (r Reader) SignOut(ctx context.Context, token ent.Token) error {
	key := os.Getenv("JWT_KEY")

	payload, err := jwt.ParseJWT(token.Value, key)
	if err != nil {
		return fmt.Errorf("error parsing token: %w", err)
	}

	accessToken, ok := payload.(ent.Token)
	if !ok {
		return exc.ErrUnexpected
	}

	if err := r.sess.Destroy(ctx, accessToken); err != nil {
		return fmt.Errorf("error destroying access token: %w", err)
	}

	refreshToken := ent.Token{ID: accessToken.UID}
	if err := r.sess.Destroy(ctx, refreshToken); err != nil {
		return fmt.Errorf("error destroying refresh token: %w", err)
	}

	return nil
}
