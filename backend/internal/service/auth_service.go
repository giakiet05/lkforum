package service

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"time"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/auth"
	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/platform/email"
	"github.com/giakiet05/lkforum/internal/repo"
	"github.com/giakiet05/lkforum/internal/util"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// GoogleAuthStatus defines the result status of a Google callback processing.
const (
	StatusLoginSuccess  = "LOGIN_SUCCESS"
	StatusSetupRequired = "SETUP_REQUIRED"
)

// GoogleAuthResult is the result of processing a Google OAuth callback.
type GoogleAuthResult struct {
	Status       string
	User         *model.User
	AccessToken  string
	RefreshToken string
	SetupToken   string
}

type AuthService interface {
	// Local Auth
	RegisterUser(username, email, password string) (*model.User, error)
	Login(identifier, password string) (*model.User, string, string, error)
	VerifyEmail(email, otp string) (*model.User, string, string, error)
	ResendVerificationEmail(email string) error
	RefreshToken(refreshToken string) (string, string, error)

	// Google OAuth
	ProcessGoogleCallback(code string) (*GoogleAuthResult, error)
	CompleteGoogleSetup(setupToken, username string) (*model.User, string, string, error)
}

type authService struct {
	userRepo    repo.UserRepo
	emailSender email.Sender
}

func NewAuthService(userRepo repo.UserRepo, emailSender email.Sender) AuthService {
	return &authService{
		userRepo:    userRepo,
		emailSender: emailSender,
	}
}

// --- Local Authentication ---

func (s *authService) RegisterUser(username, email, password string) (*model.User, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	if _, err := s.userRepo.GetByUsername(ctx, username); !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, apperror.ErrUsernameExists
	}
	if _, err := s.userRepo.GetByEmail(ctx, email); !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, apperror.ErrEmailExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	otp := generateOTP()
	otpExpiresAt := time.Now().Add(time.Duration(config.Cfg.OTPExpirationMinutes) * time.Minute)

	user := &model.User{
		Username:                  username,
		Email:                     email,
		Password:                  string(hashedPassword),
		Provider:                  model.ProviderLocal,
		Role:                      model.UserRole,
		IsVerified:                false,
		VerificationCode:          otp,
		VerificationCodeExpiresAt: &otpExpiresAt,
		CreatedAt:                 time.Now(),
	}

	createdUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	go func() {
		if err := s.emailSender.SendVerificationEmail(email, otp); err != nil {
			fmt.Printf("CRITICAL: Failed to send verification email to %s: %v\n", email, err)
		}
	}()

	return createdUser, nil
}

func (s *authService) Login(identifier, password string) (*model.User, string, string, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()
	var user *model.User
	var err error

	if isEmail(identifier) {
		user, err = s.userRepo.GetByEmail(ctx, identifier)
	} else {
		user, err = s.userRepo.GetByUsername(ctx, identifier)
	}

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, "", "", apperror.ErrInvalidCredentials
		}
		return nil, "", "", err
	}

	if user.Provider != model.ProviderLocal {
		return nil, "", "", apperror.ErrLoginMethodMismatch
	}

	if user.Password == "" || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return nil, "", "", apperror.ErrInvalidCredentials
	}

	if !user.IsVerified {
		return nil, "", "", apperror.ErrEmailNotVerified
	}

	accessToken, refreshToken, err := auth.GenerateToken(user.ID.Hex(), string(user.Role))
	if err != nil {
		return nil, "", "", err
	}
	return user, accessToken, refreshToken, nil
}

