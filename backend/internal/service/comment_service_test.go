package service

import (
	"context"
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

func TestCreateComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommentRepo := mocks.NewMockCommentRepo(ctrl)
	mockUserRepo := mocks.NewMockUserRepo(ctrl)
	mockCommunityRepo := mocks.NewMockCommunityRepo(ctrl)
	mockPostRepo := mocks.NewMockPostRepo(ctrl)
	mockEventBus := bus.NewMockEventBus(ctrl)

	svc := NewCommentService(mockCommentRepo, mockUserRepo, mockCommunityRepo, mockPostRepo, mockEventBus)

	// prepare ids
	userID := primitive.NewObjectID()
	postID := primitive.NewObjectID()
	communityID := primitive.NewObjectID()
	parentID := primitive.NewObjectID()

	now := time.Now()

	tests := []struct {
		name           string
		request        *dto.CreateCommentRequest
		userID         string
		setupMocks     func()
		wantErr        error
		validateResult func(t *testing.T, c *model.Comment)
	}{
		{
			name: "success create root comment",
			request: &dto.CreateCommentRequest{
				PostID:  postID.Hex(),
				Content: "hello",
			},
			userID: userID.Hex(),
			setupMocks: func() {
				mockPostRepo.EXPECT().GetByID(gomock.Any(), postID.Hex()).Return(&model.Post{ID: postID, CommunityID: communityID}, nil)
				mockCommunityRepo.EXPECT().IsUserBanned(gomock.Any(), userID.Hex(), model.Muted, communityID.Hex()).Return(false, nil)
				mockUserRepo.EXPECT().GetByID(gomock.Any(), userID.Hex()).Return(&model.User{ID: userID, Username: "u1", RoleContent: model.RoleContent{AsUser: &model.UserRoleContent{Avatar: &model.Image{URL: "a"}}}}, nil)
				mockCommentRepo.EXPECT().Create(gomock.Any(), gomock.AssignableToTypeOf(&model.Comment{})).DoAndReturn(func(ctx context.Context, c *model.Comment) (*model.Comment, error) {
					c.ID = primitive.NewObjectID()
					c.CreatedAt = now
					return c, nil
				})
				mockPostRepo.EXPECT().Increment(gomock.Any(), postID.Hex(), "comments_count", 1).Return(nil)
				mockEventBus.EXPECT().Publish(gomock.Any())
			},
			wantErr: nil,
			validateResult: func(t *testing.T, c *model.Comment) {
				if c == nil {
					t.Fatalf("expected comment, got nil")
				}
				if c.Content != "hello" {
					t.Errorf("content mismatch: want %v got %v", "hello", c.Content)
				}
			},
		},
		{
			name: "success create child comment with parent author fetched",
			request: &dto.CreateCommentRequest{
				PostID:   postID.Hex(),
				ParentID: func() *string { s := parentID.Hex(); return &s }(),
				Content:  "reply",
			},
			userID: userID.Hex(),
			setupMocks: func() {
				mockPostRepo.EXPECT().GetByID(gomock.Any(), postID.Hex()).Return(&model.Post{ID: postID, CommunityID: communityID}, nil)
				mockCommunityRepo.EXPECT().IsUserBanned(gomock.Any(), userID.Hex(), model.Muted, communityID.Hex()).Return(false, nil)
				mockUserRepo.EXPECT().GetByID(gomock.Any(), userID.Hex()).Return(&model.User{ID: userID, Username: "u1", RoleContent: model.RoleContent{AsUser: &model.UserRoleContent{Avatar: &model.Image{URL: "a"}}}}, nil)
				// When fetching parent, return a comment with an author
				mockCommentRepo.EXPECT().GetByID(gomock.Any(), parentID.Hex()).Return(&model.Comment{ID: parentID, Author: model.CommentAuthor{ID: primitive.NewObjectID(), Username: "parentAuthor"}}, nil)
				mockCommentRepo.EXPECT().Create(gomock.Any(), gomock.AssignableToTypeOf(&model.Comment{})).DoAndReturn(func(ctx context.Context, c *model.Comment) (*model.Comment, error) {
					c.ID = primitive.NewObjectID()
					c.CreatedAt = now
					return c, nil
				})
				mockPostRepo.EXPECT().Increment(gomock.Any(), postID.Hex(), "comments_count", 1).Return(nil)
				mockEventBus.EXPECT().Publish(gomock.Any())
			},
			wantErr: nil,
			validateResult: func(t *testing.T, c *model.Comment) {
				if c.ParentID == nil {
					t.Fatalf("expected parent id set")
				}
			},
		},
		{
			name: "error user muted",
			request: &dto.CreateCommentRequest{
				PostID:  postID.Hex(),
				Content: "x",
			},
			userID: userID.Hex(),
			setupMocks: func() {
				mockPostRepo.EXPECT().GetByID(gomock.Any(), postID.Hex()).Return(&model.Post{ID: postID, CommunityID: communityID}, nil)
				mockCommunityRepo.EXPECT().IsUserBanned(gomock.Any(), userID.Hex(), model.Muted, communityID.Hex()).Return(true, nil)
			},
			wantErr: apperror.ErrUserIsMuted,
		},
		{
			name: "error post not found",
			request: &dto.CreateCommentRequest{
				PostID:  postID.Hex(),
				Content: "x",
			},
			userID: userID.Hex(),
			setupMocks: func() {
				mockPostRepo.EXPECT().GetByID(gomock.Any(), postID.Hex()).Return(nil, mongo.ErrNoDocuments)
			},
			wantErr: mongo.ErrNoDocuments,
		},
		{
			name: "error create repo fails",
			request: &dto.CreateCommentRequest{
				PostID:  postID.Hex(),
				Content: "x",
			},
			userID: userID.Hex(),
			setupMocks: func() {
				mockPostRepo.EXPECT().GetByID(gomock.Any(), postID.Hex()).Return(&model.Post{ID: postID, CommunityID: communityID}, nil)
				mockCommunityRepo.EXPECT().IsUserBanned(gomock.Any(), userID.Hex(), model.Muted, communityID.Hex()).Return(false, nil)
				mockUserRepo.EXPECT().GetByID(gomock.Any(), userID.Hex()).Return(&model.User{ID: userID, Username: "u1", RoleContent: model.RoleContent{AsUser: &model.UserRoleContent{Avatar: &model.Image{URL: "a"}}}}, nil)
				mockCommentRepo.EXPECT().Create(gomock.Any(), gomock.AssignableToTypeOf(&model.Comment{})).Return(nil, errors.New("db error"))
			},
			wantErr: errors.New("db error"),
		},
		{
			name: "error invalid post id",
			request: &dto.CreateCommentRequest{
				PostID:  "invalid-id",
				Content: "x",
			},
			userID:     userID.Hex(),
			setupMocks: func() {},
			wantErr:    errors.New("invalid id"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			tt.setupMocks()

			c, err := svc.CreateComment(tt.request, tt.userID)

			if tt.wantErr != nil {
				// special-case invalid post id test: ensure error is non-nil
				if tt.name == "error invalid post id" {
					if err == nil {
						t.Fatalf("expected error, got nil")
					}
					return
				}

				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantErr)
				}
				if err.Error() != tt.wantErr.Error() {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			// expected success
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.validateResult != nil {
				tt.validateResult(t, c)
			}
		})
	}
}

func TestDeleteCommentByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommentRepo := mocks.NewMockCommentRepo(ctrl)
	mockUserRepo := mocks.NewMockUserRepo(ctrl)
	mockCommunityRepo := mocks.NewMockCommunityRepo(ctrl)
	mockPostRepo := mocks.NewMockPostRepo(ctrl)
	mockEventBus := bus.NewMockEventBus(ctrl)

	svc := NewCommentService(mockCommentRepo, mockUserRepo, mockCommunityRepo, mockPostRepo, mockEventBus)

	// ids
	userID := primitive.NewObjectID()
	otherUserID := primitive.NewObjectID()
	postID := primitive.NewObjectID()
	commentID := primitive.NewObjectID()

	tests := []struct {
		name       string
		comment    *model.Comment
		userID     string
		setupMocks func()
		wantErr    error
	}{
		{
			name:    "success delete",
			comment: &model.Comment{ID: commentID, Author: model.CommentAuthor{ID: userID}, PostID: postID, IsDeleted: false},
			userID:  userID.Hex(),
			setupMocks: func() {
				mockCommentRepo.EXPECT().GetByID(gomock.Any(), commentID.Hex()).Return(&model.Comment{ID: commentID, Author: model.CommentAuthor{ID: userID}, PostID: postID, IsDeleted: false}, nil)
				mockCommentRepo.EXPECT().Delete(gomock.Any(), commentID.Hex()).Return(nil)
				mockPostRepo.EXPECT().Increment(gomock.Any(), postID.Hex(), "comments_count", -1).Return(nil)
			},
			wantErr: nil,
		},
		{
			name:    "error comment not found",
			comment: nil,
			userID:  userID.Hex(),
			setupMocks: func() {
				mockCommentRepo.EXPECT().GetByID(gomock.Any(), commentID.Hex()).Return(nil, mongo.ErrNoDocuments)
			},
			wantErr: mongo.ErrNoDocuments,
		},
		{
			name:    "error forbidden not author",
			comment: &model.Comment{ID: commentID, Author: model.CommentAuthor{ID: otherUserID}, PostID: postID, IsDeleted: false},
			userID:  userID.Hex(),
			setupMocks: func() {
				mockCommentRepo.EXPECT().GetByID(gomock.Any(), commentID.Hex()).Return(&model.Comment{ID: commentID, Author: model.CommentAuthor{ID: otherUserID}, PostID: postID, IsDeleted: false}, nil)
			},
			wantErr: apperror.ErrForbidden,
		},
		{
			name:    "error repo delete fail",
			comment: &model.Comment{ID: commentID, Author: model.CommentAuthor{ID: userID}, PostID: postID, IsDeleted: false},
			userID:  userID.Hex(),
			setupMocks: func() {
				mockCommentRepo.EXPECT().GetByID(gomock.Any(), commentID.Hex()).Return(&model.Comment{ID: commentID, Author: model.CommentAuthor{ID: userID}, PostID: postID, IsDeleted: false}, nil)
				mockCommentRepo.EXPECT().Delete(gomock.Any(), commentID.Hex()).Return(errors.New("delete failed"))
			},
			wantErr: errors.New("delete failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			err := svc.DeleteCommentByID(commentID.Hex(), tt.userID)
			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected err %v, got nil", tt.wantErr)
				}
				if err.Error() != tt.wantErr.Error() {
					t.Fatalf("expected err %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
