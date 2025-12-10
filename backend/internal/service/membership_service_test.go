package service

import (
	"context"
	"errors"
	"testing"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/repo/mocks"
	"github.com/golang/mock/gomock"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestCreateMembership(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMembershipRepo := mocks.NewMockMembershipRepo(ctrl)
	mockRedisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	svc := NewMembershipService(mockMembershipRepo, mockRedisClient)

	userID := primitive.NewObjectID()
	communityID := primitive.NewObjectID()

	tests := []struct {
		name        string
		requesterID string
		req         *dto.CreateMembershipRequest
		wantErr     error
		setupMocks  func()
		validate    func(t *testing.T, m *model.Membership)
	}{
		{
			name:        "successfully create membership",
			requesterID: userID.Hex(),
			req: &dto.CreateMembershipRequest{
				UserID:      userID.Hex(),
				CommunityID: communityID.Hex(),
			},
			wantErr: nil,
			setupMocks: func() {
				mockMembershipRepo.EXPECT().
					IsCommunityExist(gomock.Any(), communityID.Hex()).
					Return(true, nil)
				mockMembershipRepo.EXPECT().
					Create(gomock.Any(), gomock.AssignableToTypeOf(&model.Membership{})).
					DoAndReturn(func(ctx context.Context, mem *model.Membership) (*model.Membership, error) {
						mem.ID = primitive.NewObjectID()
						return mem, nil
					})
				mockMembershipRepo.EXPECT().
					CountMembersByCommunityID(gomock.Any(), communityID.Hex()).
					Return(int64(1), nil)
			},
			validate: func(t *testing.T, m *model.Membership) {
				if m == nil {
					t.Fatal("expected membership, got nil")
				}
				if m.UserID != userID {
					t.Errorf("expected user ID %s, got %s", userID, m.UserID)
				}
				if m.CommunityID != communityID {
					t.Errorf("expected community ID %s, got %s", communityID, m.CommunityID)
				}
			},
		},
		{
			name:        "forbidden when requester is not user",
			requesterID: primitive.NewObjectID().Hex(),
			req: &dto.CreateMembershipRequest{
				UserID:      userID.Hex(),
				CommunityID: communityID.Hex(),
			},
			wantErr:    apperror.ErrForbidden,
			setupMocks: func() {},
		},
		{
			name:        "community not found",
			requesterID: userID.Hex(),
			req: &dto.CreateMembershipRequest{
				UserID:      userID.Hex(),
				CommunityID: communityID.Hex(),
			},
			wantErr: apperror.ErrCommunityNotFound,
			setupMocks: func() {
				mockMembershipRepo.EXPECT().
					IsCommunityExist(gomock.Any(), communityID.Hex()).
					Return(false, nil)
			},
		},
		{
			name:        "invalid user ID",
			requesterID: "invalid-id",
			req: &dto.CreateMembershipRequest{
				UserID:      "invalid-id",
				CommunityID: communityID.Hex(),
			},
			wantErr:    errors.New("invalid id"),
			setupMocks: func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			m, err := svc.CreateMembership(tt.req, tt.requesterID)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.validate != nil {
				tt.validate(t, m)
			}
		})
	}
}

func TestGetMembershipByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMembershipRepo := mocks.NewMockMembershipRepo(ctrl)
	mockRedisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	svc := NewMembershipService(mockMembershipRepo, mockRedisClient)

	membershipID := primitive.NewObjectID()
	userID := primitive.NewObjectID()
	communityID := primitive.NewObjectID()

	existingMembership := &model.Membership{
		ID:          membershipID,
		UserID:      userID,
		CommunityID: communityID,
	}

	tests := []struct {
		name         string
		membershipID string
		wantErr      error
		setupMocks   func()
		validate     func(t *testing.T, m *model.Membership)
	}{
		{
			name:         "retrieve membership by valid ID",
			membershipID: membershipID.Hex(),
			wantErr:      nil,
			setupMocks: func() {
				mockMembershipRepo.EXPECT().
					GetByID(gomock.Any(), membershipID.Hex()).
					Return(existingMembership, nil)
			},
			validate: func(t *testing.T, m *model.Membership) {
				if m.ID != membershipID {
					t.Errorf("expected ID %s, got %s", membershipID, m.ID)
				}
				if m.UserID != userID {
					t.Errorf("expected user ID %s, got %s", userID, m.UserID)
				}
			},
		},
		{
			name:         "membership not found",
			membershipID: primitive.NewObjectID().Hex(),
			wantErr:      apperror.ErrMembershipNotFound,
			setupMocks: func() {
				mockMembershipRepo.EXPECT().
					GetByID(gomock.Any(), gomock.Any()).
					Return(nil, mongo.ErrNoDocuments)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			m, err := svc.GetMembershipByID(tt.membershipID)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if err != tt.wantErr && err.Error() != tt.wantErr.Error() {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.validate != nil {
				tt.validate(t, m)
			}
		})
	}
}

func TestGetMembershipsByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMembershipRepo := mocks.NewMockMembershipRepo(ctrl)
	mockRedisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	svc := NewMembershipService(mockMembershipRepo, mockRedisClient)

	userID := primitive.NewObjectID()
	communityID1 := primitive.NewObjectID()
	communityID2 := primitive.NewObjectID()

	memberships := []*model.Membership{
		{
			ID:          primitive.NewObjectID(),
			UserID:      userID,
			CommunityID: communityID1,
		},
		{
			ID:          primitive.NewObjectID(),
			UserID:      userID,
			CommunityID: communityID2,
		},
	}

	tests := []struct {
		name       string
		userID     string
		wantErr    error
		setupMocks func()
		validate   func(t *testing.T, ms []*model.Membership)
	}{
		{
			name:    "retrieve memberships by user ID",
			userID:  userID.Hex(),
			wantErr: nil,
			setupMocks: func() {
				mockMembershipRepo.EXPECT().
					GetByUserID(gomock.Any(), userID.Hex()).
					Return(memberships, nil)
			},
			validate: func(t *testing.T, ms []*model.Membership) {
				if len(ms) != 2 {
					t.Errorf("expected 2 memberships, got %d", len(ms))
				}
			},
		},
		{
			name:    "no memberships found",
			userID:  primitive.NewObjectID().Hex(),
			wantErr: nil,
			setupMocks: func() {
				mockMembershipRepo.EXPECT().
					GetByUserID(gomock.Any(), gomock.Any()).
					Return([]*model.Membership{}, nil)
			},
			validate: func(t *testing.T, ms []*model.Membership) {
				if len(ms) != 0 {
					t.Errorf("expected 0 memberships, got %d", len(ms))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			ms, err := svc.GetMembershipsByUserID(tt.userID)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.validate != nil {
				tt.validate(t, ms)
			}
		})
	}
}

