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

	// Créer les repositories
	stageRepo := repository.NewStageRepository(db)
	filiereRepo := repository.NewFiliereRepository(db)

	// Créer les handlers en leur passant les repos
	stageHandler := handlers.NewStageHandler(stageRepo)
	filiereHandler := handlers.NewFiliereHandler(filiereRepo)

	// Setup des routes avec les deux handlers
	r := setupRoutes(stageHandler, filiereHandler)

	fmt.Println("🚀 Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
