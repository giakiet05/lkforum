package ws

import (
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/model"
)

type WebSocketMessageType string

const (
	NewNotification WebSocketMessageType = "new_notification"
	ACKMessage      WebSocketMessageType = "ack_message"
	NewMessage      WebSocketMessageType = "new_message"
	SendMessage     WebSocketMessageType = "send_message"
	TypingIndicator WebSocketMessageType = "typing"
	InChatIndicator WebSocketMessageType = "in_chat"
	ErrorMessage    WebSocketMessageType = "error"
)

type WebSocketMessage struct {
	Type    WebSocketMessageType `json:"type"`
	Payload interface{}          `json:"payload"`
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
	ChannelID string `json:"channel_id"`
	SenderID  string `json:"sender_id"`
	IsTyping  bool   `json:"is_typing"`
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
