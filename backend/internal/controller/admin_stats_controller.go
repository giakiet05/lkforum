package controller

import (
	"net/http"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/service"
	"github.com/gin-gonic/gin"
)

type AdminStatsController struct {
	adminStatsService service.AdminStatsService
}

func NewAdminStatsController(adminStatsService service.AdminStatsService) *AdminStatsController {
	return &AdminStatsController{
		adminStatsService: adminStatsService,
	}
}

// GetPlatformOverview gets the main admin dashboard statistics
func (c *AdminStatsController) GetPlatformOverview(ctx *gin.Context) {
	result, err := c.adminStatsService.GetPlatformOverview()
	if err != nil {
		dto.SendError(ctx, http.StatusInternalServerError, "Failed to get platform overview", apperror.ErrInternal.Code)
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Platform overview retrieved successfully", result)
}

// GetUserStats gets detailed user statistics
func (c *AdminStatsController) GetUserStats(ctx *gin.Context) {
	var query dto.GetUserStatsQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Invalid query parameters", apperror.ErrBadRequest.Code)
		return
	}

	// Set default period if not provided
	if query.Period == "" {
		query.Period = "week"
	}

	result, err := c.adminStatsService.GetUserStats(&query)
	if err != nil {
		dto.SendError(ctx, http.StatusInternalServerError, "Failed to get user statistics", apperror.ErrInternal.Code)
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "User statistics retrieved successfully", result)
}

// GetContentStats gets detailed content statistics
func (c *AdminStatsController) GetContentStats(ctx *gin.Context) {
	var query dto.GetContentStatsQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Invalid query parameters", apperror.ErrBadRequest.Code)
		return
	}

	// Set default period if not provided
	if query.Period == "" {
		query.Period = "week"
	}

	result, err := c.adminStatsService.GetContentStats(&query)
	if err != nil {
		dto.SendError(ctx, http.StatusInternalServerError, "Failed to get content statistics", apperror.ErrInternal.Code)
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Content statistics retrieved successfully", result)
}