package services

import (
	"errors"
	"shared/models"
	"shared/repositories"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// GetAll récupère tous les utilisateurs
func (s *UserService) GetAll() ([]*models.User, error) {
	return s.repo.GetAll()
}

// GetByID récupère un utilisateur par son ID
func (s *UserService) GetByID(id int) (*models.User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user ID")
	}
	return s.repo.GetByID(id)
}

// GetByEmail récupère un utilisateur par son email
func (s *UserService) GetByEmail(email string) (*models.User, error) {
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}
	return s.repo.GetByEmail(email)
}

// GetByUsername récupère un utilisateur par son nom d'utilisateur
func (s *UserService) GetByUsername(username string) (*models.User, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}
	return s.repo.GetByUsername(username)
}

// GetByRole récupère tous les utilisateurs ayant un rôle spécifique
func (s *UserService) GetByRole(role string) ([]*models.User, error) {
	if !s.isValidRole(role) {
		return nil, errors.New("invalid role")
	}
	return s.repo.GetByRole(role)
}

// GetByEntity récupère tous les utilisateurs d'une entité spécifique
func (s *UserService) GetByEntity(entityType string, entityID int) ([]*models.User, error) {
	if entityType == "" {
		return nil, errors.New("entity type cannot be empty")
	}
	if entityID <= 0 {
		return nil, errors.New("invalid entity ID")
	}
	return s.repo.GetByEntity(entityType, entityID)
}

// Create crée un nouvel utilisateur
func (s *UserService) Create(req *models.CreateUserRequest) (*models.User, error) {
	// Validation des données requises
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Vérifier l'unicité de l'email
	existing, _ := s.repo.GetByEmail(req.Email)
	if existing != nil {
		return nil, errors.New("email already exists")
	}

	// Vérifier l'unicité du nom d'utilisateur
	existing, _ = s.repo.GetByUsername(req.Username)
	if existing != nil {
		return nil, errors.New("username already exists")
	}

	// Hasher le mot de passe
	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Créer l'objet User
	now := time.Now()
	user := &models.User{
		Username:      req.Username,
		Email:         req.Email,
		PasswordHash:  hashedPassword,
		Role:          req.Role,
		EntityType:    req.EntityType,
		EntityID:      req.EntityID,
		IsActive:      true,
		EmailVerified: false,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Créer le profil s'il y a des données
	var profile *models.UserProfile
	if s.hasProfileData(req) {
		profile = &models.UserProfile{
			FirstName:   req.FirstName,
			LastName:    req.LastName,
			Phone:       req.Phone,
			Poste:       req.Poste,
			Departement: req.Departement,
			IsActive:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
	}

	return s.repo.Create(user, profile)
}

// Update met à jour un utilisateur
func (s *UserService) Update(id int, req *models.UpdateUserRequest) (*models.User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user ID")
	}

	// Vérifier que l'utilisateur existe
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Validation des données
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, err
	}

	// Vérifier l'unicité de l'email si modifié
	if req.Email != nil && *req.Email != user.Email {
		existing, _ := s.repo.GetByEmail(*req.Email)
		if existing != nil {
			return nil, errors.New("email already exists")
		}
	}

	// Vérifier l'unicité du nom d'utilisateur si modifié
	if req.Username != nil && *req.Username != user.Username {
		existing, _ := s.repo.GetByUsername(*req.Username)
		if existing != nil {
			return nil, errors.New("username already exists")
		}
	}

	return s.repo.Update(id, req)
}

// Delete supprime un utilisateur
func (s *UserService) Delete(id int) error {
	if id <= 0 {
		return errors.New("invalid user ID")
	}

	// Vérifier que l'utilisateur existe
	_, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(id)
}

// UpdateProfile met à jour le profil utilisateur
func (s *UserService) UpdateProfile(id int, req *models.UpdateUserProfileRequest) (*models.User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user ID")
	}

	// Vérifier que l'utilisateur existe
	_, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.repo.UpdateProfile(id, req)
}

// ChangePassword change le mot de passe d'un utilisateur
func (s *UserService) ChangePassword(id int, req *models.ChangePasswordRequest) error {
	if id <= 0 {
		return errors.New("invalid user ID")
	}

	// Récupérer l'utilisateur
	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// Vérifier l'ancien mot de passe
	if !s.checkPassword(req.OldPassword, user.PasswordHash) {
		return errors.New("old password is incorrect")
	}

	// Hasher le nouveau mot de passe
	hashedPassword, err := s.hashPassword(req.NewPassword)
	if err != nil {
		return errors.New("failed to hash new password")
	}

	return s.repo.UpdatePassword(id, hashedPassword)
}

// AuthenticateUser authentifie un utilisateur avec username/password
func (s *UserService) AuthenticateUser(username, password string) (*models.User, error) {
	if username == "" || password == "" {
		return nil, errors.New("username and password are required")
	}

	// Récupérer l'utilisateur par nom d'utilisateur
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Vérifier que l'utilisateur est actif
	if !user.IsActive {
		return nil, errors.New("user account is disabled")
	}

	// Vérifier le mot de passe
	if !s.checkPassword(password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	// Mettre à jour la dernière connexion
	now := time.Now()
	s.repo.UpdateLastLogin(user.ID, &now)

	return user, nil
}

// UpdateLastLogin met à jour la dernière connexion
func (s *UserService) UpdateLastLogin(id int) error {
	now := time.Now()
	return s.repo.UpdateLastLogin(id, &now)
}

// ActivateUser active un utilisateur
func (s *UserService) ActivateUser(id int) error {
	return s.repo.UpdateStatus(id, true)
}

// DeactivateUser désactive un utilisateur
func (s *UserService) DeactivateUser(id int) error {
	return s.repo.UpdateStatus(id, false)
}

// Méthodes utilitaires privées

func (s *UserService) validateCreateRequest(req *models.CreateUserRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}
	if len(req.Username) < 3 || len(req.Username) > 50 {
		return errors.New("username must be between 3 and 50 characters")
	}
	if req.Email == "" {
		return errors.New("email is required")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	if len(req.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if !s.isValidRole(req.Role) {
		return errors.New("invalid role")
	}
	return nil
}

func (s *UserService) validateUpdateRequest(req *models.UpdateUserRequest) error {
	if req.Username != nil {
		if len(*req.Username) < 3 || len(*req.Username) > 50 {
			return errors.New("username must be between 3 and 50 characters")
		}
	}
	if req.Role != nil && !s.isValidRole(*req.Role) {
		return errors.New("invalid role")
	}
	return nil
}

func (s *UserService) isValidRole(role string) bool {
	validRoles := []string{"admin", "tutor", "student", "enterprise"}
	for _, validRole := range validRoles {
		if role == validRole {
			return true
		}
	}
	return false
}

func (s *UserService) hasProfileData(req *models.CreateUserRequest) bool {
	return req.FirstName != nil || req.LastName != nil || req.Phone != nil ||
		req.Poste != nil || req.Departement != nil
}

func (s *UserService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (s *UserService) checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
