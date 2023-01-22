package ent

import "time"

type Token struct {
	ID         string
	UID        string
	Expiration time.Duration
}

func NewToken(id, uid string, exp time.Duration) Token {
	return Token{
		ID:         id,
		UID:        uid,
		Expiration: exp,
	}
}
