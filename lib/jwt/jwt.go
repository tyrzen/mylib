package jwt

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/delveper/mylib/app/exc"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
)

type JWT struct{}

func New() JWT {
	return JWT{}
}

type Claims struct {
	MetaData any `json:"data"`
	jwt.StandardClaims
}

func (t JWT) Parse(val, key string) (data any, err error) {
	var claims Claims

	log.Println(val)

	token, err := jwt.ParseWithClaims(val, &claims, func(*jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil || !token.Valid {
		var errV *jwt.ValidationError
		if errors.As(errV, err) {
			switch errV.Errors {
			case jwt.ValidationErrorExpired:
				return nil, exc.ErrTokenExpired
			case jwt.ValidationErrorSignatureInvalid:
				return nil, fmt.Errorf("%w: %v", exc.ErrTokenInvalidSigningMethod, err)
			default:
			}
		}

		return nil, fmt.Errorf("%w: %v", exc.ErrTokenInvalid, err)
	}

	data = claims.MetaData

	return data, nil
}

func (t JWT) Make(alg, key string, exp time.Duration, data any) (string, error) {
	method, err := selectMethod(alg)
	if err != nil {
		return "", fmt.Errorf("error parsing method: %w", err)
	}

	claims := Claims{
		MetaData:       data,
		StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Add(exp).Unix()},
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
