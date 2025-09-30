package auth

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

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
func (s *TokenService) InvalidateAllUserTokens(ctx context.Context, userID string) error {
	key := fmt.Sprintf("invalidated:user:%s", userID)
	return s.redisClient.Set(ctx, key, time.Now().Unix(), 90*24*time.Hour).Err()
}

// IsUserValid checks if a user is still valid (not invalidated)
func (s *TokenService) IsUserValid(ctx context.Context, userID string) bool {
	key := fmt.Sprintf("invalidated:user:%s", userID)
	exists, err := s.redisClient.Exists(ctx, key).Result()
	return exists == 0 && err == nil
}
