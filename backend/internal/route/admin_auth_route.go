package route

import (
	"github.com/giakiet05/lkforum/internal/controller"
	"github.com/giakiet05/lkforum/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterAdminAuthRoutes(rg *gin.RouterGroup, c *controller.AdminAuthController) {
	auth := rg.Group("/admin/auth")
	{
		// Public routes
		auth.POST("/login", c.AdminLogin)
		auth.POST("/refresh", c.AdminRefreshToken)

		// Protected routes (require admin authentication)
		auth.POST("/logout", middleware.RequireAuth(), middleware.RequireAdmin(), c.AdminLogout)
	}
}
