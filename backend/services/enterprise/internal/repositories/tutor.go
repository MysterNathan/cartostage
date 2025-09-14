package repositories

import (
	"github.com/jmoiron/sqlx"
	"shared/models"
)

type TutorRepository struct {
	db *sqlx.DB
}

func NewTutorRepository(db *sqlx.DB) *TutorRepository { // Changement ici
	return &TutorRepository{db: db}
}

func (r *TutorRepository) GetAll() ([]models.TutorWithEnterprise, error) {
	var tutors []models.TutorWithEnterprise
	query := `
        SELECT t.*, e.nom as enterprise_name 
        FROM tutors t 
        JOIN enterprises e ON t.enterprise_id = e.id 
        ORDER BY t.created_at DESC`

	err := r.db.Select(&tutors, query) // Utilisation de Select
	return tutors, err
}

func (r *TutorRepository) GetByID(id int) (*models.TutorWithEnterprise, error) {
	var tutor models.TutorWithEnterprise
	query := `
        SELECT t.*, e.nom as enterprise_name 
        FROM tutors t 
        JOIN enterprises e ON t.enterprise_id = e.id 
        WHERE t.id = $1`

	err := r.db.Get(&tutor, query, id) // Utilisation de Get
	if err != nil {
		return nil, err
	}
	return &tutor, nil
}

func (r *TutorRepository) Create(tutor *models.Tutor) error {
	query := `INSERT INTO tutors (enterprise_id, prenom, nom, email, telephone, poste, departement, is_active, created_at, updated_at) 
              VALUES (:enterprise_id, :prenom, :nom, :email, :telephone, :poste, :departement, :is_active, NOW(), NOW()) RETURNING id`

	rows, err := r.db.NamedQuery(query, tutor)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		return rows.Scan(&tutor.ID)
	}
	return nil
}

func (r *TutorRepository) Update(tutor *models.Tutor) error {
	query := `UPDATE tutors SET enterprise_id = :enterprise_id, prenom = :prenom, nom = :nom, 
              email = :email, telephone = :telephone, poste = :poste, departement = :departement, 
              is_active = :is_active, updated_at = NOW() WHERE id = :id`

	_, err := r.db.NamedExec(query, tutor)
	return err
}

func (r *TutorRepository) Delete(id int) error {
	query := `DELETE FROM tutors WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *TutorRepository) GetByEnterprise(enterpriseID int) ([]models.TutorWithEnterprise, error) {
	var tutors []models.TutorWithEnterprise
	query := `
        SELECT t.*, e.nom as enterprise_name 
        FROM tutors t 
        JOIN enterprises e ON t.enterprise_id = e.id 
        WHERE t.enterprise_id = $1 
        ORDER BY t.created_at DESC`

	err := r.db.Select(&tutors, query, enterpriseID)
	return tutors, err
}

func (r *TutorRepository) GetWithStats(id int) (*models.TutorWithStats, error) {
	var tutor models.TutorWithStats
	query := `
        SELECT t.*,
               COALESCE(student_count.active, 0) as active_students,
               COALESCE(stage_count.total, 0) as total_stages
        FROM tutors t
        LEFT JOIN (
            SELECT COUNT(*) as active 
            FROM stages s 
            WHERE s.entreprise = (SELECT e.nom FROM enterprises e WHERE e.id = t.enterprise_id)
        ) student_count ON true
        LEFT JOIN (
            SELECT COUNT(*) as total 
            FROM stages s 
            WHERE s.entreprise = (SELECT e.nom FROM enterprises e WHERE e.id = t.enterprise_id)
        ) stage_count ON true
        WHERE t.id = $1`

	err := r.db.Get(&tutor, query, id)
	if err != nil {
		return nil, err
	}
	return &tutor, nil
}
