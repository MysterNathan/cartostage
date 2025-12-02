package models

type EnterpriseStats struct {
	TotalTutors   int `json:"tutors" db:"tutors"`
	ActiveStages  int `json:"stages" db:"stages"`
	TotalStudents int `json:"student" db:"student"`
}
