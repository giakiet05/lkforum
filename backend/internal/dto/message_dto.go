package dto

import (
	"time"

	"github.com/giakiet05/lkforum/internal/model"
)

type CreateMessageRequest struct {
	ChannelID string            `json:"channel_id"`
	SenderID  string            `json:"sender_id"`
	Type      model.MessageType `json:"type"`
	Content   string            `json:"content"`
}

type GetMessageFilterQuery struct {
	ChannelID     string  `form:"channel_id" binding:"required"`
	SenderID      *string `form:"sender_id,omitempty"`
	SearchContent *string `form:"search_content,omitempty"`
	IsRead        *bool   `form:"is_read,omitempty"`
	IsSend        *bool   `form:"is_send,omitempty"`
	IsMedia       *bool   `form:"is_media,omitempty"`
	Page          int     `form:"page"`
	PageSize      int     `form:"page_size"`
}

type MessageResponse struct {
	ID             string            `json:"id"`
	ChannelID      string            `json:"channel_id"`
	SenderID       string            `json:"sender_id"`
	SenderUsername string            `json:"sender_username"`
	Type           model.MessageType `json:"type"`
	Content        string            `json:"content"`
	CreatedAt      time.Time         `json:"created_at"`
	IsRead         bool              `json:"is_read"`
}

func FromMessage(message *model.Message) *MessageResponse {
	var senderID string
	if message.SenderID != nil {
		senderID = message.SenderID.Hex()
	}

	return &MessageResponse{
		ID:             message.ID.Hex(),
		ChannelID:      message.ChannelID.Hex(),
		SenderID:       senderID,
		SenderUsername: message.SenderUsername,
		Type:           message.Type,
		Content:        message.Content,
		CreatedAt:      message.CreatedAt,
		IsRead:         message.IsRead,
	}
}

func FromMessages(messages []model.Message) []*MessageResponse {
	var messageResponses []*MessageResponse
	for _, msg := range messages {
		messageResponses = append(messageResponses, FromMessage(&msg))
	}
	return messageResponses
}
