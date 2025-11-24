package candidates

import (
	"github/com/cl0ky/e-voting-be/models"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Controller interface {
	Create(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
	ListByElectionID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type controller struct{ uc UseCase }

func NewController(uc UseCase) Controller { return &controller{uc: uc} }

func (ctl *controller) ListByElectionID(c *gin.Context) {
	electionId, err := uuid.Parse(c.Param("election_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid election_id"})
		return
	}
	resp, err := ctl.uc.ListByElectionID(c.Request.Context(), electionId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func userIDFromCtx(c *gin.Context) (uuid.UUID, bool) {
	if v, ok := c.Get("user_id"); ok {
		if s, ok := v.(string); ok {
			if id, err := uuid.Parse(s); err == nil {
				return id, true
			}
		}
	}
	return uuid.UUID{}, false
}

func (ctl *controller) Create(c *gin.Context) {
	c.Request.ParseMultipartForm(32 << 20)
	for key, values := range c.Request.PostForm {
		for _, v := range values {
			println("[DEBUG] POST ", key, "=", v)
		}
	}
	for key, files := range c.Request.MultipartForm.File {
		for _, f := range files {
			println("[DEBUG] FILE ", key, "=", f.Filename)
		}
	}
	var req CreateCandidateRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	user, ok := userVal.(*models.User)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User context error"})
		return
	}
	if req.RTId == nil {
		if user.RTId == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "RT user tidak ditemukan"})
			return
		}
		req.RTId = user.RTId
	}

	if req.Photo != nil {
		filename := uuid.New().String() + filepath.Ext(req.Photo.Filename)
		savePath := filepath.Join("uploads", "candidates", filename)
		if err := c.SaveUploadedFile(req.Photo, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal upload foto"})
			return
		}
		req.PhotoURL = "/uploads/candidates/" + filename
	}

	item, err := ctl.uc.Create(c.Request.Context(), user.Id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (ctl *controller) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	item, err := ctl.uc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (ctl *controller) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	resp, err := ctl.uc.List(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (ctl *controller) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req UpdateCandidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid, ok := userIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	item, err := ctl.uc.Update(c.Request.Context(), uid, id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (ctl *controller) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	uid, ok := userIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if err := ctl.uc.Delete(c.Request.Context(), uid, id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
