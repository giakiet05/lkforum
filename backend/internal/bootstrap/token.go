package bootstrap

import (
	"github.com/giakiet05/lkforum/internal/auth"
	"github.com/redis/go-redis/v9"
)

// InitializeTokenService sets up the token service for JWT authentication using a provided Redis client
func InitializeTokenService(redisClient *redis.Client) error {
	tokenService := auth.NewTokenService(redisClient)
	auth.SetTokenService(tokenService)

	return nil
}
