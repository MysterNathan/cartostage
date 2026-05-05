package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	sharedContext "shared/context"
	"shared/models"
	"stage/internal/services"
	"strconv"
	"strings"
)

type StageHandler struct {
	stageService services.StageServiceInterface
}

func NewStageHandler(stageService services.StageServiceInterface) *StageHandler {
	return &StageHandler{stageService: stageService}
}

func (h *StageHandler) GetStages(w http.ResponseWriter, r *http.Request) {
	claims := sharedContext.GetUserClaims(r.Context())
	var stagesData []models.Stage
	var err error

	if claims == nil {
		http.Error(w, `{"stages": []}`, http.StatusInternalServerError)
		return
	}

	stagesData, err = h.stageService.GetStages(r.Context())
	if err != nil {
		http.Error(w, `{"stages": []}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(stagesData); err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
}
func (h *StageHandler) GetStagesPublic(w http.ResponseWriter, r *http.Request) {
	var stagesData []models.StageWithDetails
	var err error

	stagesData, err = h.stageService.GetStagesPublic()
	if err != nil {
		http.Error(w, `{"stages": []}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(stagesData); err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
}
func (h *StageHandler) UpdateStage(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		http.Error(w, `{"error": "ID manquant"}`, http.StatusBadRequest)
		return
	}

	// Décoder le body en UpdateStageRequest
	var updateReq models.UpdateStageRequest
	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		http.Error(w, `{"error": "JSON invalide"}`, http.StatusBadRequest)
		return
	}

	// Vérifier qu'il y a au moins un champ à mettre à jour
	if !updateReq.HasUpdates() {
		http.Error(w, `{"error": "Aucun champ à mettre à jour"}`, http.StatusBadRequest)
		return
	}

	StageID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, `{"error": "ID manquant"}`, http.StatusBadRequest)
	}
	// Appeler le service pour l'update
	updatedStage, err := h.stageService.UpdateStage(r.Context(), StageID, updateReq)
	if err != nil {
		// Gérer les différents types d'erreurs
		if err.Error() == "stage-service not found" {
			http.Error(w, `{"error": "Stage non trouvé"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error": "Erreur serveur"}`, http.StatusInternalServerError)
		return
	}

	// Retourner le stage-service mis à jour
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"stage-service": updatedStage.ToPublic(),
	})
}

func (h *StageHandler) CreateStage(w http.ResponseWriter, r *http.Request) {
	// Décoder le body en CreateStageRequest
	var createReq models.CreateStageRequest
	if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
		http.Error(w, `{"error": "JSON invalide"}`, http.StatusBadRequest)
		return
	}

	// Validation basique
	if createReq.StageOfferID == 0 || createReq.StudentID == 0 {
		http.Error(w, `{"error": "stage_offer_id et student_id sont obligatoires"}`, http.StatusBadRequest)
		return
	}

	if createReq.Status == "" {
		http.Error(w, `{"error": "status est obligatoire"}`, http.StatusBadRequest)
		return
	}

	// Vérifier que end_date > start_date
	if !createReq.EndDate.After(createReq.StartDate) {
		http.Error(w, `{"error": "end_date doit être après start_date"}`, http.StatusBadRequest)
		return
	}

	// Appeler le service pour créer le stage-service
	createdStage, err := h.stageService.CreateStage(r.Context(), createReq)
	if err != nil {
		// Gérer les différents types d'erreurs
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "already exists") {
			http.Error(w, `{"error": "Stage déjà existant"}`, http.StatusConflict)
			return
		}
		if strings.Contains(err.Error(), "foreign key") {
			http.Error(w, `{"error": "Référence invalide (stage_offer_id, student_id, etc.)"}`, http.StatusBadRequest)
			return
		}
		http.Error(w, `{"error": "Erreur serveur"}`, http.StatusInternalServerError)
		return
	}

	// Retourner le stage-service créé
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"stage-service": createdStage.ToPublic(),
	})
}

func (h *StageHandler) DeleteStage(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID depuis l'URL
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error": "ID invalide"}`, http.StatusBadRequest)
		return
	}

	// Appeler le service pour supprimer le stage-service
	err = h.stageService.DeleteStage(r.Context(), id)
	if err != nil {
		// Gérer les différents types d'erreurs
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, `{"error": "Stage non trouvé"}`, http.StatusNotFound)
			return
		}
		if strings.Contains(err.Error(), "foreign key") {
			http.Error(w, `{"error": "Stage ne peut pas être supprimé car il a des dépendances"}`, http.StatusConflict)
			return
		}
		http.Error(w, `{"error": "Erreur serveur"}`, http.StatusInternalServerError)
		return
	}

	// Retourner une réponse de succès
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Stage supprimé avec succès",
	})
}
