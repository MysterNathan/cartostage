package repositories

import (
	"context"
	"log"
	sharedContext "shared/context"
	"shared/models"

	"github.com/jmoiron/sqlx"
)

type EnterpriseRepository struct {
	db *sqlx.DB
}

func NewEnterpriseRepository(db *sqlx.DB) *EnterpriseRepository {
	return &EnterpriseRepository{db: db}
}

//
//// GetByID - Récupère une entreprise par son ID
//func (r *EnterpriseRepository) GetByID(id int) (*models.Enterprise, error) {
//	var enterprise models.Enterprise
//	query := `
//        SELECT id, nom, adresse, secteur, taille, siret, email,
//               telephone, site_web, description, logo_url,
//               created_at, updated_at
//        FROM enterprises
//        WHERE id = $1
//    `
//
//	err := r.db.Get(&enterprise, query, id)
//	if err != nil {
//		if err == sql.ErrNoRows {
//			return nil, errors.New("enterprise not found")
//		}
//		return nil, err
//	}
//
//	return &enterprise, nil
//}
//
//// GetStagesByEnterpriseID - Récupère tous les stages d'une entreprise
//func (r *EnterpriseRepository) GetStagesByEnterpriseID(enterpriseID int) ([]*models.Stage, error) {
//	var stages []*models.Stage
//	query := `
//        SELECT s.id, s.poste, s.adresse, s.lat, s.lng,
//               s.places_disponibles, s.entreprise, s.filiere,
//               s.sector, s.commune, s.capacity_total,
//               s.capacity_filled, s.period, s.enterprise_id,
//               s.created_at, s.updated_at
//        FROM stages s
//        WHERE s.enterprise_id = $1
//        ORDER BY s.created_at DESC
//    `
//
//	err := r.db.Select(&stages, query, enterpriseID)
//	if err != nil {
//		return nil, err
//	}
//
//	return stages, nil
//}
//
//// GetStageByID - Récupère un stage par son ID
//func (r *EnterpriseRepository) GetStageByID(stageID int) (*models.Stage, error) {
//	var stage models.Stage
//	query := `
//        SELECT id, poste, adresse, lat, lng, places_disponibles,
//               entreprise, filiere, sector, commune, capacity_total,
//               capacity_filled, period, enterprise_id, created_at, updated_at
//        FROM stages
//        WHERE id = $1
//    `
//
//	err := r.db.Get(&stage, query, stageID)
//	if err != nil {
//		if err == sql.ErrNoRows {
//			return nil, errors.New("stage not found")
//		}
//		return nil, err
//	}
//
//	return &stage, nil
//}
//
//// GetAll - Récupère toutes les entreprises
//func (r *EnterpriseRepository) GetAll() ([]*models.Enterprise, error) {
//	var enterprises []*models.Enterprise
//	query := `
//        SELECT id, nom, adresse, secteur, taille, siret, email,
//               telephone, site_web, description, logo_url,
//               created_at, updated_at
//        FROM enterprises
//        ORDER BY nom ASC
//    `
//
//	err := r.db.Select(&enterprises, query)
//	if err != nil {
//		return nil, err
//	}
//
//	return enterprises, nil
//}
//
//// Create - Crée une nouvelle entreprise
//func (r *EnterpriseRepository) Create(enterprise *models.Enterprise) (*models.Enterprise, error) {
//	query := `
//        INSERT INTO enterprises (nom, adresse, secteur, taille, siret,
//                                email, telephone, site_web, description, logo_url)
//        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
//        RETURNING id, created_at, updated_at
//    `
//
//	err := r.db.QueryRow(query,
//		enterprise.Nom, enterprise.Adresse, enterprise.Secteur,
//		enterprise.Taille, enterprise.Siret, enterprise.Email,
//		enterprise.Telephone, enterprise.SiteWeb, enterprise.Description,
//		enterprise.LogoURL,
//	).Scan(&enterprise.ID, &enterprise.CreatedAt, &enterprise.UpdatedAt)
//
//	if err != nil {
//		return nil, err
//	}
//
//	return enterprise, nil
//}
//
//// Update - Met à jour une entreprise
//func (r *EnterpriseRepository) Update(id int, enterprise *models.Enterprise) (*models.Enterprise, error) {
//	query := `
//        UPDATE enterprises
//        SET nom = $1, adresse = $2, secteur = $3, taille = $4,
//            siret = $5, email = $6, telephone = $7, site_web = $8,
//            description = $9, logo_url = $10, updated_at = CURRENT_TIMESTAMP
//        WHERE id = $11
//        RETURNING created_at, updated_at
//    `
//
//	err := r.db.QueryRow(query,
//		enterprise.Nom, enterprise.Adresse, enterprise.Secteur,
//		enterprise.Taille, enterprise.Siret, enterprise.Email,
//		enterprise.Telephone, enterprise.SiteWeb, enterprise.Description,
//		enterprise.LogoURL, id,
//	).Scan(&enterprise.CreatedAt, &enterprise.UpdatedAt)
//
//	if err != nil {
//		return nil, err
//	}
//
//	enterprise.ID = id
//	return enterprise, nil
//}
//
//// Delete - Supprime une entreprise
//func (r *EnterpriseRepository) Delete(id int) error {
//	query := `DELETE FROM enterprises WHERE id = $1`
//	result, err := r.db.Exec(query, id)
//	if err != nil {
//		return err
//	}
//
//	rowsAffected, err := result.RowsAffected()
//	if err != nil {
//		return err
//	}
//
//	if rowsAffected == 0 {
//		return errors.New("enterprise not found")
//	}
//
//	return nil
//}
//
//// GetWithStats - Récupère une entreprise avec ses statistiques
//func (r *EnterpriseRepository) GetWithStats(id int) (*models.EnterpriseWithStats, error) {
//	var result models.EnterpriseWithStats
//
//	// D'abord récupérer l'entreprise
//	enterprise, err := r.GetByID(id)
//	if err != nil {
//		return nil, err
//	}
//
//	result.Enterprise = *enterprise
//
//	// Puis récupérer les stats (pour l'instant des valeurs par défaut)
//	// TODO: implémenter les vraies requêtes de stats quand les tables seront prêtes
//	result.TotalTutors = 0
//	result.ActiveStages = 0
//	result.TotalStudents = 0
//
//	return &result, nil
//}

func (r *EnterpriseRepository) GetStats(ctx context.Context) (*models.EnterpriseStats, error) {
	claims := sharedContext.GetClaimsFromContext(ctx)
	var result models.EnterpriseStats

	// Nombre de stages supervisés par l'utilisateur
	queryStages := `SELECT COUNT(*) FROM stages WHERE tutor_id = $1;`
	if err := r.db.QueryRowContext(ctx, queryStages, claims.UserID).Scan(&result.ActiveStages); err != nil {
		return nil, err
	}

	// Nombre de tuteurs dans l'établissement de l'utilisateur
	queryTutors := `
		SELECT COUNT(DISTINCT s.tutor_id) AS nb_tuteurs_uniques
		FROM stages s
         JOIN users u ON u.establishment_id = s.establishment_id
		WHERE u.id = $1
		`
	if err := r.db.QueryRowContext(ctx, queryTutors, claims.UserID).Scan(&result.TotalTutors); err != nil {
		return nil, err
	}
	log.Println(result.TotalTutors)
	return &result, nil
}
