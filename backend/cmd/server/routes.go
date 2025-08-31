package main

import (
	"net/http"

	"backend/internal/handlers"
	"github.com/gorilla/mux"
)

func setupRoutes(stageHandler *handlers.StageHandler) *mux.Router {
	r := mux.NewRouter()

	// Middleware CORS global
	r.Use(corsMiddleware)

	// Équivalent des routes NextJS
	api := r.PathPrefix("/api").Subrouter()

	// Route principale des stages - GET et POST
	api.HandleFunc("/stages", stageHandler.GetAllStages).Methods("GET", "OPTIONS")
	api.HandleFunc("/stages", stageHandler.SaveAllStages).Methods("POST", "OPTIONS")

	// Route pour un stage spécifique
	api.HandleFunc("/stages/{id:[0-9]+}", stageHandler.GetStageByID).Methods("GET", "OPTIONS")

	// Route pour les filtres
	api.HandleFunc("/stages/filters", stageHandler.GetFilterOptions).Methods("GET", "OPTIONS")

	return r
}

// Middleware CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Autoriser uniquement le frontend
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// OPTIONS preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
