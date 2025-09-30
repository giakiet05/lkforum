package route

import (
	"github.com/giakiet05/lkforum/internal/controller"
	"github.com/giakiet05/lkforum/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(rg *gin.RouterGroup, c *controller.UserController) {
	// Public routes
	users := rg.Group("/users")             // Get all users (paginated)
	users.POST("/register", c.RegisterUser) // Register new user
	users.POST("/login", c.Login)           // Login

	// Protected routes (require authentication)
	authUsers := users.Group("/")
	authUsers.Use(middleware.AuthMiddleware())
	{
		authUsers.GET("", c.GetUsers)
		authUsers.GET("/:id", c.GetUserByID)
		// Get user by ID
		authUsers.PUT("/:id", c.UpdateUser)                     // Update user
		authUsers.PUT("/:id/change-password", c.ChangePassword) // Change password

		authUsers.DELETE("/:id", c.DeleteUser) // Delete user
	}
}
