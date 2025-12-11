package controller

import (
	"net/http"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/auth"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/service"
	"github.com/gin-gonic/gin"
)

type PostController struct {
	service service.PostService
}

func NewPostController(service service.PostService) *PostController {
	return &PostController{service: service}
}

func (c *PostController) CreatePost(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	var req dto.CreatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Invalid request payload", apperror.ErrBadRequest.Code)
		return
	}

	post, err := c.service.CreatePost(authUser.(auth.AuthUser).ID, &req)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusCreated, "Post created successfully", post)
}

func (c *PostController) GetPostByID(ctx *gin.Context) {
	postID := ctx.Param("id")

	var userID string
	if val, exists := ctx.Get("authUser"); exists {
		userID = val.(auth.AuthUser).ID
	}

	post, err := c.service.GetPostByID(postID, userID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Post retrieved successfully", post)
}

func (c *PostController) GetPosts(ctx *gin.Context) {
	var query dto.GetPostsQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Invalid query parameters", apperror.ErrBadRequest.Code)
		return
	}

	var userID string
	if val, exists := ctx.Get("authUser"); exists {
		userID = val.(auth.AuthUser).ID
	}

	posts, err := c.service.GetPosts(userID, &query)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Posts retrieved successfully", posts)
}

func (c *PostController) GetMyPosts(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	var query dto.GetPostsQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Invalid query parameters", apperror.ErrBadRequest.Code)
		return
	}

	posts, err := c.service.GetMyPosts(authUser.(auth.AuthUser).ID, &query)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "My posts retrieved successfully", posts)
}

func (c *PostController) UpdatePost(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	postID := ctx.Param("id")

	var req dto.UpdatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Invalid request payload", apperror.ErrBadRequest.Code)
		return
	}

	updatedPost, err := c.service.UpdatePost(postID, authUser.(auth.AuthUser).ID, &req)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Post updated successfully", updatedPost)
}

func (c *PostController) DeletePost(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	postID := ctx.Param("id")

	err := c.service.DeletePost(postID, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Post deleted successfully", gin.H{"id": postID})
}

func (c *PostController) AddImagesToPost(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	postID := ctx.Param("id")

	form, err := ctx.MultipartForm()
	if err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Invalid form data", apperror.ErrBadRequest.Code)
		return
	}

	images, err := c.service.AddImagesToPost(authUser.(auth.AuthUser).ID, postID, form)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Images added successfully", images)
}

func (c *PostController) RemoveImagesFromPost(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	postID := ctx.Param("id")

	var req dto.RemoveImagesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Invalid request: 'public_ids' are required", apperror.ErrBadRequest.Code)
		return
	}

	if err := c.service.RemoveImagesFromPost(authUser.(auth.AuthUser).ID, postID, req.PublicIDs); err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Images removed successfully", nil)
}

func (c *PostController) AddVideosToPost(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	postID := ctx.Param("id")

	form, err := ctx.MultipartForm()
	if err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Invalid form data", apperror.ErrBadRequest.Code)
		return
	}

	videos, err := c.service.AddVideosToPost(authUser.(auth.AuthUser).ID, postID, form)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Videos added successfully", videos)
}

func (c *PostController) RemoveVideosFromPost(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	postID := ctx.Param("id")

	var req dto.RemoveImagesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Invalid request: 'public_ids' are required", apperror.ErrBadRequest.Code)
		return
	}

	if err := c.service.RemoveVideosFromPost(authUser.(auth.AuthUser).ID, postID, req.PublicIDs); err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Videos removed successfully", nil)
}

// VoteOnPost method removed - use VoteController instead
// POST /api/votes/post/:id for voting on posts

func (c *PostController) SavePost(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	postID := ctx.Param("id")

	err := c.service.SavePost(authUser.(auth.AuthUser).ID, postID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Post saved successfully", nil)
}

func (c *PostController) UnsavePost(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	postID := ctx.Param("id")

	err := c.service.UnsavePost(authUser.(auth.AuthUser).ID, postID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Post unsaved successfully", nil)
}

func (c *PostController) GetSavedPosts(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	var query dto.GetPostsQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Invalid query parameters", apperror.ErrBadRequest.Code)
		return
	}

	posts, err := c.service.GetSavedPosts(authUser.(auth.AuthUser).ID, &query)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Saved posts retrieved successfully", posts)
}

func (c *PostController) ReportPost(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	postID := ctx.Param("id")

	var req dto.ReportPostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Invalid request payload", apperror.ErrBadRequest.Code)
		return
	}

	err := c.service.ReportPost(authUser.(auth.AuthUser).ID, postID, req.Reason, req.Description)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Post reported successfully", nil)
}

func (c *PostController) HidePost(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	postID := ctx.Param("id")

	err := c.service.HidePost(authUser.(auth.AuthUser).ID, postID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Post hidden successfully", nil)
}

func (c *PostController) UnhidePost(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	postID := ctx.Param("id")

	err := c.service.UnhidePost(authUser.(auth.AuthUser).ID, postID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Post unhidden successfully", nil)
}

func (c *PostController) GetHiddenPosts(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	var query dto.GetPostsQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Invalid query parameters", apperror.ErrBadRequest.Code)
		return
	}

	posts, err := c.service.GetHiddenPosts(authUser.(auth.AuthUser).ID, &query)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Hidden posts retrieved successfully", posts)
}

func (c *PostController) VoteOnPoll(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	postID := ctx.Param("id")

	var req dto.PollVoteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Invalid request: 'option_id' is required", apperror.ErrBadRequest.Code)
		return
	}

	poll, err := c.service.VoteOnPoll(authUser.(auth.AuthUser).ID, postID, req.OptionID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Poll vote cast successfully", poll)
}

func (c *PostController) RemovePollVote(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	postID := ctx.Param("id")

	poll, err := c.service.RemovePollVote(authUser.(auth.AuthUser).ID, postID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Poll vote removed successfully", poll)
}

func (c *PostController) UpdatePoll(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	postID := ctx.Param("id")

	var req dto.UpdatePollRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Invalid request payload", apperror.ErrBadRequest.Code)
		return
	}

	poll, err := c.service.UpdatePoll(authUser.(auth.AuthUser).ID, postID, &req)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Poll updated successfully", poll)
}

func (c *PostController) AddPollOptions(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	postID := ctx.Param("id")

	var req dto.AddPollOptionsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Invalid request: 'options' are required", apperror.ErrBadRequest.Code)
		return
	}

	poll, err := c.service.AddPollOptions(authUser.(auth.AuthUser).ID, postID, req.Options)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Poll options added successfully", poll)
}

func (c *PostController) RemovePollOptions(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	postID := ctx.Param("id")

	var req dto.RemovePollOptionsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Invalid request: 'option_ids' are required", apperror.ErrBadRequest.Code)
		return
	}

	poll, err := c.service.RemovePollOptions(authUser.(auth.AuthUser).ID, postID, req.OptionIDs)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Poll options removed successfully", poll)
}

func (c *PostController) UpdatePollOption(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	postID := ctx.Param("id")
	optionID := ctx.Param("optionID")

	var req dto.UpdatePollOptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Invalid request payload", apperror.ErrBadRequest.Code)
		return
	}

	poll, err := c.service.UpdatePollOption(authUser.(auth.AuthUser).ID, postID, optionID, &req)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Poll option updated successfully", poll)
}

func (c *PostController) GetBanPosts(ctx *gin.Context) {
	var query dto.GetBanPostsQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Invalid query parameters", apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	responses, err := c.service.GetBanPosts(&query, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Banned posts retrieved successfully", responses)
}
