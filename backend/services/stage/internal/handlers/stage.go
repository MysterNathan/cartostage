package handlers

import (
	"encoding/json"
	"net/http"
	sharedContext "shared/context"
	"shared/models"
	"stage/internal/services"
)

type StageHandler struct {
	service *services.StageService
}

func NewStageHandler(service *services.StageService) *StageHandler {
	return &StageHandler{service: service}
}

func (h *StageHandler) GetStages(w http.ResponseWriter, r *http.Request) {
	claims := sharedContext.GetUserClaims(r.Context())
	var stagesData []models.Stage
	var err error

	if claims == nil {
		stagesData, err = h.service.GetStagesPublic()
		return
	}

	stagesData, err = h.service.GetStages(r.Context())

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
