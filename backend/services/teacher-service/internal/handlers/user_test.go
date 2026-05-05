package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"shared/models"
	"shared/test"
	"teacher/internal/handlers"
	"teacher/internal/services"
	"testing"
)

// --- Helpers ---

func buildTeacherHandler(svc *services.MockUserService) *handlers.UserHandler {
	return handlers.NewUserHandler(svc)
}

// --- GetAll ---

func TestGetAll_Success(t *testing.T) {
	svc := &services.MockUserService{
		Users: []models.UserPublic{
			{ID: 1, Username: "alice"},
			{ID: 2, Username: "bob"},
		},
	}
	handler := buildTeacherHandler(svc)

	req := test.MakeRequest(t, http.MethodGet, "/users", nil)
	recorder := httptest.NewRecorder()

	handler.GetAll(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Status attendu %d, obtenu %d", http.StatusOK, recorder.Code)
	}

	var users []models.UserPublic
	if err := json.NewDecoder(recorder.Body).Decode(&users); err != nil {
		t.Fatalf("Impossible de décoder la réponse : %v", err)
	}
	if len(users) != 2 {
		t.Errorf("Nombre d'utilisateurs attendu 2, obtenu %d", len(users))
	}
}

func TestGetAll_ServiceError(t *testing.T) {
	svc := &services.MockUserService{
		ErrorToReturn: errors.New("database error"),
	}
	handler := buildTeacherHandler(svc)

	req := test.MakeRequest(t, http.MethodGet, "/users", nil)
	recorder := httptest.NewRecorder()

	handler.GetAll(recorder, req)

	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("Status attendu %d, obtenu %d", http.StatusInternalServerError, recorder.Code)
	}
}

// --- GetByID ---

func TestGetByID_Success(t *testing.T) {
	svc := &services.MockUserService{
		User: &models.UserPublic{ID: 1, Username: "alice"},
	}
	handler := buildTeacherHandler(svc)

	req := test.MakeRequest(t, http.MethodGet, "/users/1", nil)
	req = test.RequestWithVars(req, map[string]string{"id": "1"})
	recorder := httptest.NewRecorder()

	handler.GetByID(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Status attendu %d, obtenu %d", http.StatusOK, recorder.Code)
	}

	var user models.UserPublic
	if err := json.NewDecoder(recorder.Body).Decode(&user); err != nil {
		t.Fatalf("Impossible de décoder la réponse : %v", err)
	}
	if user.Username != "alice" {
		t.Errorf("Username attendu 'alice', obtenu '%s'", user.Username)
	}
}

func TestGetByID_InvalidID(t *testing.T) {
	svc := &services.MockUserService{}
	handler := buildTeacherHandler(svc)

	req := test.MakeRequest(t, http.MethodGet, "/users/abc", nil)
	req = test.RequestWithVars(req, map[string]string{"id": "abc"})
	recorder := httptest.NewRecorder()

	handler.GetByID(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Status attendu %d, obtenu %d", http.StatusBadRequest, recorder.Code)
	}
}

func TestGetByID_NotFound(t *testing.T) {
	svc := &services.MockUserService{
		ErrorToReturn: errors.New("user not found"),
	}
	handler := buildTeacherHandler(svc)

	req := test.MakeRequest(t, http.MethodGet, "/users/999", nil)
	req = test.RequestWithVars(req, map[string]string{"id": "999"})
	recorder := httptest.NewRecorder()

	handler.GetByID(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("Status attendu %d, obtenu %d", http.StatusNotFound, recorder.Code)
	}
}

// --- Create ---

func TestCreate_Success(t *testing.T) {
	svc := &services.MockUserService{
		User: &models.UserPublic{ID: 1, Username: "alice"},
	}
	handler := buildTeacherHandler(svc)

	body := models.CreateUserRequest{
		Username: "alice",
		Password: "secret",
	}
	req := test.MakeRequest(t, http.MethodPost, "/users", body)
	recorder := httptest.NewRecorder()

	handler.Create(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Errorf("Status attendu %d, obtenu %d", http.StatusCreated, recorder.Code)
	}

	var user models.UserPublic
	if err := json.NewDecoder(recorder.Body).Decode(&user); err != nil {
		t.Fatalf("Impossible de décoder la réponse : %v", err)
	}
	if user.Username != "alice" {
		t.Errorf("Username attendu 'alice', obtenu '%s'", user.Username)
	}
}

