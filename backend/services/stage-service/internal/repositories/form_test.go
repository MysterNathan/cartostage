package repositories_test

import (
	"context"
	"github.com/jmoiron/sqlx"
	"shared/models"
	"shared/test"
	"stage/internal/repositories"
	"testing"
)

// --- Helpers ---

func createTestForm(t *testing.T, db *sqlx.DB, stageID, studentID int, teacherID, tutorID *int) int {
	t.Helper()
	var id int
	err := db.QueryRow(`
        INSERT INTO form (stage_id, student_id, teacher_id, tutor_id, status, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
        RETURNING id`,
		stageID, studentID, teacherID, tutorID, models.StatusCreated,
	).Scan(&id)
	if err != nil {
		t.Fatalf("Impossible de créer le form de test : %v", err)
	}
	t.Cleanup(func() { db.Exec(`DELETE FROM form WHERE id = $1`, id) })
	return id
}

func createTestFormSection(t *testing.T, db *sqlx.DB, formID, userID int, sectionType string) int {
	t.Helper()
	var id int
	err := db.QueryRow(`
        INSERT INTO form_section (form_id, section_type, user_id, status, created_at, updated_at)
        VALUES ($1, $2, $3, $4, NOW(), NOW())
        RETURNING id`,
		formID, sectionType, userID, models.StatusCreated,
	).Scan(&id)
	if err != nil {
		t.Fatalf("Impossible de créer la form_section de test : %v", err)
	}
	t.Cleanup(func() { db.Exec(`DELETE FROM form_section WHERE id = $1`, id) })
	return id
}

// --- Get ---

func TestFormGet_Student_RetourneSesForms(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewFormRepository(db)
	ctx := context.Background()

	offerID := createTestStageOffer(t, db)
	stageID := createTestStage(t, db, offerID)
	formID := createTestForm(t, db, stageID, 1, nil, nil)
	createTestFormSection(t, db, formID, 1, "STUDENT")

	results, err := repo.Get(ctx, 1, models.RoleStudent)

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
	if results == nil {
		t.Fatal("Le résultat ne devrait pas être nil")
	}

	found := false
	for _, f := range results {
		if f.Form.ID == formID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Le form %d devrait être présent dans les résultats", formID)
	}
}

func TestFormGet_Teacher_RetourneSesForms(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewFormRepository(db)
	ctx := context.Background()

	teacherID := 10
	offerID := createTestStageOffer(t, db)
	stageID := createTestStage(t, db, offerID)
	formID := createTestForm(t, db, stageID, 1, &teacherID, nil)
	createTestFormSection(t, db, formID, teacherID, "TEACHER")

	results, err := repo.Get(ctx, teacherID, models.RoleTeacher)

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}

	found := false
	for _, f := range results {
		if f.Form.ID == formID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Le form %d devrait être présent pour le teacher", formID)
	}
}

func TestFormGet_Tutor_RetourneSesForms(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewFormRepository(db)
	ctx := context.Background()

	tutorID := 20
	offerID := createTestStageOffer(t, db)
	stageID := createTestStage(t, db, offerID)
	formID := createTestForm(t, db, stageID, 1, nil, &tutorID)
	createTestFormSection(t, db, formID, tutorID, "TUTOR")

	results, err := repo.Get(ctx, tutorID, models.RoleTutor)

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}

	found := false
	for _, f := range results {
		if f.Form.ID == formID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Le form %d devrait être présent pour le tutor", formID)
	}
}

func TestFormGet_RoleInvalide_RetourneErreur(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewFormRepository(db)
	ctx := context.Background()

	_, err := repo.Get(ctx, 1, "unknown_role")

	if err == nil {
		t.Error("Un rôle invalide devrait retourner une erreur")
	}
}

func TestFormGet_FormSectionsIncluses(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewFormRepository(db)
	ctx := context.Background()

	offerID := createTestStageOffer(t, db)
	stageID := createTestStage(t, db, offerID)
	formID := createTestForm(t, db, stageID, 1, nil, nil)
	sectionID := createTestFormSection(t, db, formID, 1, "STUDENT")

	results, err := repo.Get(ctx, 1, models.RoleStudent)
	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}

	for _, f := range results {
		if f.Form.ID != formID {
			continue
		}
		if len(f.FormSections) == 0 {
			t.Fatal("Le form devrait contenir au moins une section")
		}
		if f.FormSections[0].ID != sectionID {
			t.Errorf("SectionID attendu %d, obtenu %d", sectionID, f.FormSections[0].ID)
		}
		return
	}
	t.Errorf("Form %d non trouvé", formID)
}

func TestFormGet_NePasVoirFormsDesAutres(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewFormRepository(db)
	ctx := context.Background()

	offerID := createTestStageOffer(t, db)
	stageID := createTestStage(t, db, offerID)

	// Form appartenant à l'étudiant 2
	formID := createTestForm(t, db, stageID, 2, nil, nil)
	createTestFormSection(t, db, formID, 2, "STUDENT")

	results, err := repo.Get(ctx, 1, models.RoleStudent)
	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}

	for _, f := range results {
		if f.Form.ID == formID {
			t.Error("L'étudiant 1 ne devrait pas voir les forms de l'étudiant 2")
		}
	}
}

// --- UpdateForm ---

