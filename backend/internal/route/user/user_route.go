package route

import (
	"github.com/giakiet05/lkforum/internal/controller"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(rg *gin.RouterGroup, c *controller.UserController) {
	users := rg.Group("/users")

	// Protected routes (require authentication)
	//users.Use(middleware.AuthMiddleware())
	{
		users.GET("", c.GetUsers)
		users.GET(":id", c.GetUserByID)
		users.PUT(":id", c.UpdateUser)
		users.PUT(":id/change-password", c.ChangePassword)
		users.DELETE(":id", c.DeleteUser)
	}
}
