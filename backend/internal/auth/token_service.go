package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/giakiet05/lkforum/internal/config"
	"github.com/redis/go-redis/v9"
)

// TokenServiceInterface defines the contract for token operations
type TokenServiceInterface interface {
	InvalidateAllUserTokens(ctx context.Context, userID string) error
	IsUserValid(ctx context.Context, userID string) bool
	InvalidateToken(ctx context.Context, jti string, ttl time.Duration) error
	IsTokenBlacklisted(ctx context.Context, jti string) bool
}

// TokenService handles token operations including invalidation
type TokenService struct {
	redisClient *redis.Client
}

// NewTokenService creates a new token service with Redis client
func NewTokenService(redisClient *redis.Client) *TokenService {
	return &TokenService{
		redisClient: redisClient,
	}
}

// InvalidateAllUserTokens marks a user as deleted in Redis
// Used for: Delete user account
func (s *TokenService) InvalidateAllUserTokens(ctx context.Context, userID string) error {
	key := fmt.Sprintf(config.RedisInvalidatedUserKey, userID)
	return s.redisClient.Set(ctx, key, time.Now().Unix(), 90*24*time.Hour).Err()
}

// IsUserValid checks if a user is still valid (not invalidated)
func (s *TokenService) IsUserValid(ctx context.Context, userID string) bool {
	key := fmt.Sprintf(config.RedisInvalidatedUserKey, userID)
	exists, err := s.redisClient.Exists(ctx, key).Result()
	return exists == 0 && err == nil
}

// InvalidateToken blacklists a specific token by its JTI
// Used for: Logout (invalidate only current session)
func (s *TokenService) InvalidateToken(ctx context.Context, jti string, ttl time.Duration) error {
	key := fmt.Sprintf(config.RedisBlacklistedTokenKey, jti)
	return s.redisClient.Set(ctx, key, time.Now().Unix(), ttl).Err()
}

// IsTokenBlacklisted checks if a token is blacklisted by JTI
func (s *TokenService) IsTokenBlacklisted(ctx context.Context, jti string) bool {
	key := fmt.Sprintf(config.RedisBlacklistedTokenKey, jti)
	exists, err := s.redisClient.Exists(ctx, key).Result()
	return exists > 0 && err == nil
}
