package repositories

import (
	"context"
	"shared/models"
)

type UserRepositoryInterface interface {
	GetAll(ctx context.Context, filter models.UserFilter) ([]models.User, error)
	GetByID(ctx context.Context, id int) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int) error
}
