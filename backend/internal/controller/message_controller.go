package controller

import (
	"net/http"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/auth"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/service"
	"github.com/gin-gonic/gin"
)

type MessageController struct {
	messageService service.MessageService
}

func NewMessageController(messageService service.MessageService) *MessageController {
	return &MessageController{
		messageService: messageService,
	}
}

func (mc *MessageController) GetMessageByID(ctx *gin.Context) {
	channelID := ctx.Param("channel_id")
	if channelID == "" {
		dto.SendError(ctx, http.StatusBadRequest, "Community ID is required", apperror.ErrBadRequest.Code)
		return
	}

	messageID := ctx.Param("message_id")
	if messageID == "" {
		dto.SendError(ctx, http.StatusBadRequest, "Message ID is required", apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	message, err := mc.messageService.GetMessageByID(channelID, messageID, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusCreated, "Message retrieved successfully", dto.FromMessage(message))
}

func (mc *MessageController) GetMessageFilter(ctx *gin.Context) {
	var query *dto.GetMessageFilterQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	response, err := mc.messageService.GetMessageFilter(query, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Messages retrieved successfully", response)
}

func (mc *MessageController) DeleteMessageByID(ctx *gin.Context) {
	channelID := ctx.Param("channel_id")
	if channelID == "" {
		dto.SendError(ctx, http.StatusBadRequest, "Community ID is required", apperror.ErrBadRequest.Code)
		return
	}

	messageID := ctx.Param("message_id")
	if messageID == "" {
		dto.SendError(ctx, http.StatusBadRequest, "Message ID is required", apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	err := mc.messageService.DeleteMessage(channelID, messageID, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Message deleted successfully", gin.H{"id": messageID})
}

func (mc *MessageController) MarkChannelAsRead(ctx *gin.Context) {
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

	err := mc.messageService.MarkChannelAsRead(channelID, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Messages marked as read", gin.H{"channel_id": channelID})
}
