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

func NewFormRepository(db *sqlx.DB) *FormRepository {
	return &FormRepository{db: db}
}

func (r *FormRepository) Get(ctx context.Context, id int) (*models.FormFormSection, error) {
	form, err := r.getForm(ctx, id)
	if err != nil {
		return nil, err
	}
	formSection, err := r.getFormSection(ctx, id)
	if err != nil {
		return nil, err
	}
	formFormSection := models.FormFormSection{
		Form: form, FormSection: formSection,
	}

	return &formFormSection, nil
}

func (r *FormRepository) getForm(ctx context.Context, id int) (*models.Form, error) {
	var form models.Form
	query := `
		SELECT 
		    id,
		    stage_id,
		    status,
		    content,
		    created_at,
		    updated_at,
		    completed_at
		    FROM form 
		    WHERE id = (SELECT form_id FROM form_section WHERE user_id = $1)
	`

	err := r.db.GetContext(ctx, &form, query, id)
	if err != nil {
		return nil, err
	}

	return &form, nil
}

func (r *FormRepository) getFormSection(ctx context.Context, id int) (*models.FormSection, error) {
	var formSection models.FormSection
	query := `
	SELECT 
	    id,
	    form_id,
	    section_type,
	    user_id,
	    status,
	    content,
	    created_at,
	    updated_at,
	    completed_at
	    FROM form_section WHERE user_id = $1
	`
	err := r.db.GetContext(ctx, &formSection, query, id)
	if err != nil {
		return nil, err
	}
	return &formSection, nil
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

func (r *FormRepository) createForm(ctx context.Context, data models.Form, applicantId int, tx *sqlx.Tx) (*models.Form, error) {
	var form models.Form

	setQuery := fmt.Sprintf("SET LOCAL app.user_id = %d", applicantId)
	_, err := tx.ExecContext(ctx, setQuery)
	if err != nil {
		fmt.Println("Erreur SET LOCAL:", err)
		return nil, err
	}

	query := `
        INSERT INTO form (
            stage_id,
            status,
            created_at,
            updated_at
        )
        VALUES ($1, $2, NOW(), NOW())
        RETURNING *
    `

	err = tx.GetContext(
		ctx,
		&form,
		query,
		data.StageID,
		"CREATED", // Default status CREATED
	)

	if err != nil {
		return nil, err
	}

	return &form, nil
}

func (r *FormRepository) createFormSection(ctx context.Context, data models.FormSection, applicantId int, tx *sqlx.Tx) (*models.FormSection, error) {
	var formSection models.FormSection

	// Set la variable de session pour RLS (sans paramètre préparé)
	setQuery := fmt.Sprintf("SET LOCAL app.user_id = %d", applicantId)
	_, err := tx.ExecContext(ctx, setQuery)
	if err != nil {
		fmt.Println("Erreur SET LOCAL:", err)
		return nil, err
	}

	query := `
        INSERT INTO form_section (
			form_id,
			section_type,
			user_id,
			status,
			created_at,
			updated_at
        )
        VALUES ($1, $2, $3, $4, NOW(), NOW())
        RETURNING *
    `

	err = tx.GetContext(
		ctx,
		&formSection,
		query,
		data.FormID,
		data.SectionType,
		data.UserID,
		"CREATED",
	)
	if err != nil {
		fmt.Println("Erreur INSERT:", err)
		return nil, err
	}

	return &formSection, nil
}

func (r *FormRepository) CreateForm(ctx context.Context, dataForm models.Form, dataFormSection []models.FormSection, applicantId int) error {

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		fmt.Println("Erreur BeginTxx:", err)
		return err
	}
	defer func() {
		if err != nil {
			fmt.Println("Rollback de la transaction")
			tx.Rollback()
		}
	}()

	form, err := r.createForm(ctx, dataForm, applicantId, tx)
	if err != nil {
		return err
	}

	for i := range dataFormSection {
		dataFormSection[i].FormID = form.ID
	}

	for i := range dataFormSection {
		_, err = r.createFormSection(ctx, dataFormSection[i], applicantId, tx)
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		fmt.Println("Erreur Commit:", err)
		return err
	}
	return nil
}
