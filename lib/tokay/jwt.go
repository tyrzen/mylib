package tokay

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/delveper/mylib/app/exc"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"github.com/rs/xid"
)

type TokenPair struct {
	Access  Token
	Refresh Token
}

type Token struct {
	ID         string
	Value      string
	Expiration time.Duration
}

type Claims struct {
	ID    string
	UID   string
	Admin bool
	jwt.StandardClaims
}

func NewTokenPair(uid string, admin bool) (TokenPair, error) {
	alg := os.Getenv("JWT_ALG")
	key := os.Getenv("JWT_KEY")

	// access
	expA, err := time.ParseDuration(os.Getenv("JWT_ACCESS_EXP"))
	if err != nil {
		return TokenPair{}, fmt.Errorf("error parsing access token expirity duration: %w", err)
	}

	access, err := makeToken(uid, alg, key, admin, expA)
	if err != nil {
		return TokenPair{}, fmt.Errorf("error making access token: %w", err)
	}

	// refresh
	expR, err := time.ParseDuration(os.Getenv("JWT_REFRESH_EXP"))
	if err != nil {
		return TokenPair{}, fmt.Errorf("error parsing refresh token expirity duration: %w", err)
	}

	refresh, err := makeToken(uid, alg, key, admin, expR)
	if err != nil {
		return TokenPair{}, fmt.Errorf("error making refresh token: %w", err)
	}

	// pair
	return TokenPair{
		Access:  access,
		Refresh: refresh,
	}, nil
}

func Check(val string) error {
	key := os.Getenv("JWT_KEY")

	token, err := jwt.Parse(val, func(*jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil || !token.Valid {
		var errV *jwt.ValidationError
		if errors.As(errV, err) {
			switch errV.Errors {
			case jwt.ValidationErrorExpired:
				return exc.ErrTokenExpired
			case jwt.ValidationErrorSignatureInvalid:
				return fmt.Errorf("%w: %v", exc.ErrTokenInvalidSigningMethod, err)
			default:
			}
		}

		return fmt.Errorf("%w: %v", exc.ErrTokenInvalid, err)
	}

	return nil
}

func makeToken(uid, alg, key string, admin bool, exp time.Duration) (Token, error) {
	id := xid.New().String()

	method, err := selectMethod(alg)
	if err != nil {
		return Token{}, err
	}

	expAt := time.Now().Add(exp).Unix()

	claims := Claims{
		ID:             id,
		UID:            uid,
		Admin:          admin,
		StandardClaims: jwt.StandardClaims{ExpiresAt: expAt},
	}

	token := jwt.NewWithClaims(method, claims)

	val, err := token.SignedString([]byte(key))
	if err != nil {
		return Token{}, fmt.Errorf("error creating token string: %w", err)
	}

	return Token{
		ID:         id,
		Value:      val,
		Expiration: exp,
	}, nil
}

func selectMethod(alg string) (jwt.SigningMethod, error) {
	switch strings.ToLower(alg) {
	case "eddsa":
		return jwt.SigningMethodEdDSA, nil
	case "hs256":
		return jwt.SigningMethodHS256, nil
	case "hs384":
		return jwt.SigningMethodHS384, nil
	case "hs512":
		return jwt.SigningMethodHS512, nil
	case "es256":
		return jwt.SigningMethodES256, nil
	case "es384":
		return jwt.SigningMethodES384, nil
	case "es512":
		return jwt.SigningMethodES512, nil
	case "rs256":
		return jwt.SigningMethodRS256, nil
	case "rs384":
		return jwt.SigningMethodRS384, nil
	case "rs512":
		return jwt.SigningMethodRS512, nil
	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", alg)
	}
}
