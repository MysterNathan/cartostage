package main

import (
	"enterprise/internal/handlers"
	"enterprise/internal/repositories"
	"enterprise/internal/services"

	//"enterprise/internal/handlers"
	//"enterprise/internal/repositories"
	//"enterprise/internal/services"
	"fmt"
	"log"
	"net/http"
	"os"
	"shared/config"
	sharedHandler "shared/handlers"
	"shared/middleware"
	sharedRepositories "shared/repositories"
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

	// === REPOSITORIES ===
	// Repository entreprise (existant)
	//enterpriseRepo := repositories.NewEnterpriseRepository(db)

	// Repository générique pour les utilisateurs/tuteurs
	tutorGenericRepo := sharedRepositories.NewUserRepository(db)
	enterpriseRepo := repositories.NewEnterpriseRepository(db)
	// === SERVICES ===
	// Service entreprise (existant)
	//enterpriseService := services.NewEnterpriseService(enterpriseRepo)

	// Service générique pour les tuteurs
	tutorGenericService := sharedServices.NewUserService(tutorGenericRepo)
	enterpriseService := services.NewEnterpriseService(enterpriseRepo)
	// === HANDLERS ===
	// Handler entreprise (existant)
	//enterpriseHandler := handlers.NewEnterpriseHandler(enterpriseService)

	// Handler générique pour les tuteurs
	tutorGenericHandler := sharedHandler.NewUserHandler(tutorGenericService)

	// Services partagés
	authService := sharedServices.NewAuthService(jwtSecret)
	authMiddleware := middleware.NewAuthMiddleware(authService)
	enterpriseHandler := handlers.NewEnterpriseHandler(enterpriseService)

	// Setup des routes
	r := setupRoutes(tutorGenericHandler, enterpriseHandler, authMiddleware)

	fmt.Println("Enterprise microservice running")
	log.Fatal(http.ListenAndServe(":80", r))
}
