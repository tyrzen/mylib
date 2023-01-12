package jwt

import (
	"crypto"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func New(user string) (string, error) {
	exp, err := time.ParseDuration(os.Getenv("JWT_EXP"))
	if err != nil {
		return "", fmt.Errorf("error parsing token expirity duration: %w", err)
	}

	now := time.Now()

	issuer := os.Getenv("APP_NAME")

	claims := jwt.StandardClaims{
		Audience:  user,
		ExpiresAt: now.Add(exp).Unix(),
		IssuedAt:  now.Unix(),
		Issuer:    issuer,
	}

	alg, err := parseAlg(os.Getenv("JWT_ALG"))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(alg, claims)

	key := os.Getenv("JWT_KEY")

	str, err := token.SignedString(key)
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
