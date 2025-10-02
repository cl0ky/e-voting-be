package models

import "github.com/google/uuid"

type Election struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	Title       string    `gorm:"type:varchar(255);not null;column:title"`
	Description string    `gorm:"type:text;column:description"`
	IsActive    bool      `gorm:"default:false;column:is_active"`

	Candidates     []Candidate     `gorm:"foreignKey:ElectionID"`
	EligibleVoters []EligibleVoter `gorm:"foreignKey:ElectionID"`
	Votes          []Vote          `gorm:"foreignKey:ElectionID"`

	BaseModel
}

func (Election) TableName() string {
	return "elections"
}
