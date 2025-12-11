package service

import (
	"errors"
	"testing"
	"time"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/platform/bus"
	"github.com/giakiet05/lkforum/internal/repo/mocks"
	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Manual VoteService Mock
type MockVoteService struct {
	ctrl     *gomock.Controller
	recorder *MockVoteServiceRecorder
}

type MockVoteServiceRecorder struct {
	mock *MockVoteService
}

func NewMockVoteService(ctrl *gomock.Controller) *MockVoteService {
	mock := &MockVoteService{ctrl: ctrl}
	mock.recorder = &MockVoteServiceRecorder{mock}
	return mock
}

func (m *MockVoteService) EXPECT() *MockVoteServiceRecorder {
	return m.recorder
}

func (m *MockVoteService) VoteOnTarget(userID, targetID string, targetType model.VoteTargetType, voteValue bool) (*dto.VotesCountResponse, error) {
	args := m.ctrl.Call(m, "VoteOnTarget", userID, targetID, targetType, voteValue)
	ret0, _ := args[0].(*dto.VotesCountResponse)
	ret1, _ := args[1].(error)
	return ret0, ret1
}

func (m *MockVoteService) GetUserVote(userID, targetID string, targetType model.VoteTargetType) (*model.Vote, error) {
	args := m.ctrl.Call(m, "GetUserVote", userID, targetID, targetType)
	ret0, _ := args[0].(*model.Vote)
	ret1, _ := args[1].(error)
	return ret0, ret1
}

func (m *MockVoteService) FindUserVotes(userID string, targetIDs []string, targetType model.VoteTargetType) (map[string]string, error) {
	args := m.ctrl.Call(m, "FindUserVotes", userID, targetIDs, targetType)
	ret0, _ := args[0].(map[string]string)
	ret1, _ := args[1].(error)
	return ret0, ret1
}

func (r *MockVoteServiceRecorder) VoteOnTarget(userID, targetID, targetType, voteValue interface{}) *gomock.Call {
	return r.mock.ctrl.RecordCall(r.mock, "VoteOnTarget", userID, targetID, targetType, voteValue)
}

func (r *MockVoteServiceRecorder) GetUserVote(userID, targetID, targetType interface{}) *gomock.Call {
	return r.mock.ctrl.RecordCall(r.mock, "GetUserVote", userID, targetID, targetType)
}

func (r *MockVoteServiceRecorder) FindUserVotes(userID, targetIDs, targetType interface{}) *gomock.Call {
	return r.mock.ctrl.RecordCall(r.mock, "FindUserVotes", userID, targetIDs, targetType)
}

// Test CreatePost - 7 test cases
func TestCreatePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup ALL mocks
	mockPostRepo := mocks.NewMockPostRepo(ctrl)
	mockVoteService := NewMockVoteService(ctrl)
	mockEventBus := bus.NewMockEventBus(ctrl)
	mockUserRepo := mocks.NewMockUserRepo(ctrl)
	mockCommunityRepo := mocks.NewMockCommunityRepo(ctrl)
	mockPollVoteRepo := mocks.NewMockPollVoteRepo(ctrl)
	mockMembershipRepo := mocks.NewMockMembershipRepo(ctrl)
	mockSavedPostRepo := mocks.NewMockSavedPostRepo(ctrl)
	mockReportRepo := mocks.NewMockReportRepo(ctrl)

	// Create service with ALL dependencies
	svc := NewPostService(
		mockPostRepo,
		mockVoteService,
		mockPollVoteRepo,
		mockUserRepo,
		mockCommunityRepo,
		mockMembershipRepo,
		mockSavedPostRepo,
		mockReportRepo,
		mockEventBus,
	)

	// Test data
	userID := primitive.NewObjectID()
	communityID := primitive.NewObjectID()
	postID := primitive.NewObjectID()

	createdPost := &model.Post{
		ID:          postID,
		AuthorID:    userID,
		CommunityID: communityID,
		Type:        model.PostTypeText,
		Title:       "Test Post",
		Content:     &model.PostContent{Text: "Test content"},
		VotesCount:  &model.VotesCount{Up: 0, Down: 0},
		CreatedAt:   time.Now(),
	}

	tests := []struct {
		name      string
		userID    string
		request   *dto.CreatePostRequest
		mockSetup func()
		wantErr   error
	}{
		{
			name:   "success_text_post",
			userID: userID.Hex(),
			request: &dto.CreatePostRequest{
				CommunityID: communityID.Hex(),
				Type:        model.PostTypeText,
				Title:       "Test Post",
				Text:        "Test content",
			},
			mockSetup: func() {
				// Mock post creation
				mockPostRepo.EXPECT().
					Create(gomock.Any(), gomock.AssignableToTypeOf(&model.Post{})).
					DoAndReturn(func(ctx interface{}, post *model.Post) (*model.Post, error) {
						// Verify post structure
						if post.AuthorID != userID {
							t.Errorf("Expected AuthorID %v, got %v", userID, post.AuthorID)
						}
						if post.Title != "Test Post" {
							t.Errorf("Expected title 'Test Post', got %v", post.Title)
						}
						return createdPost, nil
					})

				// Mock async vote service call
				mockVoteService.EXPECT().
					VoteOnTarget(userID.Hex(), postID.Hex(), model.VoteTargetPost, true).
					Return(&dto.VotesCountResponse{Up: 1, Down: 0}, nil).
					AnyTimes()

				// Mock event publishing
				mockEventBus.EXPECT().
					Publish(gomock.AssignableToTypeOf(&bus.PostCreatedEvent{})).
					Do(func(event interface{}) {
						if postEvent, ok := event.(*bus.PostCreatedEvent); ok {
							if postEvent.PostID != postID.Hex() {
								t.Errorf("Expected PostID %s, got %s", postID.Hex(), postEvent.PostID)
							}
						}
					})

				// Mock GetPostByID call (called at end)
				mockPostRepo.EXPECT().
					GetByID(gomock.Any(), postID.Hex()).
					Return(createdPost, nil)

				// Mock UserRepo and CommunityRepo calls in GetPostByID
				mockUserRepo.EXPECT().
					GetByID(gomock.Any(), userID.Hex()).
					Return(&model.User{ID: userID, Username: "testuser"}, nil).
					AnyTimes()

				mockCommunityRepo.EXPECT().
					GetByID(gomock.Any(), communityID.Hex()).
					Return(&model.Community{ID: communityID, Name: "testcommunity"}, nil).
					AnyTimes()

				// No PollVoteRepo mock needed for text posts

				mockVoteService.EXPECT().
					GetUserVote(userID.Hex(), postID.Hex(), model.VoteTargetPost).
					Return(nil, nil).
					AnyTimes()
			},
			wantErr: nil,
		},
		{
			name:   "success_poll_post",
			userID: userID.Hex(),
			request: &dto.CreatePostRequest{
				CommunityID: communityID.Hex(),
				Type:        model.PostTypePoll,
				Title:       "Test Poll",
				Poll: &dto.CreatePollRequest{
					Question: "What's your favorite color?",
					Options:  []string{"Red", "Blue", "Green"},
				},
			},
			mockSetup: func() {
				pollPost := *createdPost
				pollPost.Type = model.PostTypePoll
				pollPost.Content = &model.PostContent{
					Poll: &model.Poll{
						Question: "What's your favorite color?",
						Options: []model.PollOption{
							{ID: "opt1", Text: "Red", Votes: 0},
							{ID: "opt2", Text: "Blue", Votes: 0},
							{ID: "opt3", Text: "Green", Votes: 0},
						},
					},
				}

				mockPostRepo.EXPECT().
					Create(gomock.Any(), gomock.AssignableToTypeOf(&model.Post{})).
					DoAndReturn(func(ctx interface{}, post *model.Post) (*model.Post, error) {
						if post.Type != model.PostTypePoll {
							t.Errorf("Expected PostTypePoll, got %v", post.Type)
						}
						if post.Content.Poll == nil {
							t.Error("Expected poll content, got nil")
						}
						return &pollPost, nil
					})

				mockVoteService.EXPECT().VoteOnTarget(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
				mockEventBus.EXPECT().Publish(gomock.Any())
				mockPostRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(&pollPost, nil)

				// Mock UserRepo and CommunityRepo
				mockUserRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(&model.User{ID: userID, Username: "testuser"}, nil).AnyTimes()
				mockCommunityRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(&model.Community{ID: communityID, Name: "testcommunity"}, nil).AnyTimes()
				// No PollVoteRepo mock needed for text posts
				mockVoteService.EXPECT().GetUserVote(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
			},
			wantErr: nil,
		},
		{
			name:   "error_invalid_user_id",
			userID: "invalid-user-id",
			request: &dto.CreatePostRequest{
				CommunityID: communityID.Hex(),
				Type:        model.PostTypeText,
				Title:       "Test Post",
			},
			mockSetup: func() {
				// No mocks needed - should fail validation
			},
			wantErr: apperror.ErrInvalidID,
		},
		{
			name:   "error_invalid_community_id",
			userID: userID.Hex(),
			request: &dto.CreatePostRequest{
				CommunityID: "invalid-community-id",
				Type:        model.PostTypeText,
				Title:       "Test Post",
			},
			mockSetup: func() {
				// No mocks needed - should fail validation
			},
			wantErr: apperror.ErrInvalidID,
		},
		{
			name:   "error_empty_title",
			userID: userID.Hex(),
			request: &dto.CreatePostRequest{
				CommunityID: communityID.Hex(),
				Type:        model.PostTypeText,
				Title:       "", // Empty title
			},
			mockSetup: func() {
				// This would be caught by validation in the controller layer
				// For service layer, we assume valid input structure
				mockPostRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("validation error: title cannot be empty"))
			},
			wantErr: errors.New("validation error: title cannot be empty"),
		},
		{
			name:   "error_repository_failure",
			userID: userID.Hex(),
			request: &dto.CreatePostRequest{
				CommunityID: communityID.Hex(),
				Type:        model.PostTypeText,
				Title:       "Test Post",
			},
			mockSetup: func() {
				mockPostRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("database connection error"))
			},
			wantErr: errors.New("database connection error"),
		},
		{
			name:   "error_vote_service_failure",
			userID: userID.Hex(),
			request: &dto.CreatePostRequest{
				CommunityID: communityID.Hex(),
				Type:        model.PostTypeText,
				Title:       "Test Post",
			},
			mockSetup: func() {
				// Post creation succeeds
				mockPostRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(createdPost, nil)

				// Vote service fails (but doesn't prevent post creation)
				mockVoteService.EXPECT().
					VoteOnTarget(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("vote service down")).
					AnyTimes()

				// Event still publishes
				mockEventBus.EXPECT().Publish(gomock.Any())

				// GetPostByID still works (vote failure handled gracefully)
				mockPostRepo.EXPECT().
					GetByID(gomock.Any(), gomock.Any()).
					Return(createdPost, nil)

				// Mock UserRepo and CommunityRepo
				mockUserRepo.EXPECT().
					GetByID(gomock.Any(), gomock.Any()).
					Return(&model.User{ID: userID, Username: "testuser"}, nil).
					AnyTimes()

				mockCommunityRepo.EXPECT().
					GetByID(gomock.Any(), gomock.Any()).
					Return(&model.Community{ID: communityID, Name: "testcommunity"}, nil).
					AnyTimes()

				// Mock UserRepo and CommunityRepo
				mockUserRepo.EXPECT().
					GetByID(gomock.Any(), gomock.Any()).
					Return(&model.User{ID: userID, Username: "testuser"}, nil).
					AnyTimes()

				mockCommunityRepo.EXPECT().
					GetByID(gomock.Any(), gomock.Any()).
					Return(&model.Community{ID: communityID, Name: "testcommunity"}, nil).
					AnyTimes()

				// No PollVoteRepo mock needed for text posts

				mockVoteService.EXPECT().
					GetUserVote(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("vote service down")).
					AnyTimes()
			},
			wantErr: nil, // Post creation should succeed even if vote fails
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			result, err := svc.CreatePost(tt.userID, tt.request)

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
				return
			}

			// Basic validation
			if result.ID == "" {
				t.Error("expected post ID, got empty string")
			}
		})
	}
}

