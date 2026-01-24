package service

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"time"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/auth"
	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/platform/email"
	"github.com/giakiet05/lkforum/internal/repo"
	"github.com/giakiet05/lkforum/internal/util"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// GoogleAuthStatus defines the result status of a Google callback processing.
const (
	StatusLoginSuccess  = "LOGIN_SUCCESS"
	StatusSetupRequired = "SETUP_REQUIRED"
)

// GoogleAuthResult is the result of processing a Google OAuth callback.
type GoogleAuthResult struct {
	Status       string
	User         *model.User
	AccessToken  string
	RefreshToken string
	SetupToken   string
}

type AuthService interface {
	// Local Auth - New Flow (Verify Email First)
	SendEmailVerification(email string) error
	VerifyEmailCode(email, otp string) (string, error) // Returns verification_token
	CompleteRegistration(verificationToken, username, password string) (*model.User, string, string, error)
	ResendOTP(email string) error
	Login(identifier, password string) (*model.User, string, string, error)
	RefreshToken(refreshToken string) (string, string, error)
	Logout(accessToken, refreshToken string) error

	// Forgot Password Flow
	ForgotPassword(email string) error
	VerifyResetPasswordOTP(email, otp string) (string, error) // Returns reset_token
	ResetPassword(resetToken, newPassword string) error

	// Google OAuth
	ProcessGoogleCallback(code string) (*GoogleAuthResult, error)
	CompleteGoogleSetup(setupToken, username string) (*model.User, string, string, error)
}

type authService struct {
	userRepo              repo.UserRepo
	emailVerificationRepo repo.EmailVerificationRepo
	passwordResetRepo     repo.PasswordResetRepo
	emailSender           email.Sender
	redisClient           *redis.Client
	tokenService          auth.TokenServiceInterface
}

func NewAuthService(userRepo repo.UserRepo, emailVerificationRepo repo.EmailVerificationRepo, passwordResetRepo repo.PasswordResetRepo, emailSender email.Sender, redisClient *redis.Client, tokenService auth.TokenServiceInterface) AuthService {
	return &authService{
		userRepo:              userRepo,
		emailVerificationRepo: emailVerificationRepo,
		passwordResetRepo:     passwordResetRepo,
		emailSender:           emailSender,
		redisClient:           redisClient,
		tokenService:          tokenService,
	}
}

// --- Local Authentication - New Flow (Verify Email First) ---

// SendEmailVerification initiates the registration process by sending OTP to email
func (s *authService) SendEmailVerification(email string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Check if email already registered
	if _, err := s.userRepo.GetByEmail(ctx, email); !errors.Is(err, mongo.ErrNoDocuments) {
		return apperror.ErrEmailExists
	}

	// Check if there's an existing verification (delete it first to allow resend)
	if existing, err := s.emailVerificationRepo.GetByEmail(ctx, email); err == nil {
		_ = s.emailVerificationRepo.Delete(ctx, existing.Email)
	}

	// Create new verification record
	otp := generateOTP()
	nonce := generateNonce()
	otpExpiresAt := time.Now().Add(time.Duration(config.Cfg.OTPExpirationMinutes) * time.Minute)

	verification := &model.EmailVerification{
		Email:        email,
		OTP:          otp,
		OTPExpiresAt: otpExpiresAt,
		IsVerified:   false,
		Nonce:        nonce,
		CreatedAt:    time.Now(),
	}

	_, err := s.emailVerificationRepo.Create(ctx, verification)
	if err != nil {
		return err
	}

	// Send OTP email
	if s.emailSender != nil {
		go func() {
			if err := s.emailSender.SendVerificationEmail(email, otp); err != nil {
				fmt.Printf("CRITICAL: Failed to send verification email to %s: %v\n", email, err)
			}
		}()
	}

	return nil
}

