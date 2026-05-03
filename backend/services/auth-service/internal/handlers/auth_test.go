package handlers_test

import (
	"auth/internal/handlers"
	"auth/internal/services"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"shared/models"
	"testing"
)

func makeLoginRequest(t *testing.T, body any) *http.Request {
	t.Helper()
	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("Impossible de sérialiser le body : %v", err)
	}
	return httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
}

func TestLogin_Handler_Success(t *testing.T) {
	fakeService := &services.FakeAuthService{
		ResponseToReturn: &services.LoginResponse{
			Token: "un.jwt.token",
			User:  models.User{Username: "alice"},
		},
	}
	handler := handlers.NewAuthHandler(fakeService)

	req := makeLoginRequest(t, map[string]string{
		"username": "alice",
		"password": "monMotDePasse",
	})
	recorder := httptest.NewRecorder()

	handler.Login(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Status attendu %d, obtenu %d", http.StatusOK, recorder.Code)
	}

	var response services.LoginResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("Impossible de décoder la réponse : %v", err)
	}

	if response.Token != "un.jwt.token" {
		t.Errorf("Token inattendu : %s", response.Token)
	}
}

func TestLogin_Handler_InvalidCredentials(t *testing.T) {
	fakeService := &services.FakeAuthService{
		ErrorToReturn: errors.New("invalid credentials"),
	}
	handler := handlers.NewAuthHandler(fakeService)

	req := makeLoginRequest(t, map[string]string{
		"username": "alice",
		"password": "mauvais",
	})
	recorder := httptest.NewRecorder()

	handler.Login(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Errorf("Status attendu %d, obtenu %d", http.StatusUnauthorized, recorder.Code)
	}
}

func TestLogin_Handler_InvalidJSON(t *testing.T) {
	fakeService := &services.FakeAuthService{}
	handler := handlers.NewAuthHandler(fakeService)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString("ceci n'est pas du json"))
	recorder := httptest.NewRecorder()

	handler.Login(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Status attendu %d, obtenu %d", http.StatusBadRequest, recorder.Code)
	}
}
