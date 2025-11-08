package userroute

import (
	"github.com/giakiet05/lkforum/internal/controller"
	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes registers all authentication-related routes.
func RegisterAuthRoutes(rg *gin.RouterGroup, c *controller.AuthController) {
	auth := rg.Group("/auth")

	auth.POST("/refresh", c.RefreshToken)

	// Local Authentication - New Flow (Verify Email First)
	local := auth.Group("/local")
	{
		local.POST("/send-verification", c.SendEmailVerification)
		local.POST("/verify-email", c.VerifyEmailCode)
		local.POST("/complete-registration", c.CompleteRegistration)
		local.POST("/resend-otp", c.ResendOTP)
		local.POST("/login", c.Login)
	}

	// Google OAuth2
	google := auth.Group("/google")
	{
		google.GET("/login", c.GoogleLogin)
		google.GET("/callback", c.GoogleCallback)
		google.POST("/complete-setup", c.CompleteGoogleSetup)
	}
}
