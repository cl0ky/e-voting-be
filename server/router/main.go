package router

import (
	"github/com/cl0ky/e-voting-be/server/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SetupRoutesConfig struct {
	Router *gin.Engine
	DB     *gorm.DB
}

func SetupRoutes(c SetupRoutesConfig) {
	c.Router.Use(middleware.ErrorHandler())
	c.Router.Use(middleware.ReqLog())
	c.Router.Use(middleware.CORSMiddleware())
	c.Router.Use(middleware.RecoveryMiddleware())

	apiV1 := c.Router.Group("/v1")

	apiV1.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Server healthy..."})
	})

	apiV1.GET("/error", func(c *gin.Context) {
		panic("error test panic")
	})
}
