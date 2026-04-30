package services

import (
	"auth/internal/repositories"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"shared/models"
	"shared/services"
	"time"
)

type AuthService struct {
	authRepo   *repositories.AuthRepository
	jwtService *services.JWTService
}

func NewAuthService(authRepo *repositories.AuthRepository, jwtService *services.JWTService) *AuthService {
	return &AuthService{
		authRepo:   authRepo,
		jwtService: jwtService,
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

// Login - SIMPLE et DIRECT
func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	// Validation basique
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("email and password required")
	}

	// Trouver l'utilisateur
	user, err := s.authRepo.FindUserByUsername(req.Username)
	if err != nil || user == nil {
		return nil, errors.New("invalid credentials")
	}

	// Vérifier le mot de passe
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		return nil, errors.New("invalid credentials")
	}

	// Générer le token JWT SIMPLE
	claims := models.NewCustomClaims(*user, "simple-session", time.Now().Add(24*time.Hour).Unix())
	token, err := s.jwtService.GenerateToken(claims)
	if err != nil {
		return nil, err
	}

	// Nettoyer le password
	user.PasswordHash = ""

	return &LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}
