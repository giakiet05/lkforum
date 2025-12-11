package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/platform/bus"
	"github.com/giakiet05/lkforum/internal/repo/mocks"
	"github.com/giakiet05/lkforum/internal/test"
	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func TestUpdateProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepo(ctrl)
	svc := NewUserService(mockUserRepo, nil, nil)

	// Helpers
	ptrStr := func(s string) *string { return &s }
	ptrTime := func(t time.Time) *time.Time { return &t }
	ptrGender := func(g model.Gender) *model.Gender { return &g }
	ptrProvince := func(v model.VNProvince) *model.VNProvince { return &v }

	parseDate := func(s string) time.Time {
		d, _ := time.Parse("2006-01-02", s)
		return d
	}

	userID := primitive.NewObjectID()
	notFoundID := primitive.NewObjectID()

	// Existing user in DB
	oldUser := &model.User{
		ID:       userID,
		Username: "user1",
		Email:    "user1@example.com",
		Role:     model.UserRole,
		RoleContent: model.RoleContent{
			AsUser: &model.UserRoleContent{
				Bio:         ptrStr("Old bio"),
				Gender:      ptrGender(model.GenderMale),
				DateOfBirth: ptrTime(parseDate("2000-01-01")),
				Location:    ptrProvince(model.ProvinceHaNoi),
				Interests:   []model.Interest{"gaming"},
				SocialLinks: &model.SocialLinks{
					Website:  ptrStr("https://oldsite.com"),
					Facebook: ptrStr("https://fb.com/olduser"),
					YouTube:  ptrStr("https://youtube.com/olduser"),
					GitHub:   ptrStr("https://github.com/olduser"),
				},
			},
		},
	}

	tests := []struct {
		name          string
		userID        string
		updateReq     *dto.UserProfileUpdateRequest
		repoGetUser   *model.User
		repoGetErr    error
		repoUpdateErr error
		wantErr       error
		wantProfile   *dto.UserProfileResponse
	}{
		{
			name:   "success update one field",
			userID: userID.Hex(),
			updateReq: &dto.UserProfileUpdateRequest{
				Bio: ptrStr("Updated bio only"),
			},
			repoGetUser:   oldUser,
			repoGetErr:    nil,
			repoUpdateErr: nil,
			wantErr:       nil,
			wantProfile: &dto.UserProfileResponse{
				Bio:         ptrStr("Updated bio only"),
				Gender:      ptrStr("male"),
				DateOfBirth: ptrTime(parseDate("2000-01-01")),
				Location:    ptrStr("Hà Nội"),
				Interests:   []string{"gaming"},
				SocialLinks: &model.SocialLinks{
					Website:  ptrStr("https://oldsite.com"),
					Facebook: ptrStr("https://fb.com/olduser"),
					YouTube:  ptrStr("https://youtube.com/olduser"),
					GitHub:   ptrStr("https://github.com/olduser"),
				},
			},
		},
		{
			name:   "success clear field",
			userID: userID.Hex(),
			updateReq: &dto.UserProfileUpdateRequest{
				Bio:         ptrStr(""),
				Gender:      ptrStr(""),
				DateOfBirth: ptrStr(""),
				Location:    ptrStr(""),
				Interests:   nil,
				SocialLinks: &dto.SocialLinksInput{
					Website:  ptrStr(""),
					Facebook: ptrStr("https://fb.com/olduser"),
					YouTube:  ptrStr(""),
					GitHub:   ptrStr("https://github.com/olduser"),
				},
			},
			repoGetUser:   oldUser,
			repoGetErr:    nil,
			repoUpdateErr: nil,
			wantErr:       nil,
			wantProfile: &dto.UserProfileResponse{
				Bio:         nil,
				Gender:      nil,
				DateOfBirth: nil,
				Location:    nil,
				Interests:   nil,
				SocialLinks: &model.SocialLinks{
					Website:  nil,
					Facebook: ptrStr("https://fb.com/olduser"),
					YouTube:  nil,
					GitHub:   ptrStr("https://github.com/olduser"),
				},
			},
		},
		{
			name:   "success full update",
			userID: userID.Hex(),
			updateReq: &dto.UserProfileUpdateRequest{
				Bio:         ptrStr("New bio"),
				Gender:      ptrStr("female"),
				DateOfBirth: ptrStr("1995-05-15"),
				Location:    ptrStr("TPHCM"),
				Interests:   []string{"Lập trình", "Du lịch"},
				SocialLinks: &dto.SocialLinksInput{
					Website:  ptrStr("https://newsite.com"),
					Facebook: ptrStr("https://fb.com/newuser"),
					YouTube:  ptrStr("https://youtube.com/newuser"),
					GitHub:   ptrStr("https://github.com/newuser"),
				},
			},
			repoGetUser:   oldUser,
			repoGetErr:    nil,
			repoUpdateErr: nil,
			wantErr:       nil,
			wantProfile: &dto.UserProfileResponse{
				Bio:         ptrStr("New bio"),
				Gender:      ptrStr("female"),
				DateOfBirth: ptrTime(parseDate("1995-05-15")),
				Location:    ptrStr("TPHCM"),
				Interests:   []string{"Lập trình", "Du lịch"},
				SocialLinks: &model.SocialLinks{
					Website:  ptrStr("https://newsite.com"),
					Facebook: ptrStr("https://fb.com/newuser"),
					YouTube:  ptrStr("https://youtube.com/newuser"),
					GitHub:   ptrStr("https://github.com/newuser"),
				},
			},
		},
		{
			name:        "error user not found",
			userID:      notFoundID.Hex(),
			updateReq:   &dto.UserProfileUpdateRequest{Bio: ptrStr("X")},
			repoGetUser: nil,
			repoGetErr:  mongo.ErrNoDocuments,
			wantErr:     apperror.ErrUserNotFound,
		},
		{
			name:   "error string too long",
			userID: userID.Hex(),
			updateReq: &dto.UserProfileUpdateRequest{
				Bio: ptrStr(strings.Repeat("A", 15000)),
			},
			repoGetUser: oldUser,
			repoGetErr:  nil,
			wantErr:     apperror.ErrBadRequest,
		},
		{
			name:   "error invalid social link url",
			userID: userID.Hex(),
			updateReq: &dto.UserProfileUpdateRequest{
				SocialLinks: &dto.SocialLinksInput{
					Website:  ptrStr("not-a-valid-url"),
					Facebook: ptrStr("not-a-valid-url"),
					YouTube:  ptrStr("not-a-valid-url"),
					GitHub:   ptrStr("not-a-valid-url"),
				},
			},
			repoGetUser: oldUser,
			repoGetErr:  nil,
			wantErr:     apperror.ErrBadRequest,
		},
		{
			name:   "error invalid date of birth format",
			userID: userID.Hex(),
			updateReq: &dto.UserProfileUpdateRequest{
				DateOfBirth: ptrStr("31-13-2020"),
			},
			repoGetUser: oldUser,
			repoGetErr:  nil,
			wantErr:     apperror.ErrInvalidDateFormat,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			//Setup mock with deep copy
			mockUserRepo.EXPECT().
				GetByID(gomock.Any(), tt.userID).
				DoAndReturn(func(ctx context.Context, id string) (*model.User, error) {
					return model.CloneUser(tt.repoGetUser), tt.repoGetErr
				})

			// Update expectations
			if tt.repoGetUser != nil && tt.updateReq != nil && tt.repoUpdateErr == nil {
				mockUserRepo.EXPECT().
					Update(gomock.Any(), gomock.AssignableToTypeOf(&model.User{})).
					DoAndReturn(func(ctx context.Context, user *model.User) (*model.User, error) {
						// Return the modified user object that was passed in
						return user, nil
					})
			} else if tt.repoUpdateErr != nil {
				mockUserRepo.EXPECT().
					Update(gomock.Any(), gomock.AssignableToTypeOf(&model.User{})).
					Return(nil, tt.repoUpdateErr)
			}

			resp, err := svc.UpdateProfile(tt.userID, tt.updateReq)

			// Error tests
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

			// Expected success
			// ---- Field-by-field assertions ----
			got := resp.Profile
			want := tt.wantProfile

			if want == nil {
				return
			}

			// ---- Bio ----
			if (got.Bio == nil) != (want.Bio == nil) {
				t.Errorf("Bio nil mismatch: want %v, got %v", want.Bio, *got.Bio)
			} else if got.Bio != nil && *got.Bio != *want.Bio {
				t.Errorf("Bio mismatch: want %v, got %v", *want.Bio, *got.Bio)
			}

			// ---- Gender ----
			if (got.Gender == nil) != (want.Gender == nil) {
				t.Errorf("Gender nil mismatch: want %v, got %v", want.Gender, *got.Gender)
			} else if got.Gender != nil && *got.Gender != *want.Gender {
				t.Errorf("Gender mismatch: want %v, got %v", *want.Gender, *got.Gender)
			}

			// ---- Location ----
			if (got.Location == nil) != (want.Location == nil) {
				t.Errorf("Location nil mismatch: want %v, got %v", want.Location, *got.Location)
			} else if got.Location != nil && *got.Location != *want.Location {
				t.Errorf("Location mismatch: want %v, got %v", *want.Location, *got.Location)
			}

			// ---- Interests ----
			if !test.EqualStringSlices(got.Interests, want.Interests) {
				t.Errorf("Interests mismatch: want %v, got %v", want.Interests, got.Interests)
			}

			// ---- SocialLinks ----
			if (got.SocialLinks == nil) != (want.SocialLinks == nil) {
				t.Errorf("SocialLinks nil mismatch: want %v, got %v", want.SocialLinks, *got.SocialLinks)
			} else if got.SocialLinks != nil && want.SocialLinks != nil {

				// Website
				if (got.SocialLinks.Website == nil) != (want.SocialLinks.Website == nil) {
					t.Errorf("Website nil mismatch: want %v, got %v", want.SocialLinks.Website, *got.SocialLinks.Website)
				} else if got.SocialLinks.Website != nil && *got.SocialLinks.Website != *want.SocialLinks.Website {
					t.Errorf("Website mismatch: want %v, got %v", *want.SocialLinks.Website, *got.SocialLinks.Website)
				}

				// Facebook
				if (got.SocialLinks.Facebook == nil) != (want.SocialLinks.Facebook == nil) {
					t.Errorf("Facebook nil mismatch: want %v, got %v", want.SocialLinks.Facebook, *got.SocialLinks.Facebook)
				} else if got.SocialLinks.Facebook != nil && *got.SocialLinks.Facebook != *want.SocialLinks.Facebook {
					t.Errorf("Facebook mismatch: want %v, got %v", *want.SocialLinks.Facebook, *got.SocialLinks.Facebook)
				}

				// YouTube
				if (got.SocialLinks.YouTube == nil) != (want.SocialLinks.YouTube == nil) {
					t.Errorf("YouTube nil mismatch: want %v, got %v", want.SocialLinks.YouTube, *got.SocialLinks.YouTube)
				} else if got.SocialLinks.YouTube != nil && *got.SocialLinks.YouTube != *want.SocialLinks.YouTube {
					t.Errorf("YouTube mismatch: want %v, got %v", *want.SocialLinks.YouTube, *got.SocialLinks.YouTube)
				}

				// GitHub
				if (got.SocialLinks.GitHub == nil) != (want.SocialLinks.GitHub == nil) {
					t.Errorf("GitHub nil mismatch: want %v, got %v", want.SocialLinks.GitHub, *got.SocialLinks.GitHub)
				} else if got.SocialLinks.GitHub != nil && *got.SocialLinks.GitHub != *want.SocialLinks.GitHub {
					t.Errorf("GitHub mismatch: want %v, got %v", *want.SocialLinks.GitHub, *got.SocialLinks.GitHub)
				}
			}
		})
	}
}

