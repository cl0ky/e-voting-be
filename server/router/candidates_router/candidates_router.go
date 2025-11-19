package candidates_router

import (
	"github/com/cl0ky/e-voting-be/internal/candidates"
	"github/com/cl0ky/e-voting-be/server/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Init(rg *gin.RouterGroup, db *gorm.DB) {
	repo := candidates.NewRepository(db)
	uc := candidates.NewUseCase(repo)
	ctrl := candidates.NewController(uc)

	group := rg.Group("/candidates")
	group.GET("", ctrl.List)
	group.GET("/:id", ctrl.GetByID)

	// Protected write operations
	group.Use(middleware.AuthMiddleware())
	group.POST("", ctrl.Create)
	group.PUT("/:id", ctrl.Update)
	group.DELETE("/:id", ctrl.Delete)
}
