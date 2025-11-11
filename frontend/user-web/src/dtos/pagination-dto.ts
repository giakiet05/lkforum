import type { UserResponse } from './user-dto';
import type { CommunityResponse } from './community-dto';
import type { PostResponse } from './post-dto';
import type { CommentResponse } from './comment-dto';
import type { ChannelResponse } from './channel-dto';
import type { MessageResponse } from './message-dto';
import type { PostHistoryResponse } from './post-history-dto';

export interface Pagination {
    page: number;
    page_size: number;
    total: number;
}

export interface PaginatedUsersResponse {
    users: UserResponse[];
    pagination: Pagination;
}

export interface PaginatedCommunitiesResponse {
    communities: CommunityResponse[];
    pagination: Pagination;
}

export interface PaginatedPostsResponse {
    posts: PostResponse[];
    pagination: Pagination;
}

export interface PaginatedCommentsResponse {
    comments: CommentResponse[];
    pagination: Pagination;
}

export interface PaginatedChannelsResponse {
    channels: ChannelResponse[];
    pagination: Pagination;
}

export interface PaginatedMessagesResponse {
    messages: MessageResponse[];
    pagination: Pagination;
}

export interface PaginatedPostHistoryResponse {
    post_histories: PostHistoryResponse[];
    pagination: Pagination;
}
