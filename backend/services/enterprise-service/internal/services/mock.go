package services

import (
	"context"
	"shared/models"
)

type MockUserService struct {
	Users         []models.UserPublic
	User          *models.UserPublic
	ErrorToReturn error
}

func (m *MockUserService) GetAll(ctx context.Context) ([]models.UserPublic, error) {
	return m.Users, m.ErrorToReturn
}

func (m *MockUserService) GetByID(ctx context.Context, id int) (*models.UserPublic, error) {
	return m.User, m.ErrorToReturn
}

func (m *MockUserService) Create(ctx context.Context, req models.CreateUserRequest) (*models.UserPublic, error) {
	return m.User, m.ErrorToReturn
}

func (m *MockUserService) Update(ctx context.Context, id int, req models.UpdateUserRequest) (*models.UserPublic, error) {
	return m.User, m.ErrorToReturn
}

func (m *MockUserService) Delete(ctx context.Context, id int) error {
	return m.ErrorToReturn
}
