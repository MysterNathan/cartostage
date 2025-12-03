package repositories

import (
	"context"
	"fmt"
	"strings"

	"shared/models"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	GetAll(ctx context.Context, filter models.UserFilter) ([]models.User, error)
	GetByID(ctx context.Context, id int) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int) error
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetAll(ctx context.Context, filter models.UserFilter) ([]models.User, error) {
	query, args := r.buildQuery(filter)

	var users []models.User
	err := r.db.SelectContext(ctx, &users, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get teachers: %w", err)
	}

	return users, nil
}

func (r *userRepository) buildQuery(filter models.UserFilter) (string, []interface{}) {
	selectFields := `users.id, users.username, users.first_name, users.last_name, users.email, users.role, users.phone, users.establishment_id`
	baseQuery := fmt.Sprintf("SELECT DISTINCT %s FROM users", selectFields)

	var joins []string
	var conditions []string
	var args []interface{}
	argPosition := 1

	// ⭐ Condition commune : on ne veut QUE des enseignants
	conditions = append(conditions, "users.role = 'teacher'")

	switch filter.RequestorRole {
	case models.RoleAdmin:
		// Admin voit tous les enseignants (pas de filtre supplémentaire)

	case models.RoleTeacher:
		// Teacher ne voit QUE lui-même (pas les autres enseignants)
		conditions = append(conditions, fmt.Sprintf("users.id = $%d", argPosition))
		args = append(args, filter.RequestorID)
		argPosition++

	case models.RoleStudent:
		// Student voit son enseignant
		joins = append(joins, "JOIN stages ON stages.teacher_id = users.id")
		conditions = append(conditions, fmt.Sprintf("stages.student_id = $%d", argPosition))
		args = append(args, filter.RequestorID)
		argPosition++

	case models.RoleTutor:
		// Tutor voit les enseignants des élèves qu'il encadre
		joins = append(joins, "JOIN stages ON stages.teacher_id = users.id")
		conditions = append(conditions, fmt.Sprintf("stages.tutor_id = $%d", argPosition))
		args = append(args, filter.RequestorID)
		argPosition++
	}

	// Construction finale
	query := baseQuery
	if len(joins) > 0 {
		query += " " + strings.Join(joins, " ")
	}
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	return query, args
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	query := `SELECT * FROM users WHERE id = $1 AND role = 'teacher'`

	var user models.User
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get teacher: %w", err)
	}

	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	// Force le role à 'teacher'
	user.Role = string(models.RoleTeacher)

	query := `
        INSERT INTO users (username, first_name, last_name, email, password_hash, role, phone, establishment_id, is_active)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id, created_at, updated_at
    `

	err := r.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.FirstName,
		user.LastName,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.Phone,
		user.EstablishmentID,
		user.IsActive,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create teacher: %w", err)
	}

	return nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	query := `
        UPDATE users 
        SET username = $1, first_name = $2, last_name = $3, email = $4, 
            phone = $5, establishment_id = $6, is_active = $7, updated_at = CURRENT_TIMESTAMP
        WHERE id = $8 AND role = 'teacher'
    `

	result, err := r.db.ExecContext(
		ctx,
		query,
		user.Username,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Phone,
		user.EstablishmentID,
		user.IsActive,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update teacher: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("teacher not found")
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1 AND role = 'teacher'`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete teacher: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("teacher not found")
	}

	return nil
}
