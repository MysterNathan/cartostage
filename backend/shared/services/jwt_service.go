package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID    int      `json:"user_id"`
	Username  string   `json:"username"`
	Role      string   `json:"role"`      // student|teacher|enterprise|admin
	EntityID  *int     `json:"entity_id"` // Nullable pour les admins
	Scope     []string `json:"scope"`
	SessionID string   `json:"session_id"`
	jwt.RegisteredClaims
}

type JWTService struct {
	jwtSecret []byte
	issuer    string
}

func NewJWTService(jwtSecret, issuer string) *JWTService {
	return &JWTService{
		jwtSecret: []byte(jwtSecret),
		issuer:    issuer,
	}
}

// GÉNÉRATION - Uniquement dans auth-service
func (s *JWTService) GenerateToken(userID int, username, role string, entityID *int, scopes []string) (string, time.Time, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	sessionID := s.generateSessionID()

	claims := CustomClaims{
		UserID:    userID,
		Username:  username,
		Role:      role,
		EntityID:  entityID,
		Scope:     scopes,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", userID),
			Issuer:    s.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecret)

	return tokenString, expirationTime, err
}

// VALIDATION - Dans tous les services
func (s *JWTService) ValidateAndParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		// Vérification expiration
		if time.Now().After(claims.ExpiresAt.Time) {
			return nil, errors.New("token expired")
		}
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *JWTService) generateSessionID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// Helpers pour vérifier les permissions
func (c *CustomClaims) HasScope(scope string) bool {
	for _, s := range c.Scope {
		if s == scope {
			return true
		}
	}
	return false
}

func (c *CustomClaims) IsRole(role string) bool {
	return c.Role == role
}

func (c *CustomClaims) CanAccessEntity(entityID int) bool {
	// Admin peut accéder à tout
	if c.Role == "admin" {
		return true
	}
	// Sinon vérifier l'entity_id
	return c.EntityID != nil && *c.EntityID == entityID
}
