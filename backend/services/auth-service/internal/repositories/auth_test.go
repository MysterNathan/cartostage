package repositories_test

import (
	"database/sql"
	"errors"
	"shared/test"
	"testing"

	"auth/internal/repositories"
	_ "github.com/lib/pq"
)

func TestFindUserByUsername_UserExists(t *testing.T) {
	db := test.SetupTestDB(t)
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
	db := test.SetupTestDB(t)
	repo := repositories.NewAuthRepository(db)

	user, err := repo.FindUserByUsername("utilisateur_inexistant")

	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("Erreur attendue sql.ErrNoRows, obtenu : %v", err)
	}
	if user != nil {
		t.Fatal("L'utilisateur ne devrait pas exister")
	}
}
