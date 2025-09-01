package handlers

import (
	"backend/internal/models"
	"backend/internal/repository"
	"encoding/json"
	"log"
	"net/http"
)

type FiliereHandler struct {
	repo *repository.FiliereRepository
}

func NewFiliereHandler(repo *repository.FiliereRepository) *FiliereHandler {
	return &FiliereHandler{repo: repo}
}

// GetFilieres
func (h *FiliereHandler) GetFilieres(w http.ResponseWriter, r *http.Request) {

	var filiereData *models.FilieresData
	var err error

	filiereData, err = h.repo.GetFilieres()

	if err != nil {
		log.Printf("Erreur lecture: %v", err)
		http.Error(w, `{"filieres": []}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(filiereData); err != nil {
		log.Printf("Erreur encodage JSON: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
}
