package dto

import (
	"time"

	"github.com/giakiet05/lkforum/internal/model"
)

type CreateChannelRequest struct {
	Member1         string `json:"member_1"`
	Member1Username string `json:"member_1_username"`
	Member1Avatar   string `json:"member_1_avatar"`
	Member2         string `json:"member_2"`
	Member2Username string `json:"member_2_username"`
	Member2Avatar   string `json:"member_2_avatar"`
}

type GetChannelByUserIDQuery struct {
	UserID   string `form:"user_id"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}

type UpdateChannelRequest struct {
	ChannelID       string               `json:"channel_id"`
	Nickname        *string              `json:"nickname,omitempty"`
	Notification    *bool                `json:"notification,omitempty"`
	TypingIndicator *bool                `json:"typing_indicator,omitempty"`
	Status          *model.ChannelStatus `json:"status,omitempty"`
}

type ChannelResponse struct {
	ID        string                   `json:"id"`
	Members   []ChannelMemberResponse  `json:"members"`
	Settings  []ChannelSettingResponse `json:"settings"`
	Status    model.ChannelStatus      `json:"status"`
	CreatedAt time.Time                `json:"created_at"`
	UpdatedAt time.Time                `json:"updated_at"`
}

type ChannelMemberResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

type ChannelSettingResponse struct {
	UserID          string  `json:"user_id"`
	Nickname        *string `json:"nickname,omitempty"`
	Notification    bool    `json:"notification"`
	TypingIndicator bool    `json:"typing_indicator"`
}

func FromChannel(channel *model.Channel) *ChannelResponse {
	members := make([]ChannelMemberResponse, len(channel.Members))
	for i, m := range channel.Members {
		members[i] = ChannelMemberResponse{
			UserID:   m.UserID.Hex(),
			Username: m.Username,
			Avatar:   m.Avatar,
		}
	}

	settings := make([]ChannelSettingResponse, len(channel.Settings))
	for i, s := range channel.Settings {
		settings[i] = ChannelSettingResponse{
			UserID:          s.UserID.Hex(),
			Nickname:        s.Nickname,
			Notification:    s.Notification,
			TypingIndicator: s.TypingIndicator,
		}
	}

	return &ChannelResponse{
		ID:        channel.ID.Hex(),
		Members:   members,
		Settings:  settings,
		Status:    channel.Status,
		CreatedAt: channel.CreatedAt,
		UpdatedAt: channel.UpdatedAt,
	}
}

func FromChannels(channels []model.Channel) []*ChannelResponse {
	responses := make([]*ChannelResponse, len(channels))
	for i, ch := range channels {
		responses[i] = FromChannel(&ch)
	}
	return responses
}
