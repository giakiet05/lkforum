package route

import (
	"github.com/giakiet05/lkforum/internal/controller"
	"github.com/giakiet05/lkforum/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterCommentRoutes(rg *gin.RouterGroup, c *controller.CommentController) {
	comments := rg.Group("/comments")

	comments.GET(":comment_id", c.GetCommentByID)
	comments.GET("filter", c.GetCommentsFilter)
	comments.GET("post", c.GetCommentByPostID)

	// Protected routes (require authentication)
	comments.Use(middleware.RequireAuth())
	{
		comments.POST("", c.CreateComment)
		comments.DELETE("/:comment_id", c.DeleteCommentByID)
	}
}
