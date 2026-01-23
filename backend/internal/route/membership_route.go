package route

import (
	"github.com/giakiet05/lkforum/internal/controller"
	"github.com/giakiet05/lkforum/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterMembershipRoutes(rg *gin.RouterGroup, c *controller.MembershipController) {
	memberships := rg.Group("/memberships")

	// Protected routes (require authentication)
	memberships.Use(middleware.RequireAuth())
	{
		memberships.POST("", c.CreateMembership)
		memberships.GET("", c.GetAllMemberships)
		memberships.GET("/user/:user_id", c.GetMembershipByUserID)
		memberships.GET("/community/:community_id", c.GetMembershipByCommunityID)
		memberships.GET("/:membership_id", c.GetMembershipByID)
		memberships.DELETE("", c.DeleteMembership)
		memberships.DELETE("/kick/:community_id/:user_id", c.KickMember) // Moderator/Creator kick member

		// Pending members management (for moderators/creators)
		memberships.GET("/community/:community_id/pending", c.GetPendingMembers)
		memberships.GET("/community/:community_id/approved", c.GetApprovedMembers)
		memberships.POST("/community/:community_id/approve/:membership_id", c.ApproveMember)
		memberships.POST("/community/:community_id/reject/:membership_id", c.RejectMember)
	}
}
