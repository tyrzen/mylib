package models

import "time"

type Author struct {
	ID        string    `json:"id" regex:"(?i)^[0-9a-f]{8}\b-[0-9a-f]{4}\b-[0-9a-f]{4}\b-[0-9a-f]{4}\b-[0-9a-f]{12}$"`
	FirstName string    `json:"first_name" regex:"^[\p{L}&\s-\\'’.]{2,256}$"`
	LastName  string    `json:"last_name" regex:"^[\p{L}&\s-\\'’.]{2,256}$"`
	CreatedAt time.Time `json:"created_at"`
}
