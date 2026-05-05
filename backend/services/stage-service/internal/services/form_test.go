package services_test

import (
	"context"
	"fmt"
	"shared/models"
	"stage/internal/repositories"
	"stage/internal/services"
	"testing"

	sharedContext "shared/context"
)

func buildFormContext(role string, userID int) context.Context {
	claims := &models.CustomClaims{
		UserID: userID,
		Role:   models.UserRole(role),
	}
	return sharedContext.SetUserClaims(context.Background(), claims)
}

// --- Get ---

func TestFormGet_SansContext_RetourneErreur(t *testing.T) {
	repo := &repositories.FormRepositoryMock{}
	svc := services.NewFormService(repo)

	_, err := svc.Get(context.Background())

	if err == nil {
		t.Error("Devrait retourner une erreur sans claims")
	}
}

func TestFormGet_AvecClaims_AppelleRepository(t *testing.T) {
	called := false
	repo := &repositories.FormRepositoryMock{
		GetFn: func(ctx context.Context, userID int, role models.UserRole) ([]*models.FormFormSection, error) {
			called = true
			if userID != 1 {
				t.Errorf("UserID attendu 1, obtenu %d", userID)
			}
			if role != models.RoleStudent {
				t.Errorf("Role attendu student, obtenu %s", role)
			}
			return []*models.FormFormSection{}, nil
		},
	}
	svc := services.NewFormService(repo)

	_, err := svc.Get(buildFormContext("student", 1))

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
	if !called {
		t.Error("Le repository n'a pas été appelé")
	}
}

func TestFormGet_RepositoryEchoue_RetourneErreur(t *testing.T) {
	repo := &repositories.FormRepositoryMock{
		GetFn: func(ctx context.Context, userID int, role models.UserRole) ([]*models.FormFormSection, error) {
			return nil, fmt.Errorf("db error")
		},
	}
	svc := services.NewFormService(repo)

	_, err := svc.Get(buildFormContext("student", 1))

	if err == nil {
		t.Error("Devrait propager l'erreur du repository")
	}
}

// --- UpdateForm ---

func TestUpdateForm_SansContext_RetourneErreur(t *testing.T) {
	repo := &repositories.FormRepositoryMock{}
	svc := services.NewFormService(repo)

	_, err := svc.UpdateForm(context.Background(), &models.Form{})

	if err == nil {
		t.Error("Devrait retourner une erreur sans claims")
	}
}

func TestUpdateForm_AvecClaims_AppelleRepository(t *testing.T) {
	called := false
	repo := &repositories.FormRepositoryMock{
		UpdateFormFn: func(ctx context.Context, data *models.Form, userID int) (*models.Form, error) {
			called = true
			if userID != 5 {
				t.Errorf("UserID attendu 5, obtenu %d", userID)
			}
			return data, nil
		},
	}
	svc := services.NewFormService(repo)

	_, err := svc.UpdateForm(buildFormContext("teacher", 5), &models.Form{ID: 1})

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
	if !called {
		t.Error("Le repository n'a pas été appelé")
	}
}

func TestUpdateForm_RepositoryEchoue_RetourneErreur(t *testing.T) {
	repo := &repositories.FormRepositoryMock{
		UpdateFormFn: func(ctx context.Context, data *models.Form, userID int) (*models.Form, error) {
			return nil, fmt.Errorf("db error")
		},
	}
	svc := services.NewFormService(repo)

	_, err := svc.UpdateForm(buildFormContext("teacher", 5), &models.Form{})

	if err == nil {
		t.Error("Devrait propager l'erreur du repository")
	}
}

// --- UpdateFormSection ---

func TestUpdateFormSection_SansContext_RetourneErreur(t *testing.T) {
	repo := &repositories.FormRepositoryMock{}
	svc := services.NewFormService(repo)

	_, err := svc.UpdateFormSection(context.Background(), []models.FormSection{{UserID: 1}})

	if err == nil {
		t.Error("Devrait retourner une erreur sans claims")
	}
}

