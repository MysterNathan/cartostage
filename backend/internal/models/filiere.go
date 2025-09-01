package models

import "time"

type Filiere struct {
	ID        int       `json:"id" db:"id"`
	Code      string    `json:"code" db:"code"`
	Label     string    `json:"label" db:"label"`
	Color     string    `json:"color" db:"color"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type FilieresData struct {
	Filieres []Filiere `json:"filieres"`
}
