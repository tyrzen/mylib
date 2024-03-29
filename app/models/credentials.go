package models

import (
	"fmt"
	"strings"

	"github.com/delveper/mylib/app/exceptions"
	"github.com/delveper/revalid"
)

type Credentials struct {
	Email    string `json:"email" regex:"(?i)(^[a-z0-9_.+-]+@[a-z0-9-]+\.[a-z0-9-.]+$)"`
	Password string `json:"password" regex:"^[[:graph:]]{8,256}$"`
}

func (c *Credentials) Normalize() {
	c.Email = strings.ToLower(strings.TrimSpace(c.Email))
}

func (c *Credentials) OK() error {
	if err := revalid.ValidateStruct(c); err != nil {
		return fmt.Errorf("%w: %w", exceptions.ErrValidation, err)
	}

	return nil
}
