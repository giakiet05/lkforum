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
    avatar?: string;
}

export interface CreateCommunityRequest {
    name: string;
    description?: string;
    avatar?: string;
    banner?: string;
    setting?: CommunitySetting;
    moderators?: Moderator[];
    creator_name?: string;
    creator_avatar?: string;
}

export interface UpdateCommunityRequest {
    id: string;
    description?: string;
    avatar?: string;
    banner?: string;
    setting?: CommunitySetting;
}

export interface ModeratorDTO {
    id: string;
    username: string;
}

export interface AddModeratorRequest {
    id: string;
    added_moderator: ModeratorDTO[];
}

export interface RemoveModeratorRequest {
    id: string;
    removed_moderator: string[];
}

// --- Response DTOs ---

export interface CommunityResponse {
    id: string;
    name: string;
    description: string;
    avatar: string;
    banner: string;
    setting: CommunitySetting;
    moderators: Moderator[];
    post_count: number;
    member_count: number;
    create_by_id?: string;
    create_by_name?: string;
    create_by_avatar?: string;
}
