package rts_router

import (
	"github/com/cl0ky/e-voting-be/internal/rts"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Init(rg *gin.RouterGroup, db *gorm.DB) {
	rtsRepo := rts.NewRepository(db)
	rtsUseCase := rts.NewUseCase(rtsRepo)
	rtsController := rts.NewRTSController(rtsUseCase)

	rg.GET("/rts", rtsController.GetAllRT)
}
