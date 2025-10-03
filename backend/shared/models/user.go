package models

import (
	"time"
)

type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleTeacher UserRole = "teacher"
	RoleTutor   UserRole = "tutor"
	RoleStudent UserRole = "eleve"
)

// Méthode pour valider le rôle
func (r UserRole) IsValid() bool {
	switch r {
	case RoleAdmin, RoleTeacher, RoleTutor, RoleStudent:
		return true
	}
	return false
}

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

func (p UserPublic) ToPublic() UserPublic {
	return p
}

// Structures pour les requêtes de création/mise à jour
type CreateUserRequest struct {
	Username        string   `json:"username" binding:"required"`
	FirstName       string   `json:"first_name" binding:"required"`
	LastName        string   `json:"last_name" binding:"required"`
	Email           string   `json:"email" binding:"required"`
	Password        string   `json:"password" binding:"required"`
	Role            UserRole `json:"role" binding:"required"`
	Phone           *string  `json:"phone,omitempty"`
	EstablishmentID *int     `json:"establishment_id,omitempty"`
}

type UpdateUserRequest struct {
	Username        *string   `json:"username,omitempty"`
	FirstName       *string   `json:"first_name,omitempty"`
	LastName        *string   `json:"last_name,omitempty"`
	Email           *string   `json:"email,omitempty"`
	Role            *UserRole `json:"role,omitempty"`
	Phone           *string   `json:"phone,omitempty"`
	EstablishmentID *int      `json:"establishment_id,omitempty"`
	IsActive        *bool     `json:"is_active,omitempty"`
}

// Méthode pour vérifier s'il y a des champs à mettre à jour
func (r *UpdateUserRequest) HasUpdates() bool {
	return r.Username != nil || r.FirstName != nil || r.LastName != nil ||
		r.Email != nil || r.Role != nil || r.Phone != nil ||
		r.EstablishmentID != nil || r.IsActive != nil
}

// Méthode pour appliquer les modifications à un utilisateur existant
func (r *UpdateUserRequest) ApplyTo(user *User) {
	if r.Username != nil {
		user.Username = *r.Username
	}
	if r.FirstName != nil {
		user.FirstName = *r.FirstName
	}
	if r.LastName != nil {
		user.LastName = *r.LastName
	}
	if r.Email != nil {
		user.Email = *r.Email
	}
	if r.Role != nil {
		user.Role = string(*r.Role)
	}
	if r.Phone != nil {
		user.Phone = r.Phone
	}
	if r.EstablishmentID != nil {
		user.EstablishmentID = r.EstablishmentID
	}
	if r.IsActive != nil {
		user.IsActive = *r.IsActive
	}

	user.UpdatedAt = time.Now()
}

type DeleteUserRequest struct {
	Id int `json:"id" binding:"required"`
}
