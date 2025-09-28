package services

import (
	"context"
	"fmt"
	"shared/models"
	"stage/internal/repositories"
	"time"
)

type StageService struct {
	stageRepo *repositories.StageRepository
}

func NewStageService(stageRepo *repositories.StageRepository) *StageService {
	return &StageService{
		stageRepo: stageRepo,
	}
}

func (s *StageService) GetStagesPublic() ([]models.Stage, error) {
	return s.stageRepo.GetStagesPublic()
}

func (s *StageService) GetStages(ctx context.Context) ([]models.Stage, error) {
	return s.stageRepo.GetStages(ctx)
}

func (s *StageService) UpdateStage(ctx context.Context, stageID int, updateReq models.UpdateStageRequest) (*models.Stage, error) {

	// Récupérer le stage existant
	existingStage, err := s.stageRepo.GetStageByID(ctx, stageID)
	if err != nil {
		return nil, err
	}
	if existingStage == nil {
		return nil, fmt.Errorf("stage not found")
	}

	// Appliquer les modifications
	updateReq.ApplyTo(existingStage)

	// Sauvegarder en base
	if err := s.stageRepo.UpdateStage(ctx, existingStage); err != nil {
		return nil, fmt.Errorf("erreur lors de la mise à jour: %v", err)
	}

	return existingStage, nil
}
func (s *StageService) CreateStage(ctx context.Context, createReq models.CreateStageRequest) (*models.Stage, error) {
	// Créer le stage à partir de la requête
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
		return nil, fmt.Errorf("erreur lors de la création du stage: %v", err)
	}

	return createdStage, nil
}

func (s *StageService) DeleteStage(ctx context.Context, id int) error {
	// Vérifier d'abord si le stage existe
	_, err := s.stageRepo.GetStageByID(ctx, id)
	if err != nil {
		return fmt.Errorf("stage not found: %v", err)
	}

	// Supprimer le stage
	err = s.stageRepo.DeleteStage(ctx, id)
	if err != nil {
		return fmt.Errorf("erreur lors de la suppression du stage: %v", err)
	}

	return nil
}
