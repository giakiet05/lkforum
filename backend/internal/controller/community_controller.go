package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/auth"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/service"
	"github.com/gin-gonic/gin"
)

type CommunityController struct {
	communityService service.CommunityService
}

func NewCommunityController(communityService service.CommunityService) *CommunityController {
	return &CommunityController{communityService: communityService}
}

func (c *CommunityController) CreateCommunity(ctx *gin.Context) {
	var req dto.CreateCommunityRequest
	if err := ctx.ShouldBind(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	community, err := c.communityService.CreateCommunity(&req, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusCreated, "Community created successfully", dto.FromCommunity(community))
}

func (c *CommunityController) GetCommunityByID(ctx *gin.Context) {
	communityID := ctx.Param("community_id")
	if communityID == "" {
		dto.SendError(ctx, http.StatusBadRequest, "Community ID is required", apperror.ErrBadRequest.Code)
		return
	}

	community, err := c.communityService.GetCommunityByID(communityID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Community retrieved successfully", dto.FromCommunity(community))
}

func (c *CommunityController) GetCommunitiesFilter(ctx *gin.Context) {
	name := ctx.Query("name")
	description := ctx.Query("description")
	is18PlusStr := ctx.Query("is_18_plus")
	createFromStr := ctx.Query("create_from")
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")

	is18Plus, err := strconv.ParseBool(is18PlusStr)
	if err != nil {
		is18Plus = false
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	var createFrom time.Time
	if createFromStr != "" {
		t, err := time.Parse(time.RFC3339, createFromStr)
		if err != nil {
			dto.SendError(ctx, http.StatusBadRequest, "Invalid date format for create_from", apperror.ErrBadRequest.Code)
			return
		}
		createFrom = t
	}

	response, err := c.communityService.GetCommunitiesFilter(name, description, is18Plus, createFrom, page, pageSize)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Communities retrieved successfully", response)
}

func (c *CommunityController) GetCommunityByModeratorID(ctx *gin.Context) {
	moderatorID := ctx.Param("moderator_id")
	if moderatorID == "" {
		dto.SendError(ctx, http.StatusBadRequest, "Moderator ID is required", apperror.ErrBadRequest.Code)
		return
	}

	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	response, err := c.communityService.GetCommunitiesByModeratorIDPaginated(moderatorID, page, pageSize)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Communities retrieved successfully", response)
}

func (c *CommunityController) GetAllCommunities(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	response, err := c.communityService.GetAllCommunitiesPaginated(page, pageSize)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Communities retrieved successfully", response)
}

func (c *CommunityController) UpdateCommunity(ctx *gin.Context) {
	var req dto.UpdateCommunityRequest
	if err := ctx.ShouldBind(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	community, err := c.communityService.UpdateCommunity(&req, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Community updated successfully", dto.FromCommunity(community))
}

func (c *CommunityController) AddModerator(ctx *gin.Context) {
	var req *dto.AddModeratorRequest
	if err := ctx.ShouldBind(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	err := c.communityService.AddModerator(req, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Moderator added successfully", gin.H{"community_id": req.CommunityID})
}

func (c *CommunityController) RemoveModerator(ctx *gin.Context) {
	var req *dto.RemoveModeratorRequest
	if err := ctx.ShouldBind(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	err := c.communityService.RemoveModerator(req, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Moderator removed successfully", gin.H{"community_id": req.CommunityID})
}

func (c *CommunityController) DeleteCommunityByID(ctx *gin.Context) {
	communityID := ctx.Param("community_id")
	if communityID == "" {
		dto.SendError(ctx, http.StatusBadRequest, "Community ID is required", apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	err := c.communityService.DeleteCommunityByID(communityID, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Community deleted successfully", gin.H{"id": communityID})
}

func (c *CommunityController) BanUser(ctx *gin.Context) {
	var req dto.BanUserRequest
	if err := ctx.ShouldBind(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	err := c.communityService.BanUser(&req, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Ban user successfully", gin.H{"id": req.UserID})
}

func (c *CommunityController) GetBanUsers(ctx *gin.Context) {
	communityID := ctx.Param("community_id")
	banType := ctx.Param("ban_type")

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	users, err := c.communityService.GetBannedUsers(communityID, banType, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}
	userResponses := dto.FromUsers(users)

	dto.SendSuccess(ctx, http.StatusOK, "Banned user retrieved successfully", userResponses)
}

func (c *CommunityController) UnbanUser(ctx *gin.Context) {
	var req dto.UnbanUserRequest
	if err := ctx.ShouldBind(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	err := c.communityService.UnbanUser(req.UserID, req.CommunityID, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Unbanned user successfully", gin.H{"id": req.UserID})
}

func (c *CommunityController) UnmuteUser(ctx *gin.Context) {
	var req dto.UnbanUserRequest
	if err := ctx.ShouldBind(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	err := c.communityService.UnmuteUser(req.UserID, req.CommunityID, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Unmuted user successfully", gin.H{"id": req.UserID})
}
