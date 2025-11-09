package middleware

import (
	"net/http"
	"strings"

	"github.com/giakiet05/lkforum/internal/auth"
	"github.com/giakiet05/lkforum/internal/repo"
	"github.com/giakiet05/lkforum/internal/util"
	"github.com/gin-gonic/gin"
)

// userRepo is injected at startup to load user settings in middleware
var userRepo repo.UserRepo

// SetUserRepo injects the user repository for middleware to use
func SetUserRepo(repo repo.UserRepo) {
	userRepo = repo
}

// RequireAuth parse access token và nhét AuthUser vào context
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]
		user, err := auth.ParseAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Load user settings from DB once per request
		if userRepo != nil {
			ctx, cancel := util.NewDefaultDBContext()
			defer cancel()

			dbUser, err := userRepo.GetByID(ctx, user.ID)
			if err == nil && dbUser.Settings != nil {
				user.Settings = dbUser.Settings
			}
		}

		// Nhét user vào context with settings cached
		c.Set("authUser", user)
		c.Next()
	}
}

// LoadUserIfAuthenticated tries to parse the access token and load the user into the context.
// Unlike RequireAuth, it does not fail if the token is missing or invalid.
func LoadUserIfAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		token := parts[1]
		user, err := auth.ParseAccessToken(token)
		if err != nil {
			c.Next()
			return
		}

		// Load user settings from DB once per request
		if userRepo != nil {
			ctx, cancel := util.NewDefaultDBContext()
			defer cancel()

			dbUser, err := userRepo.GetByID(ctx, user.ID)
			if err == nil && dbUser.Settings != nil {
				user.Settings = dbUser.Settings
			}
		}

		// If token is valid, set the user in the context with settings cached
		c.Set("authUser", user)
		c.Next()
	}
}

// RequireAdmin check role admin
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exists := c.Get("authUser")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
			c.Abort()
			return
		}

		user, ok := val.(auth.AuthUser)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid auth context"})
			c.Abort()
			return
		}

		if user.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
