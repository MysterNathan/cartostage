package models

import (
	"time"
)

type Stage struct {
	ID              int       `json:"id" db:"id"`
	StageOfferID    int       `json:"stage_offer_id" db:"stage_offer_id"`
	StudentID       int       `json:"student_id" db:"student_id"`
	TeacherID       *int      `json:"teacher_id,omitempty" db:"teacher_id"`
	TutorID         *int      `json:"tutor_id,omitempty" db:"tutor_id"`
	EstablishmentID *int      `json:"establishment_id,omitempty" db:"establishment_id"`
	ContentID       *int      `json:"content_id,omitempty" db:"content_id"`
	Status          string    `json:"status" db:"status"`
	StartDate       time.Time `json:"start_date" db:"start_date"`
	EndDate         time.Time `json:"end_date" db:"end_date"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// Structure avec relations jointes pour l'affichage
type StageWithDetails struct {
	Stage
	StageOffer    *StageOffer    `json:"stage_offer,omitempty"`
	Student       *UserPublic    `json:"student,omitempty"`
	Teacher       *UserPublic    `json:"teacher,omitempty"`
	Tutor         *UserPublic    `json:"tutor,omitempty"`
	Establishment *Establishment `json:"establishment,omitempty"`
	Content       *Content       `json:"content,omitempty"`
}

// Méthode utilitaire pour calculer la durée
func (s *Stage) Duration() int {
	return int(s.EndDate.Sub(s.StartDate).Hours() / 24)
}

// Méthode pour vérifier si le stage est actif
func (s *Stage) IsActive() bool {
	now := time.Now()
	return s.Status == "in_progress" &&
		now.After(s.StartDate) &&
		now.Before(s.EndDate)
}

// Conversion en version publique
func (s *Stage) ToPublic() StagePublic {
	return StagePublic{
		StageOfferID: s.StageOfferID,
		StudentID:    s.StudentID,
		TeacherID:    s.TeacherID,
		TutorID:      s.TutorID,
		Status:       s.Status,
		StartDate:    s.StartDate,
		EndDate:      s.EndDate,
		CreatedAt:    s.CreatedAt,
	}
}

// Structure pour les réponses publiques
type StagePublic struct {
	ID           int       `json:"id"`
	StageOfferID int       `json:"stage_offer_id"`
	StudentID    int       `json:"student_id"`
	TeacherID    *int      `json:"teacher_id,omitempty"`
	TutorID      *int      `json:"tutor_id,omitempty"`
	Status       string    `json:"status"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	CreatedAt    time.Time `json:"created_at"`
}

// Structures pour les requêtes de création/mise à jour
type CreateStageRequest struct {
	StageOfferID    int       `json:"stage_offer_id" binding:"required"`
	StudentID       int       `json:"student_id" binding:"required"`
	TeacherID       *int      `json:"teacher_id,omitempty"`
	TutorID         *int      `json:"tutor_id,omitempty"`
	EstablishmentID *int      `json:"establishment_id,omitempty"`
	ContentID       *int      `json:"content_id,omitempty"`
	Status          string    `json:"status" binding:"required"`
	StartDate       time.Time `json:"start_date" binding:"required"`
	EndDate         time.Time `json:"end_date" binding:"required"`
}

type UpdateStageRequest struct {
	TeacherID *int       `json:"teacher_id,omitempty"`
	TutorID   *int       `json:"tutor_id,omitempty"`
	ContentID *int       `json:"content_id,omitempty"`
	Status    *string    `json:"status,omitempty"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
}

// Méthode pour vérifier s'il y a des champs à mettre à jour
func (r *UpdateStageRequest) HasUpdates() bool {
	return r.TeacherID != nil || r.TutorID != nil || r.ContentID != nil ||
		r.Status != nil || r.StartDate != nil || r.EndDate != nil
}

// Méthode pour appliquer les modifications à un stage existant
func (r *UpdateStageRequest) ApplyTo(stage *Stage) {
	if r.TeacherID != nil {
		stage.TeacherID = r.TeacherID
	}
	if r.TutorID != nil {
		stage.TutorID = r.TutorID
	}
	if r.ContentID != nil {
		stage.ContentID = r.ContentID
	}
	if r.Status != nil {
		stage.Status = *r.Status
	}
	if r.StartDate != nil {
		stage.StartDate = *r.StartDate
	}
	if r.EndDate != nil {
		stage.EndDate = *r.EndDate
	}

	stage.UpdatedAt = time.Now()
}

type DeleteStageRequest struct {
	Id int `json:"id" binding:"required"`
}

type Content struct {
	ID      int     `json:"id" db:"id"`
	Content *string `json:"content,omitempty" db:"content"`
}
