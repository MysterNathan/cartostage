package repositories

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

// SaveStage - Sauvegarde un seul stage
func (r *StageRepository) SaveStage(stage *models.Stage) error {
	query := `
        INSERT INTO stages (
            poste, adresse, lat, lng, places_disponibles,
            entreprise, filiere, sector, commune, capacity_total,
            capacity_filled, period, created_at, updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW(), NOW())
        RETURNING id
    `

	err := r.db.QueryRow(query,
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
	).Scan(&stage.ID)

	if err != nil {
		return fmt.Errorf("erreur lors de l'insertion du stage: %v", err)
	}

	return nil
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

// DeleteStage - Supprime un stage par son ID
func (r *StageRepository) DeleteStage(id int) error {
	query := `DELETE FROM stages WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erreur lors de la suppression du stage: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erreur lors de la vérification de la suppression: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("aucun stage trouvé avec l'ID %d", id)
	}

	return nil
}

// DeleteAllStages - Supprime tous les stages
func (r *StageRepository) DeleteAllStages() error {
	query := `DELETE FROM stages`

	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("erreur lors de la suppression de tous les stages: %v", err)
	}

	return nil
}

// UpdateStage - Met à jour un stage existant
func (r *StageRepository) UpdateStage(stage *models.Stage) error {
	query := `
        UPDATE stages 
        SET poste = $2, adresse = $3, lat = $4, lng = $5, 
            places_disponibles = $6, entreprise = $7, filiere = $8, 
            sector = $9, commune = $10, capacity_total = $11, 
            capacity_filled = $12, period = $13, updated_at = CURRENT_TIMESTAMP
        WHERE id = $1
        RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(
		query,
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
	).Scan(&stage.ID, &stage.CreatedAt, &stage.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("aucun stage trouvé avec l'ID %d", stage.ID)
		}
		return fmt.Errorf("erreur lors de la mise à jour du stage: %v", err)
	}

	return nil
}
