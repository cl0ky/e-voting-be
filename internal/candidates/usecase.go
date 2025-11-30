package candidates

import (
	"context"
	"errors"
	"github/com/cl0ky/e-voting-be/models"

	"github.com/google/uuid"
)

var ErrRTIdRequired = errors.New("rt_id is required")

type UseCase interface {
	Create(ctx context.Context, userID uuid.UUID, req CreateCandidateRequest) (*CandidateItem, error)
	GetByID(ctx context.Context, id uuid.UUID) (*CandidateItem, error)
	List(ctx context.Context, page, pageSize int) (*CandidateListResponse, error)
	ListByElectionID(ctx context.Context, electionId uuid.UUID) (*CandidateListByElectionResponse, error)
	Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, req UpdateCandidateRequest) (*CandidateItem, error)
	Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
}

type useCase struct {
	repo Repository
}

func NewUseCase(repo Repository) UseCase {
	return &useCase{repo: repo}
}

func (u *useCase) ListByElectionID(ctx context.Context, electionId uuid.UUID) (*CandidateListByElectionResponse, error) {
	items, err := u.repo.ListByElectionID(ctx, electionId)
	if err != nil {
		return nil, err
	}
	resp := &CandidateListByElectionResponse{Items: make([]CandidateItem, 0, len(items)), Total: int64(len(items))}
	for _, c := range items {
		resp.Items = append(resp.Items, toItem(&c))
	}
	return resp, nil
}

func toItem(m *models.Candidate) CandidateItem {
	return CandidateItem{
		Id:         m.Id,
		Name:       m.Name,
		Vision:     m.Vision,
		Mission:    m.Mission,
		PhotoURL:   m.PhotoURL,
		RTId:       m.RTId,
		ElectionId: m.ElectionId,
	}
}

func (u *useCase) Create(ctx context.Context, userID uuid.UUID, req CreateCandidateRequest) (*CandidateItem, error) {
	var rtId uuid.UUID
	if req.RTId != nil {
		rtId = *req.RTId
	} else {
		return nil, ErrRTIdRequired
	}

	if req.ElectionId == "" {
		return nil, errors.New("election_id wajib diisi")
	}
	parsed, err := uuid.Parse(req.ElectionId)
	if err != nil {
		return nil, errors.New("election_id harus UUID")
	}

	electionId := &parsed
	c := &models.Candidate{
		Id:         uuid.New(),
		Name:       req.Name,
		Vision:     req.Vision,
		Mission:    req.Mission,
		PhotoURL:   req.PhotoURL,
		RTId:       rtId,
		ElectionId: electionId,
		BaseModel:  models.BaseModel{CreatedBy: &userID, UpdatedBy: &userID},
	}
	if err := u.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	item := toItem(c)
	return &item, nil
}

func (u *useCase) GetByID(ctx context.Context, id uuid.UUID) (*CandidateItem, error) {
	c, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	item := toItem(c)
	return &item, nil
}

func (u *useCase) List(ctx context.Context, page, pageSize int) (*CandidateListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	items, total, err := u.repo.List(ctx, offset, pageSize)
	if err != nil {
		return nil, err
	}

	resp := &CandidateListResponse{Items: make([]CandidateItem, 0, len(items)), Total: total}
	for _, c := range items {
		resp.Items = append(resp.Items, toItem(&c))
	}
	return resp, nil
}

func (u *useCase) Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, req UpdateCandidateRequest) (*CandidateItem, error) {
	c, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		c.Name = *req.Name
	}
	if req.Vision != nil {
		c.Vision = *req.Vision
	}
	if req.Mission != nil {
		c.Mission = *req.Mission
	}
	if req.PhotoURL != nil {
		c.PhotoURL = *req.PhotoURL
	}
	if req.RTId != nil {
		c.RTId = *req.RTId
	}
	if req.ElectionUUID != nil {
		c.ElectionId = req.ElectionUUID
	}

	c.BaseModel.UpdatedBy = &userID

	if err := u.repo.Update(ctx, c); err != nil {
		return nil, err
	}
	item := toItem(c)
	return &item, nil
}

func (u *useCase) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	// ensure exists first
	if _, err := u.repo.GetByID(ctx, id); err != nil {
		return err
	}
	return u.repo.Delete(ctx, id, &userID)
}
