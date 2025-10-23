package controller

import (
	"net/http"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/auth"
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

// --- Local Authentication ---

func (c *AuthController) RegisterUser(ctx *gin.Context) {
	var req dto.UserRegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	user, err := c.authService.RegisterUser(req.Username, req.Email, req.Password)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusCreated, "Registration successful. Please check your email for a verification code.", gin.H{"user_id": user.ID.Hex()})
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

func (c *AuthController) VerifyEmail(ctx *gin.Context) {
	var req dto.VerifyEmailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	user, accessToken, refreshToken, err := c.authService.VerifyEmail(req.Email, req.OTP)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	data := dto.AuthResponse{
		User:         dto.FromUser(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	dto.SendSuccess(ctx, http.StatusOK, "Email verified successfully. You are now logged in.", data)
}

func (c *AuthController) ResendVerificationEmail(ctx *gin.Context) {
	var req dto.ResendVerificationEmailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	err := c.authService.ResendVerificationEmail(req.Email)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "A new verification email has been sent.", nil)
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
		dto.SendError(ctx, http.StatusBadRequest, "Authorization code not provided", "MISSING_AUTH_CODE")
		return
	}

	result, err := c.authService.ProcessGoogleCallback(code)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	switch result.Status {
	case service.StatusLoginSuccess:
		data := dto.AuthResponse{
			User:         dto.FromUser(result.User),
			AccessToken:  result.AccessToken,
			RefreshToken: result.RefreshToken,
		}
		dto.SendSuccess(ctx, http.StatusOK, "Login successful", data)

	case service.StatusSetupRequired:
		data := gin.H{"setup_token": result.SetupToken}
		ctx.AbortWithStatusJSON(http.StatusPreconditionRequired, dto.ApiResponse{
			Success:   false,
			Message:   "User requires setup. Please choose a username.",
			Data:      data,
			ErrorCode: "USER_SETUP_REQUIRED",
		})

	default:
		dto.SendError(ctx, http.StatusInternalServerError, "An unknown error occurred during Google login", apperror.ErrInternal.Code)
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
