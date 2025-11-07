package userroute

import (
	"github.com/giakiet05/lkforum/internal/controller"
	"github.com/giakiet05/lkforum/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(rg *gin.RouterGroup, c *controller.UserController) {
	users := rg.Group("/users")

	// Public routes - anyone can view a user's profile
	users.GET("/profile/:username", c.GetUserByUsername)

	// Routes for the currently authenticated user ("me")
	me := users.Group("/me")
	me.Use(middleware.RequireAuth())
	{
		me.GET("", c.GetMyProfile)
		me.PUT("/profile", c.UpdateProfile)
		me.PUT("/password", c.ChangePassword)
		me.POST("/avatar", c.UploadAvatar)
		me.POST("/cover", c.UploadCover)
	}
}
