package models

import "github.com/google/uuid"

type Candidate struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	ElectionID uuid.UUID `gorm:"type:uuid;not null;column:election_id"`
	Name       string    `gorm:"type:varchar(255);not null;column:name"`
	Vision     string    `gorm:"type:text;column:vision"`
	Mission    string    `gorm:"type:text;column:mission"`

	Election Election `gorm:"foreignKey:ElectionID;references:ID"`

	BaseModel
}

func (Candidate) TableName() string {
	return "candidates"
}
