package services

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	sharedContext "shared/context"
	"shared/models"
	"shared/repositories"
	"time"
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

func (s *UserService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
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

func (s *UserService) Delete(ctx context.Context, req *models.DeleteUserRequest) (int, error) {
	// Vérifier les permissions
	claims := sharedContext.GetUserClaims(ctx)
	if claims == nil {
		return 0, fmt.Errorf("unauthorized: no claims in context")
	}

	// Seuls les admins peuvent supprimer des utilisateurs pour l'instant
	if !claims.IsAdmin() {
		return 0, fmt.Errorf("forbidden: only admins can delete users")
	}

	// Validation basique
	if req.Id == 0 {
		return 0, fmt.Errorf("id user is required")
	}

	err := s.userRepo.Delete(ctx, req.Id)
	if err != nil {
		return 0, err
	}

	return req.Id, nil
}
func (s *UserService) UpdateUser(ctx context.Context, userID int, req *models.UpdateUserRequest) (*models.User, error) {
	// Vérifier les permissions
	claims := sharedContext.GetUserClaims(ctx)
	if claims == nil {
		return nil, fmt.Errorf("unauthorized: no claims in context")
	}

	// Validation basique de l'ID
	if userID <= 0 {
		return nil, fmt.Errorf("validation: invalid user ID")
	}

	// Vérifier s'il y a des champs à mettre à jour
	if !req.HasUpdates() {
		return nil, fmt.Errorf("validation: no fields to update")
	}

	// Récupérer l'utilisateur existant
	existing, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to load user: %w", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("not found: user with id %d not found", userID)
	}

	// Vérifier les permissions spécifiques
	// Un admin peut tout modifier
	// Un utilisateur ne peut modifier que ses propres données (et certains champs seulement)
	if !claims.IsAdmin() {
		// Un utilisateur ne peut modifier que son propre profil
		if claims.UserID != userID {
			return nil, fmt.Errorf("forbidden: can only update your own profile")
		}

		// Les utilisateurs non-admin ne peuvent pas modifier certains champs critiques
		if req.Role != nil {
			return nil, fmt.Errorf("forbidden: cannot change role")
		}
		if req.IsActive != nil {
			return nil, fmt.Errorf("forbidden: cannot change active status")
		}
		if req.EntityID != nil {
			return nil, fmt.Errorf("forbidden: cannot change entity")
		}
	}

	// Validation additionnelle du rôle si fourni
	if req.Role != nil && !req.Role.IsValid() {
		return nil, fmt.Errorf("validation: invalid role")
	}

	// Vérifier l'unicité du username et email si modifiés
	if req.Username != nil && *req.Username != existing.Username {
		existingByUsername, err := s.userRepo.GetByUsername(ctx, *req.Username)
		if err != nil {
			return nil, fmt.Errorf("failed to check username uniqueness: %w", err)
		}
		if existingByUsername != nil {
			return nil, fmt.Errorf("duplicate: username already exists")
		}
	}

	if req.Email != nil && *req.Email != existing.Email {
		existingByEmail, err := s.userRepo.GetByEmail(ctx, *req.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to check email uniqueness: %w", err)
		}
		if existingByEmail != nil {
			return nil, fmt.Errorf("duplicate: email already exists")
		}
	}

	// Appliquer les modifications en utilisant la méthode ApplyTo
	req.ApplyTo(existing)

	// Persister les changements
	err = s.userRepo.Update(ctx, existing, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Retourner la version publique de l'utilisateur mis à jour
	publicUser := existing.ToPublic()
	return &publicUser, nil
}
