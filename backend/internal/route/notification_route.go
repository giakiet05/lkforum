package route

import (
	"github.com/giakiet05/lkforum/internal/controller"
	"github.com/giakiet05/lkforum/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterNotificationRoutes(rg *gin.RouterGroup, c *controller.NotificationController) {
	notifications := rg.Group("/notifications")
	notifications.Use(middleware.RequireAuth()) // All notification routes require authentication
	{
		notifications.GET("", c.GetNotifications)
		notifications.PUT("/read-all", c.MarkAllAsRead) // Must be before /:notification_id routes
		notifications.PUT("/:notification_id/read", c.MarkAsRead)
		notifications.DELETE("/:notification_id", c.DeleteNotification)
	}
}
