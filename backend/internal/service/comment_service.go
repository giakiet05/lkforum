package service

import (
	"time"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/repo"
	"github.com/giakiet05/lkforum/internal/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentService interface {
	CreateComment(request *dto.CreateCommentRequest, userID string) (*model.Comment, error)
	GetCommentByID(commentID string) (*model.Comment, error)
	GetCommentByPostIDPaginated(query *dto.GetCommentByPostIDQuery) (*dto.PaginatedCommentsResponse, error)
	GetCommentsFilterPaginated(query *dto.GetCommentsFilterQuery) (*dto.PaginatedCommentsResponse, error)
	GetAllChildren(commentID string) ([]model.Comment, error)
	DeleteCommentByID(commentID string, userID string) error
}

type commentService struct {
	commentRepo repo.CommentRepo
}

func NewCommentService(commentRepo repo.CommentRepo) CommentService {
	return &commentService{commentRepo: commentRepo}
}

func (c *commentService) CreateComment(request *dto.CreateCommentRequest, userID string) (*model.Comment, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	postObjectID, err := primitive.ObjectIDFromHex(request.PostID)
	if err != nil {
		return nil, err
	}

	var parentObjectID *primitive.ObjectID
	if request.ParentID != nil {
		oid, err := primitive.ObjectIDFromHex(*request.ParentID)
		if err != nil {
			return nil, err
		}
		parentObjectID = &oid
	}

	author := model.CommentAuthor{
		ID:       userObjectID,
		Username: request.Username,
		Avatar:   request.UserAvatar,
	}

	comment := &model.Comment{
		Author:    author,
		PostID:    postObjectID,
		ParentID:  parentObjectID,
		Content:   request.Content,
		CreatedAt: time.Now(),
		IsDeleted: false,
	}

	comment, err = c.commentRepo.Create(ctx, comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (c *commentService) GetCommentByID(commentID string) (*model.Comment, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	return c.commentRepo.GetByID(ctx, commentID)
}

func (c *commentService) GetCommentByPostIDPaginated(query *dto.GetCommentByPostIDQuery) (*dto.PaginatedCommentsResponse, error) {
	if query.Depth < 0 || query.Depth > 2 {
		return nil, apperror.ErrDepthInvalid
	}

	if query.PageSize < 1 || query.PageSize > 500 || query.Page <= 0 {
		return nil, apperror.ErrPaginationInvalid
	}

	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	comments, total, err := c.commentRepo.GetCommentsPaginated(ctx, &query.PostID, nil, nil, query.Page, query.PageSize)
	if err != nil {
		return nil, err
	}

	if query.Depth > 0 {
		currentLevel := comments

		for d := 0; d < query.Depth; d++ {
			if len(currentLevel) == 0 {
				break
			}

			var parentIDs []string
			for _, cmt := range currentLevel {
				parentIDs = append(parentIDs, cmt.ID.Hex())
			}

			children, err := c.commentRepo.GetByParentIDsPaginated(ctx, parentIDs, 1, query.ChildrenPageSize)
			if err != nil {
				return nil, err
			}

			if len(children) == 0 {
				break
			}

			comments = append(comments, children...)
			currentLevel = children
		}
	}

	commentsResponse := dto.FromComments(comments)
	var response = &dto.PaginatedCommentsResponse{
		Comments: commentsResponse,
		Pagination: dto.Pagination{
			Total: total,
			Page:  query.Page,
		},
	}

	return response, nil
}

func (c *commentService) GetCommentsFilterPaginated(query *dto.GetCommentsFilterQuery) (*dto.PaginatedCommentsResponse, error) {
	if query.PageSize < 1 || query.PageSize > 500 || query.Page <= 0 {
		return nil, apperror.ErrPaginationInvalid
	}

	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	comments, total, err := c.commentRepo.GetCommentsPaginated(ctx, query.PostID, query.ParentID, query.UserID, query.Page, query.PageSize)
	if err != nil {
		return nil, err
	}

	commentsResponse := dto.FromComments(comments)
	var response = &dto.PaginatedCommentsResponse{
		Comments: commentsResponse,
		Pagination: dto.Pagination{
			Total: total,
			Page:  query.Page,
		},
	}

	return response, nil
}

func (c *commentService) GetAllChildren(commentID string) ([]model.Comment, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	return c.commentRepo.GetAllChildren(ctx, commentID)
}

func (c *commentService) DeleteCommentByID(commentID string, userID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	comment, err := c.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return err
	}

	if comment.IsDeleted || comment.Author.ID.Hex() != userID {
		return apperror.ErrForbidden
	}

	return c.commentRepo.Delete(ctx, commentID)
}
