package route

import (
	"github.com/giakiet05/lkforum/internal/controller"
	"github.com/giakiet05/lkforum/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterWebSocketRoutes(rg *gin.RouterGroup, c *controller.WebSocketController) {
	ws := rg.Group("/ws")
	ws.Use(middleware.RequireAuthSocket())
	{
		ws.GET("", c.HandleConnections)
	}
}
