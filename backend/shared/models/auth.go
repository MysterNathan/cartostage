package models

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
