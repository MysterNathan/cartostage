package main

import (
	"net/http"
	"shared/middleware"

	"github.com/gorilla/mux"
	"stage/internal/handlers"
)

func setupRoutes(
	stageHandler *handlers.StageHandler,
	filiereHandler *handlers.FiliereHandler,
	authMiddleware *middleware.AuthMiddleware,
) *mux.Router {
	r := mux.NewRouter()

	// Middleware CORS appliqué globalement
	r.Use(corsMiddleware)
	r.Use(authMiddleware.RequireAuth)

	api := r.PathPrefix("/api").Subrouter()

	// Routes stages - GET publiques, POST/PUT/DELETE protégées
	stagesRouter := api.PathPrefix("/stages").Subrouter()

	// Routes publiques (pas de middleware auth)
	stagesRouter.HandleFunc("", stageHandler.GetStages).Methods("GET")
	stagesRouter.HandleFunc("", stageHandler.CreateStage).Methods("POST")
	stagesRouter.HandleFunc("/{id:[0-9]+}", stageHandler.UpdateStage).Methods("PUT")
	stagesRouter.HandleFunc("/{id:[0-9]+}", stageHandler.DeleteStage).Methods("DELETE")

	// Routes protégées (avec middleware auth)
	protectedStages := stagesRouter.NewRoute().Subrouter()
	protectedStages.Use(authMiddleware.RequireAuth)
	//protectedStages.HandleFunc("", stageHandler.SaveStage).Methods("POST")
	//protectedStages.HandleFunc("/{id:[0-9]+}", stageHandler.UpdateStage).Methods("PUT")
	//protectedStages.HandleFunc("/{id:[0-9]+}", stageHandler.DeleteStage).Methods("DELETE")

	//// Routes filières public
	//filiereRouterPublic := api.PathPrefix("/filieres").Subrouter()
	//filiereRouterPublic.HandleFunc("", filiereHandler.GetFilieres).Methods("GET")

	//// Routes filieres (avec middleware auth)
	//filieresRouter := api.PathPrefix("/filieres").Subrouter()
	//filieresRouter.Use(authMiddleware.RequireAuth)
	//filieresRouter.HandleFunc("", filiereHandler.CreateFiliere).Methods("POST")
	//filieresRouter.HandleFunc("/{id:[0-9]+}", filiereHandler.UpdateFiliere).Methods("PUT")
	//filieresRouter.HandleFunc("/{id:[0-9]+}", filiereHandler.DeleteFiliere).Methods("DELETE")

	// Gestion OPTIONS pour toutes les routes
	stagesRouter.HandleFunc("", corsPreflightHandler).Methods("OPTIONS")
	//stagesRouter.HandleFunc("/filters", corsPreflightHandler).Methods("OPTIONS")
	stagesRouter.HandleFunc("/{id:[0-9]+}", corsPreflightHandler).Methods("OPTIONS")
	//filieresRouter.HandleFunc("", corsPreflightHandler).Methods("OPTIONS")
	//filieresRouter.HandleFunc("/{id:[0-9]+}", corsPreflightHandler).Methods("OPTIONS")

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
