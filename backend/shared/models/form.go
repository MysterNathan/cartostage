package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"time"
)

type FormCreationData struct {
	StageID   int  `json:"stage_id"`
	StudentID int  `json:"student_id"`
	TeacherID *int `json:"teacher_id"`
	TutorID   *int `json:"tutor_id"`
}

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

type FormSection struct {
	ID          int        `json:"id" db:"id"`
	FormID      int        `json:"form_id" db:"form_id"`
	SectionType string     `json:"section_type" db:"section_type"`
	UserID      int        `json:"user_id" db:"user_id"`
	Status      string     `json:"status" db:"status"`
	Content     *JSONB     `json:"content,omitempty" db:"content"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty" db:"completed_at"`
}

type FormFormSection struct {
	Form               *Form         `json:"form"`
	FormSections       []FormSection `json:"form_section"`
	FormSectionContent *JSONB        `json:"form_section_content"`
}

type Skill struct {
	Name  string `json:"name" validate:"required,max=100" db:"-"`
	Level int    `json:"level" validate:"omitempty,min=1,max=5" db:"-"`
}

type FormSectionContentStudent struct {
	CompanyIntegration      int     `json:"company_integration" validate:"required,min=1,max=5" db:"-"`
	UnderstandingTasks      int     `json:"understanding_tasks" validate:"required,min=1,max=5" db:"-"`
	AutonomyLevel           int     `json:"autonomy_level" validate:"required,min=1,max=5" db:"-"`
	SkillsAcquired          []Skill `json:"skills_acquired" validate:"omitempty,max=20,dive" db:"-"`
	PositivePoint           string  `json:"positive_point" validate:"omitempty,max=250" db:"-"`
	DifficultiesEncountered string  `json:"difficulties_encountered" validate:"omitempty,max=250" db:"-"`
	Comment                 string  `json:"comment" validate:"omitempty,max=250" db:"-"`
}

type FormSectionContentTeacher struct {
	Attendance          int    `json:"attendance" validate:"required,min=1,max=5" db:"-"`
	ReportQuality       int    `json:"report_quality" validate:"required,min=1,max=5" db:"-"`
	OralPresentation    int    `json:"oral_presentation" validate:"required,min=1,max=5" db:"-"`
	ProfessionalConduct int    `json:"professional_conduct" validate:"required,min=1,max=5" db:"-"`
	ObjectivesAchieved  bool   `json:"objectives_achieved" validate:"required" db:"-"`
	Grade               int    `json:"grade" validate:"required,min=1,max=20" db:"-"`
	Recommandations     string `json:"recommandations" validate:"omitempty,max=250" db:"-"`
}

type Task struct {
	Name    string  `json:"name" validate:"required,max=100" db:"-"`
	Summary string  `json:"summary" validate:"omitempty,max=250" db:"-"`
	Skills  []Skill `json:"skills" validate:"omitempty,dive" db:"-"`
}

type FormSectionContentTutor struct {
	TechnicalSkills int    `json:"technical_skills" validate:"required,min=1,max=5" db:"-"`
	WorkQuality     int    `json:"work_quality" validate:"required,min=1,max=5" db:"-"`
	Punctuality     int    `json:"punctuality" validate:"required,min=1,max=5" db:"-"`
	TeamIntegration int    `json:"team_integration" validate:"required,min=1,max=5" db:"-"`
	Autonomy        int    `json:"autonomy" validate:"required,min=1,max=5" db:"-"`
	TasksCompleted  []Task `json:"tasks_completed" validate:"omitempty,max=20,dive" db:"-"`
	Comment         string `json:"comment" validate:"omitempty,max=250" db:"-"`
}

const (
	StatusCreated    = "CREATED"
	StatusInProgress = "IN_PROGRESS"
	StatusCompleted  = "COMPLETED"
	StatusCancelled  = "CANCELLED"
)

func (fs *FormSection) GetContent() (interface{}, error) {
	if fs.Content == nil {
		return nil, fmt.Errorf("content is nil")
	}

	validate := validator.New()

	switch fs.SectionType {
	case "Student":
		var content FormSectionContentStudent
		if err := json.Unmarshal(*fs.Content, &content); err != nil {
			return nil, fmt.Errorf("error parsing JSON student: %w", err)
		}
		if err := validate.Struct(content); err != nil {
			return nil, fmt.Errorf("validation student failed: %w", err)
		}
		return content, nil

	case "Teacher":
		var content FormSectionContentTeacher
		if err := json.Unmarshal(*fs.Content, &content); err != nil {
			return nil, fmt.Errorf("error parsing JSON teacher: %w", err)
		}
		if err := validate.Struct(content); err != nil {
			return nil, fmt.Errorf("validation teacher failed: %w", err)
		}
		return content, nil

	case "Tutor":
		var content FormSectionContentTutor
		if err := json.Unmarshal(*fs.Content, &content); err != nil {
			return nil, fmt.Errorf("error parsing JSON tutor: %w", err)
		}
		if err := validate.Struct(content); err != nil {
			return nil, fmt.Errorf("validation tutor failed: %w", err)
		}
		return content, nil

	default:
		return nil, fmt.Errorf("unknown section type: %s", fs.SectionType)
	}
}

type JSONB []byte

func (j JSONB) MarshalJSON() ([]byte, error) {
	if len(j) == 0 {
		return []byte("null"), nil
	}
	return j, nil
}

func (j *JSONB) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*j = nil
		return nil
	}
	if !json.Valid(data) {
		return errors.New("JSON invalide pour le champ JSONB")
	}
	*j = data
	return nil
}
func (j JSONB) Value() (driver.Value, error) {
	if j == nil || len(j) == 0 {
		return nil, nil
	}
	return []byte(j), nil
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		*j = make([]byte, len(v))
		copy(*j, v)
		return nil
	case string:
		*j = []byte(v)
		return nil
	default:
		return fmt.Errorf("unsuported type for JSONB: %T", value)
	}
}

func (j JSONB) String() string {
	if j == nil {
		return "null"
	}
	return string(j)
}
