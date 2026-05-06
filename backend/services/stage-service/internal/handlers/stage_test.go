package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	sharedContext "shared/context"
	"shared/models"
	"stage/internal/handlers"
	"stage/internal/services"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

// --- Helpers ---

func buildHandler(t *testing.T, mockSvc *services.MockStageService) *handlers.StageHandler {
	t.Helper()
	return handlers.NewStageHandler(mockSvc)
}

func buildAuthContext() context.Context {
	claims := &models.CustomClaims{
		UserID: 1,
		Role:   "admin",
	}
	return sharedContext.SetUserClaims(context.Background(), claims)
}

func seedStage(mock *services.MockStageService) *models.Stage {
	now := time.Now()
	stage := &models.Stage{
		ID:           mock.NextID,
		StageOfferID: 1,
		StudentID:    1,
		Status:       "pending",
		StartDate:    now,
		EndDate:      now.AddDate(0, 3, 0),
	}
	mock.Stages[stage.ID] = stage
	mock.NextID++
	return stage
}

func newRequestWithVars(method, url string, body []byte, vars map[string]string, ctx context.Context) *http.Request {
	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, url, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, url, nil)
	}
	req = req.WithContext(ctx)
	req = mux.SetURLVars(req, vars)
	return req
}

func validCreatePayload(t *testing.T) []byte {
	t.Helper()
	now := time.Now()
	payload := models.CreateStageRequest{
		StageOfferID: 1,
		StudentID:    1,
		Status:       "pending",
		StartDate:    now,
		EndDate:      now.AddDate(0, 3, 0),
	}
	b, _ := json.Marshal(payload)
	return b
}

// --- GetStages ---

func TestGetStages_CasNominal(t *testing.T) {
	mock := services.NewMockStageService()
	seedStage(mock)
	h := buildHandler(t, mock)

	req := httptest.NewRequest(http.MethodGet, "/stages", nil).WithContext(buildAuthContext())
	rr := httptest.NewRecorder()

	h.GetStages(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Status attendu %d, obtenu %d", http.StatusOK, rr.Code)
	}
	var result []models.Stage
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("Impossible de décoder la réponse : %v", err)
	}
	if len(result) != 1 {
		t.Errorf("1 stage attendu, obtenu %d", len(result))
	}
}

func TestGetStages_SansAuthentification_RetourneErreur(t *testing.T) {
	mock := services.NewMockStageService()
	h := buildHandler(t, mock)

	req := httptest.NewRequest(http.MethodGet, "/stages", nil)
	rr := httptest.NewRecorder()

	h.GetStages(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Status attendu %d, obtenu %d", http.StatusInternalServerError, rr.Code)
	}
}

func TestGetStages_ErreurService_RetourneErreur(t *testing.T) {
	mock := services.NewMockStageService()
	mock.ErrorToReturn = errors.New("db error")
	h := buildHandler(t, mock)

	req := httptest.NewRequest(http.MethodGet, "/stages", nil).WithContext(buildAuthContext())
	rr := httptest.NewRecorder()

	h.GetStages(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Status attendu %d, obtenu %d", http.StatusInternalServerError, rr.Code)
	}
}

// --- GetStagesPublic ---

func TestGetStagesPublic_CasNominal(t *testing.T) {
	mock := services.NewMockStageService()
	mock.PublicStages = []models.StageWithDetails{{Stage: *seedStage(mock)}}
	h := buildHandler(t, mock)

	req := httptest.NewRequest(http.MethodGet, "/stages/public", nil)
	rr := httptest.NewRecorder()

	h.GetStagesPublic(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Status attendu %d, obtenu %d", http.StatusOK, rr.Code)
	}
}

func TestGetStagesPublic_ErreurService_RetourneErreur(t *testing.T) {
	mock := services.NewMockStageService()
	mock.ErrorToReturn = errors.New("db error")
	h := buildHandler(t, mock)

	req := httptest.NewRequest(http.MethodGet, "/stages/public", nil)
	rr := httptest.NewRecorder()

	h.GetStagesPublic(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Status attendu %d, obtenu %d", http.StatusInternalServerError, rr.Code)
	}
}

// --- UpdateStage ---

