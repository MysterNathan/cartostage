package main

import (
	"enterprise/internal/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"shared/middleware"
)

func setupRoutes(enterpriseHandler *handlers.EnterpriseHandler, tutorHandler *handlers.TutorHandler, authMiddleware *middleware.AuthMiddleware) *mux.Router {
	r := mux.NewRouter()
	r.Use(corsMiddleware)

	api := r.PathPrefix("/api").Subrouter()

	// Routes pour les entreprises
	enterprisesRouter := api.PathPrefix("/enterprises").Subrouter()
	enterprisesRouter.Use(authMiddleware.RequireAuth)

	enterprisesRouter.HandleFunc("", enterpriseHandler.GetAll).Methods("GET")
	enterprisesRouter.HandleFunc("/me", enterpriseHandler.GetMe).Methods("GET")
	enterprisesRouter.HandleFunc("/{id}", enterpriseHandler.GetByID).Methods("GET")
	enterprisesRouter.HandleFunc("/", enterpriseHandler.Create).Methods("POST")
	enterprisesRouter.HandleFunc("/{id}", enterpriseHandler.Update).Methods("PUT")
	enterprisesRouter.HandleFunc("/{id}", enterpriseHandler.Delete).Methods("DELETE")
	enterprisesRouter.HandleFunc("/{id}/stats", enterpriseHandler.GetWithStats).Methods("GET")

	// Gestion OPTIONS pour toutes les routes
	enterprisesRouter.HandleFunc("", corsPreflightHandler).Methods("OPTIONS")
	enterprisesRouter.HandleFunc("/filters", corsPreflightHandler).Methods("OPTIONS")
	enterprisesRouter.HandleFunc("/{id:[0-9]+}", corsPreflightHandler).Methods("OPTIONS")
	enterprisesRouter.HandleFunc("", corsPreflightHandler).Methods("OPTIONS")
	enterprisesRouter.HandleFunc("/{id:[0-9]+}", corsPreflightHandler).Methods("OPTIONS")
	enterprisesRouter.HandleFunc("/me", corsPreflightHandler).Methods("OPTIONS")

	tutorsRouter := api.PathPrefix("/tutors").Subrouter()

	// Routes pour les tuteurs
	tutorsRouter.HandleFunc("/tutors", tutorHandler.GetAll).Methods("GET")
	tutorsRouter.HandleFunc("/tutors/{id}", tutorHandler.GetByID).Methods("GET")
	tutorsRouter.HandleFunc("/tutors", tutorHandler.Create).Methods("POST")
	tutorsRouter.HandleFunc("/tutors/{id}", tutorHandler.Update).Methods("PUT")
	tutorsRouter.HandleFunc("/tutors/{id}", tutorHandler.Delete).Methods("DELETE")
	tutorsRouter.HandleFunc("/tutors/enterprise/{enterprise_id}", tutorHandler.GetByEnterprise).Methods("GET")
	api.HandleFunc("/tutors/{id}/stats", tutorHandler.GetWithStats).Methods("GET")

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
