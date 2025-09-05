package models

import "time"

type Stage struct {
	ID                int       `json:"id" db:"id"`
	Poste             string    `json:"poste" db:"poste"`
	Adresse           string    `json:"adresse" db:"adresse"`
	Lat               float64   `json:"lat" db:"lat"`
	Lng               float64   `json:"lng" db:"lng"`
	PlacesDisponibles int       `json:"placesDisponibles" db:"places_disponibles"`
	Entreprise        string    `json:"entreprise" db:"entreprise"`
	Filiere           string    `json:"filiere" db:"filiere"`
	Sector            string    `json:"sector" db:"sector"`
	Commune           string    `json:"commune" db:"commune"`
	CapacityTotal     int       `json:"capacity_total" db:"capacity_total"`
	CapacityFilled    int       `json:"capacity_filled" db:"capacity_filled"`
	Period            string    `json:"period" db:"period"`
	Parcours          string    `json:"parcours" db:"parcours"` // scolaire | apprentissage | mixte
	FamilleMetiers    string    `json:"famille_metiers" db:"famille_metiers"`
	NiveauScolaire    string    `json:"niveau_scolaire" db:"niveau_scolaire"` // 2de | 1re | Tle
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

type StagesData struct {
	Stages []Stage `json:"stages"`
}
