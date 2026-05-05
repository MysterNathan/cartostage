package test

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
	"testing"
)

func SetupTestDB(t *testing.T) *sqlx.DB {
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
