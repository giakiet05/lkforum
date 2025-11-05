package userroute

import (
	"github.com/giakiet05/lkforum/internal/controller"
	"github.com/giakiet05/lkforum/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterPostHistoryRoutes(rg *gin.RouterGroup, c *controller.PostHistoryController) {
	postHistories := rg.Group("/post_histories")

	// Protected routes (require authentication)
	postHistories.Use(middleware.AuthMiddleware())
	{
		postHistories.POST("", c.CreatePostHistory)
		postHistories.POST("batch", c.CreatePostHistories)
		postHistories.GET("user/:user_id", c.GetPostHistoryByUserID)
		postHistories.GET(":id", c.GetPostHistoryByID)
		postHistories.DELETE(":id", c.DeletePostHistoryByID)
	}
}
