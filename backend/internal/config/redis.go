package config

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	// Redis key patterns
	RedisInvalidatedUserKey  = "invalidated:user:%s"  // For delete user - invalidate all tokens
	RedisBlacklistedTokenKey = "blacklisted:token:%s" // For logout - invalidate specific token by JTI
	RedisActiveUsersKey      = "channel:%s:active_users"
	RedisMembersCountKey     = "community:%s:member_count"
)

// NewRedisClient creates and returns a new Redis client using the global AppConfig.
func NewRedisClient() *redis.Client {
	// Create the client with configuration from the global Cfg variable.
	client := redis.NewClient(&redis.Options{
		Addr:     Cfg.Redis.Addr,
		Password: Cfg.Redis.Password,
		DB:       Cfg.Redis.DB,
	})

	// Create a context with a timeout to test the connection.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Ping the Redis server to ensure the connection is alive.
	if err := client.Ping(ctx).Err(); err != nil {
		// Log a warning instead of a fatal error.
		// This allows the application to continue running even if Redis is unavailable.
		// Features that depend on Redis (like token invalidation) will be gracefully disabled.
		log.Printf("WARNING: Could not connect to Redis at %s. Features depending on Redis may be disabled. Error: %v", Cfg.Redis.Addr, err)
	} else {
		log.Println("Successfully connected to Redis.")
	}

	return client
}
