package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	sharedContext "shared/context"
	"shared/models"
)

type StageRepository struct {
	db *sqlx.DB
}

func NewStageRepository(db *sqlx.DB) *StageRepository {
	return &StageRepository{db: db}
}

func (r *StageRepository) GetStagesPublic() ([]models.Stage, error) {
	query := `
        SELECT s.id, s.stage_offer_id, 
               s.establishment_id, s.content_id, s.status, s.start_date, s.end_date, 
               s.created_at, s.updated_at
        FROM stages s
        ORDER BY created_at DESC
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
			&stage.StageOfferID,
			&stage.EstablishmentID,
			&stage.ContentID,
			&stage.Status,
			&stage.StartDate,
			&stage.EndDate,
			&stage.CreatedAt,
			&stage.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erreur lors du scan: %v", err)
		}
		stages = append(stages, stage)
	}
	log.Println(stages)

	return stages, rows.Err()
}

func (r *StageRepository) GetStages(ctx context.Context) ([]models.Stage, error) {
	claims := sharedContext.GetClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("aucun claims trouvé dans le contexte")
	}
	userID := claims.UserID
	query := `
        SELECT s.id, s.stage_offer_id, s.student_id, s.teacher_id, s.tutor_id, 
               s.establishment_id, s.content_id, s.status, s.start_date, s.end_date, 
               s.created_at, s.updated_at
        FROM stages s
        JOIN users u ON s.student_id = u.id
		WHERE s.student_id = $1
    `
	rows, err := r.db.Query(query, userID)
	log.Println(rows)

	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des stages: %v", err)
	}
	defer rows.Close()

	var stages []models.Stage
	for rows.Next() {
		var stage models.Stage
		err := rows.Scan(
			&stage.ID,
			&stage.StageOfferID,
			&stage.StudentID,
			&stage.TeacherID,
			&stage.TutorID,
			&stage.EstablishmentID,
			&stage.ContentID,
			&stage.Status,
			&stage.StartDate,
			&stage.EndDate,
			&stage.CreatedAt,
			&stage.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erreur lors du scan: %v", err)
		}
		stages = append(stages, stage)
	}

	return stages, rows.Err()
}
func (r *StageRepository) UpdateStage(ctx context.Context, stage *models.Stage) error {
	query := `
        UPDATE stages 
        SET teacher_id = $2, tutor_id = $3, establishment_id = $4, 
            content_id = $5, status = $6, start_date = $7, 
            end_date = $8, updated_at = $9
        WHERE id = $1
    `

	_, err := r.db.ExecContext(ctx, query,
		stage.ID,
		stage.TeacherID,
		stage.TutorID,
		stage.EstablishmentID,
		stage.ContentID,
		stage.Status,
		stage.StartDate,
		stage.EndDate,
		stage.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("erreur lors de l'update du stage: %v", err)
	}

	return nil
}

func (r *StageRepository) GetStageByID(ctx context.Context, id int) (*models.Stage, error) {
	query := `
        SELECT id, stage_offer_id, student_id, teacher_id, tutor_id, 
               establishment_id, content_id, status, start_date, 
               end_date, created_at, updated_at
        FROM stages 
        WHERE id = $1
    `

	stage := &models.Stage{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&stage.ID,
		&stage.StageOfferID,
		&stage.StudentID,
		&stage.TeacherID,
		&stage.TutorID,
		&stage.EstablishmentID,
		&stage.ContentID,
		&stage.Status,
		&stage.StartDate,
		&stage.EndDate,
		&stage.CreatedAt,
		&stage.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erreur lors de la récupération du stage: %v", err)
	}

	return stage, nil
}

func (r *StageRepository) CreateStage(ctx context.Context, stage *models.Stage) (*models.Stage, error) {
	query := `
        INSERT INTO stages (
            stage_offer_id, student_id, teacher_id, tutor_id, 
            establishment_id, content_id, status, start_date, 
            end_date, created_at, updated_at
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
        ) RETURNING id, created_at, updated_at
    `

	err := r.db.QueryRowContext(ctx, query,
		stage.StageOfferID,
		stage.StudentID,
		stage.TeacherID,
		stage.TutorID,
		stage.EstablishmentID,
		stage.ContentID,
		stage.Status,
		stage.StartDate,
		stage.EndDate,
		stage.CreatedAt,
		stage.UpdatedAt,
	).Scan(
		&stage.ID,
		&stage.CreatedAt,
		&stage.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("erreur lors de l'insertion du stage: %v", err)
	}
	log.Println("stage: ", stage.ID, stage.CreatedAt, stage.UpdatedAt)
	return stage, nil
}

func (r *StageRepository) DeleteStage(ctx context.Context, id int) error {
	query := `DELETE FROM stages WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("erreur lors de la suppression du stage: %v", err)
	}

	// Vérifier qu'une ligne a bien été affectée
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erreur lors de la vérification de la suppression: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("stage not found")
	}

	return nil
}
