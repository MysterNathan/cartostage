package models

import (
	"time"
)

// User - Modèle complet pour la base de données (interne)
type User struct {
	ID           int        `json:"id" db:"id"`
	Username     string     `json:"username" db:"username"`
	Email        *string    `json:"email" db:"email"`
	PasswordHash string     `json:"-" db:"password_hash"` // Jamais exposé en JSON
	Role         string     `json:"role" db:"role"`
	EntityID     *int       `json:"entity_id" db:"entity_id"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	LastLogin    *time.Time `json:"last_login" db:"last_login"`
}

// UserProfile - Modèle pour l'API (exposé publiquement)
type UserProfile struct {
	ID        int        `json:"id"`
	Username  string     `json:"username"`
	Email     *string    `json:"email"`
	Role      string     `json:"role"`
	EntityID  *int       `json:"entity_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	LastLogin *time.Time `json:"last_login"`
}

// ToProfile - Convertit User vers UserProfile (sans le hash)
func (u *User) ToProfile() *UserProfile {
	return &UserProfile{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role,
		EntityID:  u.EntityID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		LastLogin: u.LastLogin,
	}
}

// UserInfo - Version légère pour les réponses de login
type UserInfo struct {
	ID       int     `json:"id"`
	Username string  `json:"username"`
	Email    *string `json:"email"`
	Role     string  `json:"role"`
	EntityID *int    `json:"entity_id"`
}

// ToUserInfo - Convertit User vers UserInfo
func (u *User) ToUserInfo() UserInfo {
	return UserInfo{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Role:     u.Role,
		EntityID: u.EntityID,
	}
}

// Requests DTOs
type CreateUserRequest struct {
	Username string  `json:"username" validate:"required,min=3,max=50"`
	Email    *string `json:"email" validate:"omitempty,email"`
	Password string  `json:"password" validate:"required,min=6"`
	Role     string  `json:"role" validate:"required,oneof=student teacher enterprise admin"`
	EntityID *int    `json:"entity_id"`
}

type UpdateUserRequest struct {
	Username *string `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
	Email    *string `json:"email,omitempty" validate:"omitempty,email"`
	Password *string `json:"password,omitempty" validate:"omitempty,min=6"`
	Role     *string `json:"role,omitempty" validate:"omitempty,oneof=student teacher enterprise admin"`
	EntityID *int    `json:"entity_id,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	User      UserInfo  `json:"user"`
}
