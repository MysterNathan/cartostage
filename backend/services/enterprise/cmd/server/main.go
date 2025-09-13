package main

import (
	"backend/services/enterprise/internal/handlers"
	"backend/services/enterprise/internal/repositories"
	"backend/services/enterprise/internal/services"
	"backend/shared/config"
	"fmt"
	"log"
	"net/http"
)

func main() {
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

	// Setup des routes
	r := setupRoutes(enterpriseHandler, tutorHandler)

	fmt.Println("🚀 Enterprise microservice running")
	log.Fatal(http.ListenAndServe(":8080", r))
}
