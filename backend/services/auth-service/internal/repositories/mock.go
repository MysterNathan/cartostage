package repositories

import (
	"errors"
	"shared/models"
)

type MockAuthRepository struct {
	Users         map[string]*models.User
	ErrorToReturn error
}

func NewMockAuthRepository() *MockAuthRepository {
	return &MockAuthRepository{
		Users: make(map[string]*models.User),
	}
}

func (f *MockAuthRepository) FindUserByUsername(username string) (*models.User, error) {
	if f.ErrorToReturn != nil {
		return nil, f.ErrorToReturn
	}

	user, exists := f.Users[username]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}
