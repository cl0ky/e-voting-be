package models

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	Name     string    `gorm:"type:varchar(255);not null;column:name"`
	Email    string    `gorm:"type:varchar(255);uniqueIndex;not null;column:email"`
	Password string    `gorm:"type:varchar(255);not null;column:password"`
	NIK      string    `gorm:"type:varchar(50);uniqueIndex;not null;column:nik"`
	Role     string    `gorm:"type:varchar(50);default:'voter';not null;column:role"`

	BaseModel
}

func (User) TableName() string {
	return "users"
}
