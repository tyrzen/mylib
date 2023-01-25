package ent

import "time"

type Token struct {
	ID     string
	UID    string
	Role   string
	Value  string
	Expiry time.Duration
}

type TokenPair struct {
	Access  Token
	Refresh Token
}

func NewToken(id, uid, role, val string, exp time.Duration) Token {
	return Token{
		ID:     id,
		UID:    uid,
		Role:   role,
		Value:  val,
		Expiry: exp,
	}
}
