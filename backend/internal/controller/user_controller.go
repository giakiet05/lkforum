package controller

import (
	"net/http"
	"strconv"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/auth"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/service"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	service service.UserService
}

func NewUserController(service service.UserService) *UserController {
	return &UserController{service: service}
}

// GetUsers returns a paginated list of users
// If no pagination parameters are provided, it uses sensible defaults
func (c *UserController) GetUsers(ctx *gin.Context) {
	// Parse pagination parameters
	page := 1
	pageSize := 10

	if pageStr := ctx.Query("page"); pageStr != "" {
		if parsedPage, err := strconv.Atoi(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	if pageSizeStr := ctx.Query("pageSize"); pageSizeStr != "" {
		if parsedPageSize, err := strconv.Atoi(pageSizeStr); err == nil && parsedPageSize > 0 && parsedPageSize <= 100 {
			pageSize = parsedPageSize
		}
	}

	// Call the service
	response, err := c.service.GetUsers(page, pageSize)
	if err != nil {
		ctx.JSON(apperror.StatusFromError(err), dto.ErrorResponse{ErrorCode: apperror.Code(err), Message: apperror.Message(err)})
		return
	}
	ctx.JSON(http.StatusOK, response)
}

// RegisterUser handles user registration
func (c *UserController) RegisterUser(ctx *gin.Context) {
	var req dto.UserRegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(apperror.StatusFromError(apperror.ErrBadRequest), dto.ErrorResponse{ErrorCode: apperror.ErrBadRequest.Code, Message: apperror.Message(err)})
		return
	}

	user, accessToken, refreshToken, err := c.service.RegisterUser(req.Username, req.Email, req.Password)
	if err != nil {
		ctx.JSON(apperror.StatusFromError(err), dto.ErrorResponse{ErrorCode: apperror.Code(err), Message: apperror.Message(err)})
		return
	}

	ctx.JSON(http.StatusCreated, dto.AuthResponse{
		User:         dto.FromUser(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// Login handles user authentication
func (c *UserController) Login(ctx *gin.Context) {
	var req dto.UserLoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(apperror.StatusFromError(apperror.ErrBadRequest), dto.ErrorResponse{ErrorCode: apperror.ErrBadRequest.Code, Message: apperror.Message(err)})
		return
	}

	user, accessToken, refreshToken, err := c.service.Login(req.Identifier, req.Password)
	if err != nil {
		ctx.JSON(apperror.StatusFromError(err), dto.ErrorResponse{ErrorCode: apperror.Code(err), Message: apperror.Message(err)})
		return
	}

	ctx.JSON(http.StatusOK, dto.AuthResponse{
		User:         dto.FromUser(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// UpdateUser handles user profile updates
func (c *UserController) UpdateUser(ctx *gin.Context) {
	userID := ctx.Param("id")
	if !auth.IsOwner(ctx, userID) && !auth.IsAdmin(ctx) {
		ctx.JSON(apperror.StatusFromError(apperror.ErrForbidden), dto.ErrorResponse{ErrorCode: apperror.ErrForbidden.Code, Message: apperror.ErrForbidden.Message})
		return
	}

	var req dto.UserUpdateRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(apperror.StatusFromError(apperror.ErrBadRequest), dto.ErrorResponse{ErrorCode: apperror.ErrBadRequest.Code, Message: apperror.Message(err)})
		return
	}

	currentUser, err := c.service.GetUserByID(userID)
	if err != nil {
		ctx.JSON(apperror.StatusFromError(err), dto.ErrorResponse{ErrorCode: apperror.Code(err), Message: apperror.Message(err)})
		return
	}

	if req.Username != "" {
		currentUser.Username = req.Username
	}
	if req.Email != "" {
		currentUser.Email = req.Email
	}

	updatedUser, err := c.service.UpdateUser(currentUser)
	if err != nil {
		ctx.JSON(apperror.StatusFromError(err), dto.ErrorResponse{ErrorCode: apperror.Code(err), Message: apperror.Message(err)})
		return
	}

	ctx.JSON(http.StatusOK, dto.FromUser(updatedUser))
}

// DeleteUser handles user account deletion
func (c *UserController) DeleteUser(ctx *gin.Context) {
	userID := ctx.Param("id")
	if !auth.IsOwner(ctx, userID) && !auth.IsAdmin(ctx) {
		ctx.JSON(apperror.StatusFromError(apperror.ErrForbidden), dto.ErrorResponse{ErrorCode: apperror.ErrForbidden.Code, Message: apperror.ErrForbidden.Message})
		return
	}

	err := c.service.DeleteUser(userID)
	if err != nil {
		ctx.JSON(apperror.StatusFromError(err), dto.ErrorResponse{ErrorCode: apperror.Code(err), Message: apperror.Message(err)})
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		ID:      userID,
		Message: "User deleted successfully",
	})
}

// GetUserByID retrieves user details by ID
func (c *UserController) GetUserByID(ctx *gin.Context) {
	userID := ctx.Param("id")

	user, err := c.service.GetUserByID(userID)
	if err != nil {
		ctx.JSON(apperror.StatusFromError(err), dto.ErrorResponse{ErrorCode: apperror.Code(err), Message: apperror.Message(err)})
		return
	}

	ctx.JSON(http.StatusOK, dto.FromUser(user))
}

// GetUserByUsername retrieves user details by username
func (c *UserController) GetUserByUsername(ctx *gin.Context) {
	username := ctx.Param("username")

	user, err := c.service.GetUserByUsername(username)
	if err != nil {
		ctx.JSON(apperror.StatusFromError(err), dto.ErrorResponse{ErrorCode: apperror.Code(err), Message: apperror.Message(err)})
		return
	}

	ctx.JSON(http.StatusOK, dto.FromUser(user))
}

// ChangePassword handles password changes
func (c *UserController) ChangePassword(ctx *gin.Context) {
	userID := ctx.Param("id")
	authUser, exists := ctx.Get("authUser")
	if !exists || authUser.(auth.AuthUser).ID != userID {
		ctx.JSON(apperror.StatusFromError(apperror.ErrForbidden), dto.ErrorResponse{ErrorCode: apperror.ErrForbidden.Code, Message: apperror.ErrForbidden.Message})
		return
	}

	var req dto.ChangePasswordRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(apperror.StatusFromError(apperror.ErrBadRequest), dto.ErrorResponse{ErrorCode: apperror.ErrBadRequest.Code, Message: apperror.Message(err)})
		return
	}

	err := c.service.ChangePassword(userID, req.OldPassword, req.NewPassword)
	if err != nil {
		ctx.JSON(apperror.StatusFromError(err), dto.ErrorResponse{ErrorCode: apperror.Code(err), Message: apperror.Message(err)})
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		ID:      userID,
		Message: "Password changed successfully",
	})
}

// RefreshToken handles token refresh requests
func (c *UserController) RefreshToken(ctx *gin.Context) {
	type RefreshRequest struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	var req RefreshRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(apperror.StatusFromError(apperror.ErrBadRequest), dto.ErrorResponse{
			ErrorCode: apperror.ErrBadRequest.Code,
			Message:   apperror.ErrBadRequest.Message,
		})
		return
	}

	accessToken, refreshToken, err := c.service.RefreshToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(apperror.StatusFromError(err), dto.ErrorResponse{
			ErrorCode: apperror.Code(err),
			Message:   apperror.Message(err),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