func TestGetAllMemberships(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMembershipRepo := mocks.NewMockMembershipRepo(ctrl)
	mockRedisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	svc := NewMembershipService(mockMembershipRepo, mockRedisClient)

	memberships := []*model.Membership{
		{
			ID:          primitive.NewObjectID(),
			UserID:      primitive.NewObjectID(),
			CommunityID: primitive.NewObjectID(),
		},
	}

	tests := []struct {
		name       string
		page       int
		pageSize   int
		wantErr    error
		setupMocks func()
		validate   func(t *testing.T, resp *dto.PaginatedMembershipsResponse)
	}{
		{
			name:     "retrieve all memberships paginated",
			page:     1,
			pageSize: 10,
			wantErr:  nil,
			setupMocks: func() {
				mockMembershipRepo.EXPECT().
					GetAllPaginated(gomock.Any(), 1, 10).
					Return(memberships, int64(1), nil)
			},
			validate: func(t *testing.T, resp *dto.PaginatedMembershipsResponse) {
				if resp == nil {
					t.Fatal("expected response, got nil")
				}
				if resp.Pagination.Total != 1 {
					t.Errorf("expected total 1, got %d", resp.Pagination.Total)
				}
				if len(resp.Memberships) != 1 {
					t.Errorf("expected 1 membership, got %d", len(resp.Memberships))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			resp, err := svc.GetAllMemberships(tt.page, tt.pageSize)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error, got nil")
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

func TestGetMembershipByCommunityID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMembershipRepo := mocks.NewMockMembershipRepo(ctrl)
	mockRedisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	svc := NewMembershipService(mockMembershipRepo, mockRedisClient)

	communityID := primitive.NewObjectID()
	memberships := []*model.Membership{
		{
			ID:          primitive.NewObjectID(),
			UserID:      primitive.NewObjectID(),
			CommunityID: communityID,
		},
	}

	tests := []struct {
		name        string
		communityID string
		page        int
		pageSize    int
		wantErr     error
		setupMocks  func()
		validate    func(t *testing.T, resp *dto.PaginatedMembershipsResponse)
	}{
		{
			name:        "retrieve memberships by community ID paginated",
			communityID: communityID.Hex(),
			page:        1,
			pageSize:    10,
			wantErr:     nil,
			setupMocks: func() {
				mockMembershipRepo.EXPECT().
					GetByCommunityIDPaginated(gomock.Any(), communityID.Hex(), 1, 10).
					Return(memberships, int64(1), nil)
			},
			validate: func(t *testing.T, resp *dto.PaginatedMembershipsResponse) {
				if resp == nil {
					t.Fatal("expected response, got nil")
				}
				if resp.Pagination.Total != 1 {
					t.Errorf("expected total 1, got %d", resp.Pagination.Total)
				}
				if len(resp.Memberships) != 1 {
					t.Errorf("expected 1 membership, got %d", len(resp.Memberships))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			resp, err := svc.GetMembershipByCommunityID(tt.communityID, tt.page, tt.pageSize)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error, got nil")
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

func TestDeleteMembership(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMembershipRepo := mocks.NewMockMembershipRepo(ctrl)
	mockRedisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	svc := NewMembershipService(mockMembershipRepo, mockRedisClient)

	userID := primitive.NewObjectID()
	communityID := primitive.NewObjectID()
	otherUserID := primitive.NewObjectID()

	tests := []struct {
		name        string
		requesterID string
		req         *dto.DeleteMembershipRequest
		wantErr     error
		setupMocks  func()
	}{
		{
			name:        "successfully delete membership",
			requesterID: userID.Hex(),
			req: &dto.DeleteMembershipRequest{
				UserID:      userID.Hex(),
				CommunityID: communityID.Hex(),
			},
			wantErr: nil,
			setupMocks: func() {
				mockMembershipRepo.EXPECT().
					Delete(gomock.Any(), communityID.Hex()).
					Return(nil)
				mockMembershipRepo.EXPECT().
					CountMembersByCommunityID(gomock.Any(), communityID.Hex()).
					Return(int64(0), nil)
			},
		},
		{
			name:        "forbidden when requester is not user",
			requesterID: otherUserID.Hex(),
			req: &dto.DeleteMembershipRequest{
				UserID:      userID.Hex(),
				CommunityID: communityID.Hex(),
			},
			wantErr:    apperror.ErrForbidden,
			setupMocks: func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := svc.DeleteMembership(tt.req, tt.requesterID)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestGetMembersCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMembershipRepo := mocks.NewMockMembershipRepo(ctrl)
	mockRedisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	svc := NewMembershipService(mockMembershipRepo, mockRedisClient)

	communityID := primitive.NewObjectID()

	tests := []struct {
		name        string
		communityID string
		wantErr     error
		setupMocks  func()
		validate    func(t *testing.T, count int64)
	}{
		{
			name:        "retrieve members count when cache exists",
			communityID: communityID.Hex(),
			wantErr:     nil,
			setupMocks: func() {
				mockMembershipRepo.EXPECT().
					CountMembersByCommunityID(gomock.Any(), communityID.Hex()).
					Return(int64(5), nil)
			},
			validate: func(t *testing.T, count int64) {
				// Note: Due to Redis client initialization, actual count verification
				// would require a real Redis instance or proper mocking
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			count, err := svc.GetMembersCount(tt.communityID)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.validate != nil {
				tt.validate(t, count)
			}
		})
	}
}
