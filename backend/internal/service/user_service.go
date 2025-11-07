package service

import (
	"errors"
	"fmt"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/auth"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/platform/bus"
	"github.com/giakiet05/lkforum/internal/platform/cloudinary"
	"github.com/giakiet05/lkforum/internal/repo"
	"github.com/giakiet05/lkforum/internal/util"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// UserService handles business logic related to user management.
type UserService interface {
	UpdateProfile(userID string, req *dto.UserProfileUpdateRequest) (*dto.UserResponse, error)
	UpdateAvatar(userID string, imageURL string, publicID string) (*dto.UserResponse, error)
	UpdateCover(userID string, imageURL string, publicID string) (*dto.UserResponse, error)
	DeleteUser(id string) error
	ChangePassword(userID, oldPassword, newPassword string) error

	GetUserByID(id string) (*dto.UserResponse, error)
	GetUserByUsername(username string) (*dto.UserResponse, error)
	GetUserByEmail(email string) (*dto.UserResponse, error)
	GetUsers(page, pageSize int) (*dto.PaginatedUsersResponse, error)
	GetAllUsers() ([]*model.User, error)
}

type userService struct {
	userRepo repo.UserRepo
	eventBus *bus.EventBus
}

func NewUserService(userRepo repo.UserRepo, bus *bus.EventBus) UserService {
	return &userService{
		userRepo: userRepo,
		eventBus: bus,
	}
}

func (s *userService) UpdateProfile(userID string, req *dto.UserProfileUpdateRequest) (*dto.UserResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.RoleContent.AsUser == nil {
		user.RoleContent.AsUser = &model.UserRoleContent{}
	}

	if req.Bio != nil {
		user.RoleContent.AsUser.Bio = *req.Bio
	}
	if req.Location != nil {
		user.RoleContent.AsUser.Location = *req.Location
	}
	if req.Website != nil {
		user.RoleContent.AsUser.Website = *req.Website
	}

	updatedUser, err := s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return dto.FromUser(updatedUser), nil
}

func (s *userService) updateImage(userID string, imageURL string, publicID string, imageType string) (*dto.UserResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.RoleContent.AsUser == nil {
		user.RoleContent.AsUser = &model.UserRoleContent{}
	}

	var oldPublicID string
	newImage := model.Image{URL: imageURL, PublicID: publicID}

	if imageType == "avatar" {
		oldPublicID = user.RoleContent.AsUser.Avatar.PublicID
		user.RoleContent.AsUser.Avatar = newImage
	} else if imageType == "cover" {
		oldPublicID = user.RoleContent.AsUser.Cover.PublicID
		user.RoleContent.AsUser.Cover = newImage
	}

	updatedUser, err := s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	if oldPublicID != "" {
		go cloudinary.Delete(oldPublicID)
	}

	if imageType == "avatar" {
		s.eventBus.Publish(bus.UserChangeAvatarEventType{UserID: userID, NewAvatar: imageURL})
	}

	return dto.FromUser(updatedUser), nil
}

func (s *userService) UpdateAvatar(userID string, imageURL string, publicID string) (*dto.UserResponse, error) {
	return s.updateImage(userID, imageURL, publicID, "avatar")
}

func (s *userService) UpdateCover(userID string, imageURL string, publicID string) (*dto.UserResponse, error) {
	return s.updateImage(userID, imageURL, publicID, "cover")
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

	if user.Provider != model.ProviderLocal || user.Password == "" {
		return apperror.ErrLoginMethodMismatch
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

func (s *userService) GetUserByID(id string) (*dto.UserResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.ErrUserNotFound
		}
		return nil, err
	}
	return dto.FromUser(user), nil
}

func (s *userService) GetUserByUsername(username string) (*dto.UserResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.ErrUserNotFound
		}
		return nil, err
	}
	return dto.FromUser(user), nil
}

func (s *userService) GetUserByEmail(email string) (*dto.UserResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.ErrUserNotFound
		}
		return nil, err
	}
	return dto.FromUser(user), nil
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
