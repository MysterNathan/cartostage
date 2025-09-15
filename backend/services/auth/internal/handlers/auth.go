package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"auth/internal/services"
	"github.com/gorilla/mux"
	"shared/middleware"
	"shared/models"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validation basique
	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Login avec le service adapté
	response, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// Récupérer le token depuis l'header Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header required", http.StatusBadRequest)
		return
	}

	// Extraire le token
	var token string
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	} else {
		http.Error(w, "Invalid authorization format", http.StatusBadRequest)
		return
	}

	// Refresh avec le service
	response, err := h.authService.RefreshToken(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Récupérer les claims du context (middleware auth requis)
	claims := middleware.GetClaimsFromContext(r.Context())
	if claims == nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Logout (invalider le token si vous avez une blacklist)
	err := h.authService.Logout(claims.SessionID)
	if err != nil {
		http.Error(w, "Logout failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Successfully logged out",
	})
}

func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Récupérer les informations depuis le token
	claims := middleware.GetClaimsFromContext(r.Context())
	if claims == nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Récupérer le profil complet (retourne UserProfile sans hash)
	profile, err := h.authService.GetUserProfile(claims.UserID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

func (h *AuthHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validation des champs requis
	if req.Username == "" || req.Password == "" || req.Role == "" {
		http.Error(w, "Username, password, and role are required", http.StatusBadRequest)
		return
	}

	// Validation des rôles autorisés
	validRoles := map[string]bool{
		"student":    true,
		"teacher":    true,
		"enterprise": true,
		"admin":      true,
	}
	if !validRoles[req.Role] {
		http.Error(w, "Invalid role", http.StatusBadRequest)
		return
	}

	// Pour enterprise et teacher, entity_id peut être requis
	if (req.Role == "enterprise" || req.Role == "teacher") && req.EntityID == nil {
		http.Error(w, "entity_id is required for this role", http.StatusBadRequest)
		return
	}

	// CreateUser retourne maintenant UserProfile
	userProfile, err := h.authService.CreateUser(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userProfile)
}

func (h *AuthHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	// Pagination
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page == 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 10
	}

	// Filtres
	role := r.URL.Query().Get("role")
	entityID := r.URL.Query().Get("entity_id")

	// GetUsers retourne maintenant []*UserProfile
	users, total, err := h.authService.GetUsers(page, limit, role, entityID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"users": users,
		"total": total,
		"page":  page,
		"limit": limit,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Utiliser mux.Vars pour récupérer l'ID depuis l'URL
	vars := mux.Vars(r)
	userIDStr := vars["id"]
	if userIDStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Vérifier les permissions
	claims := middleware.GetClaimsFromContext(r.Context())
	if claims == nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Seul un admin ou l'utilisateur lui-même peut modifier
	if claims.Role != "admin" && claims.UserID != userID {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}

	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Les non-admins ne peuvent pas changer leur rôle
	if claims.Role != "admin" && req.Role != nil {
		http.Error(w, "Cannot change role", http.StatusForbidden)
		return
	}

	// UpdateUser retourne maintenant UserProfile
	userProfile, err := h.authService.UpdateUser(userID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userProfile)
}

func (h *AuthHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Utiliser mux.Vars pour récupérer l'ID depuis l'URL
	vars := mux.Vars(r)
	userIDStr := vars["id"]
	if userIDStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = h.authService.DeleteUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaimsFromContext(r.Context())
	if claims == nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Retourner les informations du token (utile pour debug)
	response := map[string]interface{}{
		"valid":      true,
		"user_id":    claims.UserID,
		"username":   claims.Username,
		"role":       claims.Role,
		"entity_id":  claims.EntityID,
		"scopes":     claims.Scope,
		"session_id": claims.SessionID,
		"expires_at": claims.ExpiresAt.Time,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Méthode helper pour obtenir un utilisateur par ID depuis l'URL (pour d'autres routes)
func (h *AuthHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["id"]
	if userIDStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Vérifier les permissions (admin ou self)
	claims := middleware.GetClaimsFromContext(r.Context())
	if claims == nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	if claims.Role != "admin" && claims.UserID != userID {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}

	userProfile, err := h.authService.GetUserProfile(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userProfile)
}
