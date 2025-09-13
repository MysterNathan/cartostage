package main

import (
	"backend/services/auth/internal/handlers"
	"backend/services/auth/internal/repositories"
	"backend/services/auth/internal/services"
	"backend/shared/config"
	"fmt"
	"log"
	"net/http"
	"os"
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

	userRepo := repositories.NewUserRepository(db)

	// Créer les services
	authService := services.NewAuthService(jwtSecret, userRepo)

	// Créer les handlers
	authHandler := handlers.NewAuthHandler(authService)

	// Setup des routes avec les deux handlers
	r := setupRoutes(authHandler)

	fmt.Println("🚀 Auth micro services running")
	log.Fatal(http.ListenAndServe(":8080", r))
}
