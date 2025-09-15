package main

import (
	"auth/internal/handlers"
	"net/http"
	"shared/middleware"
	"shared/services"

	"github.com/gorilla/mux"
)

func setupRoutes(
	authHandler *handlers.AuthHandler,
	jwtService *services.JWTService, // Ajouter le JWT service
) *mux.Router {
	r := mux.NewRouter()
	r.Use(corsMiddleware)

	// Créer le middleware d'auth avec le JWT service
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	api := r.PathPrefix("/api").Subrouter()

	// Routes publiques (non protégées)
	api.HandleFunc("/login", authHandler.Login).Methods("POST")
	api.HandleFunc("/login", corsPreflightHandler).Methods("OPTIONS")

	api.HandleFunc("/refresh", authHandler.RefreshToken).Methods("POST", "OPTIONS")

	// Routes protégées
	protected := api.PathPrefix("").Subrouter()
	protected.Use(authMiddleware.RequireAuth) // Utiliser le middleware

	protected.HandleFunc("/profile", authHandler.GetProfile).Methods("GET")
	protected.HandleFunc("/logout", authHandler.Logout).Methods("POST")
	protected.HandleFunc("/validate", authHandler.ValidateToken).Methods("GET")

	// Routes admin uniquement
	adminRoutes := api.PathPrefix("/admin").Subrouter()
	adminRoutes.Use(authMiddleware.RequireRole("admin"))

	adminRoutes.HandleFunc("/users", authHandler.CreateUser).Methods("POST")
	adminRoutes.HandleFunc("/users", authHandler.GetUsers).Methods("GET")
	adminRoutes.HandleFunc("/users/{id}", authHandler.GetUserByID).Methods("GET")
	adminRoutes.HandleFunc("/users/{id}", authHandler.UpdateUser).Methods("PUT")
	adminRoutes.HandleFunc("/users/{id}", authHandler.DeleteUser).Methods("DELETE")

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
