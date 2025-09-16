package models

import (
	"time"
)

type User struct {
	ID            int        `json:"id" db:"id"`
	Username      string     `json:"username" db:"username"`
	FirstName     string     `json:"first_name" db:"first_name"`
	LastName      string     `json:"last_name" db:"last_name"`
	Email         string     `json:"email" db:"email"`
	PasswordHash  string     `json:"-" db:"password_hash"` // Exclu du JSON
	Role          string     `json:"role" db:"role"`
	EntityType    *string    `json:"entity_type" db:"entity_type"` // "enterprise", "school", etc.
	EntityID      *int       `json:"entity_id" db:"entity_id"`
	IsActive      bool       `json:"is_active" db:"is_active"`
	EmailVerified bool       `json:"email_verified" db:"email_verified"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
	LastLogin     *time.Time `json:"last_login" db:"last_login"`

	// Relation avec le profil (chargé à la demande)
	Profile *UserProfile `json:"profile,omitempty" db:"-"`
}

type UserProfile struct {
	UserID      int       `json:"user_id" db:"user_id"`
	Phone       *string   `json:"phone" db:"phone"`
	Poste       *string   `json:"poste" db:"poste"`
	Departement *string   `json:"departement" db:"departement"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// DTOs pour les requêtes
type CreateUserRequest struct {
	Username    string  `json:"username" validate:"required,min=3,max=50"`
	FirstName   string  `json:"first_name" validate:"required,min=2,max=100"`
	LastName    string  `json:"last_name" validate:"required,min=2,max=100"`
	Email       string  `json:"email" validate:"required,email"`
	Password    string  `json:"password" validate:"required,min=8"`
	Role        string  `json:"role" validate:"required,oneof=admin tutor student teacher enterprise"`
	EntityType  *string `json:"entity_type"`
	EntityID    *int    `json:"entity_id"`
	Phone       *string `json:"phone"`
	Poste       *string `json:"poste"`
	Departement *string `json:"departement"`
}

type UpdateUserRequest struct {
	Username   *string `json:"username,omitempty"`
	FirstName  *string `json:"first_name,omitempty"`
	LastName   *string `json:"last_name,omitempty"`
	Email      *string `json:"email,omitempty"`
	Role       *string `json:"role,omitempty"`
	EntityType *string `json:"entity_type,omitempty"`
	EntityID   *int    `json:"entity_id,omitempty"`
	IsActive   *bool   `json:"is_active,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	User      User   `json:"user"`
	ExpiresAt int64  `json:"expires_at"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type UpdateUserProfileRequest struct {
	Phone       *string `json:"phone,omitempty"`
	Poste       *string `json:"poste,omitempty"`
	Departement *string `json:"departement,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

// Méthodes utilitaires
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

func (u *User) IsTutor() bool {
	return u.Role == "tutor"
}

func (u *User) IsStudent() bool {
	return u.Role == "student"
}

func (u *User) IsTeacher() bool {
	return u.Role == "teacher"
}

func (u *User) IsEnterprise() bool {
	return u.Role == "enterprise"
}

func (u *User) HasEntity() bool {
	return u.EntityType != nil && u.EntityID != nil
}

func (u *User) GetEntityType() string {
	if u.EntityType == nil {
		return ""
	}
	return *u.EntityType
}

func (u *User) GetEntityID() int {
	if u.EntityID == nil {
		return 0
	}
	return *u.EntityID
}

func (u *User) GetFullName() string {
	if u.FirstName != "" && u.LastName != "" {
		return u.FirstName + " " + u.LastName
	}
	if u.FirstName != "" {
		return u.FirstName
	}
	if u.LastName != "" {
		return u.LastName
	}
	return u.Username
}

func (u *User) GetDisplayName() string {
	fullName := u.GetFullName()
	if fullName != u.Username {
		return fullName
	}
	return u.Username
}

// ToPublic - Retourne une version publique sans données sensibles
func (u *User) ToPublic() *User {
	public := *u
	public.PasswordHash = "" // Toujours vider le hash
	return &public
}
