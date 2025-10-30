package models

import (
	"time"

	"github.com/google/uuid"
)

type Election struct {
	Id      uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	StartAt time.Time `gorm:"type:timestamp(6);not null;column:start_at"`
	EndAt   time.Time `gorm:"type:timestamp(6);not null;column:end_at"`
	Status  string    `gorm:"type:varchar(50);default:'upcoming';column:status"`

	BlockchainTxHash string `gorm:"type:varchar(128);column:blockchain_tx_hash"`

	RTId uuid.UUID `gorm:"type:uuid;not null;column:rt_id"`
	RT   RT        `gorm:"foreignKey:RTId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	BaseModel
}

func (Election) TableName() string {
	return "elections"
}
