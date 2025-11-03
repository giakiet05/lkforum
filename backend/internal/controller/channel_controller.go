package controller

import (
	"net/http"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/auth"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/service"
	"github.com/gin-gonic/gin"
)

type ChannelController struct {
	channelService service.ChannelService
}

func NewChannelController(channelService service.ChannelService) *ChannelController {
	return &ChannelController{channelService: channelService}
}

func (c *ChannelController) CreateChannel(ctx *gin.Context) {
	var req *dto.CreateChannelRequest
	if err := ctx.ShouldBind(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	channel, err := c.channelService.CreateChannel(req, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusCreated, "Channel created successfully", dto.FromChannel(channel))
}

func (c *ChannelController) GetChannelByID(ctx *gin.Context) {
	channelID := ctx.Param("channel_id")
	if channelID == "" {
		dto.SendError(ctx, http.StatusBadRequest, "Channel ID is required", apperror.ErrBadRequest.Code)
		return
	}

	channel, err := c.channelService.GetChannelByID(channelID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Channel retrieved successfully", dto.FromChannel(channel))
}

func (c *ChannelController) GetChannelsByUserID(ctx *gin.Context) {
	var query *dto.GetChannelByUserIDQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	response, err := c.channelService.GetChannelsByUserID(query.UserID, authUser.(auth.AuthUser).ID, query.Page, query.PageSize)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Channels retrieved successfully", response)
}

func (c *ChannelController) GetChannelByBothUserID(ctx *gin.Context) {
	user1ID := ctx.Param("user1")
	user2ID := ctx.Param("user2")

	if user1ID == "" || user2ID == "" {
		dto.SendError(ctx, http.StatusBadRequest, apperror.ErrBadRequest.Message, apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	channel, err := c.channelService.GetChannelByBothUserID(user1ID, user2ID, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Channel retrieved successfully", dto.FromChannel(channel))
}

func (c *ChannelController) UpdateChannel(ctx *gin.Context) {
	var req *dto.UpdateChannelRequest
	if err := ctx.ShouldBind(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	channel, err := c.channelService.UpdateChannel(req, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Channel updated successfully", dto.FromChannel(channel))
}

func (c *ChannelController) DeleteChannelByID(ctx *gin.Context) {
	channelID := ctx.Param("channel_id")
	if channelID == "" {
		dto.SendError(ctx, http.StatusBadRequest, "Channel ID is required", apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	err := c.channelService.DeleteChannel(channelID, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Channel deleted successfully", gin.H{"id": channelID})
}
