package services

import (
	"context"
	"enterprise/internal/repositories"
	"shared/models"
	//"shared/models"
	//"strconv"
)

type EnterpriseService struct {
	repo *repositories.EnterpriseRepository
}

func NewEnterpriseService(repo *repositories.EnterpriseRepository) *EnterpriseService {
	return &EnterpriseService{
		repo: repo,
	}
}

//// GetAll - Récupère toutes les entreprises (pour admin)
//func (s *EnterpriseService) GetAll() ([]*models.Enterprise, error) {
//	return s.repo.GetAll()
//}
//
//// GetByID - Récupère une entreprise par ID (pour admin)
//func (s *EnterpriseService) GetByID(id int) (*models.Enterprise, error) {
//	return s.repo.GetByID(id)
//}
//
//// GetMe - Récupère les infos de l'entreprise connectée
//func (s *EnterpriseService) GetMe(enterpriseID int) (*models.Enterprise, error) {
//	println("enterprise get me: " + strconv.Itoa(enterpriseID))
//	return s.repo.GetByID(enterpriseID)
//}
//
//// Create - Crée une nouvelle entreprise (pour admin)
//func (s *EnterpriseService) Create(enterprise *models.Enterprise) (*models.Enterprise, error) {
//	return s.repo.Create(enterprise)
//}
//
//// Update - Met à jour une entreprise (pour admin ou l'entreprise elle-même)
//func (s *EnterpriseService) Update(id int, enterprise *models.Enterprise) (*models.Enterprise, error) {
//	return s.repo.Update(id, enterprise)
//}
//
//// Delete - Supprime une entreprise (pour admin)
//func (s *EnterpriseService) Delete(id int) error {
//	return s.repo.Delete(id)
//}
//
//// GetWithStats - Récupère une entreprise avec ses statistiques
//func (s *EnterpriseService) GetWithStats(id int) (*models.EnterpriseWithStats, error) {
//	return s.repo.GetWithStats(id)
//}
//
//// GetStages - Récupère les stages de l'entreprise connectée
//func (s *EnterpriseService) GetStages(enterpriseID int) ([]*models.Stage, error) {
//	return s.repo.GetStagesByEnterpriseID(enterpriseID)
//}

func (s *EnterpriseService) GetStats(ctx context.Context) (*models.EnterpriseStats, error) {
	return s.repo.GetStats(ctx)
}
