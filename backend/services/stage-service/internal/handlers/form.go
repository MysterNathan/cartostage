package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"shared/models"
	"stage/internal/services"
)

type FormHandler struct {
	formService services.FormServiceInterface
}

func NewFormHandler(formService services.FormServiceInterface) *FormHandler {
	return &FormHandler{formService: formService}
}

func (h FormHandler) Get(w http.ResponseWriter, r *http.Request) {
	form, err := h.formService.Get(r.Context())
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	if form == nil {
		http.Error(w, `{"error": "not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(form); err != nil {
		http.Error(w, `{"error": "encoding failed"}`, http.StatusInternalServerError)
		return
	}
}

func (h FormHandler) UpdateForm(w http.ResponseWriter, r *http.Request) {
	var data models.FormFormSection
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, `{"error": "JSON invalide"}`, http.StatusBadRequest)
		return
	}
	fmt.Printf("data content: %v data status: %v\n", data.Form.Content, data.Form.Status)
	if data.Form == nil {
		http.Error(w, `{"error": "form field missing"}`, http.StatusBadRequest)
	}
	updatedForm, err := h.formService.UpdateForm(r.Context(), data.Form)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	var updatedFormSections []models.FormSection
	if len(data.FormSections) > 0 {
		updatedFormSections, err = h.formService.UpdateFormSection(r.Context(), data.FormSections)
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
	}

	response := models.FormFormSection{
		Form:         updatedForm,
		FormSections: updatedFormSections,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
