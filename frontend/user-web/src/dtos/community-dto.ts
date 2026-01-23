// --- Request DTOs ---

export interface CommunitySetting {
    is_private: boolean;              // Private community (only approved members can view)
    post_require_approval: boolean;   // Posts need moderator approval
    join_require_approval: boolean;   // Join requests need moderator approval
    max_post_length: number;          // Maximum post length
}

export interface Moderator {
    user_id: string;
    username: string;
    avatar?: {
        url: string;
        public_id?: string;
    };
    is_active?: boolean;
}

export interface CommunityRule {
    title: string;
    description: string;
}

export interface CreateCommunityRequest {
    name: string;
    description?: string;
    avatar?: string;
    banner?: string;
    setting?: CommunitySetting;
    rules?: CommunityRule[];
    moderators?: Moderator[];
    creator_name?: string;
    creator_avatar?: string;
    is_18_plus?: boolean;
}

export interface UpdateCommunityRequest {
    id: string;
    description?: string;
    avatar?: string;
    banner?: string;
    setting?: CommunitySetting;
    rules?: CommunityRule[];
}

export interface ModeratorDTO {
    id: string;
    username: string;
}

export interface AddModeratorRequest {
    id: string;
    added_moderator: string[]; // Array of user IDs
}

export interface RemoveModeratorRequest {
    id: string;
    removed_moderator: string[]; // Array of user IDs
}

export interface BanUserRequest {
    community_id: string;
    user_id: string;
    type: string; // "ban" or "mute"
    reason: string;
    length_days: number;
}

export interface UnbanUserRequest {
    community_id: string;
    user_id: string;
}

export interface UserResponse {
    id: string;
    username: string;
    email?: string;
    avatar?: string;
    banner?: string;
    bio?: string;
    created_at?: string;
    profile?: {
        avatar?: {
            url: string;
            public_id?: string;
        };
        cover?: {
            url: string;
            public_id?: string;
        };
        bio?: string;
    };
}

// --- Response DTOs ---

export interface PaginationInfo {
    current_page: number;
    total_pages: number;
    total_items: number;
    page_size: number;
}

export interface PaginatedCommunitiesResponse {
    communities: CommunityResponse[];
    pagination: PaginationInfo;
}

export interface CommunityResponse {
    id: string;
    name: string;
    description: string;
    avatar: string;
    banner: string;
    setting: CommunitySetting;
    rules?: CommunityRule[];
    moderators: Moderator[];
    post_count: number;
    member_count: number;
    create_by_id?: string;
    create_by_name?: string;
    create_by_avatar?: string;
    is_18_plus?: boolean;
    created_at?: string;
}
