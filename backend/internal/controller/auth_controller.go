package controller

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/auth"
	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthController handles authentication-related requests.
type AuthController struct {
	authService service.AuthService
}

// redirectWithHash redirects to a URL with hash fragment using HTML meta refresh
// This is necessary because HTTP redirects strip hash fragments
func redirectWithHash(ctx *gin.Context, url string) {
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Redirecting...</title>
    <script>
        window.location.href = %q;
    </script>
</head>
<body>
    <p>Redirecting...</p>
</body>
</html>
`, url)
	ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// NewAuthController creates a new AuthController.
func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// --- Local Authentication - New Flow (Verify Email First) ---

func (c *AuthController) SendEmailVerification(ctx *gin.Context) {
	var req dto.SendEmailVerificationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	err := c.authService.SendEmailVerification(req.Email)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Verification code sent to your email. Please check your inbox.", nil)
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req dto.UserLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	user, accessToken, refreshToken, err := c.authService.Login(req.Identifier, req.Password)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	data := dto.AuthResponse{
		User:         dto.FromUser(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	dto.SendSuccess(ctx, http.StatusOK, "Login successful", data)
}

func (c *AuthController) VerifyEmailCode(ctx *gin.Context) {
	var req dto.VerifyEmailCodeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	verificationToken, err := c.authService.VerifyEmailCode(req.Email, req.OTP)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	data := gin.H{"verification_token": verificationToken}
	dto.SendSuccess(ctx, http.StatusOK, "Email verified successfully. You can now complete your registration.", data)
}

func (c *AuthController) CompleteRegistration(ctx *gin.Context) {
	var req dto.CompleteRegistrationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	user, accessToken, refreshToken, err := c.authService.CompleteRegistration(req.VerificationToken, req.Username, req.Password)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	data := dto.AuthResponse{
		User:         dto.FromUser(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	dto.SendSuccess(ctx, http.StatusCreated, "Registration completed successfully. You are now logged in.", data)
}

func (c *AuthController) ResendOTP(ctx *gin.Context) {
	var req dto.ResendOTPRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	err := c.authService.ResendOTP(req.Email)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "A new verification code has been sent to your email.", nil)
}

func (c *AuthController) RefreshToken(ctx *gin.Context) {
	var req dto.RefreshRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	accessToken, refreshToken, err := c.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	data := dto.RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	dto.SendSuccess(ctx, http.StatusOK, "Tokens refreshed successfully", data)
}

func (c *AuthController) Logout(ctx *gin.Context) {
	var req dto.LogoutRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	err := c.authService.Logout(req.AccessToken, req.RefreshToken)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Logged out successfully", nil)
}

// --- Google OAuth ---

func (c *AuthController) GoogleLogin(ctx *gin.Context) {
	state := uuid.New().String()
	url := auth.GetGoogleLoginURL(state)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (c *AuthController) GoogleCallback(ctx *gin.Context) {
	code := ctx.Query("code")
	if code == "" {
		// Redirect to FE with error
		redirectURL := fmt.Sprintf("%s/#/auth/error?message=missing_auth_code", config.Cfg.FrontendURL)
		log.Printf("GoogleCallback: Missing code, redirecting to: %s", redirectURL)
		redirectWithHash(ctx, redirectURL)
		return
	}

	log.Printf("GoogleCallback: Processing code: %s", code[:10]+"...")

	result, err := c.authService.ProcessGoogleCallback(code)
	if err != nil {
		// Redirect to FE with error
		redirectURL := fmt.Sprintf("%s/#/auth/error?message=%s", config.Cfg.FrontendURL, url.QueryEscape(apperror.Message(err)))
		log.Printf("GoogleCallback: Error processing callback: %v, redirecting to: %s", err, redirectURL)
		redirectWithHash(ctx, redirectURL)
		return
	}

	log.Printf("GoogleCallback: Result status: %s", result.Status)

	switch result.Status {
	case service.StatusLoginSuccess:
		// Redirect to FE with tokens in query params (can't use hash fragment due to SPA router limitation)
		redirectURL := fmt.Sprintf("%s/#/auth/callback?access_token=%s&refresh_token=%s",
			config.Cfg.FrontendURL,
			url.QueryEscape(result.AccessToken),
			url.QueryEscape(result.RefreshToken))
		log.Printf("GoogleCallback: Login success, redirecting to: %s", redirectURL)
		redirectWithHash(ctx, redirectURL)

	case service.StatusSetupRequired:
		// Redirect to FE with setup_token in query params (can't use hash fragment due to SPA router limitation)
		redirectURL := fmt.Sprintf("%s/#/auth/google-setup?setup_token=%s",
			config.Cfg.FrontendURL,
			url.QueryEscape(result.SetupToken))
		log.Printf("GoogleCallback: Setup required, redirecting to: %s", redirectURL)
		redirectWithHash(ctx, redirectURL)

	default:
		// Redirect to FE with error
		redirectURL := fmt.Sprintf("%s/#/auth/error?message=unknown_error", config.Cfg.FrontendURL)
		log.Printf("GoogleCallback: Unknown status, redirecting to: %s", redirectURL)
		redirectWithHash(ctx, redirectURL)
	}
}

func (c *AuthController) CompleteGoogleSetup(ctx *gin.Context) {
	var req dto.CompleteGoogleSetupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	user, accessToken, refreshToken, err := c.authService.CompleteGoogleSetup(req.SetupToken, req.Username)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	data := dto.AuthResponse{
		User:         dto.FromUser(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	dto.SendSuccess(ctx, http.StatusOK, "Setup complete. You are now logged in.", data)
}

// --- Forgot Password Flow ---

func (c *AuthController) ForgotPassword(ctx *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	err := c.authService.ForgotPassword(req.Email)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Password reset code sent to your email. Please check your inbox.", nil)
}

func (c *AuthController) VerifyResetPasswordOTP(ctx *gin.Context) {
	var req dto.VerifyResetPasswordOTPRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	resetToken, err := c.authService.VerifyResetPasswordOTP(req.Email, req.OTP)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	data := dto.VerifyResetPasswordOTPResponse{
		ResetToken: resetToken,
	}
	dto.SendSuccess(ctx, http.StatusOK, "OTP verified successfully. You can now reset your password.", data)
}

func (c *AuthController) ResetPassword(ctx *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	err := c.authService.ResetPassword(req.ResetToken, req.NewPassword)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Password reset successfully. You can now login with your new password.", nil)
}
