package election

import (
	"bytes"
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
	GetDetail(ctx context.Context, id uuid.UUID) (*ElectionDetailResponse, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	FinalizeElection(ctx context.Context, electionId uuid.UUID, user *models.User) (*FinalizeElectionResponse, error)
	VerifyElectionResult(ctx context.Context, electionId uuid.UUID) (*VerifyElectionResultResponse, error)
	GetAdminDashboard(ctx context.Context, rtId uuid.UUID) (*AdminDashboardResponse, error)
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
		Id:             e.Id,
		Name:           e.Name,
		Status:         e.Status,
		FinalizeStatus: e.FinalizeStatus,
		StartAt:        e.StartAt,
		EndAt:          e.EndAt,
		FinalizedAt:    e.FinalizedAt,
		RTId:           e.RTId,
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
			Id:             e.Id,
			Name:           e.Name,
			Status:         e.Status,
			FinalizeStatus: e.FinalizeStatus,
			StartAt:        e.StartAt,
			EndAt:          e.EndAt,
			FinalizedAt:    e.FinalizedAt,
			RTId:           e.RTId,
			Year:           e.StartAt.Year(),
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
		Id:             e.Id,
		Name:           e.Name,
		Status:         e.Status,
		FinalizeStatus: e.FinalizeStatus,
		StartAt:        e.StartAt,
		EndAt:          e.EndAt,
		FinalizedAt:    e.FinalizedAt,
		RTId:           e.RTId,
		Year:           e.StartAt.Year(),
	}
	return item, nil
}

func (u *useCase) GetDetail(ctx context.Context, id uuid.UUID) (*ElectionDetailResponse, error) {
	e, err := u.repo.GetElectionByID(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := &ElectionDetailResponse{
		ElectionId:     e.Id,
		Name:           e.Name,
		Status:         e.Status,
		FinalizeStatus: e.FinalizeStatus,
		StartAt:        e.StartAt,
		EndAt:          e.EndAt,
		FinalizedAt:    e.FinalizedAt,
		SummaryHash:    e.SummaryHash,
		BlockchainTx:   e.BlockchainTxHash,
		Results:        []ElectionResultItem{},
		TotalVotes:     0,
	}

	summaryJSON := e.SummaryJSON
	if summaryJSON != nil && len(summaryJSON) > 0 {
		var summary FinalizeElectionSummary
		if err := json.Unmarshal([]byte(summaryJSON), &summary); err != nil {
			log.Printf("[GetDetail] failed to unmarshal summary json: %v", err)
		} else {
			candidateIDs := make([]uuid.UUID, 0, len(summary.Results))
			for _, r := range summary.Results {
				cid, err := uuid.Parse(r.CandidateId)
				if err != nil {
					continue
				}
				candidateIDs = append(candidateIDs, cid)
			}
			candidates, err := u.repo.GetCandidatesByIDs(ctx, candidateIDs)
			if err != nil {
				log.Printf("[GetDetail] failed to load candidates: %v", err)
			} else {
				candMap := make(map[uuid.UUID]models.Candidate)
				for _, c := range candidates {
					candMap[c.Id] = c
				}

				for _, r := range summary.Results {
					cid, err := uuid.Parse(r.CandidateId)
					if err != nil {
						continue
					}
					cand, ok := candMap[cid]
					item := ElectionResultItem{
						CandidateId: r.CandidateId,
						Votes:       r.Total,
					}
					if ok {
						item.CandidateName = cand.Name
						item.PhotoURL = cand.PhotoURL
					}
					resp.Results = append(resp.Results, item)
					resp.TotalVotes += r.Total
				}
			}
		}
	}

	return resp, nil
}

func (u *useCase) GetAdminDashboard(ctx context.Context, rtId uuid.UUID) (*AdminDashboardResponse, error) {
	total, active, unfinalized, err := u.repo.GetDashboardStatsByRT(ctx, rtId)
	if err != nil {
		return nil, err
	}

	recentModels, err := u.repo.GetRecentElectionsByRT(ctx, rtId, 3)
	if err != nil {
		return nil, err
	}

	recent := make([]AdminDashboardRecentElection, 0, len(recentModels))
	for _, e := range recentModels {
		recent = append(recent, AdminDashboardRecentElection{
			ElectionId:     e.Id,
			Name:           e.Name,
			Status:         e.Status,
			FinalizeStatus: e.FinalizeStatus,
		})
	}

	resp := &AdminDashboardResponse{
		Summary: AdminDashboardSummary{
			TotalElections:       total,
			ActiveElections:      active,
			UnfinalizedElections: unfinalized,
		},
		RecentElections: recent,
	}
	return resp, nil
}

func (u *useCase) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return u.repo.UpdateStatus(ctx, id, status)
}

