package services

import (
	"enterprise/internal/repositories"
	"shared/models"
)

type EnterpriseService struct {
	repo *repositories.EnterpriseRepository
}

func NewEnterpriseService(repo *repositories.EnterpriseRepository) *EnterpriseService {
	return &EnterpriseService{repo: repo}
}

func (s *EnterpriseService) GetAll() ([]models.Enterprise, error) {
	return s.repo.GetAll()
}

func (s *EnterpriseService) GetByID(id int) (*models.Enterprise, error) {
	return s.repo.GetByID(id)
}

func (s *EnterpriseService) Create(enterprise *models.Enterprise) error {
	return s.repo.Create(enterprise)
}

func (s *EnterpriseService) Update(enterprise *models.Enterprise) error {
	return s.repo.Update(enterprise)
}

func (s *EnterpriseService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *EnterpriseService) GetWithStats(id int) (*models.EnterpriseWithStats, error) {
	return s.repo.GetWithStats(id)
}
