package route

import (
	"github.com/giakiet05/lkforum/internal/controller"
	"github.com/giakiet05/lkforum/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterChannelRoutes(rg *gin.RouterGroup, c *controller.ChannelController) {
	channels := rg.Group("/channels")
	// Protected routes (require authentication)
	channels.Use(middleware.RequireAuth())
	{
		channels.POST("", c.CreateChannel)
		channels.GET(":channel_id", c.GetChannelByID)
		channels.GET("user", c.GetChannelsByUserID)
		channels.GET("/between/:user1/:user2", c.GetChannelByBothUserID)
		channels.PUT("", c.UpdateChannel)
		channels.DELETE(":channel_id", c.DeleteChannelByID)
	}
}
