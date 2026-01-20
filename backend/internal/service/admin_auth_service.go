package service

import (
	"errors"
	"time"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/auth"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/repo"
	"github.com/giakiet05/lkforum/internal/util"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AdminAuthService interface {
	AdminLogin(identifier, password string) (*model.User, string, string, error)
	AdminLogout(accessToken, refreshToken string) error
	AdminRefreshToken(refreshToken string) (string, string, error)
}

type adminAuthService struct {
	userRepo     repo.UserRepo
	redisClient  *redis.Client
	tokenService auth.TokenServiceInterface
}

func NewAdminAuthService(userRepo repo.UserRepo, redisClient *redis.Client, tokenService auth.TokenServiceInterface) AdminAuthService {
	return &adminAuthService{
		userRepo:     userRepo,
		redisClient:  redisClient,
		tokenService: tokenService,
	}
}

// AdminLogin authenticates an admin user and returns tokens
func (s *adminAuthService) AdminLogin(identifier, password string) (*model.User, string, string, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Get user by email or username
	user, err := s.userRepo.GetByEmail(ctx, identifier)
	if errors.Is(err, mongo.ErrNoDocuments) {
		user, err = s.userRepo.GetByUsername(ctx, identifier)
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, "", "", apperror.ErrInvalidCredentials
		}
	}
	if err != nil {
		return nil, "", "", err
	}

	// Check if user is admin
	if user.Role != model.AdminRole {
		return nil, "", "", apperror.ErrAdminAccessRequired
	}

	// Check if user is banned
	if user.IsBanned {
		return nil, "", "", apperror.ErrUserInactive
	}

	// Check if user is deleted
	if user.DeletedAt != nil {
		return nil, "", "", apperror.ErrUserNotFound
	}

	// Verify password (only for local provider)
	if user.Provider == model.ProviderLocal {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			return nil, "", "", apperror.ErrInvalidCredentials
		}
	}

	// Generate tokens
	accessToken, refreshToken, err := auth.GenerateToken(user.ID.Hex(), string(user.Role))
	if err != nil {
		return nil, "", "", apperror.ErrInternal
	}

	return user, accessToken, refreshToken, nil
}

// AdminLogout invalidates admin tokens
func (s *adminAuthService) AdminLogout(accessToken, refreshToken string) error {
	if s.tokenService == nil {
		return nil // Token invalidation not available
	}

	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Extract JTI from access token
	accessJTI, err := extractJTIFromToken(accessToken)
	if err == nil {
		accessTTL := time.Minute * time.Duration(15) // Default access token TTL
		_ = s.tokenService.InvalidateToken(ctx, accessJTI, accessTTL)
	}

	// Extract JTI from refresh token
	refreshJTI, err := extractJTIFromToken(refreshToken)
	if err == nil {
		refreshTTL := time.Hour * time.Duration(720) // Default refresh token TTL
		_ = s.tokenService.InvalidateToken(ctx, refreshJTI, refreshTTL)
	}

	return nil
}

// Helper function to extract JTI from token
func extractJTIFromToken(tokenStr string) (string, error) {
	token, _, err := jwt.NewParser().ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if jti, ok := claims["jti"].(string); ok {
			return jti, nil
		}
	}

	return "", apperror.ErrInvalidToken
}

// AdminRefreshToken refreshes an admin's access token using a refresh token
func (s *adminAuthService) AdminRefreshToken(refreshToken string) (string, string, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Parse and validate refresh token
	userID, err := auth.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", "", apperror.ErrInvalidToken
	}

	// Get user from database
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", "", apperror.ErrUserNotFound
		}
		return "", "", err
	}

	// Check if user is still admin
	if user.Role != model.AdminRole {
		return "", "", apperror.ErrAdminAccessRequired
	}

	// Check if user is banned
	if user.IsBanned {
		return "", "", apperror.ErrUserInactive
	}

	// Check if user is deleted
	if user.DeletedAt != nil {
		return "", "", apperror.ErrUserNotFound
	}

	// Check if refresh token is blacklisted
	if s.tokenService != nil {
		jti, err := extractJTIFromToken(refreshToken)
		if err == nil {
			if s.tokenService.IsTokenBlacklisted(ctx, jti) {
				return "", "", apperror.ErrInvalidToken
			}
		}
	}

	// Generate new tokens
	newAccessToken, newRefreshToken, err := auth.GenerateToken(userID, string(user.Role))
	if err != nil {
		return "", "", apperror.ErrInternal
	}

	return newAccessToken, newRefreshToken, nil
}
