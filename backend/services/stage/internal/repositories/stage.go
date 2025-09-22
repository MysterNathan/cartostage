package repositories

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
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
        SELECT id, stage_offer_id, student_id, teacher_id, tutor_id, 
               establishment_id, content_id, status, start_date, end_date, 
               created_at, updated_at
        FROM stages 
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

func (r *StageRepository) GetStages(ctx context.Context) ([]models.Stage, error) {
	claims := sharedContext.GetClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("aucun claims trouvé dans le contexte")
	}
	userID := claims.UserID
	query := `
        SELECT id, stage_offer_id, student_id, teacher_id, tutor_id, 
               establishment_id, content_id, status, start_date, end_date, 
               created_at, updated_at
        FROM stages s
        JOIN public.users ON s.student_id = public.users.id
        WHERE student_id = $1

    `
	rows, err := r.db.Query(query, userID)
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
