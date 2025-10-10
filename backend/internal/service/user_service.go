package service

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/auth"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/repo"
	"github.com/giakiet05/lkforum/internal/util"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterUser(username, email, password string) (*model.User, string, string, error)
	Login(identifier, password string) (*model.User, string, string, error)
	UpdateUser(user *model.User) (*model.User, error)
	DeleteUser(id string) error

	GetAllUsers() ([]*model.User, error)
	GetUserByID(id string) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	ChangePassword(userID, oldPassword, newPassword string) error
	GetUsers(page, pageSize int) (*dto.PaginatedUsersResponse, error)
	RefreshToken(refreshToken string) (string, string, error)
}

type userService struct {
	userRepo repo.UserRepo
}

func NewUserService(userRepo repo.UserRepo) UserService {
	return &userService{
		userRepo: userRepo,
	}
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

func (s *userService) RegisterUser(username, email, password string) (*model.User, string, string, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Check if username exists
	if user, err := s.userRepo.GetByUsername(ctx, username); err == nil && user != nil {
		return nil, "", "", apperror.ErrUsernameExists
	}
	// Check if email exists
	if user, err := s.userRepo.GetByEmail(ctx, email); err == nil && user != nil {
		return nil, "", "", apperror.ErrEmailExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", "", err
	}
	user := &model.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		Role:     model.UserRole,
	}
	createdUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, "", "", err
	}
	accessToken, refreshToken, err := auth.GenerateToken(createdUser.ID.Hex(), string(createdUser.Role))
	if err != nil {
		return nil, "", "", err
	}
	return createdUser, accessToken, refreshToken, nil
}

func (s *userService) Login(identifier, password string) (*model.User, string, string, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()
	var user *model.User
	var err error

	// Use regex to check if identifier is an email
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
	accessToken, refreshToken, err := auth.GenerateToken(user.ID.Hex(), string(user.Role))
	if err != nil {
		return nil, "", "", err
	}
	return user, accessToken, refreshToken, nil
}

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

	// First invalidate all tokens for this user
	if auth.TokenSvc != nil {
		if err := auth.TokenSvc.InvalidateAllUserTokens(ctx, id); err != nil {
			// Log the error but continue with deletion
			fmt.Printf("Failed to invalidate tokens for user %s: %v\n", id, err)
		}
	}

	// Then delete the user from the database
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

// isEmail checks if the given string is a valid email address format
func isEmail(s string) bool {
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(s)
}
