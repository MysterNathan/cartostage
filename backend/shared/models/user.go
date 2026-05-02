// models/user.go
package models

import "time"

type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleTeacher UserRole = "teacher"
	RoleTutor   UserRole = "tutor"
	RoleStudent UserRole = "student"
)

func (r UserRole) IsValid() bool {
	switch r {
	case RoleAdmin, RoleTeacher, RoleTutor, RoleStudent:
		return true
	}
	return false
}

type User struct {
	ID              int        `db:"id"`
	Username        string     `db:"username"`
	FirstName       string     `db:"first_name"`
	LastName        string     `db:"last_name"`
	Email           string     `db:"email"`
	PasswordHash    string     `db:"password_hash"`
	Role            string     `db:"role"`
	Phone           *string    `db:"phone"`
	EstablishmentID *int       `db:"establishment_id"`
	IsActive        bool       `db:"is_active"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
	LastLogin       *time.Time `db:"last_login"`
}

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
	}
}

type UserPublic struct {
	ID              int      `json:"id"`
	Username        string   `json:"username"`
	FirstName       string   `json:"first_name"`
	LastName        string   `json:"last_name"`
	Email           string   `json:"email"`
	Role            UserRole `json:"role"`
	Phone           *string  `json:"phone,omitempty"`
	EstablishmentID *int     `json:"establishment_id,omitempty"`
}

type UserFilter struct {
	RequestorRole UserRole
	RequestorID   int
}

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
func (r *UpdateUserRequest) HasUpdates() bool {
	return r.Username != nil || r.FirstName != nil || r.LastName != nil ||
		r.Email != nil || r.Role != nil || r.Phone != nil ||
		r.EstablishmentID != nil || r.IsActive != nil
}

type DeleteUserRequest struct {
	ID int `json:"id" binding:"required"`
}
