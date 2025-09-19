package services

import (
	"context"
	"fmt"
	"log"
	sharedContext "shared/context"
	"shared/models"
	"shared/repositories"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetUsersByRole(ctx context.Context, targetRole models.UserRole) ([]models.User, error) {
	// Vérifier les permissions via le contexte
	claims := sharedContext.GetUserClaims(ctx)
	if claims == nil {
		return nil, fmt.Errorf("unauthorized: no claims in context")
	}

	log.Println(claims)
	// Règles de permission : chaque rôle ne voit que son propre rôle (sauf admin)
	if !claims.IsAdmin() && claims.Role != targetRole {
		return nil, fmt.Errorf("forbidden: cannot access %s data", targetRole)
	}

	users, err := s.userRepo.GetByRole(ctx, targetRole)
	if err != nil {
		return nil, err
	}

	// Convertir en version publique (sans password_hash)
	publicUsers := make([]models.User, len(users))
	for i, user := range users {
		publicUsers[i] = user.ToPublic()
	}

	return publicUsers, nil
}

type CreateUserRequest struct {
	Username  string          `json:"username"`
	FirstName string          `json:"first_name"`
	LastName  string          `json:"last_name"`
	Email     string          `json:"email"`
	Password  string          `json:"password"`
	Role      models.UserRole `json:"role"`
}

func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*models.User, error) {
	// Vérifier les permissions
	claims := sharedContext.GetUserClaims(ctx)
	if claims == nil {
		return nil, fmt.Errorf("unauthorized: no claims in context")
	}

	// Seuls les admins peuvent créer des utilisateurs pour l'instant
	if !claims.IsAdmin() {
		return nil, fmt.Errorf("forbidden: only admins can create users")
	}

	// Validation basique
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return nil, fmt.Errorf("username, email and password are required")
	}

	if !req.Role.IsValid() {
		return nil, fmt.Errorf("invalid role: %s", req.Role)
	}

	// Hasher le mot de passe
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Créer l'utilisateur
	now := time.Now()
	user := &models.User{
		Username:     req.Username,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         req.Role,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	// Retourner la version publique
	publicUser := user.ToPublic()
	return &publicUser, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *CreateUserRequest) (*models.User, error) {
	
}
