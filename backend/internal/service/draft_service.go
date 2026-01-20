package service

import (
	"context"
	"log"
	"time"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/repo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DraftService interface {
	CreateDraft(ctx context.Context, authorID string, req *dto.CreateDraftRequest) (*dto.DraftResponse, error)
	UpdateDraft(ctx context.Context, authorID, draftID string, req *dto.UpdateDraftRequest) (*dto.DraftResponse, error)
	GetDraftByID(ctx context.Context, authorID, draftID string) (*dto.DraftResponse, error)
	GetDraftsByAuthor(ctx context.Context, authorID string, query *dto.GetPostsQuery) ([]*dto.DraftSummaryResponse, int64, error)
	DeleteDraft(ctx context.Context, authorID, draftID string) error
	PublishDraft(ctx context.Context, authorID, draftID string) (*dto.PostResponse, error)
}

type draftService struct {
	draftRepo repo.DraftRepo
	postRepo  repo.PostRepo
	postService PostService
}

func NewDraftService(draftRepo repo.DraftRepo, postRepo repo.PostRepo, postService PostService) DraftService {
	return &draftService{
		draftRepo: draftRepo,
		postRepo:  postRepo,
		postService: postService,
	}
}

func (s *draftService) CreateDraft(ctx context.Context, authorID string, req *dto.CreateDraftRequest) (*dto.DraftResponse, error) {
	authorObjID, err := primitive.ObjectIDFromHex(authorID)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}

	now := time.Now()
	postContent := &model.PostContent{
		Images: req.Images,
		Videos: req.Videos,
		Poll:   req.Poll,
	}
	if req.Text != nil {
		postContent.Text = *req.Text
	}

	draft := &model.Draft{
		AuthorID:    authorObjID,
		CommunityID: req.CommunityID,
		Type:        req.Type,
		Title:       req.Title,
		Content:     postContent,
		Tags:        req.Tags,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	createdDraft, err := s.draftRepo.Create(ctx, draft)
	if err != nil {
		return nil, err
	}

	return dto.FromDraft(createdDraft), nil
}

func (s *draftService) UpdateDraft(ctx context.Context, authorID, draftID string, req *dto.UpdateDraftRequest) (*dto.DraftResponse, error) {
	authorObjID, err := primitive.ObjectIDFromHex(authorID)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}
	draftObjID, err := primitive.ObjectIDFromHex(draftID)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}

	draft, err := s.draftRepo.GetByID(ctx, draftObjID)
	if err != nil {
		return nil, err
	}

	if draft.AuthorID != authorObjID {
		return nil, apperror.ErrForbidden
	}

	// Update fields from request
	draft.CommunityID = req.CommunityID
	draft.Type = req.Type
	draft.Title = req.Title
	draft.Tags = req.Tags
	if draft.Content == nil {
		draft.Content = &model.PostContent{}
	}
	if req.Text != nil {
		draft.Content.Text = *req.Text
	}
	draft.Content.Images = req.Images
	draft.Content.Videos = req.Videos
	draft.Content.Poll = req.Poll
	draft.UpdatedAt = time.Now()

	if err := s.draftRepo.Update(ctx, draft); err != nil {
		return nil, err
	}

	return dto.FromDraft(draft), nil
}

func (s *draftService) GetDraftByID(ctx context.Context, authorID, draftID string) (*dto.DraftResponse, error) {
	authorObjID, err := primitive.ObjectIDFromHex(authorID)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}
	draftObjID, err := primitive.ObjectIDFromHex(draftID)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}

	draft, err := s.draftRepo.GetByID(ctx, draftObjID)
	if err != nil {
		return nil, err
	}

	if draft.AuthorID != authorObjID {
		return nil, apperror.ErrForbidden
	}

	return dto.FromDraft(draft), nil
}

func (s *draftService) DeleteDraft(ctx context.Context, authorID, draftID string) error {
	authorObjID, err := primitive.ObjectIDFromHex(authorID)
	if err != nil {
		return apperror.ErrInvalidID
	}
	draftObjID, err := primitive.ObjectIDFromHex(draftID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	draft, err := s.draftRepo.GetByID(ctx, draftObjID)
	if err != nil {
		return err
	}

	if draft.AuthorID != authorObjID {
		return apperror.ErrForbidden
	}

	return s.draftRepo.Delete(ctx, draftObjID)
}

func (s *draftService) GetDraftsByAuthor(ctx context.Context, authorID string, query *dto.GetPostsQuery) ([]*dto.DraftSummaryResponse, int64, error) {
	authorObjID, err := primitive.ObjectIDFromHex(authorID)
	if err != nil {
		return nil, 0, apperror.ErrInvalidID
	}

	page := query.Page
	if page <= 0 {
		page = 1
	}
	limit := query.Limit
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	opts := &repo.FindOptions{
		Skip:  int64((page - 1) * limit),
		Limit: int64(limit),
		Sort:  map[string]int{"updated_at": -1},
	}

	drafts, total, err := s.draftRepo.GetByAuthor(ctx, authorObjID, opts)
	if err != nil {
		return nil, 0, err
	}

	return dto.FromDraftsToSummary(drafts), total, nil
}

func (s *draftService) PublishDraft(ctx context.Context, authorID, draftID string) (*dto.PostResponse, error) {
	authorObjID, err := primitive.ObjectIDFromHex(authorID)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}
	draftObjID, err := primitive.ObjectIDFromHex(draftID)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}

	draft, err := s.draftRepo.GetByID(ctx, draftObjID)
	if err != nil {
		return nil, err
	}

	if draft.AuthorID != authorObjID {
		return nil, apperror.ErrForbidden
	}

	// Validate required fields for publishing
	if draft.Title == nil || *draft.Title == "" {
		return nil, apperror.NewError(nil, apperror.ErrBadRequest.Code, "Title is required to publish a draft")
	}
	if draft.CommunityID == nil || *draft.CommunityID == "" {
		return nil, apperror.NewError(nil, apperror.ErrBadRequest.Code, "Community is required to publish a draft")
	}
	if draft.Type == nil {
		return nil, apperror.NewError(nil, apperror.ErrBadRequest.Code, "Post type is required to publish a draft")
	}

	// Create a Post from the draft
	createPostReq := &dto.CreatePostRequest{
		CommunityID: *draft.CommunityID,
		Title:       *draft.Title,
		Type:        *draft.Type,
		Tags:        draft.Tags,
	}
	if draft.Content != nil {
		createPostReq.Text = draft.Content.Text
		if draft.Content.Poll != nil {
			// This requires converting model.Poll to dto.CreatePollRequest
			// For now, let's keep it simple and assume poll is created on publish
			// A more complex implementation would handle this conversion
		}
	}

	// Use the existing PostService to create the post, which handles upvoting etc.
	postResponse, err := s.postService.CreatePost(authorID, createPostReq)
	if err != nil {
		return nil, err
	}

	// If post creation is successful, delete the draft
	if err := s.draftRepo.Delete(ctx, draftObjID); err != nil {
		// Log this error, but don't fail the whole operation since the post is already created.
		// This could lead to an orphaned draft, but it's better than failing the user's publish action.
		log.Printf("ERROR: Failed to delete draft after publishing: %v", err)
	}

	return postResponse, nil
}





