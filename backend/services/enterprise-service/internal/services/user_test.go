package services_test

import (
	"context"
	"enterprise/internal/repositories"
	"enterprise/internal/services"
	"fmt"
	"shared/models"
	"shared/test"
	"testing"
)

func buildStudentService(t *testing.T, repo *repositories.MockUserRepository) services.UserService {
	t.Helper()
	return services.NewUserService(repo)
}

// --- GetAll ---

func TestGetAll_SansClaims_RetourneErreur(t *testing.T) {
	repo := repositories.NewMockUserRepository()
	svc := buildStudentService(t, repo)

	_, err := svc.GetAll(context.Background())

	if err == nil {
		t.Fatal("Une erreur était attendue sans claims dans le contexte")
	}
}

func TestGetAll_AvecClaims_RetourneUtilisateurs(t *testing.T) {
	repo := repositories.NewMockUserRepository()
	repo.Users[1] = &models.User{ID: 1, Username: "alice", Role: "eleve"}
	repo.Users[2] = &models.User{ID: 2, Username: "bob", Role: "eleve"}
	svc := buildStudentService(t, repo)

	ctx := test.CtxWithClaims(models.RoleAdmin, 1)
	users, err := svc.GetAll(ctx)

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
	if len(users) != 2 {
		t.Errorf("Attendu 2 utilisateurs, obtenu %d", len(users))
	}
}

func TestGetAll_ErreurRepository_PropageErreur(t *testing.T) {
	repo := repositories.NewMockUserRepository()
	repo.ErrorToReturn = fmt.Errorf("db error")
	svc := buildStudentService(t, repo)

	ctx := test.CtxWithClaims(models.RoleAdmin, 1)
	_, err := svc.GetAll(ctx)

	if err == nil {
		t.Fatal("Une erreur était attendue depuis le repository")
	}
}

// --- GetByID ---

func TestGetByID_UserExistant_RetourneUser(t *testing.T) {
	repo := repositories.NewMockUserRepository()
	repo.Users[1] = &models.User{ID: 1, Username: "alice", Role: "eleve"}
	svc := buildStudentService(t, repo)

	user, err := svc.GetByID(context.Background(), 1)

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
	if user == nil {
		t.Fatal("L'utilisateur ne devrait pas être nil")
	}
	if user.Username != "alice" {
		t.Errorf("Username attendu 'alice', obtenu '%s'", user.Username)
	}
}

func TestGetByID_UserInexistant_RetourneErreur(t *testing.T) {
	repo := repositories.NewMockUserRepository()
	svc := buildStudentService(t, repo)

	_, err := svc.GetByID(context.Background(), 9999)

	if err == nil {
		t.Fatal("Une erreur était attendue pour un ID inexistant")
	}
}

// --- Create ---

func TestCreate_RequeteValide_RetourneUserPublic(t *testing.T) {
	repo := repositories.NewMockUserRepository()
	svc := buildStudentService(t, repo)

	req := models.CreateUserRequest{
		Username:  "charlie",
		FirstName: "Charlie",
		LastName:  "Dupont",
		Email:     "charlie@test.com",
		Password:  "password123",
		Role:      models.RoleStudent,
	}

	user, err := svc.Create(context.Background(), req)

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
	if user == nil {
		t.Fatal("L'utilisateur créé ne devrait pas être nil")
	}
	if user.Username != "charlie" {
		t.Errorf("Username attendu 'charlie', obtenu '%s'", user.Username)
	}
}

func TestCreate_ErreurRepository_PropageErreur(t *testing.T) {
	repo := repositories.NewMockUserRepository()
	repo.ErrorToReturn = fmt.Errorf("contrainte unique violée")
	svc := buildStudentService(t, repo)

	req := models.CreateUserRequest{
		Username: "charlie",
		Email:    "charlie@test.com",
		Password: "password123",
		Role:     models.RoleStudent,
	}

	_, err := svc.Create(context.Background(), req)

	if err == nil {
		t.Fatal("Une erreur était attendue depuis le repository")
	}
}

// --- Update ---

func TestUpdate_UserExistant_MiseAJourOK(t *testing.T) {
	repo := repositories.NewMockUserRepository()
	repo.Users[1] = &models.User{ID: 1, Username: "alice", Role: "eleve"}
	svc := buildStudentService(t, repo)

	req := models.UpdateUserRequest{
		FirstName: strPtr("Alice"),
		LastName:  strPtr("Modifiée"),
	}

	user, err := svc.Update(context.Background(), 1, req)

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
	if user == nil {
		t.Fatal("L'utilisateur mis à jour ne devrait pas être nil")
	}
}

func TestUpdate_UserInexistant_RetourneErreur(t *testing.T) {
	repo := repositories.NewMockUserRepository()
	svc := buildStudentService(t, repo)

	req := models.UpdateUserRequest{FirstName: strPtr("Alice")}

	_, err := svc.Update(context.Background(), 9999, req)

	if err == nil {
		t.Fatal("Une erreur était attendue pour un ID inexistant")
	}
}

// --- Delete ---

func TestDelete_UserExistant_SuppressionOK(t *testing.T) {
	repo := repositories.NewMockUserRepository()
	repo.Users[1] = &models.User{ID: 1, Username: "alice"}
	svc := buildStudentService(t, repo)

	err := svc.Delete(context.Background(), 1)

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
	if _, ok := repo.Users[1]; ok {
		t.Error("L'utilisateur devrait avoir été supprimé du repository")
	}
}

func TestDelete_UserInexistant_RetourneErreur(t *testing.T) {
	repo := repositories.NewMockUserRepository()
	svc := buildStudentService(t, repo)

	err := svc.Delete(context.Background(), 9999)

	if err == nil {
		t.Fatal("Une erreur était attendue pour un ID inexistant")
	}
}

func strPtr(s string) *string {
	return &s
}
