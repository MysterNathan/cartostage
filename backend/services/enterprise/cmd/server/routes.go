package main

import (
	"backend/services/enterprise/internal/handlers"

	"github.com/gorilla/mux"
)

func setupRoutes(enterpriseHandler *handlers.EnterpriseHandler, tutorHandler *handlers.TutorHandler) *mux.Router {
	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()

	// Routes pour les entreprises
	api.HandleFunc("/enterprises", enterpriseHandler.GetAll).Methods("GET")
	api.HandleFunc("/enterprises/{id}", enterpriseHandler.GetByID).Methods("GET")
	api.HandleFunc("/enterprises", enterpriseHandler.Create).Methods("POST")
	api.HandleFunc("/enterprises/{id}", enterpriseHandler.Update).Methods("PUT")
	api.HandleFunc("/enterprises/{id}", enterpriseHandler.Delete).Methods("DELETE")
	api.HandleFunc("/enterprises/{id}/stats", enterpriseHandler.GetWithStats).Methods("GET")

	// Routes pour les tuteurs
	api.HandleFunc("/tutors", tutorHandler.GetAll).Methods("GET")
	api.HandleFunc("/tutors/{id}", tutorHandler.GetByID).Methods("GET")
	api.HandleFunc("/tutors", tutorHandler.Create).Methods("POST")
	api.HandleFunc("/tutors/{id}", tutorHandler.Update).Methods("PUT")
	api.HandleFunc("/tutors/{id}", tutorHandler.Delete).Methods("DELETE")
	api.HandleFunc("/tutors/enterprise/{enterprise_id}", tutorHandler.GetByEnterprise).Methods("GET")
	api.HandleFunc("/tutors/{id}/stats", tutorHandler.GetWithStats).Methods("GET")

	return r
}
