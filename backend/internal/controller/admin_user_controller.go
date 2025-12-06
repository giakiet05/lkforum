package controller

import (
	"net/http"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/service"
	"github.com/gin-gonic/gin"
)

type AdminUserController struct {
	adminService service.AdminUserService
}

func NewAdminUserController(adminService service.AdminUserService) *AdminUserController {
	return &AdminUserController{
		adminService: adminService,
	}
}

// GetUsers gets all users with admin filters
func (c *AdminUserController) GetUsers(ctx *gin.Context) {
	var query dto.GetUsersAdminQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Invalid query parameters", apperror.ErrBadRequest.Code)
		return
	}

	users, err := c.adminService.GetUsersAdmin(&query)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Users retrieved successfully", users)
}

// BanUser bans a user
func (c *AdminUserController) BanUser(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	if userID == "" {
		dto.SendError(ctx, http.StatusBadRequest, "User ID is required", apperror.ErrBadRequest.Code)
		return
	}

	var req dto.BanUserRequest
	if err := ctx.ShouldBind(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	err := c.adminService.BanUser(userID, &req)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	banType := "permanently"
	if req.BanUntil != nil {
		banType = "until " + req.BanUntil.Format("2006-01-02")
	}

	dto.SendSuccess(ctx, http.StatusOK, "User banned "+banType, gin.H{"user_id": userID})
}

// UnbanUser unbans a user
func (c *AdminUserController) UnbanUser(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	if userID == "" {
		dto.SendError(ctx, http.StatusBadRequest, "User ID is required", apperror.ErrBadRequest.Code)
		return
	}

	err := c.adminService.UnbanUser(userID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "User unbanned successfully", gin.H{"user_id": userID})
}

// DeleteUser soft deletes a user
func (c *AdminUserController) DeleteUser(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	if userID == "" {
		dto.SendError(ctx, http.StatusBadRequest, "User ID is required", apperror.ErrBadRequest.Code)
		return
	}

	err := c.adminService.SoftDeleteUser(userID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "User deleted successfully", gin.H{"user_id": userID})
}

// RestoreUser restores a soft-deleted user
func (c *AdminUserController) RestoreUser(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	if userID == "" {
		dto.SendError(ctx, http.StatusBadRequest, "User ID is required", apperror.ErrBadRequest.Code)
		return
	}

	err := c.adminService.RestoreUser(userID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "User restored successfully", gin.H{"user_id": userID})
}
