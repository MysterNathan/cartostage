package repositories_test

import (
	"context"
	"shared/models"
	"shared/test"
	"student/internal/repositories"
	"testing"
)

func TestGetAll_Student_AdminVoitTousLesEleves(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewUserRepository(db)

	filter := models.UserFilter{
		RequestorRole: models.RoleAdmin,
		RequestorID:   1,
	}

	users, err := repo.GetAll(context.Background(), filter)

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
	for _, u := range users {
		if u.Role != string(models.RoleStudent) {
			t.Errorf("Attendu role 'eleve', obtenu '%s'", u.Role)
		}
	}
}

func TestGetAll_Student_TeacherVoitSesEleves(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewUserRepository(db)

	filter := models.UserFilter{
		RequestorRole: models.RoleTeacher,
		RequestorID:   1, // teacher seedé avec des stages liés à des élèves
	}

	users, err := repo.GetAll(context.Background(), filter)

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
	if len(users) == 0 {
		t.Fatal("Le teacher devrait voir au moins un élève")
	}
	for _, u := range users {
		if u.Role != string(models.RoleStudent) {
			t.Errorf("Attendu role 'eleve', obtenu '%s'", u.Role)
		}
	}
}

func TestGetAll_Student_StudentNeVoitQueLuiMeme(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewUserRepository(db)

	filter := models.UserFilter{
		RequestorRole: models.RoleStudent,
		RequestorID:   1, // élève seedé
	}

	users, err := repo.GetAll(context.Background(), filter)

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("L'élève devrait voir uniquement lui-même, obtenu %d résultats", len(users))
	}
	if users[0].ID != 1 {
		t.Errorf("Attendu ID 1, obtenu %d", users[0].ID)
	}
}

func TestGetAll_Student_TutorVoitSesEleves(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewUserRepository(db)

	filter := models.UserFilter{
		RequestorRole: models.RoleTutor,
		RequestorID:   1, // tutor seedé avec des stages liés à des élèves
	}

	users, err := repo.GetAll(context.Background(), filter)

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
	for _, u := range users {
		if u.Role != string(models.RoleStudent) {
			t.Errorf("Attendu role 'eleve', obtenu '%s'", u.Role)
		}
	}
}

func TestGetByID_Student_RetourneEleve(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewUserRepository(db)

	user, err := repo.GetByID(context.Background(), 1)

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
	if user == nil {
		t.Fatal("L'utilisateur devrait exister")
	}
	if user.Role != string(models.RoleStudent) {
		t.Errorf("Attendu role 'eleve', obtenu '%s'", user.Role)
	}
}

func TestGetByID_Student_IDInexistant_RetourneErreur(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewUserRepository(db)

	user, err := repo.GetByID(context.Background(), 9999)

	if err == nil {
		t.Fatal("Une erreur était attendue pour un ID inexistant")
	}
	if user != nil {
		t.Fatal("Aucun utilisateur ne devrait être retourné")
	}
}

func TestCreate_Student_CreationOK(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewUserRepository(db)
	phone := "0600000000"

	user := &models.User{
		Username:     "nouveau_eleve",
		FirstName:    "Jean",
		LastName:     "Dupont",
		Email:        "jean.dupont@test.com",
		PasswordHash: "hashed_password",
		Phone:        &phone,
		IsActive:     true,
	}

	err := repo.Create(context.Background(), user)

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
	if user.ID == 0 {
		t.Error("L'ID devrait être renseigné après création")
	}
	if user.Role != string(models.RoleStudent) {
		t.Errorf("Le role devrait être forcé à 'eleve', obtenu '%s'", user.Role)
	}
}

func TestCreate_Student_EmailDuplique_RetourneErreur(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewUserRepository(db)

	user := &models.User{
		Username:     "doublon",
		FirstName:    "Jean",
		LastName:     "Dupont",
		Email:        "student_existing@test.com", // email seedé
		PasswordHash: "hashed_password",
		IsActive:     true,
	}

	err := repo.Create(context.Background(), user)

	if err == nil {
		t.Fatal("Une erreur était attendue pour un email dupliqué")
	}
}

func TestUpdate_Student_MiseAJourOK(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewUserRepository(db)
	phone := "0611111111"
	user := &models.User{
		ID:        1,
		Username:  "eleve_modifie",
		FirstName: "Jean",
		LastName:  "Modifié",
		Email:     "jean.modifie@test.com",
		Phone:     &phone,
		IsActive:  true,
	}

	err := repo.Update(context.Background(), user)

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
}

func TestUpdate_Student_IDInexistant_RetourneErreur(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewUserRepository(db)

	user := &models.User{
		ID:       9999,
		Username: "inexistant",
		Email:    "inexistant@test.com",
	}

	err := repo.Update(context.Background(), user)

	if err == nil {
		t.Fatal("Une erreur était attendue pour un ID inexistant")
	}
}

func TestDelete_Student_SuppressionOK(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewUserRepository(db)

	err := repo.Delete(context.Background(), 1)

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
}

func TestDelete_Student_IDInexistant_RetourneErreur(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewUserRepository(db)

	err := repo.Delete(context.Background(), 9999)

	if err == nil {
		t.Fatal("Une erreur était attendue pour un ID inexistant")
	}
}
