package models

import "time"

type Enterprise struct {
	ID          int       `json:"id" db:"id"`
	Nom         string    `json:"nom" db:"nom"`
	Adresse     string    `json:"adresse" db:"adresse"`
	Secteur     string    `json:"secteur" db:"secteur"`
	Taille      string    `json:"taille" db:"taille"`
	Siret       string    `json:"siret" db:"siret"`
	Email       string    `json:"email" db:"email"`
	Telephone   string    `json:"telephone" db:"telephone"`
	SiteWeb     string    `json:"site_web" db:"site_web"`
	Description string    `json:"description" db:"description"`
	LogoURL     string    `json:"logo_url" db:"logo_url"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type EnterpriseWithStats struct {
	Enterprise
	TotalTutors   int `json:"total_tutors"`
	ActiveStages  int `json:"active_stages"`
	TotalStudents int `json:"total_students"`
}
