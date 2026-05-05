package repositories_test

import (
	"context"
	"github.com/jmoiron/sqlx"
	"shared/models"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"shared/test"
	"stage/internal/repositories"
)

// --- Helpers ---

func createTestStageOffer(t *testing.T, db *sqlx.DB) int {
	t.Helper()
	var id int
	err := db.QueryRow(`
        INSERT INTO stage_offers (position, address, lat, lng, enterprise, sector, capacity_total, capacity_filled, period, course, job_family, scolar_level)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
        RETURNING id`,
		"Développeur Go", "1 rue de la Paix, Paris", 48.8698, 2.3311,
		"Acme Corp", "Informatique", 2, 0, "2024-2025", "BTS SIO", "Développement", "BTS",
	).Scan(&id)
	if err != nil {
		t.Fatalf("Impossible de créer le stage_offer de test : %v", err)
	}
	t.Cleanup(func() {
		db.Exec(`DELETE FROM stage_offers WHERE id = $1`, id)
	})
	return id
}

func createTestStage(t *testing.T, db *sqlx.DB, stageOfferID int) int {
	t.Helper()
	var id int
	err := db.QueryRow(`
        INSERT INTO stages (stage_offer_id, student_id, status, start_date, end_date, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id`,
		stageOfferID, 1, "pending",
		time.Now(), time.Now().AddDate(0, 3, 0),
		time.Now(), time.Now(),
	).Scan(&id)
	if err != nil {
		t.Fatalf("Impossible de créer le stage de test : %v", err)
	}
	t.Cleanup(func() {
		db.Exec(`DELETE FROM stages WHERE id = $1`, id)
	})
	return id
}

// --- GetStagesPublic ---

func TestGetStagesPublic_RetourneListeSansErreur(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewStageRepository(db)

	stages, err := repo.GetStagesPublic()

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
	if stages == nil {
		t.Fatal("La liste ne devrait pas être nil, même vide")
	}
}

func TestGetStagesPublic_ContientLeStageCreé(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewStageRepository(db)

	offerID := createTestStageOffer(t, db)
	stageID := createTestStage(t, db, offerID)

	stages, err := repo.GetStagesPublic()
	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}

	found := false
	for _, s := range stages {
		if s.ID == stageID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Le stage avec l'id %d devrait être présent dans la liste", stageID)
	}
}

func TestGetStagesPublic_StageOfferJointe(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewStageRepository(db)

	offerID := createTestStageOffer(t, db)
	stageID := createTestStage(t, db, offerID)

	stages, err := repo.GetStagesPublic()
	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}

	for _, s := range stages {
		if s.ID == stageID {
			if s.StageOffer == nil {
				t.Error("Le StageOffer jointé ne devrait pas être nil")
			} else if s.StageOffer.ID != offerID {
				t.Errorf("StageOffer.ID attendu %d, obtenu %d", offerID, s.StageOffer.ID)
			}
			return
		}
	}
	t.Errorf("Stage %d non trouvé dans les résultats", stageID)
}

// --- GetStageByID ---

func TestGetStageByID_StageExistant(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewStageRepository(db)
	ctx := context.Background()

	offerID := createTestStageOffer(t, db)
	stageID := createTestStage(t, db, offerID)

	stage, err := repo.GetStageByID(ctx, stageID)

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
	if stage == nil {
		t.Fatal("Le stage devrait exister")
	}
	if stage.ID != stageID {
		t.Errorf("ID attendu %d, obtenu %d", stageID, stage.ID)
	}
	if stage.StageOfferID != offerID {
		t.Errorf("StageOfferID attendu %d, obtenu %d", offerID, stage.StageOfferID)
	}
}

func TestGetStageByID_StageInexistant_RetourneNil(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewStageRepository(db)
	ctx := context.Background()

	stage, err := repo.GetStageByID(ctx, 999999)

	if err != nil {
		t.Fatalf("Erreur inattendue pour un stage inexistant : %v", err)
	}
	if stage != nil {
		t.Error("Le stage devrait être nil pour un ID inexistant")
	}
}

