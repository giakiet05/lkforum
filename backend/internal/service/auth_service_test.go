package service

import (
	"errors"
	"testing"
	"time"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/auth"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/repo/mocks"
	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepo(ctrl)
	mockTokenService := auth.NewMockTokenServiceInterface(ctrl)

	svc := NewAuthService(mockUserRepo, nil, nil, nil, mockTokenService)

	userID := primitive.NewObjectID()
	email := "user@example.com"
	username := "testuser"
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	validUser := &model.User{
		ID:         userID,
		Username:   username,
		Email:      email,
		Password:   string(hashedPassword),
		Provider:   model.ProviderLocal,
		IsVerified: true,
		IsBanned:   false,
		Role:       model.UserRole,
	}

	bannedUser := &model.User{
		ID:         userID,
		Username:   username,
		Email:      email,
		Password:   string(hashedPassword),
		Provider:   model.ProviderLocal,
		IsVerified: true,
		IsBanned:   true,
		BanUntil:   nil, // Permanent ban
		Role:       model.UserRole,
	}

	expiredBanTime := time.Now().Add(-1 * time.Hour)
	expiredBanUser := &model.User{
		ID:         userID,
		Username:   username,
		Email:      email,
		Password:   string(hashedPassword),
		Provider:   model.ProviderLocal,
		IsVerified: true,
		IsBanned:   true,
		BanUntil:   &expiredBanTime,
		Role:       model.UserRole,
	}

	unverifiedUser := &model.User{
		ID:         userID,
		Username:   username,
		Email:      email,
		Password:   string(hashedPassword),
		Provider:   model.ProviderLocal,
		IsVerified: false,
		IsBanned:   false,
		Role:       model.UserRole,
	}

	googleUser := &model.User{
		ID:         userID,
		Username:   username,
		Email:      email,
		Provider:   model.ProviderGoogle,
		IsVerified: true,
		IsBanned:   false,
		Role:       model.UserRole,
	}

	tests := []struct {
		name       string
		identifier string
		password   string
		mockSetup  func()
		wantErr    error
		wantUser   bool
	}{
		{
			name:       "success login with email",
			identifier: email,
			password:   password,
			mockSetup: func() {
				mockUserRepo.EXPECT().
					GetByEmail(gomock.Any(), email).
					Return(validUser, nil)
			},
			wantErr:  nil,
			wantUser: true,
		},
		{
			name:       "success login with username",
			identifier: username,
			password:   password,
			mockSetup: func() {
				mockUserRepo.EXPECT().
					GetByUsername(gomock.Any(), username).
					Return(validUser, nil)
			},
			wantErr:  nil,
			wantUser: true,
		},
		{
			name:       "error user not found by email",
			identifier: email,
			password:   password,
			mockSetup: func() {
				mockUserRepo.EXPECT().
					GetByEmail(gomock.Any(), email).
					Return(nil, mongo.ErrNoDocuments)
			},
			wantErr:  apperror.ErrInvalidCredentials,
			wantUser: false,
		},
		{
			name:       "error user not found by username",
			identifier: username,
			password:   password,
			mockSetup: func() {
				mockUserRepo.EXPECT().
					GetByUsername(gomock.Any(), username).
					Return(nil, mongo.ErrNoDocuments)
			},
			wantErr:  apperror.ErrInvalidCredentials,
			wantUser: false,
		},
		{
			name:       "error wrong password",
			identifier: email,
			password:   "wrongpassword",
			mockSetup: func() {
				mockUserRepo.EXPECT().
					GetByEmail(gomock.Any(), email).
					Return(validUser, nil)
			},
			wantErr:  apperror.ErrInvalidCredentials,
			wantUser: false,
		},
		{
			name:       "error login method mismatch",
			identifier: email,
			password:   password,
			mockSetup: func() {
				mockUserRepo.EXPECT().
					GetByEmail(gomock.Any(), email).
					Return(googleUser, nil)
			},
			wantErr:  apperror.ErrLoginMethodMismatch,
			wantUser: false,
		},
		{
			name:       "error email not verified",
			identifier: email,
			password:   password,
			mockSetup: func() {
				mockUserRepo.EXPECT().
					GetByEmail(gomock.Any(), email).
					Return(unverifiedUser, nil)
			},
			wantErr:  apperror.ErrEmailNotVerified,
			wantUser: false,
		},
		{
			name:       "error user banned permanently",
			identifier: email,
			password:   password,
			mockSetup: func() {
				mockUserRepo.EXPECT().
					GetByEmail(gomock.Any(), email).
					Return(bannedUser, nil)
			},
			wantErr:  apperror.ErrUserInactive,
			wantUser: false,
		},
		{
			name:       "success user ban expired auto unban",
			identifier: email,
			password:   password,
			mockSetup: func() {
				mockUserRepo.EXPECT().
					GetByEmail(gomock.Any(), email).
					Return(expiredBanUser, nil)

				// Expect Update call to unban user
				mockUserRepo.EXPECT().
					Update(gomock.Any(), gomock.AssignableToTypeOf(&model.User{})).
					DoAndReturn(func(ctx interface{}, user *model.User) (*model.User, error) {
						if user.IsBanned {
							t.Error("Expected user to be unbanned")
						}
						if user.BanUntil != nil {
							t.Error("Expected BanUntil to be nil")
						}
						return user, nil
					})
			},
			wantErr:  nil,
			wantUser: true,
		},
		{
			name:       "error repository failure",
			identifier: email,
			password:   password,
			mockSetup: func() {
				mockUserRepo.EXPECT().
					GetByEmail(gomock.Any(), email).
					Return(nil, errors.New("database error"))
			},
			wantErr:  errors.New("database error"),
			wantUser: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			user, accessToken, refreshToken, err := svc.Login(tt.identifier, tt.password)

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

			if tt.wantUser {
				if user == nil {
					t.Error("expected user, got nil")
					return
				}
				if accessToken == "" {
					t.Error("expected access token, got empty string")
				}
				if refreshToken == "" {
					t.Error("expected refresh token, got empty string")
				}
			}
		})
	}
}

