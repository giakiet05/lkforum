package dto

import "github.com/giakiet05/lkforum/internal/model"

// UpdateSettingsRequest allows updating all or partial settings
type UpdateSettingsRequest struct {
	Appearance    *AppearanceSettingsInput   `json:"appearance,omitempty"`
	Notifications *NotificationSettingsInput `json:"notifications,omitempty"`
	Privacy       *PrivacySettingsInput      `json:"privacy,omitempty"`
	Content       *ContentSettingsInput      `json:"content,omitempty"`
}

type AppearanceSettingsInput struct {
	Theme    *string `json:"theme" binding:"omitempty,oneof=light dark auto"`
	FontSize *string `json:"font_size" binding:"omitempty,oneof=small medium large"`
}

type NotificationSettingsInput struct {
	InAppEnabled    *bool `json:"in_app_enabled"`
	EmailEnabled    *bool `json:"email_enabled"`
	NotifyOnComment *bool `json:"notify_on_comment"`
	NotifyOnMention *bool `json:"notify_on_mention"`
	NotifyOnUpvote  *bool `json:"notify_on_upvote"`
	NotifyOnMessage *bool `json:"notify_on_message"`
}

type PrivacySettingsInput struct {
	ShowProfile         *bool `json:"show_profile"`
	ShowEmail           *bool `json:"show_email"`
	ShowPostHistory     *bool `json:"show_post_history"`
	AllowDirectMessages *bool `json:"allow_direct_messages"`
	AllowMentions       *bool `json:"allow_mentions"`
}

type ContentSettingsInput struct {
	AllowNSFW *bool `json:"allow_nsfw"`
}

// SettingsResponse is the response DTO for user settings
type SettingsResponse struct {
	Appearance    model.AppearanceSettings   `json:"appearance"`
	Notifications model.NotificationSettings `json:"notifications"`
	Privacy       model.PrivacySettings      `json:"privacy"`
	Content       model.ContentSettings      `json:"content"`
}

// FromSettings converts model.UserSettings to SettingsResponse
func FromSettings(settings *model.UserSettings) *SettingsResponse {
	if settings == nil {
		settings = model.NewDefaultSettings()
	}
	return &SettingsResponse{
		Appearance:    settings.Appearance,
		Notifications: settings.Notifications,
		Privacy:       settings.Privacy,
		Content:       settings.Content,
	}
}
