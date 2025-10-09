package controller

import (
	"errors"
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
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: "INVALID_REQUEST", Error: err.Error()})
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		ctx.JSON(http.StatusForbidden, dto.ErrorResponse{Code: "FORBIDDEN", Error: "Invalid jwt token"})
		return
	}

	community, err := c.communityService.CreateCommunity(&req, authUser.(auth.AuthUser).ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: apperror.Code(err), Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, dto.SuccessResponse{
		ID:      community.ID.Hex(),
		Message: "Create community successfully",
	})
}

func (c *CommunityController) GetCommunityByID(ctx *gin.Context) {
	communityID := ctx.Param("community_id")
	if communityID == "" {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: "INVALID_REQUEST", Error: communityID})
		return
	}

	community, err := c.communityService.GetCommunityByID(communityID)
	if err != nil {
		if errors.Is(err, apperror.ErrCommunityNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{Code: apperror.Code(err), Error: err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: apperror.Code(err), Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.FromCommunity(community))
}

func (c *CommunityController) GetCommunitiesFilter(ctx *gin.Context) {
	name := ctx.Query("name")
	description := ctx.Query("description")
	createFromStr := ctx.Query("create_from")
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

	var createFrom time.Time
	if createFromStr != "" {
		t, err := time.Parse(time.RFC3339, createFromStr)
		if err == nil {
			createFrom = t
		} else {
			ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: apperror.Code(err), Error: err.Error()})
			return
		}
	}

	response, err := c.communityService.GetCommunitiesFilter(name, description, createFrom, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: apperror.Code(err), Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *CommunityController) GetCommunityByModeratorID(ctx *gin.Context) {
	moderatorID := ctx.Param("moderator_id")
	if moderatorID == "" {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: "INVALID_REQUEST", Error: moderatorID})
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
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: apperror.Code(err), Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
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
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: apperror.Code(err), Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *CommunityController) UpdateCommunity(ctx *gin.Context) {
	var req dto.UpdateCommunityRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: "INVALID_REQUEST", Error: err.Error()})
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		ctx.JSON(http.StatusForbidden, dto.ErrorResponse{Code: "FORBIDDEN", Error: "Invalid jwt token"})
		return
	}

	community, err := c.communityService.UpdateCommunity(&req, authUser.(auth.AuthUser).ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: apperror.Code(err), Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		ID:      community.ID.Hex(),
		Message: "Update community successfully",
	})
}

func (c *CommunityController) AddModerator(ctx *gin.Context) {
	var req *dto.AddModeratorRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: "INVALID_REQUEST", Error: err.Error()})
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		ctx.JSON(http.StatusForbidden, dto.ErrorResponse{Code: "FORBIDDEN", Error: "Invalid jwt token"})
		return
	}

	err := c.communityService.AddModerator(req, authUser.(auth.AuthUser).ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: apperror.Code(err), Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		ID:      req.CommunityID,
		Message: "Add moderator successfully",
	})
}

func (c *CommunityController) RemoveModerator(ctx *gin.Context) {
	var req *dto.RemoveModeratorRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: "INVALID_REQUEST", Error: err.Error()})
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		ctx.JSON(http.StatusForbidden, dto.ErrorResponse{Code: "FORBIDDEN", Error: "Invalid jwt token"})
		return
	}

	err := c.communityService.RemoveModerator(req, authUser.(auth.AuthUser).ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: apperror.Code(err), Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		ID:      req.CommunityID,
		Message: "Remove moderator successfully",
	})
}

func (c *CommunityController) DeleteCommunityByID(ctx *gin.Context) {
	communityID := ctx.Param("community_id")
	if communityID == "" {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: "INVALID_REQUEST", Error: communityID})
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		ctx.JSON(http.StatusForbidden, dto.ErrorResponse{Code: "FORBIDDEN", Error: "Invalid jwt token"})
		return
	}

	err := c.communityService.DeleteCommunityByID(communityID, authUser.(auth.AuthUser).ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: apperror.Code(err), Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		ID:      communityID,
		Message: "Delete community successfully",
	})
}
