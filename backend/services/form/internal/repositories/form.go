package repositories

import (
	"context"
	"fmt"
	"shared/models"

	"github.com/jmoiron/sqlx"
)

type FormRepository struct {
	db *sqlx.DB
}

func NewformRepository(db *sqlx.DB) *FormRepository {
	return &FormRepository{db: db}
}

func (r *FormRepository) Get(ctx context.Context) ([]models.Form, error) {
	var forms []models.Form

	query := `
        SELECT 
            id,
            stage_id,
            student_id,
            teacher_id,
            tutor_id,
            status,
            content,
            created_at,
            updated_at,
            completed_at
        FROM form
        ORDER BY created_at DESC
    `

	err := r.db.SelectContext(ctx, &forms, query)
	if err != nil {
		return nil, err
	}

	return forms, nil
}

func (r *FormRepository) UpdateForm(ctx context.Context, data models.Form, applicantId int, formId int) ([]models.Form, error) {
	var forms []models.Form

	// Démarre une transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		fmt.Println("Erreur BeginTxx:", err)
		return nil, err
	}
	defer func() {
		if err != nil {
			fmt.Println("Rollback de la transaction")
			tx.Rollback()
		}
	}()

	// Set la variable de session pour RLS (sans paramètre préparé)
	setQuery := fmt.Sprintf("SET LOCAL app.user_id = %d", applicantId)
	_, err = tx.ExecContext(ctx, setQuery)
	if err != nil {
		fmt.Println("Erreur SET LOCAL:", err)
		return nil, err
	}

	// Exécute l'UPDATE et retourne les lignes mises à jour
	query := `
        UPDATE form
        SET
            status = $2,
            content = $3
        WHERE id = $1
        RETURNING *
    `

	err = tx.SelectContext(ctx, &forms, query, formId, data.Status, data.Content)
	if err != nil {
		fmt.Println("Erreur UPDATE:", err)
		return nil, err
	}

	// Commit la transaction
	if err = tx.Commit(); err != nil {
		fmt.Println("Erreur Commit:", err)
		return nil, err
	}

	return forms, nil
}

func (r *FormRepository) CreateForm(ctx context.Context, data models.Form, applicantId int) (*models.Form, error) {
	var form models.Form // Un seul Form, pas un slice

	// Démarre une transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		fmt.Println("Erreur BeginTxx:", err)
		return nil, err
	}
	defer func() {
		if err != nil {
			fmt.Println("Rollback de la transaction")
			tx.Rollback()
		}
	}()

	// Set la variable de session pour RLS (sans paramètre préparé)
	setQuery := fmt.Sprintf("SET LOCAL app.user_id = %d", applicantId)
	_, err = tx.ExecContext(ctx, setQuery)
	if err != nil {
		fmt.Println("Erreur SET LOCAL:", err)
		return nil, err
	}

	// Exécute l'INSERT et retourne la ligne créée
	query := `
        INSERT INTO form (
            stage_id,
            student_id,
            teacher_id,
            tutor_id,
            status,
            content,
            created_at,
            updated_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
        RETURNING *
    `

	err = tx.GetContext(
		ctx,
		&form, // Pas &forms
		query,
		data.StageID,
		data.StudentID,
		data.TeacherID,
		data.TutorID,
		data.Status,
		data.Content,
	)
	if err != nil {
		fmt.Println("Erreur INSERT:", err)
		return nil, err
	}

	// Commit la transaction
	if err = tx.Commit(); err != nil {
		fmt.Println("Erreur Commit:", err)
		return nil, err
	}

	return &form, nil // Retourne un pointeur
}
