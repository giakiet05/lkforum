package ws

import (
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/model"
)

type SocketMessageType string

const (
	NewNotification SocketMessageType = "new_notification"
	ACKMessage      SocketMessageType = "ack_message"
	NewMessage      SocketMessageType = "new_message"
	SendMessage     SocketMessageType = "send_message"
	TypingIndicator SocketMessageType = "typing"
	InChatIndicator SocketMessageType = "in_chat"
	ErrorMessage    SocketMessageType = "error"
)

type WebSocketMessage struct {
	Type    SocketMessageType `json:"type"`
	Payload interface{}       `json:"payload"`
}

type NewMessagePayload struct {
	TempMessageID string            `json:"temp_message_id"`
	ChannelID     string            `json:"channel_id"`
	Type          model.MessageType `json:"type"`
	Content       string            `json:"content"`
}

type SendMessagePayload struct {
	Message dto.MessageResponse `json:"message"`
}

type ACKMessagePayload struct {
	TempMessageID string              `json:"temp_message_id"`
	Message       dto.MessageResponse `json:"message"`
}

type TypingIndicatorPayload struct {
	ConversationId string `json:"conversation_id"`
	UserID         string `json:"user_id"`
	IsTyping       bool   `json:"is_typing"`
}

type InChatIndicatorPayload struct {
	TempMessageID  string `json:"temp_message_id"`
	ConversationId string `json:"conversation_id"`
	IsInChat       bool   `json:"is_in_chat"`
}

type ErrorPayload struct {
	TempMessageID *string `json:"temp_message_id,omitempty"`
	ErrorCode     *string `json:"error_code,omitempty"`
	ErrorMsg      string  `json:"error_msg"`
}
