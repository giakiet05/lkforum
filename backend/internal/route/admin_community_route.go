package route

import (
	"github.com/giakiet05/lkforum/internal/controller"
	"github.com/giakiet05/lkforum/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterAdminCommunityRoutes(rg *gin.RouterGroup, c *controller.CommunityController, ac *controller.AdminCommunityController) {
	communities := rg.Group("/admin/communities")
	// Protected routes (require authentication)
	communities.Use(middleware.RequireAuth(), middleware.RequireAdmin())
	{
		communities.GET("moderator/:moderator_id", c.GetCommunityByModeratorID)
		communities.GET("", ac.GetCommunitiesAdmin)
		communities.POST("/:id/ban", ac.BanCommunity)
		communities.POST("/:id/unban", ac.UnbanCommunity)
	}
}
