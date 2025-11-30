package candidates

import (
	"mime/multipart"

	"github.com/google/uuid"
)

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
	Name       string                `form:"name" binding:"required,min=2"`
	Vision     string                `form:"vision"`
	Mission    string                `form:"mission"`
	Photo      *multipart.FileHeader `form:"photo"`
	PhotoURL   string                `form:"-"`
	RTId       *uuid.UUID            `form:"rt_id"`
	ElectionId string                `form:"election_id"`
}

type UpdateCandidateRequest struct {
	ElectionUUID *uuid.UUID            `json:"-"`
	Name         *string               `form:"name"`
	Vision       *string               `form:"vision"`
	Mission      *string               `form:"mission"`
	Photo        *multipart.FileHeader `form:"photo"`
	PhotoURL     *string               `form:"-"`
	RTId         *uuid.UUID            `form:"rt_id"`
	ElectionId   string                `form:"election_id"`
}

type CandidateListResponse struct {
	Items []CandidateItem `json:"items"`
	Total int64           `json:"total"`
}

type CandidateListByElectionResponse struct {
	Items []CandidateItem `json:"items"`
	Total int64           `json:"total"`
}
