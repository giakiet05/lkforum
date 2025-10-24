package dto

import "github.com/giakiet05/lkforum/internal/model"

// Request DTOs

type UserRegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserLoginRequest struct {
	Identifier string `json:"identifier" binding:"required"` // Username or Email
	Password   string `json:"password" binding:"required"`
}

type VerifyEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=6"`
}

type ResendVerificationEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type CompleteGoogleSetupRequest struct {
	SetupToken string `json:"setup_token" binding:"required"`
	Username   string `json:"username" binding:"required,min=3,max=20"`
}

type UserUpdateRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Response DTOs

type UserResponse struct {
	ID         string             `json:"id"`
	Username   string             `json:"username"`
	Email      string             `json:"email,omitempty"`
	Reputation int                `json:"reputation"`
	Title      string             `json:"title"`
	Role       model.Role         `json:"role"`
	Provider   model.AuthProvider `json:"provider"`
	IsVerified bool               `json:"is_verified"`
}

type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
}

func FromUser(u *model.User) UserResponse {
	return UserResponse{
		ID:         u.ID.Hex(),
		Username:   u.Username,
		Email:      u.Email,
		Reputation: u.Reputation,
		Title:      calculateTitle(u.Reputation),
		Role:       u.Role,
		Provider:   u.Provider,
		IsVerified: u.IsVerified,
	}
}

func FromUsers(users []*model.User) []UserResponse {
	responses := make([]UserResponse, 0, len(users))
	for _, u := range users {
		responses = append(responses, FromUser(u))
	}
	return responses
}

// calculateTitle determines the user's title based on their reputation score.
func calculateTitle(reputation int) string {
	switch {
	case reputation >= 10000:
		return "Huyền thoại"
	case reputation >= 2000:
		return "Lão làng"
	case reputation >= 500:
		return "Cây bút trẻ"
	case reputation >= 100:
		return "Thành viên tích cực"
	case reputation >= 0:
		return "Lính mới"
	default:
		return "Người qua đường"
	}
}
