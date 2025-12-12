package votes

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type UseCase interface {
	GetElectionStatus(ctx context.Context, rtId uuid.UUID, voterId uuid.UUID) (*ElectionStatusResponse, error)
	CommitVote(ctx context.Context, voterId uuid.UUID, req CommitVoteRequest) error
	RevealVote(ctx context.Context, voterId uuid.UUID, req RevealVoteRequest) error
	GetUserVoteResults(ctx context.Context, voterId uuid.UUID) (*UserVoteResultsResponse, error)
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

type UserVoteResultItem struct {
	ElectionId     uuid.UUID  `json:"election_id"`
	ElectionName   string     `json:"election_name"`
	ElectionStatus string     `json:"election_status"`
	StartAt        time.Time  `json:"start_at"`
	EndAt          time.Time  `json:"end_at"`
	HasCommitted   bool       `json:"hasCommitted"`
	HasRevealed    bool       `json:"hasRevealed"`
	CandidateId    *uuid.UUID `json:"candidate_id,omitempty"`
	CandidateName  *string    `json:"candidate_name,omitempty"`
	RevealedAt     *time.Time `json:"revealed_at,omitempty"`
}

type UserVoteResultsResponse struct {
	Items []UserVoteResultItem `json:"items"`
	Total int64                `json:"total"`
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
				StartDate: election.StartAt.Format(time.RFC3339),
				EndDate:   election.EndAt.Format(time.RFC3339),
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
				StartDate: committedElection.StartAt.Format(time.RFC3339),
				EndDate:   committedElection.EndAt.Format(time.RFC3339),
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

func (u *useCase) GetUserVoteResults(ctx context.Context, voterId uuid.UUID) (*UserVoteResultsResponse, error) {
	votes, err := u.repo.ListVotesByVoter(ctx, voterId)
	if err != nil {
		return nil, err
	}

	items := make([]UserVoteResultItem, 0, len(votes))
	for _, v := range votes {
		candidateName := (*string)(nil)
		if v.Candidate != nil {
			name := v.Candidate.Name
			candidateName = &name
		}

		item := UserVoteResultItem{
			ElectionId:     v.Election.Id,
			ElectionName:   v.Election.Name,
			ElectionStatus: v.Election.Status,
			StartAt:        v.Election.StartAt,
			EndAt:          v.Election.EndAt,
			HasCommitted:   v.HashVote != "",
			HasRevealed:    v.IsRevealed,
			CandidateId:    v.RevealedCandidateId,
			CandidateName:  candidateName,
			RevealedAt:     v.RevealedAt,
		}
		items = append(items, item)
	}

	return &UserVoteResultsResponse{
		Items: items,
		Total: int64(len(items)),
	}, nil
}

func computeVoteHash(candidateId string, nonce string) string {
	data := candidateId + nonce
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
