package auth_router

import (
	"github/com/cl0ky/e-voting-be/internal/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAuthRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	userRepo := auth.NewRepository(db)
	useCase := auth.NewUseCase(db, userRepo)
	controller := auth.NewAuthController(useCase)

	rg.POST("/register", controller.Register)
}
