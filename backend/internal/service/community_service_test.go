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
		name              string
		requesterID       string
		req               *dto.CreateCommunityRequest
		repoCommunityErr  error
		repoMembershipErr error
		wantErr           error
		validate          func(t *testing.T, comm *model.Community)
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
			repoCommunityErr:  nil,
			repoMembershipErr: nil,
			wantErr:           nil,
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
			wantErr: apperror.ErrBadRequest,
		},
		{
			name:        "duplicate community name error",
			requesterID: userID.Hex(),
			req: &dto.CreateCommunityRequest{
				Name:        "Existing Community",
				Description: ptrStr("Already exists"),
			},
			repoCommunityErr: mockDuplicateKeyMongoError(),
			wantErr:          apperror.ErrCommunityNameExists,
		},
		{
			name:        "repo community create error",
			requesterID: userID.Hex(),
			req: &dto.CreateCommunityRequest{
				Name:        "Test Community",
				Description: ptrStr("Test"),
			},
			repoCommunityErr: errors.New("db error"),
			wantErr:          errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.req != nil && tt.repoCommunityErr == nil && tt.repoMembershipErr == nil {
				mockCommunityRepo.EXPECT().
					Create(gomock.Any(), gomock.AssignableToTypeOf(&model.Community{})).
					DoAndReturn(func(ctx context.Context, comm *model.Community) (*model.Community, error) {
						comm.ID = communityID
						return comm, tt.repoCommunityErr
					})

				mockMembershipRepo.EXPECT().
					Create(gomock.Any(), gomock.AssignableToTypeOf(&model.Membership{})).
					Return(&model.Membership{}, tt.repoMembershipErr)
			} else if tt.repoCommunityErr != nil {
				mockCommunityRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil, tt.repoCommunityErr)
			}

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
		name            string
		communityID     string
		requesterID     *string
		repoGetErr      error
		repoGetCom      *model.Community
		repoIsBannedErr error
		repoIsBanned    bool
		wantErr         error
		validate        func(t *testing.T, comm *model.Community)
	}{
		{
			name:        "retrieve community by valid ID",
			communityID: communityID.Hex(),
			requesterID: nil,
			repoGetErr:  nil,
			repoGetCom:  existingCommunity,
			wantErr:     nil,
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
			repoGetErr:  mongo.ErrNoDocuments,
			wantErr:     apperror.ErrCommunityNotFound,
		},
		{
			name:            "requester is banned from community",
			communityID:     communityID.Hex(),
			requesterID:     ptrStr(userID.Hex()),
			repoGetErr:      nil,
			repoGetCom:      existingCommunity,
			repoIsBannedErr: nil,
			repoIsBanned:    true,
			wantErr:         apperror.ErrUserIsBannedFromCommunity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.requesterID != nil {
				mockCommunityRepo.EXPECT().
					IsUserBanned(gomock.Any(), *tt.requesterID, model.Banned, tt.communityID).
					Return(!tt.repoIsBanned, tt.repoIsBannedErr)
			}

			if (tt.repoGetErr == nil && !tt.repoIsBanned) || errors.Is(tt.repoGetErr, mongo.ErrNoDocuments) {
				mockCommunityRepo.EXPECT().
					GetByID(gomock.Any(), tt.communityID).
					Return(tt.repoGetCom, tt.repoGetErr)
			}

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
		repoGetErr  error
		repoGetCom  *model.Community
		repoUpdErr  error
		wantErr     error
		validate    func(t *testing.T, comm *model.Community)
	}{
		{
			name:        "successfully update community description",
			requesterID: moderatorID.Hex(),
			req: &dto.UpdateCommunityRequest{
				CommunityID: communityID.Hex(),
				Description: ptrStr("New Description"),
			},
			repoGetErr: nil,
			repoGetCom: communityWithMod,
			repoUpdErr: nil,
			wantErr:    nil,
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
			repoGetErr: nil,
			repoGetCom: communityWithoutMod,
			wantErr:    apperror.ErrForbidden,
		},
		{
			name:        "community not found",
			requesterID: moderatorID.Hex(),
			req: &dto.UpdateCommunityRequest{
				CommunityID: communityID.Hex(),
				Description: ptrStr("New"),
			},
			repoGetErr: mongo.ErrNoDocuments,
			wantErr:    mongo.ErrNoDocuments,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCommunityRepo.EXPECT().
				GetByID(gomock.Any(), tt.req.CommunityID).
				Return(tt.repoGetCom, tt.repoGetErr)

			if tt.repoGetErr == nil && tt.repoGetCom != nil && tt.wantErr == nil {
				mockCommunityRepo.EXPECT().
					Replace(gomock.Any(), gomock.AssignableToTypeOf(&model.Community{})).
					Return(tt.repoUpdErr)
			}

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
		name             string
		requesterID      string
		req              *dto.AddModeratorRequest
		repoGetErr       error
		repoGetCom       *model.Community
		repoUserExistErr error
		repoUserExists   bool
		repoMemberErr    error
		repoIsMember     bool
		repoUpdateErr    error
		wantErr          error
	}{
		{
			name:        "successfully add new moderator",
			requesterID: creatorID.Hex(),
			req: &dto.AddModeratorRequest{
				CommunityID: communityID.Hex(),
				AddedModerator: []dto.ModeratorDTO{
					{
						ModeratorID: newModID.Hex(),
						Username:    "NewMod",
					},
				},
			},
			repoGetErr:       nil,
			repoGetCom:       communityWithMod,
			repoUserExistErr: nil,
			repoUserExists:   true,
			repoMemberErr:    nil,
			repoIsMember:     true,
			repoUpdateErr:    nil,
			wantErr:          nil,
		},
		{
			name:        "attempt to add user already a moderator",
			requesterID: creatorID.Hex(),
			req: &dto.AddModeratorRequest{
				CommunityID: communityID.Hex(),
				AddedModerator: []dto.ModeratorDTO{
					{
						ModeratorID: existingModID.Hex(),
						Username:    "ExistingMod",
					},
				},
			},
			repoGetErr: nil,
			repoGetCom: communityWithMod,
			wantErr:    apperror.ErrNoFieldsToUpdate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCommunityRepo.EXPECT().
				GetByID(gomock.Any(), tt.req.CommunityID).
				Return(tt.repoGetCom, tt.repoGetErr)

			if tt.repoGetErr == nil && tt.repoGetCom != nil {
				for _, modInput := range tt.req.AddedModerator {
					if modInput.ModeratorID == existingModID.Hex() {
						continue
					}
					mockCommunityRepo.EXPECT().
						IsUserExist(gomock.Any(), modInput.ModeratorID).
						Return(tt.repoUserExists, tt.repoUserExistErr)

					mockMembershipRepo.EXPECT().
						IsMember(gomock.Any(), modInput.ModeratorID, tt.requesterID).
						Return(tt.repoIsMember, tt.repoMemberErr)

					if tt.wantErr == nil {
						mockCommunityRepo.EXPECT().
							Replace(gomock.Any(), gomock.Any()).
							Return(tt.repoUpdateErr)
					}
				}
			}

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
		repoGetErr  error
		repoGetCom  *model.Community
		repoUpdErr  error
		wantErr     error
	}{
		{
			name:        "successfully remove moderator",
			requesterID: creatorID.Hex(),
			req: &dto.RemoveModeratorRequest{
				CommunityID:      communityID.Hex(),
				RemovedModerator: []string{modID.Hex()},
			},
			repoGetErr: nil,
			repoGetCom: communityWithMod,
			repoUpdErr: nil,
			wantErr:    nil,
		},
		{
			name:        "attempt to remove self",
			requesterID: modID.Hex(),
			req: &dto.RemoveModeratorRequest{
				CommunityID:      communityID.Hex(),
				RemovedModerator: []string{modID.Hex()},
			},
			repoGetErr: nil,
			repoGetCom: communityWithMod,
			wantErr:    apperror.ErrCannotRemoveModerator,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCommunityRepo.EXPECT().
				GetByID(gomock.Any(), tt.req.CommunityID).
				Return(tt.repoGetCom, tt.repoGetErr)

			if tt.repoGetErr == nil && tt.repoGetCom != nil && tt.wantErr == nil {
				mockCommunityRepo.EXPECT().
					Replace(gomock.Any(), gomock.Any()).
					Return(tt.repoUpdErr)
			}

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

	community := &model.Community{
		ID:         communityID,
		CreateByID: creatorID,
	}

	tests := []struct {
		name        string
		communityID string
		requesterID string
		repoGetErr  error
		repoGetCom  *model.Community
		repoDelErr  error
		wantErr     error
	}{
		{
			name:        "successfully delete community",
			communityID: communityID.Hex(),
			requesterID: creatorID.Hex(),
			repoGetErr:  nil,
			repoGetCom:  community,
			repoDelErr:  nil,
			wantErr:     nil,
		},
		{
			name:        "attempt to delete without ownership",
			communityID: communityID.Hex(),
			requesterID: nonCreatorID.Hex(),
			repoGetErr:  nil,
			repoGetCom:  community,
			wantErr:     apperror.ErrForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCommunityRepo.EXPECT().
				GetByID(gomock.Any(), tt.communityID).
				Return(tt.repoGetCom, tt.repoGetErr)

			if tt.repoGetErr == nil && tt.repoGetCom != nil && tt.wantErr == nil {
				mockCommunityRepo.EXPECT().
					Delete(gomock.Any(), tt.communityID).
					Return(tt.repoDelErr)
			}

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

	community := &model.Community{
		ID:         communityID,
		CreateByID: moderatorID,
	}

	tests := []struct {
		name        string
		requesterID string
		req         *dto.CommunityBanUserRequest
		repoGetErr  error
		repoGetCom  *model.Community
		repoBanErr  error
		wantErr     error
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
			repoGetErr: nil,
			repoGetCom: community,
			repoBanErr: nil,
			wantErr:    nil,
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
			repoGetErr: nil,
			repoGetCom: community,
			wantErr:    apperror.ErrForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCommunityRepo.EXPECT().
				GetByID(gomock.Any(), tt.req.CommunityID).
				Return(tt.repoGetCom, tt.repoGetErr)

			if tt.repoGetErr == nil && tt.repoGetCom != nil && tt.wantErr == nil {
				mockCommunityRepo.EXPECT().
					BanUser(gomock.Any(), gomock.AssignableToTypeOf(&model.CommunityBan{})).
					Return(tt.repoBanErr)
			}

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

	community := &model.Community{
		ID:         communityID,
		CreateByID: moderatorID,
	}

	bannedUser := &model.User{
		ID:       bannedUserID,
		Username: "BannedUser",
	}

	tests := []struct {
		name         string
		communityID  string
		banTypeStr   string
		expired      bool
		requesterID  string
		repoGetErr   error
		repoGetCom   *model.Community
		repoUsersErr error
		repoUsers    []*model.User
		wantErr      error
		validate     func(t *testing.T, users []*model.User)
	}{
		{
			name:         "retrieve banned users list",
			communityID:  communityID.Hex(),
			banTypeStr:   "banned",
			expired:      false,
			requesterID:  moderatorID.Hex(),
			repoGetErr:   nil,
			repoGetCom:   community,
			repoUsersErr: nil,
			repoUsers:    []*model.User{bannedUser},
			wantErr:      nil,
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
			mockCommunityRepo.EXPECT().
				GetByID(gomock.Any(), tt.communityID).
				Return(tt.repoGetCom, tt.repoGetErr)

			if tt.repoGetErr == nil && tt.repoGetCom != nil && tt.wantErr == nil {
				mockCommunityRepo.EXPECT().
					GetBannedUsers(gomock.Any(), tt.communityID, tt.expired).
					Return(tt.repoUsers, tt.repoUsersErr)
			}

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

	community := &model.Community{
		ID:         communityID,
		CreateByID: moderatorID,
	}

	tests := []struct {
		name         string
		userID       string
		communityID  string
		requesterID  string
		repoGetErr   error
		repoGetCom   *model.Community
		repoUnbanErr error
		wantErr      error
	}{
		{
			name:         "successfully unban user",
			userID:       userID.Hex(),
			communityID:  communityID.Hex(),
			requesterID:  moderatorID.Hex(),
			repoGetErr:   nil,
			repoGetCom:   community,
			repoUnbanErr: nil,
			wantErr:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCommunityRepo.EXPECT().
				GetByID(gomock.Any(), tt.communityID).
				Return(tt.repoGetCom, tt.repoGetErr)

			if tt.repoGetErr == nil && tt.repoGetCom != nil && tt.wantErr == nil {
				mockCommunityRepo.EXPECT().
					UnbanUser(gomock.Any(), tt.userID, tt.communityID).
					Return(tt.repoUnbanErr)
			}

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
		name          string
		communityID   string
		moderatorID   string
		page          int
		pageSize      int
		repoGetComErr error
		repoGetCom    *model.Community
		repoFindErr   error
		repoPosts     []*model.Post
		repoTotal     int64
		repoUsersErr  error
		repoUsers     []*model.User
		wantErr       error
		validate      func(t *testing.T, resp *dto.PaginatedPostsResponse)
	}{
		{
			name:          "retrieve pending posts queue",
			communityID:   communityID.Hex(),
			moderatorID:   moderatorID.Hex(),
			page:          1,
			pageSize:      20,
			repoGetComErr: nil,
			repoGetCom:    community,
			repoFindErr:   nil,
			repoPosts:     []*model.Post{pendingPost},
			repoTotal:     1,
			repoUsersErr:  nil,
			repoUsers:     []*model.User{author},
			wantErr:       nil,
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
			mockCommunityRepo.EXPECT().
				GetByID(gomock.Any(), tt.communityID).
				Return(tt.repoGetCom, tt.repoGetComErr)

			if tt.repoGetComErr == nil && tt.repoGetCom != nil && tt.wantErr == nil {
				mockPostRepo.EXPECT().
					Find(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(tt.repoPosts, tt.repoTotal, tt.repoFindErr)

				if tt.repoFindErr == nil && len(tt.repoPosts) > 0 {
					mockUserRepo.EXPECT().
						GetByIDs(gomock.Any(), gomock.Any()).
						Return(tt.repoUsers, tt.repoUsersErr)
				}
			}

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

	community := &model.Community{
		ID:         communityID,
		CreateByID: moderatorID,
	}

	pendingPost := &model.Post{
		ID:               postID,
		CommunityID:      communityID,
		AuthorID:         authorID,
		ModerationStatus: model.ModerationPending,
	}

	tests := []struct {
		name           string
		communityID    string
		postID         string
		moderatorID    string
		approve        bool
		reason         *string
		repoGetComErr  error
		repoGetCom     *model.Community
		repoGetPostErr error
		repoGetPost    *model.Post
		repoUpdateErr  error
		wantErr        error
	}{
		{
			name:           "UTC-020: approve pending post",
			communityID:    communityID.Hex(),
			postID:         postID.Hex(),
			moderatorID:    moderatorID.Hex(),
			approve:        true,
			reason:         nil,
			repoGetComErr:  nil,
			repoGetCom:     community,
			repoGetPostErr: nil,
			repoGetPost:    pendingPost,
			repoUpdateErr:  nil,
			wantErr:        nil,
		},
		{
			name:           "UTC-021: reject pending post with reason",
			communityID:    communityID.Hex(),
			postID:         postID.Hex(),
			moderatorID:    moderatorID.Hex(),
			approve:        false,
			reason:         ptrStr("Inappropriate content"),
			repoGetComErr:  nil,
			repoGetCom:     community,
			repoGetPostErr: nil,
			repoGetPost:    pendingPost,
			repoUpdateErr:  nil,
			wantErr:        nil,
		},
		{
			name:          "UTC-022: attempt to moderate without permissions",
			communityID:   communityID.Hex(),
			postID:        postID.Hex(),
			moderatorID:   nonModeratorID.Hex(),
			approve:       true,
			reason:        nil,
			repoGetComErr: nil,
			repoGetCom:    community,
			wantErr:       apperror.ErrForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCommunityRepo.EXPECT().
				GetByID(gomock.Any(), tt.communityID).
				Return(tt.repoGetCom, tt.repoGetComErr)

			if tt.repoGetComErr == nil && tt.repoGetCom != nil && tt.wantErr == nil {
				mockPostRepo.EXPECT().
					GetByID(gomock.Any(), tt.postID).
					Return(tt.repoGetPost, tt.repoGetPostErr)

				if tt.repoGetPostErr == nil && tt.repoGetPost != nil {
					mockPostRepo.EXPECT().
						UpdateByID(gomock.Any(), tt.postID, gomock.Any()).
						Return(tt.repoUpdateErr)

					mockEventBus.EXPECT().
						Publish(gomock.Any())
				}
			}

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
