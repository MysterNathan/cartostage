package models

import (
	"time"
)

type UserRole string

const (
	RoleAdmin   UserRole = "administrateur"
	RoleTeacher UserRole = "enseignant"
	RoleTutor   UserRole = "tuteur"
	RoleStudent UserRole = "eleve"
)

type User struct {
	ID              int        `json:"id" db:"id"`
	Username        string     `json:"username" db:"username"`
	FirstName       string     `json:"first_name" db:"first_name"`
	LastName        string     `json:"last_name" db:"last_name"`
	Email           string     `json:"email" db:"email"`
	PasswordHash    string     `json:"-" db:"password_hash"` // Jamais exposé en JSON
	Role            string     `json:"role" db:"role"`
	Phone           *string    `json:"phone,omitempty" db:"phone"`
	EstablishmentID *int       `json:"establishment_id,omitempty" db:"establishment_id"`
	IsActive        bool       `json:"is_active" db:"is_active"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
	LastLogin       *time.Time `json:"last_login,omitempty" db:"last_login"`
}

// Méthode utilitaire pour obtenir le nom complet
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// Conversion en version publique
func (u *User) ToPublic() UserPublic {
	return UserPublic{
		ID:              u.ID,
		Username:        u.Username,
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		Email:           u.Email,
		Role:            UserRole(u.Role),
		Phone:           u.Phone,
		EstablishmentID: u.EstablishmentID,
		IsActive:        u.IsActive,
		CreatedAt:       u.CreatedAt,
		LastLogin:       u.LastLogin,
	}
}

// Structure pour les réponses publiques (sans données sensibles)
type UserPublic struct {
	ID              int        `json:"id"`
	Username        string     `json:"username"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	Email           string     `json:"email"`
	Role            UserRole   `json:"role"`
	Phone           *string    `json:"phone,omitempty"`
	EstablishmentID *int       `json:"establishment_id,omitempty"`
	IsActive        bool       `json:"is_active"`
	CreatedAt       time.Time  `json:"created_at"`
	LastLogin       *time.Time `json:"last_login,omitempty"`
}
