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

func TestCreateCommunity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommunityRepo := mocks.NewMockCommunityRepo(ctrl)
	mockMembershipRepo := mocks.NewMockMembershipRepo(ctrl)
	mockPostRepo := mocks.NewMockPostRepo(ctrl)
	mockUserRepo := mocks.NewMockUserRepo(ctrl)
	mockEventBus := bus.NewMockEventBus(ctrl)

	svc := NewCommunityService(mockCommunityRepo, mockMembershipRepo, mockPostRepo, mockUserRepo, mockEventBus)

	userID := primitive.NewObjectID()
	communityID := primitive.NewObjectID()

	tests := []struct {
		name        string
		requesterID string
		req         *dto.CreateCommunityRequest
		wantErr     error
		setupMocks  func()
		validate    func(t *testing.T, comm *model.Community)
	}{
		{
			name:        "successfully create a community with valid inputs",
			requesterID: userID.Hex(),
			req: &dto.CreateCommunityRequest{
				Name:          "Golang Enthusiasts",
				Description:   ptrStr("All things Go."),
				Avatar:        ptrStr("https://example.com/avatar.png"),
				Banner:        ptrStr("https://example.com/banner.png"),
				Setting:       model.CommunitySetting{},
				Rules:         []model.CommunityRule{},
				Moderators:    []model.Moderator{},
				CreatorName:   "User1",
				CreatorAvatar: "https://example.com/user1-avatar.png",
			},
			wantErr: nil,
			setupMocks: func() {
				mockCommunityRepo.EXPECT().
					Create(gomock.Any(), gomock.AssignableToTypeOf(&model.Community{})).
					DoAndReturn(func(ctx context.Context, comm *model.Community) (*model.Community, error) {
						comm.ID = communityID
						return comm, nil
					})
				mockMembershipRepo.EXPECT().
					Create(gomock.Any(), gomock.AssignableToTypeOf(&model.Membership{})).
					Return(&model.Membership{}, nil)
			},
			validate: func(t *testing.T, comm *model.Community) {
				if comm == nil {
					t.Fatal("expected community, got nil")
				}
				if comm.Name != "Golang Enthusiasts" {
					t.Errorf("expected name 'Golang Enthusiasts', got '%s'", comm.Name)
				}
				if comm.Description == nil || *comm.Description != "All things Go." {
					t.Errorf("expected description 'All things Go.', got '%v'", comm.Description)
				}
				if comm.CreateByID != userID {
					t.Errorf("expected creator ID %s, got %s", userID, comm.CreateByID)
				}
				if comm.IsDeleted {
					t.Error("expected IsDeleted to be false")
				}
			},
		},
		{
			name:        "attempt to create with empty name",
			requesterID: userID.Hex(),
			req: &dto.CreateCommunityRequest{
				Name:        "",
				Description: ptrStr("No name comm"),
			},
			wantErr:    apperror.ErrBadRequest,
			setupMocks: func() {},
		},
		{
			name:        "duplicate community name error",
			requesterID: userID.Hex(),
			req: &dto.CreateCommunityRequest{
				Name:        "Existing Community",
				Description: ptrStr("Already exists"),
			},
			wantErr: apperror.ErrCommunityNameExists,
			setupMocks: func() {
				mockCommunityRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil, mockDuplicateKeyMongoError())
			},
		},
		{
			name:        "repo community create error",
			requesterID: userID.Hex(),
			req: &dto.CreateCommunityRequest{
				Name:        "Test Community",
				Description: ptrStr("Test"),
			},
			wantErr: errors.New("db error"),
			setupMocks: func() {
				mockCommunityRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("db error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			comm, err := svc.CreateCommunity(tt.req, tt.requesterID)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantErr)
				}
				if !errors.Is(err, tt.wantErr) && err.Error() != tt.wantErr.Error() {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.validate != nil {
				tt.validate(t, comm)
			}
		})
	}
}

func TestGetCommunityByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommunityRepo := mocks.NewMockCommunityRepo(ctrl)
	mockMembershipRepo := mocks.NewMockMembershipRepo(ctrl)
	mockPostRepo := mocks.NewMockPostRepo(ctrl)
	mockUserRepo := mocks.NewMockUserRepo(ctrl)
	mockEventBus := bus.NewMockEventBus(ctrl)

	svc := NewCommunityService(mockCommunityRepo, mockMembershipRepo, mockPostRepo, mockUserRepo, mockEventBus)

	communityID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	existingCommunity := &model.Community{
		ID:          communityID,
		Name:        "Test Community",
		Description: ptrStr("Test"),
		CreateByID:  userID,
		CreateAt:    time.Now(),
		IsDeleted:   false,
		IsBanned:    false,
	}

	tests := []struct {
		name        string
		communityID string
		requesterID *string
		wantErr     error
		setupMocks  func()
		validate    func(t *testing.T, comm *model.Community)
	}{
		{
			name:        "retrieve community by valid ID",
			communityID: communityID.Hex(),
			requesterID: nil,
			wantErr:     nil,
			setupMocks: func() {
				mockCommunityRepo.EXPECT().
					GetByID(gomock.Any(), communityID.Hex()).
					Return(existingCommunity, nil)
			},
			validate: func(t *testing.T, comm *model.Community) {
				if comm.ID != communityID {
					t.Errorf("expected community ID %s, got %s", communityID, comm.ID)
				}
				if comm.Name != "Test Community" {
					t.Errorf("expected name 'Test Community', got '%s'", comm.Name)
				}
			},
		},
		{
			name:        "attempt to retrieve non-existent community",
			communityID: primitive.NewObjectID().Hex(),
			requesterID: nil,
			wantErr:     apperror.ErrCommunityNotFound,
			setupMocks: func() {
				mockCommunityRepo.EXPECT().
					GetByID(gomock.Any(), gomock.Any()).
					Return(nil, mongo.ErrNoDocuments)
			},
		},
		{
			name:        "requester is banned from community",
			communityID: communityID.Hex(),
			requesterID: ptrStr(userID.Hex()),
			wantErr:     apperror.ErrUserIsBannedFromCommunity,
			setupMocks: func() {
				mockCommunityRepo.EXPECT().
					IsUserBanned(gomock.Any(), userID.Hex(), model.Banned, communityID.Hex()).
					Return(true, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			comm, err := svc.GetCommunityByID(tt.communityID, tt.requesterID)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantErr)
				}
				if !errors.Is(err, tt.wantErr) && err.Error() != tt.wantErr.Error() {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.validate != nil {
				tt.validate(t, comm)
			}
		})
	}
}

