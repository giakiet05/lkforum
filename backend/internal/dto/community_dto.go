package dto

import "github.com/giakiet05/lkforum/internal/model"

type CreateCommunityRequest struct {
	Name        string                 `json:"name" validate:"required,min=3,max=50"`
	Description *string                `json:"description,omitempty" validate:"max=500"`
	Avatar      *string                `json:"avatar,omitempty"`
	Banner      *string                `json:"banner,omitempty"`
	Setting     model.CommunitySetting `json:"setting,omitempty"`
	Moderators  []model.Moderator      `json:"moderators,omitempty"`
}

type UpdateCommunityRequest struct {
	CommunityID string                  `json:"id" validate:"required"`
	Name        *string                 `json:"name,omitempty" validate:"omitempty,min=3,max=50"`
	Description *string                 `json:"description,omitempty" validate:"omitempty,max=500"`
	Avatar      *string                 `json:"avatar,omitempty"`
	Banner      *string                 `json:"banner,omitempty"`
	Setting     *model.CommunitySetting `json:"setting,omitempty"`
}

type AddModeratorRequest struct {
	CommunityID      string   `json:"id" validate:"required"`
	AddedModeratorID []string `json:"added_moderator,omitempty"`
}

type RemoveModeratorRequest struct {
	CommunityID        string   `json:"id" validate:"required"`
	RemovedModeratorID []string `json:"removed_moderator,omitempty"`
}
