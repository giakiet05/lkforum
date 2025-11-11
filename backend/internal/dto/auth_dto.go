package dto

type SendEmailVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyEmailCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=6"`
}

type CompleteRegistrationRequest struct {
	VerificationToken string `json:"verification_token" binding:"required"`
	Username          string `json:"username" binding:"required,min=3,max=20"`
	Password          string `json:"password" binding:"required,min=6"`
}

type ResendOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// Login
type UserLoginRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type CompleteGoogleSetupRequest struct {
	SetupToken string `json:"setup_token" binding:"required"`
	Username   string `json:"username" binding:"required,min=3,max=20"`
}
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutRequest struct {
	AccessToken  string `json:"access_token" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// AuthResponse is returned on successful login or registration.
type AuthResponse struct {
	User         *UserResponse `json:"user"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
