package models

import (
	"time"
)

type Document struct {
	MD       MetaData
	Location string
}

type MetaData struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	MimeType  string    `json:"mime"`
	File      bool      `json:"file"`
	Public    bool      `json:"public"`
	CreatedAt time.Time `json:"created"`
	Grant     []string  `json:"grant"`
	Token     string    `json:"token"`
}
