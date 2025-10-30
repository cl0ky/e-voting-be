package models

import "github.com/google/uuid"

type RT struct {
	Id      uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	Name    string     `gorm:"type:varchar(100);not null;column:name"`
	Region  string     `gorm:"type:varchar(255);not null;column:region"`
	AdminId *uuid.UUID `gorm:"type:uuid;column:admin_id"`

	Admin *User `gorm:"foreignKey:AdminId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	BaseModel
}

func (RT) TableName() string {
	return "rts"
}
