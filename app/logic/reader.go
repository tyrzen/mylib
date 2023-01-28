package logic

import (
	"context"
	"fmt"

	"github.com/delveper/mylib/app/ent"
	"github.com/delveper/mylib/app/exc"
	"github.com/delveper/mylib/lib/hash"
)

type Reader struct {
	repo ReaderRepository
	sess TokenRepository
	jwt  Tokenizer
}

func NewReader(repo ReaderRepository, sess TokenRepository, jwt Tokenizer) Reader {
	return Reader{
		repo: repo,
		sess: sess,
		jwt:  jwt,
	}
}

func (r Reader) Auth(ctx context.Context, token ent.Token) error {
	accessToken, err := r.retrieveToken(token)
	if err != nil {
		return fmt.Errorf("%v: %w", err, exc.ErrTokenInvalid)
	}

	savedToken, err := r.sess.Find(ctx, accessToken)
	if err != nil {
		return fmt.Errorf("%v: %w", err, exc.ErrTokenInvalid)
	}

	if accessToken.UID != savedToken.UID {
		return exc.ErrTokenInvalid
	}

	return nil
}

func (r Reader) SignUp(ctx context.Context, reader ent.Reader) error {
	if err := r.repo.Add(ctx, reader); err != nil {
		return fmt.Errorf("error signup reader: %w", err)
	}

	return nil
}

func (r Reader) SignIn(ctx context.Context, creds ent.Credentials) (*ent.TokenPair, error) {
	reader, err := r.repo.GetByEmail(ctx, ent.Reader{Email: creds.Email})
	if err != nil {
		return nil, fmt.Errorf("errror fetching reader: %w", err)
	}

	if err := hash.Verify(creds.Password, reader.Password); err != nil {
		return nil, exc.ErrInvalidCredits
	}

	tokenPair, err := r.newTokenPair(ctx, reader)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, exc.ErrTokenCreating)
	}

	return tokenPair, nil
}

func (r Reader) SignOut(ctx context.Context, token ent.Token) error {
	accessToken, err := r.retrieveToken(token)
	if err != nil {
		return fmt.Errorf("%v: %w", err, exc.ErrTokenInvalid)
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

func (r Reader) Refresh(ctx context.Context, token ent.Token) (*ent.TokenPair, error) {
	refreshToken, err := r.retrieveToken(token)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", exc.ErrTokenInvalid, err)
	}

	if err := r.sess.Destroy(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("error destroying refresh token: %w", err)
	}

	reader, err := r.repo.GetByID(ctx, ent.Reader{ID: refreshToken.ID})
	if err != nil {
		return nil, fmt.Errorf("errror fetching reader: %w", err)
	}

	tokenPair, err := r.newTokenPair(ctx, reader)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, exc.ErrTokenCreating)
	}

	return tokenPair, nil
}
