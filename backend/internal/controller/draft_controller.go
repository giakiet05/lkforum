package controller

import (
	"net/http"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/auth"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/service"
	"github.com/gin-gonic/gin"
)

type DraftController struct {
	service service.DraftService
}

func NewDraftController(service service.DraftService) *DraftController {
	return &DraftController{service: service}
}

func (c *DraftController) CreateDraft(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	var req dto.CreateDraftRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Invalid request payload", apperror.ErrBadRequest.Code)
		return
	}

	draft, err := c.service.CreateDraft(ctx, authUser.(auth.AuthUser).ID, &req)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusCreated, "Draft created successfully", draft)
}

func (c *DraftController) UpdateDraft(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	draftID := ctx.Param("id")

	var req dto.UpdateDraftRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Invalid request payload", apperror.ErrBadRequest.Code)
		return
	}

	draft, err := c.service.UpdateDraft(ctx, authUser.(auth.AuthUser).ID, draftID, &req)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Draft updated successfully", draft)
}

func (c *DraftController) GetDraftByID(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	draftID := ctx.Param("id")

	draft, err := c.service.GetDraftByID(ctx, authUser.(auth.AuthUser).ID, draftID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Draft retrieved successfully", draft)
}

func (c *DraftController) GetDrafts(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	var query dto.GetPostsQuery // Reusing GetPostsQuery for pagination
	if err := ctx.ShouldBindQuery(&query); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Invalid query parameters", apperror.ErrBadRequest.Code)
		return
	}

	drafts, total, err := c.service.GetDraftsByAuthor(ctx, authUser.(auth.AuthUser).ID, &query)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}
	
	limit := query.Limit
	if limit <= 0 {
		limit = 10
	}

	dto.SendSuccess(ctx, http.StatusOK, "Drafts retrieved successfully", &dto.PaginatedDraftsResponse{
		Drafts: drafts,
		Pagination: dto.Pagination{
			Page:     query.Page,
			PageSize: limit,
			Total:    total,
		},
	})
}

func (c *DraftController) DeleteDraft(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	draftID := ctx.Param("id")

	err := c.service.DeleteDraft(ctx, authUser.(auth.AuthUser).ID, draftID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Draft deleted successfully", nil)
}

func (c *DraftController) PublishDraft(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	draftID := ctx.Param("id")

	post, err := c.service.PublishDraft(ctx, authUser.(auth.AuthUser).ID, draftID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusCreated, "Draft published successfully", post)
}

