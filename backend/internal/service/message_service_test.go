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

func TestGetMessageByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMessageRepo := mocks.NewMockMessageRepo(ctrl)
	mockChannelRepo := mocks.NewMockChannelRepo(ctrl)
	mockUserRepo := mocks.NewMockUserRepo(ctrl)
	mockEventBus := bus.NewMockEventBus(ctrl)

	svc := NewMessageService(mockMessageRepo, mockChannelRepo, mockUserRepo, mockEventBus, nil)

	messageID := primitive.NewObjectID()
	channelID := primitive.NewObjectID()
	requesterID := primitive.NewObjectID()

	existingMessage := &model.Message{
		ID:        messageID,
		ChannelID: channelID,
		Content:   "hello",
		IsRead:    false,
		IsSend:    true,
		CreatedAt: time.Now(),
	}

	tests := []struct {
		name            string
		channelID       string
		messageID       string
		requesterID     string
		repoIsMember    bool
		repoIsMemberErr error
		repoGetErr      error
		repoGetMsg      *model.Message
		wantErr         error
	}{
		{
			name:         "success get message by id",
			channelID:    channelID.Hex(),
			messageID:    messageID.Hex(),
			requesterID:  requesterID.Hex(),
			repoIsMember: true,
			repoGetErr:   nil,
			repoGetMsg:   existingMessage,
			wantErr:      nil,
		},
		{
			name:         "requester not member",
			channelID:    channelID.Hex(),
			messageID:    messageID.Hex(),
			requesterID:  requesterID.Hex(),
			repoIsMember: false,
			wantErr:      apperror.ErrForbidden,
		},
		{
			name:         "message not found",
			channelID:    channelID.Hex(),
			messageID:    messageID.Hex(),
			requesterID:  requesterID.Hex(),
			repoIsMember: true,
			repoGetErr:   mongo.ErrNoDocuments,
			wantErr:      apperror.ErrNoMessageFound,
		},
		{
			name:         "repo get error",
			channelID:    channelID.Hex(),
			messageID:    messageID.Hex(),
			requesterID:  requesterID.Hex(),
			repoIsMember: true,
			repoGetErr:   errors.New("db err"),
			wantErr:      errors.New("db err"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// expectations
			mockChannelRepo.EXPECT().IsMember(gomock.Any(), tt.channelID, tt.requesterID).Return(tt.repoIsMember, tt.repoIsMemberErr)

			if tt.repoIsMember && tt.repoGetErr != nil {
				mockMessageRepo.EXPECT().GetByID(gomock.Any(), tt.messageID).Return(nil, tt.repoGetErr)
			} else if tt.repoIsMember && tt.repoGetErr == nil {
				mockMessageRepo.EXPECT().GetByID(gomock.Any(), tt.messageID).Return(tt.repoGetMsg, nil)
			}

			msg, err := svc.GetMessageByID(tt.channelID, tt.messageID, tt.requesterID)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantErr)
				}
				if !errors.Is(err, tt.wantErr) && err.Error() != tt.wantErr.Error() {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if msg == nil {
				t.Fatalf("expected message, got nil")
			}
		})
	}
}

