package main

import (
	"context"
	"enterprise/internal/handlers"
	"enterprise/internal/repositories"
	"enterprise/internal/services"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"shared/config"
	"shared/middleware"
	sharedServices "shared/services"
	"time"
)

func main() {
	cfg := config.LoadConfig()
	// Initialiser la base de données
	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Erreur connexion DB: %v", err)
	}
	defer db.Close()

	// Initialiser les repositories
	userRepo := repositories.NewUserRepository(db)
	enterpriseRepo := repositories.NewEnterpriseRepository(db)

	// Initialiser les services
	userService := services.NewUserService(userRepo)
	enterpriseService := services.NewEnterpriseService(enterpriseRepo)

	// Initialiser le middleware d'authentification
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET must be set")
	}

	authService := sharedServices.NewAuthService(jwtSecret)
	authMiddleware := middleware.NewAuthMiddleware(authService)
	// Initialiser les handlers
	userHandler := handlers.NewUserHandler(userService)
	enterpriseHandler := handlers.NewEnterpriseHandler(enterpriseService)

	// Configurer les routes
	router := setupRoutes(userHandler, enterpriseHandler, authMiddleware)

	// Configuration du serveur
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Démarrage du serveur avec graceful shutdown
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server : %v", err)
		}
	}()

	// Attendre un signal d'arrêt
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown avec timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}
