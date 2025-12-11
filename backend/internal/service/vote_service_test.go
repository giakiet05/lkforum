package service

import (
	"errors"
	"testing"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/platform/bus"
	"github.com/giakiet05/lkforum/internal/repo/mocks"
	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Mock EventBus for testing
type mockEventBus struct{}

func (m *mockEventBus) Subscribe(topic string, ch bus.EventListener) {}
func (m *mockEventBus) Publish(event bus.Event)                      {}

// TestVoteOnPost tests voting on posts
func TestVoteOnPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVoteRepo := mocks.NewMockVoteRepo(ctrl)
	mockPostRepo := mocks.NewMockPostRepo(ctrl)
	mockCommentRepo := mocks.NewMockCommentRepo(ctrl)
	mockEventBus := &mockEventBus{}
	svc := NewVoteService(mockVoteRepo, mockPostRepo, mockCommentRepo, mockEventBus)

	userID := primitive.NewObjectID()
	authorID := primitive.NewObjectID()
	postID := primitive.NewObjectID()

	approvedPost := &model.Post{
		ID:               postID,
		AuthorID:         authorID,
		ModerationStatus: model.ModerationApproved,
		VotesCount: &model.VotesCount{
			Up:   10,
			Down: 2,
		},
	}

	pendingPost := &model.Post{
		ID:               postID,
		AuthorID:         authorID,
		ModerationStatus: model.ModerationPending,
		VotesCount:       &model.VotesCount{Up: 0, Down: 0},
	}

	tests := []struct {
		name      string
		userID    string
		postID    string
		voteValue bool
		mockSetup func()
		wantErr   error
		wantScore int
	}{
		{
			name:      "success upvote on post (new vote)",
			userID:    userID.Hex(),
			postID:    postID.Hex(),
			voteValue: true,
			mockSetup: func() {
				// GetByID will be called 3 times (check moderation, get author, get updated counts)
				// First 2 calls return original post, last call returns updated post
				gomock.InOrder(
					mockPostRepo.EXPECT().
						GetByID(gomock.Any(), postID.Hex()).
						Return(approvedPost, nil),

					mockPostRepo.EXPECT().
						GetByID(gomock.Any(), postID.Hex()).
						Return(approvedPost, nil),

					mockPostRepo.EXPECT().
						GetByID(gomock.Any(), postID.Hex()).
						DoAndReturn(func(ctx interface{}, id string) (*model.Post, error) {
							updatedPost := *approvedPost
							updatedPost.VotesCount = &model.VotesCount{Up: 11, Down: 2}
							return &updatedPost, nil
						}),
				)

				// Get previous vote (no previous vote)
				mockVoteRepo.EXPECT().
					GetUserVote(gomock.Any(), userID.Hex(), postID.Hex(), model.VoteTargetPost).
					Return(nil, mongo.ErrNoDocuments)

				// Perform vote
				mockVoteRepo.EXPECT().
					Vote(gomock.Any(), userID.Hex(), postID.Hex(), model.VoteTargetPost, true).
					Return(nil)
			},
			wantErr:   nil,
			wantScore: 9, // 11 - 2
		},
		{
			name:      "success downvote on post",
			userID:    userID.Hex(),
			postID:    postID.Hex(),
			voteValue: false,
			mockSetup: func() {
				gomock.InOrder(
					mockPostRepo.EXPECT().
						GetByID(gomock.Any(), postID.Hex()).
						Return(approvedPost, nil),

					mockPostRepo.EXPECT().
						GetByID(gomock.Any(), postID.Hex()).
						Return(approvedPost, nil),

					mockPostRepo.EXPECT().
						GetByID(gomock.Any(), postID.Hex()).
						DoAndReturn(func(ctx interface{}, id string) (*model.Post, error) {
							updatedPost := *approvedPost
							updatedPost.VotesCount = &model.VotesCount{Up: 10, Down: 3}
							return &updatedPost, nil
						}),
				)

				mockVoteRepo.EXPECT().
					GetUserVote(gomock.Any(), userID.Hex(), postID.Hex(), model.VoteTargetPost).
					Return(nil, mongo.ErrNoDocuments)

				mockVoteRepo.EXPECT().
					Vote(gomock.Any(), userID.Hex(), postID.Hex(), model.VoteTargetPost, false).
					Return(nil)
			},
			wantErr:   nil,
			wantScore: 7, // 10 - 3
		},
		{
			name:      "success change vote from up to down",
			userID:    userID.Hex(),
			postID:    postID.Hex(),
			voteValue: false,
			mockSetup: func() {
				gomock.InOrder(
					mockPostRepo.EXPECT().
						GetByID(gomock.Any(), postID.Hex()).
						Return(approvedPost, nil),

					mockPostRepo.EXPECT().
						GetByID(gomock.Any(), postID.Hex()).
						Return(approvedPost, nil),

					mockPostRepo.EXPECT().
						GetByID(gomock.Any(), postID.Hex()).
						DoAndReturn(func(ctx interface{}, id string) (*model.Post, error) {
							updatedPost := *approvedPost
							updatedPost.VotesCount = &model.VotesCount{Up: 9, Down: 3}
							return &updatedPost, nil
						}),
				)

				// Previous vote was upvote
				prevVote := &model.Vote{
					UserID:     userID,
					TargetID:   postID,
					TargetType: model.VoteTargetPost,
					Value:      true,
				}
				mockVoteRepo.EXPECT().
					GetUserVote(gomock.Any(), userID.Hex(), postID.Hex(), model.VoteTargetPost).
					Return(prevVote, nil)

				mockVoteRepo.EXPECT().
					Vote(gomock.Any(), userID.Hex(), postID.Hex(), model.VoteTargetPost, false).
					Return(nil)
			},
			wantErr:   nil,
			wantScore: 6, // 9 - 3
		},
		{
			name:      "error post not approved",
			userID:    userID.Hex(),
			postID:    postID.Hex(),
			voteValue: true,
			mockSetup: func() {
				// Post is pending moderation
				mockPostRepo.EXPECT().
					GetByID(gomock.Any(), postID.Hex()).
					Return(pendingPost, nil)
			},
			wantErr:   apperror.ErrForbidden,
			wantScore: 0,
		},
		{
			name:      "error post not found",
			userID:    userID.Hex(),
			postID:    primitive.NewObjectID().Hex(),
			voteValue: true,
			mockSetup: func() {
				mockPostRepo.EXPECT().
					GetByID(gomock.Any(), gomock.Any()).
					Return(nil, mongo.ErrNoDocuments)
			},
			wantErr:   mongo.ErrNoDocuments,
			wantScore: 0,
		},
		{
			name:      "error vote repository failure",
			userID:    userID.Hex(),
			postID:    postID.Hex(),
			voteValue: true,
			mockSetup: func() {
				mockPostRepo.EXPECT().
					GetByID(gomock.Any(), postID.Hex()).
					Return(approvedPost, nil)

				mockPostRepo.EXPECT().
					GetByID(gomock.Any(), postID.Hex()).
					Return(approvedPost, nil)

				mockVoteRepo.EXPECT().
					GetUserVote(gomock.Any(), userID.Hex(), postID.Hex(), model.VoteTargetPost).
					Return(nil, mongo.ErrNoDocuments)

				mockVoteRepo.EXPECT().
					Vote(gomock.Any(), userID.Hex(), postID.Hex(), model.VoteTargetPost, true).
					Return(errors.New("database error"))
			},
			wantErr:   errors.New("database error"),
			wantScore: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			result, err := svc.VoteOnTarget(tt.userID, tt.postID, model.VoteTargetPost, tt.voteValue)

			// Verify error
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

			// Verify success
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Error("expected result, got nil")
				return
			}

			if result.Score != tt.wantScore {
				t.Errorf("expected score %d, got %d", tt.wantScore, result.Score)
			}
		})
	}
}

