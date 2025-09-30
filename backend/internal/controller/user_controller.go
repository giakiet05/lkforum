package controller

import (
	"errors"
	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/auth"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: apperror.Code(err), Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// RegisterUser handles user registration
func (c *UserController) RegisterUser(ctx *gin.Context) {
	var req dto.UserRegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: "INVALID_REQUEST", Error: err.Error()})
		return
	}

	user, accessToken, refreshToken, err := c.service.RegisterUser(req.Username, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, apperror.ErrUsernameExists) || errors.Is(err, apperror.ErrEmailExists) {
			ctx.JSON(http.StatusConflict, dto.ErrorResponse{Code: apperror.Code(err), Error: err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: apperror.Code(err), Error: err.Error()})
		}
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
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: "INVALID_REQUEST", Error: err.Error()})
		return
	}

	user, accessToken, refreshToken, err := c.service.Login(req.Identifier, req.Password, req.LoginType)
	if err != nil {
		if errors.Is(err, apperror.ErrInvalidCredentials) {
			ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{Code: apperror.Code(err), Error: err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: apperror.Code(err), Error: err.Error()})
		}
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
		ctx.JSON(http.StatusForbidden, dto.ErrorResponse{Code: "FORBIDDEN", Error: "Not authorized to update this user"})
		return
	}

	var req dto.UserUpdateRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: "INVALID_REQUEST", Error: err.Error()})
		return
	}

	// Get the current user first
	currentUser, err := c.service.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, apperror.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{Code: apperror.Code(err), Error: err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: apperror.Code(err), Error: err.Error()})
		}
		return
	}

	// Update fields that were provided
	if req.Username != "" {
		currentUser.Username = req.Username
	}
	if req.Email != "" {
		currentUser.Email = req.Email
	}

	updatedUser, err := c.service.UpdateUser(currentUser)
	if err != nil {
		if errors.Is(err, apperror.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{Code: apperror.Code(err), Error: err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: apperror.Code(err), Error: err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, dto.FromUser(updatedUser))
}

// DeleteUser handles user account deletion
func (c *UserController) DeleteUser(ctx *gin.Context) {
	userID := ctx.Param("id")
	if !auth.IsOwner(ctx, userID) && !auth.IsAdmin(ctx) {
		ctx.JSON(http.StatusForbidden, dto.ErrorResponse{Code: "FORBIDDEN", Error: "Not authorized to delete this user"})
		return
	}

	err := c.service.DeleteUser(userID)
	if err != nil {
		if errors.Is(err, apperror.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{Code: apperror.Code(err), Error: err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: apperror.Code(err), Error: err.Error()})
		}
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
		if errors.Is(err, apperror.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{Code: apperror.Code(err), Error: err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: apperror.Code(err), Error: err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, dto.FromUser(user))
}

// GetUserByUsername retrieves user details by username
func (c *UserController) GetUserByUsername(ctx *gin.Context) {
	username := ctx.Param("username")

	user, err := c.service.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, apperror.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{Code: apperror.Code(err), Error: err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: apperror.Code(err), Error: err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, dto.FromUser(user))
}

// ChangePassword handles password changes
func (c *UserController) ChangePassword(ctx *gin.Context) {
	userID := ctx.Param("id")
	// Check if the authenticated user is the owner
	authUser, exists := ctx.Get("authUser")
	if !exists || authUser.(auth.AuthUser).ID != userID {
		ctx.JSON(http.StatusForbidden, dto.ErrorResponse{Code: "FORBIDDEN", Error: "Not authorized to change this user's password"})
		return
	}

	var req dto.ChangePasswordRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: "INVALID_REQUEST", Error: err.Error()})
		return
	}

	err := c.service.ChangePassword(userID, req.OldPassword, req.NewPassword)
	if err != nil {
		if errors.Is(err, apperror.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{Code: apperror.Code(err), Error: err.Error()})
		} else if errors.Is(err, apperror.ErrInvalidCredentials) {
			ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{Code: apperror.Code(err), Error: err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: apperror.Code(err), Error: err.Error()})
		}
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
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:  "INVALID_REQUEST",
			Error: "Invalid request payload",
		})
		return
	}

	accessToken, refreshToken, err := c.service.RefreshToken(req.RefreshToken)
	if err != nil {
		code := apperror.Code(err)
		status := http.StatusInternalServerError
		if code == "INVALID_TOKEN" {
			status = http.StatusUnauthorized
		} else if code == "USER_NOT_FOUND" {
			status = http.StatusNotFound
		} else if err.Error() == "user account is inactive" {
			status = http.StatusForbidden
			code = "USER_INACTIVE"
		}
		ctx.JSON(status, dto.ErrorResponse{
			Code:  code,
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
