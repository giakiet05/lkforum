// --- Request DTOs ---

export interface CreateChannelRequest {
    member_1: string;
    member_1_username: string;
    member_1_avatar: string;
    member_2: string;
    member_2_username: string;
    member_2_avatar: string;
}

export interface GetChannelByUserIDQuery {
    user_id: string;
    page?: number;
    page_size?: number;
}

export type ChannelStatus = "active" | "archived" | "blocked";

export interface UpdateChannelRequest {
    channel_id: string;
    nickname?: string;
    notification?: boolean;
    typing_indicator?: boolean;
    status?: ChannelStatus;
}

// --- Response DTOs ---

export interface ChannelMemberResponse {
    user_id: string;
    username: string;
    avatar: string;
}

export interface ChannelSettingResponse {
    user_id: string;
    nickname?: string;
    notification: boolean;
    typing_indicator: boolean;
}

export interface ChannelResponse {
    id: string;
    members: ChannelMemberResponse[];
    settings: ChannelSettingResponse[];
    status: ChannelStatus;
    created_at: string; // ISO 8601 format
    updated_at: string; // ISO 8601 format
}

export interface Pagination {
    page: number;
    page_size: number;
    total: number;
}

export interface PaginatedChannelsResponse {
    channels: ChannelResponse[];
    pagination: Pagination;
}