// TestVoteOnComment tests voting on comments
func TestVoteOnComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVoteRepo := mocks.NewMockVoteRepo(ctrl)
	mockPostRepo := mocks.NewMockPostRepo(ctrl)
	mockCommentRepo := mocks.NewMockCommentRepo(ctrl)
	mockEventBus := &mockEventBus{}
	svc := NewVoteService(mockVoteRepo, mockPostRepo, mockCommentRepo, mockEventBus)

	userID := primitive.NewObjectID()
	authorID := primitive.NewObjectID()
	commentID := primitive.NewObjectID()

	comment := &model.Comment{
		ID: commentID,
		Author: model.CommentAuthor{
			ID: authorID,
		},
		VotesCount: &model.VotesCount{
			Up:   5,
			Down: 1,
		},
	}

	tests := []struct {
		name      string
		userID    string
		commentID string
		voteValue bool
		mockSetup func()
		wantErr   error
		wantScore int
	}{
		{
			name:      "success upvote on comment",
			userID:    userID.Hex(),
			commentID: commentID.Hex(),
			voteValue: true,
			mockSetup: func() {
				// GetByID will be called 3 times
				gomock.InOrder(
					mockCommentRepo.EXPECT().
						GetByID(gomock.Any(), commentID.Hex()).
						Return(comment, nil),

					mockCommentRepo.EXPECT().
						GetByID(gomock.Any(), commentID.Hex()).
						Return(comment, nil),

					mockCommentRepo.EXPECT().
						GetByID(gomock.Any(), commentID.Hex()).
						DoAndReturn(func(ctx interface{}, id string) (*model.Comment, error) {
							updatedComment := *comment
							updatedComment.VotesCount = &model.VotesCount{Up: 6, Down: 1}
							return &updatedComment, nil
						}),
				)

				// Get previous vote
				mockVoteRepo.EXPECT().
					GetUserVote(gomock.Any(), userID.Hex(), commentID.Hex(), model.VoteTargetComment).
					Return(nil, mongo.ErrNoDocuments)

				// Perform vote
				mockVoteRepo.EXPECT().
					Vote(gomock.Any(), userID.Hex(), commentID.Hex(), model.VoteTargetComment, true).
					Return(nil)
			},
			wantErr:   nil,
			wantScore: 5, // 6 - 1
		},
		{
			name:      "error comment not found",
			userID:    userID.Hex(),
			commentID: primitive.NewObjectID().Hex(),
			voteValue: true,
			mockSetup: func() {
				mockCommentRepo.EXPECT().
					GetByID(gomock.Any(), gomock.Any()).
					Return(nil, mongo.ErrNoDocuments)
			},
			wantErr:   mongo.ErrNoDocuments,
			wantScore: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			result, err := svc.VoteOnTarget(tt.userID, tt.commentID, model.VoteTargetComment, tt.voteValue)

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
				return
			}

			if result == nil {
				t.Error("expected result, got nil")
				return
			}

			if result.Score != tt.wantScore {
				t.Errorf("expected score %d, got %d", tt.wantScore, result.Score)
			}
		})
	}
}

