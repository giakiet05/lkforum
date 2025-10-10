package controller

import (
	"errors"
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

// === CRUD Cơ bản ===

func (c *PostController) CreatePost(ctx *gin.Context) {
	var req dto.CreatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handlePostServiceError(ctx, apperror.NewError(err, "INVALID_REQUEST", "Invalid request payload"))
		return
	}

	authUser := c.getAuthUser(ctx)
	if authUser == nil {
		return
	}
	userID, ok := c.getAuthUserID(ctx)
	if !ok {
		return
	}
	post, err := c.service.CreatePost(ctx.Request.Context(), userID, &req)
	if err != nil {
		handlePostServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, post)
}

func (c *PostController) GetPostByID(ctx *gin.Context) {
	postID, ok := c.parseObjectID(ctx, ctx.Param("id"))
	if !ok {
		return
	}

	userID, ok := c.getAuthUserID(ctx)
	if !ok {
		return
	}

	post, err := c.service.GetPostByID(ctx.Request.Context(), postID, userID)
	if err != nil {
		handlePostServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, post)
}

func (c *PostController) GetPosts(ctx *gin.Context) {
	var query dto.GetPostsQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		handlePostServiceError(ctx, apperror.NewError(err, "INVALID_QUERY", "Invalid query parameters"))
		return
	}

	// Set defaults
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 20
	} else if query.Limit > 100 {
		query.Limit = 100
	}

	userID, ok := c.getAuthUserID(ctx)
	if !ok {
		return
	}

	posts, err := c.service.GetPosts(ctx.Request.Context(), userID, &query)
	if err != nil {
		handlePostServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, posts)
}

func (c *PostController) UpdatePost(ctx *gin.Context) {
	postID, ok := c.parseObjectID(ctx, ctx.Param("id"))
	if !ok {
		return
	}
	authUser := c.getAuthUser(ctx)
	if authUser == nil {
		return
	}

	var req dto.UpdatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handlePostServiceError(ctx, apperror.NewError(err, "INVALID_REQUEST", "Invalid request payload"))
		return
	}
	if req.Title == "" && req.Text == "" {
		handlePostServiceError(ctx, apperror.NewError(nil, "INVALID_REQUEST", "Invalid request payload"))
		return
	}
	userID, ok := c.getAuthUserID(ctx)
	if !ok {
		return
	}
	updatedPost, err := c.service.UpdatePost(ctx.Request.Context(), postID, userID, &req)
	if err != nil {
		handlePostServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, updatedPost)
}

func (c *PostController) DeletePost(ctx *gin.Context) {
	postID, ok := c.parseObjectID(ctx, ctx.Param("id"))
	if !ok {
		return
	}
	authUser := c.getAuthUser(ctx)
	if authUser == nil {
		return
	}
	userID, ok := c.getAuthUserID(ctx)
	if !ok {
		return
	}
	err := c.service.DeletePost(ctx.Request.Context(), postID, userID)
	if err != nil {
		handlePostServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse{Message: "Post deleted successfully"})
}

// === Tương tác ===

func (c *PostController) VoteOnPost(ctx *gin.Context) {
	postID, ok := c.parseObjectID(ctx, ctx.Param("id"))
	if !ok {
		return
	}
	authUser := c.getAuthUser(ctx)
	if authUser == nil {
		return
	}

	var req struct {
		Value bool `json:"value"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handlePostServiceError(ctx, apperror.NewError(err, "INVALID_REQUEST", "Vote value is required"))
		return
	}
	userID, ok := c.getAuthUserID(ctx)
	if !ok {
		return
	}
	votesCount, err := c.service.VoteOnPost(ctx.Request.Context(), userID, postID, req.Value)
	if err != nil {
		handlePostServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, votesCount)
}

func (c *PostController) VoteOnPoll(ctx *gin.Context) {
	postID, ok := c.parseObjectID(ctx, ctx.Param("id"))
	if !ok {
		return
	}
	authUser := c.getAuthUser(ctx)
	if authUser == nil {
		return
	}

	var req struct {
		OptionID string `json:"option_id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handlePostServiceError(ctx, apperror.NewError(err, "INVALID_REQUEST", "Option ID is required"))
		return
	}

	optionID, ok := c.parseObjectID(ctx, req.OptionID)
	if !ok {
		return
	}
	userID, ok := c.getAuthUserID(ctx)
	if !ok {
		return
	}
	poll, err := c.service.VoteOnPoll(ctx.Request.Context(), userID, postID, optionID)
	if err != nil {
		handlePostServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, poll)
}

