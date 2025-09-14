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

	// Charger la config
	cfg := config.LoadConfig()

	// Connexion à la DB
	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Erreur connexion DB: %v", err)
	}
	defer db.Close()

	// Créer les repositories
	stageRepo := repositories.NewStageRepository(db)
	filiereRepo := repositories.NewFiliereRepository(db)

	// Créer les services
	authService := services.NewAuthService(jwtSecret)

	// Créer les handlers en leur passant les repos
	stageHandler := handlers.NewStageHandler(stageRepo)
	filiereHandler := handlers.NewFiliereHandler(filiereRepo)

	// Créer les middlewares en leur passant es services
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Setup des routes avec les deux handlers
	r := setupRoutes(stageHandler, filiereHandler, authMiddleware)

	fmt.Println("🚀 Server running on :80")
	log.Fatal(http.ListenAndServe(":80", r))

}
