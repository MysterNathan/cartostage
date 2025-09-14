package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"shared/models"
	"stage-map/internal/repositories"
	"strconv"
	"strings"
)

type StageHandler struct {
	repo *repositories.StageRepository
}

func NewStageHandler(repo *repositories.StageRepository) *StageHandler {
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

// SaveStage - Pour sauvegarder un seul stage
func (h *StageHandler) SaveStage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var stage models.Stage
	if err := json.NewDecoder(r.Body).Decode(&stage); err != nil {
		log.Printf("Erreur décodage JSON: %v", err)
		http.Error(w, `{"error": "Format invalide"}`, http.StatusBadRequest)
		return
	}

	// Validation des champs obligatoires
	if stage.Poste == "" || stage.Adresse == "" || stage.Entreprise == "" {
		http.Error(w, `{"error": "Champs obligatoires manquants"}`, http.StatusBadRequest)
		return
	}

	// Sauvegarde
	if err := h.repo.SaveStage(&stage); err != nil {
		log.Printf("Erreur sauvegarde: %v", err)
		http.Error(w, `{"error": "Erreur serveur"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"id":      stage.ID,
	})
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
	// Cette méthode pourrait être ajoutée au repositories si nécessaire
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

// DeleteStage - Supprime un stage spécifique par ID
func (h *StageHandler) DeleteStage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Récupérer l'ID depuis l'URL ou les paramètres
	// Supposons que l'ID soit passé dans l'URL comme /api/stages/123
	// Vous devrez adapter selon votre routeur
	idStr := r.URL.Path[len("/api/stages/"):]
	if idStr == "" {
		http.Error(w, `{"error": "ID manquant"}`, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error": "ID invalide"}`, http.StatusBadRequest)
		return
	}

	// Supprimer le stage
	if err := h.repo.DeleteStage(id); err != nil {
		if strings.Contains(err.Error(), "aucun stage trouvé") {
			http.Error(w, `{"error": "Stage non trouvé"}`, http.StatusNotFound)
			return
		}
		log.Printf("Erreur suppression: %v", err)
		http.Error(w, `{"error": "Erreur serveur"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true, "message": "Stage supprimé avec succès"}`))
}

// UpdateStage - Met à jour un stage existant
func (h *StageHandler) UpdateStage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Récupérer l'ID depuis l'URL avec mux
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

	// Décoder le stage depuis le body
	var stage models.Stage
	if err := json.NewDecoder(r.Body).Decode(&stage); err != nil {
		log.Printf("Erreur décodage JSON: %v", err)
		http.Error(w, `{"error": "Format invalide"}`, http.StatusBadRequest)
		return
	}

	// Validation des champs obligatoires
	if stage.Poste == "" || stage.Adresse == "" || stage.Entreprise == "" {
		http.Error(w, `{"error": "Champs obligatoires manquants"}`, http.StatusBadRequest)
		return
	}

	// S'assurer que l'ID du stage correspond à celui de l'URL
	stage.ID = id

	// Mettre à jour le stage
	if err := h.repo.UpdateStage(&stage); err != nil {
		if strings.Contains(err.Error(), "aucun stage trouvé") {
			http.Error(w, `{"error": "Stage non trouvé"}`, http.StatusNotFound)
			return
		}
		log.Printf("Erreur mise à jour: %v", err)
		http.Error(w, `{"error": "Erreur serveur"}`, http.StatusInternalServerError)
		return
	}

	// Retourner le stage mis à jour
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"success": true,
		"message": "Stage mis à jour avec succès",
		"stage":   stage,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Erreur encodage JSON: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
	}
}
