package main

import (
	"context"
	"fmt"
	"form/internal/handlers"
	"form/internal/repositories"
	"form/internal/services"
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
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	issuer := os.Getenv("JWT_ISSUER")
	if issuer == "" {
		issuer = "stage-service"
	}

	// Charger la config
	cfg := config.LoadConfig()

	// Connexion à la DB
	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Erreur connexion DB: %v", err)
	}
	defer db.Close()

	stageSheetRepository := repositories.NewformRepository(db)
	stageSheetService := services.NewFormService(stageSheetRepository)
	stageSheetHandler := handlers.NewformHandler(stageSheetService)

	authService := sharedServices.NewAuthService(jwtSecret)
	authMiddleware := middleware.NewAuthMiddleware(authService)

	router := setupRoutes(authMiddleware, stageSheetHandler)

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
			log.Fatalf("Failed to start server: %v", err)
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