func (s *authService) VerifyEmail(email, otp string) (*model.User, string, string, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, "", "", apperror.ErrUserNotFound
		}
		return nil, "", "", err
	}

	if user.IsVerified {
		return nil, "", "", apperror.ErrEmailAlreadyVerified
	}

	if user.VerificationCode != otp {
		return nil, "", "", apperror.ErrInvalidOTP
	}

	if user.VerificationCodeExpiresAt == nil || user.VerificationCodeExpiresAt.Before(time.Now()) {
		return nil, "", "", apperror.ErrOTPExpired
	}

	user.IsVerified = true
	user.VerificationCode = ""
	user.VerificationCodeExpiresAt = nil

	updatedUser, err := s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, "", "", err
	}

	accessToken, refreshToken, err := auth.GenerateToken(updatedUser.ID.Hex(), string(updatedUser.Role))
	if err != nil {
		return nil, "", "", err
	}

	return updatedUser, accessToken, refreshToken, nil
}

func (s *authService) ResendVerificationEmail(email string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return apperror.ErrUserNotFound
		}
		return err
	}

	if user.IsVerified {
		return apperror.ErrEmailAlreadyVerified
	}

	otp := generateOTP()
	otpExpiresAt := time.Now().Add(time.Duration(config.Cfg.OTPExpirationMinutes) * time.Minute)
	user.VerificationCode = otp
	user.VerificationCodeExpiresAt = &otpExpiresAt

	_, err = s.userRepo.Update(ctx, user)
	if err != nil {
		return err
	}

	go func() {
		if err := s.emailSender.SendVerificationEmail(email, otp); err != nil {
			fmt.Printf("CRITICAL: Failed to resend verification email to %s: %v\n", email, err)
		}
	}()

	return nil
}

func (s *authService) RefreshToken(refreshToken string) (string, string, error) {
	userID, err := auth.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", "", apperror.ErrUserNotFound
		}
		return "", "", err
	}

	accessToken, newRefreshToken, err := auth.GenerateToken(user.ID.Hex(), string(user.Role))
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

// --- Google OAuth ---

func (s *authService) ProcessGoogleCallback(code string) (*GoogleAuthResult, error) {
	userInfo, err := auth.GetGoogleUserInfo(code)
	if err != nil {
		return nil, err
	}

	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	user, err := s.userRepo.GetByEmail(ctx, userInfo.Email)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			setupToken, err := auth.CreateSetupToken(userInfo)
			if err != nil {
				return nil, err
			}
			return &GoogleAuthResult{Status: StatusSetupRequired, SetupToken: setupToken}, nil
		}
		return nil, err
	}

	if user.Provider != model.ProviderGoogle {
		return nil, apperror.ErrLoginMethodMismatch
	}

	accessToken, refreshToken, err := auth.GenerateToken(user.ID.Hex(), string(user.Role))
	if err != nil {
		return nil, err
	}

	return &GoogleAuthResult{
		Status:       StatusLoginSuccess,
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) CompleteGoogleSetup(setupToken, username string) (*model.User, string, string, error) {
	claims, err := auth.ParseSetupToken(setupToken)
	if err != nil {
		return nil, "", "", err
	}

	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	if _, err := s.userRepo.GetByUsername(ctx, username); !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, "", "", apperror.ErrUsernameExists
	}

	if _, err := s.userRepo.GetByEmail(ctx, claims.Email); !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, "", "", apperror.ErrEmailExists
	}

	newUser := &model.User{
		Username:   username,
		Email:      claims.Email,
		Provider:   model.ProviderGoogle,
		ProviderID: claims.GoogleID,
		Role:       model.UserRole,
		IsVerified: true,
		CreatedAt:  time.Now(),
		RoleContent: model.RoleContent{
			AsUser: &model.UserRoleContent{
				Avatar: model.Image{URL: claims.Picture},
			},
		},
	}

	createdUser, err := s.userRepo.Create(ctx, newUser)
	if err != nil {
		return nil, "", "", err
	}

	accessToken, refreshToken, err := auth.GenerateToken(createdUser.ID.Hex(), string(createdUser.Role))
	if err != nil {
		return nil, "", "", err
	}

	return createdUser, accessToken, refreshToken, nil
}

// --- Helpers ---

func isEmail(s string) bool {
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(s)
}

func generateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}
