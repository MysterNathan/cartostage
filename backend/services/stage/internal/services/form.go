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

func (s FormService) Get(ctx context.Context) (*models.FormFormSection, error) {
	claims := sharedContext.GetClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no claims found in context")
	}
	return s.formRepository.Get(ctx, claims.UserID)
}

func (s FormService) UpdateForm(ctx context.Context, data models.Form, formId int) ([]models.Form, error) {
	claims := sharedContext.GetClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no claims found in context")
	}
	return s.formRepository.UpdateForm(ctx, data, claims.UserID, formId)
}

func (s FormService) CreateForm(ctx context.Context, data models.FormCreationData) error {
	claims := sharedContext.GetClaimsFromContext(ctx)
	if claims == nil {
		return fmt.Errorf("no claims found in context")
	}
	userID := claims.UserID
	dataForm := models.Form{
		StageID:   data.StageID,
		StudentID: data.StudentID,
		TeacherID: data.TeacherID,
		TutorID:   data.TutorID,
		Status:    "CREATED",
	}
	dataFormSection := []models.FormSection{
		{SectionType: "STUDENT", UserID: data.StudentID, Status: "CREATED"},
	}

	if data.TeacherID != nil {
		dataFormSection = append(dataFormSection, models.FormSection{
			SectionType: "TEACHER", UserID: *data.TeacherID, Status: "CREATED",
		})
	}

	if data.TutorID != nil {
		dataFormSection = append(dataFormSection, models.FormSection{
			SectionType: "TUTOR", UserID: *data.TutorID, Status: "CREATED",
		})
	}

	return s.formRepository.CreateForm(ctx, dataForm, dataFormSection, userID)
}
