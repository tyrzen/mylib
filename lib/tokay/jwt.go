package tokay

import (
	"fmt"
	"strings"
	"time"

	"github.com/delveper/mylib/app/exceptions"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
)

type Claims[T any] struct {
	MetaData T `json:"data"`
	jwt.RegisteredClaims
}

func Parse[T any](val, key string) (data T, err error) {
	var claims Claims[T]

	token, err := jwt.ParseWithClaims(val, &claims, func(*jwt.Token) (interface{}, error) { return []byte(key), nil })

	if err != nil || !token.Valid {
		var errV *jwt.ValidationError
		if errors.As(err, &errV) {
			switch errV.Errors {
			case jwt.ValidationErrorExpired:
				return data, exceptions.ErrTokenExpired
			case jwt.ValidationErrorSignatureInvalid:
				return data, fmt.Errorf("%w: %v", exceptions.ErrTokenInvalidSigningMethod, err)
			default:
			}
		}

		return data, fmt.Errorf("%w: %v", exceptions.ErrTokenInvalid, err)
	}

	data = claims.MetaData

	return data, nil
}

func Make[T any](alg, key string, exp time.Duration, data T) (string, error) {
	method, err := selectMethod(alg)
	if err != nil {
		return "", fmt.Errorf("error parsing method: %w", err)
	}

	claims := Claims[T]{
		MetaData:         data,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp))},
	}

	token := jwt.NewWithClaims(method, claims)

	val, err := token.SignedString([]byte(key))
	if err != nil {
		return "", fmt.Errorf("error creating token string: %w", err)
	}

	return val, nil
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
