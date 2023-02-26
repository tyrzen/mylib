package models

import (
	"fmt"

	"github.com/delveper/mylib/app/exceptions"
	"github.com/delveper/revalid"
)

type Book struct {
	ID       string `json:"id" sql:"id"`
	AuthorID string `json:"author_id" sql:"author_id" regex:"(?i)^[0-9a-f]{8}\b-[0-9a-f]{4}\b-[0-9a-f]{4}\b-[0-9a-f]{4}\b-[0-9a-f]{12}$"`
	Title    string `json:"title" sql:"title" regex:"^[[:graph:]]{1,256}$"`
	// TODO: Add ISBN
	Genre string `json:"genre" sql:"genre" regex:"^[[:graph:]]{1,256}$"`
	Rate  int    `json:"rate" sql:"rate" regex:"^([[:digit:]]|10)$"`
	Size  int    `json:"size" sql:"size" regex:"^[[:digit:]]{1,256}$"`
	Year  int    `json:"year" sql:"year" regex:"^[[:digit:]]{4}$"`
}

func (b *Book) OK() error {
	if err := revalid.ValidateStruct(b); err != nil {
		return fmt.Errorf("%w: %w", exceptions.ErrValidation, err)
	}

	return nil
}
