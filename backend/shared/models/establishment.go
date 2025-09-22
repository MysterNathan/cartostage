package models

import (
	"time"
)

type Establishment struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Address     string    `json:"address" db:"address"`
	Lat         float64   `json:"lat" db:"lat"`
	Lng         float64   `json:"lng" db:"lng"`
	Sector      string    `json:"sector" db:"sector"`
	Size        *int      `json:"size,omitempty" db:"size"`
	Siret       *string   `json:"siret,omitempty" db:"siret"`
	Email       *string   `json:"email,omitempty" db:"email"`
	Phone       *string   `json:"phone,omitempty" db:"phone"`
	Website     *string   `json:"website,omitempty" db:"website"`
	Description *string   `json:"description,omitempty" db:"description"`
	LogoURL     *string   `json:"logo_url,omitempty" db:"logo_url"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
