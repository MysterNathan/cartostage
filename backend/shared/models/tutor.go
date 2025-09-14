package models

import "time"

type Tutor struct {
	ID           int       `json:"id" db:"id"`
	EnterpriseID int       `json:"enterprise_id" db:"enterprise_id"`
	Prenom       string    `json:"prenom" db:"prenom"`
	Nom          string    `json:"nom" db:"nom"`
	Email        string    `json:"email" db:"email"`
	Telephone    *string   `json:"telephone" db:"telephone"`
	Poste        *string   `json:"poste" db:"poste"`
	Departement  *string   `json:"departement" db:"departement"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`

	// Données de l'enterprise (pour les jointures)
	EnterpriseName *string `json:"enterprise_name,omitempty" db:"enterprise_name"`
}

type TutorWithEnterprise struct {
	Tutor
	EnterpriseName string `json:"enterprise_name"`
}

type TutorWithStats struct {
	Tutor
	ActiveStudents int `json:"active_students"`
	TotalStages    int `json:"total_stages"`
}
