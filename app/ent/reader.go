package ent

import (
	"fmt"
	"strings"
	"time"

	"github.com/delveper/mylib/app/exc"
	"github.com/delveper/mylib/lib/hash"
	"github.com/delveper/mylib/lib/revalid"
)

type Reader struct {
	ID        string    `json:"id"`                                                            // revalid:"(?i)^[0-9a-f]{8}\b-[0-9a-f]{4}\b-[0-9a-f]{4}\b-[0-9a-f]{4}\b-[0-9a-f]{12}$"
	FirstName string    `json:"first_name" revalid:"^[\p{L}&\s-\\'’.]{2,256}$"`                //nolint:govet
	LastName  string    `json:"last_name" revalid:"^[\p{L}&\s-\\'’.]{2,256}$"`                 //nolint:govet
	Email     string    `json:"email" revalid:"(?i)(^[a-z0-9_.+-]+@[a-z0-9-]+\.[a-z0-9-.]+$)"` //nolint:govet
	Password  string    `json:"password" revalid:"^[[:graph:]]{8,256}$"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
}

func (r *Reader) OK() error {
	if err := revalid.ValidateStruct(r); err != nil {
		return fmt.Errorf("%w: %v", exc.ErrValidation, err)
	}

	return nil
}

func (r *Reader) Normalize() {
	r.FirstName = strings.TrimSpace(r.FirstName)
	r.LastName = strings.TrimSpace(r.LastName)
	r.Email = strings.TrimSpace(strings.ToLower(r.Email))
}

func (r *Reader) HashPassword() (err error) {
	if r.Password, err = hash.Make(r.Password); err != nil {
		return fmt.Errorf("%w: %v", exc.ErrHashing, err)
	}

	return nil
}

func (r *Reader) CheckPassword(password string) error {
	if err := hash.Check(password, r.Password); err != nil {
		return fmt.Errorf("%w: %v", exc.ErrComparingHash, err)
	}

	return nil
}
