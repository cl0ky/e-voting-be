package models

import "github.com/google/uuid"

type Candidate struct {
	Id       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	Name     string    `gorm:"type:varchar(255);not null;column:name"`
	Vision   string    `gorm:"type:text;column:vision"`
	Mission  string    `gorm:"type:text;column:mission"`
	PhotoURL string    `gorm:"type:varchar(255);column:photo_url"`

	RTId uuid.UUID `gorm:"type:uuid;not null;column:rt_id"`
	RT   RT        `gorm:"foreignKey:RTId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	ElectionId *uuid.UUID `gorm:"type:uuid;column:election_id"`
	Election   *Election  `gorm:"foreignKey:ElectionId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	BaseModel
}

func (Candidate) TableName() string {
	return "candidates"
}
