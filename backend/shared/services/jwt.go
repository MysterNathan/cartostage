// shared/services/jwt.go
package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"shared/models"
)

type JWTService struct {
	secretKey []byte
	issuer    string
}

func NewJWTService(secret, issuer string) *JWTService {
	return &JWTService{
		secretKey: []byte(secret),
		issuer:    issuer,
	}
}

func (j *JWTService) ValidateAndParseToken(tokenString string) (*models.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
		// Vérifier l'expiration
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			return nil, errors.New("token expired")
		}

		// Vérifier l'issuer si configuré
		if j.issuer != "" && claims.Issuer != j.issuer {
			return nil, errors.New("invalid issuer")
		}

		return claims, nil
	}

	return nil, errors.New("invalid token claims")
}

func (j *JWTService) GenerateToken(claims *models.CustomClaims) (string, error) {
	// Si les RegisteredClaims ne sont pas encore définis, les définir
	if claims.ExpiresAt == nil {
		claims.RegisteredClaims = jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    j.issuer,
		}
	} else {
		// Si déjà définis, juste s'assurer que l'issuer est correct
		if claims.Issuer == "" {
			claims.Issuer = j.issuer
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}
