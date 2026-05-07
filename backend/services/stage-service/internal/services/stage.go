package services

import (
	"context"
	"fmt"
	"shared/models"
	"stage/internal/repositories"
	"time"
)

type StageService struct {
	stageRepo   repositories.StageRepositoryInterface
	formService *FormService
}

func NewStageService(stageRepo repositories.StageRepositoryInterface, formService *FormService) *StageService {
	return &StageService{
		stageRepo:   stageRepo,
		formService: formService,
	}
}

func (s *StageService) GetStagesPublic() ([]models.StageWithDetails, error) {
	return s.stageRepo.GetStagesPublic()
}

func (s *StageService) GetStages(ctx context.Context) ([]models.Stage, error) {
	return s.stageRepo.GetStages(ctx)
}

func (s *StageService) UpdateStage(ctx context.Context, stageID int, updateReq models.UpdateStageRequest) (*models.Stage, error) {

	// Récupérer le stage-service existant
	existingStage, err := s.stageRepo.GetStageByID(ctx, stageID)
	if err != nil {
		return nil, err
	}
	if existingStage == nil {
		return nil, fmt.Errorf("stage-service not found")
	}

	// Appliquer les modifications
	updateReq.ApplyTo(existingStage)

	// Sauvegarder en base
	if err := s.stageRepo.UpdateStage(ctx, existingStage); err != nil {
		return nil, fmt.Errorf("error while updating: %v", err)
	}

	return existingStage, nil
}
func (s *StageService) CreateStage(ctx context.Context, createReq models.CreateStageRequest) (*models.Stage, error) {
	// Créer le stage-service à partir de la requête
	stage := &models.Stage{
		StageOfferID:    createReq.StageOfferID,
		StudentID:       createReq.StudentID,
		TeacherID:       createReq.TeacherID,
		TutorID:         createReq.TutorID,
		EstablishmentID: createReq.EstablishmentID,
		ContentID:       createReq.ContentID,
		Status:          createReq.Status,
		StartDate:       createReq.StartDate,
		EndDate:         createReq.EndDate,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Utiliser le repository pour sauvegarder
	createdStage, err := s.stageRepo.CreateStage(ctx, stage)
	if err != nil {
		return nil, fmt.Errorf("error when trying to create stage-service: %v", err)
	}

	formData := models.FormCreationData{
		StageID:   createdStage.ID,
		StudentID: createReq.StudentID,
		TeacherID: createReq.TeacherID,
		TutorID:   createReq.TutorID,
	}
	if s.formService == nil {
		return nil, fmt.Errorf("formService is nil")
	}
	err = s.formService.CreateForm(ctx, formData)

	if err != nil {
		return nil, fmt.Errorf("error when trying to create formular stage-service: %v", err)
	}

	return createdStage, nil
}

func (s *StageService) DeleteStage(ctx context.Context, id int) error {
	// Vérifier d'abord si le stage-service existe
	_, err := s.stageRepo.GetStageByID(ctx, id)
	if err != nil {
		return fmt.Errorf("stage-service not found: %v", err)
	}

	// Supprimer le stage-service
	err = s.stageRepo.DeleteStage(ctx, id)
	if err != nil {
		return fmt.Errorf("error when trying to delete stage-service: %v", err)
	}

	return nil
}
