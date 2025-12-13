package controller

import (
	"net/http"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/auth"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/service"
	"github.com/gin-gonic/gin"
)

type VoteController struct {
	voteService service.VoteService
}

func NewVoteController(voteService service.VoteService) *VoteController {
	return &VoteController{voteService: voteService}
}

// VoteOnTarget handles voting on posts or comments
// POST /api/votes/:target_type/:target_id
func (c *VoteController) VoteOnTarget(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	targetType := ctx.Param("target_type")
	targetID := ctx.Param("target_id")

	// Validate target type
	var voteTargetType model.VoteTargetType
	switch targetType {
	case "post":
		voteTargetType = model.VoteTargetPost
	case "comment":
		voteTargetType = model.VoteTargetComment
	default:
		dto.SendError(ctx, http.StatusBadRequest, "Invalid target type. Must be 'post' or 'comment'", apperror.ErrBadRequest.Code)
		return
	}

	var req dto.PostVoteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Invalid request payload", apperror.ErrBadRequest.Code)
		return
	}

	if req.Value == nil {
		dto.SendError(ctx, http.StatusBadRequest, "Vote value is required", apperror.ErrBadRequest.Code)
		return
	}

	userID := authUser.(auth.AuthUser).ID
	votesCount, err := c.voteService.VoteOnTarget(userID, targetID, voteTargetType, *req.Value)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Vote cast successfully", votesCount)
}

// RemoveVote removes a user's vote from a target
// DELETE /api/votes/:target_type/:target_id
func (c *VoteController) RemoveVote(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	targetType := ctx.Param("target_type")
	targetID := ctx.Param("target_id")

	// Validate target type
	var voteTargetType model.VoteTargetType
	switch targetType {
	case "post":
		voteTargetType = model.VoteTargetPost
	case "comment":
		voteTargetType = model.VoteTargetComment
	default:
		dto.SendError(ctx, http.StatusBadRequest, "Invalid target type. Must be 'post' or 'comment'", apperror.ErrBadRequest.Code)
		return
	}

	userID := authUser.(auth.AuthUser).ID

	// Get existing vote first
	existingVote, err := c.voteService.GetUserVote(userID, targetID, voteTargetType)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	if existingVote == nil {
		dto.SendError(ctx, http.StatusNotFound, "No vote found to remove", apperror.ErrVoteNotFound.Code)
		return
	}

	// Remove vote by voting opposite then removing
	_, err = c.voteService.VoteOnTarget(userID, targetID, voteTargetType, !existingVote.Value)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Vote removed successfully", gin.H{"removed": true})
}

// GetUserVote gets the current user's vote on a target
// GET /api/votes/:target_type/:target_id
func (c *VoteController) GetUserVote(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	targetType := ctx.Param("target_type")
	targetID := ctx.Param("target_id")

	// Validate target type
	var voteTargetType model.VoteTargetType
	switch targetType {
	case "post":
		voteTargetType = model.VoteTargetPost
	case "comment":
		voteTargetType = model.VoteTargetComment
	default:
		dto.SendError(ctx, http.StatusBadRequest, "Invalid target type. Must be 'post' or 'comment'", apperror.ErrBadRequest.Code)
		return
	}

	userID := authUser.(auth.AuthUser).ID
	vote, err := c.voteService.GetUserVote(userID, targetID, voteTargetType)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	if vote == nil {
		dto.SendSuccess(ctx, http.StatusOK, "No vote found", gin.H{"vote": nil})
		return
	}

	voteResponse := gin.H{
		"vote": gin.H{
			"target_id":   vote.TargetID.Hex(),
			"target_type": string(vote.TargetType),
			"value":       vote.Value,
			"created_at":  vote.CreateAt,
		},
	}

	dto.SendSuccess(ctx, http.StatusOK, "Vote retrieved successfully", voteResponse)
}

// GetBulkUserVotes gets user's votes for multiple targets of same type
// POST /api/votes/bulk/:target_type
func (c *VoteController) GetBulkUserVotes(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	targetType := ctx.Param("target_type")

	// Validate target type
	var voteTargetType model.VoteTargetType
	switch targetType {
	case "post":
		voteTargetType = model.VoteTargetPost
	case "comment":
		voteTargetType = model.VoteTargetComment
	default:
		dto.SendError(ctx, http.StatusBadRequest, "Invalid target type. Must be 'post' or 'comment'", apperror.ErrBadRequest.Code)
		return
	}

	var req struct {
		TargetIDs []string `json:"target_ids" binding:"required,min=1"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Invalid request payload", apperror.ErrBadRequest.Code)
		return
	}

	if len(req.TargetIDs) > 100 {
		dto.SendError(ctx, http.StatusBadRequest, "Too many target IDs. Maximum 100 allowed", apperror.ErrBadRequest.Code)
		return
	}

	userID := authUser.(auth.AuthUser).ID
	votes, err := c.voteService.FindUserVotes(userID, req.TargetIDs, voteTargetType)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Bulk votes retrieved successfully", gin.H{"votes": votes})
}
