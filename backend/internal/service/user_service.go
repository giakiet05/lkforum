package service

import (
	"errors"
	"fmt"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/auth"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/repo"
	"github.com/giakiet05/lkforum/internal/util"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// UserService handles business logic related to user management.
type UserService interface {
	UpdateUser(user *model.User) (*model.User, error)
	DeleteUser(id string) error
	ChangePassword(userID, oldPassword, newPassword string) error

	GetUserByID(id string) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	GetUsers(page, pageSize int) (*dto.PaginatedUsersResponse, error)
	GetAllUsers() ([]*model.User, error)
}

type userService struct {
	userRepo repo.UserRepo
}

func NewUserService(userRepo repo.UserRepo) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) UpdateUser(user *model.User) (*model.User, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()
	// TODO: Add validation logic here. For example, check if the new username or email is already taken by another user.
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
			// Log the error but continue with deletion
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

	// Ensure this user is a local user and has a password
	if user.Provider != model.ProviderLocal || user.Password == "" {
		return apperror.ErrLoginMethodMismatch // Cannot change password for OAuth users this way
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

func (s *userService) GetAllUsers() ([]*model.User, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()
	return s.userRepo.GetAll(ctx)
}
