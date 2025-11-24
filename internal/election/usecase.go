package election

import (
	"context"
	"fmt"
	"github/com/cl0ky/e-voting-be/models"

	"github.com/google/uuid"
)

type UseCase interface {
	Create(ctx context.Context, req CreateElectionRequest) (*ElectionItem, error)
	GetAllByRTId(ctx context.Context, rtId uuid.UUID) ([]ElectionItem, error)
	GetById(ctx context.Context, id uuid.UUID) (*ElectionItem, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
}

type useCase struct {
	repo Repository
}

func NewUseCase(repo Repository) UseCase {
	return &useCase{repo: repo}
}

func (u *useCase) Create(ctx context.Context, req CreateElectionRequest) (*ElectionItem, error) {
	e := &models.Election{
		Name:    req.Name,
		StartAt: req.StartAt,
		EndAt:   req.EndAt,
		Status:  req.Status,
		RTId:    req.RTId,
	}
	if req.CreatedBy != nil {
		e.CreatedBy = req.CreatedBy
	}

	if err := u.repo.Create(ctx, e); err != nil {
		return nil, err
	}
	item := ElectionItem{
		Id:      e.Id,
		Name:    e.Name,
		StartAt: e.StartAt,
		EndAt:   e.EndAt,
		Status:  e.Status,
		RTId:    e.RTId,
	}
	return &item, nil
}

func (u *useCase) GetAllByRTId(ctx context.Context, rtId uuid.UUID) ([]ElectionItem, error) {
	elections, err := u.repo.GetAllByRTId(ctx, rtId)
	if err != nil {
		return []ElectionItem{}, err
	}
	items := make([]ElectionItem, 0)
	for _, e := range elections {
		items = append(items, ElectionItem{
			Id:      e.Id,
			Name:    e.Name,
			StartAt: e.StartAt,
			EndAt:   e.EndAt,
			Status:  e.Status,
			RTId:    e.RTId,
			Year:    e.StartAt.Year(),
		})
	}

	return items, nil
}

func (u *useCase) GetById(ctx context.Context, id uuid.UUID) (*ElectionItem, error) {
	e, err := u.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	fmt.Println("Election found:", e)
	item := &ElectionItem{
		Id:      e.Id,
		Name:    e.Name,
		StartAt: e.StartAt,
		EndAt:   e.EndAt,
		Status:  e.Status,
		RTId:    e.RTId,
		Year:    e.StartAt.Year(),
	}
	return item, nil
}

func (u *useCase) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return u.repo.UpdateStatus(ctx, id, status)
}
