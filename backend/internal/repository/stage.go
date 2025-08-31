package repository

import (
	"backend/internal/models"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

type StageRepository struct {
	db *sqlx.DB
}

func NewStageRepository(db *sqlx.DB) *StageRepository {
	return &StageRepository{db: db}
}

// GetAllStages - Équivalent du GET NextJS
func (r *StageRepository) GetAllStages() (*models.StagesData, error) {
	query := `
        SELECT id, poste, adresse, lat, lng, places_disponibles, 
               entreprise, filiere, sector, commune, capacity_total, 
               capacity_filled, period, created_at, updated_at
        FROM stages 
        ORDER BY commune, entreprise
    `

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des stages: %v", err)
	}
	defer rows.Close()

	var stages []models.Stage
	for rows.Next() {
		var stage models.Stage
		err := rows.Scan(
			&stage.ID,
			&stage.Poste,
			&stage.Adresse,
			&stage.Lat,
			&stage.Lng,
			&stage.PlacesDisponibles,
			&stage.Entreprise,
			&stage.Filiere,
			&stage.Sector,
			&stage.Commune,
			&stage.CapacityTotal,
			&stage.CapacityFilled,
			&stage.Period,
			&stage.CreatedAt,
			&stage.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erreur lors du scan: %v", err)
		}
		stages = append(stages, stage)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erreur lors de l'itération: %v", err)
	}

	return &models.StagesData{Stages: stages}, nil
}

// SaveAllStages - Équivalent du POST NextJS (mise à jour complète)
func (r *StageRepository) SaveAllStages(stagesData *models.StagesData) error {
	// Démarrer une transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("erreur lors du démarrage de la transaction: %v", err)
	}
	defer tx.Rollback()

	// Vider la table existante
	_, err = tx.Exec("DELETE FROM stages")
	if err != nil {
		return fmt.Errorf("erreur lors de la suppression: %v", err)
	}

	// Insérer les nouveaux stages
	query := `
        INSERT INTO stages (
            id, poste, adresse, lat, lng, places_disponibles,
            entreprise, filiere, sector, commune, capacity_total,
            capacity_filled, period
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
    `

	for _, stage := range stagesData.Stages {
		_, err = tx.Exec(query,
			stage.ID,
			stage.Poste,
			stage.Adresse,
			stage.Lat,
			stage.Lng,
			stage.PlacesDisponibles,
			stage.Entreprise,
			stage.Filiere,
			stage.Sector,
			stage.Commune,
			stage.CapacityTotal,
			stage.CapacityFilled,
			stage.Period,
		)
		if err != nil {
			return fmt.Errorf("erreur lors de l'insertion du stage %d: %v", stage.ID, err)
		}
	}

	// Réinitialiser la séquence pour les prochains ID
	if len(stagesData.Stages) > 0 {
		maxID := 0
		for _, stage := range stagesData.Stages {
			if stage.ID > maxID {
				maxID = stage.ID
			}
		}
		_, err = tx.Exec("SELECT setval('stages_id_seq', $1)", maxID)
		if err != nil {
			return fmt.Errorf("erreur lors de la réinitialisation de la séquence: %v", err)
		}
	}

	// Valider la transaction
	return tx.Commit()
}

// Méthodes additionnelles pour des opérations plus granulaires
func (r *StageRepository) GetStageByID(id int) (*models.Stage, error) {
	query := `
        SELECT id, poste, adresse, lat, lng, places_disponibles, 
               entreprise, filiere, sector, commune, capacity_total, 
               capacity_filled, period, created_at, updated_at
        FROM stages 
        WHERE id = $1
    `

	var stage models.Stage
	err := r.db.QueryRow(query, id).Scan(
		&stage.ID,
		&stage.Poste,
		&stage.Adresse,
		&stage.Lat,
		&stage.Lng,
		&stage.PlacesDisponibles,
		&stage.Entreprise,
		&stage.Filiere,
		&stage.Sector,
		&stage.Commune,
		&stage.CapacityTotal,
		&stage.CapacityFilled,
		&stage.Period,
		&stage.CreatedAt,
		&stage.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("stage avec l'ID %d introuvable", id)
		}
		return nil, fmt.Errorf("erreur lors de la récupération du stage: %v", err)
	}

	return &stage, nil
}

func (r *StageRepository) GetStagesWithFilters(filiere, commune string, availableOnly bool) (*models.StagesData, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	baseQuery := `
        SELECT id, poste, adresse, lat, lng, places_disponibles, 
               entreprise, filiere, sector, commune, capacity_total, 
               capacity_filled, period, created_at, updated_at
        FROM stages
    `

	// Construire les conditions WHERE
	if filiere != "" {
		conditions = append(conditions, fmt.Sprintf("filiere = $%d", argIndex))
		args = append(args, filiere)
		argIndex++
	}

	if commune != "" {
		conditions = append(conditions, fmt.Sprintf("commune = $%d", argIndex))
		args = append(args, commune)
		argIndex++
	}

	if availableOnly {
		conditions = append(conditions, "places_disponibles > 0")
	}

	// Ajouter les conditions à la requête
	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	baseQuery += " ORDER BY commune, entreprise"

	rows, err := r.db.Query(baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des stages filtrés: %v", err)
	}
	defer rows.Close()

	var stages []models.Stage
	for rows.Next() {
		var stage models.Stage
		err := rows.Scan(
			&stage.ID,
			&stage.Poste,
			&stage.Adresse,
			&stage.Lat,
			&stage.Lng,
			&stage.PlacesDisponibles,
			&stage.Entreprise,
			&stage.Filiere,
			&stage.Sector,
			&stage.Commune,
			&stage.CapacityTotal,
			&stage.CapacityFilled,
			&stage.Period,
			&stage.CreatedAt,
			&stage.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erreur lors du scan: %v", err)
		}
		stages = append(stages, stage)
	}

	return &models.StagesData{Stages: stages}, nil
}
