package usecases

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/delveper/mylib/app/exceptions"
	"github.com/delveper/mylib/app/models"
	"github.com/delveper/mylib/lib/tokay"
	"github.com/google/uuid"
)

func (r Reader) newTokenPair(ctx context.Context, reader models.Reader) (*models.TokenPair, error) {
	refreshToken, refreshTokenVal, err := newRefreshToken(reader.ID)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, exceptions.ErrTokenNotCreated)
	}
	if err = r.sess.Create(ctx, refreshToken); err != nil {
		return nil, err
	}

	accessToken, accessTokenVal, err := newAccessToken(reader.ID, refreshToken.ID, reader.Role)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, exceptions.ErrTokenNotCreated)
	}

	if err = r.sess.Create(ctx, accessToken); err != nil {
		return nil, err
	}
	return &models.TokenPair{
		AccessToken:  accessTokenVal,
		TokenType:    "Bearer",
		ExpiresIn:    time.Duration(accessToken.Expiry.Seconds()),
		RefreshToken: refreshTokenVal,
	}, nil
}

func newAccessToken(readerID, refreshTokenID, role string) (token models.Token, val string, err error) {
	alg := os.Getenv("JWT_ALG")
	key := os.Getenv("JWT_KEY")

	exp, err := time.ParseDuration(os.Getenv("JWT_ACCESS_EXP"))

	if err != nil {
		return models.Token{}, "", fmt.Errorf("error parsing access token expirity: %w", err)
	}

	data := models.AccessToken{
		ReaderID:       readerID,
		RefreshTokenID: refreshTokenID,
		Role:           role,
		Expiry:         exp,
	}

	val, err = tokay.Make[models.AccessToken](alg, key, exp, data)
	if err != nil {
		return models.Token{}, "", fmt.Errorf("error making access token: %w", err)
	}

	token = models.Token{
		ID:     readerID,
		UID:    refreshTokenID,
		Expiry: exp,
	}

	return token, val, nil
}

func newRefreshToken(uid string) (token models.Token, val string, err error) {
	id := uuid.New().String()

	alg := os.Getenv("JWT_ALG")
	key := os.Getenv("JWT_KEY")

	exp, err := time.ParseDuration(os.Getenv("JWT_REFRESH_EXP"))
	if err != nil {
		return models.Token{}, "", fmt.Errorf("error parsing refresh token expirity: %w", err)
	}

	data := models.RefreshToken{
		ID:       id,
		ReaderID: uid,
		Expiry:   exp,
	}

	val, err = tokay.Make[models.RefreshToken](alg, key, exp, data)
	if err != nil {
		return models.Token{}, "", fmt.Errorf("error making refresh token: %w", err)
	}

	token = models.Token{
		ID:     id,
		UID:    uid,
		Expiry: exp,
	}

	return token, val, nil
}