// VerifyEmailCode verifies the OTP and returns a verification_token
func (s *authService) VerifyEmailCode(email, otp string) (string, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	verification, err := s.emailVerificationRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", apperror.ErrInvalidOTP
		}
		return "", err
	}

	if verification.IsVerified {
		return "", apperror.ErrEmailAlreadyVerified
	}

	if verification.OTP != otp {
		return "", apperror.ErrInvalidOTP
	}

	if verification.OTPExpiresAt.Before(time.Now()) {
		return "", apperror.ErrOTPExpired
	}

	// Mark as verified
	verification.IsVerified = true
	_, err = s.emailVerificationRepo.Update(ctx, verification)
	if err != nil {
		return "", err
	}

	// Generate verification token (valid 15 min)
	verificationToken, err := auth.CreateVerificationToken(email, verification.Nonce)
	if err != nil {
		return "", err
	}

	return verificationToken, nil
}

// CompleteRegistration creates the user account after email verification
func (s *authService) CompleteRegistration(verificationToken, username, password string) (*model.User, string, string, error) {
	// Validate username length (3-20 characters)
	if len(username) < 3 || len(username) > 20 {
		return nil, "", "", apperror.ErrInvalidUsername
	}

	// Validate password length (minimum 8 characters)
	if len(password) < 8 {
		return nil, "", "", apperror.ErrPasswordTooShort
	}

	// Validate password strength (must contain uppercase, lowercase, number, special char)
	if !isStrongPassword(password) {
		return nil, "", "", apperror.ErrPasswordTooWeak
	}

	// Parse verification token
	claims, err := auth.ParseVerificationToken(verificationToken)
	if err != nil {
		return nil, "", "", err
	}

	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Verify the nonce matches (prevent replay)
	verification, err := s.emailVerificationRepo.GetByEmail(ctx, claims.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, "", "", apperror.ErrInvalidToken
		}
		return nil, "", "", err
	}

	if !verification.IsVerified {
		return nil, "", "", apperror.ErrEmailNotVerified
	}

	if verification.Nonce != claims.Nonce {
		return nil, "", "", apperror.ErrInvalidToken
	}

	// Check username availability
	if _, err := s.userRepo.GetByUsername(ctx, username); !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, "", "", apperror.ErrUsernameExists
	}

	// Double-check email not taken (race condition prevention)
	if _, err := s.userRepo.GetByEmail(ctx, claims.Email); !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, "", "", apperror.ErrEmailExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", "", err
	}

	// Create user (already verified)
	user := &model.User{
		Username:   username,
		Email:      claims.Email,
		Password:   string(hashedPassword),
		Provider:   model.ProviderLocal,
		Role:       model.UserRole,
		Settings:   model.NewDefaultSettings(),
		IsVerified: true,  // Always true since we verified email first
		IsBanned:   false, // Initialize ban status
		CreatedAt:  time.Now(),
		RoleContent: model.RoleContent{
			AsUser: &model.UserRoleContent{
				Avatar: &model.Image{}, // Initialize with an empty Image struct
				Stats: &model.ActivityStats{
					PostCount:    0,
					CommentCount: 0,
					TotalUpvotes: 0,
					JoinedAt:     time.Now(),
					LastActiveAt: time.Now(),
				},
			},
		},
	}

	createdUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, "", "", err
	}

	// Delete verification record (cleanup)
	_ = s.emailVerificationRepo.Delete(ctx, claims.Email)

	// Invalidate username cache
	s.invalidateUsernameCache(username)

	// Generate access & refresh tokens
	accessToken, refreshToken, err := auth.GenerateToken(createdUser.ID.Hex(), string(createdUser.Role))
	if err != nil {
		return nil, "", "", err
	}

	return createdUser, accessToken, refreshToken, nil
}

