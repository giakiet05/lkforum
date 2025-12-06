package route

import (
	"github.com/giakiet05/lkforum/internal/controller"
	"github.com/giakiet05/lkforum/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterAdminStatsRoutes(rg *gin.RouterGroup, ac *controller.AdminStatsController) {
	stats := rg.Group("/admin/stats")
	// Protected routes (require admin authentication)
	stats.Use(middleware.RequireAuth(), middleware.RequireAdmin())
	{
		stats.GET("/overview", ac.GetPlatformOverview) // Main dashboard
		stats.GET("/users", ac.GetUserStats)           // User statistics with growth
		stats.GET("/content", ac.GetContentStats)      // Content statistics with growth
	}
}
