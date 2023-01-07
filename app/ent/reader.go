package ent

import (
	"time"

	"github.com/delveper/mylib/lib/revalid"
)

type Reader struct {
	ID        string    `json:"id"` // revalid:"(?i)^[0-9a-f]{8}\b-[0-9a-f]{4}\b-[0-9a-f]{4}\b-[0-9a-f]{4}\b-[0-9a-f]{12}$"
	FirstName string    `json:"first_name" revalid:"^[\p{L}&\s-\\'’.]{2,256}$"`
	LastName  string    `json:"last_name" revalid:"^[\p{L}&\s-\\'’.]{2,256}$"`
	Email     string    `json:"email" revalid:"(?i)(^[a-z0-9_.+-]+@[a-z0-9-]+\.[a-z0-9-.]+$)"`
	Password  string    `json:"password" revalid:"^[[:graph:]]{8,256}$"`
	CreatedAt time.Time `json:"created_at"`
}

func (r Reader) Validate() error {
	return revalid.ValidateStruct(r)
}