func TestUpdateCommunity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommunityRepo := mocks.NewMockCommunityRepo(ctrl)
	mockMembershipRepo := mocks.NewMockMembershipRepo(ctrl)
	mockPostRepo := mocks.NewMockPostRepo(ctrl)
	mockUserRepo := mocks.NewMockUserRepo(ctrl)
	mockEventBus := bus.NewMockEventBus(ctrl)

	svc := NewCommunityService(mockCommunityRepo, mockMembershipRepo, mockPostRepo, mockUserRepo, mockEventBus)

	communityID := primitive.NewObjectID()
	moderatorID := primitive.NewObjectID()
	nonModeratorID := primitive.NewObjectID()

	communityWithMod := &model.Community{
		ID:          communityID,
		Name:        "Test Community",
		Description: ptrStr("Old Description"),
		CreateByID:  moderatorID,
		Moderators:  []model.Moderator{},
		IsDeleted:   false,
	}

	communityWithoutMod := &model.Community{
		ID:          communityID,
		Name:        "Test Community",
		Description: ptrStr("Old Description"),
		CreateByID:  moderatorID,
		Moderators:  []model.Moderator{},
		IsDeleted:   false,
	}

	tests := []struct {
		name        string
		requesterID string
		req         *dto.UpdateCommunityRequest
		wantErr     error
		setupMocks  func()
		validate    func(t *testing.T, comm *model.Community)
	}{
		{
			name:        "successfully update community description",
			requesterID: moderatorID.Hex(),
			req: &dto.UpdateCommunityRequest{
				CommunityID: communityID.Hex(),
				Description: ptrStr("New Description"),
			},
			wantErr: nil,
			setupMocks: func() {
				mockCommunityRepo.EXPECT().
					GetByID(gomock.Any(), communityID.Hex()).
					Return(communityWithMod, nil)
				mockCommunityRepo.EXPECT().
					IsModerator(gomock.Any(), communityID.Hex(), moderatorID.Hex()).
					Return(true, nil)
				mockCommunityRepo.EXPECT().
					Replace(gomock.Any(), gomock.AssignableToTypeOf(&model.Community{})).
					Return(nil)
			},
			validate: func(t *testing.T, comm *model.Community) {
				if comm == nil {
					t.Fatal("expected community, got nil")
				}
				if comm.Description == nil || *comm.Description != "New Description" {
					t.Errorf("expected description 'New Description', got '%v'", comm.Description)
				}
			},
		},
		{
			name:        "attempt to update without moderator permissions",
			requesterID: nonModeratorID.Hex(),
			req: &dto.UpdateCommunityRequest{
				CommunityID: communityID.Hex(),
				Description: ptrStr("Hacked"),
			},
			wantErr: apperror.ErrForbidden,
			setupMocks: func() {
				mockCommunityRepo.EXPECT().
					GetByID(gomock.Any(), communityID.Hex()).
					Return(communityWithoutMod, nil)
				mockCommunityRepo.EXPECT().
					IsModerator(gomock.Any(), communityID.Hex(), nonModeratorID.Hex()).
					Return(false, nil)
			},
		},
		{
			name:        "community not found",
			requesterID: moderatorID.Hex(),
			req: &dto.UpdateCommunityRequest{
				CommunityID: communityID.Hex(),
				Description: ptrStr("New"),
			},
			wantErr: mongo.ErrNoDocuments,
			setupMocks: func() {
				mockCommunityRepo.EXPECT().
					GetByID(gomock.Any(), communityID.Hex()).
					Return(nil, mongo.ErrNoDocuments)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			comm, err := svc.UpdateCommunity(tt.req, tt.requesterID)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.validate != nil {
				tt.validate(t, comm)
			}
		})
	}
}

