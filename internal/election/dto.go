package election

import (
	"time"

	"github.com/google/uuid"
)

type ElectionItem struct {
	Id      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at"`
	Status  string    `json:"status"`
	RTId    uuid.UUID `json:"rt_id"`
	Year    int       `json:"year"`
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
