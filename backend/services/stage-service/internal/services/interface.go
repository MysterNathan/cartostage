package services

import (
	"context"
	"shared/models"
)

type StageServiceInterface interface {
	GetStages(context context.Context) ([]models.Stage, error)
	GetStagesPublic() ([]models.StageWithDetails, error)
	UpdateStage(ctx context.Context, id int, req models.UpdateStageRequest) (*models.Stage, error)
	CreateStage(ctx context.Context, req models.CreateStageRequest) (*models.Stage, error)
	DeleteStage(ctx context.Context, id int) error
}

type FormServiceInterface interface {
	Get(ctx context.Context) ([]*models.FormFormSection, error)
	UpdateForm(ctx context.Context, form *models.Form) (*models.Form, error)
	UpdateFormSection(ctx context.Context, sections []models.FormSection) ([]models.FormSection, error)
}
