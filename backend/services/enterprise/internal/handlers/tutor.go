package handlers

import (
	"encoding/json"
	"enterprise/internal/services"
	"net/http"
	"shared/models"
	"strconv"

	"github.com/gorilla/mux"
)

type TutorHandler struct {
	service *services.TutorService
}

func NewTutorHandler(service *services.TutorService) *TutorHandler {
	return &TutorHandler{service: service}
}

func (h *TutorHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	tutors, err := h.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tutors)
}

func (h *TutorHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	tutor, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tutor)
}

func (h *TutorHandler) Create(w http.ResponseWriter, r *http.Request) {
	var tutor models.Tutor
	if err := json.NewDecoder(r.Body).Decode(&tutor); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.Create(&tutor); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tutor)
}

func (h *TutorHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var tutor models.Tutor
	if err := json.NewDecoder(r.Body).Decode(&tutor); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tutor.ID = id

	if err := h.service.Update(&tutor); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tutor)
}

func (h *TutorHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

func (h *TutorHandler) GetByEnterprise(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	enterpriseID, err := strconv.Atoi(vars["enterprise_id"])
	if err != nil {
		http.Error(w, "Invalid enterprise ID", http.StatusBadRequest)
		return
	}

	tutors, err := h.service.GetByEnterprise(enterpriseID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tutors)
}

func (h *TutorHandler) GetWithStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	tutor, err := h.service.GetWithStats(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tutor)
}
