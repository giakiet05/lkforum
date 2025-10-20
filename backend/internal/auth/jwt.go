package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AuthUser represents the user parsed from a token.
type AuthUser struct {
	ID   string
	Role string
}

// Global token service instance
var TokenSvc *TokenService

// SetTokenService sets the token service for JWT operations.
func SetTokenService(service *TokenService) {
	TokenSvc = service
}

// ====== CREATE ======

// createAccessToken creates a short-lived access token.
func createAccessToken(userID, role string) (string, error) {
	jti := uuid.New().String()

	claims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"iss":  config.Cfg.JWTIssuer,
		"aud":  config.Cfg.JWTAudience,
		"iat":  time.Now().UTC().Unix(),
		"exp":  time.Now().Add(time.Minute * time.Duration(config.Cfg.TokenTTL)).Unix(),
		"jti":  jti,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Cfg.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// createRefreshToken creates a long-lived refresh token.
func createRefreshToken(userID string) (string, error) {
	jti := uuid.New().String()

	claims := jwt.MapClaims{
		"sub":  userID,
		"type": "refresh",
		"iss":  config.Cfg.JWTIssuer,
		"aud":  config.Cfg.JWTAudience,
		"iat":  time.Now().UTC().Unix(),
		"exp":  time.Now().Add(time.Hour * time.Duration(config.Cfg.RefreshTokenTTL)).Unix(),
		"jti":  jti,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Cfg.JWTSecret)) // Using the same secret for simplicity
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateToken creates a new pair of access and refresh tokens.
func GenerateToken(id string, role string) (accessToken string, refreshToken string, err error) {
	accessToken, err = createAccessToken(id, role)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = createRefreshToken(id)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ====== PARSE ======

// ParseAccessToken parses and validates an access token string.
func ParseAccessToken(tokenStr string) (AuthUser, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.Cfg.JWTSecret), nil
	})

	if err != nil {
		return AuthUser{}, apperror.ErrInvalidToken
	}
	if !token.Valid {
		return AuthUser{}, apperror.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return AuthUser{}, apperror.ErrInvalidClaims
	}

	if iss, ok := claims["iss"].(string); !ok || iss != config.Cfg.JWTIssuer {
		return AuthUser{}, apperror.ErrInvalidIssuer
	}

	if aud, ok := claims["aud"].(string); !ok || aud != config.Cfg.JWTAudience {
		return AuthUser{}, apperror.ErrInvalidAudience
	}

	userID, _ := claims["sub"].(string)
	role, _ := claims["role"].(string)

	if TokenSvc != nil {
		ctx := context.Background()
		if !TokenSvc.IsUserValid(ctx, userID) {
			return AuthUser{}, apperror.ErrTokenInvalidated
		}
	}

	return AuthUser{ID: userID, Role: role}, nil
}

// ParseRefreshToken parses and validates a refresh token string.
func ParseRefreshToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.Cfg.JWTSecret), nil
	})

	if err != nil {
		return "", apperror.ErrInvalidToken
	}
	if !token.Valid {
		return "", apperror.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", apperror.ErrInvalidClaims
	}

	if iss, ok := claims["iss"].(string); !ok || iss != config.Cfg.JWTIssuer {
		return "", apperror.ErrInvalidIssuer
	}

	if aud, ok := claims["aud"].(string); !ok || aud != config.Cfg.JWTAudience {
		return "", apperror.ErrInvalidAudience
	}

	userID, _ := claims["sub"].(string)

	if TokenSvc != nil {
		ctx := context.Background()
		if !TokenSvc.IsUserValid(ctx, userID) {
			return "", apperror.ErrTokenInvalidated
		}
	}

	return userID, nil
}

// IsOwner checks if the authenticated user is the owner of a resource.
func IsOwner(c *gin.Context, ownerID string) bool {
	authUser, exists := c.Get("authUser")
	if !exists {
		return false
	}
	user := authUser.(AuthUser)
	return user.ID == ownerID
}

// IsAdmin checks if the authenticated user has the admin role.
func IsAdmin(c *gin.Context) bool {
	authUser, exists := c.Get("authUser")
	if !exists {
		return false
	}
	return authUser.(AuthUser).Role == "admin"
}
