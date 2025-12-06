package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/auth"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/platform/bus"
	"github.com/giakiet05/lkforum/internal/platform/cloudinary"
	"github.com/giakiet05/lkforum/internal/repo"
	"github.com/giakiet05/lkforum/internal/util"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// UserService handles business logic related to user management.
type UserService interface {
	UpdateProfile(userID string, req *dto.UserProfileUpdateRequest) (*dto.UserResponse, error)
	UpdateAvatar(userID string, imageURL string, publicID string) (*dto.UserResponse, error)
	UpdateCover(userID string, imageURL string, publicID string) (*dto.UserResponse, error)
	DeleteAvatar(userID string) (*dto.UserResponse, error)
	DeleteCover(userID string) (*dto.UserResponse, error)
	DeleteUser(id string) error
	ChangePassword(userID, oldPassword, newPassword string) error

	GetUserByID(id string) (*dto.UserResponse, error)
	GetUserByUsername(username string, requesterID string) (*dto.UserResponse, error)
	GetUserByEmail(email string) (*dto.UserResponse, error)
	GetUsers(query *dto.GetUsersQuery) (*dto.PaginatedUsersResponse, error)

	GetSettings(userID string) (*dto.SettingsResponse, error)
	UpdateSettings(userID string, req *dto.UpdateSettingsRequest) (*dto.SettingsResponse, error)

	CheckUsernameAvailability(username string) (bool, error)
}

type userService struct {
	userRepo    repo.UserRepo
	eventBus    bus.EventBus
	redisClient *redis.Client
}

