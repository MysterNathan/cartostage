package main

import (
	"backend/internal/middleware"
	"backend/internal/services"
	"fmt"
	"log"
	"net/http"
	"os"

	"backend/internal/config"
	"backend/internal/handlers"
	"backend/internal/repositories"
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
	userRepo := repositories.NewUserRepository(db)

	// Créer les services
	authService := services.NewAuthService(jwtSecret, userRepo)

	// Créer les handlers en leur passant les repos
	stageHandler := handlers.NewStageHandler(stageRepo)
	filiereHandler := handlers.NewFiliereHandler(filiereRepo)
	authHandler := handlers.NewAuthHandler(authService)

	// Créer les middlewares en leur passant es services
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Setup des routes avec les deux handlers
	r := setupRoutes(stageHandler, filiereHandler, authHandler, authMiddleware)

	fmt.Println("🚀 Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))

}
