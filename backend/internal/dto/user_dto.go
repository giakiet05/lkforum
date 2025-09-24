package dto

import "github.com/giakiet05/lkforum/internal/model"

type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
type UserRegisterRequest struct {
}

func FromUser(u *model.User) UserResponse {
	return UserResponse{
		ID:       u.ID.Hex(),
		Username: u.Username,
		Email:    u.Email,
	}
}

func FromUsers(users []*model.User) []UserResponse {
	responses := make([]UserResponse, 0, len(users))
	for _, u := range users {
		responses = append(responses, FromUser(u))
	}
	return responses
}
