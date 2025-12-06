package route

import (
	"github.com/giakiet05/lkforum/internal/controller"
	"github.com/giakiet05/lkforum/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterAdminReportRoutes(rg *gin.RouterGroup, rc *controller.ReportController) {
	reports := rg.Group("/admin/reports")
	// Protected routes (require admin authentication)
	reports.Use(middleware.RequireAuth(), middleware.RequireAdmin())
	{
		reports.GET("", rc.GetReportsFilter)                    // Get all reports with filters
		reports.GET("/:report_id", rc.GetReportByIDAdmin)       // Get report by ID (admin)
		reports.DELETE("/:report_id", rc.DeleteReportByIDAdmin) // Delete single report (admin)
		reports.DELETE("/batch", rc.DeleteReportsByIDAdmin)     // Delete multiple reports (admin)
	}
}
