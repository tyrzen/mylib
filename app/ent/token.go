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
