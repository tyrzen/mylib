package hash

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Hasher struct{}

func NewHash() Hasher { return Hasher{} }

func (h Hasher) Hash(password string) (string, error) {
	bts, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MaxCost)
	if err != nil {
		return "", fmt.Errorf("error hashing: %w", err)
	}

	return string(bts), nil
}

func (h Hasher) Compare(password string, hash string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(hash)); err != nil {
		return false
	}

	return true
}
