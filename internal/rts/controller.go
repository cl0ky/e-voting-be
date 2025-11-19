package rts

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RTSController interface {
	GetAllRT(c *gin.Context)
}

type controller struct {
	useCase UseCase
}

func NewRTSController(useCase UseCase) RTSController {
	return &controller{useCase: useCase}
}

func (rc *controller) GetAllRT(c *gin.Context) {
	rts, err := rc.useCase.GetAllRT()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data RT"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"rts": rts})
}
