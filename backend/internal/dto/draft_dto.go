package dto

import (
	"time"

	"github.com/giakiet05/lkforum/internal/model"
)

// --- Request DTOs ---

type CreateDraftRequest struct {
	CommunityID *string          `json:"community_id,omitempty"`
	Type        *model.PostType  `json:"type,omitempty"`
	Title       *string          `json:"title,omitempty"`
	Text        *string          `json:"text,omitempty"`
	Tags        []string         `json:"tags,omitempty"`
	Images      []model.Image    `json:"images,omitempty"`
	Videos      []*model.Video   `json:"videos,omitempty"`
	Poll        *model.Poll      `json:"poll,omitempty"`
}

type UpdateDraftRequest struct {
	CommunityID *string          `json:"community_id,omitempty"`
	Type        *model.PostType  `json:"type,omitempty"`
	Title       *string          `json:"title,omitempty"`
	Text        *string          `json:"text,omitempty"`
	Tags        []string         `json:"tags,omitempty"`
	Images      []model.Image    `json:"images,omitempty"`
	Videos      []*model.Video   `json:"videos,omitempty"`
	Poll        *model.Poll      `json:"poll,omitempty"`
}

// --- Response DTOs ---

// DraftSummaryResponse for listing drafts.
type DraftSummaryResponse struct {
	ID          string     `json:"id"`
	Title       *string    `json:"title,omitempty"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// DraftResponse is the detailed draft object returned to the client.
type DraftResponse struct {
	ID          string           `json:"id"`
	AuthorID    string           `json:"author_id"`
	CommunityID *string          `json:"community_id,omitempty"`
	Type        *model.PostType  `json:"type,omitempty"`
	Title       *string          `json:"title,omitempty"`
	Content     *model.PostContent `json:"content,omitempty"`
	Tags        []string         `json:"tags,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

type PaginatedDraftsResponse struct {
	Drafts     []*DraftSummaryResponse `json:"drafts"`
	Pagination Pagination              `json:"pagination"`
}


// --- Mapping Functions ---

func FromDraft(d *model.Draft) *DraftResponse {
	if d == nil {
		return nil
	}
	return &DraftResponse{
		ID:          d.ID.Hex(),
		AuthorID:    d.AuthorID.Hex(),
		CommunityID: d.CommunityID,
		Type:        d.Type,
		Title:       d.Title,
		Content:     d.Content,
		Tags:        d.Tags,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}

func FromDraftsToSummary(drafts []*model.Draft) []*DraftSummaryResponse {
	res := make([]*DraftSummaryResponse, len(drafts))
	for i, d := range drafts {
		res[i] = &DraftSummaryResponse{
			ID:        d.ID.Hex(),
			Title:     d.Title,
			UpdatedAt: d.UpdatedAt,
		}
	}
	return res
}
