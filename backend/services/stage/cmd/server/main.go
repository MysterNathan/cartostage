package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"shared/config"
	"shared/middleware"
	"shared/services"
	"stage-map/internal/handlers"
	"stage-map/internal/repositories"
)

func main() {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	issuer := os.Getenv("JWT_ISSUER")
	if issuer == "" {
		issuer = "stage-service"
	}

	// Charger la config
	cfg := config.LoadConfig()

	// Connexion à la DB
	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Erreur connexion DB: %v", err)
	}
	defer db.Close()

	// === SIMPLE SETUP ===
	stageRepository := repositories.NewStageRepository(db)
	filiereRepository := repositories.NewFiliereRepository(db)

	stageHandler := handlers.NewStageHandler(stageRepository)
	filiereHandler := handlers.NewFiliereHandler(filiereRepository)

	authService := services.NewAuthService(jwtSecret)
	authMiddleware := middleware.NewAuthMiddleware(authService)
	// Setup des routes
	r := setupRoutes(stageHandler, filiereHandler, authMiddleware)

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	fmt.Printf("🗺️ Stage-map service running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
