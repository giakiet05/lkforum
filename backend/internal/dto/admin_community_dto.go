package dto

// AdminBanCommunityRequest represents request to ban a community
type AdminBanCommunityRequest struct {
	Reason string `json:"reason" binding:"required,max=500"`
}

// GetCommunitiesAdminQuery represents query parameters for admin community listing
type GetCommunitiesAdminQuery struct {
	Status   string `form:"status"`    // active, banned, deleted, all
	Name     string `form:"name"`      // search by name
	Page     int    `form:"page"`      // page number
	PageSize int    `form:"page_size"` // items per page
}