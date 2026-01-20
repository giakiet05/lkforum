package route

import (
	"github.com/giakiet05/lkforum/internal/controller"
	"github.com/giakiet05/lkforum/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterVoteRoutes(rg *gin.RouterGroup, vc *controller.VoteController) {
	votes := rg.Group("/votes")

	// All vote endpoints require authentication
	votes.Use(middleware.RequireAuth())

	{
		// Vote on a target (post/comment)
		votes.POST("/:target_type/:target_id", vc.VoteOnTarget)

		// Get user's vote on a target
		votes.GET("/:target_type/:target_id", vc.GetUserVote)

		// Get user's votes for multiple targets (bulk operation)
		votes.POST("/bulk/:target_type", vc.GetBulkUserVotes)
	}
}
