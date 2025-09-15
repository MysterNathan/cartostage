package repositories

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"shared/models"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetByUsername - Récupère un utilisateur par nom (avec hash pour auth)
func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	query := `
        SELECT id, username, email, password_hash, role, entity_id, created_at, updated_at, last_login
        FROM users 
        WHERE username = $1
    `
	err := r.db.Get(&user, query, username)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByID - Récupère un utilisateur par ID (avec hash pour auth interne)
func (r *UserRepository) GetByID(id int) (*models.User, error) {
	var user models.User
	query := `
        SELECT id, username, email, password_hash, role, entity_id, created_at, updated_at, last_login
        FROM users 
        WHERE id = $1
    `
	err := r.db.Get(&user, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser - Crée un utilisateur (attend un User avec hash)
func (r *UserRepository) CreateUser(user *models.User) error {
	query := `
        INSERT INTO users (username, email, password_hash, role, entity_id)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at, updated_at
    `
	err := r.db.QueryRow(query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.EntityID,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	return err
}

// UpdateUser - Met à jour un utilisateur
func (r *UserRepository) UpdateUser(user *models.User) error {
	query := `
        UPDATE users 
        SET username = $2, email = $3, password_hash = $4, role = $5, entity_id = $6, updated_at = CURRENT_TIMESTAMP
        WHERE id = $1
        RETURNING updated_at
    `
	err := r.db.QueryRow(query,
		user.ID,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.EntityID,
	).Scan(&user.UpdatedAt)

	return err
}

// UpdateLastLogin - Met à jour la dernière connexion
func (r *UserRepository) UpdateLastLogin(userID int) error {
	query := `UPDATE users SET last_login = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := r.db.Exec(query, userID)
	return err
}

// UserExists - Vérifie si un utilisateur existe
func (r *UserRepository) UserExists(username string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE username = $1`
	err := r.db.Get(&count, query, username)
	return count > 0, err
}

// DeleteUser - Supprime un utilisateur
func (r *UserRepository) DeleteUser(userID int) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.db.Exec(query, userID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return fmt.Errorf("user with id %d not found", userID)
	}

	return nil
}

// GetUsers - Liste les utilisateurs (pour admin) - retourne des UserProfile
func (r *UserRepository) GetUsers(page, limit int, roleFilter, entityIDFilter string) ([]*models.UserProfile, int, error) {
	offset := (page - 1) * limit
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argCount := 0

	if roleFilter != "" {
		argCount++
		whereClause += fmt.Sprintf(" AND role = $%d", argCount)
		args = append(args, roleFilter)
	}

	if entityIDFilter != "" {
		argCount++
		whereClause += fmt.Sprintf(" AND entity_id = $%d", argCount)
		args = append(args, entityIDFilter)
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM users %s", whereClause)
	var total int
	err := r.db.Get(&total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// Get users (sans password_hash pour les listes)
	argCount++
	limitArg := argCount
	argCount++
	offsetArg := argCount

	query := fmt.Sprintf(`
        SELECT id, username, email, role, entity_id, created_at, updated_at, last_login
        FROM users %s
        ORDER BY created_at DESC
        LIMIT $%d OFFSET $%d
    `, whereClause, limitArg, offsetArg)

	args = append(args, limit, offset)

	var users []*models.UserProfile
	err = r.db.Select(&users, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
