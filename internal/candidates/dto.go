package candidates

import "github.com/google/uuid"

// DTOs for Candidate entity

type CandidateItem struct {
	Id         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	Vision     string     `json:"vision"`
	Mission    string     `json:"mission"`
	PhotoURL   string     `json:"photo_url"`
	RTId       uuid.UUID  `json:"rt_id"`
	ElectionId *uuid.UUID `json:"election_id,omitempty"`
}

type CreateCandidateRequest struct {
	Name       string     `json:"name" binding:"required,min=2"`
	Vision     string     `json:"vision"`
	Mission    string     `json:"mission"`
	PhotoURL   string     `json:"photo_url"`
	RTId       uuid.UUID  `json:"rt_id" binding:"required"`
	ElectionId *uuid.UUID `json:"election_id"`
}

type UpdateCandidateRequest struct {
	Name       *string    `json:"name"`
	Vision     *string    `json:"vision"`
	Mission    *string    `json:"mission"`
	PhotoURL   *string    `json:"photo_url"`
	RTId       *uuid.UUID `json:"rt_id"`
	ElectionId *uuid.UUID `json:"election_id"`
}

type CandidateListResponse struct {
	Items []CandidateItem `json:"items"`
	Total int64           `json:"total"`
}