// ResendOTP resends OTP for email verification
func (s *authService) ResendOTP(email string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	verification, err := s.emailVerificationRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return apperror.ErrUserNotFound
		}
		return err
	}

	if verification.IsVerified {
		return apperror.ErrEmailAlreadyVerified
	}

	// Generate new OTP
	otp := generateOTP()
	otpExpiresAt := time.Now().Add(time.Duration(config.Cfg.OTPExpirationMinutes) * time.Minute)
	verification.OTP = otp
	verification.OTPExpiresAt = otpExpiresAt

	_, err = s.emailVerificationRepo.Update(ctx, verification)
	if err != nil {
		return err
	}

	// Send email
	if s.emailSender != nil {
		go func() {
			if err := s.emailSender.SendVerificationEmail(email, otp); err != nil {
				fmt.Printf("CRITICAL: Failed to resend verification email to %s: %v\n", email, err)
			}
		}()
	}

	return nil
}

func (s *authService) Login(identifier, password string) (*model.User, string, string, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()
	var user *model.User
	var err error

	if isEmail(identifier) {
		user, err = s.userRepo.GetByEmail(ctx, identifier)
	} else {
		user, err = s.userRepo.GetByUsername(ctx, identifier)
	}

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, "", "", apperror.ErrInvalidCredentials
		}
		return nil, "", "", err
	}

	if user.Provider != model.ProviderLocal {
		return nil, "", "", apperror.ErrLoginMethodMismatch
	}

	if user.Password == "" || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return nil, "", "", apperror.ErrInvalidCredentials
	}

	if !user.IsVerified {
		return nil, "", "", apperror.ErrEmailNotVerified
	}

	// Check if user is banned
	if user.IsBanned {
		// Check if ban has expired
		if user.BanUntil != nil && time.Now().After(*user.BanUntil) {
			// Ban expired, unban user
			user.IsBanned = false
			user.BanUntil = nil
			user.BanReason = nil
			s.userRepo.Update(ctx, user)
		} else {
			// Still banned
			return nil, "", "", apperror.ErrUserInactive
		}
	}

	accessToken, refreshToken, err := auth.GenerateToken(user.ID.Hex(), string(user.Role))
	if err != nil {
		return nil, "", "", err
	}
	return user, accessToken, refreshToken, nil
}

func (s *authService) RefreshToken(refreshToken string) (string, string, error) {
	userID, err := auth.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", "", apperror.ErrUserNotFound
		}
		return "", "", err
	}

	// Check if user is banned
	if user.IsBanned {
		// Check if ban has expired
		if user.BanUntil != nil && time.Now().After(*user.BanUntil) {
			// Ban expired, unban user
			user.IsBanned = false
			user.BanUntil = nil
			user.BanReason = nil
			s.userRepo.Update(ctx, user)
		} else {
			// Still banned
			return "", "", apperror.ErrUserInactive
		}
	}

	accessToken, newRefreshToken, err := auth.GenerateToken(user.ID.Hex(), string(user.Role))
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

func (s *authService) Logout(accessToken, refreshToken string) error {
	if s.tokenService == nil {
		return apperror.ErrInternal
	}

	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Parse access token to get JTI
	accessJTI, err := extractJTI(accessToken)
	if err != nil {
		return apperror.ErrInvalidToken
	}

	// Parse refresh token to get JTI
	refreshJTI, err := extractJTI(refreshToken)
	if err != nil {
		return apperror.ErrInvalidToken
	}

	// Blacklist access token
	accessTTL := time.Minute * time.Duration(config.Cfg.TokenTTL)
	if err := s.tokenService.InvalidateToken(ctx, accessJTI, accessTTL); err != nil {
		return err
	}

	// Blacklist refresh token
	refreshTTL := time.Hour * time.Duration(config.Cfg.RefreshTokenTTL)
	if err := s.tokenService.InvalidateToken(ctx, refreshJTI, refreshTTL); err != nil {
		return err
	}

	return nil
}

// --- Forgot Password Flow ---

