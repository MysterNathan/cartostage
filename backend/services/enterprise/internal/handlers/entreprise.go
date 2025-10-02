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

//func (h *EnterpriseHandler) GetAll(w http.ResponseWriter, r *http.Request) {
//	enterprises, err := h.service.GetAll()
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	json.NewEncoder(w).Encode(enterprises)
//}
//
//func (h *EnterpriseHandler) GetByID(w http.ResponseWriter, r *http.Request) {
//	vars := mux.Vars(r)
//	id, err := strconv.Atoi(vars["id"])
//	if err != nil {
//		http.Error(w, "Invalid ID", http.StatusBadRequest)
//		return
//	}
//
//	enterprise, err := h.service.GetByID(id)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusNotFound)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	json.NewEncoder(w).Encode(enterprise)
//}
//
//func (h *EnterpriseHandler) GetMe(w http.ResponseWriter, r *http.Request) {
//	log.Println("Entrer dans le Handler getMe") //todo supprimer les logs
//	print("Entrer dans le Handlser getMe")
//
//	// Utiliser la fonction utilitaire pour récupérer l'ID de l'entreprise
//	enterpriseID, ok := middleware.GetEnterpriseIDFromContext(r.Context())
//	if !ok {
//		http.Error(w, "Enterprise ID not found in context", http.StatusInternalServerError)
//		return
//	}
//
//	// Récupérer les données de l'entreprise
//	enterprise, err := h.service.GetMe(enterpriseID)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusNotFound)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	json.NewEncoder(w).Encode(enterprise)
//}
//
//func (h *EnterpriseHandler) Create(w http.ResponseWriter, r *http.Request) {
//	var enterprise models.Enterprise
//	if err := json.NewDecoder(r.Body).Decode(&enterprise); err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	// Le service Create retourne (*models.Enterprise, error)
//	createdEnterprise, err := h.service.Create(&enterprise)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusCreated)
//	json.NewEncoder(w).Encode(createdEnterprise)
//}
//
//func (h *EnterpriseHandler) Update(w http.ResponseWriter, r *http.Request) {
//	vars := mux.Vars(r)
//	id, err := strconv.Atoi(vars["id"])
//	if err != nil {
//		http.Error(w, "Invalid ID", http.StatusBadRequest)
//		return
//	}
//
//	var enterprise models.Enterprise
//	if err := json.NewDecoder(r.Body).Decode(&enterprise); err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	// Le service Update prend (id int, enterprise *models.Enterprise) et retourne (*models.Enterprise, error)
//	updatedEnterprise, err := h.service.Update(id, &enterprise)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	json.NewEncoder(w).Encode(updatedEnterprise)
//}
//
//func (h *EnterpriseHandler) Delete(w http.ResponseWriter, r *http.Request) {
//	vars := mux.Vars(r)
//	id, err := strconv.Atoi(vars["id"])
//	if err != nil {
//		http.Error(w, "Invalid ID", http.StatusBadRequest)
//		return
//	}
//
//	if err := h.service.Delete(id); err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	w.WriteHeader(http.StatusNoContent)
//}
//
//func (h *EnterpriseHandler) GetWithStats(w http.ResponseWriter, r *http.Request) {
//	vars := mux.Vars(r)
//	id, err := strconv.Atoi(vars["id"])
//	if err != nil {
//		http.Error(w, "Invalid ID", http.StatusBadRequest)
//		return
//	}
//
//	enterprise, err := h.service.GetWithStats(id)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusNotFound)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	json.NewEncoder(w).Encode(enterprise)
//}
//
//// Handlers pour les stages (à décommenter quand tu implémenteras)
///*
//func (h *EnterpriseHandler) GetMyStages(w http.ResponseWriter, r *http.Request) {
//    enterpriseID, ok := middleware.GetEnterpriseIDFromContext(r.Context())
//    if !ok {
//        http.Error(w, "Enterprise ID not found in context", http.StatusInternalServerError)
//        return
//    }
//
//    stages, err := h.service.GetStages(enterpriseID)
//    if err != nil {
//        http.Error(w, err.Error(), http.StatusInternalServerError)
//        return
//    }
//
//    w.Header().Set("Content-Type", "application/json")
//    json.NewEncoder(w).Encode(stages)
//}
//*/

func (h *EnterpriseHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	claims := sharedContext.GetClaimsFromContext(r.Context())
	if claims.Role != "tutor" {
		http.Error(w, "User not allowed", http.StatusForbidden)
	}
	stats, _ := h.service.GetStats(r.Context())
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
	return
}
