package main

import (
	"net/http"

	"backend/internal/handlers"
	"github.com/gorilla/mux"
)

func setupRoutes(stageHandler *handlers.StageHandler, filiereHandler *handlers.FiliereHandler) *mux.Router {
	r := mux.NewRouter()

	// Middleware CORS appliqué en premier
	r.Use(corsMiddleware)

	api := r.PathPrefix("/api").Subrouter()

	// Route principale des stages - GET et POST
	stages := api.PathPrefix("/stages").Subrouter()
	stages.HandleFunc("", stageHandler.GetAllStages).Methods("GET")
	stages.HandleFunc("", stageHandler.SaveStage).Methods("POST")
	stages.HandleFunc("/{id}", stageHandler.DeleteStage).Methods("DELETE")
	stages.HandleFunc("/{id}", stageHandler.UpdateStage).Methods("PUT")

	// Support explicite des requêtes OPTIONS pour toutes les routes stages
	stages.HandleFunc("", corsPreflightHandler).Methods("OPTIONS")

	// Route pour un stage spécifique
	stages.HandleFunc("/{id:[0-9]+}", stageHandler.GetStageByID).Methods("GET")
	stages.HandleFunc("/{id:[0-9]+}", corsPreflightHandler).Methods("OPTIONS")

	// Route pour les filtres
	stages.HandleFunc("/filters", stageHandler.GetFilterOptions).Methods("GET")
	stages.HandleFunc("/filters", corsPreflightHandler).Methods("OPTIONS")

	// Routes des filières
	filieres := api.PathPrefix("/filieres").Subrouter()
	filieres.HandleFunc("", filiereHandler.GetFilieres).Methods("GET")
	filieres.HandleFunc("", filiereHandler.CreateFiliere).Methods("POST")
	filieres.HandleFunc("/{id}", filiereHandler.UpdateFiliere).Methods("PUT")
	filieres.HandleFunc("/{id}", filiereHandler.DeleteFiliere).Methods("DELETE")

	filieres.HandleFunc("/{id}", corsPreflightHandler).Methods("OPTIONS")

	filieres.HandleFunc("", corsPreflightHandler).Methods("OPTIONS")

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
