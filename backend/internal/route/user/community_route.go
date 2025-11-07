package userroute

import (
	"github.com/giakiet05/lkforum/internal/controller"
	"github.com/giakiet05/lkforum/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterCommunityRoutes(rg *gin.RouterGroup, c *controller.CommunityController) {
	communities := rg.Group("/communities")

	communities.GET("filter", c.GetCommunitiesFilter)
	communities.GET(":community_id", c.GetCommunityByID)

	// Protected routes (require authentication)
	communities.Use(middleware.RequireAuth())
	{
		communities.POST("", c.CreateCommunity)
		communities.PUT("", c.UpdateCommunity)
		communities.PUT("/add_moderator", c.AddModerator)
		communities.PUT("/remove_moderator", c.RemoveModerator)
		communities.DELETE("/:community_id", c.DeleteCommunityByID)
	}
}
