package models

import (
	"time"

	"github.com/google/uuid"
)

type Election struct {
	Id      uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	Name    string    `gorm:"type:varchar(100);not null;column:name"`
	StartAt time.Time `gorm:"type:timestamp(6);not null;column:start_at"`
	EndAt   time.Time `gorm:"type:timestamp(6);not null;column:end_at"`
	Status  string    `gorm:"type:varchar(50);default:'upcoming';column:status"`

	BlockchainTxHash string `gorm:"type:varchar(128);column:blockchain_tx_hash"`

	RTId uuid.UUID `gorm:"type:uuid;not null;column:rt_id"`
	RT   RT        `gorm:"foreignKey:RTId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	Votes      []Vote      `gorm:"foreignKey:ElectionId;references:Id"`
	Candidates []Candidate `gorm:"foreignKey:RTId;references:Id"`

	BaseModel
}

func (Election) TableName() string {
	return "elections"
}
