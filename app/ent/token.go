package ent

import "time"

type Token struct {
	ID         string
	UID        string
	Expiration time.Duration
}
