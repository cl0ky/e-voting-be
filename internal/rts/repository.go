package rts

import (
	"github/com/cl0ky/e-voting-be/models"

	"gorm.io/gorm"
)

type Repository interface {
	GetAllRT() ([]models.RT, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAllRT() ([]models.RT, error) {
	var rts []models.RT
	if err := r.db.Find(&rts).Error; err != nil {
		return nil, err
	}
	return rts, nil
}
