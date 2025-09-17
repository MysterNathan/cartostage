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
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	issuer := os.Getenv("JWT_ISSUER")
	if issuer == "" {
		issuer = "auth-service"
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
	authRepository := repositories.NewAuthRepository(db)
	jwtService := sharedServices.NewJWTService(jwtSecret, issuer)
	authService := services.NewAuthService(authRepository, jwtService)
	authHandler := handlers.NewAuthHandler(authService)

	// Setup des routes
	r := setupRoutes(authHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	fmt.Printf("🔐 Auth service running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
