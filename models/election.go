package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Election struct {
	Id      uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	Name    string    `gorm:"type:varchar(100);not null;column:name"`
	StartAt time.Time `gorm:"type:timestamptz;not null;column:start_at"`
	EndAt   time.Time `gorm:"type:timestamptz;not null;column:end_at"`
	Status  string    `gorm:"type:varchar(50);default:'upcoming';column:status"`

	BlockchainTxHash string         `gorm:"type:varchar(128);column:blockchain_tx_hash"`
	SummaryHash      string         `gorm:"type:text;column:summary_hash"`
	SummaryJSON      datatypes.JSON `gorm:"type:jsonb;column:summary_json"`
	FinalizeStatus   string         `gorm:"type:varchar(20);default:'pending';column:finalize_status"`
	FinalizeError    string         `gorm:"type:text;column:finalize_error"`
	FinalizedAt      *time.Time     `gorm:"type:timestamptz;column:finalized_at"`

	RTId uuid.UUID `gorm:"type:uuid;not null;column:rt_id"`
	RT   RT        `gorm:"foreignKey:RTId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	Votes      []Vote      `gorm:"foreignKey:ElectionId;references:Id"`
	Candidates []Candidate `gorm:"foreignKey:ElectionId;references:Id"`

	BaseModel
}

func (Election) TableName() string {
	return "elections"
}
