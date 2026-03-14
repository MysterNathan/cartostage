package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"shared/models"
	"stage/internal/services"
	"strconv"

	"github.com/gorilla/mux"
)

type FormHandler struct {
	formService *services.FormService
}

func NewFormHandler(formService *services.FormService) *FormHandler {
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
	fmt.Printf("data content:", data.Content, "data status:", data.Status, "\n")
	form, err := h.formService.UpdateForm(r.Context(), data, formId)
	if err != nil {
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(form)
	return
}
