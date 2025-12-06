package controller

import (
	"net/http"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/service"
	"github.com/gin-gonic/gin"
)

type AdminCommunityController struct {
	adminCommunityService service.AdminCommunityService
}

func NewAdminCommunityController(adminCommunityService service.AdminCommunityService) *AdminCommunityController {
	return &AdminCommunityController{
		adminCommunityService: adminCommunityService,
	}
}

// GetCommunitiesAdmin gets all communities with admin filters
func (c *AdminCommunityController) GetCommunitiesAdmin(ctx *gin.Context) {
	var query dto.GetCommunitiesAdminQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Invalid query parameters", apperror.ErrBadRequest.Code)
		return
	}

	result, err := c.adminCommunityService.GetCommunitiesAdmin(&query)
	if err != nil {
		dto.SendError(ctx, http.StatusInternalServerError, "Failed to get communities", apperror.ErrInternal.Code)
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Communities retrieved successfully", result)
}

// BanCommunity bans a community
func (c *AdminCommunityController) BanCommunity(ctx *gin.Context) {
	communityID := ctx.Param("id")
	if communityID == "" {
		dto.SendError(ctx, http.StatusBadRequest, "Community ID is required", apperror.ErrBadRequest.Code)
		return
	}

	var req dto.AdminBanCommunityRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Invalid request body", apperror.ErrBadRequest.Code)
		return
	}

	err := c.adminCommunityService.BanCommunity(communityID, &req)
	if err != nil {
		if err == apperror.ErrCommunityNotFound {
			dto.SendError(ctx, http.StatusNotFound, "Community not found", apperror.ErrCommunityNotFound.Code)
			return
		}
		dto.SendError(ctx, http.StatusInternalServerError, "Failed to ban community", apperror.ErrInternal.Code)
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Community banned successfully", nil)
}

// UnbanCommunity unbans a community
func (c *AdminCommunityController) UnbanCommunity(ctx *gin.Context) {
	communityID := ctx.Param("id")
	if communityID == "" {
		dto.SendError(ctx, http.StatusBadRequest, "Community ID is required", apperror.ErrBadRequest.Code)
		return
	}

	err := c.adminCommunityService.UnbanCommunity(communityID)
	if err != nil {
		if err == apperror.ErrCommunityNotFound {
			dto.SendError(ctx, http.StatusNotFound, "Community not found", apperror.ErrCommunityNotFound.Code)
			return
		}
		dto.SendError(ctx, http.StatusInternalServerError, "Failed to unban community", apperror.ErrInternal.Code)
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Community unbanned successfully", nil)
}