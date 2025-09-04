package main

import (
	"backend/internal/middleware"
	"net/http"

	"backend/internal/handlers"
	"github.com/gorilla/mux"
)

func setupRoutes(
	stageHandler *handlers.StageHandler,
	filiereHandler *handlers.FiliereHandler,
	authHandler *handlers.AuthHandler,
	authMiddleware *middleware.AuthMiddleware,
) *mux.Router {
	r := mux.NewRouter()

	// Middleware CORS appliqué en premier
	r.Use(corsMiddleware)

	api := r.PathPrefix("/api").Subrouter()

	// Route de login (non protégée)
	api.HandleFunc("/login", authHandler.Login).Methods("POST")
	api.HandleFunc("/login", corsPreflightHandler).Methods("OPTIONS")

	// Routes publiques des stages (lecture seule)
	stagesPublic := api.PathPrefix("/stages").Subrouter()
	stagesPublic.HandleFunc("", stageHandler.GetAllStages).Methods("GET")
	stagesPublic.HandleFunc("/{id:[0-9]+}", stageHandler.GetStageByID).Methods("GET")
	stagesPublic.HandleFunc("/filters", stageHandler.GetFilterOptions).Methods("GET")

	// Routes admin protégées (CRUD complet)
	adminRoutes := api.PathPrefix("/admin").Subrouter()
	adminRoutes.Use(authMiddleware.RequireAuth)

	// Routes stages admin
	adminStages := adminRoutes.PathPrefix("/stages").Subrouter()
	adminStages.HandleFunc("", stageHandler.SaveStage).Methods("POST")
	adminStages.HandleFunc("/{id}", stageHandler.DeleteStage).Methods("DELETE")
	adminStages.HandleFunc("/{id}", stageHandler.UpdateStage).Methods("PUT")

	// Routes filieres admin
	adminFilieres := adminRoutes.PathPrefix("/filieres").Subrouter()
	adminFilieres.HandleFunc("", filiereHandler.GetFilieres).Methods("GET")
	adminFilieres.HandleFunc("", filiereHandler.CreateFiliere).Methods("POST")
	adminFilieres.HandleFunc("/{id}", filiereHandler.UpdateFiliere).Methods("PUT")
	adminFilieres.HandleFunc("/{id}", filiereHandler.DeleteFiliere).Methods("DELETE")

	// Options pour les routes publiques
	stagesPublic.HandleFunc("", corsPreflightHandler).Methods("OPTIONS")
	stagesPublic.HandleFunc("/{id:[0-9]+}", corsPreflightHandler).Methods("OPTIONS")
	stagesPublic.HandleFunc("/filters", corsPreflightHandler).Methods("OPTIONS")

	// Options pour les routes admin
	adminStages.HandleFunc("", corsPreflightHandler).Methods("OPTIONS")
	adminStages.HandleFunc("/{id}", corsPreflightHandler).Methods("OPTIONS")
	adminFilieres.HandleFunc("", corsPreflightHandler).Methods("OPTIONS")
	adminFilieres.HandleFunc("/{id}", corsPreflightHandler).Methods("OPTIONS")

	return r
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Configuration CORS plus spécifique pour le développement
		origin := r.Header.Get("Origin")
		if origin == "http://localhost:3000" || origin == "http://127.0.0.1:3000" || origin == "crissime.freeboxos.fr" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			// En développement, on peut autoriser tous les origins
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Gestion des requêtes preflight OPTIONS
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Handler spécifique pour les requêtes OPTIONS (preflight)
func corsPreflightHandler(w http.ResponseWriter, r *http.Request) {
	// Les headers CORS sont déjà définis par le middleware
	w.WriteHeader(http.StatusNoContent)
}
