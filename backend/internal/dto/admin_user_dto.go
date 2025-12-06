package dto

import "time"

// BanUserRequest is the request to ban a user
type BanUserRequest struct {
	Reason   string     `json:"reason" binding:"required,max=500"`
	BanUntil *time.Time `json:"ban_until,omitempty"` // null = permanent ban
}

// GetUsersAdminQuery is the query for admin to get all users
type GetUsersAdminQuery struct {
	Username string `form:"username"`
	Status   string `form:"status"` // all, active, banned, deleted
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}
