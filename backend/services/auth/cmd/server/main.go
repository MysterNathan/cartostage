package main

import (
	"auth/internal/handlers"
	"auth/internal/repositories"
	"auth/internal/services"
	"fmt"
	"log"
	"net/http"
	"os"
	"shared/config"
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
	log.Fatal(http.ListenAndServe(":80", r))
}
