package repository

import (
	"backend/internal/models"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type FiliereRepository struct {
	db *sqlx.DB
}

func NewFiliereRepository(db *sqlx.DB) *FiliereRepository {
	return &FiliereRepository{db: db}
}

// GetFilieres récupère toutes les filières
func (r *FiliereRepository) GetFilieres() (*models.FilieresData, error) {
	query := `
        SELECT id, code, label, color, created_at, updated_at
        FROM filieres 
        ORDER BY code
    `

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des filières: %v", err)
	}
	defer rows.Close()

	var filieres []models.Filiere
	for rows.Next() {
		var filiere models.Filiere
		err := rows.Scan(
			&filiere.ID,
			&filiere.Code,
			&filiere.Label,
			&filiere.Color,
			&filiere.CreatedAt,
			&filiere.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erreur lors du scan: %v", err)
		}
		filieres = append(filieres, filiere)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erreur lors de l'itération: %v", err)
	}

	return &models.FilieresData{Filieres: filieres}, nil
}

// GetFiliereByCode récupère une filière par son code
func (r *FiliereRepository) GetFiliereByCode(code string) (*models.Filiere, error) {
	query := `
        SELECT id, code, label, color, created_at, updated_at
        FROM filieres 
        WHERE code = $1
    `

	var filiere models.Filiere
	err := r.db.QueryRow(query, code).Scan(
		&filiere.ID,
		&filiere.Code,
		&filiere.Label,
		&filiere.Color,
		&filiere.CreatedAt,
		&filiere.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("filière non trouvée: %v", err)
	}

	return &filiere, nil
}

// CreateFiliere crée une nouvelle filière
func (r *FiliereRepository) CreateFiliere(filiere *models.Filiere) error {
	query := `
        INSERT INTO filieres (code, label, color)
        VALUES ($1, $2, $3)
        RETURNING id, created_at, updated_at
    `

	err := r.db.QueryRow(query, filiere.Code, filiere.Label, filiere.Color).Scan(
		&filiere.ID,
		&filiere.CreatedAt,
		&filiere.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("erreur lors de la création de la filière: %v", err)
	}

	return nil
}

// UpdateFiliere met à jour une filière existante
func (r *FiliereRepository) UpdateFiliere(filiere *models.Filiere) error {
	query := `
        UPDATE filieres 
        SET label = $2, color = $3
        WHERE code = $1
        RETURNING updated_at
    `

	err := r.db.QueryRow(query, filiere.Code, filiere.Label, filiere.Color).Scan(
		&filiere.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("erreur lors de la mise à jour de la filière: %v", err)
	}

	return nil
}

// DeleteFiliere supprime une filière
func (r *FiliereRepository) DeleteFiliere(code string) error {
	query := `DELETE FROM filieres WHERE code = $1`

	result, err := r.db.Exec(query, code)
	if err != nil {
		return fmt.Errorf("erreur lors de la suppression de la filière: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erreur lors de la vérification de la suppression: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("filière non trouvée pour suppression")
	}

	return nil
}
