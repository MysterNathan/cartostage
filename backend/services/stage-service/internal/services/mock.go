package services

import (
	"context"
	"errors"
	"shared/models"
)

/*
----------STAGE-------------
*/
type MockStageService struct {
	Stages        map[int]*models.Stage
	PublicStages  []models.StageWithDetails
	ErrorToReturn error
	NextID        int
}

func NewMockStageService() *MockStageService {
	return &MockStageService{
		Stages: make(map[int]*models.Stage),
		NextID: 1,
	}
}

func (m *MockStageService) GetStages(ctx context.Context) ([]models.Stage, error) {
	if m.ErrorToReturn != nil {
		return nil, m.ErrorToReturn
	}
	stages := make([]models.Stage, 0, len(m.Stages))
	for _, s := range m.Stages {
		stages = append(stages, *s)
	}
	return stages, nil
}

func (m *MockStageService) GetStagesPublic() ([]models.StageWithDetails, error) {
	if m.ErrorToReturn != nil {
		return nil, m.ErrorToReturn
	}
	return m.PublicStages, nil
}

func (m *MockStageService) UpdateStage(ctx context.Context, id int, req models.UpdateStageRequest) (*models.Stage, error) {
	if m.ErrorToReturn != nil {
		return nil, m.ErrorToReturn
	}
	stage, exists := m.Stages[id]
	if !exists {
		return nil, errors.New("stage-service not found")
	}
	req.ApplyTo(stage)
	return stage, nil
}

func (m *MockStageService) CreateStage(ctx context.Context, req models.CreateStageRequest) (*models.Stage, error) {
	if m.ErrorToReturn != nil {
		return nil, m.ErrorToReturn
	}
	stage := &models.Stage{
		ID:           m.NextID,
		StageOfferID: req.StageOfferID,
		StudentID:    req.StudentID,
		TeacherID:    req.TeacherID,
		TutorID:      req.TutorID,
		Status:       req.Status,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
	}
	m.NextID++
	m.Stages[stage.ID] = stage
	return stage, nil
}

func (m *MockStageService) DeleteStage(ctx context.Context, id int) error {
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
----------FORM-------------
*/
type FormServiceMock struct {
	GetFn               func(ctx context.Context) ([]*models.FormFormSection, error)
	UpdateFormFn        func(ctx context.Context, form *models.Form) (*models.Form, error)
	UpdateFormSectionFn func(ctx context.Context, sections []models.FormSection) ([]models.FormSection, error)
}

func (m *FormServiceMock) Get(ctx context.Context) ([]*models.FormFormSection, error) {
	return m.GetFn(ctx)
}

func (m *FormServiceMock) UpdateForm(ctx context.Context, form *models.Form) (*models.Form, error) {
	return m.UpdateFormFn(ctx, form)
}

func (m *FormServiceMock) UpdateFormSection(ctx context.Context, sections []models.FormSection) ([]models.FormSection, error) {
	return m.UpdateFormSectionFn(ctx, sections)
}
