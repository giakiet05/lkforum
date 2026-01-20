package route

import (
	"github.com/giakiet05/lkforum/internal/controller"
	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes registers all authentication-related routes.
func RegisterAuthRoutes(rg *gin.RouterGroup, authCtrl *controller.AuthController, userCtrl *controller.UserController) {
	auth := rg.Group("/auth")

	auth.POST("/refresh", authCtrl.RefreshToken)
	auth.POST("/logout", authCtrl.Logout)
	auth.POST("/check-username", userCtrl.CheckUsername) // Public endpoint for username availability check

	// Local Authentication - New Flow (Verify Email First)
	local := auth.Group("/local")
	{
		local.POST("/send-verification", authCtrl.SendEmailVerification)
		local.POST("/verify-email", authCtrl.VerifyEmailCode)
		local.POST("/complete-registration", authCtrl.CompleteRegistration)
		local.POST("/resend-otp", authCtrl.ResendOTP)
		local.POST("/login", authCtrl.Login)

		// Forgot Password Flow (only for local auth)
		local.POST("/forgot-password", authCtrl.ForgotPassword)
		local.POST("/verify-reset-otp", authCtrl.VerifyResetPasswordOTP)
		local.POST("/reset-password", authCtrl.ResetPassword)
	}

	// Google OAuth2
	google := auth.Group("/google")
	{
		google.GET("/login", authCtrl.GoogleLogin)
		google.GET("/callback", authCtrl.GoogleCallback)
		google.POST("/complete-setup", authCtrl.CompleteGoogleSetup)
	}
}
