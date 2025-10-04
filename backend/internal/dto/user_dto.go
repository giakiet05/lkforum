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

type UserUpdateRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// Response DTOs

type UserResponse struct {
	ID       string     `json:"id"`
	Username string     `json:"username"`
	Email    string     `json:"email,omitempty"`
	Role     model.Role `json:"role"`
}

type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
}

func FromUser(u *model.User) UserResponse {
	return UserResponse{
		ID:       u.ID.Hex(),
		Username: u.Username,
		Email:    u.Email,
		Role:     u.Role,
	}
}

func FromUsers(users []*model.User) []UserResponse {
	responses := make([]UserResponse, 0, len(users))
	for _, u := range users {
		responses = append(responses, FromUser(u))
	}
	return responses
}
