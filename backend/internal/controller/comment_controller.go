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
		ctx.JSON(apperror.StatusFromError(apperror.ErrBadRequest), dto.ErrorResponse{ErrorCode: apperror.ErrBadRequest.Code, Message: apperror.Message(err)})
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		ctx.JSON(apperror.StatusFromError(apperror.ErrForbidden), dto.ErrorResponse{ErrorCode: apperror.ErrForbidden.Code, Message: apperror.ErrForbidden.Message})
		return
	}

	comment, err := c.commentService.CreateComment(req, authUser.(auth.AuthUser).ID)
	if err != nil {
		ctx.JSON(apperror.StatusFromError(err), dto.ErrorResponse{ErrorCode: apperror.Code(err), Message: apperror.Message(err)})
		return
	}

	ctx.JSON(http.StatusCreated, dto.SuccessResponse{
		ID:      comment.ID.Hex(),
		Message: "Create comment successfully",
	})
}

func (c *CommentController) GetCommentByID(ctx *gin.Context) {
	commentID := ctx.Param("comment_id")
	if commentID == "" {
		ctx.JSON(apperror.StatusFromError(apperror.ErrBadRequest), dto.ErrorResponse{ErrorCode: apperror.ErrBadRequest.Code, Message: apperror.ErrBadRequest.Message})
		return
	}

	comment, err := c.commentService.GetCommentByID(commentID)
	if err != nil {
		ctx.JSON(apperror.StatusFromError(err), dto.ErrorResponse{ErrorCode: apperror.Code(err), Message: apperror.Message(err)})
		return
	}

	ctx.JSON(http.StatusOK, dto.FromComment(comment))
}

func (c *CommentController) GetCommentByPostID(ctx *gin.Context) {
	var query *dto.GetCommentByPostIDQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(apperror.StatusFromError(apperror.ErrBadRequest), dto.ErrorResponse{ErrorCode: apperror.ErrBadRequest.Code, Message: apperror.Message(err)})
		return
	}

	response, err := c.commentService.GetCommentByPostIDPaginated(query)
	if err != nil {
		ctx.JSON(apperror.StatusFromError(err), dto.ErrorResponse{ErrorCode: apperror.Code(err), Message: apperror.Message(err)})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *CommentController) GetCommentsFilter(ctx *gin.Context) {
	var query *dto.GetCommentsFilterQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(apperror.StatusFromError(apperror.ErrBadRequest), dto.ErrorResponse{ErrorCode: apperror.ErrBadRequest.Code, Message: apperror.Message(err)})
		return
	}

	response, err := c.commentService.GetCommentsFilterPaginated(query)
	if err != nil {
		ctx.JSON(apperror.StatusFromError(err), dto.ErrorResponse{ErrorCode: apperror.Code(err), Message: apperror.Message(err)})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *CommentController) DeleteCommentByID(ctx *gin.Context) {
	commentID := ctx.Param("comment_id")
	if commentID == "" {
		ctx.JSON(apperror.StatusFromError(apperror.ErrBadRequest), dto.ErrorResponse{ErrorCode: apperror.ErrBadRequest.Code, Message: apperror.ErrBadRequest.Message})
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		ctx.JSON(apperror.StatusFromError(apperror.ErrForbidden), dto.ErrorResponse{ErrorCode: apperror.ErrForbidden.Code, Message: apperror.ErrForbidden.Message})
		return
	}

	err := c.commentService.DeleteCommentByID(commentID, authUser.(auth.AuthUser).ID)
	if err != nil {
		ctx.JSON(apperror.StatusFromError(err), dto.ErrorResponse{ErrorCode: apperror.Code(err), Message: apperror.Message(err)})
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		ID:      commentID,
		Message: "Delete comment successfully",
	})
}