func TestUpdateStage_CasNominal(t *testing.T) {
	mock := services.NewMockStageService()
	existing := seedStage(mock)
	h := buildHandler(t, mock)

	newStatus := "active"
	payload, _ := json.Marshal(models.UpdateStageRequest{Status: &newStatus})
	req := newRequestWithVars(http.MethodPut, "/stages/1", payload,
		map[string]string{"id": "1"}, buildAuthContext())
	rr := httptest.NewRecorder()

	h.UpdateStage(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Status attendu %d, obtenu %d — body: %s", http.StatusOK, rr.Code, rr.Body.String())
	}
	_ = existing
}

func TestUpdateStage_IDManquant_RetourneBadRequest(t *testing.T) {
	mock := services.NewMockStageService()
	h := buildHandler(t, mock)

	newStatus := "active"
	payload, _ := json.Marshal(models.UpdateStageRequest{Status: &newStatus})
	req := newRequestWithVars(http.MethodPut, "/stages/", payload,
		map[string]string{"id": ""}, buildAuthContext())
	rr := httptest.NewRecorder()

	h.UpdateStage(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Status attendu %d, obtenu %d", http.StatusBadRequest, rr.Code)
	}
}

func TestUpdateStage_IDInvalide_RetourneBadRequest(t *testing.T) {
	mock := services.NewMockStageService()
	h := buildHandler(t, mock)

	newStatus := "active"
	payload, _ := json.Marshal(models.UpdateStageRequest{Status: &newStatus})
	req := newRequestWithVars(http.MethodPut, "/stages/abc", payload,
		map[string]string{"id": "abc"}, buildAuthContext())
	rr := httptest.NewRecorder()

	h.UpdateStage(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Status attendu %d, obtenu %d", http.StatusBadRequest, rr.Code)
	}
}

func TestUpdateStage_AucunChamp_RetourneBadRequest(t *testing.T) {
	mock := services.NewMockStageService()
	h := buildHandler(t, mock)

	payload, _ := json.Marshal(models.UpdateStageRequest{})
	req := newRequestWithVars(http.MethodPut, "/stages/1", payload,
		map[string]string{"id": "1"}, buildAuthContext())
	rr := httptest.NewRecorder()

	h.UpdateStage(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Status attendu %d, obtenu %d", http.StatusBadRequest, rr.Code)
	}
}

func TestUpdateStage_StageInexistant_RetourneNotFound(t *testing.T) {
	mock := services.NewMockStageService()
	h := buildHandler(t, mock)

	newStatus := "active"
	payload, _ := json.Marshal(models.UpdateStageRequest{Status: &newStatus})
	req := newRequestWithVars(http.MethodPut, "/stages/999", payload,
		map[string]string{"id": "999"}, buildAuthContext())
	rr := httptest.NewRecorder()

	h.UpdateStage(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Status attendu %d, obtenu %d", http.StatusNotFound, rr.Code)
	}
}

// --- CreateStage ---

func TestCreateStage_CasNominal(t *testing.T) {
	mock := services.NewMockStageService()
	h := buildHandler(t, mock)

	req := newRequestWithVars(http.MethodPost, "/stages", validCreatePayload(t),
		map[string]string{}, buildAuthContext())
	rr := httptest.NewRecorder()

	h.CreateStage(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Status attendu %d, obtenu %d — body: %s", http.StatusCreated, rr.Code, rr.Body.String())
	}
}

func TestCreateStage_ChampsObligatoiresManquants_RetourneBadRequest(t *testing.T) {
	tests := []struct {
		name    string
		payload models.CreateStageRequest
	}{
		{
			"stage_offer_id manquant",
			models.CreateStageRequest{StudentID: 1, Status: "pending", StartDate: time.Now(), EndDate: time.Now().AddDate(0, 3, 0)},
		},
		{
			"student_id manquant",
			models.CreateStageRequest{StageOfferID: 1, Status: "pending", StartDate: time.Now(), EndDate: time.Now().AddDate(0, 3, 0)},
		},
		{
			"status manquant",
			models.CreateStageRequest{StageOfferID: 1, StudentID: 1, StartDate: time.Now(), EndDate: time.Now().AddDate(0, 3, 0)},
		},
		{
			"end_date avant start_date",
			models.CreateStageRequest{StageOfferID: 1, StudentID: 1, Status: "pending", StartDate: time.Now().AddDate(0, 3, 0), EndDate: time.Now()},
		},
	}

	mock := services.NewMockStageService()
	h := buildHandler(t, mock)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, _ := json.Marshal(tt.payload)
			req := newRequestWithVars(http.MethodPost, "/stages", payload,
				map[string]string{}, buildAuthContext())
			rr := httptest.NewRecorder()

			h.CreateStage(rr, req)

			if rr.Code != http.StatusBadRequest {
				t.Errorf("Cas '%s' : status attendu %d, obtenu %d", tt.name, http.StatusBadRequest, rr.Code)
			}
		})
	}
}

