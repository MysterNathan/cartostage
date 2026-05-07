package repositories_test

import (
	"context"
	"shared/models"
	"shared/test"
	"testing"

	"teacher/internal/repositories"
)

func TestGetAll_Teacher_AdminVoitTousLesTeachers(t *testing.T) {
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
		if u.Role != string(models.RoleTeacher) {
			t.Errorf("Attendu role 'teacher-service', obtenu '%s'", u.Role)
		}
	}
}

func TestGetAll_Teacher_TeacherNeVoitQueLuiMeme(t *testing.T) {
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
	if len(users) != 1 {
		t.Fatalf("Le teacher devrait voir uniquement lui-même, obtenu %d résultats", len(users))
	}
	if users[0].ID != 1 {
		t.Errorf("Attendu ID 1, obtenu %d", users[0].ID)
	}
}

func TestGetAll_Teacher_StudentVoitSonTeacher(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewUserRepository(db)

	filter := models.UserFilter{
		RequestorRole: models.RoleStudent,
		RequestorID:   1, // élève seedé avec un stage lié à un teacher
	}

	users, err := repo.GetAll(context.Background(), filter)

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
	if len(users) == 0 {
		t.Fatal("L'élève devrait voir au moins un teacher")
	}
	for _, u := range users {
		if u.Role != string(models.RoleTeacher) {
			t.Errorf("Attendu role 'teacher-service', obtenu '%s'", u.Role)
		}
	}
}

func TestGetAll_Teacher_TutorVoitTeachersDesSesEleves(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewUserRepository(db)

	filter := models.UserFilter{
		RequestorRole: models.RoleTutor,
		RequestorID:   1,
	}

	users, err := repo.GetAll(context.Background(), filter)

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
	for _, u := range users {
		if u.Role != string(models.RoleTeacher) {
			t.Errorf("Attendu role 'teacher-service', obtenu '%s'", u.Role)
		}
	}
}

func TestGetByID_Teacher_RetourneTeacher(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewUserRepository(db)

	user, err := repo.GetByID(context.Background(), 1)

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
	if user == nil {
		t.Fatal("L'utilisateur devrait exister")
	}
	if user.Role != string(models.RoleTeacher) {
		t.Errorf("Attendu role 'teacher-service', obtenu '%s'", user.Role)
	}
}

func TestGetByID_Teacher_IDInexistant_RetourneErreur(t *testing.T) {
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

func TestCreate_Teacher_CreationOK(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewUserRepository(db)
	phone := "0600000001"
	user := &models.User{
		Username:     "nouveau_teacher",
		FirstName:    "Marie",
		LastName:     "Curie",
		Email:        "marie.curie@test.com",
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
	if user.Role != string(models.RoleTeacher) {
		t.Errorf("Le role devrait être forcé à 'teacher-service', obtenu '%s'", user.Role)
	}
}

func TestCreate_Teacher_EmailDuplique_RetourneErreur(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewUserRepository(db)

	user := &models.User{
		Username:     "doublon_teacher",
		FirstName:    "Marie",
		LastName:     "Curie",
		Email:        "teacher_existing@test.com", // email seedé
		PasswordHash: "hashed_password",
		IsActive:     true,
	}

	err := repo.Create(context.Background(), user)

	if err == nil {
		t.Fatal("Une erreur était attendue pour un email dupliqué")
	}
}

func TestUpdate_Teacher_MiseAJourOK(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewUserRepository(db)
	phone := "0622222222"
	user := &models.User{
		ID:        1,
		Username:  "teacher_modifie",
		FirstName: "Marie",
		LastName:  "Modifiée",
		Email:     "marie.modifiee@test.com",
		Phone:     &phone,
		IsActive:  true,
	}

	err := repo.Update(context.Background(), user)

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
}

func TestUpdate_Teacher_IDInexistant_RetourneErreur(t *testing.T) {
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

func TestDelete_Teacher_SuppressionOK(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewUserRepository(db)

	err := repo.Delete(context.Background(), 1)

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
}

func TestDelete_Teacher_IDInexistant_RetourneErreur(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewUserRepository(db)

	err := repo.Delete(context.Background(), 9999)

	if err == nil {
		t.Fatal("Une erreur était attendue pour un ID inexistant")
	}
}
