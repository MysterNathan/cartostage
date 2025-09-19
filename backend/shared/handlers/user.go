package handlers

import (
	"encoding/json"
	"net/http"
	sharedContext "shared/context"
	"shared/models"
	"shared/services"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Extraire le rôle depuis l'URL ou utiliser le rôle de l'utilisateur actuel
	vars := mux.Vars(r)
	roleParam := vars["role"]

	var targetRole models.UserRole
	if roleParam != "" {
		targetRole = models.UserRole(roleParam)
	} else {
		// Si pas de rôle spécifié, utiliser le rôle de l'utilisateur actuel
		targetRole = sharedContext.GetUserRoleFromContext(r.Context())
	}

	if !targetRole.IsValid() {
		http.Error(w, "Invalid role", http.StatusBadRequest)
		return
	}

	users, err := h.userService.GetUsersByRole(r.Context(), targetRole)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req services.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userService.CreateUser(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	var req DeleteUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	res, err := h.userService.Delete(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusGone)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
