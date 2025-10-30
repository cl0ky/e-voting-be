package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Register(c *gin.Context)
}

type controller struct {
	useCase UseCase
}

func NewAuthController(useCase UseCase) AuthController {
	return &controller{
		useCase: useCase,
	}
}

func (rc *controller) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := rc.useCase.Register(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}
