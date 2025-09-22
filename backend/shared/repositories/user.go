package repositories

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	sharedContext "shared/context"
	"shared/models"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByRole(ctx context.Context, role models.UserRole) ([]models.User, error) {
	claims := sharedContext.GetUserClaims(ctx)
	query := `
        SELECT id, username, first_name, last_name, email, role, phone, 
               establishment_id, is_active, last_login, created_at, updated_at 
        FROM users 
        WHERE role = $1 
          AND establishment_id = (SELECT establishment_id FROM users WHERE id = $2)
        ORDER BY last_name, first_name
    `

	var users []models.User
	err := r.db.SelectContext(ctx, &users, query, role, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get users by role %s: %w", role, err)
	}

	return users, nil
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
        INSERT INTO users (username, first_name, last_name, email, password_hash, role, 
                          phone, establishment_id, is_active, created_at, updated_at)
        VALUES (:username, :first_name, :last_name, :email, :password_hash, :role,
                :phone, :establishment_id, :is_active, :created_at, :updated_at)
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
        SELECT id, username, first_name, last_name, email, role, phone, 
               establishment_id, is_active, last_login, created_at, updated_at 
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

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
        SELECT id, username, first_name, last_name, email, password_hash, role, phone, 
               establishment_id, is_active, last_login, created_at, updated_at 
        FROM users 
        WHERE username = $1
    `

	var user models.User
	err := r.db.GetContext(ctx, &user, query, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by username %s: %w", username, err)
	}

	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
        SELECT id, username, first_name, last_name, email, password_hash, role, phone, 
               establishment_id, is_active, last_login, created_at, updated_at 
        FROM users 
        WHERE email = $1
    `

	var user models.User
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email %s: %w", email, err)
	}

	return &user, nil
}

func (r *UserRepository) Delete(ctx context.Context, userID int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no user found with ID %d", userID)
	}

	return nil
}

func (r *UserRepository) Update(ctx context.Context, user *models.User, userID int) error {
	query := `
        UPDATE users
        SET username = :username,
            first_name = :first_name,
            last_name = :last_name,
            email = :email,
            password_hash = :password_hash,
            role = :role,
            phone = :phone,
            establishment_id = :establishment_id,
            is_active = :is_active,
            updated_at = :updated_at
        WHERE id = :id
        RETURNING id
    `

	params := map[string]interface{}{
		"id":               userID,
		"username":         user.Username,
		"first_name":       user.FirstName,
		"last_name":        user.LastName,
		"email":            user.Email,
		"password_hash":    user.PasswordHash,
		"role":             user.Role,
		"phone":            user.Phone,
		"establishment_id": user.EstablishmentID,
		"is_active":        user.IsActive,
		"updated_at":       user.UpdatedAt,
	}

	rows, err := r.db.NamedQueryContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&user.ID)
		if err != nil {
			return fmt.Errorf("failed to scan updated user ID: %w", err)
		}
	} else {
		return fmt.Errorf("no user found with id %d", userID)
	}

	return nil
}

// Méthode utilitaire pour mettre à jour le dernier login
func (r *UserRepository) UpdateLastLogin(ctx context.Context, userID int) error {
	query := `UPDATE users SET last_login = NOW() WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	return nil
}