func TestLogout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTokenService := auth.NewMockTokenServiceInterface(ctrl)
	svc := NewAuthService(nil, nil, nil, nil, mockTokenService)

	// Set global TokenSvc to avoid JWT parse validation errors
	// The extractJTI function only needs to parse, not validate
	oldTokenSvc := auth.TokenSvc
	auth.TokenSvc = nil // Disable validation in extractJTI
	defer func() { auth.TokenSvc = oldTokenSvc }()

	// Generate valid tokens for testing
	userID := primitive.NewObjectID()
	accessToken, refreshToken, _ := auth.GenerateToken(userID.Hex(), string(model.UserRole))

	tests := []struct {
		name         string
		accessToken  string
		refreshToken string
		mockSetup    func()
		wantErr      error
	}{
		{
			name:         "success logout",
			accessToken:  accessToken,
			refreshToken: refreshToken,
			mockSetup: func() {
				// Expect InvalidateToken to be called twice (access + refresh)
				mockTokenService.EXPECT().
					InvalidateToken(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).
					Times(2)
			},
			wantErr: nil,
		},
		{
			name:         "error invalid access token",
			accessToken:  "invalid-token",
			refreshToken: refreshToken,
			mockSetup:    func() {},
			wantErr:      apperror.ErrInvalidToken,
		},
		{
			name:         "error invalid refresh token",
			accessToken:  accessToken,
			refreshToken: "invalid-token",
			mockSetup: func() {
				// Access token invalidation succeeds, but refresh token parsing fails
			},
			wantErr: apperror.ErrInvalidToken,
		},
		{
			name:         "error token service failure on access token",
			accessToken:  accessToken,
			refreshToken: refreshToken,
			mockSetup: func() {
				mockTokenService.EXPECT().
					InvalidateToken(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("redis error"))
			},
			wantErr: errors.New("redis error"),
		},
		{
			name:         "error token service failure on refresh token",
			accessToken:  accessToken,
			refreshToken: refreshToken,
			mockSetup: func() {
				// First call (access token) succeeds
				mockTokenService.EXPECT().
					InvalidateToken(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)

				// Second call (refresh token) fails
				mockTokenService.EXPECT().
					InvalidateToken(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("redis error"))
			},
			wantErr: errors.New("redis error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := svc.Logout(tt.accessToken, tt.refreshToken)

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
			}
		})
	}
}

func TestLogoutNilTokenService(t *testing.T) {
	svc := NewAuthService(nil, nil, nil, nil, nil)

	// Set global TokenSvc to nil for this test
	oldTokenSvc := auth.TokenSvc
	auth.TokenSvc = nil
	defer func() { auth.TokenSvc = oldTokenSvc }()

	userID := primitive.NewObjectID()
	accessToken, refreshToken, _ := auth.GenerateToken(userID.Hex(), string(model.UserRole))

	err := svc.Logout(accessToken, refreshToken)

	if err != apperror.ErrInternal {
		t.Errorf("expected ErrInternal when tokenService is nil, got %v", err)
	}
}
