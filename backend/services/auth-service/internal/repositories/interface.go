package repositories

import "shared/models"

type AuthRepositoryInterface interface {
	FindUserByUsername(username string) (*models.User, error)
}
