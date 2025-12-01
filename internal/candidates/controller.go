package candidates

import (
	"fmt"
	"net/http"

	"github/com/cl0ky/e-voting-be/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CandidateController interface {
	ListByElectionID(c *gin.Context)
	List(c *gin.Context)
	GetByID(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type candidateController struct {
	useCase UseCase
}

func NewCandidateController(useCase UseCase) CandidateController {
	return &candidateController{useCase: useCase}
}

func (cc *candidateController) ListByElectionID(c *gin.Context) {
	electionIdStr := c.Param("election_id")
	electionId, err := uuid.Parse(electionIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid election_id"})
		return
	}
	resp, err := cc.useCase.ListByElectionID(c.Request.Context(), electionId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (cc *candidateController) List(c *gin.Context) {
	page := 1
	pageSize := 20
	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if ps := c.Query("page_size"); ps != "" {
		fmt.Sscanf(ps, "%d", &pageSize)
	}
	resp, err := cc.useCase.List(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (cc *candidateController) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	item, err := cc.useCase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (cc *candidateController) Create(c *gin.Context) {
	var req CreateCandidateRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user tidak ditemukan di context"})
		return
	}
	user, ok := userVal.(*models.User)
	if !ok || user.Id == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user tidak valid"})
		return
	}
	if req.RTId == nil {
		req.RTId = user.RTId
	}
	item, err := cc.useCase.Create(c.Request.Context(), user.Id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (cc *candidateController) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req UpdateCandidateRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user tidak ditemukan di context"})
		return
	}
	user, ok := userVal.(*models.User)
	if !ok || user.Id == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user tidak valid"})
		return
	}
	item, err := cc.useCase.Update(c.Request.Context(), user.Id, id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (cc *candidateController) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	userVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user tidak ditemukan di context"})
		return
	}
	user, ok := userVal.(*models.User)
	if !ok || user.Id == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user tidak valid"})
		return
	}
	err = cc.useCase.Delete(c.Request.Context(), user.Id, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
