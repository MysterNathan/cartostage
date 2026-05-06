package services

import (
	"context"
	"shared/models"
)

type UserService interface {
	GetAll(ctx context.Context) ([]models.UserPublic, error)
	GetByID(ctx context.Context, id int) (*models.UserPublic, error)
	Create(ctx context.Context, req models.CreateUserRequest) (*models.UserPublic, error)
	Update(ctx context.Context, id int, req models.UpdateUserRequest) (*models.UserPublic, error)
	Delete(ctx context.Context, id int) error
}
