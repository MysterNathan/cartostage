package services_test

import (
	"context"
	"errors"
	"shared/models"
	"stage/internal/repositories"
	"stage/internal/services"
	"testing"
	"time"
)

// --- Helpers ---

func buildStageService(t *testing.T, mockRepo *repositories.MockStageRepository) *services.StageService {
	t.Helper()
	mockFormService := &services.FormService{} // adapter si FormService a un constructeur
	return services.NewStageService(mockRepo, mockFormService)
}

func baseCreateRequest() models.CreateStageRequest {
	return models.CreateStageRequest{
		StageOfferID: 1,
		StudentID:    10,
		Status:       "pending",
		StartDate:    time.Now(),
		EndDate:      time.Now().AddDate(0, 3, 0),
	}
}

func seedStage(mock *repositories.MockStageRepository) *models.Stage {
	stage := &models.Stage{
		ID:           1,
		StageOfferID: 1,
		StudentID:    10,
		Status:       "pending",
		StartDate:    time.Now(),
		EndDate:      time.Now().AddDate(0, 3, 0),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	mock.Stages[stage.ID] = stage
	return stage
}

// --- GetStagesPublic ---

func TestGetStagesPublic_RetourneListe(t *testing.T) {
	mock := repositories.NewMockStageRepository()
	mock.PublicStages = []models.StageWithDetails{
		{Stage: models.Stage{ID: 1}},
		{Stage: models.Stage{ID: 2}},
	}
	svc := buildStageService(t, mock)

	result, err := svc.GetStagesPublic()

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
	if len(result) != 2 {
		t.Errorf("Attendu 2 stages, obtenu %d", len(result))
	}
}

func TestGetStagesPublic_ErreurRepository_PropageeErreur(t *testing.T) {
	mock := repositories.NewMockStageRepository()
	mock.ErrorToReturn = errors.New("db error")
	svc := buildStageService(t, mock)

	_, err := svc.GetStagesPublic()

	if err == nil {
		t.Fatal("Une erreur était attendue")
	}
}

// --- GetStages ---

func TestGetStages_RetourneListe(t *testing.T) {
	mock := repositories.NewMockStageRepository()
	seedStage(mock)
	svc := buildStageService(t, mock)

	result, err := svc.GetStages(context.Background())

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
	if len(result) != 1 {
		t.Errorf("Attendu 1 stage, obtenu %d", len(result))
	}
}

func TestGetStages_ErreurRepository_PropageeErreur(t *testing.T) {
	mock := repositories.NewMockStageRepository()
	mock.ErrorToReturn = errors.New("db error")
	svc := buildStageService(t, mock)

	_, err := svc.GetStages(context.Background())

	if err == nil {
		t.Fatal("Une erreur était attendue")
	}
}

// --- CreateStage ---

func TestCreateStage_CasNominal(t *testing.T) {
	mock := repositories.NewMockStageRepository()
	svc := buildStageService(t, mock)

	req := baseCreateRequest()
	result, err := svc.CreateStage(context.Background(), req)

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
	if result == nil {
		t.Fatal("Le stage créé ne devrait pas être nil")
	}
	if result.ID == 0 {
		t.Error("L'ID devrait être renseigné")
	}
	if result.StageOfferID != req.StageOfferID {
		t.Errorf("StageOfferID attendu %d, obtenu %d", req.StageOfferID, result.StageOfferID)
	}
	if result.CreatedAt.IsZero() {
		t.Error("CreatedAt ne devrait pas être zéro")
	}
	if result.UpdatedAt.IsZero() {
		t.Error("UpdatedAt ne devrait pas être zéro")
	}
}

func TestCreateStage_ErreurRepository_PropageeErreur(t *testing.T) {
	mock := repositories.NewMockStageRepository()
	mock.ErrorToReturn = errors.New("db error")
	svc := buildStageService(t, mock)

	_, err := svc.CreateStage(context.Background(), baseCreateRequest())

	if err == nil {
		t.Fatal("Une erreur était attendue")
	}
}

func TestCreateStage_FormServiceNil_RetourneErreur(t *testing.T) {
	mock := repositories.NewMockStageRepository()
	svc := services.NewStageService(mock, nil)

	_, err := svc.CreateStage(context.Background(), baseCreateRequest())

	if err == nil {
		t.Fatal("Une erreur était attendue avec un FormService nil")
	}
}

// --- UpdateStage ---

func TestUpdateStage_CasNominal(t *testing.T) {
	mock := repositories.NewMockStageRepository()
	existing := seedStage(mock)
	svc := buildStageService(t, mock)

	newStatus := "active"
	req := models.UpdateStageRequest{Status: &newStatus}

	result, err := svc.UpdateStage(context.Background(), existing.ID, req)

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
	if result.Status != newStatus {
		t.Errorf("Status attendu '%s', obtenu '%s'", newStatus, result.Status)
	}
}

func TestUpdateStage_StageInexistant_RetourneErreur(t *testing.T) {
	mock := repositories.NewMockStageRepository()
	svc := buildStageService(t, mock)

	newStatus := "active"
	req := models.UpdateStageRequest{Status: &newStatus}

	_, err := svc.UpdateStage(context.Background(), 999999, req)

	if err == nil {
		t.Fatal("Une erreur était attendue pour un stage inexistant")
	}
}

func TestUpdateStage_ErreurRepository_PropageeErreur(t *testing.T) {
	mock := repositories.NewMockStageRepository()
	existing := seedStage(mock)
	svc := buildStageService(t, mock)

	mock.ErrorToReturn = errors.New("db error")
	newStatus := "active"

	_, err := svc.UpdateStage(context.Background(), existing.ID, models.UpdateStageRequest{Status: &newStatus})

	if err == nil {
		t.Fatal("Une erreur était attendue")
	}
}

// --- DeleteStage ---

func TestDeleteStage_CasNominal(t *testing.T) {
	mock := repositories.NewMockStageRepository()
	existing := seedStage(mock)
	svc := buildStageService(t, mock)

	err := svc.DeleteStage(context.Background(), existing.ID)

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
	if _, exists := mock.Stages[existing.ID]; exists {
		t.Error("Le stage devrait avoir été supprimé")
	}
}

func TestDeleteStage_StageInexistant_RetourneErreur(t *testing.T) {
	mock := repositories.NewMockStageRepository()
	svc := buildStageService(t, mock)

	err := svc.DeleteStage(context.Background(), 999999)

	if err == nil {
		t.Fatal("Une erreur était attendue pour un stage inexistant")
	}
}

func TestDeleteStage_ErreurRepository_PropageeErreur(t *testing.T) {
	mock := repositories.NewMockStageRepository()
	existing := seedStage(mock)
	mock.ErrorToReturn = errors.New("db error")
	svc := buildStageService(t, mock)

	err := svc.DeleteStage(context.Background(), existing.ID)

	if err == nil {
		t.Fatal("Une erreur était attendue")
	}
}
