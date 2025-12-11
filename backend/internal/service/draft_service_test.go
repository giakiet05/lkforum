package service

import (
	"context"
	"errors"
	"mime/multipart"
	"testing"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/repo/mocks"
	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// TestCreateDraft tests the CreateDraft method
func TestCreateDraft(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDraftRepo := mocks.NewMockDraftRepo(ctrl)
	svc := NewDraftService(mockDraftRepo, nil, nil)

	ctx := context.Background()
	authorID := primitive.NewObjectID()
	communityID := "community123"
	title := "Test Draft"
	text := "Draft content"
	postType := model.PostTypeText

	tests := []struct {
		name      string
		authorID  string
		request   *dto.CreateDraftRequest
		mockSetup func()
		wantErr   error
	}{
		{
			name:     "success create text draft",
			authorID: authorID.Hex(),
			request: &dto.CreateDraftRequest{
				CommunityID: &communityID,
				Type:        &postType,
				Title:       &title,
				Text:        &text,
			},
			mockSetup: func() {
				mockDraftRepo.EXPECT().
					Create(gomock.Any(), gomock.AssignableToTypeOf(&model.Draft{})).
					DoAndReturn(func(ctx context.Context, draft *model.Draft) (*model.Draft, error) {
						// Verify draft structure
						if draft.AuthorID != authorID {
							t.Errorf("Expected AuthorID %v, got %v", authorID, draft.AuthorID)
						}
						if *draft.Title != title {
							t.Errorf("Expected title %s, got %s", title, *draft.Title)
						}
						// Set ID for return
						draft.ID = primitive.NewObjectID()
						return draft, nil
					})
			},
			wantErr: nil,
		},
		{
			name:     "error invalid author id",
			authorID: "invalid-id",
			request: &dto.CreateDraftRequest{
				Title: &title,
			},
			mockSetup: func() {
				// No mock needed - should fail validation
			},
			wantErr: apperror.ErrInvalidID,
		},
		{
			name:     "error repository failure",
			authorID: authorID.Hex(),
			request: &dto.CreateDraftRequest{
				Title: &title,
			},
			mockSetup: func() {
				mockDraftRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("database error"))
			},
			wantErr: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			result, err := svc.CreateDraft(ctx, tt.authorID, tt.request)

			// Verify error
			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.wantErr)
					return
				}
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			// Verify success
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Error("expected result, got nil")
			}
		})
	}
}

// TestGetDraftByID tests the GetDraftByID method
func TestGetDraftByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDraftRepo := mocks.NewMockDraftRepo(ctrl)
	svc := NewDraftService(mockDraftRepo, nil, nil)

	ctx := context.Background()
	authorID := primitive.NewObjectID()
	otherAuthorID := primitive.NewObjectID()
	draftID := primitive.NewObjectID()

	title := "Test Draft"
	existingDraft := &model.Draft{
		ID:       draftID,
		AuthorID: authorID,
		Title:    &title,
	}

	draftByOtherAuthor := &model.Draft{
		ID:       draftID,
		AuthorID: otherAuthorID,
		Title:    &title,
	}

	tests := []struct {
		name      string
		authorID  string
		draftID   string
		mockSetup func()
		wantErr   error
	}{
		{
			name:     "success get draft by id",
			authorID: authorID.Hex(),
			draftID:  draftID.Hex(),
			mockSetup: func() {
				mockDraftRepo.EXPECT().
					GetByID(gomock.Any(), draftID).
					Return(existingDraft, nil)
			},
			wantErr: nil,
		},
		{
			name:     "error invalid author id",
			authorID: "invalid-id",
			draftID:  draftID.Hex(),
			mockSetup: func() {
				// No mock needed
			},
			wantErr: apperror.ErrInvalidID,
		},
		{
			name:     "error invalid draft id",
			authorID: authorID.Hex(),
			draftID:  "invalid-id",
			mockSetup: func() {
				// No mock needed
			},
			wantErr: apperror.ErrInvalidID,
		},
		{
			name:     "error draft not found",
			authorID: authorID.Hex(),
			draftID:  primitive.NewObjectID().Hex(),
			mockSetup: func() {
				mockDraftRepo.EXPECT().
					GetByID(gomock.Any(), gomock.Any()).
					Return(nil, mongo.ErrNoDocuments)
			},
			wantErr: mongo.ErrNoDocuments,
		},
		{
			name:     "error not draft author (forbidden)",
			authorID: otherAuthorID.Hex(),
			draftID:  draftID.Hex(),
			mockSetup: func() {
				mockDraftRepo.EXPECT().
					GetByID(gomock.Any(), draftID).
					Return(draftByOtherAuthor, nil)
			},
			wantErr: apperror.ErrForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			result, err := svc.GetDraftByID(ctx, tt.authorID, tt.draftID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.wantErr)
					return
				}
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Error("expected result, got nil")
			}
		})
	}
}

