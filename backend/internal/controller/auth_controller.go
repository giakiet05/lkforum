package controller

import (
	"fmt"
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

	data := gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}
	dto.SendSuccess(ctx, http.StatusOK, "Tokens refreshed successfully", data)
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
		redirectURL := fmt.Sprintf("%s/auth/error?message=missing_auth_code", config.Cfg.FrontendURL)
		ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)
		return
	}

	result, err := c.authService.ProcessGoogleCallback(code)
	if err != nil {
		// Redirect to FE with error
		redirectURL := fmt.Sprintf("%s/auth/error?message=%s", config.Cfg.FrontendURL, url.QueryEscape(apperror.Message(err)))
		ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)
		return
	}

	switch result.Status {
	case service.StatusLoginSuccess:
		// Redirect to FE with tokens in hash fragment
		redirectURL := fmt.Sprintf("%s/auth/callback#access_token=%s&refresh_token=%s",
			config.Cfg.FrontendURL,
			url.QueryEscape(result.AccessToken),
			url.QueryEscape(result.RefreshToken))
		ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)

	case service.StatusSetupRequired:
		// Redirect to FE with setup_token in hash fragment
		redirectURL := fmt.Sprintf("%s/auth/google-setup#setup_token=%s",
			config.Cfg.FrontendURL,
			url.QueryEscape(result.SetupToken))
		ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)

	default:
		// Redirect to FE with error
		redirectURL := fmt.Sprintf("%s/auth/error?message=unknown_error", config.Cfg.FrontendURL)
		ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)
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