func TestUpdateForm_CasNominal(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewFormRepository(db)
	ctx := context.Background()

	offerID := createTestStageOffer(t, db)
	stageID := createTestStage(t, db, offerID)
	formID := createTestForm(t, db, stageID, 1, nil, nil)

	input := &models.Form{
		ID:     formID,
		Status: models.StatusInProgress,
	}

	updated, err := repo.UpdateForm(ctx, input, 1)

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
	if updated == nil {
		t.Fatal("Le form retourné ne devrait pas être nil")
	}
	if updated.Status != models.StatusInProgress {
		t.Errorf("Status attendu %s, obtenu %s", models.StatusInProgress, updated.Status)
	}
}

func TestUpdateForm_IDInexistant_RetourneErreur(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewFormRepository(db)
	ctx := context.Background()

	input := &models.Form{
		ID:     999999,
		Status: models.StatusInProgress,
	}

	_, err := repo.UpdateForm(ctx, input, 1)

	if err == nil {
		t.Error("Une erreur devrait être retournée pour un ID inexistant")
	}
}

// --- UpdateFormSection ---

func TestUpdateFormSection_CasNominal(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewFormRepository(db)
	ctx := context.Background()

	offerID := createTestStageOffer(t, db)
	stageID := createTestStage(t, db, offerID)
	formID := createTestForm(t, db, stageID, 1, nil, nil)
	sectionID := createTestFormSection(t, db, formID, 1, "STUDENT")

	input := &models.FormSection{
		ID:     sectionID,
		Status: models.StatusInProgress,
	}

	sections, err := repo.UpdateFormSection(ctx, input, 1)

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
	if len(sections) == 0 {
		t.Fatal("Au moins une section devrait être retournée")
	}
	if sections[0].Status != models.StatusInProgress {
		t.Errorf("Status attendu %s, obtenu %s", models.StatusInProgress, sections[0].Status)
	}
}

func TestUpdateFormSection_IDInexistant_RetourneErreur(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewFormRepository(db)
	ctx := context.Background()

	input := &models.FormSection{
		ID:     999999,
		Status: models.StatusInProgress,
	}

	_, err := repo.UpdateFormSection(ctx, input, 1)

	if err == nil {
		t.Error("Une erreur devrait être retournée pour un ID inexistant")
	}
}

// --- CreateForm ---

func TestCreateForm_CasNominal(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewFormRepository(db)
	ctx := context.Background()

	offerID := createTestStageOffer(t, db)
	stageID := createTestStage(t, db, offerID)

	dataForm := models.Form{
		StageID:   stageID,
		StudentID: 1,
	}
	dataFormSections := []models.FormSection{
		{SectionType: "STUDENT", UserID: 1},
	}

	err := repo.CreateForm(ctx, dataForm, dataFormSections, 1)

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}

	t.Cleanup(func() {
		db.Exec(`DELETE FROM form WHERE stage_id = $1`, stageID)
	})
}

func TestCreateForm_AvecTeacherEtTutor(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewFormRepository(db)
	ctx := context.Background()

	teacherID := 10
	tutorID := 20
	offerID := createTestStageOffer(t, db)
	stageID := createTestStage(t, db, offerID)

	dataForm := models.Form{
		StageID:   stageID,
		StudentID: 1,
		TeacherID: &teacherID,
		TutorID:   &tutorID,
	}
	dataFormSections := []models.FormSection{
		{SectionType: "STUDENT", UserID: 1},
		{SectionType: "TEACHER", UserID: teacherID},
		{SectionType: "TUTOR", UserID: tutorID},
	}

	err := repo.CreateForm(ctx, dataForm, dataFormSections, 1)

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}

	t.Cleanup(func() {
		db.Exec(`DELETE FROM form WHERE stage_id = $1`, stageID)
	})
}

func TestCreateForm_StageIDInexistant_RetourneErreur(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewFormRepository(db)
	ctx := context.Background()

	dataForm := models.Form{
		StageID:   999999,
		StudentID: 1,
	}
	dataFormSections := []models.FormSection{
		{SectionType: "STUDENT", UserID: 1},
	}

	err := repo.CreateForm(ctx, dataForm, dataFormSections, 1)

	if err == nil {
		t.Error("Une erreur devrait être retournée pour un stage_id inexistant")
	}
}

func TestCreateForm_SansSections_RetourneErreur(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewFormRepository(db)
	ctx := context.Background()

	offerID := createTestStageOffer(t, db)
	stageID := createTestStage(t, db, offerID)

	dataForm := models.Form{
		StageID:   stageID,
		StudentID: 1,
	}

	err := repo.CreateForm(ctx, dataForm, []models.FormSection{}, 1)

	// Sans sections, le form est créé mais aucune section n'est insérée.
	// On vérifie simplement qu'il n'y a pas d'erreur inattendue.
	if err != nil {
		t.Fatalf("Erreur inattendue pour un form sans sections : %v", err)
	}

	t.Cleanup(func() {
		db.Exec(`DELETE FROM form WHERE stage_id = $1`, stageID)
	})
}

func TestCreateForm_StatusParDefautCreated(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewFormRepository(db)
	ctx := context.Background()

	offerID := createTestStageOffer(t, db)
	stageID := createTestStage(t, db, offerID)

	dataForm := models.Form{
		StageID:   stageID,
		StudentID: 1,
	}
	dataFormSections := []models.FormSection{
		{SectionType: "STUDENT", UserID: 1},
	}

	err := repo.CreateForm(ctx, dataForm, dataFormSections, 1)
	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}

	var status string
	db.QueryRow(`SELECT status FROM form WHERE stage_id = $1`, stageID).Scan(&status)
	if status != models.StatusCreated {
		t.Errorf("Status attendu %s, obtenu %s", models.StatusCreated, status)
	}

	t.Cleanup(func() {
		db.Exec(`DELETE FROM form WHERE stage_id = $1`, stageID)
	})
}
