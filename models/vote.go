package models

import (
	"time"

	"github.com/google/uuid"
)

type Vote struct {
	Id                  uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	VoterId             uuid.UUID  `gorm:"type:uuid;not null;column:voter_id"`
	ElectionId          uuid.UUID  `gorm:"type:uuid;not null;column:election_id"`
	HashVote            string     `gorm:"type:text;not null;column:hash_vote"`
	RevealedCandidateId *uuid.UUID `gorm:"type:uuid;column:revealed_candidate_id"`
	IsRevealed          bool       `gorm:"type:boolean;default:false;column:is_revealed"`
	RevealedAt          *time.Time `gorm:"type:timestamptz;column:revealed_at"`

	Voter     User       `gorm:"foreignKey:VoterId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Candidate *Candidate `gorm:"foreignKey:RevealedCandidateId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Election  Election   `gorm:"foreignKey:ElectionId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	BaseModel
}

func (Vote) TableName() string {
	return "votes"
}
