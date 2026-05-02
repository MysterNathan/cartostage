package repositories

import (
	"github.com/jmoiron/sqlx"
	"shared/models"
)

type AuthRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

// FindUserByUsername avec sqlx - ULTRA SIMPLE
func (r *AuthRepository) FindUserByUsername(username string) (*models.User, error) {
	var user models.User

	query := `
        SELECT id, username, password_hash, role, first_name, last_name, 
               created_at, updated_at
        FROM users 
        WHERE username = $1
    `

	err := r.db.Get(&user, query, username)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
