package models

import (
	"time"

	"github.com/google/uuid"
)

type Vote struct {
	Id          uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	VoterId     uuid.UUID  `gorm:"type:uuid;not null;column:voter_id"`
	CandidateId *uuid.UUID `gorm:"type:uuid;column:candidate_id"`
	ElectionId  uuid.UUID  `gorm:"type:uuid;not null;column:election_id"`

	HashVote     string     `gorm:"type:varchar(128);not null;column:hash_vote"`
	Nonce        string     `gorm:"type:text;column:nonce"`
	TxHashCommit string     `gorm:"type:varchar(128);column:tx_hash_commit"`
	TxHashReveal string     `gorm:"type:varchar(128);column:tx_hash_reveal"`
	IsRevealed   bool       `gorm:"type:boolean;default:false;column:is_revealed"`
	RevealedAt   *time.Time `gorm:"type:timestamptz;column:revealed_at"`

	Voter     User       `gorm:"foreignKey:VoterId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Candidate *Candidate `gorm:"foreignKey:CandidateId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Election  Election   `gorm:"foreignKey:ElectionId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	BaseModel
}

func (Vote) TableName() string {
	return "votes"
}
