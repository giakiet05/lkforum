package route

import (
	"github.com/giakiet05/lkforum/internal/controller"
	"github.com/giakiet05/lkforum/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterCommunityRoutes(rg *gin.RouterGroup, c *controller.CommunityController) {
	communities := rg.Group("/communities")

	communities.GET("filter", c.GetCommunitiesFilter)
	communities.GET("user/:user_id", c.GetCommunitiesByUserID)
	communities.GET("name/:name", c.GetCommunityByName)
	communities.GET(":community_id", c.GetCommunityByID)

	// Protected routes (require authentication)
	communities.Use(middleware.RequireAuth())
	{
		communities.POST("", c.CreateCommunity)
		communities.PUT("", c.UpdateCommunity)
		communities.PUT("add_moderator", c.AddModerator)
		communities.PUT("activate_moderator/:community_id", c.ActivateModerator)
		communities.PUT("remove_moderator", c.RemoveModerator)
		communities.DELETE(":community_id", c.DeleteCommunityByID)

		communities.POST("ban/user", c.BanUser)
		communities.GET("banned_user", c.GetBanUsers)
		communities.POST("unban/user", c.UnbanUser)
		communities.POST("unmute/user", c.UnbanUser)

		communities.PUT("ban/post", c.BanPost)
		communities.PUT("unban/post", c.UnbanPost)

		// Manual moderation routes (moderator only)
		communities.GET(":community_id/posts/pending", c.GetPendingPosts)
		communities.GET(":community_id/posts/edited", c.GetEditedPosts)
		communities.PUT(":community_id/posts/:post_id/moderate", c.ModeratePost)
	}
}
