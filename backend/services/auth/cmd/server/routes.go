package main

import (
	"auth/internal/handlers"
	"github.com/gorilla/mux"
)

func setupRoutes(
	authHandler *handlers.AuthHandler,
) *mux.Router {
	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()

	// Route de login (non protégée)
	api.HandleFunc("/login", authHandler.Login).Methods("POST", "OPTIONS")

	return r
}
