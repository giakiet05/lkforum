package dto

import (
	"fmt"
	"time"

	"github.com/giakiet05/lkforum/internal/model"
)

// --- Request DTOs ---

// New Registration Flow (Verify Email First)

// GetUsersQuery contains query parameters for searching and paginating users
type GetUsersQuery struct {
	Username string `form:"username"`
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
}

// UserProfileUpdateRequest defines the fields a user can update for their own profile.
type UserProfileUpdateRequest struct {
	Bio         *string           `json:"bio" binding:"omitempty,max=500"`
	Gender      *string           `json:"gender" binding:"omitempty"`
	DateOfBirth *string           `json:"date_of_birth"` // ISO 8601 format: "2000-01-15"
	Location    *string           `json:"location"`
	Interests   []string          `json:"interests" binding:"omitempty,max=10,dive"`
	SocialLinks *SocialLinksInput `json:"social_links"`
}

type SocialLinksInput struct {
	Website  *string `json:"website" binding:"omitempty,url"`
	Facebook *string `json:"facebook"`
	YouTube  *string `json:"youtube"`
	GitHub   *string `json:"github"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// --- Response DTOs ---

// UserProfileResponse contains public profile information.
type UserProfileResponse struct {
	Avatar      *model.Image           `json:"avatar,omitempty"`
	Cover       *model.Image           `json:"cover,omitempty"`
	Bio         *string                `json:"bio,omitempty"`
	Gender      *string                `json:"gender,omitempty"`
	DateOfBirth *time.Time             `json:"date_of_birth,omitempty"`
	Age         *int                   `json:"age,omitempty"`
	Location    *string                `json:"location,omitempty"`
	Interests   []string               `json:"interests,omitempty"`
	SocialLinks *model.SocialLinks     `json:"social_links,omitempty"`
	Stats       *ActivityStatsResponse `json:"stats,omitempty"`
}

type ActivityStatsResponse struct {
	PostCount    int    `json:"post_count"`
	CommentCount int    `json:"comment_count"`
	TotalUpvotes int    `json:"total_upvotes"`
	MemberSince  string `json:"member_since"`
	LastActive   string `json:"last_active"`
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
		profile := UserProfileResponse{
			Avatar:      u.RoleContent.AsUser.Avatar,
			Cover:       u.RoleContent.AsUser.Cover,
			Bio:         u.RoleContent.AsUser.Bio,
			SocialLinks: u.RoleContent.AsUser.SocialLinks,
		}

		// Convert Gender to string
		if u.RoleContent.AsUser.Gender != nil {
			genderStr := string(*u.RoleContent.AsUser.Gender)
			profile.Gender = &genderStr
		}

		// Calculate Age from DateOfBirth
		if u.RoleContent.AsUser.DateOfBirth != nil {
			age := calculateAge(*u.RoleContent.AsUser.DateOfBirth)
			profile.Age = &age
		}

		// Convert Location to string
		if u.RoleContent.AsUser.Location != nil {
			locationStr := string(*u.RoleContent.AsUser.Location)
			profile.Location = &locationStr
		}

		// Convert Interests to []string
		if len(u.RoleContent.AsUser.Interests) > 0 {
			interests := make([]string, len(u.RoleContent.AsUser.Interests))
			for i, interest := range u.RoleContent.AsUser.Interests {
				interests[i] = string(interest)
			}
			profile.Interests = interests
		}

		// Map ActivityStats
		if u.RoleContent.AsUser.Stats != nil {
			profile.Stats = &ActivityStatsResponse{
				PostCount:    u.RoleContent.AsUser.Stats.PostCount,
				CommentCount: u.RoleContent.AsUser.Stats.CommentCount,
				TotalUpvotes: u.RoleContent.AsUser.Stats.TotalUpvotes,
				MemberSince:  formatMemberSince(u.RoleContent.AsUser.Stats.JoinedAt),
				LastActive:   formatLastActive(u.RoleContent.AsUser.Stats.LastActiveAt),
			}
		} else {
			// Provide default stats if not set (for backward compatibility with old users)
			profile.Stats = &ActivityStatsResponse{
				PostCount:    0,
				CommentCount: 0,
				TotalUpvotes: 0,
				MemberSince:  formatMemberSince(u.CreatedAt),
				LastActive:   formatLastActive(u.CreatedAt),
			}
		}

		resp.Profile = profile
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

func calculateAge(birthDate time.Time) int {
	now := time.Now()
	age := now.Year() - birthDate.Year()
	if now.YearDay() < birthDate.YearDay() {
		age--
	}
	return age
}

func formatMemberSince(joinedAt time.Time) string {
	monthNames := map[time.Month]string{
		time.January: "Jan", time.February: "Feb", time.March: "Mar",
		time.April: "Apr", time.May: "May", time.June: "Jun",
		time.July: "Jul", time.August: "Aug", time.September: "Sep",
		time.October: "Oct", time.November: "Nov", time.December: "Dec",
	}
	return fmt.Sprintf("Member since %s %d", monthNames[joinedAt.Month()], joinedAt.Year())
}

func formatLastActive(lastActive time.Time) string {
	now := time.Now()
	duration := now.Sub(lastActive)

	switch {
	case duration < time.Minute:
		return "Active now"
	case duration < time.Hour:
		minutes := int(duration.Minutes())
		if minutes == 1 {
			return "Active 1 minute ago"
		}
		return fmt.Sprintf("Active %d minutes ago", minutes)
	case duration < 24*time.Hour:
		hours := int(duration.Hours())
		if hours == 1 {
			return "Active 1 hour ago"
		}
		return fmt.Sprintf("Active %d hours ago", hours)
	case duration < 7*24*time.Hour:
		days := int(duration.Hours() / 24)
		if days == 1 {
			return "Active 1 day ago"
		}
		return fmt.Sprintf("Active %d days ago", days)
	default:
		return fmt.Sprintf("Active on %s %d", lastActive.Month().String()[:3], lastActive.Day())
	}
}
