package repositories

import (
	"github.com/jmoiron/sqlx"
	"shared/models"
)

type EnterpriseRepository struct {
	db *sqlx.DB
}

func NewEnterpriseRepository(db *sqlx.DB) *EnterpriseRepository {
	return &EnterpriseRepository{db: db}
}

func (r *EnterpriseRepository) GetAll() ([]models.Enterprise, error) {
	query := `
        SELECT id, nom, adresse, secteur, taille, siret, email_contact, 
               telephone, site_web, description, logo_url, created_at, updated_at
        FROM enterprises ORDER BY nom
    `
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enterprises []models.Enterprise
	for rows.Next() {
		var e models.Enterprise
		err := rows.Scan(
			&e.ID, &e.Nom, &e.Adresse, &e.Secteur, &e.Taille, &e.Siret,
			&e.Email, &e.Telephone, &e.SiteWeb, &e.Description, &e.LogoURL,
			&e.CreatedAt, &e.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		enterprises = append(enterprises, e)
	}
	return enterprises, nil
}

func (r *EnterpriseRepository) GetByID(id int) (*models.Enterprise, error) {
	query := `
        SELECT id, nom, adresse, secteur, taille, siret, email_contact, 
               telephone, site_web, description, logo_url, created_at, updated_at
        FROM enterprises WHERE id = $1
    `
	var e models.Enterprise
	err := r.db.QueryRow(query, id).Scan(
		&e.ID, &e.Nom, &e.Adresse, &e.Secteur, &e.Taille, &e.Siret,
		&e.Email, &e.Telephone, &e.SiteWeb, &e.Description, &e.LogoURL,
		&e.CreatedAt, &e.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *EnterpriseRepository) Create(e *models.Enterprise) error {
	query := `
        INSERT INTO enterprises (nom, adresse, secteur, taille, siret, email_contact, 
                                telephone, site_web, description, logo_url)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING id, created_at, updated_at
    `
	return r.db.QueryRow(
		query, e.Nom, e.Adresse, e.Secteur, e.Taille, e.Siret,
		e.Email, e.Telephone, e.SiteWeb, e.Description, e.LogoURL,
	).Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt)
}

func (r *EnterpriseRepository) Update(e *models.Enterprise) error {
	query := `
        UPDATE enterprises 
        SET nom = $2, adresse = $3, secteur = $4, taille = $5, siret = $6,
            email_contact = $7, telephone = $8, site_web = $9, 
            description = $10, logo_url = $11, updated_at = NOW()
        WHERE id = $1
        RETURNING updated_at
    `
	return r.db.QueryRow(
		query, e.ID, e.Nom, e.Adresse, e.Secteur, e.Taille, e.Siret,
		e.Email, e.Telephone, e.SiteWeb, e.Description, e.LogoURL,
	).Scan(&e.UpdatedAt)
}

func (r *EnterpriseRepository) Delete(id int) error {
	query := `DELETE FROM enterprises WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *EnterpriseRepository) GetWithStats(id int) (*models.EnterpriseWithStats, error) {
	query := `
        SELECT e.id, e.nom, e.adresse, e.secteur, e.taille, e.siret, e.email_contact,
               e.telephone, e.site_web, e.description, e.logo_url, e.created_at, e.updated_at,
               COALESCE(COUNT(DISTINCT t.id), 0) as total_tutors,
               COALESCE(COUNT(DISTINCT s.id), 0) as active_stages,
               COALESCE(COUNT(DISTINCT st.id), 0) as total_students
        FROM enterprises e
        LEFT JOIN tutors t ON e.id = t.enterprise_id AND t.is_active = true
        LEFT JOIN stages s ON e.nom = s.enterprise
        LEFT JOIN students st ON s.id = st.stage_id
        WHERE e.id = $1
        GROUP BY e.id
    `
	var ews models.EnterpriseWithStats
	err := r.db.QueryRow(query, id).Scan(
		&ews.ID, &ews.Nom, &ews.Adresse, &ews.Secteur, &ews.Taille, &ews.Siret,
		&ews.Email, &ews.Telephone, &ews.SiteWeb, &ews.Description, &ews.LogoURL,
		&ews.CreatedAt, &ews.UpdatedAt, &ews.TotalTutors, &ews.ActiveStages, &ews.TotalStudents,
	)
	if err != nil {
		return nil, err
	}
	return &ews, nil
}

// GetMe - Alias pour GetWithStats pour l'entreprise connectée
func (r *EnterpriseRepository) GetMe(enterpriseID int) (*models.EnterpriseWithStats, error) {
	return r.GetWithStats(enterpriseID)
}
