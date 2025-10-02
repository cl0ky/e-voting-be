package models

import "github.com/google/uuid"

type Vote struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	ElectionID  uuid.UUID `gorm:"type:uuid;not null;column:election_id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;column:user_id"`
	CandidateID uuid.UUID `gorm:"type:uuid;not null;column:candidate_id"`

	Election  Election  `gorm:"foreignKey:ElectionID;references:ID"`
	User      User      `gorm:"foreignKey:UserID;references:ID"`
	Candidate Candidate `gorm:"foreignKey:CandidateID;references:ID"`

	BaseModel
}

func (Vote) TableName() string {
	return "votes"
}
