package election

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github/com/cl0ky/e-voting-be/models"
	"log"
	"sort"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type UseCase interface {
	Create(ctx context.Context, req CreateElectionRequest) (*ElectionItem, error)
	GetAllByRTId(ctx context.Context, rtId uuid.UUID) ([]ElectionItem, error)
	GetById(ctx context.Context, id uuid.UUID) (*ElectionItem, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	FinalizeElection(ctx context.Context, electionId uuid.UUID, user *models.User) (*FinalizeElectionResponse, error)
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

type BlockchainService interface {
	StoreResultHash(hash string) (string, error)
}

var blockchainService BlockchainService

func SetBlockchainService(svc BlockchainService) {
	blockchainService = svc
}

func (u *useCase) FinalizeElection(ctx context.Context, electionId uuid.UUID, user *models.User) (*FinalizeElectionResponse, error) {
	if blockchainService == nil {
		return nil, fmt.Errorf("blockchain service not initialized")
	}
	log.Printf("[DEBUG] FinalizeElection usecase electionId: %s", electionId.String())
	election, err := u.repo.GetElectionByID(ctx, electionId)
	if err != nil {
		return nil, fmt.Errorf("election not found")
	}

	if user.Role != "admin" || user.RTId == nil || election.RTId != *user.RTId {
		return nil, fmt.Errorf("forbidden: not election admin")
	}

	if election.FinalizeStatus == "success" {
		return &FinalizeElectionResponse{
			SummaryHash:      election.SummaryHash,
			BlockchainTxHash: election.BlockchainTxHash,
		}, nil
	}
	if election.Status != "ended" {
		return nil, fmt.Errorf("election status must be 'ended'")
	}

	setOk, err := u.repo.SetElectionFinalizeStatusIfPending(ctx, electionId, "finalizing")
	if err != nil {
		return nil, fmt.Errorf("failed to lock election for finalizing: %w", err)
	}
	if !setOk {
		return nil, fmt.Errorf("finalize already in progress or completed")
	}

	votes, err := u.repo.GetRevealedVotesByElection(ctx, electionId)
	if err != nil {
		_ = u.repo.SetElectionFinalizeFailed(ctx, electionId, "failed to read votes: "+err.Error())
		return nil, fmt.Errorf("failed to get revealed votes: %w", err)
	}

	counts := map[string]int{}
	for _, v := range votes {
		if v.RevealedCandidateId != nil {
			cid := v.RevealedCandidateId.String()
			counts[cid]++
		}
	}

	type resultItem struct {
		CandidateId string `json:"candidate_id"`
		Votes       int    `json:"votes"`
	}
	var candidateIDs []string
	for cid := range counts {
		candidateIDs = append(candidateIDs, cid)
	}
	sort.Strings(candidateIDs)

	var results []resultItem
	winner := ""
	maxVote := -1
	for _, cid := range candidateIDs {
		total := counts[cid]
		results = append(results, resultItem{CandidateId: cid, Votes: total})
		if total > maxVote || (total == maxVote && cid < winner) {
			winner = cid
			maxVote = total
		}
	}
	if maxVote == -1 {
		winner = ""
		maxVote = 0
	}

	summary := map[string]interface{}{
		"election_id":    electionId.String(),
		"total_revealed": len(votes),
		"results":        results,
		"winner":         winner,
		"timestamp":      time.Now().UTC().Format(time.RFC3339),
	}

	summaryJSON, err := json.Marshal(summary)
	if err != nil {
		_ = u.repo.SetElectionFinalizeFailed(ctx, electionId, "marshal summary error: "+err.Error())
		return nil, fmt.Errorf("failed to marshal summary: %w", err)
	}

	hashBytes := sha256.Sum256(summaryJSON)
	summaryHash := "0x" + hex.EncodeToString(hashBytes[:])

	if err := u.repo.UpdateElectionSummaryAndStatus(ctx, electionId, datatypes.JSON(summaryJSON), summaryHash); err != nil {
		_ = u.repo.SetElectionFinalizeFailed(ctx, electionId, "db update error: "+err.Error())
		return nil, fmt.Errorf("failed to update election summary: %w", err)
	}

	txHash, err := blockchainService.StoreResultHash(summaryHash)
	if err != nil {
		_ = u.repo.SetElectionFinalizeFailed(ctx, electionId, "blockchain error: "+err.Error())
		return nil, fmt.Errorf("blockchain error: %w", err)
	}

	if err := u.repo.UpdateElectionFinalizeResult(ctx, electionId, txHash); err != nil {
		return nil, fmt.Errorf("failed to finalize el ction: %w", err)
	}

	var finalizeResults []FinalizeElectionResult
	for _, r := range results {
		finalizeResults = append(finalizeResults, FinalizeElectionResult{
			CandidateId: r.CandidateId,
			Total:       r.Votes,
		})
	}
	responseSummary := FinalizeElectionSummary{
		ElectionId:    electionId.String(),
		TotalRevealed: len(votes),
		Results:       finalizeResults,
		Winner:        winner,
		Timestamp:     summary["timestamp"].(string),
	}
	return &FinalizeElectionResponse{
		Summary:          responseSummary,
		SummaryHash:      summaryHash,
		BlockchainTxHash: txHash,
		Winner:           winner,
	}, nil
}
