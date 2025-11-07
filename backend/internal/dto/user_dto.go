package dto

import (
	"github.com/giakiet05/lkforum/internal/model"
)

// --- Request DTOs ---

type UserRegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserLoginRequest struct {
	Identifier string `json:"identifier" binding:"required"`
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

// UserProfileUpdateRequest defines the fields a user can update for their own profile.
type UserProfileUpdateRequest struct {
	Bio      *string `json:"bio"`
	Location *string `json:"location"`
	Website  *string `json:"website"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// --- Response DTOs ---

// UserProfileResponse contains public profile information.
type UserProfileResponse struct {
	Avatar   model.Image `json:"avatar,omitempty"`
	Cover    model.Image `json:"cover,omitempty"`
	Bio      string      `json:"bio,omitempty"`	
	Location string      `json:"location,omitempty"`
	Website  string      `json:"website,omitempty"`
}

// UserResponse is the main user object returned in API responses.
type UserResponse struct {
	ID         string              `json:"id"`	
	Username   string              `json:"username"`
	Email      string              `json:"email,omitempty"`
	Reputation int                 `json:"reputation"`
	Title      string              `json:"title"`
	Role       model.Role          `json:"role"`
	Provider   model.AuthProvider  `json:"provider"`
	IsVerified bool                `json:"is_verified"`
	Profile    UserProfileResponse `json:"profile"`
}

// AuthResponse is returned on successful login or registration.
type AuthResponse struct {
	User         *UserResponse `json:"user"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
}

func FromUser(u *model.User) *UserResponse {
	if u == nil {
		return nil
	}
	resp := &UserResponse{
		ID:         u.ID.Hex(),
		Username:   u.Username,
		Email:      u.Email,
		Reputation: u.Reputation,
		Title:      calculateTitle(u.Reputation),
		Role:       u.Role,
		Provider:   u.Provider,
		IsVerified: u.IsVerified,
	}

	if u.RoleContent.AsUser != nil {
		resp.Profile = UserProfileResponse{
			Avatar:   u.RoleContent.AsUser.Avatar,
			Cover:    u.RoleContent.AsUser.Cover,
			Bio:      u.RoleContent.AsUser.Bio,
			Location: u.RoleContent.AsUser.Location,
			Website:  u.RoleContent.AsUser.Website,
		}
	}

	return resp
}

func FromUsers(users []*model.User) []*UserResponse {
	responses := make([]*UserResponse, len(users))
	for i, u := range users {
		userResponse := FromUser(u)
		userResponse.Email = ""
		responses[i] = userResponse
	}
	return responses
}

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