// === Quản lý Image ===

func (c *PostController) AddImagesToPost(ctx *gin.Context) {
	postID, ok := c.parseObjectID(ctx, ctx.Param("id"))
	if !ok {
		return
	}
	authUser := c.getAuthUser(ctx)
	if authUser == nil {
		return
	}

	var req dto.AddImageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handlePostServiceError(ctx, apperror.NewError(err, "INVALID_REQUEST", "Invalid image data"))
		return
	}
	userID, ok := c.getAuthUserID(ctx)
	if !ok {
		return
	}
	images, err := c.service.AddImagesToPost(ctx.Request.Context(), userID, postID, &req)
	if err != nil {
		handlePostServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, images)
}

func (c *PostController) RemoveImagesFromPost(ctx *gin.Context) {
	postID, ok := c.parseObjectID(ctx, ctx.Param("id"))
	if !ok {
		return
	}
	authUser := c.getAuthUser(ctx)
	if authUser == nil {
		return
	}

	var req dto.RemoveImageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handlePostServiceError(ctx, apperror.NewError(err, "INVALID_REQUEST", "Image IDs are required"))
		return
	}
	userID, ok := c.getAuthUserID(ctx)
	if !ok {
		return
	}
	if err := c.service.RemoveImagesFromPost(ctx.Request.Context(), userID, postID, &req); err != nil {
		handlePostServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse{Message: "Images removed successfully"})
}

// === Quản lý Poll ===

func (c *PostController) UpdatePollDetails(ctx *gin.Context) {
	postID, ok := c.parseObjectID(ctx, ctx.Param("id"))
	if !ok {
		return
	}
	authUser := c.getAuthUser(ctx)
	if authUser == nil {
		return
	}

	var req dto.UpdatePollRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handlePostServiceError(ctx, apperror.NewError(err, "INVALID_REQUEST", "Invalid poll data"))
		return
	}
	userID, ok := c.getAuthUserID(ctx)
	if !ok {
		return
	}
	poll, err := c.service.UpdatePollDetails(ctx.Request.Context(), postID, userID, &req)
	if err != nil {
		handlePostServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, poll)
}

func (c *PostController) AddPollOptions(ctx *gin.Context) {
	postID, ok := c.parseObjectID(ctx, ctx.Param("id"))
	if !ok {
		return
	}
	authUser := c.getAuthUser(ctx)
	if authUser == nil {
		return
	}

	var req dto.AddPollOptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handlePostServiceError(ctx, apperror.NewError(err, "INVALID_REQUEST", "Invalid poll options"))
		return
	}
	userID, ok := c.getAuthUserID(ctx)
	if !ok {
		return
	}
	poll, err := c.service.AddPollOptions(ctx.Request.Context(), userID, postID, &req)
	if err != nil {
		handlePostServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, poll)
}

