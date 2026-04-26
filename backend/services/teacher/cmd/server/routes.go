package main

import (
	"net/http"
	"shared/middleware"
	"teacher/internal/handlers"

	"github.com/gorilla/mux"
)

func setupRoutes(
	userHandler *handlers.UserHandler,
	authMiddleware *middleware.AuthMiddleware,
) *mux.Router {
	r := mux.NewRouter()

	// Middleware CORS global
	r.Use(corsMiddleware)

	// API v1
	api := r.PathPrefix("/api").Subrouter()

	// === ROUTES STUDENT ===
	teacher := api.PathPrefix("/teacher").Subrouter()
	teacher.Use(authMiddleware.RequireAuth)

	// Preflight
	teacher.HandleFunc("", corsPreflightHandler).Methods("OPTIONS")
	teacher.HandleFunc("/{id}", corsPreflightHandler).Methods("OPTIONS")

	// === ROUTES USERS (TEACHER) ===
	users := teacher.PathPrefix("/users").Subrouter()
	users.Use(authMiddleware.RequireAuth)

	// Routes CRUD avec filtres automatiques basés sur les rôles
	users.HandleFunc("", userHandler.GetAll).Methods("GET")
	users.HandleFunc("/{id}", userHandler.GetByID).Methods("GET")
	users.HandleFunc("", userHandler.Create).Methods("POST")
	users.HandleFunc("/{id}", userHandler.Update).Methods("PUT")
	users.HandleFunc("/{id}", userHandler.Delete).Methods("DELETE")

	// Preflight
	users.HandleFunc("", corsPreflightHandler).Methods("OPTIONS")
	users.HandleFunc("/{id}", corsPreflightHandler).Methods("OPTIONS")

	// === HEALTH CHECK ===
	r.HandleFunc("/health", healthCheckHandler).Methods("GET")

	return r
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		allowedOrigins := map[string]bool{
			"http://localhost:3000":            true,
			"http://127.0.0.1:3000":            true,
			"http://localhost":                 true,
			"http://127.0.0.1":                 true,
			"http://mysternathan.freeboxos.fr": true,
		}

		// Vérifier si l'origin est autorisée
		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		} else if origin == "" {
			// Si pas d'origin (requête serveur à serveur), autoriser
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept, Origin")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Gestion des requêtes preflight OPTIONS
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func corsPreflightHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","service":"enterprise-api"}`))
}
