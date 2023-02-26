package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/delveper/mylib/app/exceptions"
	"github.com/delveper/mylib/lib/hash"
	"github.com/delveper/mylib/lib/revalid"
)

type Reader struct {
	ID        string    `json:"id"` // regex:"(?i)^[0-9a-f]{8}\b-[0-9a-f]{4}\b-[0-9a-f]{4}\b-[0-9a-f]{4}\b-[0-9a-f]{12}$"
	FirstName string    `json:"first_name" regex:"^[\p{L}&\s-\\'’.]{2,256}$"`
	LastName  string    `json:"last_name" regex:"^[\p{L}&\s-\\'’.]{2,256}$"`
	Email     string    `json:"email" regex:"(?i)(^[a-z0-9_.+-]+@[a-z0-9-]+\.[a-z0-9-.]+$)"`
	Password  string    `json:"password" regex:"^[[:graph:]]{8,256}$"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

func (r *Reader) OK() error {
	if err := revalid.ValidateStruct(r); err != nil {
		return fmt.Errorf("%w: %w", exceptions.ErrValidation, err)
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
		return fmt.Errorf("%w: %w", exceptions.ErrHashing, err)
	}

	return nil
}

func (r *Reader) CheckPassword(password string) error {
	if err := hash.Verify(password, r.Password); err != nil {
		return fmt.Errorf("%w: %w", exceptions.ErrComparingHash, err)
	}

	return nil
}
