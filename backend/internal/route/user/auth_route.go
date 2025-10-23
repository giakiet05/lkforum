package userroute

import (
	"github.com/giakiet05/lkforum/internal/controller"
	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes registers all authentication-related routes.
func RegisterAuthRoutes(rg *gin.RouterGroup, c *controller.AuthController) {
	auth := rg.Group("/auth")

	// Local Authentication
	auth.POST("/register", c.RegisterUser)
	auth.POST("/login", c.Login)
	auth.POST("/refresh", c.RefreshToken)
	auth.POST("/verify-email", c.VerifyEmail)
	auth.POST("/resend-verification-email", c.ResendVerificationEmail)

	// Google OAuth2
	google := auth.Group("/google")
	{
		google.GET("/login", c.GoogleLogin)
		google.GET("/callback", c.GoogleCallback)
		google.POST("/complete-setup", c.CompleteGoogleSetup)
	}
}
