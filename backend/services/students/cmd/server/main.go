package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"shared/config"
	"shared/middleware"
)

func main() {
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

	// Setup des routes
	r := setupRoutes(tutorGenericHandler, enterpriseHandler, authMiddleware)

	fmt.Println("Enterprise microservice running")
	log.Fatal(http.ListenAndServe(":80", r))
}
