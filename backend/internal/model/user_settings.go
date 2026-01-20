package model

// UserSettings contains all user preference settings
type UserSettings struct {
	Appearance    AppearanceSettings   `bson:"appearance" json:"appearance"`
	Notifications NotificationSettings `bson:"notifications" json:"notifications"`
	Privacy       PrivacySettings      `bson:"privacy" json:"privacy"`
	Content       ContentSettings      `bson:"content" json:"content"`
}

// AppearanceSettings controls visual preferences
type AppearanceSettings struct {
	Theme    string `bson:"theme" json:"theme"`         // "light", "dark", "auto"
	FontSize string `bson:"font_size" json:"font_size"` // "small", "medium", "large"
}

// NotificationSettings controls notification preferences
// InApp and Email are channels - notify options apply to both
type NotificationSettings struct {
	InAppEnabled    bool `bson:"in_app_enabled" json:"in_app_enabled"`       // Show notifications in website
	EmailEnabled    bool `bson:"email_enabled" json:"email_enabled"`         // Send email notifications
	NotifyOnComment bool `bson:"notify_on_comment" json:"notify_on_comment"` // Applies to both in-app & email
	NotifyOnMention bool `bson:"notify_on_mention" json:"notify_on_mention"` // Applies to both in-app & email
	NotifyOnUpvote  bool `bson:"notify_on_upvote" json:"notify_on_upvote"`   // Applies to both in-app & email
	NotifyOnMessage bool `bson:"notify_on_message" json:"notify_on_message"` // Applies to both in-app & email
}

// PrivacySettings controls privacy and interaction preferences
type PrivacySettings struct {
	ShowProfile         bool `bson:"show_profile" json:"show_profile"`                   // If false, entire profile is private
	ShowEmail           bool `bson:"show_email" json:"show_email"`                       // Show email on profile
	ShowPostHistory     bool `bson:"show_post_history" json:"show_post_history"`         // Show post history
	AllowDirectMessages bool `bson:"allow_direct_messages" json:"allow_direct_messages"` // Allow others to send DMs
	AllowMentions       bool `bson:"allow_mentions" json:"allow_mentions"`               // Allow others to mention you
}

// ContentSettings controls content filtering preferences
type ContentSettings struct {
	AllowNSFW bool `bson:"allow_nsfw" json:"allow_nsfw"`
}

// Theme constants
const (
	ThemeLight = "light"
	ThemeDark  = "dark"
	ThemeAuto  = "auto"
)

// FontSize constants
const (
	FontSizeSmall  = "small"
	FontSizeMedium = "medium"
	FontSizeLarge  = "large"
)

// NewDefaultSettings returns default user settings
func NewDefaultSettings() *UserSettings {
	return &UserSettings{
		Appearance: AppearanceSettings{
			Theme:    ThemeLight,
			FontSize: FontSizeMedium,
		},
		Notifications: NotificationSettings{
			InAppEnabled:    true,
			EmailEnabled:    false, // Disabled by default to avoid spam
			NotifyOnComment: true,
			NotifyOnMention: true,
			NotifyOnUpvote:  true,
			NotifyOnMessage: true,
		},
		Privacy: PrivacySettings{
			ShowProfile:         true,  // Profile public by default
			ShowEmail:           false, // Email private by default
			ShowPostHistory:     true,
			AllowDirectMessages: true,
			AllowMentions:       true,
		},
		Content: ContentSettings{
			AllowNSFW: false,
		},
	}
}

// IsValidTheme checks if theme is valid
func IsValidTheme(theme string) bool {
	return theme == ThemeLight || theme == ThemeDark || theme == ThemeAuto
}

// IsValidFontSize checks if font size is valid
func IsValidFontSize(fontSize string) bool {
	return fontSize == FontSizeSmall || fontSize == FontSizeMedium || fontSize == FontSizeLarge
}
