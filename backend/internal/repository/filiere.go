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
func (r *FiliereRepository) CreateFiliere(filiere models.Filiere) (*models.Filiere, error) {
	query := `
        INSERT INTO filieres (code, label, color)
        VALUES ($1, $2, $3)
        RETURNING id, created_at, updated_at
    `

	var createdFiliere models.Filiere
	// Copier les données d'entrée
	createdFiliere.Code = filiere.Code
	createdFiliere.Label = filiere.Label
	createdFiliere.Color = filiere.Color

	err := r.db.QueryRow(query, filiere.Code, filiere.Label, filiere.Color).Scan(
		&createdFiliere.ID,
		&createdFiliere.CreatedAt,
		&createdFiliere.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la création de la filière: %v", err)
	}

	return &createdFiliere, nil
}

// UpdateFiliere met à jour une filière existante
func (r *FiliereRepository) UpdateFiliere(filiere models.Filiere) (*models.Filiere, error) {
	query := `
        UPDATE filieres 
        SET code = $1, label = $2, color = $3, updated_at = NOW()
        WHERE id = $4
        RETURNING id, code, label, color, created_at, updated_at
    `

	var updatedFiliere models.Filiere
	err := r.db.QueryRow(query, filiere.Code, filiere.Label, filiere.Color, filiere.ID).Scan(
		&updatedFiliere.ID,
		&updatedFiliere.Code,
		&updatedFiliere.Label,
		&updatedFiliere.Color,
		&updatedFiliere.CreatedAt,
		&updatedFiliere.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la mise à jour de la filière: %v", err)
	}

	return &updatedFiliere, nil
}

// DeleteFiliere supprime une filière
func (r *FiliereRepository) DeleteFiliere(id int) error {
	query := `DELETE FROM filieres WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erreur lors de la suppression de la filière: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erreur lors de la vérification de la suppression: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("aucune filière trouvée avec l'ID %d", id)
	}

	return nil
}
