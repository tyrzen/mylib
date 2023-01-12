package tokay

import (
	"crypto"
	"fmt"
	"os"
	"time"

	"github.com/delveper/mylib/app/exc"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
)

func New(aud string) (string, error) {
	exp, err := time.ParseDuration(os.Getenv("JWT_EXP"))
	if err != nil {
		return "", fmt.Errorf("error parsing token expirity duration: %w", err)
	}

	now := time.Now()

	iss := os.Getenv("APP_NAME")

	claims := jwt.StandardClaims{
		Audience:  aud,
		ExpiresAt: now.Add(exp).Unix(),
		IssuedAt:  now.Unix(),
		Issuer:    iss,
	}

	alg, err := parseAlg(os.Getenv("JWT_ALG"))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(alg, claims)

	key := os.Getenv("JWT_KEY")

	str, err := token.SignedString([]byte(key))
	if err != nil {
		return "", fmt.Errorf("error creating token string: %v", err)
	}

	return str, nil
}

func parseAlg(alg string) (*jwt.SigningMethodHMAC, error) {
	switch alg {
	case "HS256":
		return &jwt.SigningMethodHMAC{Name: alg, Hash: crypto.SHA256}, nil
	case "HS384":
		return &jwt.SigningMethodHMAC{Name: alg, Hash: crypto.SHA512}, nil
	case "HS512":
		return &jwt.SigningMethodHMAC{Name: alg, Hash: crypto.SHA512}, nil
	default:
		return nil, fmt.Errorf("unsupported alg for HMAC: %s", alg)
	}
}

func Parse(str string) error {
	token, err := jwt.Parse(str, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_KEY")), nil
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
