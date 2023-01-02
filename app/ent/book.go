package ent

import "time"

type Book struct {
	ID        string    `json:"id" reval:"(?i)^[0-9a-f]{8}\b-[0-9a-f]{4}\b-[0-9a-f]{4}\b-[0-9a-f]{4}\b-[0-9a-f]{12}$"`
	AuthorID  string    `json:"author_id" reval:"(?i)^[0-9a-f]{8}\b-[0-9a-f]{4}\b-[0-9a-f]{4}\b-[0-9a-f]{4}\b-[0-9a-f]{12}$"`
	Title     string    `json:"title" reval:"^[[:graph:]]{1,256}$"`
	Genre     string    `json:"genre"`
	Rate      int       `json:"rate"`
	CreatedAt time.Time `json:"created_at"`
}