// ForgotPassword initiates password reset by sending OTP to email
func (s *authService) ForgotPassword(email string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Check if email is registered and user uses local auth
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return apperror.ErrEmailNotRegistered
		}
		return err
	}

	// Only allow password reset for local auth users
	if user.Provider != model.ProviderLocal {
		return apperror.ErrLoginMethodMismatch
	}

	// Check if there's an existing reset request (delete it first to allow resend)
	if existing, err := s.passwordResetRepo.GetByEmail(ctx, email); err == nil {
		_ = s.passwordResetRepo.Delete(ctx, existing.Email)
	}

	// Create new password reset record
	otp := generateOTP()
	nonce := generateNonce()
	otpExpiresAt := time.Now().Add(time.Duration(config.Cfg.OTPExpirationMinutes) * time.Minute)

	reset := &model.PasswordReset{
		Email:        email,
		OTP:          otp,
		OTPExpiresAt: otpExpiresAt,
		IsVerified:   false,
		Nonce:        nonce,
		CreatedAt:    time.Now(),
	}

	_, err = s.passwordResetRepo.Create(ctx, reset)
	if err != nil {
		return err
	}

	// Send OTP email
	if s.emailSender != nil {
		go func() {
			if err := s.emailSender.SendPasswordResetEmail(email, otp); err != nil {
				fmt.Printf("CRITICAL: Failed to send password reset email to %s: %v\n", email, err)
			}
		}()
	}

	return nil
}

// VerifyResetPasswordOTP verifies the OTP and returns a reset_token
func (s *authService) VerifyResetPasswordOTP(email, otp string) (string, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	reset, err := s.passwordResetRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", apperror.ErrInvalidOTP
		}
		return "", err
	}

	if reset.IsVerified {
		return "", apperror.ErrEmailAlreadyVerified
	}

	if reset.OTP != otp {
		return "", apperror.ErrInvalidOTP
	}

	if reset.OTPExpiresAt.Before(time.Now()) {
		return "", apperror.ErrOTPExpired
	}

	// Mark as verified
	reset.IsVerified = true
	_, err = s.passwordResetRepo.Update(ctx, reset)
	if err != nil {
		return "", err
	}

	// Generate reset token (valid 15 min) - reuse verification token logic
	resetToken, err := auth.CreateVerificationToken(email, reset.Nonce)
	if err != nil {
		return "", err
	}

	return resetToken, nil
}

// ResetPassword changes the user's password using reset token
func (s *authService) ResetPassword(resetToken, newPassword string) error {
	// Validate password length (minimum 8 characters)
	if len(newPassword) < 8 {
		return apperror.ErrPasswordTooShort
	}

	// Validate password strength
	if !isStrongPassword(newPassword) {
		return apperror.ErrPasswordTooWeak
	}

	// Parse reset token
	claims, err := auth.ParseVerificationToken(resetToken)
	if err != nil {
		return err
	}

	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Verify the nonce matches (prevent replay)
	reset, err := s.passwordResetRepo.GetByEmail(ctx, claims.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return apperror.ErrInvalidToken
		}
		return err
	}

	if !reset.IsVerified {
		return apperror.ErrInvalidOTP
	}

	if reset.Nonce != claims.Nonce {
		return apperror.ErrInvalidToken
	}

	// Get user
	user, err := s.userRepo.GetByEmail(ctx, claims.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return apperror.ErrUserNotFound
		}
		return err
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update user password
	user.Password = string(hashedPassword)
	_, err = s.userRepo.Update(ctx, user)
	if err != nil {
		return err
	}

	// Delete reset record (cleanup)
	_ = s.passwordResetRepo.Delete(ctx, claims.Email)

	return nil
}

// --- Google OAuth ---

