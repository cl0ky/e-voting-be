package models

import "github.com/google/uuid"

type User struct {
	Id         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	Name       string    `gorm:"type:varchar(255);not null;column:name"`
	Email      string    `gorm:"type:varchar(255);uniqueIndex;not null;column:email"`
	Password   string    `gorm:"type:varchar(255);not null;column:password"`
	Role       string    `gorm:"type:varchar(50);default:'Voter';not null;column:role"`
	NIK        string    `gorm:"type:varchar(50);uniqueIndex;not null;column:nik"`
	IsVerified bool      `gorm:"type:boolean;default:false;column:is_verified"`

	RTId *uuid.UUID `gorm:"type:uuid;column:rt_id"`

	BaseModel
}

func (User) TableName() string {
	return "users"
}
