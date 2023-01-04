package ent

import (
	"time"

	"github.com/delveper/mylib/lib/reval"
)

type Reader struct {
	ID        string    `json:"id" reval:"(?i)^[0-9a-f]{8}\b-[0-9a-f]{4}\b-[0-9a-f]{4}\b-[0-9a-f]{4}\b-[0-9a-f]{12}$"`
	FirstName string    `json:"first_name" reval:"^[\p{L}&\s-\\'’.]{2,256}$"`
	LastName  string    `json:"last_name" reval:"^[\p{L}&\s-\\'’.]{2,256}$"`
	Email     string    `json:"email" reval:"(?i)(^[a-z0-9_.+-]+@[a-z0-9-]+\.[a-z0-9-.]+$)"`
	Password  string    `json:"password" reval:"^[[:graph:]]{8,256}$"`
	CreatedAt time.Time `json:"created_at"`
}

func (r Reader) Validate() error {
	return reval.ValidateStruct(r)
}
