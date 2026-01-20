package controller

import (
	"log"
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
		log.Printf("[CreateCommunity] ShouldBind error: %v", err)
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

	var userID *string
	authUser, exists := ctx.Get("authUser")
	if exists {
		userID = &authUser.(*auth.AuthUser).ID
	}

	community, err := c.communityService.GetCommunityByID(communityID, userID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Community retrieved successfully", dto.FromCommunity(community))
}

func (c *CommunityController) GetCommunityByName(ctx *gin.Context) {
	name := ctx.Param("name")
	if name == "" {
		dto.SendError(ctx, http.StatusBadRequest, "Community name is required", apperror.ErrBadRequest.Code)
		return
	}

	var userID *string
	authUser, exists := ctx.Get("authUser")
	if exists {
		userID = &authUser.(*auth.AuthUser).ID
	}

	community, err := c.communityService.GetCommunityByName(name, userID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Community retrieved successfully", dto.FromCommunity(community))
}

func (c *CommunityController) GetCommunitiesByUserID(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	if userID == "" {
		dto.SendError(ctx, http.StatusBadRequest, "User ID is required", apperror.ErrBadRequest.Code)
		return
	}

	communities, err := c.communityService.GetCommunitiesByUserID(userID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Communities retrieved successfully", communities)
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

	var userID *string
	authUser, exists := ctx.Get("authUser")
	if exists {
		userID = &authUser.(*auth.AuthUser).ID
	}

	responses, err := c.communityService.GetCommunitiesFilter(userID, name, description, is18Plus, createFrom, page, pageSize)
	if err != nil {
		log.Printf("ERROR: GetCommunitiesFilter failed: %v", err)
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Communities retrieved successfully", responses)
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

	var userID *string
	authUser, exists := ctx.Get("authUser")
	if exists {
		userID = &authUser.(*auth.AuthUser).ID
	}

	response, err := c.communityService.GetAllCommunitiesPaginated(userID, page, pageSize)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Communities retrieved successfully", response)
}

func (c *CommunityController) UpdateCommunity(ctx *gin.Context) {
	var req dto.UpdateCommunityRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.Printf("UpdateCommunity binding error: %v", err)
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

func (c *CommunityController) ActivateModerator(ctx *gin.Context) {
	communityID := ctx.Param("community_id")
	if communityID == "" {
		dto.SendError(ctx, http.StatusBadRequest, apperror.ErrBadRequest.Message, apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	err := c.communityService.ActivateModerator(communityID, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Moderator activated successfully", gin.H{"community_id": communityID})
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

func (c *CommunityController) BanPost(ctx *gin.Context) {
	var req dto.CommunityBanPostRequest
	if err := ctx.ShouldBind(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	userID := authUser.(auth.AuthUser).ID
	err := c.communityService.BanPost(&req, userID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Post banned successfully", gin.H{"id": req.PostID})
}

func (c *CommunityController) UnbanPost(ctx *gin.Context) {
	var req dto.CommunityUnbanPostRequest
	if err := ctx.ShouldBind(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	userID := authUser.(auth.AuthUser).ID
	err := c.communityService.UnbanPost(&req, userID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Post unbanned successfully", gin.H{"id": req.PostID})
}

func (c *CommunityController) BanUser(ctx *gin.Context) {
	var req dto.CommunityBanUserRequest
	if err := ctx.ShouldBind(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	userID := authUser.(auth.AuthUser).ID
	err := c.communityService.BanUser(&req, userID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Ban user successfully", gin.H{"id": userID})
}

func (c *CommunityController) GetBanUsers(ctx *gin.Context) {
	communityID := ctx.Query("community_id")
	banType := ctx.Query("ban_type")

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	expiredStr := ctx.Query("expired")
	expired := false

	if expiredStr != "" {
		parsed, err := strconv.ParseBool(expiredStr)
		if err != nil {
			dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
			return
		}
		expired = parsed
	}

	users, err := c.communityService.GetBannedUsers(communityID, banType, expired, authUser.(auth.AuthUser).ID)

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

func (c *CommunityController) GetPendingPosts(ctx *gin.Context) {
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

	response, err := c.communityService.GetPendingPosts(communityID, authUser.(auth.AuthUser).ID, page, pageSize)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Pending posts retrieved successfully", response)
}

func (c *CommunityController) GetEditedPosts(ctx *gin.Context) {
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

	response, err := c.communityService.GetEditedPosts(communityID, authUser.(auth.AuthUser).ID, page, pageSize)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Edited posts retrieved successfully", response)
}

func (c *CommunityController) ModeratePost(ctx *gin.Context) {
	communityID := ctx.Param("community_id")
	postID := ctx.Param("post_id")

	if communityID == "" || postID == "" {
		dto.SendError(ctx, http.StatusBadRequest, "Community ID and Post ID are required", apperror.ErrBadRequest.Code)
		return
	}

	var req dto.ModeratePostRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.Printf("❌ ModeratePost binding error: %v", err)
		log.Printf("   Request body: %+v", req)
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	log.Printf("✅ ModeratePost request parsed successfully: approve=%v, reason=%v", req.Approve, req.Reason)

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	err := c.communityService.ModeratePost(communityID, postID, authUser.(auth.AuthUser).ID, req.Approve, req.Reason)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	status := "approved"
	if !req.Approve {
		status = "rejected"
	}

	dto.SendSuccess(ctx, http.StatusOK, "Post "+status+" successfully", gin.H{"post_id": postID, "status": status})
}
