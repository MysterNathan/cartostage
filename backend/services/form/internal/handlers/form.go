package handlers

import (
	"encoding/json"
	"form/internal/services"
	"net/http"
	"shared/models"
	"strconv"

	"github.com/gorilla/mux"
)

type FormHandler struct {
	formService *services.FormService
}

func NewformHandler(formService *services.FormService) *FormHandler {
	return &FormHandler{formService: formService}
}

func (h FormHandler) Get(w http.ResponseWriter, r *http.Request) {
	form, err := h.formService.Get(r.Context())
	if err != nil {
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(form)
	return
}

func (h FormHandler) UpdateForm(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		http.Error(w, `{"error": "ID manquant"}`, http.StatusBadRequest)
		return
	}
	formId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, `{"error": "ID manquant"}`, http.StatusBadRequest)
	}
	var data models.Form
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, `{"error": "JSON invalide"}`, http.StatusBadRequest)
		return
	}
	form, err := h.formService.UpdateForm(r.Context(), data, formId)
	if err != nil {
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(form)
	return
}

func (h FormHandler) CreateForm(w http.ResponseWriter, r *http.Request) {
	var data models.Form
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, `{"error": "JSON invalide"}`, http.StatusBadRequest)
		return
	}
	form, err := h.formService.CreateForm(r.Context(), data)
	if err != nil {
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(form)
	return
}
