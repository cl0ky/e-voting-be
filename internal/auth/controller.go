package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
	GetProfile(c *gin.Context)
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

func (rc *controller) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := rc.useCase.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("token", resp.Token, 60*60*24, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"role":    resp.Role,
		"message": resp.Message,
	})
}

func (rc *controller) Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logout berhasil"})
}

func (rc *controller) GetProfile(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id tidak ditemukan"})
		return
	}

	user, err := rc.useCase.GetProfile(userId.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    user.Id,
		"name":  user.Name,
		"email": user.Email,
		"nik":   user.NIK,
		"role":  user.Role,
	})
}
