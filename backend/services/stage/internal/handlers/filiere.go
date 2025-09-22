package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"shared/models"
	"stage/internal/repositories"
	"strconv"
)

type FiliereHandler struct {
	repo *repositories.FiliereRepository
}

func NewFiliereHandler(repo *repositories.FiliereRepository) *FiliereHandler {
	return &FiliereHandler{repo: repo}
}

// GetFilieres
func (h *FiliereHandler) GetFilieres(w http.ResponseWriter, r *http.Request) {

	var filiereData *models.Filieres
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

// CreateFiliere
func (h *FiliereHandler) CreateFiliere(w http.ResponseWriter, r *http.Request) {
	// Vérifier que la méthode est POST
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Décoder le JSON de la requête
	var filiere models.Filiere
	if err := json.NewDecoder(r.Body).Decode(&filiere); err != nil {
		log.Printf("Erreur décodage JSON: %v", err)
		http.Error(w, "Format JSON invalide", http.StatusBadRequest)
		return
	}

	// Validation basique des champs obligatoires
	if filiere.Code == "" || filiere.Label == "" || filiere.Color == "" {
		http.Error(w, "Les champs code, label et color sont obligatoires", http.StatusBadRequest)
		return
	}

	// Appeler le repositories pour créer la filière
	createdFiliere, err := h.repo.CreateFiliere(filiere)
	if err != nil {
		log.Printf("Erreur création filière: %v", err)
		http.Error(w, "Erreur lors de la création de la filière", http.StatusInternalServerError)
		return
	}

	// Retourner la filière créée
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(createdFiliere); err != nil {
		log.Printf("Erreur encodage JSON: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
}

// UpdateFiliere
func (h *FiliereHandler) UpdateFiliere(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID depuis l'URL
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	// Décoder le JSON de la requête
	var filiere models.Filiere
	if err := json.NewDecoder(r.Body).Decode(&filiere); err != nil {
		log.Printf("Erreur décodage JSON: %v", err)
		http.Error(w, "Format JSON invalide", http.StatusBadRequest)
		return
	}

	// Validation basique des champs obligatoires
	if filiere.Code == "" || filiere.Label == "" || filiere.Color == "" {
		http.Error(w, "Les champs code, label et color sont obligatoires", http.StatusBadRequest)
		return
	}

	// Assigner l'ID
	filiere.ID = id

	// Appeler le repositories pour mettre à jour la filière
	updatedFiliere, err := h.repo.UpdateFiliere(filiere)
	if err != nil {
		log.Printf("Erreur mise à jour filière: %v", err)
		http.Error(w, "Erreur lors de la mise à jour de la filière", http.StatusInternalServerError)
		return
	}

	// Retourner la filière mise à jour
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(updatedFiliere); err != nil {
		log.Printf("Erreur encodage JSON: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
}

// DeleteFiliere
func (h *FiliereHandler) DeleteFiliere(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID depuis l'URL
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	print(id)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	// Appeler le repositories pour supprimer la filière
	err = h.repo.DeleteFiliere(id)
	if err != nil {
		log.Printf("Erreur suppression filière: %v", err)
		http.Error(w, "Erreur lors de la suppression de la filière", http.StatusInternalServerError)
		return
	}

	// Retourner un statut 204 No Content
	w.WriteHeader(http.StatusNoContent)
}
