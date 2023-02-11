package models

import (
	"fmt"

	"github.com/delveper/mylib/app/exceptions"
	"github.com/delveper/mylib/lib/revalid"
)

type Book struct {
	ID       string `json:"id" sql:"id"`
	AuthorID string `json:"author_id" sql:"author_id" revalid:"(?i)^[0-9a-f]{8}\b-[0-9a-f]{4}\b-[0-9a-f]{4}\b-[0-9a-f]{4}\b-[0-9a-f]{12}$"`
	Title    string `json:"title" sql:"title" revalid:"^[[:graph:]]{1,256}$"`
	Genre    string `json:"genre" sql:"genre" revalid:"^[[:graph:]]{1,256}$"`
	Rate     int    `json:"rate" sql:"rate" revalid:"^([[:digit:]]|10)$"`
	Size     int    `json:"size" sql:"size" revalid:"^[[:digit:]]{1,256}$"`
	Year     int    `json:"year" sql:"year" revalid:"^[[:digit:]]{4}$"`
}

func (b *Book) OK() error {
	if err := revalid.ValidateStruct(b); err != nil {
		return fmt.Errorf("%w: %v", exceptions.ErrValidation, err)
	}

	return nil
}
