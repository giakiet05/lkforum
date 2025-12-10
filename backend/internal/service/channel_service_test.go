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

func TestCreateChannel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChannelRepo := mocks.NewMockChannelRepo(ctrl)
	mockEventBus := bus.NewMockEventBus(ctrl)

	svc := NewChannelService(mockChannelRepo, mockEventBus)

	user1 := primitive.NewObjectID()
	user2 := primitive.NewObjectID()

	tests := []struct {
		name      string
		requester string
		req       *dto.CreateChannelRequest
		repoErr   error
		wantErr   error
		validate  func(t *testing.T, ch *model.Channel)
	}{
		{
			name:      "success create channel",
			requester: user1.Hex(),
			req: &dto.CreateChannelRequest{
				Member1:         user1.Hex(),
				Member1Username: "user1",
				Member1Avatar:   "a1",
				Member2:         user2.Hex(),
				Member2Username: "user2",
				Member2Avatar:   "a2",
			},
			repoErr: nil,
			wantErr: nil,
			validate: func(t *testing.T, ch *model.Channel) {
				if ch == nil {
					t.Fatal("expected channel, got nil")
				}
				if len(ch.Members) != 2 {
					t.Fatalf("expected 2 members, got %d", len(ch.Members))
				}
				if ch.Members[0].Username != "user1" {
					t.Errorf("expected member1 username user1, got %s", ch.Members[0].Username)
				}
			},
		},
		{
			name:      "forbidden when requester not member",
			requester: primitive.NewObjectID().Hex(),
			req: &dto.CreateChannelRequest{
				Member1:         user1.Hex(),
				Member1Username: "user1",
				Member1Avatar:   "a1",
				Member2:         user2.Hex(),
				Member2Username: "user2",
				Member2Avatar:   "a2",
			},
			wantErr: apperror.ErrForbidden,
		},
		{
			name:      "invalid member id",
			requester: user1.Hex(),
			req: &dto.CreateChannelRequest{
				Member1:         "not-a-hex",
				Member1Username: "user1",
				Member1Avatar:   "a1",
				Member2:         user2.Hex(),
				Member2Username: "user2",
				Member2Avatar:   "a2",
			},
			wantErr: errors.New("invalid hex"), // will match by non-nil and not ErrForbidden
		},
		{
			name:      "repo create error",
			requester: user1.Hex(),
			req: &dto.CreateChannelRequest{
				Member1:         user1.Hex(),
				Member1Username: "user1",
				Member1Avatar:   "a1",
				Member2:         user2.Hex(),
				Member2Username: "user2",
				Member2Avatar:   "a2",
			},
			repoErr: errors.New("db error"),
			wantErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr == nil {
				mockChannelRepo.EXPECT().
					Create(gomock.Any(), gomock.AssignableToTypeOf(&model.Channel{})).
					DoAndReturn(func(_ interface{}, ch *model.Channel) (*model.Channel, error) {
						return ch, tt.repoErr
					})
			} else if tt.repoErr != nil {
				mockChannelRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil, tt.repoErr)
			}

			ch, err := svc.CreateChannel(tt.req, tt.requester)

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
				tt.validate(t, ch)
			}
		})
	}
}

func TestGetChannelByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChannelRepo := mocks.NewMockChannelRepo(ctrl)
	mockEventBus := bus.NewMockEventBus(ctrl)

	svc := NewChannelService(mockChannelRepo, mockEventBus)

	chID := primitive.NewObjectID()
	ch := &model.Channel{ID: chID}

	tests := []struct {
		name    string
		id      string
		repoCh  *model.Channel
		repoErr error
		wantErr error
	}{
		{"found", chID.Hex(), ch, nil, nil},
		{"not found", primitive.NewObjectID().Hex(), nil, mongo.ErrNoDocuments, mongo.ErrNoDocuments},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.repoErr == nil {
				mockChannelRepo.EXPECT().GetByID(gomock.Any(), tt.id).Return(tt.repoCh, tt.repoErr)
			} else {
				mockChannelRepo.EXPECT().GetByID(gomock.Any(), tt.id).Return(nil, tt.repoErr)
			}

			res, err := svc.GetChannelByID(tt.id)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if res == nil {
				t.Fatalf("expected channel, got nil")
			}
		})
	}
}

func TestGetChannelsByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChannelRepo := mocks.NewMockChannelRepo(ctrl)
	mockEventBus := bus.NewMockEventBus(ctrl)

	svc := NewChannelService(mockChannelRepo, mockEventBus)

	userID := primitive.NewObjectID()
	channel := model.Channel{ID: primitive.NewObjectID()}

	tests := []struct {
		name        string
		userID      string
		requesterID string
		page        int
		pageSize    int
		repoCh      []model.Channel
		repoTotal   int64
		repoErr     error
		wantErr     error
		validate    func(t *testing.T, resp *dto.PaginatedChannelsResponse)
	}{
		{
			name:        "success",
			userID:      userID.Hex(),
			requesterID: userID.Hex(),
			page:        1,
			pageSize:    10,
			repoCh:      []model.Channel{channel},
			repoTotal:   1,
			repoErr:     nil,
			wantErr:     nil,
			validate: func(t *testing.T, resp *dto.PaginatedChannelsResponse) {
				if resp == nil {
					t.Fatal("expected response, got nil")
				}
				if resp.Pagination.Total != 1 {
					t.Errorf("expected total 1, got %d", resp.Pagination.Total)
				}
			},
		},
		{
			name:        "forbidden",
			userID:      userID.Hex(),
			requesterID: primitive.NewObjectID().Hex(),
			page:        1,
			pageSize:    10,
			wantErr:     apperror.ErrForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr == nil {
				mockChannelRepo.EXPECT().GetByUserID(gomock.Any(), tt.userID, tt.page, tt.pageSize).Return(tt.repoCh, tt.repoTotal, tt.repoErr)
			}

			resp, err := svc.GetChannelsByUserID(tt.userID, tt.requesterID, tt.page, tt.pageSize)

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

func TestGetChannelByBothUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChannelRepo := mocks.NewMockChannelRepo(ctrl)
	mockEventBus := bus.NewMockEventBus(ctrl)

	svc := NewChannelService(mockChannelRepo, mockEventBus)

	user1 := primitive.NewObjectID()
	user2 := primitive.NewObjectID()
	ch := &model.Channel{ID: primitive.NewObjectID()}

	tests := []struct {
		name        string
		user1       string
		user2       string
		requesterID string
		repoCh      *model.Channel
		repoErr     error
		wantErr     error
	}{
		{"success", user1.Hex(), user2.Hex(), user1.Hex(), ch, nil, nil},
		{"forbidden", user1.Hex(), user2.Hex(), primitive.NewObjectID().Hex(), nil, nil, apperror.ErrForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr == nil {
				mockChannelRepo.EXPECT().GetByBothUserID(gomock.Any(), tt.user1, tt.user2).Return(tt.repoCh, tt.repoErr)
			}

			res, err := svc.GetChannelByBothUserID(tt.user1, tt.user2, tt.requesterID)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if res == nil {
				t.Fatalf("expected channel, got nil")
			}
		})
	}
}

