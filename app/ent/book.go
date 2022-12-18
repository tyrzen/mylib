package ent

import "time"

type Book struct {
	ID       string    `json:"id"`
	AuthorID string    `json:"author_id"`
	Title    string    `json:"title"`
	Genre    string    `json:"genre"`
	Rate     int       `json:"rate"`
	WroteAt  time.Time `json:"wrote_at"`
}
