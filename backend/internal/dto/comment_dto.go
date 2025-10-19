package dto

import "github.com/giakiet05/lkforum/internal/model"

type CreateCommentRequest struct {
	UserID     string  `json:"user_id"`
	Username   string  `json:"username"`
	UserAvatar string  `json:"user_avatar"`
	PostID     string  `json:"post_id"`
	ParentID   *string `json:"parent_id,omitempty"`
	Content    string  `json:"content"`
}

type GetCommentsFilterQuery struct {
	PostID   *string `form:"post_id,omitempty"`
	ParentID *string `form:"parent_id,omitempty"`
	UserID   *string `form:"user_id,omitempty"`
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
	ID        string              `json:"id"`
	Author    model.CommentAuthor `json:"author"`
	PostID    string              `json:"post_id"`
	ParentID  *string             `json:"parent_id,omitempty"`
	Children  *[]CommentResponse  `json:"children,omitempty"`
	Content   string              `json:"content"`
	CreatedAt string              `json:"created_at"`
	IsDeleted bool                `json:"is_deleted"`
}

func FromComments(comments []model.Comment) []CommentResponse {
	commentMap := make(map[string]*CommentResponse)
	for _, c := range comments {
		resp := FromComment(&c)
		commentMap[resp.ID] = resp
	}

	var roots []CommentResponse
	for _, c := range comments {
		resp := commentMap[c.ID.Hex()]
		if c.ParentID != nil {
			parentResp, ok := commentMap[c.ParentID.Hex()]
			if ok {
				if parentResp.Children == nil {
					parentResp.Children = &[]CommentResponse{}
				}
				*parentResp.Children = append(*parentResp.Children, *resp)
			}
		} else {
			roots = append(roots, *resp)
		}
	}

	return roots
}

func FromComment(comment *model.Comment) *CommentResponse {
	var parentID *string
	if comment.ParentID != nil {
		pid := comment.ParentID.Hex()
		parentID = &pid
	}

	return &CommentResponse{
		ID:        comment.ID.Hex(),
		Author:    comment.Author,
		PostID:    comment.PostID.Hex(),
		ParentID:  parentID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt.Format("2006-01-02 15:04:05"),
		IsDeleted: comment.IsDeleted,
	}
}

func FromCommentWithChildren(comments []model.Comment) *CommentResponse {
	if len(comments) == 0 {
		return nil
	}

	commentMap := make(map[string]*CommentResponse)
	for _, c := range comments {
		resp := FromComment(&c)
		commentMap[resp.ID] = resp
	}

	var root *CommentResponse
	for _, c := range comments {
		resp := commentMap[c.ID.Hex()]
		if c.ParentID != nil {
			parentResp, ok := commentMap[c.ParentID.Hex()]
			if ok {
				if parentResp.Children == nil {
					parentResp.Children = &[]CommentResponse{}
				}
				*parentResp.Children = append(*parentResp.Children, *resp)
			}
		} else {
			root = resp
		}
	}
	return root
}
