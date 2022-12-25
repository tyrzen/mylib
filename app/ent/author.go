package ent

import "time"

type Author struct {
	ID        string    `json:"id" valid:"(?i)^[0-9a-f]{8}\b-[0-9a-f]{4}\b-[0-9a-f]{4}\b-[0-9a-f]{4}\b-[0-9a-f]{12}$"`
	FirstName string    `json:"first_name" valid:"^[\p{L}&\s-\\'’.]{2,256}$"` //nolint:govet
	LastName  string    `json:"last_name" valid:"^[\p{L}&\s-\\'’.]{2,256}$"`  //nolint:govet
	CreatedAt time.Time `json:"created_at"`
}
