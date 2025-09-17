package models

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type CustomClaims struct {
	UserID    int      `json:"user_id"`
	Username  string   `json:"username"`
	Role      UserRole `json:"role"`
	SessionID string   `json:"session_id"`
	jwt.RegisteredClaims
}

// Méthodes pour CustomClaims
func (c *CustomClaims) IsRole(role UserRole) bool {
	return c.Role == role
}

func (c *CustomClaims) IsAdmin() bool {
	return c.Role == RoleAdmin
}

func (c *CustomClaims) IsTeacher() bool {
	return c.Role == RoleTeacher
}

func (c *CustomClaims) IsTutor() bool {
	return c.Role == RoleTutor
}

func (c *CustomClaims) IsStudent() bool {
	return c.Role == RoleStudent
}

// Helper pour créer des claims à partir d'un User
func NewCustomClaims(user User, sessionID string, expiresAt int64) *CustomClaims {
	return &CustomClaims{
		UserID:    user.ID,
		Username:  user.Username,
		Role:      user.Role,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(expiresAt, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
}
