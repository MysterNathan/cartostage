package main

import (
	"fmt"
	"log"
	"net/http"

	"backend/internal/config"
	"backend/internal/handlers"
	"backend/internal/repository"
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

	// Créer un repository
	stageRepo := repository.NewStageRepository(db)

	// Créer un handler en lui passant le repo
	stageHandler := handlers.NewStageHandler(stageRepo)

	// Setup des routes
	r := setupRoutes(stageHandler)

	fmt.Println("🚀 Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
