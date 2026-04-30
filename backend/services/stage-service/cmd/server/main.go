package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"shared/config"
	"shared/middleware"
	sharedServices "shared/services"
	"stage/internal/handlers"
	"stage/internal/repositories"
	"stage/internal/services"
)

func main() {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	issuer := os.Getenv("JWT_ISSUER")
	if issuer == "" {
		issuer = "stage-service-service"
	}

	// Charger la config
	cfg := config.LoadConfig()

	// Connexion à la DB
	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Erreur connexion DB: %v", err)
	}
	defer db.Close()

	stageRepository := repositories.NewStageRepository(db)
	filiereRepository := repositories.NewFiliereRepository(db)
	formRepository := repositories.NewFormRepository(db)

	formService := services.NewFormService(formRepository)
	stageService := services.NewStageService(stageRepository, formService)

	stageHandler := handlers.NewStageHandler(stageService)
	filiereHandler := handlers.NewFiliereHandler(filiereRepository)
	formHandler := handlers.NewFormHandler(formService)
	authService := sharedServices.NewAuthService(jwtSecret)
	authMiddleware := middleware.NewAuthMiddleware(authService)
	// Setup des routes
	r := setupRoutes(stageHandler, filiereHandler, formHandler, authMiddleware)

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	fmt.Printf("🗺️ Stage service running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
