package repositories

import (
	"context"
	"shared/models"
)

type StageRepositoryInterface interface {
	GetStages(ctx context.Context) ([]models.Stage, error)
	GetStagesPublic() ([]models.StageWithDetails, error)
	GetStageByID(context.Context, int) (*models.Stage, error)
	UpdateStage(context.Context, *models.Stage) error
	CreateStage(context.Context, *models.Stage) (*models.Stage, error)
	DeleteStage(context.Context, int) error
}

type FormRepositoryInterface interface {
	Get(context.Context, int, models.UserRole) ([]*models.FormFormSection, error)
	UpdateForm(context.Context, *models.Form, int) (*models.Form, error)
	UpdateFormSection(context.Context, *models.FormSection, int) ([]models.FormSection, error)
	CreateForm(context.Context, models.Form, []models.FormSection, int) error
}