func TestAddModerator(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommunityRepo := mocks.NewMockCommunityRepo(ctrl)
	mockMembershipRepo := mocks.NewMockMembershipRepo(ctrl)
	mockPostRepo := mocks.NewMockPostRepo(ctrl)
	mockUserRepo := mocks.NewMockUserRepo(ctrl)
	mockEventBus := bus.NewMockEventBus(ctrl)

	svc := NewCommunityService(mockCommunityRepo, mockMembershipRepo, mockPostRepo, mockUserRepo, mockEventBus)

	communityID := primitive.NewObjectID()
	creatorID := primitive.NewObjectID()
	newModID := primitive.NewObjectID()
	existingModID := primitive.NewObjectID()

	communityWithMod := &model.Community{
		ID:         communityID,
		Name:       "Test",
		CreateByID: creatorID,
		Moderators: []model.Moderator{
			{
				UserID:     existingModID,
				Username:   "ExistingMod",
				AssignedAt: time.Now(),
			},
		},
	}

	tests := []struct {
		name        string
		requesterID string
		req         *dto.AddModeratorRequest
		wantErr     error
		setupMocks  func()
	}{
		{
			name:        "successfully add new moderator",
			requesterID: creatorID.Hex(),
			req: &dto.AddModeratorRequest{
				CommunityID:    communityID.Hex(),
				AddedModerator: []string{newModID.Hex()},
			},
			wantErr: nil,
			setupMocks: func() {
				mockCommunityRepo.EXPECT().
					GetByID(gomock.Any(), communityID.Hex()).
					Return(communityWithMod, nil)
				mockCommunityRepo.EXPECT().
					IsCreator(gomock.Any(), communityID.Hex(), creatorID.Hex()).
					Return(true, nil)
				mockUserRepo.EXPECT().
					GetByID(gomock.Any(), newModID.Hex()).
					Return(&model.User{ID: primitive.NewObjectID(), Username: "NewMod", RoleContent: model.RoleContent{AsUser: &model.UserRoleContent{Avatar: &model.Image{URL: "https://example.com/avatar.png"}}}}, nil)
				mockCommunityRepo.EXPECT().
					Replace(gomock.Any(), gomock.Any()).
					Return(nil)
				mockEventBus.EXPECT().
					Publish(gomock.Any())
			},
		},
		{
			name:        "attempt to add user already a moderator",
			requesterID: creatorID.Hex(),
			req: &dto.AddModeratorRequest{
				CommunityID:    communityID.Hex(),
				AddedModerator: []string{existingModID.Hex()},
			},
			wantErr: apperror.ErrModeratorAlreadyExists,
			setupMocks: func() {
				mockCommunityRepo.EXPECT().
					GetByID(gomock.Any(), communityID.Hex()).
					Return(communityWithMod, nil)
				mockCommunityRepo.EXPECT().
					IsCreator(gomock.Any(), communityID.Hex(), creatorID.Hex()).
					Return(true, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := svc.AddModerator(tt.req, tt.requesterID)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestRemoveModerator(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommunityRepo := mocks.NewMockCommunityRepo(ctrl)
	mockMembershipRepo := mocks.NewMockMembershipRepo(ctrl)
	mockPostRepo := mocks.NewMockPostRepo(ctrl)
	mockUserRepo := mocks.NewMockUserRepo(ctrl)
	mockEventBus := bus.NewMockEventBus(ctrl)

	svc := NewCommunityService(mockCommunityRepo, mockMembershipRepo, mockPostRepo, mockUserRepo, mockEventBus)

	communityID := primitive.NewObjectID()
	creatorID := primitive.NewObjectID()
	modID := primitive.NewObjectID()

	communityWithMod := &model.Community{
		ID:         communityID,
		CreateByID: creatorID,
		Moderators: []model.Moderator{
			{
				UserID:     modID,
				Username:   "ModUser",
				AssignedAt: time.Now(),
			},
		},
	}

	tests := []struct {
		name        string
		requesterID string
		req         *dto.RemoveModeratorRequest
		wantErr     error
		setupMocks  func()
	}{
		{
			name:        "successfully remove moderator",
			requesterID: creatorID.Hex(),
			req: &dto.RemoveModeratorRequest{
				CommunityID:      communityID.Hex(),
				RemovedModerator: []string{modID.Hex()},
			},
			wantErr: nil,
			setupMocks: func() {
				mockCommunityRepo.EXPECT().
					GetByID(gomock.Any(), communityID.Hex()).
					Return(communityWithMod, nil)
				mockCommunityRepo.EXPECT().
					IsCreator(gomock.Any(), communityID.Hex(), creatorID.Hex()).
					Return(true, nil)
				mockCommunityRepo.EXPECT().
					Replace(gomock.Any(), gomock.Any()).
					Return(nil)
			},
		},
		{
			name:        "attempt to remove self",
			requesterID: modID.Hex(),
			req: &dto.RemoveModeratorRequest{
				CommunityID:      communityID.Hex(),
				RemovedModerator: []string{modID.Hex()},
			},
			wantErr: apperror.ErrCannotRemoveModerator,
			setupMocks: func() {
				mockCommunityRepo.EXPECT().
					GetByID(gomock.Any(), communityID.Hex()).
					Return(communityWithMod, nil)
				mockCommunityRepo.EXPECT().
					IsCreator(gomock.Any(), communityID.Hex(), modID.Hex()).
					Return(true, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := svc.RemoveModerator(tt.req, tt.requesterID)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestDeleteCommunityByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommunityRepo := mocks.NewMockCommunityRepo(ctrl)
	mockMembershipRepo := mocks.NewMockMembershipRepo(ctrl)
	mockPostRepo := mocks.NewMockPostRepo(ctrl)
	mockUserRepo := mocks.NewMockUserRepo(ctrl)
	mockEventBus := bus.NewMockEventBus(ctrl)

	svc := NewCommunityService(mockCommunityRepo, mockMembershipRepo, mockPostRepo, mockUserRepo, mockEventBus)

	communityID := primitive.NewObjectID()
	creatorID := primitive.NewObjectID()
	nonCreatorID := primitive.NewObjectID()

	tests := []struct {
		name        string
		communityID string
		requesterID string
		wantErr     error
		setupMocks  func()
	}{
		{
			name:        "successfully delete community",
			communityID: communityID.Hex(),
			requesterID: creatorID.Hex(),
			wantErr:     nil,
			setupMocks: func() {
				mockCommunityRepo.EXPECT().
					IsCreator(gomock.Any(), communityID.Hex(), creatorID.Hex()).
					Return(true, nil)
				mockCommunityRepo.EXPECT().
					Delete(gomock.Any(), communityID.Hex()).
					Return(nil)
			},
		},
		{
			name:        "attempt to delete without ownership",
			communityID: communityID.Hex(),
			requesterID: nonCreatorID.Hex(),
			wantErr:     apperror.ErrForbidden,
			setupMocks: func() {
				mockCommunityRepo.EXPECT().
					IsCreator(gomock.Any(), communityID.Hex(), nonCreatorID.Hex()).
					Return(false, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := svc.DeleteCommunityByID(tt.communityID, tt.requesterID)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestBanUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommunityRepo := mocks.NewMockCommunityRepo(ctrl)
	mockMembershipRepo := mocks.NewMockMembershipRepo(ctrl)
	mockPostRepo := mocks.NewMockPostRepo(ctrl)
	mockUserRepo := mocks.NewMockUserRepo(ctrl)
	mockEventBus := bus.NewMockEventBus(ctrl)

	svc := NewCommunityService(mockCommunityRepo, mockMembershipRepo, mockPostRepo, mockUserRepo, mockEventBus)

	communityID := primitive.NewObjectID()
	moderatorID := primitive.NewObjectID()
	nonModeratorID := primitive.NewObjectID()
	userToBanID := primitive.NewObjectID()

	tests := []struct {
		name        string
		requesterID string
		req         *dto.CommunityBanUserRequest
		wantErr     error
		setupMocks  func()
	}{
		{
			name:        "successfully ban user",
			requesterID: moderatorID.Hex(),
			req: &dto.CommunityBanUserRequest{
				CommunityID: communityID.Hex(),
				UserID:      userToBanID.Hex(),
				Type:        "banned",
				Reason:      "Spam",
				LengthDays:  30,
			},
			wantErr: nil,
			setupMocks: func() {
				mockCommunityRepo.EXPECT().
					IsModerator(gomock.Any(), communityID.Hex(), moderatorID.Hex()).
					Return(true, nil)
				mockCommunityRepo.EXPECT().
					BanUser(gomock.Any(), gomock.AssignableToTypeOf(&model.CommunityBan{})).
					Return(nil)
			},
		},
		{
			name:        "attempt to ban without moderator permissions",
			requesterID: nonModeratorID.Hex(),
			req: &dto.CommunityBanUserRequest{
				CommunityID: communityID.Hex(),
				UserID:      userToBanID.Hex(),
				Type:        "banned",
				Reason:      "Spam",
				LengthDays:  30,
			},
			wantErr: apperror.ErrForbidden,
			setupMocks: func() {
				mockCommunityRepo.EXPECT().
					IsModerator(gomock.Any(), communityID.Hex(), nonModeratorID.Hex()).
					Return(false, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := svc.BanUser(tt.req, tt.requesterID)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestGetBannedUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommunityRepo := mocks.NewMockCommunityRepo(ctrl)
	mockMembershipRepo := mocks.NewMockMembershipRepo(ctrl)
	mockPostRepo := mocks.NewMockPostRepo(ctrl)
	mockUserRepo := mocks.NewMockUserRepo(ctrl)
	mockEventBus := bus.NewMockEventBus(ctrl)

	svc := NewCommunityService(mockCommunityRepo, mockMembershipRepo, mockPostRepo, mockUserRepo, mockEventBus)

	communityID := primitive.NewObjectID()
	moderatorID := primitive.NewObjectID()
	bannedUserID := primitive.NewObjectID()

	bannedUser := &model.User{
		ID:       bannedUserID,
		Username: "BannedUser",
	}

	tests := []struct {
		name        string
		communityID string
		banTypeStr  string
		expired     bool
		requesterID string
		wantErr     error
		setupMocks  func()
		validate    func(t *testing.T, users []*model.User)
	}{
		{
			name:        "retrieve banned users list",
			communityID: communityID.Hex(),
			banTypeStr:  "banned",
			expired:     false,
			requesterID: moderatorID.Hex(),
			wantErr:     nil,
			setupMocks: func() {
				mockCommunityRepo.EXPECT().
					IsModerator(gomock.Any(), communityID.Hex(), moderatorID.Hex()).
					Return(true, nil)
				mockCommunityRepo.EXPECT().
					GetBannedUsers(gomock.Any(), communityID.Hex(), false).
					Return([]*model.User{bannedUser}, nil)
			},
			validate: func(t *testing.T, users []*model.User) {
				if len(users) != 1 {
					t.Errorf("expected 1 user, got %d", len(users))
				}
				if users[0].Username != "BannedUser" {
					t.Errorf("expected username 'BannedUser', got '%s'", users[0].Username)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			users, err := svc.GetBannedUsers(tt.communityID, tt.banTypeStr, tt.expired, tt.requesterID)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.validate != nil {
				tt.validate(t, users)
			}
		})
	}
}

func TestUnbanUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommunityRepo := mocks.NewMockCommunityRepo(ctrl)
	mockMembershipRepo := mocks.NewMockMembershipRepo(ctrl)
	mockPostRepo := mocks.NewMockPostRepo(ctrl)
	mockUserRepo := mocks.NewMockUserRepo(ctrl)
	mockEventBus := bus.NewMockEventBus(ctrl)

	svc := NewCommunityService(mockCommunityRepo, mockMembershipRepo, mockPostRepo, mockUserRepo, mockEventBus)

	communityID := primitive.NewObjectID()
	moderatorID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	tests := []struct {
		name        string
		userID      string
		communityID string
		requesterID string
		wantErr     error
		setupMocks  func()
	}{
		{
			name:        "successfully unban user",
			userID:      userID.Hex(),
			communityID: communityID.Hex(),
			requesterID: moderatorID.Hex(),
			wantErr:     nil,
			setupMocks: func() {
				mockCommunityRepo.EXPECT().
					IsModerator(gomock.Any(), communityID.Hex(), moderatorID.Hex()).
					Return(true, nil)
				mockCommunityRepo.EXPECT().
					UnbanUser(gomock.Any(), userID.Hex(), communityID.Hex()).
					Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := svc.UnbanUser(tt.userID, tt.communityID, tt.requesterID)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestGetPendingPosts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommunityRepo := mocks.NewMockCommunityRepo(ctrl)
	mockMembershipRepo := mocks.NewMockMembershipRepo(ctrl)
	mockPostRepo := mocks.NewMockPostRepo(ctrl)
	mockUserRepo := mocks.NewMockUserRepo(ctrl)
	mockEventBus := bus.NewMockEventBus(ctrl)

	svc := NewCommunityService(mockCommunityRepo, mockMembershipRepo, mockPostRepo, mockUserRepo, mockEventBus)

	communityID := primitive.NewObjectID()
	moderatorID := primitive.NewObjectID()
	authorID := primitive.NewObjectID()
	postID := primitive.NewObjectID()

	community := &model.Community{
		ID:         communityID,
		CreateByID: moderatorID,
	}

	pendingPost := &model.Post{
		ID:               postID,
		CommunityID:      communityID,
		AuthorID:         authorID,
		ModerationStatus: model.ModerationPending,
		CreatedAt:        time.Now(),
	}

	author := &model.User{
		ID:       authorID,
		Username: "Author",
	}

	tests := []struct {
		name        string
		communityID string
		moderatorID string
		page        int
		pageSize    int
		wantErr     error
		setupMocks  func()
		validate    func(t *testing.T, resp *dto.PaginatedPostsResponse)
	}{
		{
			name:        "retrieve pending posts queue",
			communityID: communityID.Hex(),
			moderatorID: moderatorID.Hex(),
			page:        1,
			pageSize:    20,
			wantErr:     nil,
			setupMocks: func() {
				mockCommunityRepo.EXPECT().
					GetByID(gomock.Any(), communityID.Hex()).
					Return(community, nil)
				mockCommunityRepo.EXPECT().
					IsModerator(gomock.Any(), communityID.Hex(), moderatorID.Hex()).
					Return(true, nil)
				mockPostRepo.EXPECT().
					Find(gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]*model.Post{pendingPost}, int64(1), nil)
				mockUserRepo.EXPECT().
					GetByIDs(gomock.Any(), gomock.Any()).
					Return([]*model.User{author}, nil)
			},
			validate: func(t *testing.T, resp *dto.PaginatedPostsResponse) {
				if resp == nil {
					t.Fatal("expected response, got nil")
				}
				if resp.Pagination.Total != 1 {
					t.Errorf("expected total 1, got %d", resp.Pagination.Total)
				}
				if len(resp.Posts) != 1 {
					t.Errorf("expected 1 post, got %d", len(resp.Posts))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			resp, err := svc.GetPendingPosts(tt.communityID, tt.moderatorID, tt.page, tt.pageSize)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.validate != nil {
				tt.validate(t, resp)
			}
		})
	}
}

func TestModeratePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommunityRepo := mocks.NewMockCommunityRepo(ctrl)
	mockMembershipRepo := mocks.NewMockMembershipRepo(ctrl)
	mockPostRepo := mocks.NewMockPostRepo(ctrl)
	mockUserRepo := mocks.NewMockUserRepo(ctrl)
	mockEventBus := bus.NewMockEventBus(ctrl)

	svc := NewCommunityService(mockCommunityRepo, mockMembershipRepo, mockPostRepo, mockUserRepo, mockEventBus)

	communityID := primitive.NewObjectID()
	postID := primitive.NewObjectID()
	authorID := primitive.NewObjectID()
	moderatorID := primitive.NewObjectID()
	nonModeratorID := primitive.NewObjectID()

	pendingPost := &model.Post{
		ID:               postID,
		CommunityID:      communityID,
		AuthorID:         authorID,
		ModerationStatus: model.ModerationPending,
	}

	tests := []struct {
		name        string
		communityID string
		postID      string
		moderatorID string
		approve     bool
		reason      *string
		wantErr     error
		setupMocks  func()
	}{
		{
			name:        "UTC-020: approve pending post",
			communityID: communityID.Hex(),
			postID:      postID.Hex(),
			moderatorID: moderatorID.Hex(),
			approve:     true,
			reason:      nil,
			wantErr:     nil,
			setupMocks: func() {
				mockCommunityRepo.EXPECT().
					IsModerator(gomock.Any(), communityID.Hex(), moderatorID.Hex()).
					Return(true, nil)
				mockPostRepo.EXPECT().
					GetByID(gomock.Any(), postID.Hex()).
					Return(pendingPost, nil)
				mockPostRepo.EXPECT().
					UpdateByID(gomock.Any(), postID.Hex(), gomock.Any()).
					Return(nil)
				mockEventBus.EXPECT().
					Publish(gomock.Any())
			},
		},
		{
			name:        "UTC-021: reject pending post with reason",
			communityID: communityID.Hex(),
			postID:      postID.Hex(),
			moderatorID: moderatorID.Hex(),
			approve:     false,
			reason:      ptrStr("Inappropriate content"),
			wantErr:     nil,
			setupMocks: func() {
				mockCommunityRepo.EXPECT().
					IsModerator(gomock.Any(), communityID.Hex(), moderatorID.Hex()).
					Return(true, nil)
				mockPostRepo.EXPECT().
					GetByID(gomock.Any(), postID.Hex()).
					Return(pendingPost, nil)
				mockPostRepo.EXPECT().
					UpdateByID(gomock.Any(), postID.Hex(), gomock.Any()).
					Return(nil)
				mockEventBus.EXPECT().
					Publish(gomock.Any())
			},
		},
		{
			name:        "UTC-022: attempt to moderate without permissions",
			communityID: communityID.Hex(),
			postID:      postID.Hex(),
			moderatorID: nonModeratorID.Hex(),
			approve:     true,
			reason:      nil,
			wantErr:     apperror.ErrForbidden,
			setupMocks: func() {
				mockCommunityRepo.EXPECT().
					IsModerator(gomock.Any(), communityID.Hex(), nonModeratorID.Hex()).
					Return(false, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := svc.ModeratePost(tt.communityID, tt.postID, tt.moderatorID, tt.approve, tt.reason)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestGetAllCommunitiesPaginated_FilterOutBanned(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommunityRepo := mocks.NewMockCommunityRepo(ctrl)
	mockMembershipRepo := mocks.NewMockMembershipRepo(ctrl)
	mockPostRepo := mocks.NewMockPostRepo(ctrl)
	mockUserRepo := mocks.NewMockUserRepo(ctrl)
	mockEventBus := bus.NewMockEventBus(ctrl)

	svc := NewCommunityService(mockCommunityRepo, mockMembershipRepo, mockPostRepo, mockUserRepo, mockEventBus)

	// Prepare communities
	comm1ID := primitive.NewObjectID()
	comm2ID := primitive.NewObjectID()

	comm1 := &model.Community{ID: comm1ID, Name: "Comm1"}
	comm2 := &model.Community{ID: comm2ID, Name: "Comm2"}

	// Case A: requesterID is nil -> no filtering
	mockCommunityRepo.EXPECT().
		GetAllPaginated(gomock.Any(), 1, 10).
		Return([]*model.Community{comm1, comm2}, int64(2), nil)

	resp, err := svc.GetAllCommunitiesPaginated(nil, 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Pagination.Total != 2 {
		t.Fatalf("expected total 2, got %d", resp.Pagination.Total)
	}

	// Case B: requesterID provided -> one community is banned for requester
	requesterID := primitive.NewObjectID().Hex()

	mockCommunityRepo.EXPECT().
		GetAllPaginated(gomock.Any(), 1, 10).
		Return([]*model.Community{comm1, comm2}, int64(2), nil)

	// Return comm2 as banned for this requester
	mockCommunityRepo.EXPECT().
		GetBannedCommunityIDs(gomock.Any(), requesterID, model.Banned, gomock.Any()).
		Return([]string{comm2ID.Hex()}, nil)

	resp2, err := svc.GetAllCommunitiesPaginated(ptrStr(requesterID), 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp2.Pagination.Total != 2 {
		t.Fatalf("expected total 2 (unchanged), got %d", resp2.Pagination.Total)
	}
	if len(resp2.Communities) != 1 {
		t.Fatalf("expected 1 community after filtering, got %d", len(resp2.Communities))
	}
	if resp2.Communities[0].Name != "Comm1" {
		t.Fatalf("expected remaining community to be Comm1, got %s", resp2.Communities[0].Name)
	}
}

// Helper function for pointer strings
func ptrStr(s string) *string {
	return &s
}

func mockDuplicateKeyMongoError() error {
	return mongo.WriteException{
		WriteErrors: []mongo.WriteError{
			{Code: 11000, Message: "duplicate key error"},
		},
	}
}
