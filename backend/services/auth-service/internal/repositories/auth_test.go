package repositories_test

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"

	"auth/internal/repositories"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	t.Helper()

	host := getEnvOrDefault("TEST_DB_HOST", "192.168.1.11")
	port := getEnvOrDefault("TEST_DB_PORT", "5432")
	user := getEnvOrDefault("TEST_DB_USER", "stages_user")
	password := getEnvOrDefault("TEST_DB_PASSWORD", "stages_password")
	dbname := getEnvOrDefault("TEST_DB_NAME", "postgres")
	sslmode := getEnvOrDefault("TEST_DB_SSLMODE", "disable")

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbname, sslmode,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		t.Skipf("Base de données non disponible, test ignoré: %v", err)
	}

	t.Cleanup(func() { db.Close() })
	return db
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func TestFindUserByUsername_UserExists(t *testing.T) {
	db := setupTestDB(t)
	repo := repositories.NewAuthRepository(db)

	// admin_user doit exister dans create_users_table.sql
	user, err := repo.FindUserByUsername("admin_user")

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
	if user == nil {
		t.Fatal("L'utilisateur devrait exister")
	}
	if user.Username != "admin_user" {
		t.Errorf("Username attendu 'admin_user', obtenu '%s'", user.Username)
	}
}

func TestFindUserByUsername_UserNotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := repositories.NewAuthRepository(db)

	user, err := repo.FindUserByUsername("utilisateur_inexistant")

	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("Erreur attendue sql.ErrNoRows, obtenu : %v", err)
	}
	if user != nil {
		t.Fatal("L'utilisateur ne devrait pas exister")
	}
}
