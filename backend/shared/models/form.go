package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type Form struct {
	ID          int        `json:"id" db:"id"`
	StageID     int        `json:"stage_id" db:"stage_id"`
	StudentID   int        `json:"student_id" db:"student_id"`
	TeacherID   *int       `json:"teacher_id,omitempty" db:"teacher_id"`
	TutorID     *int       `json:"tutor_id,omitempty" db:"tutor_id"`
	Status      string     `json:"status" db:"status"`
	Content     *JSONB     `json:"content,omitempty" db:"content"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty" db:"completed_at"`
}

const (
	StatusCreated    = "CREATED"
	StatusInProgress = "IN_PROGRESS"
	StatusCompleted  = "COMPLETED"
	StatusCancelled  = "CANCELLED"
)

type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}
