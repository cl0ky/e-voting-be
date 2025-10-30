package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	CreatedAt time.Time      `gorm:"type:timestamp(6);default:now();column:created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamp(6);default:now();column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamp(6);column:deleted_at;index"`

	CreatedBy *uuid.UUID `gorm:"type:uuid;column:created_by"`
	UpdatedBy *uuid.UUID `gorm:"type:uuid;column:updated_by"`
	DeletedBy *uuid.UUID `gorm:"type:uuid;column:deleted_by"`

	UserCreatedBy *User `gorm:"foreignKey:CreatedBy;references:Id;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UserUpdatedBy *User `gorm:"foreignKey:UpdatedBy;references:Id;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UserDeletedBy *User `gorm:"foreignKey:DeletedBy;references:Id;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}
