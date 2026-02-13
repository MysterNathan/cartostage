package services

import (
	"context"
	"fmt"
	sharedContext "shared/context"
	"shared/models"
	"stage/internal/repositories"
)

type FormService struct {
	formRepository *repositories.FormRepository
}

func NewFormService(repository *repositories.FormRepository) *FormService {
	return &FormService{formRepository: repository}
}

func (s FormService) Get(ctx context.Context) ([]models.Form, error) {
	return s.formRepository.Get(ctx)
}

func (s FormService) UpdateForm(ctx context.Context, data models.Form, formId int) ([]models.Form, error) {
	claims := sharedContext.GetClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no claims found in context")
	}
	return s.formRepository.UpdateForm(ctx, data, claims.UserID, formId)
}

func (s FormService) CreateForm(ctx context.Context, data models.Form) (*models.Form, error) {
	claims := sharedContext.GetClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no claims found in context")
	}
	if claims.UserID != *data.TeacherID {
		return nil, fmt.Errorf("invalid teacher ID")
	}
	return s.formRepository.CreateForm(ctx, data, claims.UserID)
}
