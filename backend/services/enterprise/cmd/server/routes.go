package main

import (
	"github.com/gorilla/mux"
	"net/http"
	sharedHandler "shared/handlers"
	"shared/middleware"
)

func setupRoutes(
	userHandler *sharedHandler.UserHandler,
	authMiddleware *middleware.AuthMiddleware,
) *mux.Router {
	r := mux.NewRouter()
	r.Use(corsMiddleware)

	api := r.PathPrefix("/api").Subrouter()

	// === ROUTES TUTEURS (avec handler générique) ===
	tutorsRouter := api.PathPrefix("/tutors").Subrouter()
	tutorsRouter.Use(authMiddleware.RequireAuth)

	// Routes CRUD génériques avec filtres automatiques basés sur les rôles
	tutorsRouter.HandleFunc("", userHandler.GetAll).Methods("GET")
	tutorsRouter.HandleFunc("", userHandler.Create).Methods("POST")

	tutorsRouter.HandleFunc("/{id}", userHandler.Delete).Methods("DELETE")

	tutorsRouter.HandleFunc("", corsPreflightHandler).Methods("OPTIONS")
	return r
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Configuration CORS
		origin := r.Header.Get("Origin")
		allowedOrigins := []string{
			"http://localhost:3000",
			"http://127.0.0.1:3000",
			"http://localhost",
			"http://127.0.0.1",
			"crissime.freeboxos.fr",
		}

		originAllowed := false
		for _, allowed := range allowedOrigins {
			if origin == allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				originAllowed = true
				break
			}
		}

		// En développement, autoriser tous les origins si non trouvé
		if !originAllowed {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept, Origin")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
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