func TestCreate_InvalidJSON(t *testing.T) {
	svc := &services.MockUserService{}
	handler := buildTeacherHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString("invalid json"))
	recorder := httptest.NewRecorder()

	handler.Create(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Status attendu %d, obtenu %d", http.StatusBadRequest, recorder.Code)
	}
}

func TestCreate_ServiceError(t *testing.T) {
	svc := &services.MockUserService{
		ErrorToReturn: errors.New("creation failed"),
	}
	handler := buildTeacherHandler(svc)

	req := test.MakeRequest(t, http.MethodPost, "/users", models.CreateUserRequest{Username: "alice"})
	recorder := httptest.NewRecorder()

	handler.Create(recorder, req)

	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("Status attendu %d, obtenu %d", http.StatusInternalServerError, recorder.Code)
	}
}

// --- Update ---

func TestUpdate_Success(t *testing.T) {
	svc := &services.MockUserService{
		User: &models.UserPublic{ID: 1, Username: "alice-updated"},
	}
	handler := buildTeacherHandler(svc)

	req := test.MakeRequest(t, http.MethodPut, "/users/1", models.UpdateUserRequest{})
	req = test.RequestWithVars(req, map[string]string{"id": "1"})
	recorder := httptest.NewRecorder()

	handler.Update(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Status attendu %d, obtenu %d", http.StatusOK, recorder.Code)
	}
}

func TestUpdate_InvalidID(t *testing.T) {
	svc := &services.MockUserService{}
	handler := buildTeacherHandler(svc)

	req := test.MakeRequest(t, http.MethodPut, "/users/abc", models.UpdateUserRequest{})
	req = test.RequestWithVars(req, map[string]string{"id": "abc"})
	recorder := httptest.NewRecorder()

	handler.Update(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Status attendu %d, obtenu %d", http.StatusBadRequest, recorder.Code)
	}
}

func TestUpdate_ServiceError(t *testing.T) {
	svc := &services.MockUserService{
		ErrorToReturn: errors.New("update failed"),
	}
	handler := buildTeacherHandler(svc)

	req := test.MakeRequest(t, http.MethodPut, "/users/1", models.UpdateUserRequest{})
	req = test.RequestWithVars(req, map[string]string{"id": "1"})
	recorder := httptest.NewRecorder()

	handler.Update(recorder, req)

	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("Status attendu %d, obtenu %d", http.StatusInternalServerError, recorder.Code)
	}
}

// --- Delete ---

func TestDelete_Success(t *testing.T) {
	svc := &services.MockUserService{}
	handler := buildTeacherHandler(svc)

	req := test.MakeRequest(t, http.MethodDelete, "/users/1", nil)
	req = test.RequestWithVars(req, map[string]string{"id": "1"})
	recorder := httptest.NewRecorder()

	handler.Delete(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Status attendu %d, obtenu %d", http.StatusOK, recorder.Code)
	}

	var response map[string]string
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("Impossible de décoder la réponse : %v", err)
	}
	if response["message"] != "user deleted successfully" {
		t.Errorf("Message inattendu : %s", response["message"])
	}
}

func TestDelete_InvalidID(t *testing.T) {
	svc := &services.MockUserService{}
	handler := buildTeacherHandler(svc)

	req := test.MakeRequest(t, http.MethodDelete, "/users/abc", nil)
	req = test.RequestWithVars(req, map[string]string{"id": "abc"})
	recorder := httptest.NewRecorder()

	handler.Delete(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Status attendu %d, obtenu %d", http.StatusBadRequest, recorder.Code)
	}
}

func TestDelete_ServiceError(t *testing.T) {
	svc := &services.MockUserService{
		ErrorToReturn: errors.New("delete failed"),
	}
	handler := buildTeacherHandler(svc)

	req := test.MakeRequest(t, http.MethodDelete, "/users/1", nil)
	req = test.RequestWithVars(req, map[string]string{"id": "1"})
	recorder := httptest.NewRecorder()

	handler.Delete(recorder, req)

	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("Status attendu %d, obtenu %d", http.StatusInternalServerError, recorder.Code)
	}
}
