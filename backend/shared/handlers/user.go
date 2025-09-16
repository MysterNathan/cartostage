package handlers

import (
	"encoding/json"
	"net/http"
	"shared/models"
	"shared/services"
	"strconv"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// GetAll récupère tous les utilisateurs
func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// GetByID récupère un utilisateur par son ID
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	user, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user.ToPublic())
}

// Create crée un nouvel utilisateur
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var createReq models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdUser, err := h.service.Create(&createReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser.ToPublic())
}

// Update met à jour un utilisateur
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updateReq models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedUser, err := h.service.Update(id, &updateReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser.ToPublic())
}

// Delete supprime un utilisateur
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

// GetByEmail récupère un utilisateur par son email
func (h *UserHandler) GetByEmail(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email parameter is required", http.StatusBadRequest)
		return
	}

	user, err := h.service.GetByEmail(email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user.ToPublic())
}

// GetByUsername récupère un utilisateur par son nom d'utilisateur
func (h *UserHandler) GetByUsername(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username parameter is required", http.StatusBadRequest)
		return
	}

	user, err := h.service.GetByUsername(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user.ToPublic())
}

// GetByRole récupère tous les utilisateurs ayant un rôle spécifique
func (h *UserHandler) GetByRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	role := vars["role"]
	if role == "" {
		http.Error(w, "Role parameter is required", http.StatusBadRequest)
		return
	}

	users, err := h.service.GetByRole(role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convertir tous les utilisateurs en version publique
	publicUsers := make([]*models.User, len(users))
	for i, user := range users {
		publicUsers[i] = user.ToPublic()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(publicUsers)
}

// GetByEntity récupère tous les utilisateurs d'une entité spécifique
func (h *UserHandler) GetByEntity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	entityType := vars["entity_type"]
	entityIDStr := vars["entity_id"]

	if entityType == "" || entityIDStr == "" {
		http.Error(w, "Entity type and ID are required", http.StatusBadRequest)
		return
	}

	entityID, err := strconv.Atoi(entityIDStr)
	if err != nil {
		http.Error(w, "Invalid entity ID", http.StatusBadRequest)
		return
	}

	users, err := h.service.GetByEntity(entityType, entityID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convertir tous les utilisateurs en version publique
	publicUsers := make([]*models.User, len(users))
	for i, user := range users {
		publicUsers[i] = user.ToPublic()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(publicUsers)
}

// UpdateProfile met à jour le profil utilisateur
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var profileReq models.UpdateUserProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&profileReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedUser, err := h.service.UpdateProfile(id, &profileReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser.ToPublic())
}

// ChangePassword change le mot de passe d'un utilisateur
func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var changeReq models.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&changeReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.ChangePassword(id, &changeReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