func TestCreateStage_JSONInvalide_RetourneBadRequest(t *testing.T) {
	mock := services.NewMockStageService()
	h := buildHandler(t, mock)

	req := newRequestWithVars(http.MethodPost, "/stages", []byte(`{invalid}`),
		map[string]string{}, buildAuthContext())
	rr := httptest.NewRecorder()

	h.CreateStage(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Status attendu %d, obtenu %d", http.StatusBadRequest, rr.Code)
	}
}

func TestCreateStage_Duplicate_RetourneConflict(t *testing.T) {
	mock := services.NewMockStageService()
	mock.ErrorToReturn = errors.New("duplicate key")
	h := buildHandler(t, mock)

	req := newRequestWithVars(http.MethodPost, "/stages", validCreatePayload(t),
		map[string]string{}, buildAuthContext())
	rr := httptest.NewRecorder()

	h.CreateStage(rr, req)

	if rr.Code != http.StatusConflict {
		t.Errorf("Status attendu %d, obtenu %d", http.StatusConflict, rr.Code)
	}
}

func TestCreateStage_ForeignKey_RetourneBadRequest(t *testing.T) {
	mock := services.NewMockStageService()
	mock.ErrorToReturn = errors.New("foreign key violation")
	h := buildHandler(t, mock)

	req := newRequestWithVars(http.MethodPost, "/stages", validCreatePayload(t),
		map[string]string{}, buildAuthContext())
	rr := httptest.NewRecorder()

	h.CreateStage(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Status attendu %d, obtenu %d", http.StatusBadRequest, rr.Code)
	}
}

// --- DeleteStage ---

func TestDeleteStage_CasNominal(t *testing.T) {
	mock := services.NewMockStageService()
	seedStage(mock)
	h := buildHandler(t, mock)

	req := newRequestWithVars(http.MethodDelete, "/stages/1", nil,
		map[string]string{"id": "1"}, buildAuthContext())
	rr := httptest.NewRecorder()

	h.DeleteStage(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Status attendu %d, obtenu %d", http.StatusOK, rr.Code)
	}
}

func TestDeleteStage_IDInvalide_RetourneBadRequest(t *testing.T) {
	mock := services.NewMockStageService()
	h := buildHandler(t, mock)

	req := newRequestWithVars(http.MethodDelete, "/stages/abc", nil,
		map[string]string{"id": "abc"}, buildAuthContext())
	rr := httptest.NewRecorder()

	h.DeleteStage(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Status attendu %d, obtenu %d", http.StatusBadRequest, rr.Code)
	}
}

func TestDeleteStage_StageInexistant_RetourneNotFound(t *testing.T) {
	mock := services.NewMockStageService()
	h := buildHandler(t, mock)

	req := newRequestWithVars(http.MethodDelete, "/stages/999", nil,
		map[string]string{"id": "999"}, buildAuthContext())
	rr := httptest.NewRecorder()

	h.DeleteStage(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Status attendu %d, obtenu %d", http.StatusNotFound, rr.Code)
	}
}

func TestDeleteStage_ErreurForeignKey_RetourneConflict(t *testing.T) {
	mock := services.NewMockStageService()
	seedStage(mock)
	mock.ErrorToReturn = errors.New("foreign key constraint")
	h := buildHandler(t, mock)

	req := newRequestWithVars(http.MethodDelete, "/stages/1", nil,
		map[string]string{"id": "1"}, buildAuthContext())
	rr := httptest.NewRecorder()

	h.DeleteStage(rr, req)

	if rr.Code != http.StatusConflict {
		t.Errorf("Status attendu %d, obtenu %d", http.StatusConflict, rr.Code)
	}
}