// --- CreateStage ---

func TestCreateStage_CasNominal(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewStageRepository(db)
	ctx := context.Background()

	offerID := createTestStageOffer(t, db)

	now := time.Now().UTC().Truncate(time.Second)
	input := &models.Stage{
		StageOfferID: offerID,
		StudentID:    1,
		Status:       "pending",
		StartDate:    now,
		EndDate:      now.AddDate(0, 3, 0),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	created, err := repo.CreateStage(ctx, input)

	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
	if created == nil {
		t.Fatal("Le stage créé ne devrait pas être nil")
	}
	if created.ID == 0 {
		t.Error("L'ID du stage créé devrait être renseigné")
	}
	if created.StageOfferID != offerID {
		t.Errorf("StageOfferID attendu %d, obtenu %d", offerID, created.StageOfferID)
	}

	t.Cleanup(func() {
		db.Exec(`DELETE FROM stages WHERE id = $1`, created.ID)
	})
}

func TestCreateStage_ChampObligatoireManquant_RetourneErreur(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewStageRepository(db)
	ctx := context.Background()

	// stage_offer_id inexistant → violation de contrainte FK
	input := &models.Stage{
		StageOfferID: 999999,
		StudentID:    1,
		Status:       "pending",
		StartDate:    time.Now(),
		EndDate:      time.Now().AddDate(0, 3, 0),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err := repo.CreateStage(ctx, input)

	if err == nil {
		t.Error("Une erreur devrait être retournée pour une FK invalide")
	}
}

// --- DeleteStage ---

func TestDeleteStage_StageExistant(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewStageRepository(db)
	ctx := context.Background()

	offerID := createTestStageOffer(t, db)
	stageID := createTestStage(t, db, offerID)

	err := repo.DeleteStage(ctx, stageID)

	if err != nil {
		t.Fatalf("Erreur inattendue lors de la suppression : %v", err)
	}

	// Vérification que le stage n'existe plus
	stage, err := repo.GetStageByID(ctx, stageID)
	if err != nil {
		t.Fatalf("Erreur inattendue lors de la vérification : %v", err)
	}
	if stage != nil {
		t.Error("Le stage devrait avoir été supprimé")
	}
}

func TestDeleteStage_StageInexistant_RetourneErreur(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewStageRepository(db)
	ctx := context.Background()

	err := repo.DeleteStage(ctx, 999999)

	if err == nil {
		t.Error("Une erreur devrait être retournée pour un stage inexistant")
	}
}

// --- GetStagesPublic : cas NULL ---

func TestGetStagesPublic_SansStageOffer_StageOfferNil(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewStageRepository(db)

	// Stage sans stage_offer_id (NULL)
	var stageID int
	err := db.QueryRow(`
        INSERT INTO stages (student_id, status, start_date, end_date, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id`,
		1, "pending", time.Now(), time.Now().AddDate(0, 3, 0), time.Now(), time.Now(),
	).Scan(&stageID)
	if err != nil {
		t.Fatalf("Impossible de créer le stage sans offer : %v", err)
	}
	t.Cleanup(func() { db.Exec(`DELETE FROM stages WHERE id = $1`, stageID) })

	stages, err := repo.GetStagesPublic()
	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}

	for _, s := range stages {
		if s.ID == stageID {
			if s.StageOffer != nil {
				t.Error("StageOffer devrait être nil pour un stage sans offer")
			}
			return
		}
	}
	t.Errorf("Stage %d non trouvé", stageID)
}

func TestGetStagesPublic_ChampsStageOfferCorrects(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewStageRepository(db)

	offerID := createTestStageOffer(t, db)
	stageID := createTestStage(t, db, offerID)

	stages, err := repo.GetStagesPublic()
	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}

	for _, s := range stages {
		if s.ID != stageID {
			continue
		}
		if s.StageOffer == nil {
			t.Fatal("StageOffer ne devrait pas être nil")
		}
		if s.StageOffer.Position != "Développeur Go" {
			t.Errorf("Position attendue 'Développeur Go', obtenu '%s'", s.StageOffer.Position)
		}
		if s.StageOffer.Enterprise != "Acme Corp" {
			t.Errorf("Enterprise attendu 'Acme Corp', obtenu '%s'", s.StageOffer.Enterprise)
		}
		if s.StageOffer.Lat != 48.8698 {
			t.Errorf("Lat attendu 48.8698, obtenu %f", s.StageOffer.Lat)
		}
		if s.StageOffer.CapacityTotal != 2 {
			t.Errorf("CapacityTotal attendu 2, obtenu %d", s.StageOffer.CapacityTotal)
		}
		return
	}
	t.Errorf("Stage %d non trouvé", stageID)
}

