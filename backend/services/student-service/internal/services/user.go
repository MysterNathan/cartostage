// service/user_service.go
package services

import (
	"context"
	"fmt"
	sharedContext "shared/context"

	"shared/models"
	"student/internal/repositories"
)

type UserService interface {
	GetAll(ctx context.Context) ([]models.UserPublic, error)
	GetByID(ctx context.Context, id int) (*models.UserPublic, error)
	Create(ctx context.Context, req models.CreateUserRequest) (*models.UserPublic, error)
	Update(ctx context.Context, id int, req models.UpdateUserRequest) (*models.UserPublic, error)
	Delete(ctx context.Context, id int) error
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetAll(ctx context.Context) ([]models.UserPublic, error) {
	claims := sharedContext.GetClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no claims found in context")
	}

	filter := models.UserFilter{
		RequestorRole: models.UserRole(claims.Role),
		RequestorID:   claims.UserID,
	}

	users, err := s.userRepo.GetAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	userPublics := make([]models.UserPublic, len(users))
	for i, user := range users {
		userPublics[i] = user.ToPublic()
	}

	return userPublics, nil
}

func (s *userService) GetByID(ctx context.Context, id int) (*models.UserPublic, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	userPublic := user.ToPublic()
	return &userPublic, nil
}

func (s *userService) Create(ctx context.Context, req models.CreateUserRequest) (*models.UserPublic, error) {
	// TODO: Hash le password
	user := &models.User{
		Username:        req.Username,
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Email:           req.Email,
		PasswordHash:    req.Password, // À hasher !
		Role:            string(req.Role),
		Phone:           req.Phone,
		EstablishmentID: req.EstablishmentID,
		IsActive:        true,
	}

	err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	userPublic := user.ToPublic()
	return &userPublic, nil
}

func (s *userService) Update(ctx context.Context, id int, req models.UpdateUserRequest) (*models.UserPublic, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	req.ApplyTo(user)

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	userPublic := user.ToPublic()
	return &userPublic, nil
}

func (s *userService) Delete(ctx context.Context, id int) error {
	return s.userRepo.Delete(ctx, id)
}
