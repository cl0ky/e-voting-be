package votes

import (
	"context"
	"github/com/cl0ky/e-voting-be/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type controller struct {
	useCase votesUseCase
}

type votesUseCase interface {
	GetElectionStatus(ctx context.Context, rtId uuid.UUID, voterId uuid.UUID) (*ElectionStatusResponse, error)
	CommitVote(ctx context.Context, voterId uuid.UUID, req CommitVoteRequest) error
	RevealVote(ctx context.Context, voterId uuid.UUID, req RevealVoteRequest) error
}

func NewVotesController(useCase votesUseCase) *controller {
	return &controller{useCase: useCase}
}

func (vc *controller) GetElectionStatus(c *gin.Context) {
	userVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user tidak ditemukan di context"})
		return
	}
	user, ok := userVal.(*models.User)
	if !ok || user.RTId == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "RT user tidak ditemukan"})
		return
	}
	result, err := vc.useCase.GetElectionStatus(c.Request.Context(), *user.RTId, user.Id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "active election not found or already voted"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (vc *controller) CommitVote(c *gin.Context) {
	userVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user tidak ditemukan di context"})
		return
	}
	user, ok := userVal.(*models.User)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user tidak valid"})
		return
	}
	var req CommitVoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := vc.useCase.CommitVote(c.Request.Context(), user.Id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Vote committed successfully"})
}

func (vc *controller) RevealVote(c *gin.Context) {
	userVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user tidak ditemukan di context"})
		return
	}
	user, ok := userVal.(*models.User)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user tidak valid"})
		return
	}
	var req RevealVoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := vc.useCase.RevealVote(c.Request.Context(), user.Id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Vote revealed successfully"})
}
