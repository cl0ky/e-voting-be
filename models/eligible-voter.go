package models

import "github.com/google/uuid"

type EligibleVoter struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	ElectionID uuid.UUID `gorm:"type:uuid;not null;column:election_id"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;column:user_id"`
	HasVoted   bool      `gorm:"default:false;column:has_voted"`

	Election *Election `gorm:"foreignKey:ElectionID;references:ID"`
	User     *User     `gorm:"foreignKey:UserID;references:ID"`

	BaseModel
}

func (EligibleVoter) TableName() string {
	return "eligible_voters"
}