// Test GetPostByID - 6 test cases
func TestGetPostByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := mocks.NewMockPostRepo(ctrl)
	mockVoteService := NewMockVoteService(ctrl)
	mockUserRepo := mocks.NewMockUserRepo(ctrl)
	mockCommunityRepo := mocks.NewMockCommunityRepo(ctrl)
	mockPollVoteRepo := mocks.NewMockPollVoteRepo(ctrl)

	svc := NewPostService(
		mockPostRepo,      // postRepo
		mockVoteService,   // voteService
		mockPollVoteRepo,  // pollVoteRepo
		mockUserRepo,      // userRepo
		mockCommunityRepo, // communityRepo
		nil,               // membershipRepo
		nil,               // savedPostRepo
		nil,               // reportRepo
		nil,               // bus
	)

	postID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	existingPost := &model.Post{
		ID:       postID,
		AuthorID: primitive.NewObjectID(),
		Title:    "Existing Post",
		Type:     model.PostTypeText, // NOT poll type
		Content:  &model.PostContent{Text: "Post content"},
	}

	userVote := &model.Vote{
		UserID:     userID,
		TargetID:   postID,
		TargetType: model.VoteTargetPost,
		Value:      true,
	}

	tests := []struct {
		name      string
		postID    string
		userID    string
		mockSetup func()
		wantErr   error
	}{
		{
			name:   "success_with_user_vote",
			postID: postID.Hex(),
			userID: userID.Hex(),
			mockSetup: func() {
				mockPostRepo.EXPECT().
					GetByID(gomock.Any(), postID.Hex()).
					Return(existingPost, nil)

				mockVoteService.EXPECT().
					GetUserVote(userID.Hex(), postID.Hex(), model.VoteTargetPost).
					Return(userVote, nil)

				// Mock UserRepo and CommunityRepo
				mockUserRepo.EXPECT().
					GetByID(gomock.Any(), gomock.Any()).
					Return(&model.User{ID: userID, Username: "testuser"}, nil)

				mockCommunityRepo.EXPECT().
					GetByID(gomock.Any(), gomock.Any()).
					Return(&model.Community{ID: primitive.NewObjectID(), Name: "testcommunity"}, nil)

				// No PollVoteRepo mock needed for text posts
			},
			wantErr: nil,
		},
		{
			name:   "success_without_user_vote",
			postID: postID.Hex(),
			userID: "", // Anonymous user
			mockSetup: func() {
				mockPostRepo.EXPECT().
					GetByID(gomock.Any(), postID.Hex()).
					Return(existingPost, nil)

				// Mock UserRepo and CommunityRepo for anonymous user
				mockUserRepo.EXPECT().
					GetByID(gomock.Any(), gomock.Any()).
					Return(&model.User{Username: "author"}, nil)

				mockCommunityRepo.EXPECT().
					GetByID(gomock.Any(), gomock.Any()).
					Return(&model.Community{Name: "testcommunity"}, nil)

				// Mock PollVoteRepo (not called for anonymous but needed for service)
				// No expectations needed since userID is empty

				// No vote service call for anonymous users
			},
			wantErr: nil,
		},
		{
			name:   "error_invalid_post_id",
			postID: "invalid-post-id",
			userID: userID.Hex(),
			mockSetup: func() {
				// Mock repo to return invalid ID error
				mockPostRepo.EXPECT().
					GetByID(gomock.Any(), "invalid-post-id").
					Return(nil, apperror.ErrInvalidID)
			},
			wantErr: apperror.ErrInvalidID,
		},
		{
			name:   "error_post_not_found",
			postID: primitive.NewObjectID().Hex(),
			userID: userID.Hex(),
			mockSetup: func() {
				mockPostRepo.EXPECT().
					GetByID(gomock.Any(), gomock.Any()).
					Return(nil, mongo.ErrNoDocuments)
			},
			wantErr: mongo.ErrNoDocuments,
		},
		{
			name:   "error_repository_failure",
			postID: postID.Hex(),
			userID: userID.Hex(),
			mockSetup: func() {
				mockPostRepo.EXPECT().
					GetByID(gomock.Any(), postID.Hex()).
					Return(nil, errors.New("database connection error"))
			},
			wantErr: errors.New("database connection error"),
		},
		{
			name:   "error_vote_service_failure_graceful",
			postID: postID.Hex(),
			userID: userID.Hex(),
			mockSetup: func() {
				mockPostRepo.EXPECT().
					GetByID(gomock.Any(), postID.Hex()).
					Return(existingPost, nil)

				mockVoteService.EXPECT().
					GetUserVote(userID.Hex(), postID.Hex(), model.VoteTargetPost).
					Return(nil, errors.New("vote service down"))

				// Mock UserRepo and CommunityRepo
				mockUserRepo.EXPECT().
					GetByID(gomock.Any(), gomock.Any()).
					Return(&model.User{Username: "author"}, nil)

				mockCommunityRepo.EXPECT().
					GetByID(gomock.Any(), gomock.Any()).
					Return(&model.Community{Name: "testcommunity"}, nil)

				mockPollVoteRepo.EXPECT().
					GetUserVoteIDs(gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]string{}, nil)

				// VoteService should fail gracefully
				mockVoteService.EXPECT().
					GetUserVote(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("vote service down"))

				// Should handle vote service error gracefully
			},
			wantErr: nil, // Should succeed even if vote service fails
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			result, err := svc.GetPostByID(tt.postID, tt.userID)

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

// Test UpdatePost - 5 test cases
func TestUpdatePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup ALL mocks
	mockPostRepo := mocks.NewMockPostRepo(ctrl)
	mockVoteService := NewMockVoteService(ctrl)
	mockUserRepo := mocks.NewMockUserRepo(ctrl)
	mockCommunityRepo := mocks.NewMockCommunityRepo(ctrl)
	mockPollVoteRepo := mocks.NewMockPollVoteRepo(ctrl)
	mockEventBus := bus.NewMockEventBus(ctrl)

	svc := NewPostService(
		mockPostRepo,
		mockVoteService,
		mockPollVoteRepo,
		mockUserRepo,
		mockCommunityRepo,
		nil, nil, nil,
		mockEventBus, // Fix nil bus
	)

	postID := primitive.NewObjectID()
	userID := primitive.NewObjectID()
	otherUserID := primitive.NewObjectID()

	existingPost := &model.Post{
		ID:       postID,
		AuthorID: userID,
		Title:    "Original Title",
		Content:  &model.PostContent{Text: "Original text"},
		Type:     model.PostTypeText,
	}

	postByOtherUser := &model.Post{
		ID:       postID,
		AuthorID: otherUserID, // Different author
		Title:    "Other User's Post",
		Type:     model.PostTypeText,
	}

	ptrStr := func(s string) *string { return &s }

	tests := []struct {
		name      string
		postID    string
		userID    string
		request   *dto.UpdatePostRequest
		mockSetup func()
		wantErr   error
	}{
		{
			name:   "success_update_title",
			postID: postID.Hex(),
			userID: userID.Hex(),
			request: &dto.UpdatePostRequest{
				Title: ptrStr("Updated Title"),
			},
			mockSetup: func() {
				mockPostRepo.EXPECT().
					GetByID(gomock.Any(), postID.Hex()).
					Return(existingPost, nil)

				mockPostRepo.EXPECT().
					UpdateByID(gomock.Any(), postID.Hex(), gomock.Any()).
					Return(nil)

				// Mock the GetPostByID call at the end
				updatedPost := *existingPost
				updatedPost.Title = "Updated Title"
				mockPostRepo.EXPECT().
					GetByID(gomock.Any(), postID.Hex()).
					Return(&updatedPost, nil)

				// Mock UserRepo and CommunityRepo for GetPostByID calls (called twice)
				mockUserRepo.EXPECT().
					GetByID(gomock.Any(), gomock.Any()).
					Return(&model.User{ID: userID, Username: "testuser"}, nil).
					Times(2)

				mockCommunityRepo.EXPECT().
					GetByID(gomock.Any(), gomock.Any()).
					Return(&model.Community{ID: primitive.NewObjectID(), Name: "testcommunity"}, nil).
					Times(2)

				// No PollVoteRepo mock needed for text posts

				mockVoteService.EXPECT().
					GetUserVote(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, nil).
					AnyTimes()

				// Mock EventBus for PostUpdatedEvent
				mockEventBus.EXPECT().
					Publish(gomock.Any()).
					Return()
			},
			wantErr: nil,
		},
		{
			name:   "error_invalid_post_id",
			postID: "invalid-post-id",
			userID: userID.Hex(),
			request: &dto.UpdatePostRequest{
				Title: ptrStr("Updated Title"),
			},
			mockSetup: func() {
				// No mocks needed - should fail validation
			},
			wantErr: apperror.ErrInvalidID,
		},
		{
			name:   "error_not_post_author",
			postID: postID.Hex(),
			userID: otherUserID.Hex(),
			request: &dto.UpdatePostRequest{
				Title: ptrStr("Updated Title"),
			},
			mockSetup: func() {
				mockPostRepo.EXPECT().
					GetByID(gomock.Any(), postID.Hex()).
					Return(postByOtherUser, nil)
			},
			wantErr: apperror.ErrForbidden,
		},
		{
			name:   "error_post_not_found",
			postID: primitive.NewObjectID().Hex(),
			userID: userID.Hex(),
			request: &dto.UpdatePostRequest{
				Title: ptrStr("Updated Title"),
			},
			mockSetup: func() {
				mockPostRepo.EXPECT().
					GetByID(gomock.Any(), gomock.Any()).
					Return(nil, mongo.ErrNoDocuments)
			},
			wantErr: mongo.ErrNoDocuments,
		},
		{
			name:   "error_repository_failure",
			postID: postID.Hex(),
			userID: userID.Hex(),
			request: &dto.UpdatePostRequest{
				Title: ptrStr("Updated Title"),
			},
			mockSetup: func() {
				mockPostRepo.EXPECT().
					GetByID(gomock.Any(), postID.Hex()).
					Return(existingPost, nil)

				mockPostRepo.EXPECT().
					UpdateByID(gomock.Any(), postID.Hex(), gomock.Any()).
					Return(errors.New("database connection error"))
			},
			wantErr: errors.New("database connection error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			result, err := svc.UpdatePost(tt.postID, tt.userID, tt.request)

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

// Test DeletePost - 4 test cases
func TestDeletePost(t *testing.T) {
	tests := []struct {
		name      string
		postID    string
		userID    string
		mockSetup func(*mocks.MockPostRepo, string, string)
		wantErr   error
	}{
		{
			name:   "success_delete_own_post",
			postID: primitive.NewObjectID().Hex(),
			userID: primitive.NewObjectID().Hex(),
			mockSetup: func(mockPostRepo *mocks.MockPostRepo, postID, userID string) {
				authorID, _ := primitive.ObjectIDFromHex(userID)
				postObjectID, _ := primitive.ObjectIDFromHex(postID)

				existingPost := &model.Post{
					ID:       postObjectID,
					AuthorID: authorID,
					Title:    "Test Post",
					Type:     model.PostTypeText,
				}

				mockPostRepo.EXPECT().
					GetByID(gomock.Any(), postID).
					Return(existingPost, nil)

				mockPostRepo.EXPECT().
					Delete(gomock.Any(), postID).
					Return(nil)
			},
			wantErr: nil,
		},
		{
			name:   "error_not_post_author",
			postID: primitive.NewObjectID().Hex(),
			userID: primitive.NewObjectID().Hex(),
			mockSetup: func(mockPostRepo *mocks.MockPostRepo, postID, userID string) {
				postObjectID, _ := primitive.ObjectIDFromHex(postID)
				otherUserID := primitive.NewObjectID() // Different user

				postByOtherUser := &model.Post{
					ID:       postObjectID,
					AuthorID: otherUserID,
					Title:    "Other User's Post",
					Type:     model.PostTypeText,
				}

				mockPostRepo.EXPECT().
					GetByID(gomock.Any(), postID).
					Return(postByOtherUser, nil)
			},
			wantErr: apperror.ErrForbidden,
		},
		{
			name:   "error_post_not_found",
			postID: primitive.NewObjectID().Hex(),
			userID: primitive.NewObjectID().Hex(),
			mockSetup: func(mockPostRepo *mocks.MockPostRepo, postID, userID string) {
				mockPostRepo.EXPECT().
					GetByID(gomock.Any(), postID).
					Return(nil, mongo.ErrNoDocuments)
			},
			wantErr: mongo.ErrNoDocuments,
		},
		{
			name:   "error_repository_failure",
			postID: primitive.NewObjectID().Hex(),
			userID: primitive.NewObjectID().Hex(),
			mockSetup: func(mockPostRepo *mocks.MockPostRepo, postID, userID string) {
				authorID, _ := primitive.ObjectIDFromHex(userID)
				postObjectID, _ := primitive.ObjectIDFromHex(postID)

				existingPost := &model.Post{
					ID:       postObjectID,
					AuthorID: authorID,
					Title:    "Test Post",
					Type:     model.PostTypeText,
				}

				mockPostRepo.EXPECT().
					GetByID(gomock.Any(), postID).
					Return(existingPost, nil)

				mockPostRepo.EXPECT().
					Delete(gomock.Any(), postID).
					Return(errors.New("database connection error"))
			},
			wantErr: errors.New("database connection error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPostRepo := mocks.NewMockPostRepo(ctrl)

			svc := NewPostService(
				mockPostRepo,
				nil, // voteService not needed for delete
				nil, nil, nil, nil, nil, nil, nil,
			)

			tt.mockSetup(mockPostRepo, tt.postID, tt.userID)

			err := svc.DeletePost(tt.postID, tt.userID)

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
