package auth_router

import (
	"github/com/cl0ky/e-voting-be/internal/auth"
	"github/com/cl0ky/e-voting-be/server/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Init(rg *gin.RouterGroup, db *gorm.DB) {
	userRepo := auth.NewRepository(db)
	useCase := auth.NewUseCase(db, userRepo)
	controller := auth.NewAuthController(useCase)

	rg.POST("/register", controller.Register)
	rg.POST("/login", controller.Login)

	authGroup := rg.Group("")
	authGroup.Use(middleware.AuthMiddleware())

	authGroup.GET("/profile", controller.GetProfile)
	authGroup.POST("/logout", controller.Logout)
}
