package ent

import "time"

type Reader struct {
	ID        string    `json:"id" valid:"(?i)^[0-9a-f]{8}\b-[0-9a-f]{4}\b-[0-9a-f]{4}\b-[0-9a-f]{4}\b-[0-9a-f]{12}$"`
	FirstName string    `json:"first_name" valid:"^[\p{L}&\s-\\'’.]{2,256}$"` //nolint:govet
	LastName  string    `json:"last_name" valid:"^[\p{L}&\s-\\'’.]{2,256}$"`  //nolint:govet
	Email     string    `json:"email" valid:"(?i)(^[a-z0-9_.+-]+@[a-z0-9-]+\.[a-z0-9-.]+$)"`
	Password  string    `json:"password" valid:"^[[:graph:]]{8,256}$"`
	CreatedAt time.Time `json:"created_at"`
}

func (r Reader) Validate() error {
	// TODO implement me
	panic("implement me")
}
