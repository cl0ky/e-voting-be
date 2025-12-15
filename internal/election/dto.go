package election

import (
	"time"

	"github.com/google/uuid"
)

type ElectionItem struct {
	Id             uuid.UUID  `json:"election_id"`
	Name           string     `json:"name"`
	Status         string     `json:"status"`
	FinalizeStatus string     `json:"finalize_status"`
	StartAt        time.Time  `json:"start_at"`
	EndAt          time.Time  `json:"end_at"`
	FinalizedAt    *time.Time `json:"finalized_at"`
	RTId           uuid.UUID  `json:"-"`
	Year           int        `json:"-"`
}

type CreateElectionRequest struct {
	Name      string     `json:"name" binding:"required,min=2"`
	StartAt   time.Time  `json:"start_at" binding:"required"`
	EndAt     time.Time  `json:"end_at" binding:"required"`
	Status    string     `json:"status"`
	RTId      uuid.UUID  `json:"rt_id"`
	CreatedBy *uuid.UUID `json:"-"`
}

type UpdateElectionStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type FinalizeElectionRequest struct{}

type FinalizeElectionResult struct {
	CandidateId string `json:"candidate_id"`
	Total       int    `json:"total"`
}

type FinalizeElectionResponse struct {
	Summary          FinalizeElectionSummary `json:"summary"`
	SummaryHash      string                  `json:"summary_hash"`
	BlockchainTxHash string                  `json:"blockchain_tx_hash"`
	Winner           string                  `json:"winner"`
}

type FinalizeElectionSummary struct {
	ElectionId    string                   `json:"election_id"`
	TotalRevealed int                      `json:"total_revealed"`
	Results       []FinalizeElectionResult `json:"results"`
	Winner        string                   `json:"winner"`
	Timestamp     string                   `json:"timestamp"`
}

type VerifyElectionResultResponse struct {
	LocalHash        string `json:"local_hash"`
	DBHash           string `json:"db_hash"`
	BlockchainHash   string `json:"blockchain_hash"`
	BlockchainTxHash string `json:"blockchain_tx_hash,omitempty"`
	Valid            bool   `json:"valid"`
	Message          string `json:"message,omitempty"`
}

type ElectionResultItem struct {
	CandidateId   string `json:"candidate_id"`
	CandidateName string `json:"candidate_name"`
	PhotoURL      string `json:"photo_url"`
	Votes         int    `json:"votes"`
}

type ElectionDetailResponse struct {
	ElectionId     uuid.UUID            `json:"election_id"`
	Name           string               `json:"name"`
	Status         string               `json:"status"`
	FinalizeStatus string               `json:"finalize_status"`
	StartAt        time.Time            `json:"start_at"`
	EndAt          time.Time            `json:"end_at"`
	FinalizedAt    *time.Time           `json:"finalized_at"`
	SummaryHash    string               `json:"summary_hash,omitempty"`
	BlockchainTx   string               `json:"blockchain_tx_hash,omitempty"`
	Results        []ElectionResultItem `json:"results"`
	TotalVotes     int                  `json:"total_votes"`
}
