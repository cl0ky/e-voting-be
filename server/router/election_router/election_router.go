package election_router

import (
	"github/com/cl0ky/e-voting-be/internal/election"
	"github/com/cl0ky/e-voting-be/server/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Init(rg *gin.RouterGroup, db *gorm.DB) {
	repo := election.NewRepository(db)
	uc := election.NewUseCase(repo)
	ctrl := election.NewElectionController(uc)

	group := rg.Group("/elections")

	group.Use(middleware.AuthMiddleware())
	group.POST("", ctrl.Create)
	group.GET("", ctrl.GetAll)
	group.GET(":id", ctrl.GetDetail)
	group.PATCH(":id", ctrl.UpdateStatus)
}
