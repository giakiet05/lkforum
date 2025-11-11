// --- Request DTOs ---

export interface AppearanceSettingsInput {
    theme?: "light" | "dark" | "auto";
    font_size?: "small" | "medium" | "large";
}

export interface NotificationSettingsInput {
    in_app_enabled?: boolean;
    email_enabled?: boolean;
    notify_on_comment?: boolean;
    notify_on_mention?: boolean;
    notify_on_upvote?: boolean;
    notify_on_message?: boolean;
}

export interface PrivacySettingsInput {
    show_profile?: boolean;
    show_email?: boolean;
    show_post_history?: boolean;
    allow_direct_messages?: boolean;
    allow_mentions?: boolean;
}

export interface ContentSettingsInput {
    allow_nsfw?: boolean;
}

export interface UpdateSettingsRequest {
    appearance?: AppearanceSettingsInput;
    notifications?: NotificationSettingsInput;
    privacy?: PrivacySettingsInput;
    content?: ContentSettingsInput;
}

// --- Response DTOs ---

export interface AppearanceSettings {
    theme: "light" | "dark" | "auto";
    font_size: "small" | "medium" | "large";
}

export interface NotificationSettings {
    in_app_enabled: boolean;
    email_enabled: boolean;
    notify_on_comment: boolean;
    notify_on_mention: boolean;
    notify_on_upvote: boolean;
    notify_on_message: boolean;
}

export interface PrivacySettings {
    show_profile: boolean;
    show_email: boolean;
    show_post_history: boolean;
    allow_direct_messages: boolean;
    allow_mentions: boolean;
}

export interface ContentSettings {
    allow_nsfw: boolean;
}

export interface SettingsResponse {
    appearance: AppearanceSettings;
    notifications: NotificationSettings;
    privacy: PrivacySettings;
    content: ContentSettings;
}