// TestGetUserVote tests getting a user's vote on a target
func TestGetUserVote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVoteRepo := mocks.NewMockVoteRepo(ctrl)
	svc := NewVoteService(mockVoteRepo, nil, nil, nil)

	userID := primitive.NewObjectID()
	postID := primitive.NewObjectID()

	existingVote := &model.Vote{
		UserID:     userID,
		TargetID:   postID,
		TargetType: model.VoteTargetPost,
		Value:      true,
	}

	tests := []struct {
		name      string
		userID    string
		targetID  string
		mockSetup func()
		wantErr   error
		wantVote  bool
	}{
		{
			name:     "success get existing vote",
			userID:   userID.Hex(),
			targetID: postID.Hex(),
			mockSetup: func() {
				mockVoteRepo.EXPECT().
					GetUserVote(gomock.Any(), userID.Hex(), postID.Hex(), model.VoteTargetPost).
					Return(existingVote, nil)
			},
			wantErr:  nil,
			wantVote: true,
		},
		{
			name:     "vote not found",
			userID:   userID.Hex(),
			targetID: primitive.NewObjectID().Hex(),
			mockSetup: func() {
				mockVoteRepo.EXPECT().
					GetUserVote(gomock.Any(), gomock.Any(), gomock.Any(), model.VoteTargetPost).
					Return(nil, mongo.ErrNoDocuments)
			},
			wantErr:  mongo.ErrNoDocuments,
			wantVote: false,
		},
		{
			name:     "error repository failure",
			userID:   userID.Hex(),
			targetID: postID.Hex(),
			mockSetup: func() {
				mockVoteRepo.EXPECT().
					GetUserVote(gomock.Any(), userID.Hex(), postID.Hex(), model.VoteTargetPost).
					Return(nil, errors.New("database error"))
			},
			wantErr:  errors.New("database error"),
			wantVote: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			result, err := svc.GetUserVote(tt.userID, tt.targetID, model.VoteTargetPost)

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
				return
			}

			if tt.wantVote && result == nil {
				t.Error("expected vote, got nil")
			}
		})
	}
}

// TestFindUserVotes tests finding user votes for multiple targets
func TestFindUserVotes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVoteRepo := mocks.NewMockVoteRepo(ctrl)
	svc := NewVoteService(mockVoteRepo, nil, nil, nil)

	userID := primitive.NewObjectID()
	postID1 := primitive.NewObjectID().Hex()
	postID2 := primitive.NewObjectID().Hex()
	postID3 := primitive.NewObjectID().Hex()
	targetIDs := []string{postID1, postID2, postID3}

	votesMap := map[string]string{
		postID1: "up",
		postID2: "down",
		// postID3 has no vote
	}

	tests := []struct {
		name      string
		userID    string
		targetIDs []string
		mockSetup func()
		wantErr   error
		wantCount int
	}{
		{
			name:      "success find user votes",
			userID:    userID.Hex(),
			targetIDs: targetIDs,
			mockSetup: func() {
				mockVoteRepo.EXPECT().
					FindUserVotes(gomock.Any(), userID.Hex(), targetIDs, model.VoteTargetPost).
					Return(votesMap, nil)
			},
			wantErr:   nil,
			wantCount: 2, // 2 votes found
		},
		{
			name:      "success no votes found",
			userID:    userID.Hex(),
			targetIDs: targetIDs,
			mockSetup: func() {
				mockVoteRepo.EXPECT().
					FindUserVotes(gomock.Any(), userID.Hex(), targetIDs, model.VoteTargetPost).
					Return(map[string]string{}, nil)
			},
			wantErr:   nil,
			wantCount: 0,
		},
		{
			name:      "error repository failure",
			userID:    userID.Hex(),
			targetIDs: targetIDs,
			mockSetup: func() {
				mockVoteRepo.EXPECT().
					FindUserVotes(gomock.Any(), userID.Hex(), targetIDs, model.VoteTargetPost).
					Return(nil, errors.New("database error"))
			},
			wantErr:   errors.New("database error"),
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			result, err := svc.FindUserVotes(tt.userID, tt.targetIDs, model.VoteTargetPost)

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
				return
			}

			if len(result) != tt.wantCount {
				t.Errorf("expected %d votes, got %d", tt.wantCount, len(result))
			}
		})
	}
}
