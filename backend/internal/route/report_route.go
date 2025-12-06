package route

import (
	"github.com/giakiet05/lkforum/internal/controller"
	"github.com/giakiet05/lkforum/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterReportRoutes(rg *gin.RouterGroup, c *controller.ReportController) {
	reports := rg.Group("/reports")

	// Protected routes (require authentication)
	reports.Use(middleware.RequireAuth())
	{
		reports.GET("reporter", c.GetReportByReporterID)
		reports.GET(":report_id", c.GetReportByID)
		reports.DELETE("batch", c.DeleteReportsByID)
		reports.DELETE(":report_id", c.DeleteReportByID)
	}
}
