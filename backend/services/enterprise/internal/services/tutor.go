package services

import (
	"enterprise/internal/repositories"
	"shared/models"
)

type TutorService struct {
	repo *repositories.TutorRepository
}

func NewTutorService(repo *repositories.TutorRepository) *TutorService {
	return &TutorService{repo: repo}
}

func (s *TutorService) GetAll() ([]models.TutorWithEnterprise, error) {
	return s.repo.GetAll()
}

func (s *TutorService) GetByID(id int) (*models.TutorWithEnterprise, error) {
	return s.repo.GetByID(id)
}

func (s *TutorService) Create(tutor *models.Tutor) error {
	return s.repo.Create(tutor)
}

func (s *TutorService) Update(tutor *models.Tutor) error {
	return s.repo.Update(tutor)
}

func (s *TutorService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *TutorService) GetByEnterprise(enterpriseID int) ([]models.TutorWithEnterprise, error) {
	return s.repo.GetByEnterprise(enterpriseID)
}

func (s *TutorService) GetWithStats(id int) (*models.TutorWithStats, error) {
	return s.repo.GetWithStats(id)
}
