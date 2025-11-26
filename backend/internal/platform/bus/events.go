package bus

import (
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/model"
)

// Event Topics
const (
	TopicBroadcast = "broadcast"

	TopicUserChangeAvatar = "user.avatar"

	TopicPostCreated         = "post.created"
	TopicPostUpvoted         = "post.upvoted"
	TopicPostDownvoted       = "post.downvoted"
	TopicPostUpvoteRemoved   = "post.upvote_removed"
	TopicPostDownvoteRemoved = "post.downvote_removed"

	TopicCommentCreated         = "comment.created"
	TopicCommentUpvoted         = "comment.upvoted"
	TopicCommentDownvoted       = "comment.downvoted"
	TopicCommentUpvoteRemoved   = "comment.upvote_removed"
	TopicCommentDownvoteRemoved = "comment.downvote_removed"

	TopicNotificationCreated = "notification.created"

	TopicNewMessage    = "message.new"
	TopicMessageError  = "message.error"
	TopicTypingMessage = "message.typing"
	TopicInChatMessage = "message.in_chat"
)

type BroadcastEventType string

const (
	// ---- Message-related ----
	BroadcastEventMessageCreated BroadcastEventType = "message_created"
	BroadcastEventMessageDeleted BroadcastEventType = "message_deleted"
	BroadcastEventTypingStart    BroadcastEventType = "typing_start"
	BroadcastEventTypingStop     BroadcastEventType = "typing_stop"
	BroadcastEventMessageRead    BroadcastEventType = "message_read"

	// ---- Notification-related ----
	BroadcastEventMessageNotification BroadcastEventType = "message_notification"
)

type BroadcastEvent struct {
	RecipientIDs []string           `json:"recipient_ids"`
	EventType    BroadcastEventType `json:"event_type"`
	TempID       string             `json:"temp_id"`
	Data         interface{}        `json:"data"`
}

func (e BroadcastEvent) Topic() string { return TopicBroadcast }
func (e BroadcastEvent) Payload() map[string]interface{} {
	return map[string]interface{}{
		"recipient_ids": e.RecipientIDs,
		"event_type":    e.EventType,
		"temp_id":       e.TempID,
		"data":          e.Data,
	}
}

type UserChangeAvatarEventType struct {
	UserID    string
	NewAvatar string
}

func (e UserChangeAvatarEventType) Topic() string { return TopicUserChangeAvatar }
func (e UserChangeAvatarEventType) Payload() map[string]interface{} {
	return map[string]interface{}{"user_id": e.UserID, "new_avatar": e.NewAvatar}
}

// --- Post Events ---

type PostCreatedEvent struct {
	AuthorID string
}

func (e PostCreatedEvent) Topic() string { return TopicPostCreated }
func (e PostCreatedEvent) Payload() map[string]interface{} {
	return map[string]interface{}{"authorId": e.AuthorID}
}

type PostUpvotedEvent struct {
	AuthorID string
	VoterID  string
	PostID   string
}

func (e PostUpvotedEvent) Topic() string { return TopicPostUpvoted }
func (e PostUpvotedEvent) Payload() map[string]interface{} {
	return map[string]interface{}{"authorId": e.AuthorID, "voterId": e.VoterID, "postId": e.PostID}
}

type PostDownvotedEvent struct {
	AuthorID string
	VoterID  string
	PostID   string
}

func (e PostDownvotedEvent) Topic() string { return TopicPostDownvoted }
func (e PostDownvotedEvent) Payload() map[string]interface{} {
	return map[string]interface{}{"authorId": e.AuthorID, "voterId": e.VoterID, "postId": e.PostID}
}

type PostUpvoteRemovedEvent struct {
	AuthorID string
	VoterID  string
	PostID   string
}

func (e PostUpvoteRemovedEvent) Topic() string { return TopicPostUpvoteRemoved }
func (e PostUpvoteRemovedEvent) Payload() map[string]interface{} {
	return map[string]interface{}{"authorId": e.AuthorID, "voterId": e.VoterID, "postId": e.PostID}
}

type PostDownvoteRemovedEvent struct {
	AuthorID string
	VoterID  string
	PostID   string
}

func (e PostDownvoteRemovedEvent) Topic() string { return TopicPostDownvoteRemoved }
func (e PostDownvoteRemovedEvent) Payload() map[string]interface{} {
	return map[string]interface{}{"authorId": e.AuthorID, "voterId": e.VoterID, "postId": e.PostID}
}

// --- Comment Events ---

type CommentCreatedEvent struct {
	AuthorID       string
	PostID         string
	CommentID      string
	ParentAuthorID string
}

func (e CommentCreatedEvent) Topic() string { return TopicCommentCreated }
func (e CommentCreatedEvent) Payload() map[string]interface{} {
	return map[string]interface{}{"authorId": e.AuthorID, "postId": e.PostID, "commentId": e.CommentID, "parentAuthorId": e.ParentAuthorID}
}

// --- Notification Events ---

type NotificationCreatedEvent struct {
	RecipientID  string
	Notification dto.NotificationResponse
}

func (e NotificationCreatedEvent) Topic() string { return TopicNotificationCreated }
func (e NotificationCreatedEvent) Payload() map[string]interface{} {
	return map[string]interface{}{"recipientId": e.RecipientID, "notification": e.Notification}
}

// --- Message Events ---

type NewMessageEvent struct {
	TempMessageID  string            `json:"temp_message_id"`
	ChannelID      string            `json:"channel_id"`
	SenderID       string            `json:"sender_id"`
	SenderUsername string            `json:"sender_username"`
	Type           model.MessageType `json:"type"`
	Content        string            `json:"content"`
}

func (e NewMessageEvent) Topic() string { return TopicNewMessage }
func (e NewMessageEvent) Payload() map[string]interface{} {
	return map[string]interface{}{
		"temp_message_id": e.TempMessageID,
		"channel_id":      e.ChannelID,
		"sender_id":       e.SenderID,
		"sender_username": e.SenderUsername,
		"type":            e.Type,
		"content":         e.Content,
	}
}

type MessageErrorEvent struct {
	SenderID      string `json:"sender_id"`
	ChannelID     string `json:"channel_id"`
	TempMessageID string `json:"temp_message_id"`
	ErrorCode     string `json:"error_code"`
	ErrorMsg      string `json:"error_msg"`
}

func (e MessageErrorEvent) Topic() string { return TopicMessageError }
func (e MessageErrorEvent) Payload() map[string]interface{} {
	return map[string]interface{}{
		"sender_id":       e.SenderID,
		"channel_id":      e.ChannelID,
		"temp_message_id": e.TempMessageID,
		"error_code":      e.ErrorCode,
		"error_msg":       e.ErrorMsg,
	}
}

type TypingMessageEvent struct {
	ChannelID string `json:"channel_id"`
	SenderID  string `json:"sender_id"`
	IsTyping  bool   `json:"is_typing"`
}

func (e TypingMessageEvent) Topic() string { return TopicTypingMessage }
func (e TypingMessageEvent) Payload() map[string]interface{} {
	return map[string]interface{}{"channel_id": e.ChannelID, "sender_id": e.SenderID, "is_typing": e.IsTyping}
}

type InChatMessageEvent struct {
	ChannelID string `json:"channel_id"`
	UserID    string `json:"user_id"`
	IsInChat  bool   `json:"is_in_chat"`
}

func (e InChatMessageEvent) Topic() string { return TopicInChatMessage }
func (e InChatMessageEvent) Payload() map[string]interface{} {
	return map[string]interface{}{"channel_id": e.ChannelID, "user_id": e.UserID, "is_in_chat": e.IsInChat}
}
