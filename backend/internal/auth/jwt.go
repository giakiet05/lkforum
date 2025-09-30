package auth

import (
	"context"
	"fmt"
	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"os"
	"time"
)

// AuthUser đại diện cho user sau khi parse token
type AuthUser struct {
	ID   string
	Role string
}

// Global token service instance
var TokenSvc *TokenService

// SetTokenService sets the token service for JWT operations
func SetTokenService(service *TokenService) {
	TokenSvc = service
}

var (
	accessSecret  = []byte(os.Getenv("ACCESS_TOKEN_SECRET"))
	refreshSecret = []byte(os.Getenv("REFRESH_TOKEN_SECRET"))
	issuer        = os.Getenv("JWT_ISS")
	audience      = os.Getenv("JWT_AUD")
)

// ====== CREATE ======

// Tạo access token ngắn hạn
func createAccessToken(userID, role string) (string, error) {
	expMinutes := config.GetEnvIntWithDefault("ACCESS_TOKEN_EXP_MIN", 15)
	jti := uuid.New().String()

	claims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"iss":  issuer,
		"aud":  audience,
		"iat":  time.Now().UTC().Unix(),
		"exp":  time.Now().Add(time.Minute * time.Duration(expMinutes)).Unix(),
		"jti":  jti,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(accessSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Tạo refresh token dài hạn
func createRefreshToken(userID string) (string, error) {
	expDays := config.GetEnvIntWithDefault("REFRESH_TOKEN_EXP_DAYS", 7)
	jti := uuid.New().String()

	claims := jwt.MapClaims{
		"sub":  userID,
		"type": "refresh",
		"iss":  issuer,
		"aud":  audience,
		"iat":  time.Now().UTC().Unix(),
		"exp":  time.Now().Add(24 * time.Hour * time.Duration(expDays)).Unix(),
		"jti":  jti,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(refreshSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

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

// Parse + validate access token
func ParseAccessToken(tokenStr string) (AuthUser, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return accessSecret, nil
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

	// Verify issuer and audience explicitly
	if iss, ok := claims["iss"].(string); !ok || iss != issuer {
		return AuthUser{}, apperror.ErrInvalidIssuer
	}

	if aud, ok := claims["aud"].(string); !ok || aud != audience {
		return AuthUser{}, apperror.ErrInvalidAudience
	}

	userID, _ := claims["sub"].(string)
	role, _ := claims["role"].(string)

	// Check if token has been invalidated (if token service is available)
	if TokenSvc != nil {
		ctx := context.Background()
		if !TokenSvc.IsUserValid(ctx, userID) {
			return AuthUser{}, apperror.ErrTokenInvalidated
		}
	}

	return AuthUser{ID: userID, Role: role}, nil
}

// Parse + validate refresh token
func ParseRefreshToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return refreshSecret, nil
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

	// Verify issuer and audience explicitly
	if iss, ok := claims["iss"].(string); !ok || iss != issuer {
		return "", apperror.ErrInvalidIssuer
	}

	if aud, ok := claims["aud"].(string); !ok || aud != audience {
		return "", apperror.ErrInvalidAudience
	}

	userID, _ := claims["sub"].(string)

	// Check if token has been invalidated (if token service is available)
	if TokenSvc != nil {
		ctx := context.Background()
		if !TokenSvc.IsUserValid(ctx, userID) {
			return "", apperror.ErrTokenInvalidated
		}
	}

	return userID, nil
}

func IsOwner(c *gin.Context, ownerID string) bool {
	authUser, exists := c.Get("authUser")
	if !exists {
		return false
	}
	user := authUser.(AuthUser)
	return user.ID == ownerID
}

func IsAdmin(c *gin.Context) bool {
	authUser, exists := c.Get("authUser")
	if !exists {
		return false
	}
	return authUser.(AuthUser).Role == "admin"
}
