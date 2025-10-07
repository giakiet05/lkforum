package route

import (
	"github.com/giakiet05/lkforum/internal/controller"
	"github.com/giakiet05/lkforum/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterMembershipRoutes(rg *gin.RouterGroup, c *controller.MembershipController) {
	memberships := rg.Group("/memberships")

	// Protected routes (require authentication)
	memberships.Use(middleware.AuthMiddleware())
	{
		memberships.POST("", c.CreateMembership)
		memberships.GET("", c.GetAllMemberships)
		memberships.GET("/user/:user_id", c.GetMembershipByUserID)
		memberships.GET("/community/:community_id", c.GetMembershipByCommunityID)
		memberships.GET("/:membership_id", c.GetMembershipByID)
		memberships.DELETE("", c.DeleteMembership)
	}
}
