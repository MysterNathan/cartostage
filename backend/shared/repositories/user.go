package repositories

import (
	"database/sql"
	"fmt"
	"shared/models"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetAll récupère tous les utilisateurs avec leurs profils
func (r *UserRepository) GetAll() ([]*models.User, error) {
	query := `
        SELECT 
            u.id, u.username, u.first_name, u.last_name, u.email, u.password_hash, 
            u.role, u.entity_type, u.entity_id, u.is_active, u.email_verified, 
            u.last_login, u.created_at, u.updated_at,
            p.phone, p.poste, p.departement,
            p.is_active as profile_is_active, p.created_at as profile_created_at, 
            p.updated_at as profile_updated_at
        FROM users u
        LEFT JOIN user_profiles p ON u.id = p.user_id
        ORDER BY u.created_at DESC
    `

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanUsers(rows)
}

// GetByID récupère un utilisateur par son ID
func (r *UserRepository) GetByID(id int) (*models.User, error) {
	query := `
        SELECT 
            u.id, u.username, u.first_name, u.last_name, u.email, u.password_hash, 
            u.role, u.entity_type, u.entity_id, u.is_active, u.email_verified, 
            u.last_login, u.created_at, u.updated_at,
            p.phone, p.poste, p.departement,
            p.is_active as profile_is_active, p.created_at as profile_created_at, 
            p.updated_at as profile_updated_at
        FROM users u
        LEFT JOIN user_profiles p ON u.id = p.user_id
        WHERE u.id = $1
    `

	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users, err := r.scanUsers(rows)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, sql.ErrNoRows
	}

	return users[0], nil
}

// GetByUsername récupère un utilisateur par son nom d'utilisateur
func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	query := `
        SELECT 
            u.id, u.username, u.first_name, u.last_name, u.email, u.password_hash, 
            u.role, u.entity_type, u.entity_id, u.is_active, u.email_verified, 
            u.last_login, u.created_at, u.updated_at,
            p.phone, p.poste, p.departement,
            p.is_active as profile_is_active, p.created_at as profile_created_at, 
            p.updated_at as profile_updated_at
        FROM users u
        LEFT JOIN user_profiles p ON u.id = p.user_id
        WHERE u.username = $1
    `

	rows, err := r.db.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users, err := r.scanUsers(rows)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, sql.ErrNoRows
	}

	return users[0], nil
}