func TestGetMessageFilter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMessageRepo := mocks.NewMockMessageRepo(ctrl)
	mockChannelRepo := mocks.NewMockChannelRepo(ctrl)
	mockUserRepo := mocks.NewMockUserRepo(ctrl)
	mockEventBus := bus.NewMockEventBus(ctrl)

	svc := NewMessageService(mockMessageRepo, mockChannelRepo, mockUserRepo, mockEventBus, nil)

	channelID := primitive.NewObjectID()
	req := &dto.GetMessageFilterQuery{
		ChannelID: channelID.Hex(),
		Page:      1,
		PageSize:  10,
	}

	messages := []model.Message{
		{
			ID:        primitive.NewObjectID(),
			ChannelID: channelID,
			Content:   "one",
			CreatedAt: time.Now(),
		},
		{
			ID:        primitive.NewObjectID(),
			ChannelID: channelID,
			Content:   "two",
			CreatedAt: time.Now(),
		},
	}

	tests := []struct {
		name            string
		query           *dto.GetMessageFilterQuery
		requesterID     string
		repoIsMember    bool
		repoIsMemberErr error
		repoGetErr      error
		repoMsgs        []model.Message
		repoTotal       int64
		wantErr         error
	}{
		{
			name:         "success get filter",
			query:        req,
			requesterID:  "user1",
			repoIsMember: true,
			repoGetErr:   nil,
			repoMsgs:     messages,
			repoTotal:    int64(len(messages)),
			wantErr:      nil,
		},
		{
			name:         "not member",
			query:        req,
			requesterID:  "user1",
			repoIsMember: false,
			wantErr:      apperror.ErrForbidden,
		},
		{
			name:         "repo error",
			query:        req,
			requesterID:  "user1",
			repoIsMember: true,
			repoGetErr:   errors.New("db err"),
			wantErr:      errors.New("db err"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockChannelRepo.EXPECT().IsMember(gomock.Any(), tt.query.ChannelID, tt.requesterID).Return(tt.repoIsMember, tt.repoIsMemberErr)

			if tt.repoIsMember && tt.repoGetErr != nil {
				mockMessageRepo.EXPECT().GetFilter(gomock.Any(), tt.query.ChannelID, tt.query.SenderID, tt.query.SearchContent, tt.query.IsRead, tt.query.IsMedia, tt.query.IsSend, tt.query.Page, tt.query.PageSize).Return(nil, int64(0), tt.repoGetErr)
			} else if tt.repoIsMember && tt.repoGetErr == nil {
				mockMessageRepo.EXPECT().GetFilter(gomock.Any(), tt.query.ChannelID, tt.query.SenderID, tt.query.SearchContent, tt.query.IsRead, tt.query.IsMedia, tt.query.IsSend, tt.query.Page, tt.query.PageSize).Return(tt.repoMsgs, tt.repoTotal, nil)
			}

			res, err := svc.GetMessageFilter(tt.query, tt.requesterID)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantErr)
				}
				if !errors.Is(err, tt.wantErr) && err.Error() != tt.wantErr.Error() {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if res == nil {
				t.Fatalf("expected response, got nil")
			}
			if res.Pagination.Total != tt.repoTotal {
				t.Fatalf("expected total %d, got %d", tt.repoTotal, res.Pagination.Total)
			}
		})
	}
}

func TestDeleteMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMessageRepo := mocks.NewMockMessageRepo(ctrl)
	mockChannelRepo := mocks.NewMockChannelRepo(ctrl)
	mockUserRepo := mocks.NewMockUserRepo(ctrl)
	mockEventBus := bus.NewMockEventBus(ctrl)

	svc := NewMessageService(mockMessageRepo, mockChannelRepo, mockUserRepo, mockEventBus, nil)

	messageID := primitive.NewObjectID()
	channelID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	tests := []struct {
		name            string
		channelID       string
		messageID       string
		requesterID     string
		repoIsMember    bool
		repoIsMemberErr error
		repoIsSend      bool
		repoIsSendErr   error
		repoDeleteErr   error
		wantErr         error
	}{
		{
			name:          "success delete",
			channelID:     channelID.Hex(),
			messageID:     messageID.Hex(),
			requesterID:   userID.Hex(),
			repoIsMember:  true,
			repoIsSend:    true,
			repoDeleteErr: nil,
			wantErr:       nil,
		},
		{
			name:         "not member",
			channelID:    channelID.Hex(),
			messageID:    messageID.Hex(),
			requesterID:  userID.Hex(),
			repoIsMember: false,
			wantErr:      apperror.ErrForbidden,
		},
		{
			name:         "not sent by user",
			channelID:    channelID.Hex(),
			messageID:    messageID.Hex(),
			requesterID:  userID.Hex(),
			repoIsMember: true,
			repoIsSend:   false,
			wantErr:      apperror.ErrForbidden,
		},
		{
			name:          "repo delete error",
			channelID:     channelID.Hex(),
			messageID:     messageID.Hex(),
			requesterID:   userID.Hex(),
			repoIsMember:  true,
			repoIsSend:    true,
			repoDeleteErr: errors.New("db err"),
			wantErr:       errors.New("db err"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockChannelRepo.EXPECT().IsMember(gomock.Any(), tt.channelID, tt.requesterID).Return(tt.repoIsMember, tt.repoIsMemberErr)
			if !tt.repoIsMember {
				// no further expectations
			} else {
				mockMessageRepo.EXPECT().IsSendByUser(gomock.Any(), tt.messageID, tt.requesterID).Return(tt.repoIsSend, tt.repoIsSendErr)
				if tt.repoIsSend && tt.repoDeleteErr != nil {
					mockMessageRepo.EXPECT().Delete(gomock.Any(), tt.messageID).Return(tt.repoDeleteErr)
				} else if tt.repoIsSend && tt.repoDeleteErr == nil {
					mockMessageRepo.EXPECT().Delete(gomock.Any(), tt.messageID).Return(nil)
				}
			}

			err := svc.DeleteMessage(tt.channelID, tt.messageID, tt.requesterID)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantErr)
				}
				if !errors.Is(err, tt.wantErr) && err.Error() != tt.wantErr.Error() {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
