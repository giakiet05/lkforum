package dto

import (
	"time"

	"github.com/giakiet05/lkforum/internal/model"
)

type CreatePostHistoryRequest struct {
	UserID string `json:"user_id"`
	PostID string `json:"post_id"`
}

type PostHistoryResponse struct {
	ID       string    `json:"id"`
	PostID   string    `json:"post_id"`
	UserID   string    `json:"user_id"`
	ViewedAt time.Time `json:"viewed_at"`
}

func FromPostHistory(postHistory *model.PostHistory) *PostHistoryResponse {
	return &PostHistoryResponse{
		ID:       postHistory.ID.Hex(),
		PostID:   postHistory.PostID.Hex(),
		UserID:   postHistory.UserID.Hex(),
		ViewedAt: postHistory.ViewedAt,
	}
}

func FromPostHistories(postHistories []*model.PostHistory) []*PostHistoryResponse {
	responses := make([]*PostHistoryResponse, len(postHistories))
	for i, ph := range postHistories {
		responses[i] = FromPostHistory(ph)
	}
	return responses
}
