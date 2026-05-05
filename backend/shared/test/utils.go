package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"net/http"
	"net/http/httptest"
	"os"
	sharedContext "shared/context"
	"shared/models"
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

func MakeRequest(t *testing.T, method, url string, body any) *http.Request {
	t.Helper()
	if body == nil {
		return httptest.NewRequest(method, url, nil)
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("Impossible de sérialiser le body : %v", err)
	}
	return httptest.NewRequest(method, url, bytes.NewBuffer(jsonBody))
}

func RequestWithVars(r *http.Request, vars map[string]string) *http.Request {
	return mux.SetURLVars(r, vars)
}

func CtxWithClaims(role models.UserRole, userID int) context.Context {
	claims := &models.CustomClaims{
		Role:   role,
		UserID: userID,
	}
	return sharedContext.SetClaimsInContext(context.Background(), claims)
}