func TestUpdateFormSection_SectionUtilisateurTrouvee_AppelleRepository(t *testing.T) {
	called := false
	repo := &repositories.FormRepositoryMock{
		UpdateFormSectionFn: func(ctx context.Context, data *models.FormSection, userID int) ([]models.FormSection, error) {
			called = true
			if data.UserID != 3 {
				t.Errorf("UserID attendu 3, obtenu %d", data.UserID)
			}
			return []models.FormSection{*data}, nil
		},
	}
	svc := services.NewFormService(repo)

	sections := []models.FormSection{
		{UserID: 3, SectionType: "STUDENT"},
	}
	_, err := svc.UpdateFormSection(buildFormContext("student", 3), sections)

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
	if !called {
		t.Error("Le repository n'a pas été appelé")
	}
}

func TestUpdateFormSection_AucuneSectionUtilisateur_RetourneErreur(t *testing.T) {
	repo := &repositories.FormRepositoryMock{}
	svc := services.NewFormService(repo)

	sections := []models.FormSection{
		{UserID: 99, SectionType: "TEACHER"},
	}
	_, err := svc.UpdateFormSection(buildFormContext("student", 3), sections)

	if err == nil {
		t.Error("Devrait retourner une erreur si aucune section n'appartient à l'utilisateur")
	}
}

// --- CreateForm ---

func TestCreateForm_SansContext_RetourneErreur(t *testing.T) {
	repo := &repositories.FormRepositoryMock{}
	svc := services.NewFormService(repo)

	err := svc.CreateForm(context.Background(), models.FormCreationData{})

	if err == nil {
		t.Error("Devrait retourner une erreur sans claims")
	}
}

func TestCreateForm_SansTeacherNiTutor_CreeSectionEtudiantSeulement(t *testing.T) {
	var capturedSections []models.FormSection
	repo := &repositories.FormRepositoryMock{
		CreateFormFn: func(ctx context.Context, form models.Form, sections []models.FormSection, userID int) error {
			capturedSections = sections
			return nil
		},
	}
	svc := services.NewFormService(repo)

	err := svc.CreateForm(buildFormContext("admin", 1), models.FormCreationData{
		StageID:   1,
		StudentID: 2,
	})

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
	if len(capturedSections) != 1 {
		t.Errorf("Attendu 1 section, obtenu %d", len(capturedSections))
	}
	if capturedSections[0].SectionType != "STUDENT" {
		t.Errorf("Section attendue STUDENT, obtenu %s", capturedSections[0].SectionType)
	}
}

func TestCreateForm_AvecTeacherEtTutor_Cree3Sections(t *testing.T) {
	var capturedSections []models.FormSection
	repo := &repositories.FormRepositoryMock{
		CreateFormFn: func(ctx context.Context, form models.Form, sections []models.FormSection, userID int) error {
			capturedSections = sections
			return nil
		},
	}
	svc := services.NewFormService(repo)

	teacherID := 10
	tutorID := 20
	err := svc.CreateForm(buildFormContext("admin", 1), models.FormCreationData{
		StageID:   1,
		StudentID: 2,
		TeacherID: &teacherID,
		TutorID:   &tutorID,
	})

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
	if len(capturedSections) != 3 {
		t.Errorf("Attendu 3 sections, obtenu %d", len(capturedSections))
	}
}

func TestCreateForm_RepositoryEchoue_RetourneErreur(t *testing.T) {
	repo := &repositories.FormRepositoryMock{
		CreateFormFn: func(ctx context.Context, form models.Form, sections []models.FormSection, userID int) error {
			return fmt.Errorf("db error")
		},
	}
	svc := services.NewFormService(repo)

	err := svc.CreateForm(buildFormContext("admin", 1), models.FormCreationData{StageID: 1, StudentID: 2})

	if err == nil {
		t.Error("Devrait propager l'erreur du repository")
	}
}
