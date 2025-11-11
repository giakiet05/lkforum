// --- Enums ---

export enum Role {
    User = 'user',
    Admin = 'admin',
    Moderator = 'moderator',
}

export enum AuthProvider {
    Local = 'local',
    Google = 'google',
}

// --- Shared Types ---

export interface Image {
    url: string;
    public_id?: string;
    width?: number;
    height?: number;
}

export interface SocialLinks {
    website?: string;
    facebook?: string;
    youtube?: string;
    github?: string;
}

// --- Request DTOs ---

export interface SocialLinksInput {
    website?: string;
    facebook?: string;
    youtube?: string;
    github?: string;
}

export interface UserProfileUpdateRequest {
    bio?: string;
    gender?: string;
    date_of_birth?: string; // ISO 8601 format: "2000-01-15"
    location?: string;
    interests?: string[];
    social_links?: SocialLinksInput;
}

export interface ChangePasswordRequest {
    old_password: string;
    new_password: string;
}

// --- Response DTOs ---

export interface ActivityStatsResponse {
    post_count: number;
    comment_count: number;
    total_upvotes: number;
    member_since: string;
    last_active: string;
}

export interface UserProfileResponse {
    avatar?: Image;
    cover?: Image;
    bio?: string;
    gender?: string;
    age?: number;
    location?: string;
    interests?: string[];
    social_links?: SocialLinks;
    stats?: ActivityStatsResponse;
}

export interface UserResponse {
    id: string;
    username: string;
    email?: string;
    reputation: number;
    title: string;
    role: Role;
    provider: AuthProvider;
    is_verified: boolean;
    profile: UserProfileResponse;
}
