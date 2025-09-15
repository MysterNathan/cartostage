package main

import (
	"enterprise/internal/handlers"
	"enterprise/internal/repositories"
	"enterprise/internal/services"
	"fmt"
	"log"
	"net/http"
	"os"
	"shared/config"
	"shared/middleware"
	sharedServices "shared/services"
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

	// Créer les repositories
	enterpriseRepo := repositories.NewEnterpriseRepository(db)
	tutorRepo := repositories.NewTutorRepository(db)

	// Créer les services
	enterpriseService := services.NewEnterpriseService(enterpriseRepo)
	tutorService := services.NewTutorService(tutorRepo)

	// Créer les handlers
	enterpriseHandler := handlers.NewEnterpriseHandler(enterpriseService)
	tutorHandler := handlers.NewTutorHandler(tutorService)

	jwtService := sharedServices.NewJWTService(jwtSecret, issuer)

	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	// Setup des routes
	r := setupRoutes(enterpriseHandler, tutorHandler, authMiddleware)

	fmt.Println("Enterprise microservice running")
	log.Fatal(http.ListenAndServe(":80", r))
}
