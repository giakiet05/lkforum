package dto

type CreateMembershipRequest struct {
	UserID      string `json:"user_id"`
	CommunityID string `json:"community_id"`
}

type DeleteMembershipRequest struct {
	UserID      string `json:"user_id"`
	CommunityID string `json:"community_id"`
}
