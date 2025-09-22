package models

import (
	"time"
)

type StageOffer struct {
	ID             int       `json:"id" db:"id"`
	Position       string    `json:"position" db:"position"`
	Address        string    `json:"address" db:"address"`
	Lat            float64   `json:"lat" db:"lat"`
	Lng            float64   `json:"lng" db:"lng"`
	Enterprise     string    `json:"enterprise" db:"enterprise"`
	Sector         string    `json:"sector" db:"sector"`
	CapacityTotal  int       `json:"capacity_total" db:"capacity_total"`
	CapacityFilled int       `json:"capacity_filled" db:"capacity_filled"`
	Period         string    `json:"period" db:"period"`
	Course         string    `json:"course" db:"course"`
	JobFamily      string    `json:"job_family" db:"job_family"`
	ScolarLevel    string    `json:"scolar_level" db:"scolar_level"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}
