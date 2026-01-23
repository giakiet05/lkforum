package route

import (
	"github.com/giakiet05/lkforum/internal/controller"
	"github.com/giakiet05/lkforum/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterMessageRoutes(rg *gin.RouterGroup, c *controller.MessageController) {
	messages := rg.Group("/messages")

	// Protected routes (require authentication)
	messages.Use(middleware.RequireAuth())
	{
		messages.GET(":message_id", c.GetMessageByID)
		messages.GET("filter", c.GetMessageFilter)
		messages.DELETE(":message_id", c.DeleteMessageByID)
		messages.PUT("/channels/:channel_id/read", c.MarkChannelAsRead)
	}
}