func (c *PostController) UpdatePollOption(ctx *gin.Context) {
	// 1. Lấy userID từ context
	authUser := c.getAuthUser(ctx)
	if authUser == nil {
		return
	}

	// 2. Lấy postID và optionID từ URL
	postID, ok := c.parseObjectID(ctx, ctx.Param("id"))
	if !ok {
		return
	}
	optionID, ok := c.parseObjectID(ctx, ctx.Param("optionID"))
	if !ok {
		return
	}

	// 3. Lấy newText từ request body (sử dụng DTO đã sửa)
	var req dto.UpdatePollOptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handlePostServiceError(ctx, apperror.NewError(err, "INVALID_REQUEST", "Invalid option text"))
		return
	}
	userID, ok := c.getAuthUserID(ctx)
	if !ok {
		return
	}
	// 4. Gọi đến service với các tham số đã được tách biệt rõ ràng
	pollResponse, err := c.service.UpdatePollOption(ctx.Request.Context(), userID, postID, optionID, req.Text)
	if err != nil {
		handlePostServiceError(ctx, err)
		return
	}

	// 5. Trả về kết quả thành công
	ctx.JSON(http.StatusOK, pollResponse)
}

func (c *PostController) RemovePollOptions(ctx *gin.Context) {
	postID, ok := c.parseObjectID(ctx, ctx.Param("id"))
	if !ok {
		return
	}
	authUser := c.getAuthUser(ctx)
	if authUser == nil {
		return
	}

	var req dto.RemovePollOptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handlePostServiceError(ctx, apperror.NewError(err, "INVALID_REQUEST", "Option IDs are required"))
		return
	}
	userID, ok := c.getAuthUserID(ctx)
	if !ok {
		return
	}
	poll, err := c.service.RemovePollOptions(ctx.Request.Context(), userID, postID, &req)
	if err != nil {
		handlePostServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, poll)
}

// === Helpers ===

// getAuthUser là một helper để lấy thông tin user đã xác thực.
func (c *PostController) getAuthUser(ctx *gin.Context) *auth.AuthUser {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		handlePostServiceError(ctx, apperror.NewError(nil, "UNAUTHORIZED", "Authentication required"))
		return nil
	}
	user, ok := authUser.(auth.AuthUser)
	if !ok {
		handlePostServiceError(ctx, apperror.NewError(nil, "INTERNAL_ERROR", "Invalid auth user type in context"))
		return nil
	}
	return &user
}

// parseObjectID là một helper để parse ObjectID từ string.
func (c *PostController) parseObjectID(ctx *gin.Context, idStr string) (primitive.ObjectID, bool) {
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		handlePostServiceError(ctx, apperror.NewError(err, "INVALID_ID", "Invalid ID format: "+idStr))
		return primitive.NilObjectID, false
	}
	return id, true
}

// handlePostServiceError là hàm helper cục bộ để xử lý lỗi từ PostService.
func handlePostServiceError(ctx *gin.Context, err error) {
	code := apperror.Code(err)
	status := http.StatusInternalServerError

	switch {
	case errors.Is(err, service.ErrPostNotFound):
		status = http.StatusNotFound
	case errors.Is(err, service.ErrPermissionDenied):
		status = http.StatusForbidden
	case errors.Is(err, service.ErrInvalidInput):
		status = http.StatusBadRequest
	case errors.Is(err, service.ErrPollCannotEdit):
		status = http.StatusConflict
	case code == "UNAUTHORIZED":
		status = http.StatusUnauthorized
	}

	ctx.JSON(status, dto.ErrorResponse{Code: code, Error: err.Error()})
}
func (c *PostController) getAuthUserID(ctx *gin.Context) (primitive.ObjectID, bool) {
	// Bước 1: Lấy thông tin user từ context
	authUser, exists := ctx.Get("authUser")
	if !exists {
		handlePostServiceError(ctx, apperror.NewError(nil, "UNAUTHORIZED", "Authentication required"))
		return primitive.NilObjectID, false
	}

	// Bước 2: Kiểm tra kiểu dữ liệu
	user, ok := authUser.(auth.AuthUser)
	if !ok {
		handlePostServiceError(ctx, apperror.NewError(nil, "INTERNAL_ERROR", "Invalid auth user type in context"))
		return primitive.NilObjectID, false
	}

	// Bước 3: Chuyển đổi ID từ string sang ObjectID
	userID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		handlePostServiceError(ctx, apperror.NewError(err, "INVALID_TOKEN", "Invalid user ID in token"))
		return primitive.NilObjectID, false
	}

	return userID, true
}
