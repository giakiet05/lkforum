package controller

import (
	"net/http"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/service"
	"github.com/gin-gonic/gin"
)

type AdminAuthController struct {
	adminAuthService service.AdminAuthService
}

func NewAdminAuthController(adminAuthService service.AdminAuthService) *AdminAuthController {
	return &AdminAuthController{
		adminAuthService: adminAuthService,
	}
}

// AdminLogin handles admin login with role check
func (c *AdminAuthController) AdminLogin(ctx *gin.Context) {
	var req dto.UserLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	user, accessToken, refreshToken, err := c.adminAuthService.AdminLogin(req.Identifier, req.Password)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	data := dto.AuthResponse{
		User:         dto.FromUser(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	dto.SendSuccess(ctx, http.StatusOK, "Admin login successful", data)
}

// AdminLogout handles admin logout
func (c *AdminAuthController) AdminLogout(ctx *gin.Context) {
	// Get tokens from request
	accessToken := ctx.GetHeader("Authorization")
	if len(accessToken) > 7 && accessToken[:7] == "Bearer " {
		accessToken = accessToken[7:]
	}

	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	_ = ctx.ShouldBindJSON(&req)

	err := c.adminAuthService.AdminLogout(accessToken, req.RefreshToken)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Admin logout successful", nil)
}

// AdminRefreshToken handles admin token refresh
func (c *AdminAuthController) AdminRefreshToken(ctx *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	newAccessToken, newRefreshToken, err := c.adminAuthService.AdminRefreshToken(req.RefreshToken)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	data := gin.H{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	}
	dto.SendSuccess(ctx, http.StatusOK, "Token refreshed successfully", data)
}
