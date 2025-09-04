package repositories

import (
	"backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}

	query := `
        SELECT id, username, password_hash, role, created_at, updated_at 
        FROM users 
        WHERE username = $1
    `

	err := r.db.Get(user, query, username)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil // Utilisateur non trouvé
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetUserByID(id int) (*models.User, error) {
	user := &models.User{}

	query := `
        SELECT id, username, password_hash, role, created_at, updated_at 
        FROM users 
        WHERE id = $1
    `

	err := r.db.Get(user, query, id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetAllUsers() ([]*models.User, error) {
	users := []*models.User{}

	query := `
        SELECT id, username, password_hash, role, created_at, updated_at 
        FROM users 
        ORDER BY created_at DESC
    `

	err := r.db.Select(&users, query)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) CreateUser(user *models.User) error {
	query := `
        INSERT INTO users (username, password_hash, role) 
        VALUES (:username, :password_hash, :role) 
        RETURNING id, created_at, updated_at
    `

	rows, err := r.db.NamedQuery(query, user)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *UserRepository) UpdateUser(user *models.User) error {
	query := `
        UPDATE users 
        SET username = :username, password_hash = :password_hash, role = :role, updated_at = CURRENT_TIMESTAMP
        WHERE id = :id
        RETURNING updated_at
    `

	rows, err := r.db.NamedQuery(query, user)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&user.UpdatedAt)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *UserRepository) UpdatePassword(userID int, newPasswordHash string) error {
	query := `
        UPDATE users 
        SET password_hash = $1, updated_at = CURRENT_TIMESTAMP
        WHERE id = $2
    `

	_, err := r.db.Exec(query, newPasswordHash, userID)
	return err
}

func (r *UserRepository) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return nil // Ou retourner une erreur spécifique si souhaité
	}

	return nil
}

func (r *UserRepository) UserExists(username string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE username = $1`

	err := r.db.Get(&count, query, username)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
