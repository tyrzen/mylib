package ent

type Author struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name" revalid:"^[\p{L}&\s-\\'’.]{2,256}$"`
	LastName  string `json:"last_name" revalid:"^[\p{L}&\s-\\'’.]{2,256}$"`
	BornAt    string `json:"born_at"`
}
