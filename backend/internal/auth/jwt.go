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

// AuthUser represents the authenticated user with their settings cached.
// This is populated in middleware by querying DB once per request.
type AuthUser struct {
	ID       string
	Role     string
	Settings interface{} // Will hold *model.UserSettings, using interface{} to avoid circular import
}

// SetupTokenClaims holds the claims for the short-lived token used for completing Google user setup.
type SetupTokenClaims struct {
	GoogleID string `json:"google_id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
	jwt.RegisteredClaims
}

// VerificationTokenClaims holds the claims for email verification after OTP is verified.
// This token allows the user to complete registration within 15 minutes.
type VerificationTokenClaims struct {
	Email string `json:"email"`
	Nonce string `json:"nonce"` // Prevents replay attacks
	jwt.RegisteredClaims
}

// Global token service instance
var TokenSvc *TokenService

// SetTokenService sets the token service for JWT operations.
func SetTokenService(service *TokenService) {
	TokenSvc = service
}

// ====== Login/Refresh Tokens ======

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
	return token.SignedString([]byte(config.Cfg.JWTSecret))
}

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
	return token.SignedString([]byte(config.Cfg.JWTSecret))
}

// ====== Setup Token (for Google OAuth) ======

// CreateSetupToken creates a short-lived token to complete user registration.
func CreateSetupToken(userInfo *GoogleUserInfo) (string, error) {
	claims := SetupTokenClaims{
		GoogleID: userInfo.ID,
		Email:    userInfo.Email,
		Name:     userInfo.Name,
		Picture:  userInfo.Picture,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)), // Token is valid for 15 minutes
			Issuer:    config.Cfg.JWTIssuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Use a slightly different secret for setup tokens for security.
	return token.SignedString([]byte(config.Cfg.JWTSecret + "-setup"))
}

// ParseSetupToken validates the setup token and returns the claims.
func ParseSetupToken(tokenStr string) (*SetupTokenClaims, error) {
	var claims SetupTokenClaims
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.Cfg.JWTSecret + "-setup"), nil
	})

	if err != nil {
		return nil, apperror.ErrInvalidToken
	}

	if !token.Valid {
		return nil, apperror.ErrInvalidToken
	}

	return &claims, nil
}

// ====== Verification Token (for Email Verification) ======

// CreateVerificationToken creates a short-lived token after email OTP is verified.
func CreateVerificationToken(email, nonce string) (string, error) {
	claims := VerificationTokenClaims{
		Email: email,
		Nonce: nonce,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)), // Valid for 15 minutes
			Issuer:    config.Cfg.JWTIssuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Cfg.JWTSecret + "-verification"))
}

// ParseVerificationToken validates the verification token and returns the claims.
func ParseVerificationToken(tokenStr string) (*VerificationTokenClaims, error) {
	var claims VerificationTokenClaims
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.Cfg.JWTSecret + "-verification"), nil
	})

	if err != nil {
		return nil, apperror.ErrInvalidToken
	}

	if !token.Valid {
		return nil, apperror.ErrInvalidToken
	}

	return &claims, nil
}

// ====== PARSE ======

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
	jti, _ := claims["jti"].(string)

	if TokenSvc != nil {
		// Use context with timeout for Redis operations
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		// Check if user is deleted/invalidated
		if !TokenSvc.IsUserValid(ctx, userID) {
			return AuthUser{}, apperror.ErrTokenInvalidated
		}

		// Check if this specific token is blacklisted (logout)
		if jti != "" && TokenSvc.IsTokenBlacklisted(ctx, jti) {
			return AuthUser{}, apperror.ErrTokenInvalidated
		}
	}

	// Settings will be loaded by middleware through DB query
	return AuthUser{ID: userID, Role: role, Settings: nil}, nil
}

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
	jti, _ := claims["jti"].(string)

	if TokenSvc != nil {
		// Use context with timeout for Redis operations
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		// Check if user is deleted/invalidated
		if !TokenSvc.IsUserValid(ctx, userID) {
			return "", apperror.ErrTokenInvalidated
		}

		// Check if this specific refresh token is blacklisted (logout)
		if jti != "" && TokenSvc.IsTokenBlacklisted(ctx, jti) {
			return "", apperror.ErrTokenInvalidated
		}
	}

	return userID, nil
}

// ====== HELPERS ======

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
