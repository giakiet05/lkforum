package service

import (
	"time"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/platform/bus"
	"github.com/giakiet05/lkforum/internal/repo"
	"github.com/giakiet05/lkforum/internal/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentService interface {
	CreateComment(request *dto.CreateCommentRequest, userID string) (*model.Comment, error)
	GetCommentByID(commentID string) (*model.Comment, error)
	GetCommentByPostIDPaginated(query *dto.GetCommentByPostIDQuery, currentUserID *string) (*dto.PaginatedCommentsResponse, error)
	GetCommentsFilterPaginated(query *dto.GetCommentsFilterQuery, currentUserID *string) (*dto.PaginatedCommentsResponse, error)
	GetAllChildren(commentID string) ([]model.Comment, error)
	DeleteCommentByID(commentID string, userID string) error
	VoteOnComment(userID, commentID string, voteValue bool) (*dto.VotesCountResponse, error)
}

type commentService struct {
	commentRepo   repo.CommentRepo
	voteService   VoteService
	userRepo      repo.UserRepo
	communityRepo repo.CommunityRepo
	postRepo      repo.PostRepo
	bus           bus.EventBus
}

func NewCommentService(
	commentRepo repo.CommentRepo,
	voteService VoteService,
	userRepo repo.UserRepo,
	communityRepo repo.CommunityRepo,
	postRepo repo.PostRepo,
	bus bus.EventBus,
) CommentService {
	return &commentService{
		commentRepo:   commentRepo,
		voteService:   voteService,
		userRepo:      userRepo,
		communityRepo: communityRepo,
		postRepo:      postRepo,
		bus:           bus,
	}
}

func (s *commentService) CreateComment(request *dto.CreateCommentRequest, userID string) (*model.Comment, error) {
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

	post, err := s.postRepo.GetByID(ctx, request.PostID)
	if err != nil {
		return nil, err
	}

	ok, err := s.communityRepo.IsUserBanned(ctx, userID, model.Muted, post.CommunityID.Hex())
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, apperror.ErrUserIsMuted
	}

	var parentAuthorID string
	var parentObjectID *primitive.ObjectID
	if request.ParentID != nil {
		oid, err := primitive.ObjectIDFromHex(*request.ParentID)
		if err != nil {
			return nil, err
		}
		parentObjectID = &oid

		// Fetch the parent comment to get its author's ID for the notification
		parentComment, err := s.commentRepo.GetByID(ctx, *request.ParentID)
		if err == nil { // If parent comment exists
			parentAuthorID = parentComment.Author.ID.Hex()
		}
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	author := model.CommentAuthor{
		ID:       userObjectID,
		Username: user.Username,
		Avatar:   user.RoleContent.AsUser.Avatar,
	}

	comment := &model.Comment{
		Author:           author,
		PostID:           postObjectID,
		ParentID:         parentObjectID,
		Content:          request.Content,
		CreatedAt:        time.Now(),
		IsDeleted:        false,
		ModerationStatus: model.ModerationPending, // Set pending for moderation
	}

	createdComment, err := s.commentRepo.Create(ctx, comment)
	if err != nil {
		return nil, err
	}

	// Publish event for moderation
	s.bus.Publish(&bus.CommentCreatedEvent{
		CommentID:      createdComment.ID.Hex(),
		PostID:         request.PostID,
		AuthorID:       userID,
		ParentAuthorID: &parentAuthorID,
	})

	return createdComment, nil
}

func (s *commentService) GetCommentByID(commentID string) (*model.Comment, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	return s.commentRepo.GetByID(ctx, commentID)
}

func (s *commentService) GetCommentByPostIDPaginated(query *dto.GetCommentByPostIDQuery, currentUserID *string) (*dto.PaginatedCommentsResponse, error) {
	if query.Depth < 0 || query.Depth > 2 {
		return nil, apperror.ErrDepthInvalid
	}

	if query.PageSize < 1 || query.PageSize > 500 || query.Page <= 0 || query.ChildrenPageSize > 100 {
		return nil, apperror.ErrPaginationInvalid
	}

	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	if query.PostID == "" {
		return nil, apperror.ErrInvalidID
	}

	comments, total, err := s.commentRepo.GetCommentsFilterPaginated(ctx, &query.PostID, nil, nil, nil, currentUserID, query.Page, query.PageSize)
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

			children, err := s.commentRepo.GetByParentIDsPaginated(ctx, parentIDs, 1, query.ChildrenPageSize)
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

	commentsResponse := dto.FromComments(comments, currentUserID)
	var response = &dto.PaginatedCommentsResponse{
		Comments: commentsResponse,
		Pagination: dto.Pagination{
			Total: total,
			Page:  query.Page,
		},
	}

	return response, nil
}

func (s *commentService) GetCommentsFilterPaginated(query *dto.GetCommentsFilterQuery, currentUserID *string) (*dto.PaginatedCommentsResponse, error) {
	if query.PageSize < 1 || query.PageSize > 500 || query.Page <= 0 {
		return nil, apperror.ErrPaginationInvalid
	}

	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	comments, total, err := s.commentRepo.GetCommentsFilterPaginated(ctx, query.PostID, query.ParentID, query.UserID, query.Content, currentUserID, query.Page, query.PageSize)
	if err != nil {
		return nil, err
	}

	commentsResponse := dto.FromComments(comments, currentUserID)
	var response = &dto.PaginatedCommentsResponse{
		Comments: commentsResponse,
		Pagination: dto.Pagination{
			Total: total,
			Page:  query.Page,
		},
	}

	return response, nil
}

func (s *commentService) GetAllChildren(commentID string) ([]model.Comment, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	return s.commentRepo.GetAllChildren(ctx, commentID)
}

func (s *commentService) DeleteCommentByID(commentID string, userID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return err
	}

	if comment.IsDeleted || comment.Author.ID.Hex() != userID {
		return apperror.ErrForbidden
	}

	return s.commentRepo.Delete(ctx, commentID)
}

func (s *commentService) VoteOnComment(userID, commentID string, voteValue bool) (*dto.VotesCountResponse, error) {
	return s.voteService.VoteOnTarget(userID, commentID, model.VoteTargetComment, voteValue)
}
