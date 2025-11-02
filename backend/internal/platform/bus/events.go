package bus

import (
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/model"
)

// Event Topics
const (
	TopicPostCreated         = "post.created"
	TopicPostUpvoted         = "post.upvoted"
	TopicPostDownvoted       = "post.downvoted"
	TopicCommentCreated      = "comment.created"
	TopicCommentUpvoted      = "comment.upvoted"
	TopicCommentDownvoted    = "comment.downvoted"
	TopicNotificationCreated = "notification.created"

	TopicNewMessage     = "message.new"
	TopicMessageCreated = "message.created"
	TopicMessageSend    = "message.send"
	TopicMessageError   = "message.error"
)

// --- Post Events ---

type PostCreatedEvent struct {
	AuthorID string
}

func (e PostCreatedEvent) Topic() string {
	return TopicPostCreated
}
func (e PostCreatedEvent) Payload() map[string]interface{} {
	return map[string]interface{}{"authorId": e.AuthorID}
}

type PostUpvotedEvent struct {
	AuthorID string
	VoterID  string
	PostID   string
}

func (e PostUpvotedEvent) Topic() string {
	return TopicPostUpvoted
}
func (e PostUpvotedEvent) Payload() map[string]interface{} {
	return map[string]interface{}{
		"authorId": e.AuthorID,
		"voterId":  e.VoterID,
		"postId":   e.PostID,
	}
}

type PostDownvotedEvent struct {
	AuthorID string
	VoterID  string
	PostID   string
}

func (e PostDownvotedEvent) Topic() string {
	return TopicPostDownvoted
}
func (e PostDownvotedEvent) Payload() map[string]interface{} {
	return map[string]interface{}{
		"authorId": e.AuthorID,
		"voterId":  e.VoterID,
		"postId":   e.PostID,
	}
}

// --- Comment Events ---

type CommentCreatedEvent struct {
	AuthorID       string
	PostID         string
	CommentID      string
	ParentAuthorID string
}

func (e CommentCreatedEvent) Topic() string {
	return TopicCommentCreated
}
func (e CommentCreatedEvent) Payload() map[string]interface{} {
	return map[string]interface{}{
		"authorId":       e.AuthorID,
		"postId":         e.PostID,
		"commentId":      e.CommentID,
		"parentAuthorId": e.ParentAuthorID,
	}
}

// --- Notification Events ---

type NotificationCreatedEvent struct {
	RecipientID  string
	Notification dto.NotificationResponse
}

func (e NotificationCreatedEvent) Topic() string {
	return TopicNotificationCreated
}
func (e NotificationCreatedEvent) Payload() map[string]interface{} {
	return map[string]interface{}{
		"recipientId":  e.RecipientID,
		"notification": e.Notification,
	}
}

// --- Message Events ---

type NewMessageEvent struct {
	TempMessageID string            `json:"temp_message_id"`
	ChannelID     string            `json:"channel_id"`
	SenderID      string            `json:"sender_id"`
	Type          model.MessageType `json:"type"`
	Content       string            `json:"content"`
}

func (e NewMessageEvent) Topic() string {
	return TopicNewMessage
}

func (e NewMessageEvent) Payload() map[string]interface{} {
	return map[string]interface{}{
		"temp_message_id": e.TempMessageID,
		"channel_id":      e.ChannelID,
		"sender_id":       e.SenderID,
		"type":            e.Type,
		"content":         e.Content,
	}
}

type MessageCreatedEvent struct {
	RecipientIDs  []string
	TempMessageID string
	Message       dto.MessageResponse
}

func (e MessageCreatedEvent) Topic() string {
	return TopicMessageCreated
}

func (e MessageCreatedEvent) Payload() map[string]interface{} {
	return map[string]interface{}{
		"recipient_ids":   e.RecipientIDs,
		"temp_message_id": e.TempMessageID,
		"message":         e.Message,
	}
}

type MessageSendEvent struct {
	MessageID  string   `json:"message_id"`
	ChannelID  string   `json:"channel_id"`
	SenderID   string   `json:"sender_id"`
	ReceivedID []string `json:"received_id"`
}

func (e MessageSendEvent) Topic() string {
	return TopicMessageSend
}

func (e MessageSendEvent) Payload() map[string]interface{} {
	return map[string]interface{}{
		"message_id":  e.MessageID,
		"channel_id":  e.ChannelID,
		"sender_id":   e.SenderID,
		"received_id": e.ReceivedID,
	}
}

type MessageErrorEvent struct {
	SenderID      string `json:"sender_id"`
	ChannelID     string `json:"channel_id"`
	TempMessageID string `json:"temp_message_id"`
	ErrorCode     string `json:"error_code"`
	ErrorMsg      string `json:"error_msg"`
}

func (e MessageErrorEvent) Topic() string {
	return TopicMessageError
}

func (e MessageErrorEvent) Payload() map[string]interface{} {
	return map[string]interface{}{
		"sender_id":       e.SenderID,
		"channel_id":      e.ChannelID,
		"temp_message_id": e.TempMessageID,
		"error_code":      e.ErrorCode,
		"error_msg":       e.ErrorMsg,
	}
}
