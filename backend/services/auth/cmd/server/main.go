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
	sharedServices "shared/services"
)

func main() {
	// Variables d'environnement
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	issuer := os.Getenv("JWT_ISSUER")
	if issuer == "" {
		issuer = "auth-service" // valeur par défaut
	}

	// Charger la config
	cfg := config.LoadConfig()

	// Connexion à la DB
	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Erreur connexion DB: %v", err)
	}
	defer db.Close()

	// Créer le repository
	userRepo := repositories.NewUserRepository(db)

	// Créer le service JWT (shared service)
	jwtService := sharedServices.NewJWTService(jwtSecret, issuer)

	// Créer le service auth avec le JWT service
	authService := services.NewAuthService(userRepo, jwtService)

	// Créer les handlers
	authHandler := handlers.NewAuthHandler(authService)

	// Setup des routes avec le JWT service (pas authService)
	r := setupRoutes(authHandler, jwtService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "80" // port par défaut pour auth service
	}

	fmt.Printf("🚀 Auth service running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
