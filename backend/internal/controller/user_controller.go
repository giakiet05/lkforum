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

// UserController handles requests related to user management.
type UserController struct {
	service service.UserService
}

// NewUserController creates a new UserController.
func NewUserController(service service.UserService) *UserController {
	return &UserController{service: service}
}

func (c *UserController) GetUsers(ctx *gin.Context) {
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

	response, err := c.service.GetUsers(page, pageSize)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}
	dto.SendSuccess(ctx, http.StatusOK, "Users retrieved successfully", response)
}

func (c *UserController) UpdateUser(ctx *gin.Context) {
	userID := ctx.Param("id")
	if !auth.IsOwner(ctx, userID) && !auth.IsAdmin(ctx) {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	var req dto.UserUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	currentUser, err := c.service.GetUserByID(userID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	if req.Username != "" {
		currentUser.Username = req.Username
	}
	if req.Email != "" {
		// In a real app, changing email should trigger a re-verification process.
		currentUser.Email = req.Email
	}

	updatedUser, err := c.service.UpdateUser(currentUser)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "User updated successfully", dto.FromUser(updatedUser))
}

func (c *UserController) DeleteUser(ctx *gin.Context) {
	userID := ctx.Param("id")
	if !auth.IsOwner(ctx, userID) && !auth.IsAdmin(ctx) {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	err := c.service.DeleteUser(userID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "User deleted successfully", gin.H{"id": userID})
}

func (c *UserController) GetUserByID(ctx *gin.Context) {
	userID := ctx.Param("id")

	user, err := c.service.GetUserByID(userID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "User retrieved successfully", dto.FromUser(user))
}

func (c *UserController) GetUserByUsername(ctx *gin.Context) {
	username := ctx.Param("username")

	user, err := c.service.GetUserByUsername(username)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "User retrieved successfully", dto.FromUser(user))
}

func (c *UserController) ChangePassword(ctx *gin.Context) {
	userID := ctx.Param("id")
	authUser, exists := ctx.Get("authUser")
	if !exists || authUser.(auth.AuthUser).ID != userID {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	var req dto.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	err := c.service.ChangePassword(userID, req.OldPassword, req.NewPassword)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Password changed successfully", nil)
}
