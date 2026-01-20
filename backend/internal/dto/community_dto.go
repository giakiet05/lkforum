package dto

import (
	"time"

	"github.com/giakiet05/lkforum/internal/model"
)

type CreateCommunityRequest struct {
	Name          string                 `json:"name" binding:"required,min=3,max=50"`
	Description   *string                `json:"description,omitempty" binding:"max=500"`
	Avatar        *string                `json:"avatar,omitempty"`
	Banner        *string                `json:"banner,omitempty"`
	Setting       model.CommunitySetting `json:"setting,omitempty"`
	Rules         []model.CommunityRule  `json:"rules,omitempty"`
	Moderators    []model.Moderator      `json:"moderators,omitempty"`
	CreatorName   string                 `json:"creator_name,omitempty"`
	CreatorAvatar string                 `json:"creator_avatar,omitempty"`
	Is18Plus      bool                   `json:"is_18_plus"`
}

type UpdateCommunityRequest struct {
	CommunityID string                  `json:"id" binding:"required"`
	Description *string                 `json:"description,omitempty"`
	Avatar      *string                 `json:"avatar,omitempty"`
	Banner      *string                 `json:"banner,omitempty"`
	Setting     *model.CommunitySetting `json:"setting,omitempty"`
	Rules       *[]model.CommunityRule  `json:"rules,omitempty"`
}

type AddModeratorRequest struct {
	CommunityID    string   `json:"id" binding:"required"`
	AddedModerator []string `json:"added_moderator" binding:"required"`
}

type RemoveModeratorRequest struct {
	CommunityID      string   `json:"id" binding:"required"`
	RemovedModerator []string `json:"removed_moderator" binding:"required"`
}

type CommunityBanPostRequest struct {
	CommunityID string  `json:"community_id" binding:"required"`
	PostID      string  `json:"post_id" binding:"required"`
	Reason      *string `json:"reason,omitempty" binding:"max=500"`
}

type CommunityUnbanPostRequest struct {
	CommunityID string `json:"community_id" binding:"required"`
	PostID      string `json:"post_id" binding:"required"`
}

type CommunityBanUserRequest struct {
	CommunityID string `json:"community_id" binding:"required"`
	UserID      string `json:"user_id" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Reason      string `json:"reason" binding:"max=500"`
	LengthDays  int    `json:"length_days" binding:"required"`
}

type UnbanUserRequest struct {
	CommunityID string `json:"community_id" binding:"required"`
	UserID      string `json:"user_id" binding:"required"`
}

type ModeratePostRequest struct {
	Approve bool    `json:"approve"`
	Reason  *string `json:"reason,omitempty"`
}

type CommunityResponse struct {
	ID             string                 `json:"id"`
	Name           string                 `json:"name"`
	Description    *string                `json:"description"`
	Avatar         *string                `json:"avatar"`
	Banner         *string                `json:"banner"`
	Setting        model.CommunitySetting `json:"setting"`
	Rules          []model.CommunityRule  `json:"rules"`
	Moderators     []model.Moderator      `json:"moderators"`
	PostCount      int64                  `json:"post_count"`
	MemberCount    int64                  `json:"member_count"`
	Is18Plus       bool                   `json:"is_18_plus"`
	CreateByID     string                 `json:"create_by_id,omitempty"`
	CreateByName   string                 `json:"create_by_name,omitempty"`
	CreateByAvatar string                 `json:"create_by_avatar,omitempty"`
	CreateAt       time.Time              `json:"create_at,omitempty"`
}

func FromCommunities(communities []*model.Community) []*CommunityResponse {
	var communityResponses []*CommunityResponse
	for _, community := range communities {
		communityResponses = append(communityResponses, FromCommunity(community))
	}
	return communityResponses
}

func FromCommunity(community *model.Community) *CommunityResponse {
	return &CommunityResponse{
		ID:             community.ID.Hex(),
		Name:           community.Name,
		Description:    community.Description,
		Avatar:         community.Avatar,
		Banner:         community.Banner,
		Setting:        community.Setting,
		Rules:          community.Rules,
		Moderators:     community.Moderators,
		PostCount:      community.PostCount,
		MemberCount:    community.MemberCount,
		Is18Plus:       community.Is18Plus,
		CreateByID:     community.CreateByID.Hex(),
		CreateByName:   community.CreateByName,
		CreateByAvatar: community.CreateByAvatar,
		CreateAt:       community.CreateAt,
	}
}
