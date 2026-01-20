package dto

import (
	"time"

	"github.com/giakiet05/lkforum/internal/model"
)

// NotificationResponse defines the structure for a notification returned to the client.
type NotificationResponse struct {
	ID        string                 `json:"id"`
	Type      model.NotificationType `json:"type"`
	Message   string                 `json:"message"`
	Link      string                 `json:"link"`
	IsRead    bool                   `json:"is_read"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	// We can add actor information here later if needed
	// Actor     *ShortUserResponse `json:"actor,omitempty"`
}

// PaginatedNotificationsResponse defines the structure for a paginated list of notifications.
type PaginatedNotificationsResponse struct {
	Notifications []NotificationResponse `json:"notifications"`
	Pagination    Pagination             `json:"pagination"`
}

// FromNotification converts a model.Notification to a NotificationResponse DTO.
func FromNotification(n *model.Notification) NotificationResponse {
	return NotificationResponse{
		ID:        n.ID.Hex(),
		Type:      n.Type,
		Message:   n.Message,
		Link:      n.Link,
		IsRead:    n.IsRead,
		Metadata:  n.Metadata,
		CreatedAt: n.CreatedAt,
	}
}

// FromNotifications converts a slice of model.Notification to a slice of NotificationResponse DTOs.
func FromNotifications(notifications []*model.Notification) []NotificationResponse {
	responses := make([]NotificationResponse, len(notifications))
	for i, n := range notifications {
		responses[i] = FromNotification(n)
	}
	return responses
}
