package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	sharedContext "shared/context"
	"shared/models"
	"shared/services"
	"strconv"
	"strings"

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
	log.Println(targetRole)
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

func (h *UserHandler) GetStudents(w http.ResponseWriter, r *http.Request) {
	claims := sharedContext.GetUserClaims(r.Context())
	if claims.Role != "tutor" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
	students, err := h.userService.GetStudentByTutor(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest
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

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	// Vérifier que l'ID est fourni
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Convertir l'ID en entier
	id, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// Décoder la requête d'update avec la bonne structure
	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Vérifier s'il y a des champs à mettre à jour
	if !req.HasUpdates() {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}

	// Valider le rôle si fourni
	if req.Role != nil && !req.Role.IsValid() {
		http.Error(w, "Invalid role", http.StatusBadRequest)
		return
	}

	// Appeler le service d'update
	user, err := h.userService.UpdateUser(r.Context(), id, &req)
	if err != nil {
		// Gérer les différents types d'erreurs
		switch {
		case strings.Contains(err.Error(), "unauthorized"):
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		case strings.Contains(err.Error(), "forbidden"):
			http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
			return
		case strings.Contains(err.Error(), "not found"):
			http.Error(w, "User not found", http.StatusNotFound)
			return
		case strings.Contains(err.Error(), "validation"):
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		case strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "already exists"):
			http.Error(w, "User with this username or email already exists", http.StatusConflict)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Retourner l'utilisateur mis à jour (version publique)
	publicUser := user.ToPublic()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200 OK au lieu de 201 Created pour un update
	json.NewEncoder(w).Encode(publicUser)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID depuis les paramètres d'URL
	vars := mux.Vars(r)
	userID := vars["id"]

	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Convertir l'ID en entier
	id, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// Créer la requête de suppression
	req := &models.DeleteUserRequest{
		Id: id,
		// Le Role n'est pas utilisé dans le service, on peut l'omettre
		// ou laisser la valeur par défaut
	}

	// Appeler le service de suppression
	deletedID, err := h.userService.Delete(r.Context(), req)
	if err != nil {
		// Gérer les différents types d'erreurs
		switch {
		case strings.Contains(err.Error(), "unauthorized"):
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		case strings.Contains(err.Error(), "forbidden"):
			http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
			return
		case strings.Contains(err.Error(), "not found"):
			http.Error(w, "User not found", http.StatusNotFound)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Réponse de succès avec l'ID supprimé
	response := map[string]interface{}{
		"message":    "User deleted successfully",
		"deleted_id": deletedID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
