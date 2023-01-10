package hash

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func Make(password string) (string, error) {
	bts, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing: %w", err)
	}

	return string(bts), nil
}

func Compare(password string, hash string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}

	return true
}
