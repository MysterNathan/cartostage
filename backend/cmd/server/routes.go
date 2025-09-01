package main

import (
	"net/http"

	"backend/internal/handlers"
	"github.com/gorilla/mux"
)

func setupRoutes(stageHandler *handlers.StageHandler) *mux.Router {
	r := mux.NewRouter()

	// Équivalent des routes NextJS
	api := r.PathPrefix("/api").Subrouter()

	// Route principale des stages - GET et POST
	stages := api.PathPrefix("/stages").Subrouter()
	stages.HandleFunc("", stageHandler.GetAllStages).Methods("GET")
	stages.HandleFunc("", stageHandler.SaveAllStages).Methods("POST")

	// Route pour un stage spécifique
	stages.HandleFunc("/{id:[0-9]+}", stageHandler.GetStageByID).Methods("GET")

	// Route pour les filtres
	stages.HandleFunc("/filters", stageHandler.GetFilterOptions).Methods("GET")

	// Middleware CORS si nécessaire
	r.Use(corsMiddleware)

	return r
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
