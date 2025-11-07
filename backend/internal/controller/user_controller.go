package controller

import (
	"net/http"
	"strconv"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/auth"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/platform/cloudinary"
	"github.com/giakiet05/lkforum/internal/service"
	"github.com/gin-gonic/gin"
)

// UserController handles requests related to user management.
type UserController struct {
	service service.UserService
}

// NewUserController creates a new UserController.
func NewUserController(service service.UserService) *UserController {
	return &UserController{service: service}
}

// GetUsers retrieves a paginated list of users (for admin or public listing).
func (c *UserController) GetUsers(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	response, err := c.service.GetUsers(page, pageSize)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}
	dto.SendSuccess(ctx, http.StatusOK, "Users retrieved successfully", response)
}

// GetUserByUsername retrieves a user's public profile by their username.
func (c *UserController) GetUserByUsername(ctx *gin.Context) {
	username := ctx.Param("username")

	user, err := c.service.GetUserByUsername(username)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	publicProfile := dto.FromUser(user)
	publicProfile.Email = ""
	dto.SendSuccess(ctx, http.StatusOK, "User profile retrieved successfully", publicProfile)
}

// GetMyProfile retrieves the profile of the currently authenticated user.
func (c *UserController) GetMyProfile(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	user, err := c.service.GetUserByID(authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Profile retrieved successfully", dto.FromUser(user))
}

// UpdateProfile allows a user to update their own profile information.
func (c *UserController) UpdateProfile(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	var req dto.UserProfileUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	updatedUser, err := c.service.UpdateProfile(authUser.(auth.AuthUser).ID, &req)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Profile updated successfully", dto.FromUser(updatedUser))
}

// UploadAvatar handles avatar image uploads.
func (c *UserController) UploadAvatar(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	file, header, err := ctx.Request.FormFile("avatar")
	if err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Missing avatar file", "MISSING_FILE")
		return
	}
	defer file.Close()

	result, err := cloudinary.Upload(file, header)
	if err != nil {
		dto.SendError(ctx, http.StatusInternalServerError, "Failed to upload image", "UPLOAD_FAILED")
		return
	}

	if err := c.service.UpdateAvatar(authUser.(auth.AuthUser).ID, result.SecureURL, result.PublicID); err != nil {
		dto.SendError(ctx, http.StatusInternalServerError, "Failed to update avatar", "DB_UPDATE_FAILED")
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Avatar updated successfully", gin.H{"avatar_url": result.SecureURL})
}

// UploadCover handles cover image uploads.
func (c *UserController) UploadCover(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	file, header, err := ctx.Request.FormFile("cover")
	if err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Missing cover file", "MISSING_FILE")
		return
	}
	defer file.Close()

	result, err := cloudinary.Upload(file, header)
	if err != nil {
		dto.SendError(ctx, http.StatusInternalServerError, "Failed to upload image", "UPLOAD_FAILED")
		return
	}

	if err := c.service.UpdateCover(authUser.(auth.AuthUser).ID, result.SecureURL, result.PublicID); err != nil {
		dto.SendError(ctx, http.StatusInternalServerError, "Failed to update cover", "DB_UPDATE_FAILED")
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Cover updated successfully", gin.H{"cover_url": result.SecureURL})
}

func (c *UserController) ChangePassword(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	var req dto.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	err := c.service.ChangePassword(authUser.(auth.AuthUser).ID, req.OldPassword, req.NewPassword)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Password changed successfully", nil)
}

// --- Admin-only actions ---

func (c *UserController) DeleteUser(ctx *gin.Context) {
	userID := ctx.Param("id")
	err := c.service.DeleteUser(userID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "User deleted successfully", gin.H{"id": userID})
}
