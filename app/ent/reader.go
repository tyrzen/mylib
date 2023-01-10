package ent

import (
	"strings"
	"time"

	"github.com/delveper/mylib/lib/revalid"
)

type Reader struct {
	ID        string    `json:"id"`                                                            // revalid:"(?i)^[0-9a-f]{8}\b-[0-9a-f]{4}\b-[0-9a-f]{4}\b-[0-9a-f]{4}\b-[0-9a-f]{12}$"
	FirstName string    `json:"first_name" revalid:"^[\p{L}&\s-\\'’.]{2,256}$"`                //nolint:govet
	LastName  string    `json:"last_name" revalid:"^[\p{L}&\s-\\'’.]{2,256}$"`                 //nolint:govet
	Email     string    `json:"email" revalid:"(?i)(^[a-z0-9_.+-]+@[a-z0-9-]+\.[a-z0-9-.]+$)"` //nolint:govet
	Password  string    `json:"password" revalid:"^[[:graph:]]{8,256}$"`
	CreatedAt time.Time `json:"created_at"`
}

func (r *Reader) OK() error {
	return revalid.ValidateStruct(r) //nolint:wrapcheck
}

func (r *Reader) Normalize() {
	r.FirstName = strings.TrimSpace(r.FirstName)
	r.LastName = strings.TrimSpace(r.LastName)
	r.Email = strings.TrimSpace(strings.ToLower(r.Email))
}
