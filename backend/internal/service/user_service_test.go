package service

import (
	"errors"
	"testing"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/repo/mocks"
	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func TestUpdateProfile(t *testing.T) {
	// TODO
}

func TestUpdateAvatar(t *testing.T) {
	// TODO
}

func TestUpdateCover(t *testing.T) {
	// TODO
}

func TestChangePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepo(ctrl)
	svc := NewUserService(mockRepo, nil, nil)

	oldPwd := "correct-old"
	oldHashBytes, err := bcrypt.GenerateFromPassword([]byte(oldPwd), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to bcrypt hash oldPwd: %v", err)
	}
	oldHash := string(oldHashBytes)

	// Prepare ObjectIDs for test users
	user1ID := primitive.NewObjectID()
	user2ID := primitive.NewObjectID()
	user3ID := primitive.NewObjectID()
	user4ID := primitive.NewObjectID()
	noUserID := primitive.NewObjectID()

	tests := []struct {
		name          string
		userID        string
		repoGetErr    error
		repoGetUser   *model.User
		incomingOld   string
		newPassword   string
		repoUpdateErr error
		wantErr       error
	}{
		{
			name:       "success change password",
			userID:     user1ID.Hex(),
			repoGetErr: nil,
			repoGetUser: &model.User{
				ID:       user1ID,
				Username: "user1",
				Password: oldHash,
				Provider: model.ProviderLocal,
			},
			incomingOld:   oldPwd,
			newPassword:   "new-secret",
			repoUpdateErr: nil,
			wantErr:       nil,
		},
		{
			name:        "user not found error",
			userID:      noUserID.Hex(),
			repoGetErr:  mongo.ErrNoDocuments,
			repoGetUser: nil,
			wantErr:     apperror.ErrUserNotFound,
		},
		{
			name:       "login method mismatch (not local provider)",
			userID:     user2ID.Hex(),
			repoGetErr: nil,
			repoGetUser: &model.User{
				ID:       user2ID,
				Username: "user2",
				Provider: model.ProviderGoogle,
				Password: "",
			},
			incomingOld: oldPwd,
			newPassword: "whatever",
			wantErr:     apperror.ErrLoginMethodMismatch,
		},
		{
			name:       "invalid old password",
			userID:     user3ID.Hex(),
			repoGetErr: nil,
			repoGetUser: &model.User{
				ID:       user3ID,
				Username: "user3",
				Password: oldHash,
				Provider: model.ProviderLocal,
			},
			incomingOld: "wrong-old-password",
			newPassword: "newone",
			wantErr:     apperror.ErrInvalidCredentials,
		},
		{
			name:       "update repo returns error",
			userID:     user4ID.Hex(),
			repoGetErr: nil,
			repoGetUser: &model.User{
				ID:       user4ID,
				Username: "user4",
				Password: oldHash,
				Provider: model.ProviderLocal,
			},
			incomingOld:   oldPwd,
			newPassword:   "new-secret",
			repoUpdateErr: errors.New("db update failed"),
			wantErr:       errors.New("db update failed"),
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable

		t.Run(tt.name, func(t *testing.T) {
			mockRepo.EXPECT().
				GetByID(gomock.Any(), tt.userID).
				Return(tt.repoGetUser, tt.repoGetErr)

			if tt.repoGetErr == nil && tt.repoGetUser != nil {
				if tt.wantErr == nil || (tt.wantErr != nil && tt.wantErr.Error() == "db update failed") {
					mockRepo.EXPECT().
						Update(gomock.Any(), gomock.AssignableToTypeOf(&model.User{})).
						Return(tt.repoGetUser, tt.repoUpdateErr)
				}
			}

			err := svc.ChangePassword(tt.userID, tt.incomingOld, tt.newPassword)

			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.wantErr)
				} else if err.Error() != tt.wantErr.Error() {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
			}
		})
	}
}

func TestGetSingleUser(t *testing.T) {
	// TODO
}

func TestGetMultipleUsers(t *testing.T) {
	// TODO
}

func TestUserSetting(t *testing.T) {
	// TODO
}