// --- GetStagesPublic : timestamps ---

func TestGetStagesPublic_TimestampsRemplis(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewStageRepository(db)

	offerID := createTestStageOffer(t, db)
	stageID := createTestStage(t, db, offerID)

	stages, err := repo.GetStagesPublic()
	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}

	for _, s := range stages {
		if s.ID != stageID {
			continue
		}
		if s.CreatedAt.IsZero() {
			t.Error("CreatedAt ne devrait pas être zéro")
		}
		if s.UpdatedAt.IsZero() {
			t.Error("UpdatedAt ne devrait pas être zéro")
		}
		return
	}
	t.Errorf("Stage %d non trouvé", stageID)
}

// --- UpdateStage ---

func TestUpdateStage_CasNominal(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewStageRepository(db)
	ctx := context.Background()

	offerID := createTestStageOffer(t, db)
	stageID := createTestStage(t, db, offerID)

	stage, err := repo.GetStageByID(ctx, stageID)
	if err != nil || stage == nil {
		t.Fatalf("Stage introuvable : %v", err)
	}

	stage.Status = "in_progress"
	stage.UpdatedAt = time.Now().UTC().Truncate(time.Second)

	err = repo.UpdateStage(ctx, stage)
	if err != nil {
		t.Fatalf("Erreur inattendue lors de l'update : %v", err)
	}

	updated, err := repo.GetStageByID(ctx, stageID)
	if err != nil || updated == nil {
		t.Fatalf("Stage introuvable après update : %v", err)
	}
	if updated.Status != "in_progress" {
		t.Errorf("Status attendu 'in_progress', obtenu '%s'", updated.Status)
	}
}

func TestUpdateStage_IDInexistant_RetourneErreur(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewStageRepository(db)
	ctx := context.Background()

	stage := &models.Stage{
		ID:        999999,
		Status:    "in_progress",
		UpdatedAt: time.Now(),
	}

	err := repo.UpdateStage(ctx, stage)
	if err == nil {
		t.Error("Une erreur devrait être retournée pour un ID inexistant")
	}
}

func TestUpdateStage_UpdatedAtEstMisAJour(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewStageRepository(db)
	ctx := context.Background()

	offerID := createTestStageOffer(t, db)
	stageID := createTestStage(t, db, offerID)

	stage, _ := repo.GetStageByID(ctx, stageID)
	oldUpdatedAt := stage.UpdatedAt

	time.Sleep(time.Second) // garantit un écart de temps
	stage.Status = "completed"
	stage.UpdatedAt = time.Now().UTC().Truncate(time.Second)

	repo.UpdateStage(ctx, stage)

	updated, _ := repo.GetStageByID(ctx, stageID)
	if !updated.UpdatedAt.After(oldUpdatedAt) {
		t.Error("UpdatedAt devrait être postérieur à l'ancienne valeur")
	}
}

// --- CreateStage : timestamps ---