func (s *authService) ProcessGoogleCallback(code string) (*GoogleAuthResult, error) {
	userInfo, err := auth.GetGoogleUserInfo(code)
	if err != nil {
		return nil, err
	}

	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	user, err := s.userRepo.GetByEmail(ctx, userInfo.Email)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			setupToken, err := auth.CreateSetupToken(userInfo)
			if err != nil {
				return nil, err
			}
			return &GoogleAuthResult{Status: StatusSetupRequired, SetupToken: setupToken}, nil
		}
		return nil, err
	}

	if user.Provider != model.ProviderGoogle {
		return nil, apperror.ErrLoginMethodMismatch
	}

	// Check if user is banned
	if user.IsBanned {
		// Check if ban has expired
		if user.BanUntil != nil && time.Now().After(*user.BanUntil) {
			// Ban expired, unban user
			user.IsBanned = false
			user.BanUntil = nil
			user.BanReason = nil
			s.userRepo.Update(ctx, user)
		} else {
			// Still banned
			return nil, apperror.ErrUserInactive
		}
	}

	accessToken, refreshToken, err := auth.GenerateToken(user.ID.Hex(), string(user.Role))
	if err != nil {
		return nil, err
	}

	return &GoogleAuthResult{
		Status:       StatusLoginSuccess,
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) CompleteGoogleSetup(setupToken, username string) (*model.User, string, string, error) {
	// Validate username length (3-20 characters)
	if len(username) < 3 || len(username) > 20 {
		return nil, "", "", apperror.ErrInvalidUsername
	}

	claims, err := auth.ParseSetupToken(setupToken)
	if err != nil {
		return nil, "", "", err
	}

	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	if _, err := s.userRepo.GetByUsername(ctx, username); !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, "", "", apperror.ErrUsernameExists
	}

	if _, err := s.userRepo.GetByEmail(ctx, claims.Email); !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, "", "", apperror.ErrEmailExists
	}

	newUser := &model.User{
		Username:   username,
		Email:      claims.Email,
		Provider:   model.ProviderGoogle,
		ProviderID: claims.GoogleID,
		Role:       model.UserRole,
		Settings:   model.NewDefaultSettings(),
		IsVerified: true,
		IsBanned:   false, // Initialize ban status
		CreatedAt:  time.Now(),
		RoleContent: model.RoleContent{
			AsUser: &model.UserRoleContent{
				Avatar: &model.Image{URL: claims.Picture},
				Stats: &model.ActivityStats{
					PostCount:    0,
					CommentCount: 0,
					TotalUpvotes: 0,
					JoinedAt:     time.Now(),
					LastActiveAt: time.Now(),
				},
			},
		},
	}

	createdUser, err := s.userRepo.Create(ctx, newUser)
	if err != nil {
		return nil, "", "", err
	}

	// Invalidate username cache
	s.invalidateUsernameCache(username)

	accessToken, refreshToken, err := auth.GenerateToken(createdUser.ID.Hex(), string(createdUser.Role))
	if err != nil {
		return nil, "", "", err
	}

	return createdUser, accessToken, refreshToken, nil
}

// --- Helpers ---

func isEmail(s string) bool {
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(s)
}

func generateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func generateNonce() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%032d", rand.Int63())
}

// isStrongPassword checks if password meets strength requirements:
// - At least 8 characters
// - Contains at least one uppercase letter
// - Contains at least one lowercase letter
// - Contains at least one digit
// - Contains at least one special character (@$!%*?&)
func isStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasDigit   = false
		hasSpecial = false
	)

	specialChars := "@$!%*?&"

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasDigit = true
		default:
			for _, special := range specialChars {
				if char == special {
					hasSpecial = true
					break
				}
			}
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial
}

func extractJTI(tokenStr string) (string, error) {
	// Use ParseUnverified to extract JTI without validating the token
	// This is safe here because we only need the JTI for blacklisting
	token, _, err := jwt.NewParser().ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		return "", apperror.ErrInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if jti, ok := claims["jti"].(string); ok {
			return jti, nil
		}
	}

	return "", apperror.ErrInvalidToken
}

// invalidateUsernameCache removes the cached username availability check
func (s *authService) invalidateUsernameCache(username string) {
	if s.redisClient == nil {
		return
	}

	ctx, cancel := util.NewDefaultRedisContext()
	defer cancel()

	cacheKey := fmt.Sprintf("username_exists:%s", username)
	// Ignore error, cache invalidation is not critical
	_ = s.redisClient.Del(ctx, cacheKey).Err()
}
