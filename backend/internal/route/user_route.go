package route

import (
	"github.com/giakiet05/lkforum/internal/controller"
	"github.com/giakiet05/lkforum/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(rg *gin.RouterGroup, c *controller.UserController) {
	users := rg.Group("/users")

	// Public search endpoint - anyone can search users
	users.GET("", c.GetUsers)

	// Public routes - anyone can view a user's profile
	users.GET("/", c.GetUsers)
	users.GET("/profile/:username", c.GetUserByUsername)
	// Metadata endpoints - public, used for dropdowns
	users.GET("/provinces", c.GetProvinces)
	users.GET("/interests", c.GetInterests)
	users.GET("/genders", c.GetGenders)

	// Routes for the currently authenticated user ("me")
	me := users.Group("/me")
	me.Use(middleware.RequireAuth())
	{
		me.GET("", c.GetMyProfile)
		me.PUT("/profile", c.UpdateProfile)
		me.PUT("/password", c.ChangePassword)
		me.POST("/avatar", c.UploadAvatar)
		me.DELETE("/avatar", c.DeleteAvatar)
		me.POST("/cover", c.UploadCover)
		me.DELETE("/cover", c.DeleteCover)
		me.GET("/settings", c.GetSettings)
		me.PUT("/settings", c.UpdateSettings)
	}
}
