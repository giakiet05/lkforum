package dto

import (
	"time"

	"github.com/giakiet05/lkforum/internal/model"
)

type CreateCommentRequest struct {
	PostID   string  `json:"post_id" binding:"required"`
	ParentID *string `json:"parent_id,omitempty"`
	Content  string  `json:"content" binding:"required"`
}

type GetCommentsFilterQuery struct {
	PostID   *string `form:"post_id,omitempty"`
	ParentID *string `form:"parent_id,omitempty"`
	UserID   *string `form:"user_id,omitempty"`
	Content  *string `form:"content,omitempty"`
	Page     int     `form:"page"`
	PageSize int     `form:"page_size"`
}

type GetCommentByPostIDQuery struct {
	PostID           string `form:"post_id"`
	Depth            int    `form:"depth"`
	ChildrenPageSize int    `form:"children_page_size"`
	Page             int    `form:"page"`
	PageSize         int    `form:"page_size"`
}

type CommentResponse struct {
	ID               string              `json:"id"`
	Author           model.CommentAuthor `json:"author"`
	PostID           string              `json:"post_id"`
	ParentID         *string             `json:"parent_id,omitempty"`
	Children         []*CommentResponse  `json:"children"`
	Content          string              `json:"content"`
	CreatedAt        string              `json:"created_at"`
	IsDeleted        bool                `json:"is_deleted"`
	ModerationStatus *string             `json:"moderation_status,omitempty"`
	ModerationReason *string             `json:"moderation_reason,omitempty"`
	ModeratedAt      *time.Time          `json:"moderated_at,omitempty"`
}

func FromComments(comments []model.Comment, currentUserID *string) []*CommentResponse {
	commentMap := make(map[string]*CommentResponse)
	for _, c := range comments {
		resp := FromComment(&c, currentUserID)
		commentMap[resp.ID] = resp
	}

	var roots []*CommentResponse
	addedToParent := make(map[string]bool) // Track which comments were added as children

	for _, c := range comments {
		resp := commentMap[c.ID.Hex()]
		if c.ParentID != nil {
			parentResp, ok := commentMap[c.ParentID.Hex()]
			if ok {
				if parentResp.Children == nil {
					parentResp.Children = []*CommentResponse{}
				}
				parentResp.Children = append(parentResp.Children, resp)
				addedToParent[c.ID.Hex()] = true
			} else {
				// Parent not found, treat as root (orphaned comment)
				roots = append(roots, resp)
			}
		} else {
			// Actual root comment (no parent)
			roots = append(roots, resp)
		}
	}

	result := make([]*CommentResponse, len(roots))
	for i, r := range roots {
		result[i] = r
	}
	return result
}

func FromComment(comment *model.Comment, currentUserID *string) *CommentResponse {
	var parentID *string
	if comment.ParentID != nil {
		pid := comment.ParentID.Hex()
		parentID = &pid
	}

	resp := &CommentResponse{
		ID:        comment.ID.Hex(),
		Author:    comment.Author,
		PostID:    comment.PostID.Hex(),
		ParentID:  parentID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt.Format("2006-01-02 15:04:05"),
		IsDeleted: comment.IsDeleted,
	}

	return resp
}

func FromCommentWithChildren(comments []model.Comment, currentUserID *string) *CommentResponse {
	if len(comments) == 0 {
		return nil
	}

	commentMap := make(map[string]*CommentResponse)
	for _, c := range comments {
		resp := FromComment(&c, currentUserID)
		commentMap[resp.ID] = resp
	}

	var root *CommentResponse
	for _, c := range comments {
		resp := commentMap[c.ID.Hex()]
		if c.ParentID != nil {
			parentResp, ok := commentMap[c.ParentID.Hex()]
			if ok {
				if parentResp.Children == nil {
					parentResp.Children = []*CommentResponse{}
				}
				parentResp.Children = append(parentResp.Children, resp)
			}
		} else {
			root = resp
		}
	}
	return root
}
