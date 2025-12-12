package votes_router

import (
	votes "github/com/cl0ky/e-voting-be/internal/votes"

	middleware "github/com/cl0ky/e-voting-be/server/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Init(rg *gin.RouterGroup, db *gorm.DB) {
	repo := votes.NewRepository(db)
	uc := votes.NewUseCase(repo)
	ctrl := votes.NewVotesController(uc)

	group := rg.Group("/votes")
	group.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("Voter"))
	group.POST("/commit", ctrl.CommitVote)
	group.GET("/election-status", ctrl.GetElectionStatus)
	group.GET("/results", ctrl.GetUserVoteResults)
	group.POST("/reveal", ctrl.RevealVote)
}
