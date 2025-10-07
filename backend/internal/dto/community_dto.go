package dto

import (
	"github.com/giakiet05/lkforum/internal/model"
)

type CreateCommunityRequest struct {
	Name          string                 `json:"name" binding:"required,min=3,max=50"`
	Description   *string                `json:"description,omitempty" binding:"max=500"`
	Avatar        *string                `json:"avatar,omitempty"`
	Banner        *string                `json:"banner,omitempty"`
	Setting       model.CommunitySetting `json:"setting,omitempty"`
	Moderators    []model.Moderator      `json:"moderators,omitempty"`
	CreatorName   string                 `json:"creator_name,omitempty"`
	CreatorAvatar string                 `json:"creator_avatar,omitempty"`
}

type UpdateCommunityRequest struct {
	CommunityID string                  `json:"id" binding:"required"`
	Name        *string                 `json:"name,omitempty" binding:"min=3,max=50"`
	Description *string                 `json:"description,omitempty" binding:"max=500"`
	Avatar      *string                 `json:"avatar,omitempty"`
	Banner      *string                 `json:"banner,omitempty"`
	Setting     *model.CommunitySetting `json:"setting,omitempty"`
}

type ModeratorDTO struct {
	ModeratorID string `json:"id" binding:"required"`
	Username    string `json:"username" binding:"required"`
}

type AddModeratorRequest struct {
	CommunityID    string         `json:"id" binding:"required"`
	AddedModerator []ModeratorDTO `json:"added_moderator" binding:"required"`
}

type RemoveModeratorRequest struct {
	CommunityID      string         `json:"id" binding:"required"`
	RemovedModerator []ModeratorDTO `json:"removed_moderator" binding:"required"`
}

type CommunityResponse struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Avatar      string                 `json:"avatar"`
	Banner      string                 `json:"banner"`
	Setting     model.CommunitySetting `json:"setting"`
	Moderators  []model.Moderator      `json:"moderators"`
	PostCount   int64                  `json:"post_count"`
	MemberCount int64                  `json:"member_count"`
}

func FromCommunities(communities []model.Community) []CommunityResponse {
	var communityResponses []CommunityResponse
	for _, community := range communities {
		communityResponses = append(communityResponses, *FromCommunity(&community))
	}
	return communityResponses
}

func FromCommunity(community *model.Community) *CommunityResponse {
	return &CommunityResponse{
		ID:          community.ID.Hex(),
		Name:        community.Name,
		Description: *community.Description,
		Avatar:      *community.Avatar,
		Banner:      *community.Banner,
		Setting:     community.Setting,
		Moderators:  community.Moderators,
		PostCount:   community.PostCount,
		MemberCount: community.MemberCount,
	}
}
