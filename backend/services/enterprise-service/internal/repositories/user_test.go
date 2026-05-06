package repositories_test

import (
	"context"
	"shared/models"
	"shared/test"
	"testing"

	"enterprise/internal/repositories"
)

func TestGetAll_Tutor_AdminVoitTousLesTutors(t *testing.T) {
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
		if u.Role != string(models.RoleTutor) {
			t.Errorf("Attendu role 'tutor', obtenu '%s'", u.Role)
		}
	}
}

func TestGetAll_Tutor_TeacherVoitTutorsDeSesEleves(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewUserRepository(db)

	filter := models.UserFilter{
		RequestorRole: models.RoleTeacher,
		RequestorID:   1,
	}

	users, err := repo.GetAll(context.Background(), filter)

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
	for _, u := range users {
		if u.Role != string(models.RoleTutor) {
			t.Errorf("Attendu role 'tutor', obtenu '%s'", u.Role)
		}
	}
}

func TestGetAll_Tutor_StudentVoitSonTutor(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewUserRepository(db)

	filter := models.UserFilter{
		RequestorRole: models.RoleStudent,
		RequestorID:   1,
	}

	users, err := repo.GetAll(context.Background(), filter)

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
	if len(users) == 0 {
		t.Fatal("L'élève devrait voir au moins un tuteur")
	}
	for _, u := range users {
		if u.Role != string(models.RoleTutor) {
			t.Errorf("Attendu role 'tutor', obtenu '%s'", u.Role)
		}
	}
}

func TestGetAll_Tutor_TutorVoitTutorsDeSonEntreprise(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewUserRepository(db)

	filter := models.UserFilter{
		RequestorRole: models.RoleTutor,
		RequestorID:   1, // tutor seedé avec establishment_id
	}

	users, err := repo.GetAll(context.Background(), filter)

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
	for _, u := range users {
		if u.Role != string(models.RoleTutor) {
			t.Errorf("Attendu role 'tutor', obtenu '%s'", u.Role)
		}
	}
}

func TestGetByID_Tutor_RetourneTutor(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewUserRepository(db)

	user, err := repo.GetByID(context.Background(), 1)

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
	if user == nil {
		t.Fatal("L'utilisateur devrait exister")
	}
}

func TestGetByID_Tutor_IDInexistant_RetourneErreur(t *testing.T) {
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

func TestCreate_Tutor_CreationOK(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewUserRepository(db)
	phone := "0600000002"
	user := &models.User{
		Username:     "nouveau_tutor",
		FirstName:    "Pierre",
		LastName:     "Martin",
		Email:        "pierre.martin@test.com",
		PasswordHash: "hashed_password",
		Role:         string(models.RoleTutor),
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
}

func TestCreate_Tutor_EmailDuplique_RetourneErreur(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewUserRepository(db)

	user := &models.User{
		Username:     "doublon_tutor",
		FirstName:    "Pierre",
		LastName:     "Martin",
		Email:        "tutor_existing@test.com", // email seedé
		PasswordHash: "hashed_password",
		Role:         string(models.RoleTutor),
		IsActive:     true,
	}

	err := repo.Create(context.Background(), user)

	if err == nil {
		t.Fatal("Une erreur était attendue pour un email dupliqué")
	}
}

func TestUpdate_Tutor_MiseAJourOK(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewUserRepository(db)
	phone := "0633333333"
	user := &models.User{
		ID:        1,
		Username:  "tutor_modifie",
		FirstName: "Pierre",
		LastName:  "Modifié",
		Email:     "pierre.modifie@test.com",
		Role:      string(models.RoleTutor),
		Phone:     &phone,
		IsActive:  true,
	}

	err := repo.Update(context.Background(), user)

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
}

func TestUpdate_Tutor_IDInexistant_RetourneErreur(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewUserRepository(db)

	user := &models.User{
		ID:       9999,
		Username: "inexistant",
		Email:    "inexistant@test.com",
		Role:     string(models.RoleTutor),
	}

	err := repo.Update(context.Background(), user)

	if err == nil {
		t.Fatal("Une erreur était attendue pour un ID inexistant")
	}
}

func TestDelete_Tutor_SuppressionOK(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewUserRepository(db)

	err := repo.Delete(context.Background(), 1)

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
}

func TestDelete_Tutor_IDInexistant_RetourneErreur(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewUserRepository(db)

	err := repo.Delete(context.Background(), 9999)

	if err == nil {
		t.Fatal("Une erreur était attendue pour un ID inexistant")
	}
}
