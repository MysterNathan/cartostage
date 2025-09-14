package services

import (
	"errors"
	"time"

	"auth/internal/repositories"
	"shared/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	jwtSecret      []byte
	userRepository *repositories.UserRepository
}

func NewAuthService(jwtSecret string, userRepository *repositories.UserRepository) *AuthService {
	return &AuthService{
		jwtSecret:      []byte(jwtSecret),
		userRepository: userRepository,
	}
}

func (s *AuthService) ValidateCredentials(username, password string) (*models.User, error) {
	user, err := s.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *AuthService) GenerateToken(username string) (string, time.Time, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &jwt.RegisteredClaims{
		Subject:   username,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecret)

	return tokenString, expirationTime, err
}

// Utilitaire pour hasher un mot de passe
func (s *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Créer un nouvel utilisateur
func (s *AuthService) CreateUser(username, password, role string) (*models.User, error) {
	// Vérifier si l'utilisateur existe déjà
	existingUser, err := s.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Hasher le mot de passe
	hashedPassword, err := s.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Créer l'utilisateur
	user := &models.User{
		Username:     username,
		PasswordHash: hashedPassword,
		Role:         role,
	}

	err = s.userRepository.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
