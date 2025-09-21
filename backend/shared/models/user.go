package models

import (
	"time"
)

// Types d'enum pour plus de clarté
type UserRole string

const (
	RoleAdmin   UserRole = "administrateur"
	RoleTeacher UserRole = "enseignant"
	RoleTutor   UserRole = "tuteur"
	RoleStudent UserRole = "eleve"
)

// Structure principale User
type User struct {
	ID           int        `json:"id" db:"id"`
	Username     string     `json:"username" db:"username"`
	FirstName    string     `json:"first_name" db:"first_name"`
	LastName     string     `json:"last_name" db:"last_name"`
	Email        string     `json:"email" db:"email"`
	PasswordHash string     `json:"-" db:"password_hash"`
	Role         UserRole   `json:"role" db:"role"`
	Phone        *string    `json:"phone,omitempty" db:"phone"`
	EntityID     int        `json:"entity_id" db:"entity_id"`
	IsActive     bool       `json:"is_active" db:"is_active"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	LastLogin    *time.Time `json:"last_login,omitempty" db:"last_login"`
}

// DTOs pour les requêtes
type CreateUserRequest struct {
	Username  string   `json:"username" validate:"required,min=3,max=50"`
	FirstName string   `json:"first_name" validate:"required,min=2,max=100"`
	LastName  string   `json:"last_name" validate:"required,min=2,max=100"`
	Email     string   `json:"email" validate:"required,email"`
	Password  string   `json:"password" validate:"required,min=8"`
	Role      UserRole `json:"role" validate:"required"`
	Phone     *string  `json:"phone,omitempty" validate:"omitempty,min=10,max=15"`
	EntityID  int      `json:"entity_id" validate:"required"`
}

// Structure pour l'update - tous les champs optionnels avec des pointeurs
type UpdateUserRequest struct {
	Username  *string   `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
	FirstName *string   `json:"first_name,omitempty" validate:"omitempty,min=2,max=100"`
	LastName  *string   `json:"last_name,omitempty" validate:"omitempty,min=2,max=100"`
	Email     *string   `json:"email,omitempty" validate:"omitempty,email"`
	Role      *UserRole `json:"role,omitempty"`
	Phone     *string   `json:"phone,omitempty" validate:"omitempty,min=10,max=15"`
	EntityID  *int      `json:"entity_id,omitempty"`
	IsActive  *bool     `json:"is_active,omitempty"`
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

type DeleteUserRequest struct {
	Id   int      `json:"id"`
	Role UserRole `json:"role"`
}

// Méthodes utilitaires pour UpdateUserRequest
func (u *UpdateUserRequest) HasUpdates() bool {
	return u.Username != nil || u.FirstName != nil || u.LastName != nil ||
		u.Email != nil || u.Role != nil || u.Phone != nil ||
		u.EntityID != nil || u.IsActive != nil
}

// Méthode pour appliquer les updates à un user existant
func (u *UpdateUserRequest) ApplyTo(user *User) {
	if u.Username != nil {
		user.Username = *u.Username
	}
	if u.FirstName != nil {
		user.FirstName = *u.FirstName
	}
	if u.LastName != nil {
		user.LastName = *u.LastName
	}
	if u.Email != nil {
		user.Email = *u.Email
	}
	if u.Role != nil {
		user.Role = *u.Role
	}
	if u.Phone != nil {
		user.Phone = u.Phone
	}
	if u.EntityID != nil {
		user.EntityID = *u.EntityID
	}
	if u.IsActive != nil {
		user.IsActive = *u.IsActive
	}
	user.UpdatedAt = time.Now()
}

// Méthodes utilitaires existantes
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u *User) IsTeacher() bool {
	return u.Role == RoleTeacher
}

func (u *User) IsTutor() bool {
	return u.Role == RoleTutor
}

func (u *User) IsStudent() bool {
	return u.Role == RoleStudent
}

func (u *User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}

func (u *User) GetDisplayName() string {
	fullName := u.GetFullName()
	if fullName != " " {
		return fullName
	}
	return u.Username
}

func (u *User) ToPublic() User {
	public := *u
	public.PasswordHash = ""
	return public
}

func (r UserRole) IsValid() bool {
	switch r {
	case RoleAdmin, RoleTeacher, RoleTutor, RoleStudent:
		return true
	}
	return false
}

func ParseUserRole(role string) (UserRole, bool) {
	switch role {
	case "administrateur", "admin":
		return RoleAdmin, true
	case "enseignant", "teacher":
		return RoleTeacher, true
	case "tuteur", "tutor":
		return RoleTutor, true
	case "eleve", "student":
		return RoleStudent, true
	}
	return "", false
}
