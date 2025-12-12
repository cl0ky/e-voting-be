package election

import (
	"github/com/cl0ky/e-voting-be/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ElectionController interface {
	Create(c *gin.Context)
	GetAll(c *gin.Context)
	GetDetail(c *gin.Context)
	UpdateStatus(c *gin.Context)
	FinalizeElection(c *gin.Context)
	VerifyElectionResult(c *gin.Context)
}

type controller struct {
	useCase UseCase
}

func NewElectionController(useCase UseCase) ElectionController {
	return &controller{useCase: useCase}
}

func (rc *controller) Create(c *gin.Context) {
	var req CreateElectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
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
	req.RTId = *user.RTId
	req.CreatedBy = &user.Id
	item, err := rc.useCase.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (rc *controller) GetAll(c *gin.Context) {
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
	items, err := rc.useCase.GetAllByRTId(c.Request.Context(), *user.RTId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, items)
}

func (rc *controller) GetDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	item, err := rc.useCase.GetById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (rc *controller) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req UpdateElectionStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = rc.useCase.UpdateStatus(c.Request.Context(), id, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	item, err := rc.useCase.GetById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (rc *controller) FinalizeElection(c *gin.Context) {
	userVal, exists := c.Get("user")
	if !exists {
		c.JSON(401, gin.H{"error": "user tidak ditemukan di context"})
		return
	}
	user, ok := userVal.(*models.User)
	if !ok || user.Role != "admin" {
		c.JSON(403, gin.H{"error": "hanya admin yang boleh finalize"})
		return
	}
	idStr := c.Param("id")
	log.Printf("[DEBUG] FinalizeElection controller id param: '%s'", idStr)
	electionId, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid election id"})
		return
	}
	log.Printf("[DEBUG] FinalizeElection controller parsed UUID: %s", electionId.String())
	resp, err := rc.useCase.FinalizeElection(c.Request.Context(), electionId, user)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, resp)
}

func (rc *controller) VerifyElectionResult(c *gin.Context) {
	idStr := c.Param("id")
	electionId, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid election id"})
		return
	}
	resp, err := rc.useCase.VerifyElectionResult(c.Request.Context(), electionId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
