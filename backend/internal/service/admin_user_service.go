package service

import (
	"errors"
	"time"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/repo"
	"github.com/giakiet05/lkforum/internal/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AdminUserService interface {
	// User management
	GetUsersAdmin(query *dto.GetUsersAdminQuery) (*dto.PaginatedUsersResponse, error)
	BanUser(userID string, req *dto.BanUserRequest) error
	UnbanUser(userID string) error
	SoftDeleteUser(userID string) error
	RestoreUser(userID string) error
}

type adminUserService struct {
	userRepo repo.UserRepo
}

func NewAdminUserService(userRepo repo.UserRepo) AdminUserService {
	return &adminUserService{
		userRepo: userRepo,
	}
}

func (s *adminUserService) GetUsersAdmin(query *dto.GetUsersAdminQuery) (*dto.PaginatedUsersResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Build filter based on status
	filter := repo.Filter{}

	switch query.Status {
	case "active":
		filter["is_banned"] = false
		filter["deleted_at"] = bson.M{"$exists": false}
	case "banned":
		filter["is_banned"] = true
	case "deleted":
		filter["deleted_at"] = bson.M{"$exists": true}
	case "all":
		// No filter - get all users
	default:
		// Default: active users only
		filter["is_banned"] = false
		filter["deleted_at"] = bson.M{"$exists": false}
	}

	// Add username search if provided
	if query.Username != "" {
		filter["username"] = bson.M{"$regex": primitive.Regex{Pattern: query.Username, Options: "i"}}
	}

	// Pagination
	page := query.Page
	if page < 1 {
		page = 1
	}
	pageSize := query.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	findOptions := &repo.FindOptions{
		Skip:  int64((page - 1) * pageSize),
		Limit: int64(pageSize),
		Sort:  map[string]int{"created_at": -1},
	}

	users, total, err := s.userRepo.Find(ctx, filter, findOptions)
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

func (s *adminUserService) BanUser(userID string, req *dto.BanUserRequest) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return apperror.ErrUserNotFound
		}
		return err
	}

	// Cannot ban admin
	if user.Role == model.AdminRole {
		return apperror.NewError(nil, apperror.ErrForbidden.Code, "cannot ban admin user")
	}

	// Update ban fields
	user.IsBanned = true
	user.BanUntil = req.BanUntil // null = permanent
	user.BanReason = &req.Reason

	_, err = s.userRepo.Update(ctx, user)
	return err
}

func (s *adminUserService) UnbanUser(userID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return apperror.ErrUserNotFound
		}
		return err
	}

	// Unban user
	user.IsBanned = false
	user.BanUntil = nil
	user.BanReason = nil

	_, err = s.userRepo.Update(ctx, user)
	return err
}

func (s *adminUserService) SoftDeleteUser(userID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return apperror.ErrUserNotFound
		}
		return err
	}

	// Cannot delete admin
	if user.Role == model.AdminRole {
		return apperror.NewError(nil, apperror.ErrForbidden.Code, "cannot delete admin user")
	}

	// Soft delete
	now := time.Now()
	user.DeletedAt = &now

	_, err = s.userRepo.Update(ctx, user)
	return err
}

func (s *adminUserService) RestoreUser(userID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return apperror.ErrUserNotFound
		}
		return err
	}

	// Check if user is deleted
	if user.DeletedAt == nil {
		return apperror.NewError(nil, apperror.ErrBadRequest.Code, "user is not deleted")
	}

	// Restore user
	user.DeletedAt = nil

	_, err = s.userRepo.Update(ctx, user)
	return err
}
