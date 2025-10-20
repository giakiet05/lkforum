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
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/email"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/repo"
	"github.com/giakiet05/lkforum/internal/util"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	// Auth
	RegisterUser(username, email, password string) (*model.User, error)
	Login(identifier, password string) (*model.User, string, string, error)
	VerifyEmail(email, otp string) (*model.User, string, string, error)
	ResendVerificationEmail(email string) error
	RefreshToken(refreshToken string) (string, string, error)

	// User Management
	UpdateUser(user *model.User) (*model.User, error)
	DeleteUser(id string) error
	ChangePassword(userID, oldPassword, newPassword string) error

	// User Retrieval
	GetAllUsers() ([]*model.User, error)
	GetUserByID(id string) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	GetUsers(page, pageSize int) (*dto.PaginatedUsersResponse, error)
}

type userService struct {
	userRepo    repo.UserRepo
	emailSender email.Sender
}

func NewUserService(userRepo repo.UserRepo, emailSender email.Sender) UserService {
	return &userService{
		userRepo:    userRepo,
		emailSender: emailSender,
	}
}

func (s *userService) RegisterUser(username, email, password string) (*model.User, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	if user, err := s.userRepo.GetByUsername(ctx, username); err == nil && user != nil {
		return nil, apperror.ErrUsernameExists
	}
	if user, err := s.userRepo.GetByEmail(ctx, email); err == nil && user != nil {
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
		Role:                      model.UserRole,
		IsVerified:                false,
		VerificationCode:          otp,
		VerificationCodeExpiresAt: &otpExpiresAt,
		CreateAt:                  time.Now(),
	}

	createdUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	// Asynchronously send the verification email
	go func() {
		if err := s.emailSender.SendVerificationEmail(email, otp); err != nil {
			fmt.Printf("CRITICAL: Failed to send verification email to %s: %v\n", email, err)
		}
	}()

	return createdUser, nil
}

func (s *userService) VerifyEmail(email, otp string) (*model.User, string, string, error) {
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

func (s *userService) ResendVerificationEmail(email string) error {
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

	// Asynchronously send the new verification email
	go func() {
		if err := s.emailSender.SendVerificationEmail(email, otp); err != nil {
			fmt.Printf("CRITICAL: Failed to resend verification email to %s: %v\n", email, err)
		}
	}()

	return nil
}

func (s *userService) Login(identifier, password string) (*model.User, string, string, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()
	var user *model.User
	var err error

	isEmail := isEmail(identifier)
	if isEmail {
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
	if user == nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
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

// ... (rest of the functions are unchanged)

func (s *userService) UpdateUser(user *model.User) (*model.User, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()
	updatedUser, err := s.userRepo.Update(ctx, user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.ErrUserNotFound
		}
		return nil, err
	}
	return updatedUser, nil
}

func (s *userService) DeleteUser(id string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	if auth.TokenSvc != nil {
		if err := auth.TokenSvc.InvalidateAllUserTokens(ctx, id); err != nil {
			fmt.Printf("Failed to invalidate tokens for user %s: %v\n", id, err)
		}
	}

	err := s.userRepo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return apperror.ErrUserNotFound
		}
		return err
	}
	return nil
}

func (s *userService) ChangePassword(userID, oldPassword, newPassword string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return apperror.ErrUserNotFound
		}
		return err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)) != nil {
		return apperror.ErrInvalidCredentials
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	_, err = s.userRepo.Update(ctx, user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return apperror.ErrUserNotFound
		}
		return err
	}
	return nil
}

func (s *userService) GetUsers(page, pageSize int) (*dto.PaginatedUsersResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	users, total, err := s.userRepo.GetPaginated(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	userResponses := dto.FromUsers(users)

	return &dto.PaginatedUsersResponse{
		Users: userResponses,
		Pagination: dto.Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	}, nil
}

func (s *userService) RefreshToken(refreshToken string) (string, string, error) {
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

func (s *userService) GetAllUsers() ([]*model.User, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()
	return s.userRepo.GetAll(ctx)
}

func (s *userService) GetUserByID(id string) (*model.User, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *userService) GetUserByUsername(username string) (*model.User, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *userService) GetUserByEmail(email string) (*model.User, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func isEmail(s string) bool {
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(s)
}

func generateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}