// TestDeleteDraft tests the DeleteDraft method
func TestDeleteDraft(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDraftRepo := mocks.NewMockDraftRepo(ctrl)
	svc := NewDraftService(mockDraftRepo, nil, nil)

	ctx := context.Background()
	authorID := primitive.NewObjectID()
	otherAuthorID := primitive.NewObjectID()
	draftID := primitive.NewObjectID()

	title := "Test Draft"
	existingDraft := &model.Draft{
		ID:       draftID,
		AuthorID: authorID,
		Title:    &title,
	}

	draftByOtherAuthor := &model.Draft{
		ID:       draftID,
		AuthorID: otherAuthorID,
		Title:    &title,
	}

	tests := []struct {
		name      string
		authorID  string
		draftID   string
		mockSetup func()
		wantErr   error
	}{
		{
			name:     "success delete draft",
			authorID: authorID.Hex(),
			draftID:  draftID.Hex(),
			mockSetup: func() {
				mockDraftRepo.EXPECT().
					GetByID(gomock.Any(), draftID).
					Return(existingDraft, nil)

				mockDraftRepo.EXPECT().
					Delete(gomock.Any(), draftID).
					Return(nil)
			},
			wantErr: nil,
		},
		{
			name:     "error invalid author id",
			authorID: "invalid-id",
			draftID:  draftID.Hex(),
			mockSetup: func() {
				// No mock needed
			},
			wantErr: apperror.ErrInvalidID,
		},
		{
			name:     "error draft not found",
			authorID: authorID.Hex(),
			draftID:  primitive.NewObjectID().Hex(),
			mockSetup: func() {
				mockDraftRepo.EXPECT().
					GetByID(gomock.Any(), gomock.Any()).
					Return(nil, mongo.ErrNoDocuments)
			},
			wantErr: mongo.ErrNoDocuments,
		},
		{
			name:     "error not draft author (forbidden)",
			authorID: otherAuthorID.Hex(),
			draftID:  draftID.Hex(),
			mockSetup: func() {
				mockDraftRepo.EXPECT().
					GetByID(gomock.Any(), draftID).
					Return(draftByOtherAuthor, nil)
			},
			wantErr: apperror.ErrForbidden,
		},
		{
			name:     "error repository delete failure",
			authorID: authorID.Hex(),
			draftID:  draftID.Hex(),
			mockSetup: func() {
				mockDraftRepo.EXPECT().
					GetByID(gomock.Any(), draftID).
					Return(existingDraft, nil)

				mockDraftRepo.EXPECT().
					Delete(gomock.Any(), draftID).
					Return(errors.New("database error"))
			},
			wantErr: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := svc.DeleteDraft(ctx, tt.authorID, tt.draftID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.wantErr)
					return
				}
				if tt.wantErr == mongo.ErrNoDocuments {
					if err != mongo.ErrNoDocuments {
						t.Errorf("expected mongo.ErrNoDocuments, got %v", err)
					}
					return
				}
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// Mock PostService for testing PublishDraft
type mockPostService struct {
	createPostFunc func(userID string, req *dto.CreatePostRequest) (*dto.PostResponse, error)
}

func (m *mockPostService) CreatePost(userID string, req *dto.CreatePostRequest) (*dto.PostResponse, error) {
	if m.createPostFunc != nil {
		return m.createPostFunc(userID, req)
	}
	return nil, errors.New("createPostFunc not set")
}

// Implement other PostService methods (not used in tests)
func (m *mockPostService) GetPostByID(postID string, userID string) (*dto.PostResponse, error) {
	return nil, nil
}
func (m *mockPostService) GetPosts(userID string, query *dto.GetPostsQuery) (*dto.PaginatedPostsResponse, error) {
	return nil, nil
}
func (m *mockPostService) GetMyPosts(userID string, query *dto.GetPostsQuery) (*dto.PaginatedPostsResponse, error) {
	return nil, nil
}
func (m *mockPostService) UpdatePost(postID string, userID string, req *dto.UpdatePostRequest) (*dto.PostResponse, error) {
	return nil, nil
}
func (m *mockPostService) DeletePost(postID string, userID string) error {
	return nil
}
func (m *mockPostService) AddImagesToPost(userID, postID string, form *multipart.Form) ([]*model.Image, error) {
	return nil, nil
}
func (m *mockPostService) RemoveImagesFromPost(userID, postID string, publicIDs []string) error {
	return nil
}
func (m *mockPostService) AddVideosToPost(userID, postID string, form *multipart.Form) ([]*model.Video, error) {
	return nil, nil
}
func (m *mockPostService) RemoveVideosFromPost(userID, postID string, publicIDs []string) error {
	return nil
}
func (m *mockPostService) VoteOnPost(userID, postID string, voteValue bool) (*dto.VotesCountResponse, error) {
	return nil, nil
}
func (m *mockPostService) SavePost(userID, postID string) error {
	return nil
}
func (m *mockPostService) UnsavePost(userID, postID string) error {
	return nil
}
func (m *mockPostService) GetSavedPosts(userID string, query *dto.GetPostsQuery) (*dto.PaginatedPostsResponse, error) {
	return nil, nil
}
func (m *mockPostService) ReportPost(reporterID, postID, reason, description string) error {
	return nil
}
func (m *mockPostService) HidePost(userID, postID string) error {
	return nil
}
func (m *mockPostService) UnhidePost(userID, postID string) error {
	return nil
}
func (m *mockPostService) BanPost(postID string, reason *string) error {
	return nil
}
func (m *mockPostService) UnbanPost(postID string) error {
	return nil
}
func (m *mockPostService) GetBanPosts(query *dto.GetBanPostsQuery, requesterID string) (*dto.PaginatedPostsResponse, error) {
	return nil, nil
}
func (m *mockPostService) GetHiddenPosts(userID string, query *dto.GetPostsQuery) (*dto.PaginatedPostsResponse, error) {
	return nil, nil
}
func (m *mockPostService) VoteOnPoll(userID, postID, optionID string) (*dto.PollResponse, error) {
	return nil, nil
}
func (m *mockPostService) RemovePollVote(userID, postID string) (*dto.PollResponse, error) {
	return nil, nil
}
func (m *mockPostService) AddPollOptions(userID, postID string, options []string) (*dto.PollResponse, error) {
	return nil, nil
}
func (m *mockPostService) RemovePollOptions(userID, postID string, optionIDs []string) (*dto.PollResponse, error) {
	return nil, nil
}
func (m *mockPostService) UpdatePoll(userID, postID string, req *dto.UpdatePollRequest) (*dto.PollResponse, error) {
	return nil, nil
}
func (m *mockPostService) UpdatePollOption(userID, postID, optionID string, req *dto.UpdatePollOptionRequest) (*dto.PollResponse, error) {
	return nil, nil
}

// TestPublishDraft tests the PublishDraft method
func TestPublishDraft(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDraftRepo := mocks.NewMockDraftRepo(ctrl)
	mockPostSvc := &mockPostService{}
	svc := NewDraftService(mockDraftRepo, nil, mockPostSvc)

	ctx := context.Background()
	authorID := primitive.NewObjectID()
	otherAuthorID := primitive.NewObjectID()
	draftID := primitive.NewObjectID()
	communityID := "community123"
	title := "Test Draft"
	text := "Draft content"
	postType := model.PostTypeText

	completeDraft := &model.Draft{
		ID:          draftID,
		AuthorID:    authorID,
		CommunityID: &communityID,
		Title:       &title,
		Type:        &postType,
		Content: &model.PostContent{
			Text: text,
		},
	}

	incompleteDraft := &model.Draft{
		ID:       draftID,
		AuthorID: authorID,
		Title:    nil, // Missing title
	}

	draftByOtherAuthor := &model.Draft{
		ID:          draftID,
		AuthorID:    otherAuthorID,
		CommunityID: &communityID,
		Title:       &title,
		Type:        &postType,
	}

	tests := []struct {
		name      string
		authorID  string
		draftID   string
		mockSetup func()
		wantErr   error
		wantPost  bool
	}{
		{
			name:     "success publish draft",
			authorID: authorID.Hex(),
			draftID:  draftID.Hex(),
			mockSetup: func() {
				// Get draft
				mockDraftRepo.EXPECT().
					GetByID(gomock.Any(), draftID).
					Return(completeDraft, nil)

				// Create post - setup manual mock function
				mockPostSvc.createPostFunc = func(userID string, req *dto.CreatePostRequest) (*dto.PostResponse, error) {
					return &dto.PostResponse{
						ID:    primitive.NewObjectID().Hex(),
						Title: title,
					}, nil
				}

				// Delete draft after publish
				mockDraftRepo.EXPECT().
					Delete(gomock.Any(), draftID).
					Return(nil)
			},
			wantErr:  nil,
			wantPost: true,
		},
		{
			name:     "error invalid author id",
			authorID: "invalid-id",
			draftID:  draftID.Hex(),
			mockSetup: func() {
				// No mock needed
			},
			wantErr:  apperror.ErrInvalidID,
			wantPost: false,
		},
		{
			name:     "error draft not found",
			authorID: authorID.Hex(),
			draftID:  primitive.NewObjectID().Hex(),
			mockSetup: func() {
				mockDraftRepo.EXPECT().
					GetByID(gomock.Any(), gomock.Any()).
					Return(nil, mongo.ErrNoDocuments)
			},
			wantErr:  mongo.ErrNoDocuments,
			wantPost: false,
		},
		{
			name:     "error not draft author (forbidden)",
			authorID: otherAuthorID.Hex(),
			draftID:  draftID.Hex(),
			mockSetup: func() {
				mockDraftRepo.EXPECT().
					GetByID(gomock.Any(), draftID).
					Return(draftByOtherAuthor, nil)
			},
			wantErr:  apperror.ErrForbidden,
			wantPost: false,
		},
		{
			name:     "error missing required fields",
			authorID: authorID.Hex(),
			draftID:  draftID.Hex(),
			mockSetup: func() {
				mockDraftRepo.EXPECT().
					GetByID(gomock.Any(), draftID).
					Return(incompleteDraft, nil)
			},
			wantErr:  apperror.ErrBadRequest,
			wantPost: false,
		},
		{
			name:     "error post creation failure",
			authorID: authorID.Hex(),
			draftID:  draftID.Hex(),
			mockSetup: func() {
				mockDraftRepo.EXPECT().
					GetByID(gomock.Any(), draftID).
					Return(completeDraft, nil)

				// Setup manual mock to return error
				mockPostSvc.createPostFunc = func(userID string, req *dto.CreatePostRequest) (*dto.PostResponse, error) {
					return nil, errors.New("post creation failed")
				}
			},
			wantErr:  errors.New("post creation failed"),
			wantPost: false,
		},
		{
			name:     "success publish even if draft delete fails",
			authorID: authorID.Hex(),
			draftID:  draftID.Hex(),
			mockSetup: func() {
				mockDraftRepo.EXPECT().
					GetByID(gomock.Any(), draftID).
					Return(completeDraft, nil)

				// Setup manual mock to return success
				mockPostSvc.createPostFunc = func(userID string, req *dto.CreatePostRequest) (*dto.PostResponse, error) {
					return &dto.PostResponse{
						ID:    primitive.NewObjectID().Hex(),
						Title: title,
					}, nil
				}

				// Delete fails but shouldn't fail the whole operation
				mockDraftRepo.EXPECT().
					Delete(gomock.Any(), draftID).
					Return(errors.New("delete failed"))
			},
			wantErr:  nil,
			wantPost: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			result, err := svc.PublishDraft(ctx, tt.authorID, tt.draftID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.wantErr)
					return
				}
				// Handle mongo error specially
				if tt.wantErr == mongo.ErrNoDocuments {
					if err != mongo.ErrNoDocuments {
						t.Errorf("expected mongo.ErrNoDocuments, got %v", err)
					}
					return
				}
				// Check if error contains expected code
				if tt.wantErr == apperror.ErrBadRequest {
					appErr, ok := err.(*apperror.AppError)
					if !ok || appErr.Code != apperror.ErrBadRequest.Code {
						t.Errorf("expected ErrBadRequest, got %v", err)
					}
					return
				}
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if tt.wantPost {
				if result == nil {
					t.Error("expected post response, got nil")
				}
			}
		})
	}
}
