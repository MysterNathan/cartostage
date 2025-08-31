package handlers

import (
	"backend/internal/models"
	"backend/internal/repository"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type StageHandler struct {
	repo *repository.StageRepository
}

func NewStageHandler(repo *repository.StageRepository) *StageHandler {
	return &StageHandler{repo: repo}
}

// GetAllStages - Équivalent du GET NextJS
func (h *StageHandler) GetAllStages(w http.ResponseWriter, r *http.Request) {
	// Vérifier les paramètres de requête pour les filtres
	filiere := r.URL.Query().Get("filiere")
	commune := r.URL.Query().Get("commune")
	availableOnly := r.URL.Query().Get("available") == "true"

	var stagesData *models.StagesData
	var err error

	// Si des filtres sont présents, utiliser la méthode filtrée
	if filiere != "" || commune != "" || availableOnly {
		stagesData, err = h.repo.GetStagesWithFilters(filiere, commune, availableOnly)
	} else {
		stagesData, err = h.repo.GetAllStages()
	}

	if err != nil {
		log.Printf("Erreur lecture: %v", err)
		http.Error(w, `{"stages": []}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(stagesData); err != nil {
		log.Printf("Erreur encodage JSON: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
}

// SaveAllStages - Équivalent du POST NextJS
func (h *StageHandler) SaveAllStages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var stagesData models.StagesData
	if err := json.NewDecoder(r.Body).Decode(&stagesData); err != nil {
		log.Printf("Erreur décodage JSON: %v", err)
		http.Error(w, `{"error": "Format invalide"}`, http.StatusBadRequest)
		return
	}

	// Validation basique
	if stagesData.Stages == nil {
		http.Error(w, `{"error": "Format invalide"}`, http.StatusBadRequest)
		return
	}

	// Sauvegarde
	if err := h.repo.SaveAllStages(&stagesData); err != nil {
		log.Printf("Erreur sauvegarde: %v", err)
		http.Error(w, `{"error": "Erreur serveur"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success": true}`))
}

// GetStageByID - Handler pour récupérer un stage spécifique
func (h *StageHandler) GetStageByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, `{"error": "ID manquant"}`, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error": "ID invalide"}`, http.StatusBadRequest)
		return
	}

	stage, err := h.repo.GetStageByID(id)
	if err != nil {
		log.Printf("Erreur récupération stage: %v", err)
		http.Error(w, `{"error": "Stage introuvable"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]*models.Stage{"stage": stage}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Erreur encodage JSON: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
	}
}

// GetFilterOptions - Handler pour récupérer les options de filtres
func (h *StageHandler) GetFilterOptions(w http.ResponseWriter, r *http.Request) {
	// Cette méthode pourrait être ajoutée au repository si nécessaire
	stagesData, err := h.repo.GetAllStages()
	if err != nil {
		log.Printf("Erreur récupération stages: %v", err)
		http.Error(w, `{"error": "Erreur serveur"}`, http.StatusInternalServerError)
		return
	}

	// Extraire les filières et communes uniques
	filieresSet := make(map[string]bool)
	communesSet := make(map[string]bool)

	for _, stage := range stagesData.Stages {
		filieresSet[stage.Filiere] = true
		communesSet[stage.Commune] = true
	}

	// Convertir en slices
	var filieres []string
	var communes []string

	for filiere := range filieresSet {
		filieres = append(filieres, filiere)
	}

	for commune := range communesSet {
		communes = append(communes, commune)
	}

	response := map[string]interface{}{
		"filieres": filieres,
		"communes": communes,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Erreur encodage JSON: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
	}
}