func TestUpdateAvatarAndCover(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepo(ctrl)
	mockEventBus := bus.NewMockEventBus(ctrl) // Create mock event bus
	svc := NewUserService(mockUserRepo, mockEventBus, nil)

	userID := primitive.NewObjectID()
	notFoundID := primitive.NewObjectID()

	now := time.Now()

	oldUser := &model.User{
		ID:       userID,
		Username: "user1",
		Email:    "user1@example.com",
		Role:     model.UserRole,
		RoleContent: model.RoleContent{
			AsUser: &model.UserRoleContent{
				Avatar: &model.Image{
					URL:        "https://old.com/a.png",
					PublicID:   "old-avatar-id",
					UploadedAt: now,
				},
				Cover: &model.Image{
					URL:        "https://old.com/c.png",
					PublicID:   "old-cover-id",
					UploadedAt: now,
				},
			},
		},
	}

	tests := []struct {
		name          string
		action        string // "avatar", "cover", "delete_avatar", "delete_cover"
		userID        string
		imageURL      string
		publicID      string
		repoGetUser   *model.User
		repoGetErr    error
		repoUpdateErr error
		wantErr       error
		wantAvatar    *model.Image
		wantCover     *model.Image
	}{
		{
			name:        "success update avatar",
			action:      "avatar",
			userID:      userID.Hex(),
			imageURL:    "https://cdn.com/new-avatar.png",
			publicID:    "new-avatar-id",
			repoGetUser: oldUser,
			wantAvatar: &model.Image{
				URL:      "https://cdn.com/new-avatar.png",
				PublicID: "new-avatar-id",
			},
			wantCover: oldUser.RoleContent.AsUser.Cover,
		},
		{
			name:        "error invalid avatar URL",
			action:      "avatar",
			userID:      userID.Hex(),
			imageURL:    "bad-url",
			publicID:    "id-x",
			repoGetUser: oldUser,
			wantErr:     apperror.ErrBadRequest,
		},
		{
			name:       "error user not found updating avatar",
			action:     "avatar",
			userID:     notFoundID.Hex(),
			imageURL:   "https://cdn.com/x.png",
			publicID:   "x",
			repoGetErr: mongo.ErrNoDocuments,
			wantErr:    apperror.ErrUserNotFound,
		},
		{
			name:        "success update cover",
			action:      "cover",
			userID:      userID.Hex(),
			imageURL:    "https://cdn.com/new-cover.png",
			publicID:    "new-cover-id",
			repoGetUser: oldUser,
			wantAvatar:  oldUser.RoleContent.AsUser.Avatar,
			wantCover: &model.Image{
				URL:      "https://cdn.com/new-cover.png",
				PublicID: "new-cover-id",
			},
		},
		{
			name:        "error invalid cover URL",
			action:      "cover",
			userID:      userID.Hex(),
			imageURL:    "invalid",
			publicID:    "pid-xx",
			repoGetUser: oldUser,
			wantErr:     apperror.ErrBadRequest,
		},
		{
			name:        "success delete avatar",
			action:      "delete_avatar",
			userID:      userID.Hex(),
			repoGetUser: oldUser,
			wantAvatar:  nil,
			wantCover:   oldUser.RoleContent.AsUser.Cover,
		},
		{
			name:       "delete avatar user not found",
			action:     "delete_avatar",
			userID:     notFoundID.Hex(),
			repoGetErr: mongo.ErrNoDocuments,
			wantErr:    apperror.ErrUserNotFound,
		},
		{
			name:        "success delete cover",
			action:      "delete_cover",
			userID:      userID.Hex(),
			repoGetUser: oldUser,
			wantAvatar:  oldUser.RoleContent.AsUser.Avatar,
			wantCover:   nil,
		},
		{
			name:       "delete cover user not found",
			action:     "delete_cover",
			userID:     notFoundID.Hex(),
			repoGetErr: mongo.ErrNoDocuments,
			wantErr:    apperror.ErrUserNotFound,
		},
		{
			name:          "repo update error in avatar update",
			action:        "avatar",
			userID:        userID.Hex(),
			imageURL:      "https://cdn.com/err.png",
			publicID:      "pid-err",
			repoGetUser:   oldUser,
			repoUpdateErr: errors.New("db error"),
			wantErr:       apperror.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup GetByID mock with deep copy
			mockUserRepo.EXPECT().
				GetByID(gomock.Any(), tt.userID).
				DoAndReturn(func(ctx context.Context, id string) (*model.User, error) {
					return model.CloneUser(tt.repoGetUser), tt.repoGetErr
				})

			// Setup Update mock - only when we expect Update to be called
			shouldCallUpdate := tt.repoGetUser != nil && tt.repoGetErr == nil
			if shouldCallUpdate {
				if tt.repoUpdateErr != nil {
					// Update returns error
					mockUserRepo.EXPECT().
						Update(gomock.Any(), gomock.Any()).
						Return(nil, tt.repoUpdateErr)
				} else {
					// Update succeeds
					mockUserRepo.EXPECT().
						Update(gomock.Any(), gomock.Any()).
						DoAndReturn(func(ctx context.Context, user *model.User) (*model.User, error) {
							return user, nil
						})

					// Mock EventBus.Publish for avatar updates only (when successful)
					if tt.action == "avatar" || tt.action == "delete_avatar" {
						mockEventBus.EXPECT().
							Publish(gomock.Any())
					}
				}
			}

			var resp *dto.UserResponse
			var err error

			// Call the appropriate service method
			switch tt.action {
			case "avatar":
				resp, err = svc.UpdateAvatar(tt.userID, tt.imageURL, tt.publicID)
			case "cover":
				resp, err = svc.UpdateCover(tt.userID, tt.imageURL, tt.publicID)
			case "delete_avatar":
				resp, err = svc.DeleteAvatar(tt.userID)
			case "delete_cover":
				resp, err = svc.DeleteCover(tt.userID)
			}

			// Verify errors
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected err %v, got %v", tt.wantErr, err)
			}
			if tt.wantErr != nil {
				return
			}

			// Nil-safe avatar check
			gotAvatar := resp.Profile.Avatar
			wantAvatar := tt.wantAvatar

			if (gotAvatar == nil) != (wantAvatar == nil) {
				t.Fatalf("Avatar nil mismatch: want %v, got %v", wantAvatar, gotAvatar)
			}
			if gotAvatar != nil && wantAvatar != nil {
				if gotAvatar.URL != wantAvatar.URL {
					t.Errorf("Avatar URL mismatch: want %v, got %v", wantAvatar.URL, gotAvatar.URL)
				}
				if gotAvatar.PublicID != wantAvatar.PublicID {
					t.Errorf("Avatar PublicID mismatch: want %v, got %v", wantAvatar.PublicID, gotAvatar.PublicID)
				}
			}

			// Nil-safe cover check
			gotCover := resp.Profile.Cover
			wantCover := tt.wantCover

			if (gotCover == nil) != (wantCover == nil) {
				t.Fatalf("Cover nil mismatch: want %v, got %v", wantCover, gotCover)
			}
			if gotCover != nil && wantCover != nil {
				if gotCover.URL != wantCover.URL {
					t.Errorf("Cover URL mismatch: want %v, got %v", wantCover.URL, gotCover.URL)
				}
				if gotCover.PublicID != wantCover.PublicID {
					t.Errorf("Cover PublicID mismatch: want %v, got %v", wantCover.PublicID, gotCover.PublicID)
				}
			}
		})
	}
}

func TestChangePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepo(ctrl)
	svc := NewUserService(mockUserRepo, nil, nil)

	oldPwd := "OLD_PASSWORD"
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
			newPassword:   "NEW_PASSWORD",
			repoUpdateErr: nil,
			wantErr:       nil,
		},
		{
			name:        "error user not found",
			userID:      noUserID.Hex(),
			repoGetErr:  mongo.ErrNoDocuments,
			repoGetUser: nil,
			wantErr:     apperror.ErrUserNotFound,
		},
		{
			name:       "error login method mismatch (not local provider)",
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
			name:       "error invalid old password",
			userID:     user3ID.Hex(),
			repoGetErr: nil,
			repoGetUser: &model.User{
				ID:       user3ID,
				Username: "user3",
				Password: oldHash,
				Provider: model.ProviderLocal,
			},
			incomingOld: "WRONG_PASS",
			newPassword: "NEW_PASSWORD",
			wantErr:     apperror.ErrInvalidCredentials,
		},
		{
			name:       "error update repo returns",
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
			mockUserRepo.EXPECT().
				GetByID(gomock.Any(), tt.userID).
				Return(tt.repoGetUser, tt.repoGetErr)

			if tt.repoGetErr == nil && tt.repoGetUser != nil {
				if tt.wantErr == nil || (tt.wantErr != nil && tt.wantErr.Error() == "db update failed") {
					mockUserRepo.EXPECT().
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

func TestUpdateUserSetting(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepo(ctrl)
	mockEventBus := bus.NewMockEventBus(ctrl)
	svc := NewUserService(mockUserRepo, mockEventBus, nil)

	userID := primitive.NewObjectID()
	notFoundID := primitive.NewObjectID()

	// Helper functions for pointer values
	ptrString := func(s string) *string { return &s }
	ptrBool := func(b bool) *bool { return &b }

	userWithSettings := &model.User{
		ID:         userID,
		Username:   "user1",
		Email:      "user1@example.com",
		Reputation: 100,
		Provider:   model.ProviderLocal,
		Role:       model.UserRole,
		IsVerified: true,
		Settings: &model.UserSettings{
			Appearance: model.AppearanceSettings{
				Theme:    "light",
				FontSize: "medium",
			},
			Notifications: model.NotificationSettings{
				InAppEnabled:    true,
				EmailEnabled:    true,
				NotifyOnComment: true,
				NotifyOnMention: true,
				NotifyOnUpvote:  false,
				NotifyOnMessage: true,
			},
			Privacy: model.PrivacySettings{
				ShowProfile:         true,
				ShowEmail:           false,
				ShowPostHistory:     true,
				AllowDirectMessages: true,
				AllowMentions:       true,
			},
			Content: model.ContentSettings{
				AllowNSFW: false,
			},
		},
	}

	userWithoutSettings := &model.User{
		ID:         userID,
		Username:   "user2",
		Email:      "user2@example.com",
		Reputation: 50,
		Provider:   model.ProviderLocal,
		Role:       model.UserRole,
		IsVerified: true,
		Settings:   nil,
	}

	tests := []struct {
		name          string
		userID        string
		request       *dto.UpdateSettingsRequest
		repoGetUser   *model.User
		repoGetErr    error
		repoUpdateErr error
		wantErr       error
		validate      func(t *testing.T, resp *dto.SettingsResponse)
	}{
		{
			name:        "success update appearance theme",
			userID:      userID.Hex(),
			repoGetUser: userWithSettings,
			request: &dto.UpdateSettingsRequest{
				Appearance: &dto.AppearanceSettingsInput{
					Theme: ptrString("dark"),
				},
			},
			validate: func(t *testing.T, resp *dto.SettingsResponse) {
				if resp.Appearance.Theme != "dark" {
					t.Errorf("expected theme 'dark', got '%s'", resp.Appearance.Theme)
				}
				if resp.Appearance.FontSize != "medium" {
					t.Errorf("expected fontSize 'medium', got '%s'", resp.Appearance.FontSize)
				}
			},
		},
		{
			name:        "success update appearance font size",
			userID:      userID.Hex(),
			repoGetUser: userWithSettings,
			request: &dto.UpdateSettingsRequest{
				Appearance: &dto.AppearanceSettingsInput{
					FontSize: ptrString("large"),
				},
			},
			validate: func(t *testing.T, resp *dto.SettingsResponse) {
				if resp.Appearance.FontSize != "large" {
					t.Errorf("expected fontSize 'large', got '%s'", resp.Appearance.FontSize)
				}
				if resp.Appearance.Theme != "light" {
					t.Errorf("expected theme 'light', got '%s'", resp.Appearance.Theme)
				}
			},
		},
		{
			name:        "error invalid theme",
			userID:      userID.Hex(),
			repoGetUser: userWithSettings,
			request: &dto.UpdateSettingsRequest{
				Appearance: &dto.AppearanceSettingsInput{
					Theme: ptrString("invalid-theme"),
				},
			},
			wantErr: apperror.ErrBadRequest,
		},
		{
			name:        "error invalid font size",
			userID:      userID.Hex(),
			repoGetUser: userWithSettings,
			request: &dto.UpdateSettingsRequest{
				Appearance: &dto.AppearanceSettingsInput{
					FontSize: ptrString("invalid-size"),
				},
			},
			wantErr: apperror.ErrBadRequest,
		},
		{
			name:        "success update notification settings",
			userID:      userID.Hex(),
			repoGetUser: userWithSettings,
			request: &dto.UpdateSettingsRequest{
				Notifications: &dto.NotificationSettingsInput{
					InAppEnabled:    ptrBool(false),
					EmailEnabled:    ptrBool(false),
					NotifyOnComment: ptrBool(false),
					NotifyOnUpvote:  ptrBool(true),
				},
			},
			validate: func(t *testing.T, resp *dto.SettingsResponse) {
				if resp.Notifications.InAppEnabled != false {
					t.Errorf("expected InAppEnabled false, got %v", resp.Notifications.InAppEnabled)
				}
				if resp.Notifications.EmailEnabled != false {
					t.Errorf("expected EmailEnabled false, got %v", resp.Notifications.EmailEnabled)
				}
				if resp.Notifications.NotifyOnComment != false {
					t.Errorf("expected NotifyOnComment false, got %v", resp.Notifications.NotifyOnComment)
				}
				if resp.Notifications.NotifyOnUpvote != true {
					t.Errorf("expected NotifyOnUpvote true, got %v", resp.Notifications.NotifyOnUpvote)
				}
				// Unchanged values
				if resp.Notifications.NotifyOnMention != true {
					t.Errorf("expected NotifyOnMention true, got %v", resp.Notifications.NotifyOnMention)
				}
			},
		},
		{
			name:        "success update privacy settings",
			userID:      userID.Hex(),
			repoGetUser: userWithSettings,
			request: &dto.UpdateSettingsRequest{
				Privacy: &dto.PrivacySettingsInput{
					ShowProfile:         ptrBool(false),
					ShowEmail:           ptrBool(true),
					AllowDirectMessages: ptrBool(false),
				},
			},
			validate: func(t *testing.T, resp *dto.SettingsResponse) {
				if resp.Privacy.ShowProfile != false {
					t.Errorf("expected ShowProfile false, got %v", resp.Privacy.ShowProfile)
				}
				if resp.Privacy.ShowEmail != true {
					t.Errorf("expected ShowEmail true, got %v", resp.Privacy.ShowEmail)
				}
				if resp.Privacy.AllowDirectMessages != false {
					t.Errorf("expected AllowDirectMessages false, got %v", resp.Privacy.AllowDirectMessages)
				}
				// Unchanged values
				if resp.Privacy.ShowPostHistory != true {
					t.Errorf("expected ShowPostHistory true, got %v", resp.Privacy.ShowPostHistory)
				}
			},
		},
		{
			name:        "success update content settings",
			userID:      userID.Hex(),
			repoGetUser: userWithSettings,
			request: &dto.UpdateSettingsRequest{
				Content: &dto.ContentSettingsInput{
					AllowNSFW: ptrBool(true),
				},
			},
			validate: func(t *testing.T, resp *dto.SettingsResponse) {
				if resp.Content.AllowNSFW != true {
					t.Errorf("expected AllowNSFW true, got %v", resp.Content.AllowNSFW)
				}
			},
		},
		{
			name:        "success update multiple categories",
			userID:      userID.Hex(),
			repoGetUser: userWithSettings,
			request: &dto.UpdateSettingsRequest{
				Appearance: &dto.AppearanceSettingsInput{
					Theme: ptrString("dark"),
				},
				Notifications: &dto.NotificationSettingsInput{
					EmailEnabled: ptrBool(false),
				},
				Privacy: &dto.PrivacySettingsInput{
					ShowProfile: ptrBool(false),
				},
			},
			validate: func(t *testing.T, resp *dto.SettingsResponse) {
				if resp.Appearance.Theme != "dark" {
					t.Errorf("expected theme 'dark', got '%s'", resp.Appearance.Theme)
				}
				if resp.Notifications.EmailEnabled != false {
					t.Errorf("expected EmailEnabled false, got %v", resp.Notifications.EmailEnabled)
				}
				if resp.Privacy.ShowProfile != false {
					t.Errorf("expected ShowProfile false, got %v", resp.Privacy.ShowProfile)
				}
			},
		},
		{
			name:        "success initialize nil settings",
			userID:      userID.Hex(),
			repoGetUser: userWithoutSettings,
			request: &dto.UpdateSettingsRequest{
				Appearance: &dto.AppearanceSettingsInput{
					Theme: ptrString("dark"),
				},
			},
			validate: func(t *testing.T, resp *dto.SettingsResponse) {
				if resp.Appearance.Theme != "dark" {
					t.Errorf("expected theme 'dark', got '%s'", resp.Appearance.Theme)
				}
				// Should have default values for other settings
				if resp == nil {
					t.Error("expected non-nil settings response")
				}
			},
		},
		{
			name:       "error user not found",
			userID:     notFoundID.Hex(),
			repoGetErr: mongo.ErrNoDocuments,
			request: &dto.UpdateSettingsRequest{
				Appearance: &dto.AppearanceSettingsInput{
					Theme: ptrString("dark"),
				},
			},
			wantErr: apperror.ErrUserNotFound,
		},
		{
			name:        "success empty request no changes",
			userID:      userID.Hex(),
			repoGetUser: userWithSettings,
			request:     &dto.UpdateSettingsRequest{},
			validate: func(t *testing.T, resp *dto.SettingsResponse) {
				// All values should remain unchanged
				if resp.Appearance.Theme != "light" {
					t.Errorf("expected theme 'light', got '%s'", resp.Appearance.Theme)
				}
				if resp.Notifications.InAppEnabled != true {
					t.Errorf("expected InAppEnabled true, got %v", resp.Notifications.InAppEnabled)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup GetByID mock with deep copy
			mockUserRepo.EXPECT().
				GetByID(gomock.Any(), tt.userID).
				DoAndReturn(func(ctx context.Context, id string) (*model.User, error) {
					return model.CloneUser(tt.repoGetUser), tt.repoGetErr
				})

			// Setup Update mock - only when we expect Update to be called
			shouldCallUpdate := tt.repoGetUser != nil && tt.repoGetErr == nil && tt.wantErr == nil
			if shouldCallUpdate {
				if tt.repoUpdateErr != nil {
					// Update returns error
					mockUserRepo.EXPECT().
						Update(gomock.Any(), gomock.Any()).
						Return(nil, tt.repoUpdateErr)
				} else {
					// Update succeeds
					mockUserRepo.EXPECT().
						Update(gomock.Any(), gomock.Any()).
						DoAndReturn(func(ctx context.Context, user *model.User) (*model.User, error) {
							return user, nil
						})
				}
			}

			// Call the service method
			resp, err := svc.UpdateSettings(tt.userID, tt.request)

			// Verify errors
			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantErr)
				}
				if !errors.Is(err, tt.wantErr) && err.Error() != tt.wantErr.Error() {
					t.Fatalf("expected err %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Run custom validation if provided
			if tt.validate != nil {
				tt.validate(t, resp)
			}
		})
	}
}