// GetByEmail récupère un utilisateur par son email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `
        SELECT 
            u.id, u.username, u.first_name, u.last_name, u.email, u.password_hash, 
            u.role, u.entity_type, u.entity_id, u.is_active, u.email_verified, 
            u.last_login, u.created_at, u.updated_at,
            p.phone, p.poste, p.departement,
            p.is_active as profile_is_active, p.created_at as profile_created_at, 
            p.updated_at as profile_updated_at
        FROM users u
        LEFT JOIN user_profiles p ON u.id = p.user_id
        WHERE u.email = $1
    `

	rows, err := r.db.Query(query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users, err := r.scanUsers(rows)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, sql.ErrNoRows
	}

	return users[0], nil
}

// GetByRole récupère les utilisateurs par rôle
func (r *UserRepository) GetByRole(role string) ([]*models.User, error) {
	query := `
        SELECT 
            u.id, u.username, u.first_name, u.last_name, u.email, u.password_hash, 
            u.role, u.entity_type, u.entity_id, u.is_active, u.email_verified, 
            u.last_login, u.created_at, u.updated_at,
            p.phone, p.poste, p.departement,
            p.is_active as profile_is_active, p.created_at as profile_created_at, 
            p.updated_at as profile_updated_at
        FROM users u
        LEFT JOIN user_profiles p ON u.id = p.user_id
        WHERE u.role = $1
        ORDER BY u.created_at DESC
    `

	rows, err := r.db.Query(query, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanUsers(rows)
}

// GetByEntity récupère les utilisateurs par entité
func (r *UserRepository) GetByEntity(entityType string, entityID int) ([]*models.User, error) {
	query := `
        SELECT 
            u.id, u.username, u.first_name, u.last_name, u.email, u.password_hash, 
            u.role, u.entity_type, u.entity_id, u.is_active, u.email_verified, 
            u.last_login, u.created_at, u.updated_at,
            p.phone, p.poste, p.departement,
            p.is_active as profile_is_active, p.created_at as profile_created_at, 
            p.updated_at as profile_updated_at
        FROM users u
        LEFT JOIN user_profiles p ON u.id = p.user_id
        WHERE u.entity_type = $1 AND u.entity_id = $2
        ORDER BY u.created_at DESC
    `

	rows, err := r.db.Query(query, entityType, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanUsers(rows)
}

// Create crée un nouvel utilisateur avec son profil
func (r *UserRepository) Create(req models.CreateUserRequest, passwordHash string) (*models.User, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Insérer l'utilisateur principal
	var userID int
	query := `
        INSERT INTO users (username, first_name, last_name, email, password_hash, role, entity_type, entity_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id
    `

	err = tx.QueryRow(query, req.Username, req.FirstName, req.LastName, req.Email,
		passwordHash, req.Role, req.EntityType, req.EntityID).Scan(&userID)
	if err != nil {
		return nil, err
	}

	// Créer le profil si des données de profil sont fournies
	if req.Phone != nil || req.Poste != nil || req.Departement != nil {
		insertQuery := `
            INSERT INTO user_profiles (user_id, phone, poste, departement, is_active, created_at, updated_at)
            VALUES ($1, $2, $3, $4, $5, $6, $7)
        `
		now := time.Now()
		_, err = tx.Exec(insertQuery, userID, req.Phone, req.Poste, req.Departement, true, now, now)
		if err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return r.GetByID(userID)
}

// Update met à jour un utilisateur et son profil
func (r *UserRepository) Update(userID int, req models.UpdateUserRequest) (*models.User, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Mettre à jour les champs de l'utilisateur
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Username != nil {
		setParts = append(setParts, fmt.Sprintf("username = $%d", argIndex))
		args = append(args, *req.Username)
		argIndex++
	}
	if req.FirstName != nil {
		setParts = append(setParts, fmt.Sprintf("first_name = $%d", argIndex))
		args = append(args, *req.FirstName)
		argIndex++
	}
	if req.LastName != nil {
		setParts = append(setParts, fmt.Sprintf("last_name = $%d", argIndex))
		args = append(args, *req.LastName)
		argIndex++
	}
	if req.Email != nil {
		setParts = append(setParts, fmt.Sprintf("email = $%d", argIndex))
		args = append(args, *req.Email)
		argIndex++
	}
	if req.Role != nil {
		setParts = append(setParts, fmt.Sprintf("role = $%d", argIndex))
		args = append(args, *req.Role)
		argIndex++
	}
	if req.EntityType != nil {
		setParts = append(setParts, fmt.Sprintf("entity_type = $%d", argIndex))
		args = append(args, *req.EntityType)
		argIndex++
	}
	if req.EntityID != nil {
		setParts = append(setParts, fmt.Sprintf("entity_id = $%d", argIndex))
		args = append(args, *req.EntityID)
		argIndex++
	}
	if req.IsActive != nil {
		setParts = append(setParts, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, *req.IsActive)
		argIndex++
	}

	if len(setParts) > 0 {
		setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
		args = append(args, time.Now())
		argIndex++

		args = append(args, userID)

		query := fmt.Sprintf(`
            UPDATE users 
            SET %s 
            WHERE id = $%d
        `, strings.Join(setParts, ", "), argIndex)

		_, err = tx.Exec(query, args...)
		if err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return r.GetByID(userID)
}

// UpdateProfile met à jour uniquement le profil utilisateur
func (r *UserRepository) UpdateProfile(userID int, req models.UpdateUserProfileRequest) (*models.User, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Vérifier si le profil existe
	var exists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM user_profiles WHERE user_id = $1)", userID).Scan(&exists)
	if err != nil {
		return nil, err
	}

	if !exists {
		// Créer un nouveau profil
		insertQuery := `
            INSERT INTO user_profiles (user_id, phone, poste, departement, is_active, created_at, updated_at)
            VALUES ($1, $2, $3, $4, $5, $6, $7)
        `
		now := time.Now()
		isActive := true
		if req.IsActive != nil {
			isActive = *req.IsActive
		}
		_, err = tx.Exec(insertQuery, userID, req.Phone, req.Poste, req.Departement, isActive, now, now)
		if err != nil {
			return nil, err
		}
	} else {
		// Mettre à jour le profil existant
		setParts := []string{}
		args := []interface{}{}
		argIndex := 1

		if req.Phone != nil {
			setParts = append(setParts, fmt.Sprintf("phone = $%d", argIndex))
			args = append(args, *req.Phone)
			argIndex++
		}
		if req.Poste != nil {
			setParts = append(setParts, fmt.Sprintf("poste = $%d", argIndex))
			args = append(args, *req.Poste)
			argIndex++
		}
		if req.Departement != nil {
			setParts = append(setParts, fmt.Sprintf("departement = $%d", argIndex))
			args = append(args, *req.Departement)
			argIndex++
		}
		if req.IsActive != nil {
			setParts = append(setParts, fmt.Sprintf("is_active = $%d", argIndex))
			args = append(args, *req.IsActive)
			argIndex++
		}

		if len(setParts) > 0 {
			setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
			args = append(args, time.Now())
			argIndex++

			args = append(args, userID)

			query := fmt.Sprintf(`
                UPDATE user_profiles 
                SET %s 
                WHERE user_id = $%d
            `, strings.Join(setParts, ", "), argIndex)

			_, err = tx.Exec(query, args...)
			if err != nil {
				return nil, err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return r.GetByID(userID)
}

// Delete supprime un utilisateur et son profil
func (r *UserRepository) Delete(id int) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Supprimer le profil d'abord (clé étrangère)
	_, err = tx.Exec("DELETE FROM user_profiles WHERE user_id = $1", id)
	if err != nil {
		return err
	}

	// Supprimer l'utilisateur
	result, err := tx.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return tx.Commit()
}

// UpdatePassword met à jour le mot de passe d'un utilisateur
func (r *UserRepository) UpdatePassword(id int, passwordHash string) error {
	query := `
        UPDATE users 
        SET password_hash = $1, updated_at = $2 
        WHERE id = $3
    `

	result, err := r.db.Exec(query, passwordHash, time.Now(), id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// UpdateLastLogin met à jour la dernière connexion
func (r *UserRepository) UpdateLastLogin(id int, lastLogin *time.Time) error {
	query := `
        UPDATE users 
        SET last_login = $1, updated_at = $2 
        WHERE id = $3
    `

	_, err := r.db.Exec(query, lastLogin, time.Now(), id)
	return err
}

// UpdateStatus met à jour le statut actif/inactif
func (r *UserRepository) UpdateStatus(id int, isActive bool) error {
	query := `
        UPDATE users 
        SET is_active = $1, updated_at = $2 
        WHERE id = $3
    `

	result, err := r.db.Exec(query, isActive, time.Now(), id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// scanUsers fonction utilitaire pour scanner les résultats SQL
func (r *UserRepository) scanUsers(rows *sql.Rows) ([]*models.User, error) {
	users := []*models.User{}

	for rows.Next() {
		var user models.User
		var profile models.UserProfile
		var profileIsActive sql.NullBool
		var profileCreatedAt, profileUpdatedAt sql.NullTime

		err := rows.Scan(
			&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Email,
			&user.PasswordHash, &user.Role, &user.EntityType, &user.EntityID,
			&user.IsActive, &user.EmailVerified, &user.LastLogin, &user.CreatedAt, &user.UpdatedAt,
			&profile.Phone, &profile.Poste, &profile.Departement,
			&profileIsActive, &profileCreatedAt, &profileUpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Associer le profil s'il y a des données (au moins un champ non null)
		if profile.Phone != nil || profile.Poste != nil || profile.Departement != nil {
			profile.UserID = user.ID
			if profileIsActive.Valid {
				profile.IsActive = profileIsActive.Bool
			}
			if profileCreatedAt.Valid {
				profile.CreatedAt = profileCreatedAt.Time
			}
			if profileUpdatedAt.Valid {
				profile.UpdatedAt = profileUpdatedAt.Time
			}
			user.Profile = &profile
		}

		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