func TestCreateStage_TimestampsRetournés(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewStageRepository(db)
	ctx := context.Background()

	offerID := createTestStageOffer(t, db)
	now := time.Now().UTC().Truncate(time.Second)

	input := &models.Stage{
		StageOfferID: offerID,
		StudentID:    1,
		Status:       "pending",
		StartDate:    now,
		EndDate:      now.AddDate(0, 3, 0),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	created, err := repo.CreateStage(ctx, input)
	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
	t.Cleanup(func() { db.Exec(`DELETE FROM stages WHERE id = $1`, created.ID) })

	if created.CreatedAt.IsZero() {
		t.Error("CreatedAt ne devrait pas être zéro")
	}
	if created.UpdatedAt.IsZero() {
		t.Error("UpdatedAt ne devrait pas être zéro")
	}
}

// --- GetStages (avec contexte / rôle) ---

func buildContext(role models.UserRole, userID int) context.Context {
	claims := &models.CustomClaims{
		UserID: userID,
		Role:   role,
	}
	return context.WithValue(context.Background(), "claims", claims)
}

func TestGetStages_Etudiant_VoitSeulementSesStages(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewStageRepository(db)

	offerID := createTestStageOffer(t, db)

	// Stage appartenant à l'étudiant 1
	stageID := createTestStage(t, db, offerID)

	// Stage appartenant à l'étudiant 2
	var otherStageID int
	db.QueryRow(`
        INSERT INTO stages (stage_offer_id, student_id, status, start_date, end_date, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		offerID, 2, "pending", time.Now(), time.Now().AddDate(0, 3, 0), time.Now(), time.Now(),
	).Scan(&otherStageID)
	t.Cleanup(func() { db.Exec(`DELETE FROM stages WHERE id = $1`, otherStageID) })

	ctx := buildContext("student", 1)
	stages, err := repo.GetStages(ctx)
	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}

	for _, s := range stages {
		if s.ID == otherStageID {
			t.Error("L'étudiant ne devrait pas voir les stages des autres étudiants")
		}
	}

	found := false
	for _, s := range stages {
		if s.ID == stageID {
			found = true
		}
	}
	if !found {
		t.Errorf("L'étudiant devrait voir son propre stage %d", stageID)
	}
}

func TestGetStages_RoleInvalide_RetourneErreur(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewStageRepository(db)

	ctx := buildContext("unknown_role", 1)
	_, err := repo.GetStages(ctx)

	if err == nil {
		t.Error("Un rôle invalide devrait retourner une erreur")
	}
}

func TestGetStages_ContextSansClaims_RetourneErreur(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewStageRepository(db)

	_, err := repo.GetStages(context.Background())

	if err == nil {
		t.Error("Un contexte sans claims devrait retourner une erreur")
	}
}

func TestGetStages_Teacher_VoitSeulementSesStages(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewStageRepository(db)

	offerID := createTestStageOffer(t, db)

	var stageID int
	db.QueryRow(`
        INSERT INTO stages (stage_offer_id, student_id, teacher_id, status, start_date, end_date, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		offerID, 1, 10, "pending", time.Now(), time.Now().AddDate(0, 3, 0), time.Now(), time.Now(),
	).Scan(&stageID)
	t.Cleanup(func() { db.Exec(`DELETE FROM stages WHERE id = $1`, stageID) })

	ctx := buildContext("teacher", 10)
	stages, err := repo.GetStages(ctx)
	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}

	found := false
	for _, s := range stages {
		if s.ID == stageID {
			found = true
		}
	}
	if !found {
		t.Errorf("Le teacher devrait voir son stage %d", stageID)
	}
}

func TestGetStages_Tutor_VoitSeulementSesStages(t *testing.T) {
	db := test.SetupTestDB(t)
	repo := repositories.NewStageRepository(db)

	offerID := createTestStageOffer(t, db)

	var stageID int
	db.QueryRow(`
        INSERT INTO stages (stage_offer_id, student_id, tutor_id, status, start_date, end_date, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		offerID, 1, 20, "pending", time.Now(), time.Now().AddDate(0, 3, 0), time.Now(), time.Now(),
	).Scan(&stageID)
	t.Cleanup(func() { db.Exec(`DELETE FROM stages WHERE id = $1`, stageID) })

	ctx := buildContext("tutor", 20)
	stages, err := repo.GetStages(ctx)
	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}

	found := false
	for _, s := range stages {
		if s.ID == stageID {
			found = true
		}
	}
	if !found {
		t.Errorf("Le tutor devrait voir son stage %d", stageID)
	}
}
