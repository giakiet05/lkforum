package adminroute

import (
	"github.com/giakiet05/lkforum/internal/controller"
	"github.com/giakiet05/lkforum/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterAdminCommunityRoutes(rg *gin.RouterGroup, c *controller.CommunityController) {
	communities := rg.Group("/admin/communities")
	// Protected routes (require authentication)
	communities.Use(middleware.RequireAdmin())
	{
		communities.GET("moderator/:moderator_id", c.GetCommunityByModeratorID)
		communities.GET("", c.GetAllCommunities)
	}
}
