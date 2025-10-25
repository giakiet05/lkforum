package bus

import "github.com/giakiet05/lkforum/internal/dto"

// Event Topics
const (
	TopicPostCreated         = "post.created"
	TopicPostUpvoted         = "post.upvoted"
	TopicPostDownvoted       = "post.downvoted"
	TopicCommentCreated      = "comment.created"
	TopicCommentUpvoted      = "comment.upvoted"
	TopicCommentDownvoted    = "comment.downvoted"
	TopicNotificationCreated = "notification.created"
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
