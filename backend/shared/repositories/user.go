package repositories

import (
	"context"
	"fmt"
	"shared/models"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByRole(ctx context.Context, role models.UserRole) ([]models.User, error) {
	query := `
        SELECT id, username, first_name, last_name, email, role, created_at, updated_at 
        FROM users 
        WHERE role = $1
        ORDER BY last_name, first_name
    `

	var users []models.User
	err := r.db.SelectContext(ctx, &users, query, role)
	if err != nil {
		return nil, fmt.Errorf("failed to get users by role %s: %w", role, err)
	}

	return users, nil
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
        INSERT INTO users (username, first_name, last_name, email, password_hash, role, created_at, updated_at)
        VALUES (:username, :first_name, :last_name, :email, :password_hash, :role, :created_at, :updated_at)
        RETURNING id
    `

	rows, err := r.db.NamedQueryContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&user.ID)
		if err != nil {
			return fmt.Errorf("failed to scan created user ID: %w", err)
		}
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	query := `
        SELECT id, username, first_name, last_name, email, role, created_at, updated_at 
        FROM users 
        WHERE id = $1
    `

	var user models.User
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID %d: %w", id, err)
	}

	return &user, nil
}
