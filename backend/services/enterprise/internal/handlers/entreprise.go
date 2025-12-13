package handlers

import (
	"encoding/json"
	"enterprise/internal/services"
	//"log"
	"net/http"
	sharedContext "shared/context"
	//"shared/middleware"
	//"shared/models"
	//"strconv"
	//
	//"github.com/gorilla/mux"
)

type EnterpriseHandler struct {
	service *services.EnterpriseService
}

func NewEnterpriseHandler(service *services.EnterpriseService) *EnterpriseHandler {
	return &EnterpriseHandler{service: service}
}

func (h *EnterpriseHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	claims := sharedContext.GetClaimsFromContext(r.Context())
	if claims.Role != "tutor" && claims.Role != "admin" {
		http.Error(w, "User not allowed", http.StatusForbidden)
	}
	stats, _ := h.service.GetStats(r.Context())
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
	return
}