func TestUpdateChannel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChannelRepo := mocks.NewMockChannelRepo(ctrl)
	mockEventBus := bus.NewMockEventBus(ctrl)

	svc := NewChannelService(mockChannelRepo, mockEventBus)

	member1 := primitive.NewObjectID()
	member2 := primitive.NewObjectID()

	channel := &model.Channel{
		ID: primitive.NewObjectID(),
		Members: []model.ChannelMember{
			{UserID: member1, Username: "m1", Avatar: "a1"},
			{UserID: member2, Username: "m2", Avatar: "a2"},
		},
		Settings: []model.ChannelSetting{
			{UserID: member1, Notification: true, TypingIndicator: true},
			{UserID: member2, Notification: true, TypingIndicator: true},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name       string
		requester  string
		req        *dto.UpdateChannelRequest
		repoGetErr error
		repoGetCh  *model.Channel
		repoUpdErr error
		wantErr    error
		validate   func(t *testing.T, ch *model.Channel)
	}{
		{
			name:      "success update notification",
			requester: member1.Hex(),
			req: &dto.UpdateChannelRequest{
				ChannelID:    channel.ID.Hex(),
				Notification: ptrBool(false),
			},
			repoGetErr: nil,
			repoGetCh:  channel,
			repoUpdErr: nil,
			wantErr:    nil,
			validate: func(t *testing.T, ch *model.Channel) {
				if ch == nil {
					t.Fatal("expected channel, got nil")
				}
				// find setting for requester
				for _, s := range ch.Settings {
					if s.UserID.Hex() == member1.Hex() {
						if s.Notification != false {
							t.Errorf("expected notification false, got %v", s.Notification)
						}
					}
				}
			},
		},
		{
			name:      "forbidden when not a member",
			requester: primitive.NewObjectID().Hex(),
			req: &dto.UpdateChannelRequest{
				ChannelID: channel.ID.Hex(),
				Nickname:  ptrStr("nick"),
			},
			repoGetErr: nil,
			repoGetCh:  channel,
			wantErr:    apperror.ErrForbidden,
		},
		{
			name:      "no fields to update",
			requester: member1.Hex(),
			req: &dto.UpdateChannelRequest{
				ChannelID: channel.ID.Hex(),
				// all nil
			},
			repoGetErr: nil,
			repoGetCh:  channel,
			wantErr:    apperror.ErrNoFieldsToUpdate,
		},
		{
			name:      "repo get error",
			requester: member1.Hex(),
			req: &dto.UpdateChannelRequest{
				ChannelID: channel.ID.Hex(),
				Nickname:  ptrStr("n"),
			},
			repoGetErr: mongo.ErrNoDocuments,
			wantErr:    mongo.ErrNoDocuments,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockChannelRepo.EXPECT().GetByID(gomock.Any(), tt.req.ChannelID).Return(tt.repoGetCh, tt.repoGetErr)

			if tt.repoGetErr == nil && tt.repoGetCh != nil && tt.wantErr == nil {
				mockChannelRepo.EXPECT().Update(gomock.Any(), gomock.AssignableToTypeOf(&model.Channel{})).Return(tt.repoGetCh, tt.repoUpdErr)
			}

			ch, err := svc.UpdateChannel(tt.req, tt.requester)

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
				tt.validate(t, ch)
			}
		})
	}
}

func TestDeleteChannel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChannelRepo := mocks.NewMockChannelRepo(ctrl)
	mockEventBus := bus.NewMockEventBus(ctrl)

	svc := NewChannelService(mockChannelRepo, mockEventBus)

	channelID := primitive.NewObjectID().Hex()
	userID := primitive.NewObjectID().Hex()

	tests := []struct {
		name    string
		chID    string
		userID  string
		repoErr error
		wantErr error
	}{
		{"success", channelID, userID, nil, nil},
		{"repo error", channelID, userID, errors.New("db"), errors.New("db")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.repoErr == nil {
				mockChannelRepo.EXPECT().Delete(gomock.Any(), tt.chID, tt.userID).Return(nil)
			} else {
				mockChannelRepo.EXPECT().Delete(gomock.Any(), tt.chID, tt.userID).Return(tt.repoErr)
			}

			err := svc.DeleteChannel(tt.chID, tt.userID)

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

func TestHandleNewAvatar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChannelRepo := mocks.NewMockChannelRepo(ctrl)
	mockEventBus := bus.NewMockEventBus(ctrl)

	svc := NewChannelService(mockChannelRepo, mockEventBus).(*channelService)

	mockEvent := bus.NewMockEvent(ctrl)
	userID := primitive.NewObjectID().Hex()
	newAvatar := "avatar.png"

	mockEvent.EXPECT().Topic().Return(bus.TopicUserChangeAvatar)
	mockEvent.EXPECT().Payload().Return(map[string]interface{}{"user_id": userID, "new_avatar": newAvatar})
	mockChannelRepo.EXPECT().UpdateUserAvatar(gomock.Any(), userID, newAvatar).Return(nil)

	svc.handleNewAvatar(mockEvent)
}

func ptrBool(b bool) *bool { return &b }
