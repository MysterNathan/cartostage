package repositories

import (
	"context"
	"fmt"
	"shared/models"
)

type MockUserRepository struct {
	Users         map[int]*models.User
	ErrorToReturn error
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		Users: make(map[int]*models.User),
	}
}

func (m *MockUserRepository) GetAll(ctx context.Context, filter models.UserFilter) ([]models.User, error) {
	if m.ErrorToReturn != nil {
		return nil, m.ErrorToReturn
	}
	users := make([]models.User, 0, len(m.Users))
	for _, u := range m.Users {
		users = append(users, *u)
	}
	return users, nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	if m.ErrorToReturn != nil {
		return nil, m.ErrorToReturn
	}
	user, ok := m.Users[id]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	if m.ErrorToReturn != nil {
		return m.ErrorToReturn
	}
	user.ID = len(m.Users) + 1
	m.Users[user.ID] = user
	return nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	if m.ErrorToReturn != nil {
		return m.ErrorToReturn
	}
	if _, ok := m.Users[user.ID]; !ok {
		return fmt.Errorf("user not found")
	}
	m.Users[user.ID] = user
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id int) error {
	if m.ErrorToReturn != nil {
		return m.ErrorToReturn
	}
	if _, ok := m.Users[id]; !ok {
		return fmt.Errorf("user not found")
	}
	delete(m.Users, id)
	return nil
}
