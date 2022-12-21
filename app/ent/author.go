package ent

type Author struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name" valid:"^[\p{L}&\s-\\'’.]{2,256}$"` //nolint:govet
	LastName  string `json:"last_name" valid:"^[\p{L}&\s-\\'’.]{2,256}$"`  //nolint:govet
	BornAt    string `json:"born_at"`
}
