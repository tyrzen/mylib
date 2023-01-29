package models

import (
	"fmt"
	"time"

	"github.com/delveper/mylib/app/exceptions"
	"github.com/delveper/mylib/lib/revalid"
)

type Book struct {
	ID        string    `json:"id"`
	AuthorID  string    `json:"author_id" revalid:"(?i)^[0-9a-f]{8}\b-[0-9a-f]{4}\b-[0-9a-f]{4}\b-[0-9a-f]{4}\b-[0-9a-f]{12}$"`
	Title     string    `json:"title" revalid:"^[[:graph:]]{1,256}$"`
	Genre     string    `json:"genre" revalid:"^[[:graph:]]{1,256}$"`
	Rate      int       `json:"rate" revalid:"^([[:digit:]]|10)$"`
	Size      int       `json:"size" revalid:"^[[:digit:]]{1,256}$"`
	CreatedAt time.Time `json:"created_at"`
}

func (b *Book) OK() error {
	if err := revalid.ValidateStruct(b); err != nil {
		return fmt.Errorf("%w: %v", exceptions.ErrValidation, err)
	}

	return nil
}
