package repositories

import (
	"context"
	"errors"
	"shared/models"
)

/*
----------MOCK STAGE----------
*/
type MockStageRepository struct {
	Stages        map[int]*models.Stage
	PublicStages  []models.StageWithDetails
	ErrorToReturn error
	nextID        int
}

func NewMockStageRepository() *MockStageRepository {
	return &MockStageRepository{
		Stages: make(map[int]*models.Stage),
		nextID: 1,
	}
}

func (m *MockStageRepository) GetStages(ctx context.Context) ([]models.Stage, error) {
	if m.ErrorToReturn != nil {
		return nil, m.ErrorToReturn
	}
	stages := make([]models.Stage, 0, len(m.Stages))
	for _, s := range m.Stages {
		stages = append(stages, *s)
	}
	return stages, nil
}

func (m *MockStageRepository) GetStagesPublic() ([]models.StageWithDetails, error) {
	if m.ErrorToReturn != nil {
		return nil, m.ErrorToReturn
	}
	return m.PublicStages, nil
}

func (m *MockStageRepository) GetStageByID(ctx context.Context, id int) (*models.Stage, error) {
	if m.ErrorToReturn != nil {
		return nil, m.ErrorToReturn
	}
	stage, exists := m.Stages[id]
	if !exists {
		return nil, nil
	}
	return stage, nil
}

func (m *MockStageRepository) UpdateStage(ctx context.Context, stage *models.Stage) error {
	if m.ErrorToReturn != nil {
		return m.ErrorToReturn
	}
	if _, exists := m.Stages[stage.ID]; !exists {
		return errors.New("stage not found")
	}
	m.Stages[stage.ID] = stage
	return nil
}

func (m *MockStageRepository) CreateStage(ctx context.Context, stage *models.Stage) (*models.Stage, error) {
	if m.ErrorToReturn != nil {
		return nil, m.ErrorToReturn
	}
	stage.ID = m.nextID
	m.nextID++
	m.Stages[stage.ID] = stage
	return stage, nil
}

func (m *MockStageRepository) DeleteStage(ctx context.Context, id int) error {
	if m.ErrorToReturn != nil {
		return m.ErrorToReturn
	}
	if _, exists := m.Stages[id]; !exists {
		return errors.New("stage not found")
	}
	delete(m.Stages, id)
	return nil
}

/*
----------MOCK FORM----------
*/
type FormRepositoryMock struct {
	GetFn               func(ctx context.Context, userID int, role models.UserRole) ([]*models.FormFormSection, error)
	UpdateFormFn        func(ctx context.Context, data *models.Form, userID int) (*models.Form, error)
	UpdateFormSectionFn func(ctx context.Context, data *models.FormSection, userID int) ([]models.FormSection, error)
	CreateFormFn        func(ctx context.Context, form models.Form, sections []models.FormSection, userID int) error
}

func (m *FormRepositoryMock) Get(ctx context.Context, userID int, role models.UserRole) ([]*models.FormFormSection, error) {
	return m.GetFn(ctx, userID, role)
}

func (m *FormRepositoryMock) UpdateForm(ctx context.Context, data *models.Form, userID int) (*models.Form, error) {
	return m.UpdateFormFn(ctx, data, userID)
}

func (m *FormRepositoryMock) UpdateFormSection(ctx context.Context, data *models.FormSection, userID int) ([]models.FormSection, error) {
	return m.UpdateFormSectionFn(ctx, data, userID)
}

func (m *FormRepositoryMock) CreateForm(ctx context.Context, form models.Form, sections []models.FormSection, userID int) error {
	return m.CreateFormFn(ctx, form, sections, userID)
}