func NewUserService(userRepo repo.UserRepo, bus bus.EventBus, redisClient *redis.Client) UserService {
	return &userService{
		userRepo:    userRepo,
		eventBus:    bus,
		redisClient: redisClient,
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

	// Update Bio
	if req.Bio != nil {
		if *req.Bio == "" {
			user.RoleContent.AsUser.Bio = nil // Delete bio
		} else {
			user.RoleContent.AsUser.Bio = req.Bio
		}
	}

	// Update Gender
	if req.Gender != nil {
		if *req.Gender == "" {
			user.RoleContent.AsUser.Gender = nil // Delete gender
		} else {
			gender := model.Gender(*req.Gender)
			if !model.IsValidGender(gender) {
				return nil, apperror.ErrInvalidGender
			}
			user.RoleContent.AsUser.Gender = &gender
		}
	}

	// Update DateOfBirth
	if req.DateOfBirth != nil {
		if *req.DateOfBirth == "" {
			user.RoleContent.AsUser.DateOfBirth = nil // Delete date of birth
		} else {
			dob, err := time.Parse("2006-01-02", *req.DateOfBirth)
			if err != nil {
				return nil, apperror.ErrInvalidDateFormat
			}
			// Validate age >= 13
			age := time.Now().Year() - dob.Year()
			if age < 13 {
				return nil, apperror.ErrAgeTooYoung
			}
			if age > 150 {
				return nil, apperror.ErrInvalidBirthDate
			}
			user.RoleContent.AsUser.DateOfBirth = &dob
		}
	}

	// Update Location
	if req.Location != nil {
		if *req.Location == "" {
			user.RoleContent.AsUser.Location = nil // Delete location
		} else {
			province := model.VNProvince(*req.Location)
			if !model.IsValidProvince(province) {
				return nil, apperror.ErrInvalidProvince
			}
			user.RoleContent.AsUser.Location = &province
		}
	}

	// Update Interests
	if req.Interests != nil {
		if len(req.Interests) == 0 {
			user.RoleContent.AsUser.Interests = nil // Delete all interests
		} else {
			if len(req.Interests) > 10 {
				return nil, apperror.ErrTooManyInterests
			}
			interests := make([]model.Interest, len(req.Interests))
			for i, interestStr := range req.Interests {
				interest := model.Interest(interestStr)
				if !model.IsValidInterest(interest) {
					return nil, fmt.Errorf("%w: %s", apperror.ErrInvalidInterest, interestStr)
				}
				interests[i] = interest
			}
			user.RoleContent.AsUser.Interests = interests
		}
	}

	// Update Social Links
	if req.SocialLinks != nil {
		if user.RoleContent.AsUser.SocialLinks == nil {
			user.RoleContent.AsUser.SocialLinks = &model.SocialLinks{}
		}
		if req.SocialLinks.Website != nil {
			user.RoleContent.AsUser.SocialLinks.Website = req.SocialLinks.Website
		}
		if req.SocialLinks.Facebook != nil {
			user.RoleContent.AsUser.SocialLinks.Facebook = req.SocialLinks.Facebook
		}
		if req.SocialLinks.YouTube != nil {
			user.RoleContent.AsUser.SocialLinks.YouTube = req.SocialLinks.YouTube
		}
		if req.SocialLinks.GitHub != nil {
			user.RoleContent.AsUser.SocialLinks.GitHub = req.SocialLinks.GitHub
		}
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
	newImage := &model.Image{URL: imageURL, PublicID: publicID}

	if imageType == "avatar" {
		if user.RoleContent.AsUser.Avatar != nil {
			oldPublicID = user.RoleContent.AsUser.Avatar.PublicID
		}
		user.RoleContent.AsUser.Avatar = newImage
	} else if imageType == "cover" {
		if user.RoleContent.AsUser.Cover != nil {
			oldPublicID = user.RoleContent.AsUser.Cover.PublicID
		}
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

func (s *userService) DeleteAvatar(userID string) (*dto.UserResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.RoleContent.AsUser == nil {
		return dto.FromUser(user), nil
	}

	var oldPublicID string
	if user.RoleContent.AsUser.Avatar != nil {
		oldPublicID = user.RoleContent.AsUser.Avatar.PublicID
	}
	user.RoleContent.AsUser.Avatar = nil

	updatedUser, err := s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	if oldPublicID != "" {
		go cloudinary.Delete(oldPublicID)
	}

	s.eventBus.Publish(bus.UserChangeAvatarEventType{UserID: userID, NewAvatar: ""})

	return dto.FromUser(updatedUser), nil
}

func (s *userService) DeleteCover(userID string) (*dto.UserResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.RoleContent.AsUser == nil {
		return dto.FromUser(user), nil
	}

	var oldPublicID string
	if user.RoleContent.AsUser.Cover != nil {
		oldPublicID = user.RoleContent.AsUser.Cover.PublicID
	}
	user.RoleContent.AsUser.Cover = nil

	updatedUser, err := s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	if oldPublicID != "" {
		go cloudinary.Delete(oldPublicID)
	}

	return dto.FromUser(updatedUser), nil
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

func (s *userService) GetUserByUsername(username string, requesterID string) (*dto.UserResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.ErrUserNotFound
		}
		return nil, err
	}

	// Check if profile is private
	if user.Settings != nil && !user.Settings.Privacy.ShowProfile {
		// Allow user to view their own private profile
		if requesterID != "" && user.ID.Hex() == requesterID {
			return dto.FromUser(user), nil
		}

		// Return limited profile info (username, avatar, cover only) for others
		limitedProfile := &dto.UserResponse{
			ID:       user.ID.Hex(),
			Username: user.Username,
			Profile:  dto.UserProfileResponse{},
		}

		// Include avatar and cover if they exist
		if user.RoleContent.AsUser != nil {
			if user.RoleContent.AsUser.Avatar != nil {
				limitedProfile.Profile.Avatar = user.RoleContent.AsUser.Avatar
			}
			if user.RoleContent.AsUser.Cover != nil {
				limitedProfile.Profile.Cover = user.RoleContent.AsUser.Cover
			}
		}

		// Return error with limited data
		return limitedProfile, apperror.ErrForbidden
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

func (s *userService) GetUsers(query *dto.GetUsersQuery) (*dto.PaginatedUsersResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Build filter - exclude deleted and banned users for regular users
	filter := repo.Filter{
		"deleted_at": bson.M{"$exists": false},
		"is_banned":  false, // Only show non-banned users
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

func (s *userService) GetSettings(userID string) (*dto.SettingsResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.ErrUserNotFound
		}
		return nil, err
	}

	// Return user settings or default settings
	return dto.FromSettings(user.Settings), nil
}

func (s *userService) UpdateSettings(userID string, req *dto.UpdateSettingsRequest) (*dto.SettingsResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.ErrUserNotFound
		}
		return nil, err
	}

	// Initialize settings if nil
	if user.Settings == nil {
		user.Settings = model.NewDefaultSettings()
	}

	// Update Appearance settings
	if req.Appearance != nil {
		if req.Appearance.Theme != nil {
			if !model.IsValidTheme(*req.Appearance.Theme) {
				return nil, apperror.ErrBadRequest
			}
			user.Settings.Appearance.Theme = *req.Appearance.Theme
		}
		if req.Appearance.FontSize != nil {
			if !model.IsValidFontSize(*req.Appearance.FontSize) {
				return nil, apperror.ErrBadRequest
			}
			user.Settings.Appearance.FontSize = *req.Appearance.FontSize
		}
	}

	// Update Notification settings
	if req.Notifications != nil {
		if req.Notifications.InAppEnabled != nil {
			user.Settings.Notifications.InAppEnabled = *req.Notifications.InAppEnabled
		}
		if req.Notifications.EmailEnabled != nil {
			user.Settings.Notifications.EmailEnabled = *req.Notifications.EmailEnabled
		}
		if req.Notifications.NotifyOnComment != nil {
			user.Settings.Notifications.NotifyOnComment = *req.Notifications.NotifyOnComment
		}
		if req.Notifications.NotifyOnMention != nil {
			user.Settings.Notifications.NotifyOnMention = *req.Notifications.NotifyOnMention
		}
		if req.Notifications.NotifyOnUpvote != nil {
			user.Settings.Notifications.NotifyOnUpvote = *req.Notifications.NotifyOnUpvote
		}
		if req.Notifications.NotifyOnMessage != nil {
			user.Settings.Notifications.NotifyOnMessage = *req.Notifications.NotifyOnMessage
		}
	}

	// Update Privacy settings
	if req.Privacy != nil {
		if req.Privacy.ShowProfile != nil {
			user.Settings.Privacy.ShowProfile = *req.Privacy.ShowProfile
		}
		if req.Privacy.ShowEmail != nil {
			user.Settings.Privacy.ShowEmail = *req.Privacy.ShowEmail
		}
		if req.Privacy.ShowPostHistory != nil {
			user.Settings.Privacy.ShowPostHistory = *req.Privacy.ShowPostHistory
		}
		if req.Privacy.AllowDirectMessages != nil {
			user.Settings.Privacy.AllowDirectMessages = *req.Privacy.AllowDirectMessages
		}
		if req.Privacy.AllowMentions != nil {
			user.Settings.Privacy.AllowMentions = *req.Privacy.AllowMentions
		}
	}

	// Update Content settings
	if req.Content != nil {
		if req.Content.AllowNSFW != nil {
			user.Settings.Content.AllowNSFW = *req.Content.AllowNSFW
		}
	}

	// Save updated user
	updatedUser, err := s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return dto.FromSettings(updatedUser.Settings), nil
}

func (s *userService) CheckUsernameAvailability(username string) (bool, error) {
	// Try cache first
	if s.redisClient != nil {
		ctx, cancel := util.NewDefaultRedisContext()
		defer cancel()

		cacheKey := fmt.Sprintf("username_exists:%s", username)
		cached, err := s.redisClient.Get(ctx, cacheKey).Result()
		if err == nil {
			// Cache hit - "false" means available, "true" means taken
			return cached == "false", nil
		}
	}

	// Cache miss - query database
	dbCtx, cancel := util.NewDefaultDBContext()
	defer cancel()

	_, err := s.userRepo.GetByUsername(dbCtx, username)
	exists := !errors.Is(err, mongo.ErrNoDocuments)

	// Cache the result (5 minutes TTL)
	if s.redisClient != nil {
		ctx, cancel := util.NewDefaultRedisContext()
		defer cancel()

		cacheKey := fmt.Sprintf("username_exists:%s", username)
		value := "false"
		if exists {
			value = "true"
		}
		// Ignore cache write errors, not critical
		_ = s.redisClient.Set(ctx, cacheKey, value, 5*time.Minute).Err()
	}

	return !exists, nil
}
