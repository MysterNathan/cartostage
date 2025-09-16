package main

import (
	"enterprise/internal/handlers"
	"github.com/gorilla/mux"
	"net/http"
	sharedHandler "shared/handlers"
	"shared/middleware"
)

func setupRoutes(enterpriseHandler *handlers.EnterpriseHandler, userHandler *sharedHandler.UserHandler, authMiddleware *middleware.AuthMiddleware) *mux.Router {
	r := mux.NewRouter()
	r.Use(corsMiddleware)
	r.Use(authMiddleware.RequireAuth)

	api := r.PathPrefix("/api").Subrouter()

	// Routes pour les entreprises
	enterprisesRouter := api.PathPrefix("/enterprises").Subrouter()

	enterprisesRouter.HandleFunc("", enterpriseHandler.GetAll).Methods("GET")
	enterprisesRouter.HandleFunc("/me", enterpriseHandler.GetMe).Methods("GET")
	enterprisesRouter.HandleFunc("/{id}", enterpriseHandler.GetByID).Methods("GET")
	enterprisesRouter.HandleFunc("/", enterpriseHandler.Create).Methods("POST")
	enterprisesRouter.HandleFunc("/{id}", enterpriseHandler.Update).Methods("PUT")
	enterprisesRouter.HandleFunc("/{id}", enterpriseHandler.Delete).Methods("DELETE")
	enterprisesRouter.HandleFunc("/{id}/stats", enterpriseHandler.GetWithStats).Methods("GET")

	// Gestion OPTIONS pour toutes les routes entreprise
	enterprisesRouter.HandleFunc("", corsPreflightHandler).Methods("OPTIONS")
	enterprisesRouter.HandleFunc("/filters", corsPreflightHandler).Methods("OPTIONS")
	enterprisesRouter.HandleFunc("/{id:[0-9]+}", corsPreflightHandler).Methods("OPTIONS")
	enterprisesRouter.HandleFunc("", corsPreflightHandler).Methods("OPTIONS")
	enterprisesRouter.HandleFunc("/{id:[0-9]+}", corsPreflightHandler).Methods("OPTIONS")
	enterprisesRouter.HandleFunc("/me", corsPreflightHandler).Methods("OPTIONS")

	tutorsRouter := api.PathPrefix("/tutors").Subrouter()

	// Routes pour les tuteurs
	tutorsRouter.HandleFunc("", userHandler.GetAll).Methods("GET")
	tutorsRouter.HandleFunc("/{id}", userHandler.GetByID).Methods("GET")
	tutorsRouter.HandleFunc("", userHandler.Create).Methods("POST")
	tutorsRouter.HandleFunc("/{id}", userHandler.Update).Methods("PUT")
	tutorsRouter.HandleFunc("/{id}", userHandler.Delete).Methods("DELETE")

	// Gestion OPTIONS pour toutes les routes tutors
	tutorsRouter.HandleFunc("", corsPreflightHandler).Methods("OPTIONS")
	tutorsRouter.HandleFunc("/filters", corsPreflightHandler).Methods("OPTIONS")
	tutorsRouter.HandleFunc("/{id:[0-9]+}", corsPreflightHandler).Methods("OPTIONS")
	tutorsRouter.HandleFunc("", corsPreflightHandler).Methods("OPTIONS")
	tutorsRouter.HandleFunc("/{id:[0-9]+}", corsPreflightHandler).Methods("OPTIONS")
	tutorsRouter.HandleFunc("/me", corsPreflightHandler).Methods("OPTIONS")

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
