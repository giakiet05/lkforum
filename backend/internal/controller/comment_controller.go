package controller

import (
	"net/http"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/auth"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/service"
	"github.com/gin-gonic/gin"
)

type CommentController struct {
	commentService service.CommentService
}

func NewCommentController(commentService service.CommentService) *CommentController {
	return &CommentController{commentService: commentService}
}

func (c *CommentController) CreateComment(ctx *gin.Context) {
	var req *dto.CreateCommentRequest
	if err := ctx.ShouldBind(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	comment, err := c.commentService.CreateComment(req, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	userID := authUser.(auth.AuthUser).ID
	dto.SendSuccess(ctx, http.StatusCreated, "Comment created successfully", dto.FromComment(comment, &userID))
}

func (c *CommentController) GetCommentByID(ctx *gin.Context) {
	commentID := ctx.Param("comment_id")
	if commentID == "" {
		dto.SendError(ctx, http.StatusBadRequest, "Comment ID is required", apperror.ErrBadRequest.Code)
		return
	}

	comment, err := c.commentService.GetCommentByID(commentID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	var userID *string
	if val, exists := ctx.Get("authUser"); exists {
		uid := val.(auth.AuthUser).ID
		userID = &uid
	}

	dto.SendSuccess(ctx, http.StatusOK, "Comment retrieved successfully", dto.FromComment(comment, userID))
}

func (c *CommentController) GetCommentByPostID(ctx *gin.Context) {
	var query *dto.GetCommentByPostIDQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	var currentUserID *string
	if val, exists := ctx.Get("authUser"); exists {
		uid := val.(auth.AuthUser).ID
		currentUserID = &uid
	}

	response, err := c.commentService.GetCommentByPostIDPaginated(query, currentUserID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Comments retrieved successfully", response)
}

func (c *CommentController) GetCommentsFilter(ctx *gin.Context) {
	var query *dto.GetCommentsFilterQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	var currentUserID *string
	if val, exists := ctx.Get("authUser"); exists {
		uid := val.(auth.AuthUser).ID
		currentUserID = &uid
	}

	response, err := c.commentService.GetCommentsFilterPaginated(query, currentUserID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Comments retrieved successfully", response)
}

func (c *CommentController) DeleteCommentByID(ctx *gin.Context) {
	commentID := ctx.Param("comment_id")
	if commentID == "" {
		dto.SendError(ctx, http.StatusBadRequest, "Comment ID is required", apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	err := c.commentService.DeleteCommentByID(commentID, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Comment deleted successfully", gin.H{"id": commentID})
}
