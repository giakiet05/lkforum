package controller

import (
	"net/http"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/auth"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostController struct {
	service service.PostService
}

func NewPostController(service service.PostService) *PostController {
	return &PostController{service: service}
}

// === CRUD Operations ===

func (c *PostController) CreatePost(ctx *gin.Context) {
	var req dto.CreatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Invalid request payload", apperror.ErrBadRequest.Code)
		return
	}

	userID, err := c.getAuthUserID(ctx)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	post, err := c.service.CreatePost(ctx.Request.Context(), userID, &req)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusCreated, "Post created successfully", post)
}

func (c *PostController) GetPostByID(ctx *gin.Context) {
	postID, err := c.parseObjectID(ctx.Param("id"))
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	// For GET requests, a missing auth user is not a hard error.
	// We pass a NilObjectID to the service to indicate an anonymous user.
	userID, err := c.getAuthUserID(ctx)
	if err != nil {
		userID = primitive.NilObjectID
	}

	post, err := c.service.GetPostByID(ctx.Request.Context(), postID, userID)
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

	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 20
	} else if query.Limit > 100 {
		query.Limit = 100
	}

	userID, err := c.getAuthUserID(ctx)
	if err != nil {
		userID = primitive.NilObjectID
	}

	posts, err := c.service.GetPosts(ctx.Request.Context(), userID, &query)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Posts retrieved successfully", posts)
}

func (c *PostController) UpdatePost(ctx *gin.Context) {
	postID, err := c.parseObjectID(ctx.Param("id"))
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	userID, err := c.getAuthUserID(ctx)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	var req dto.UpdatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Invalid request payload", apperror.ErrBadRequest.Code)
		return
	}

	updatedPost, err := c.service.UpdatePost(ctx.Request.Context(), postID, userID, &req)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Post updated successfully", updatedPost)
}

func (c *PostController) DeletePost(ctx *gin.Context) {
	postID, err := c.parseObjectID(ctx.Param("id"))
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	userID, err := c.getAuthUserID(ctx)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	err = c.service.DeletePost(ctx.Request.Context(), postID, userID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Post deleted successfully", gin.H{"id": postID.Hex()})
}

// === Interactions ===

func (c *PostController) VoteOnPost(ctx *gin.Context) {
	postID, err := c.parseObjectID(ctx.Param("id"))
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	userID, err := c.getAuthUserID(ctx)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	var req struct {
		Value *bool `json:"value" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Invalid request: 'value' field is required", apperror.ErrBadRequest.Code)
		return
	}

	votesCount, err := c.service.VoteOnPost(ctx.Request.Context(), userID, postID, *req.Value)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Vote cast successfully", votesCount)
}

func (c *PostController) VoteOnPoll(ctx *gin.Context) {
	postID, err := c.parseObjectID(ctx.Param("id"))
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	userID, err := c.getAuthUserID(ctx)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	var req struct {
		OptionID string `json:"option_id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Invalid request: 'option_id' is required", apperror.ErrBadRequest.Code)
		return
	}

	optionID, err := c.parseObjectID(req.OptionID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	poll, err := c.service.VoteOnPoll(ctx.Request.Context(), userID, postID, optionID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Poll vote cast successfully", poll)
}

// === Image Management ===

func (c *PostController) AddImagesToPost(ctx *gin.Context) {
	postID, err := c.parseObjectID(ctx.Param("id"))
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	userID, err := c.getAuthUserID(ctx)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	var req dto.AddImageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Invalid image data", apperror.ErrBadRequest.Code)
		return
	}

	images, err := c.service.AddImagesToPost(ctx.Request.Context(), userID, postID, &req)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Images added successfully", images)
}

func (c *PostController) RemoveImagesFromPost(ctx *gin.Context) {
	postID, err := c.parseObjectID(ctx.Param("id"))
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	userID, err := c.getAuthUserID(ctx)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	var req dto.RemoveImageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Image IDs are required", apperror.ErrBadRequest.Code)
		return
	}

	if err := c.service.RemoveImagesFromPost(ctx.Request.Context(), userID, postID, &req); err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Images removed successfully", nil)
}

// === Poll Management ===

func (c *PostController) UpdatePollDetails(ctx *gin.Context) {
	postID, err := c.parseObjectID(ctx.Param("id"))
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	userID, err := c.getAuthUserID(ctx)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	var req dto.UpdatePollRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Invalid poll data", apperror.ErrBadRequest.Code)
		return
	}

	poll, err := c.service.UpdatePollDetails(ctx.Request.Context(), postID, userID, &req)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Poll details updated successfully", poll)
}

func (c *PostController) AddPollOptions(ctx *gin.Context) {
	postID, err := c.parseObjectID(ctx.Param("id"))
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	userID, err := c.getAuthUserID(ctx)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	var req dto.AddPollOptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Invalid poll options", apperror.ErrBadRequest.Code)
		return
	}

	poll, err := c.service.AddPollOptions(ctx.Request.Context(), userID, postID, &req)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Poll options added successfully", poll)
}

func (c *PostController) UpdatePollOption(ctx *gin.Context) {
	userID, err := c.getAuthUserID(ctx)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	postID, err := c.parseObjectID(ctx.Param("id"))
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	optionID, err := c.parseObjectID(ctx.Param("optionID"))
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	var req dto.UpdatePollOptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Invalid option text", apperror.ErrBadRequest.Code)
		return
	}

	pollResponse, err := c.service.UpdatePollOption(ctx.Request.Context(), userID, postID, optionID, req.Text)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Poll option updated successfully", pollResponse)
}

func (c *PostController) RemovePollOptions(ctx *gin.Context) {
	postID, err := c.parseObjectID(ctx.Param("id"))
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	userID, err := c.getAuthUserID(ctx)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	var req dto.RemovePollOptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Option IDs are required", apperror.ErrBadRequest.Code)
		return
	}

	poll, err := c.service.RemovePollOptions(ctx.Request.Context(), userID, postID, &req)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Poll options removed successfully", poll)
}

// === Helpers ===

// getAuthUserID is a helper to get the authenticated user's ObjectID.
// It returns an error if the user is not authenticated or the ID is invalid.
func (c *PostController) getAuthUserID(ctx *gin.Context) (primitive.ObjectID, error) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		return primitive.NilObjectID, apperror.ErrForbidden // Or a more specific "UNAUTHORIZED" error
	}

	user, ok := authUser.(auth.AuthUser)
	if !ok {
		return primitive.NilObjectID, apperror.ErrInternal
	}

	userID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return primitive.NilObjectID, apperror.ErrInvalidToken // Invalid ID in token
	}

	return userID, nil
}

// parseObjectID is a helper to parse an ObjectID from a string.
func (c *PostController) parseObjectID(idStr string) (primitive.ObjectID, error) {
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return primitive.NilObjectID, apperror.ErrInvalidID
	}
	return id, nil
}
