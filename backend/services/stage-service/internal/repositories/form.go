package repositories

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"shared/models"
	"time"
)

type FormRepository struct {
	db *sqlx.DB
}

func NewFormRepository(db *sqlx.DB) *FormRepository {
	return &FormRepository{db: db}
}

func (r *FormRepository) Get(ctx context.Context, userID int, userRole models.UserRole) ([]*models.FormFormSection, error) {
	var column string
	switch userRole {
	case models.RoleStudent:
		column = "f.student_id"
	case models.RoleTeacher:
		column = "f.teacher_id"
	case models.RoleTutor:
		column = "f.tutor_id"
	default:
		return nil, fmt.Errorf("rôle non supporté: %s", userRole)
	}

	query := fmt.Sprintf(`
        SELECT 
            f.id, f.stage_id, f.status, f.content,
            f.created_at, f.updated_at, f.completed_at,
            fs.id AS fs_id, fs.form_id, fs.section_type, fs.user_id,
            fs.status AS fs_status, fs.content AS fs_content,
            fs.created_at AS fs_created_at, fs.updated_at AS fs_updated_at,
            fs.completed_at AS fs_completed_at
        FROM form f
        INNER JOIN form_section fs ON fs.form_id = f.id
        WHERE %s = $1
        ORDER BY f.id
    `, column)

	rows, err := r.db.QueryxContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return mapFormRows(rows)
}
func mapFormRows(rows *sqlx.Rows) ([]*models.FormFormSection, error) {
	formsMap := make(map[int]*models.FormFormSection)
	var orderedIDs []int

	for rows.Next() {
		var row struct {
			// Form
			ID          int          `db:"id"`
			StageID     int          `db:"stage_id"`
			Status      string       `db:"status"`
			Content     models.JSONB `db:"content"`
			CreatedAt   time.Time    `db:"created_at"`
			UpdatedAt   time.Time    `db:"updated_at"`
			CompletedAt *time.Time   `db:"completed_at"`
			// FormSection
			FsID          int          `db:"fs_id"`
			FormID        int          `db:"form_id"`
			SectionType   string       `db:"section_type"`
			UserID        int          `db:"user_id"`
			FsStatus      string       `db:"fs_status"`
			FsContent     models.JSONB `db:"fs_content"`
			FsCreatedAt   time.Time    `db:"fs_created_at"`
			FsUpdatedAt   time.Time    `db:"fs_updated_at"`
			FsCompletedAt *time.Time   `db:"fs_completed_at"`
		}

		if err := rows.StructScan(&row); err != nil {
			return nil, err
		}

		if _, exists := formsMap[row.ID]; !exists {
			formsMap[row.ID] = &models.FormFormSection{
				Form: &models.Form{
					ID:          row.ID,
					StageID:     row.StageID,
					Status:      row.Status,
					Content:     &row.Content,
					CreatedAt:   row.CreatedAt,
					UpdatedAt:   row.UpdatedAt,
					CompletedAt: row.CompletedAt,
				},
				FormSections: []models.FormSection{},
			}
			orderedIDs = append(orderedIDs, row.ID)
		}

		formsMap[row.ID].FormSections = append(formsMap[row.ID].FormSections, models.FormSection{
			ID:          row.FsID,
			FormID:      row.FormID,
			SectionType: row.SectionType,
			UserID:      row.UserID,
			Status:      row.FsStatus,
			Content:     &row.FsContent,
			CreatedAt:   row.FsCreatedAt,
			UpdatedAt:   row.FsUpdatedAt,
			CompletedAt: row.FsCompletedAt,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Préserve l'ordre des résultats
	result := make([]*models.FormFormSection, 0, len(orderedIDs))
	for _, id := range orderedIDs {
		result = append(result, formsMap[id])
	}

	return result, nil
}

func (r *FormRepository) UpdateForm(ctx context.Context, data *models.Form, applicantId int) (*models.Form, error) {
	var form models.Form

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

	setQuery := fmt.Sprintf("SET LOCAL app.user_id = %d", applicantId)
	_, err = tx.ExecContext(ctx, setQuery)
	if err != nil {
		fmt.Println("Erreur SET LOCAL:", err)
		return nil, err
	}

	query := `
        UPDATE form
        SET
            status = $2,
            content = $3
        WHERE id = $1
        RETURNING *
    `

	err = tx.GetContext(ctx, &form, query, data.ID, data.Status, data.Content)
	if err != nil {
		fmt.Println("Erreur UPDATE Form:", err)
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		fmt.Println("Erreur Commit:", err)
		return nil, err
	}

	return &form, nil
}

func (r *FormRepository) UpdateFormSection(ctx context.Context, data *models.FormSection, applicantId int) ([]models.FormSection, error) {
	var forms []models.FormSection

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
        UPDATE form_section
        SET
            status = $2,
            content = $3
        WHERE id = $1
        RETURNING *
    `

	err = tx.SelectContext(ctx, &forms, query, data.ID, data.Status, data.Content)
	if err != nil {
		fmt.Println("Erreur UPDATE FormSection:", err)
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
        	student_id,
        	teacher_id,
        	tutor_id,
            created_at,
            updated_at
        )
        VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
        RETURNING *
    `

	err = tx.GetContext(
		ctx,
		&form,
		query,
		data.StageID,
		"CREATED", // Default status CREATED
		data.StudentID,
		data.TeacherID,
		data.TutorID,
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

	fmt.Printf("Form: %+v\n", dataForm)
	fmt.Printf("Sections: %+v\n", dataFormSection)
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		fmt.Println("Erreur BeginTxx:", err)
		return err
	}
	defer func() {
		if err != nil {
			fmt.Println("Rollback de la transaction", err)
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
