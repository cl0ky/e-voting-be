package votes

import (
	"context"
	"github/com/cl0ky/e-voting-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	GetActiveElectionForVoter(ctx context.Context, rtId uuid.UUID, voterId uuid.UUID) (*models.Election, error)
	GetCommittedActiveElectionForVoter(ctx context.Context, rtId uuid.UUID, voterId uuid.UUID) (*models.Election, error)
	GetVoteStatus(ctx context.Context, electionId uuid.UUID, voterId uuid.UUID) (hasCommitted bool, hasRevealed bool, err error)
	CreateVoteCommit(ctx context.Context, voterId uuid.UUID, electionId uuid.UUID, hashVote string) error
	RevealVote(ctx context.Context, voterId uuid.UUID, electionId uuid.UUID, candidateId uuid.UUID, nonce string) error
	GetVoteByVoterAndElection(ctx context.Context, electionId uuid.UUID, voterId uuid.UUID) (*models.Vote, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetActiveElectionForVoter(ctx context.Context, rtId uuid.UUID, voterId uuid.UUID) (*models.Election, error) {
	var election models.Election
	err := r.db.WithContext(ctx).
		Where("rt_id = ? AND status = ?", rtId, "active").
		Where("NOT EXISTS (SELECT 1 FROM votes WHERE votes.election_id = elections.id AND votes.voter_id = ?)", voterId).
		First(&election).Error
	if err != nil {
		return nil, err
	}
	return &election, nil
}

func (r *repository) GetCommittedActiveElectionForVoter(ctx context.Context, rtId uuid.UUID, voterId uuid.UUID) (*models.Election, error) {
	var election models.Election
	err := r.db.WithContext(ctx).
		Where("rt_id = ?", rtId).
		Joins("JOIN votes ON votes.election_id = elections.id AND votes.voter_id = ?", voterId).
		Order("elections.start_at DESC").
		First(&election).Error
	if err != nil {
		return nil, err
	}
	return &election, nil
}

func (r *repository) GetVoteStatus(ctx context.Context, electionId uuid.UUID, voterId uuid.UUID) (bool, bool, error) {
	var vote models.Vote
	err := r.db.WithContext(ctx).Where("election_id = ? AND voter_id = ?", electionId, voterId).First(&vote).Error
	if err != nil {
		return false, false, err
	}
	hasCommitted := vote.HashVote != ""
	hasRevealed := vote.IsRevealed
	return hasCommitted, hasRevealed, nil
}

func (r *repository) CreateVoteCommit(ctx context.Context, voterId uuid.UUID, electionId uuid.UUID, hashVote string) error {
	vote := models.Vote{
		Id:         uuid.New(),
		VoterId:    voterId,
		ElectionId: electionId,
		HashVote:   hashVote,
		IsRevealed: false,
		BaseModel:  models.BaseModel{CreatedBy: &voterId},
	}
	return r.db.WithContext(ctx).Create(&vote).Error
}

func (r *repository) RevealVote(ctx context.Context, voterId uuid.UUID, electionId uuid.UUID, candidateId uuid.UUID, nonce string) error {
	now := gorm.Expr("NOW()")
	return r.db.WithContext(ctx).Model(&models.Vote{}).
		Where("voter_id = ? AND election_id = ?", voterId, electionId).
		Updates(map[string]interface{}{
			"is_revealed":           true,
			"revealed_candidate_id": candidateId,
			"revealed_at":           now,
			"updated_by":            voterId,
		}).Error
}

func (r *repository) GetVoteByVoterAndElection(ctx context.Context, electionId uuid.UUID, voterId uuid.UUID) (*models.Vote, error) {
	var vote models.Vote
	err := r.db.WithContext(ctx).Where("election_id = ? AND voter_id = ?", electionId, voterId).First(&vote).Error
	if err != nil {
		return nil, err
	}
	return &vote, nil
}
