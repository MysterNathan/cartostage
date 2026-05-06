package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"shared/models"
	"stage/internal/handlers"
	servicemock "stage/internal/services"
	"testing"
)

// --- Get ---

func TestFormGet_ServiceEchoue_Retourne500(t *testing.T) {
	svc := &servicemock.FormServiceMock{
		GetFn: func(ctx context.Context) ([]*models.FormFormSection, error) {
			return nil, fmt.Errorf("db error")
		},
	}
	h := handlers.NewFormHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/form", nil)
	w := httptest.NewRecorder()

	h.Get(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Attendu 500, obtenu %d", w.Code)
	}
}

func TestFormGet_RetourneNil_Retourne404(t *testing.T) {
	svc := &servicemock.FormServiceMock{
		GetFn: func(ctx context.Context) ([]*models.FormFormSection, error) {
			return nil, nil
		},
	}
	h := handlers.NewFormHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/form", nil)
	w := httptest.NewRecorder()

	h.Get(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Attendu 404, obtenu %d", w.Code)
	}
}

func TestFormGet_Success_Retourne200(t *testing.T) {
	expected := []*models.FormFormSection{
		{Form: &models.Form{Status: "CREATED"}},
	}
	svc := &servicemock.FormServiceMock{
		GetFn: func(ctx context.Context) ([]*models.FormFormSection, error) {
			return expected, nil
		},
	}
	h := handlers.NewFormHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/form", nil)
	w := httptest.NewRecorder()

	h.Get(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Attendu 200, obtenu %d", w.Code)
	}
	if w.Header().Get("Content-Type") != "application/json" {
		t.Error("Content-Type devrait être application/json")
	}
}

// --- UpdateForm ---

func TestUpdateForm_JSONInvalide_Retourne400(t *testing.T) {
	svc := &servicemock.FormServiceMock{}
	h := handlers.NewFormHandler(svc)

	req := httptest.NewRequest(http.MethodPut, "/form", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	h.UpdateForm(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Attendu 400, obtenu %d", w.Code)
	}
}

func TestUpdateForm_ServiceUpdateFormEchoue_Retourne500(t *testing.T) {
	svc := &servicemock.FormServiceMock{
		UpdateFormFn: func(ctx context.Context, form *models.Form) (*models.Form, error) {
			return nil, fmt.Errorf("db error")
		},
	}
	h := handlers.NewFormHandler(svc)

	body := models.FormFormSection{Form: &models.Form{Status: "CREATED"}}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, "/form", bytes.NewBuffer(b))
	w := httptest.NewRecorder()

	h.UpdateForm(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Attendu 500, obtenu %d", w.Code)
	}
}

func TestUpdateForm_ServiceUpdateSectionEchoue_Retourne500(t *testing.T) {
	svc := &servicemock.FormServiceMock{
		UpdateFormFn: func(ctx context.Context, form *models.Form) (*models.Form, error) {
			return &models.Form{Status: "CREATED"}, nil
		},
		UpdateFormSectionFn: func(ctx context.Context, sections []models.FormSection) ([]models.FormSection, error) {
			return nil, fmt.Errorf("section error")
		},
	}
	h := handlers.NewFormHandler(svc)

	body := models.FormFormSection{
		Form:         &models.Form{Status: "CREATED"},
		FormSections: []models.FormSection{{SectionType: "STUDENT"}},
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, "/form", bytes.NewBuffer(b))
	w := httptest.NewRecorder()

	h.UpdateForm(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Attendu 500, obtenu %d", w.Code)
	}
}

func TestUpdateForm_Success_Retourne200(t *testing.T) {
	svc := &servicemock.FormServiceMock{
		UpdateFormFn: func(ctx context.Context, form *models.Form) (*models.Form, error) {
			return &models.Form{Status: "UPDATED"}, nil
		},
		UpdateFormSectionFn: func(ctx context.Context, sections []models.FormSection) ([]models.FormSection, error) {
			return sections, nil
		},
	}
	h := handlers.NewFormHandler(svc)

	body := models.FormFormSection{
		Form:         &models.Form{Status: "CREATED"},
		FormSections: []models.FormSection{{SectionType: "STUDENT"}},
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, "/form", bytes.NewBuffer(b))
	w := httptest.NewRecorder()

	h.UpdateForm(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Attendu 200, obtenu %d", w.Code)
	}

	var response models.FormFormSection
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Impossible de décoder la réponse : %v", err)
	}
	if response.Form.Status != "UPDATED" {
		t.Errorf("Status attendu UPDATED, obtenu %s", response.Form.Status)
	}
}

func TestUpdateForm_SansSections_NAppellePasUpdateSection(t *testing.T) {
	sectionCalled := false
	svc := &servicemock.FormServiceMock{
		UpdateFormFn: func(ctx context.Context, form *models.Form) (*models.Form, error) {
			return &models.Form{Status: "UPDATED"}, nil
		},
		UpdateFormSectionFn: func(ctx context.Context, sections []models.FormSection) ([]models.FormSection, error) {
			sectionCalled = true
			return sections, nil
		},
	}
	h := handlers.NewFormHandler(svc)

	body := models.FormFormSection{Form: &models.Form{Status: "CREATED"}}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, "/form", bytes.NewBuffer(b))
	w := httptest.NewRecorder()

	h.UpdateForm(w, req)

	if sectionCalled {
		t.Error("UpdateFormSection ne devrait pas être appelé sans sections")
	}
	if w.Code != http.StatusOK {
		t.Errorf("Attendu 200, obtenu %d", w.Code)
	}
}
