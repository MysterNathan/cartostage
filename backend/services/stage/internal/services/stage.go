package services

import (
	"context"
	"shared/models"
	"stage/internal/repositories"
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
