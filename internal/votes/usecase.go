package votes

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
)

type UseCase interface {
	GetElectionStatus(ctx context.Context, rtId uuid.UUID, voterId uuid.UUID) (*ElectionStatusResponse, error)
	CommitVote(ctx context.Context, voterId uuid.UUID, req CommitVoteRequest) error
	RevealVote(ctx context.Context, voterId uuid.UUID, req RevealVoteRequest) error
}

type useCase struct {
	repo Repository
}

func NewUseCase(repo Repository) UseCase {
	return &useCase{repo: repo}
}

type ElectionStatusResponse struct {
	Election     *ElectionStatusItem `json:"election"`
	HasCommitted bool                `json:"hasCommitted"`
	HasRevealed  bool                `json:"hasRevealed"`
}

type ElectionStatusItem struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	StartDate string    `json:"start_date"`
	EndDate   string    `json:"end_date"`
}

func (u *useCase) GetElectionStatus(ctx context.Context, rtId uuid.UUID, voterId uuid.UUID) (*ElectionStatusResponse, error) {
	election, err := u.repo.GetActiveElectionForVoter(ctx, rtId, voterId)
	if err == nil && election != nil {
		hasCommitted, hasRevealed, _ := u.repo.GetVoteStatus(ctx, election.Id, voterId)
		resp := &ElectionStatusResponse{
			Election: &ElectionStatusItem{
				Id:        election.Id,
				Name:      election.Name,
				Status:    election.Status,
				StartDate: election.StartAt.Format("2006-01-02"),
				EndDate:   election.EndAt.Format("2006-01-02"),
			},
			HasCommitted: hasCommitted,
			HasRevealed:  hasRevealed,
		}
		return resp, nil
	}
	committedElection, err := u.repo.GetCommittedActiveElectionForVoter(ctx, rtId, voterId)
	if err == nil && committedElection != nil {
		hasCommitted, hasRevealed, _ := u.repo.GetVoteStatus(ctx, committedElection.Id, voterId)
		resp := &ElectionStatusResponse{
			Election: &ElectionStatusItem{
				Id:        committedElection.Id,
				Name:      committedElection.Name,
				Status:    committedElection.Status,
				StartDate: committedElection.StartAt.Format("2006-01-02"),
				EndDate:   committedElection.EndAt.Format("2006-01-02"),
			},
			HasCommitted: hasCommitted,
			HasRevealed:  hasRevealed,
		}
		return resp, nil
	}

	return &ElectionStatusResponse{
		Election:     nil,
		HasCommitted: false,
		HasRevealed:  false,
	}, nil
}

func (u *useCase) CommitVote(ctx context.Context, voterId uuid.UUID, req CommitVoteRequest) error {
	electionId, err := uuid.Parse(req.ElectionId)
	if err != nil {
		return err
	}
	return u.repo.CreateVoteCommit(ctx, voterId, electionId, req.HashVote)
}

func (u *useCase) RevealVote(ctx context.Context, voterId uuid.UUID, req RevealVoteRequest) error {
	electionId, err := uuid.Parse(req.ElectionId)
	if err != nil {
		return err
	}
	candidateId, err := uuid.Parse(req.CandidateId)
	if err != nil {
		return err
	}
	vote, err := u.repo.GetVoteByVoterAndElection(ctx, electionId, voterId)
	if err != nil {
		return err
	}
	if vote.IsRevealed {
		return fmt.Errorf("vote has already been revealed")
	}
	computedHash := computeVoteHash(req.CandidateId, req.Nonce)
	if vote.HashVote != computedHash {
		return fmt.Errorf("hash verification failed: invalid candidate_id or nonce")
	}
	return u.repo.RevealVote(ctx, voterId, electionId, candidateId, req.Nonce)
}

func computeVoteHash(candidateId string, nonce string) string {
	data := candidateId + nonce
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
