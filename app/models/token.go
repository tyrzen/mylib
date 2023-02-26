package models

import "time"

type Token struct {
	ID     string
	UID    string
	Expiry time.Duration
}

type AccessToken struct {
	ReaderID       string
	RefreshTokenID string
	Role           string
	Expiry         time.Duration
}

type RefreshToken struct {
	ID       string
	ReaderID string
	Expiry   time.Duration
}

type TokenPair struct {
	AccessToken  string        `json:"access_token"`
	TokenType    string        `json:"token_type"`
	ExpiresIn    time.Duration `json:"expires_in"`
	RefreshToken string        `json:"refresh_token"`
}
