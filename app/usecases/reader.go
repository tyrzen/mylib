package usecases

import (
	"context"
	"fmt"

	"github.com/delveper/mylib/app/exceptions"
	"github.com/delveper/mylib/app/models"
	"github.com/delveper/mylib/lib/hash"
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

func (r Reader) Auth(ctx context.Context, token models.AccessToken) error {
	accessToken := models.Token{ID: token.ID, UID: token.RefreshTokenID}

	savedToken, err := r.sess.Find(ctx, accessToken)
	if err != nil {
		return fmt.Errorf("%v: %w", err, exceptions.ErrTokenInvalid)
	}

	if accessToken.UID != savedToken.UID {
		return exceptions.ErrTokenInvalid
	}

	return nil
}

func (r Reader) SignUp(ctx context.Context, reader models.Reader) error {
	if err := r.repo.Add(ctx, reader); err != nil {
		return fmt.Errorf("error signup reader: %w", err)
	}

	return nil
}

func (r Reader) SignIn(ctx context.Context, creds models.Credentials) (*models.TokenPair, error) {
	reader, err := r.repo.GetByEmail(ctx, models.Reader{Email: creds.Email})
	if err != nil {
		return nil, fmt.Errorf("errror fetching reader: %w", err)
	}

	if err := hash.Verify(creds.Password, reader.Password); err != nil {
		return nil, exceptions.ErrInvalidCredits
	}

	tokenPair, err := r.newTokenPair(ctx, reader)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, exceptions.ErrTokenNotCreated)
	}

	return tokenPair, nil
}

func (r Reader) SignOut(ctx context.Context, token models.AccessToken) error {
	accessToken := models.Token{ID: token.ID, UID: token.RefreshTokenID}

	if err := r.sess.Destroy(ctx, accessToken); err != nil {
		return fmt.Errorf("error destroying access token: %w", err)
	}

	refreshToken := models.Token{ID: accessToken.UID}

	if err := r.sess.Destroy(ctx, refreshToken); err != nil {
		return fmt.Errorf("error destroying refresh token: %w", err)
	}

	return nil
}

func (r Reader) Refresh(ctx context.Context, token models.RefreshToken) (*models.TokenPair, error) {
	refreshToken := models.Token{ID: token.ID}

	if err := r.sess.Destroy(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("error destroying refresh token: %w", err)
	}

	reader, err := r.repo.GetByID(ctx, models.Reader{ID: token.ReaderID})
	if err != nil {
		return nil, fmt.Errorf("errror fetching reader: %w", err)
	}

	tokenPair, err := r.newTokenPair(ctx, reader)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, exceptions.ErrTokenNotCreated)
	}

	return tokenPair, nil
}
