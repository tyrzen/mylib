package ent

import "time"

type Reader struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email" valid:"(^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$)"`
	Password  string    `json:"password" valid:"^[[:graph:]]{8,256}$"`
	BornAt    time.Time `json:"born_at"`
}
