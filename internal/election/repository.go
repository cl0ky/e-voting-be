package election

import (
	"context"
	"github/com/cl0ky/e-voting-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, e *models.Election) error
	GetAllByRTId(ctx context.Context, rtId uuid.UUID) ([]models.Election, error)
	GetById(ctx context.Context, id uuid.UUID) (*models.Election, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, e *models.Election) error {
	return r.db.WithContext(ctx).Create(e).Error
}

func (r *repository) GetAllByRTId(ctx context.Context, rtId uuid.UUID) ([]models.Election, error) {
	var elections []models.Election
	err := r.db.WithContext(ctx).Where("rt_id = ?", rtId).Find(&elections).Error
	return elections, err
}

func (r *repository) GetById(ctx context.Context, id uuid.UUID) (*models.Election, error) {
	var election models.Election
	err := r.db.WithContext(ctx).First(&election, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &election, nil
}

func (r *repository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return r.db.WithContext(ctx).Model(&models.Election{}).Where("id = ?", id).Update("status", status).Error
}
