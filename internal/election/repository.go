package election

import (
	"context"
	"github/com/cl0ky/e-voting-be/models"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, e *models.Election) error
	GetAllByRTId(ctx context.Context, rtId uuid.UUID) ([]models.Election, error)
	GetById(ctx context.Context, id uuid.UUID) (*models.Election, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	GetElectionByID(ctx context.Context, electionId uuid.UUID) (*models.Election, error)
	GetRevealedVotesByElection(ctx context.Context, electionId uuid.UUID) ([]models.Vote, error)
	SetElectionFinalizeStatusIfPending(ctx context.Context, electionId uuid.UUID, newStatus string) (bool, error)
	SetElectionFinalizeFailed(ctx context.Context, electionId uuid.UUID, reason string) error
	UpdateElectionSummaryAndStatus(ctx context.Context, electionId uuid.UUID, summary datatypes.JSON, summaryHash string) error
	UpdateElectionFinalizeResult(ctx context.Context, electionId uuid.UUID, blockchainTxHash string) error
	GetCandidatesByIDs(ctx context.Context, ids []uuid.UUID) ([]models.Candidate, error)
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

func (r *repository) GetElectionByID(ctx context.Context, electionId uuid.UUID) (*models.Election, error) {
	var election models.Election
	err := r.db.WithContext(ctx).Where("id = ?", electionId).First(&election).Error
	if err != nil {
		return nil, err
	}
	return &election, nil
}

func (r *repository) GetRevealedVotesByElection(ctx context.Context, electionId uuid.UUID) ([]models.Vote, error) {
	var votes []models.Vote
	err := r.db.WithContext(ctx).Where("election_id = ? AND is_revealed = true", electionId).Find(&votes).Error
	if err != nil {
		return nil, err
	}
	return votes, nil
}

func (r *repository) SetElectionFinalizeStatusIfPending(ctx context.Context, electionId uuid.UUID, newStatus string) (bool, error) {
	log.Printf("[FinalizeStatus] Attempting to set finalize_status to '%s' for electionId: %s", newStatus, electionId.String())
	res := r.db.WithContext(ctx).
		Model(&models.Election{}).
		Where("id = ? AND (finalize_status = 'pending' OR finalize_status = 'failed' OR finalize_status IS NULL)", electionId).
		Update("finalize_status", newStatus)
	// DEBUG LOG
	if res.Error != nil {
		log.Printf("[FinalizeStatus] Update error: %v", res.Error)
		return false, res.Error
	}
	log.Printf("[FinalizeStatus] Rows affected: %d for electionId: %s (status awal NULL/pending/failed)", res.RowsAffected, electionId.String())
	return res.RowsAffected > 0, nil
}

func (r *repository) SetElectionFinalizeFailed(ctx context.Context, electionId uuid.UUID, reason string) error {
	return r.db.WithContext(ctx).Model(&models.Election{}).
		Where("id = ?", electionId).
		Updates(map[string]interface{}{
			"finalize_status": "failed",
			"finalize_error":  reason,
		}).Error
}

func (r *repository) UpdateElectionSummaryAndStatus(ctx context.Context, electionId uuid.UUID, summary datatypes.JSON, summaryHash string) error {
	return r.db.WithContext(ctx).Model(&models.Election{}).
		Where("id = ?", electionId).
		Updates(map[string]interface{}{
			"summary_json":    summary,
			"summary_hash":    summaryHash,
			"finalize_status": "finalizing",
		}).Error
}

func (r *repository) UpdateElectionFinalizeResult(ctx context.Context, electionId uuid.UUID, blockchainTxHash string) error {
	now := time.Now().UTC()
	return r.db.WithContext(ctx).Model(&models.Election{}).
		Where("id = ?", electionId).
		Updates(map[string]interface{}{
			"blockchain_tx_hash": blockchainTxHash,
			"finalized_at":       now,
			"finalize_status":    "success",
			"status":             "finalized",
		}).Error
}

func (r *repository) GetCandidatesByIDs(ctx context.Context, ids []uuid.UUID) ([]models.Candidate, error) {
	var candidates []models.Candidate
	if len(ids) == 0 {
		return candidates, nil
	}
	err := r.db.WithContext(ctx).
		Where("id IN ?", ids).
		Find(&candidates).Error
	if err != nil {
		return nil, err
	}
	return candidates, nil
}
