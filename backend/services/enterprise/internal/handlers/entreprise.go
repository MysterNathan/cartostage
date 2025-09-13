package handlers

import (
	"backend/services/enterprise/internal/services"
	"backend/shared/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type EnterpriseHandler struct {
	service *services.EnterpriseService
}

func NewEnterpriseHandler(service *services.EnterpriseService) *EnterpriseHandler {
	return &EnterpriseHandler{service: service}
}

func (h *EnterpriseHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	enterprises, err := h.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(enterprises)
}

func (h *EnterpriseHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	enterprise, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(enterprise)
}

func (h *EnterpriseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var enterprise models.Enterprise
	if err := json.NewDecoder(r.Body).Decode(&enterprise); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.Create(&enterprise); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(enterprise)
}

func (h *EnterpriseHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var enterprise models.Enterprise
	if err := json.NewDecoder(r.Body).Decode(&enterprise); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	enterprise.ID = id

	if err := h.service.Update(&enterprise); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(enterprise)
}

func (h *EnterpriseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *EnterpriseHandler) GetWithStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	enterprise, err := h.service.GetWithStats(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(enterprise)
}
