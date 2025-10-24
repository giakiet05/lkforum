package userroute

import (
	"github.com/giakiet05/lkforum/internal/controller"
	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes registers all authentication-related routes.
func RegisterAuthRoutes(rg *gin.RouterGroup, c *controller.AuthController) {
	auth := rg.Group("/auth")

	auth.POST("/refresh", c.RefreshToken)

	// Local Authentication
	local := auth.Group("/local")
	{
		local.POST("/register", c.RegisterUser)
		local.POST("/login", c.Login)
		local.POST("/verify-email", c.VerifyEmail)
		local.POST("/resend-verification-email", c.ResendVerificationEmail)
	}

	// Google OAuth2
	google := auth.Group("/google")
	{
		google.GET("/login", c.GoogleLogin)
		google.GET("/callback", c.GoogleCallback)
		google.POST("/complete-setup", c.CompleteGoogleSetup)
	}
}
