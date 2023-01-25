package logic

import (
	"fmt"
	"os"
	"time"

	"github.com/delveper/mylib/app/ent"
	"github.com/delveper/mylib/lib/jwt"
	"github.com/rs/xid"
)

type AccessToken struct {
	ID             string
	Role           string
	RefreshTokenID string
	Expiry         time.Duration
}

type RefreshToken struct {
	ID     string
	UID    string
	Expiry time.Duration
}

func NewTokenPair(access AccessToken, refresh RefreshToken) ent.TokenPair {
	tokenPair := ent.TokenPair{
		Access: ent.Token{
			Value:  "",
			Expiry: access.Expiry,
		},
		Refresh: ent.Token{
			Value:  "",
			Expiry: refresh.Expiry,
		},
	}
	return tokenPair

}

func NewAccessToken(uid, role string) (*ent.Token, error) {
	id := xid.New().String()

	payload := ent.Token{ID: id, UID: uid, Role: role}

	alg := os.Getenv("JWT_ALG")
	key := os.Getenv("JWT_ALG")

	exp, err := time.ParseDuration(os.Getenv("JWT_ACCESS_EXP"))
	if err != nil {
		return nil, fmt.Errorf("error parsing access token expirity: %w", err)
	}

	val, err := jwt.MakeJWT(alg, key, exp, payload)
	if err != nil {
		return nil, fmt.Errorf("error making access token: %w", err)
	}

	return &ent.Token{
		ID:     id,
		UID:    uid,
		Role:   role,
		Value:  val,
		Expiry: exp,
	}, nil
}

func NewRefreshToken(uid string) (*ent.Token, error) {
	id := xid.New().String()

	payload := ent.Token{ID: id, UID: uid}

	alg := os.Getenv("JWT_ALG")
	key := os.Getenv("JWT_ALG")

	exp, err := time.ParseDuration(os.Getenv("JWT_REFRESH_EXP"))
	if err != nil {
		return nil, fmt.Errorf("error parsing refresh token expirity: %w", err)
	}

	val, err := jwt.MakeJWT(alg, key, exp, payload)
	if err != nil {
		return nil, fmt.Errorf("error making refresh token: %w", err)
	}

	return &ent.Token{
		ID:     id,
		UID:    uid,
		Value:  val,
		Expiry: exp,
	}, nil
}
