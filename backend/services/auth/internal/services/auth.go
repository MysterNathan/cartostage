package services

import (
	"auth/internal/repositories"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"shared/models"
	"shared/services"
)

type AuthService struct {
	userRepo   *repositories.UserRepository
	jwtService *services.JWTService
}

func NewAuthService(userRepo *repositories.UserRepository, jwtService *services.JWTService) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (s *AuthService) Login(username, password string) (*models.LoginResponse, error) {
	// Récupérer l'utilisateur complet (avec hash)
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	// Vérifier le mot de passe
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Mettre à jour last login
	s.userRepo.UpdateLastLogin(user.ID)

	// Générer le token
	scopes := s.getScopesForRole(user.Role)
	token, expiresAt, err := s.jwtService.GenerateToken(
		user.ID,
		user.Username,
		user.Role,
		user.EntityID,
		scopes,
	)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      user.ToUserInfo(), // Conversion sans hash
	}, nil
}

func (s *AuthService) GetUserProfile(userID int) (*models.UserProfile, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return user.ToProfile(), nil // Conversion sans hash
}

func (s *AuthService) CreateUser(req models.CreateUserRequest) (*models.UserProfile, error) {
	// Vérifier si l'utilisateur existe déjà
	exists, err := s.userRepo.UserExists(req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	// Hasher le mot de passe
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Créer l'utilisateur (modèle User avec hash)
	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         req.Role,
		EntityID:     req.EntityID,
	}

	err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return user.ToProfile(), nil // Retourner UserProfile sans hash
}

func (s *AuthService) UpdateUser(userID int, req models.UpdateUserRequest) (*models.UserProfile, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	// Mettre à jour les champs fournis
	if req.Username != nil {
		user.Username = *req.Username
	}
	if req.Email != nil {
		user.Email = req.Email
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.EntityID != nil {
		user.EntityID = req.EntityID
	}
	if req.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = string(hashedPassword)
	}

	err = s.userRepo.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	return user.ToProfile(), nil // Retourner UserProfile sans hash
}
func (s *AuthService) RefreshToken(oldToken string) (*models.LoginResponse, error) {
	claims, err := s.jwtService.ValidateAndParseToken(oldToken)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	// Nouveau token avec mêmes permissions
	newToken, expiresAt, err := s.jwtService.GenerateToken(
		user.ID,
		user.Username,
		user.Role,
		user.EntityID,
		claims.Scope,
	)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Token:     newToken,
		ExpiresAt: expiresAt,
		User: models.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
			EntityID: user.EntityID,
		},
	}, nil
}

func (s *AuthService) Logout(sessionID string) error {
	// TODO: Ajouter à une blacklist de tokens si nécessaire
	return nil
}

func (s *AuthService) DeleteUser(userID int) error {
	return s.userRepo.DeleteUser(userID)
}

func (s *AuthService) GetUsers(page, limit int, roleFilter, entityIDFilter string) ([]*models.UserProfile, int, error) {
	return s.userRepo.GetUsers(page, limit, roleFilter, entityIDFilter)
}

func (s *AuthService) getScopesForRole(role string) []string {
	switch role {
	case "student":
		return []string{
			"read:profile",
			"write:profile",
			"read:stages",
			"write:applications",
			"read:applications",
		}
	case "teacher":
		return []string{
			"read:profile",
			"write:profile",
			"read:stages",
			"write:stages",
			"read:students",
			"read:applications",
		}
	case "enterprise":
		return []string{
			"read:profile",
			"write:profile",
			"read:stages",
			"write:stages",
			"read:applications",
			"write:applications",
		}
	case "admin":
		return []string{
			"read:*",
			"write:*",
			"delete:*",
		}
	default:
		return []string{"read:profile"}
	}
}
