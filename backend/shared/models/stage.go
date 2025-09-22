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

type Content struct {
	ID      int     `json:"id" db:"id"`
	Content *string `json:"content,omitempty" db:"content"`
}
