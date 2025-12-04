package dto

import "time"

// PlatformStatsResponse represents the main admin dashboard overview
type PlatformStatsResponse struct {
	Users      UserStatsData      `json:"users"`
	Communities CommunityStatsData `json:"communities"`
	Content    ContentStatsData   `json:"content"`
	Moderation ModerationStatsData `json:"moderation"`
	UpdatedAt  time.Time          `json:"updated_at"`
}

// UserStatsData contains user-related statistics
type UserStatsData struct {
	Total         int64 `json:"total"`
	ActiveToday   int64 `json:"active_today"`
	NewThisWeek   int64 `json:"new_this_week"`
	NewThisMonth  int64 `json:"new_this_month"`
	Banned        int64 `json:"banned"`
	Verified      int64 `json:"verified"`
}

// CommunityStatsData contains community-related statistics
type CommunityStatsData struct {
	Total         int64 `json:"total"`
	Active        int64 `json:"active"`
	NewThisWeek   int64 `json:"new_this_week"`
	NewThisMonth  int64 `json:"new_this_month"`
	Banned        int64 `json:"banned"`
	Private       int64 `json:"private"`
}

// ContentStatsData contains posts and comments statistics
type ContentStatsData struct {
	TotalPosts     int64 `json:"total_posts"`
	TotalComments  int64 `json:"total_comments"`
	PostsToday     int64 `json:"posts_today"`
	CommentsToday  int64 `json:"comments_today"`
	PostsThisWeek  int64 `json:"posts_this_week"`
	CommentsThisWeek int64 `json:"comments_this_week"`
}

// ModerationStatsData contains moderation-related statistics
type ModerationStatsData struct {
	PendingReports      int64 `json:"pending_reports"`
	ReportsToday        int64 `json:"reports_today"`
	PostsNeedApproval   int64 `json:"posts_need_approval"`
	BannedUsers         int64 `json:"banned_users"`
	BannedCommunities   int64 `json:"banned_communities"`
}

// GetUserStatsQuery for filtering user statistics
type GetUserStatsQuery struct {
	Period string `form:"period"` // today, week, month, year
}

// GetContentStatsQuery for filtering content statistics  
type GetContentStatsQuery struct {
	Period string `form:"period"` // today, week, month, year
	Type   string `form:"type"`   // posts, comments, all
}

// UserGrowthData for user growth trends
type UserGrowthData struct {
	Date  time.Time `json:"date"`
	Count int64     `json:"count"`
}

// ContentGrowthData for content growth trends
type ContentGrowthData struct {
	Date     time.Time `json:"date"`
	Posts    int64     `json:"posts"`
	Comments int64     `json:"comments"`
}

// UserStatsResponse for detailed user statistics
type UserStatsResponse struct {
	Overview UserStatsData    `json:"overview"`
	Growth   []UserGrowthData `json:"growth"`
}

// ContentStatsResponse for detailed content statistics
type ContentStatsResponse struct {
	Overview ContentStatsData    `json:"overview"`
	Growth   []ContentGrowthData `json:"growth"`
}