package candidates

import (
	"context"

	"github/com/cl0ky/e-voting-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, c *models.Candidate) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Candidate, error)
	List(ctx context.Context, offset, limit int) ([]models.Candidate, int64, error)
	Update(ctx context.Context, c *models.Candidate) error
	Delete(ctx context.Context, id uuid.UUID, deletedBy *uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, c *models.Candidate) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*models.Candidate, error) {
	var c models.Candidate
	if err := r.db.WithContext(ctx).Preload("RT").Preload("Election").First(&c, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *repository) List(ctx context.Context, offset, limit int) ([]models.Candidate, int64, error) {
	var (
		items []models.Candidate
		total int64
	)
	q := r.db.WithContext(ctx).Model(&models.Candidate{})
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Preload("RT").Preload("Election").Order("created_at desc").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *repository) Update(ctx context.Context, c *models.Candidate) error {
	return r.db.WithContext(ctx).Model(&models.Candidate{}).Where("id = ?", c.Id).Updates(map[string]any{
		"name":        c.Name,
		"vision":      c.Vision,
		"mission":     c.Mission,
		"photo_url":   c.PhotoURL,
		"rt_id":       c.RTId,
		"election_id": c.ElectionId,
		"updated_by":  c.BaseModel.UpdatedBy,
	}).Error
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID, deletedBy *uuid.UUID) error {
	if err := r.db.WithContext(ctx).Model(&models.Candidate{}).Where("id = ?", id).Update("deleted_by", deletedBy).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Candidate{}).Error
}
