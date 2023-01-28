package logic

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/delveper/mylib/app/ent"
	"github.com/delveper/mylib/app/exc"
	"github.com/rs/xid"
)

func (r Reader) retrieveToken(token ent.Token) (ent.Token, error) {
	key := os.Getenv("JWT_KEY")

	data, err := r.jwt.Parse(token.Value, key)
	if err != nil {
		return ent.Token{}, fmt.Errorf("error parsing token: %w", err)
	}

	token, ok := data.(ent.Token)
	if !ok {
		return ent.Token{}, fmt.Errorf("error parsing token metadata: %+v", token)
	}

	return token, nil
}

func (r Reader) newTokenPair(ctx context.Context, reader ent.Reader) (*ent.TokenPair, error) {
	accessToken, err := r.newAccessToken(reader.ID, reader.Role)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, exc.ErrTokenCreating)
	}

	if err = r.sess.Create(ctx, accessToken); err != nil {
		return nil, err
	}

	refreshToken, err := r.newRefreshToken(accessToken.ID)
	if err != nil {
		return nil, err
	}

	if err = r.sess.Create(ctx, refreshToken); err != nil {
		return nil, err
	}

	return &ent.TokenPair{
		Access:  accessToken,
		Refresh: refreshToken,
	}, nil
}

func (r Reader) newAccessToken(uid, role string) (ent.Token, error) {
	id := xid.New().String()

	alg := os.Getenv("JWT_ALG")
	key := os.Getenv("JWT_KEY")

	exp, err := time.ParseDuration(os.Getenv("JWT_ACCESS_EXP"))
	if err != nil {
		return ent.Token{}, fmt.Errorf("error parsing access token expirity: %w", err)
	}

	payload := ent.Token{ID: id, UID: uid, Role: role, Expiry: exp}

	val, err := r.jwt.Make(alg, key, exp, payload)
	if err != nil {
		return ent.Token{}, fmt.Errorf("error making access token: %w", err)
	}

	return ent.Token{
		ID:     id,
		UID:    uid,
		Role:   role,
		Value:  val,
		Expiry: exp,
	}, nil
}

func (r Reader) newRefreshToken(uid string) (ent.Token, error) {
	id := xid.New().String()

	alg := os.Getenv("JWT_ALG")
	key := os.Getenv("JWT_KEY")

	exp, err := time.ParseDuration(os.Getenv("JWT_REFRESH_EXP"))
	if err != nil {
		return ent.Token{}, fmt.Errorf("error parsing refresh token expirity: %w", err)
	}

	payload := ent.Token{ID: id, UID: uid, Expiry: exp}

	val, err := r.jwt.Make(alg, key, exp, payload)
	if err != nil {
		return ent.Token{}, fmt.Errorf("error making refresh token: %w", err)
	}

	return ent.Token{
		ID:     id,
		UID:    uid,
		Value:  val,
		Expiry: exp,
	}, nil
}