type BlockchainService interface {
	StoreResultHash(electionId string, hash string) (string, error)
	GetHash(electionId string) (string, error)
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

	if user.Role != "Admin" || user.RTId == nil || election.RTId != *user.RTId {
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
		if total > maxVote || (total == maxVote && (winner == "" || cid < winner)) {
			winner = cid
			maxVote = total
		}
	}
	if maxVote == -1 {
		winner = ""
		maxVote = 0
	}

	timestamp := time.Now().UTC().Format(time.RFC3339)
	var finalizeResults []FinalizeElectionResult
	for _, r := range results {
		finalizeResults = append(finalizeResults, FinalizeElectionResult{
			CandidateId: r.CandidateId,
			Total:       r.Votes,
		})
	}

	summary := FinalizeElectionSummary{
		ElectionId:    electionId.String(),
		TotalRevealed: len(votes),
		Results:       finalizeResults,
		Winner:        winner,
		Timestamp:     timestamp,
	}

	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)

	if err := enc.Encode(summary); err != nil {
		_ = u.repo.SetElectionFinalizeFailed(ctx, electionId, "marshal summary error: "+err.Error())
		return nil, fmt.Errorf("failed to marshal summary: %w", err)
	}

	summaryJSONBytes := buf.Bytes()
	if len(summaryJSONBytes) > 0 && summaryJSONBytes[len(summaryJSONBytes)-1] == '\n' {
		summaryJSONBytes = summaryJSONBytes[:len(summaryJSONBytes)-1]
	}

	hashBytes := sha256.Sum256(summaryJSONBytes)
	summaryHash := "0x" + hex.EncodeToString(hashBytes[:])

	if err := u.repo.UpdateElectionSummaryAndStatus(ctx, electionId, datatypes.JSON(summaryJSONBytes), summaryHash); err != nil {
		_ = u.repo.SetElectionFinalizeFailed(ctx, electionId, "db update error: "+err.Error())
		return nil, fmt.Errorf("failed to update election summary: %w", err)
	}

	txHash, err := blockchainService.StoreResultHash(electionId.String(), summaryHash)

	if err != nil {
		_ = u.repo.SetElectionFinalizeFailed(ctx, electionId, "blockchain error: "+err.Error())
		return nil, fmt.Errorf("blockchain error: %w", err)
	}
	if err := u.repo.UpdateElectionFinalizeResult(ctx, electionId, txHash); err != nil {
		return nil, fmt.Errorf("failed to finalize election: %w", err)
	}
	responseSummary := summary
	return &FinalizeElectionResponse{
		Summary:          responseSummary,
		SummaryHash:      summaryHash,
		BlockchainTxHash: txHash,
		Winner:           winner,
	}, nil
}

func (u *useCase) VerifyElectionResult(ctx context.Context, electionId uuid.UUID) (*VerifyElectionResultResponse, error) {
	election, err := u.repo.GetElectionByID(ctx, electionId)
	if err != nil {
		return nil, fmt.Errorf("election not found")
	}

	dbHash := election.SummaryHash
	blockchainTxHash := election.BlockchainTxHash
	summaryJSON := election.SummaryJSON

	var summaryBytes []byte
	if summaryJSON != nil && len(summaryJSON) > 0 {
		var summaryStruct FinalizeElectionSummary
		if err := json.Unmarshal([]byte(summaryJSON), &summaryStruct); err != nil {
			log.Printf("[VERIFY] warning: unable to unmarshal summary json: %v", err)
			summaryBytes = []byte(summaryJSON)
		} else {
			buf := &bytes.Buffer{}
			enc := json.NewEncoder(buf)
			enc.SetEscapeHTML(false)
			if err := enc.Encode(summaryStruct); err != nil {
				log.Printf("[VERIFY] warning: unable to encode canonical summary: %v", err)
				summaryBytes = []byte(summaryJSON)
			} else {
				b := buf.Bytes()
				if len(b) > 0 && b[len(b)-1] == '\n' {
					b = b[:len(b)-1]
				}
				summaryBytes = b
			}
		}
	} else {
		summaryBytes = []byte{}
	}

	hashBytes := sha256.Sum256(summaryBytes)
	localHash := "0x" + hex.EncodeToString(hashBytes[:])

	log.Printf("[VERIFY] canonical_summary_json: %s", string(summaryBytes))
	log.Printf("[VERIFY] localHash: %s", localHash)
	log.Printf("[VERIFY] dbHash: %s", dbHash)

	blockchainHash := ""
	if blockchainService != nil {
		blockchainHash, err = blockchainService.GetHash(electionId.String())
		if err != nil {
			blockchainHash = ""
		}
	}
	log.Printf("[VERIFY] blockchainHash: %s", blockchainHash)

	valid := (localHash == dbHash) && (dbHash == blockchainHash)

	if valid {
		return &VerifyElectionResultResponse{
			LocalHash:        localHash,
			DBHash:           dbHash,
			BlockchainHash:   blockchainHash,
			BlockchainTxHash: blockchainTxHash,
			Valid:            true,
		}, nil
	} else {
		return &VerifyElectionResultResponse{
			LocalHash:      localHash,
			DBHash:         dbHash,
			BlockchainHash: blockchainHash,
			Valid:          false,
			Message:        "Hash mismatch: database or blockchain data has been modified.",
		}, nil
	}
}
