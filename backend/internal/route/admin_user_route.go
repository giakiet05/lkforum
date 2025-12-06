package route

import (
	"github.com/giakiet05/lkforum/internal/controller"
	"github.com/giakiet05/lkforum/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterAdminUserRoutes(rg *gin.RouterGroup, c *controller.AdminUserController) {
	admin := rg.Group("/admin/users")

	// All admin routes require authentication AND admin role
	admin.Use(middleware.RequireAuth(), middleware.RequireAdmin())
	{
		// User management
		admin.GET("", c.GetUsers)
		admin.POST("/:user_id/ban", c.BanUser)
		admin.POST("/:user_id/unban", c.UnbanUser)
		admin.DELETE("/:user_id", c.DeleteUser)
		admin.POST("/:user_id/restore", c.RestoreUser)
	}
}
