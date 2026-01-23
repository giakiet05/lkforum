// --- Request DTOs ---

export interface CreatePostRequest {
    community_id: string;
    title: string;
    type: "text" | "poll" | "image" | "video";
    text?: string;
    tags?: string[];
    poll?: CreatePollRequest;
}

export interface CreatePollRequest {
    question: string;
    options: string[];
    expires_at?: string; // ISO 8601 format
    allow_multiple?: boolean;
}

export interface UpdatePostRequest {
    title?: string;
    text?: string;
    tags?: string[];
}

export interface UpdatePollRequest {
    question?: string;
    expires_at?: string; // ISO 8601 format
    allow_multiple?: boolean;
}

export interface UpdatePollOptionRequest {
    text: string;
}

export interface AddPollOptionsRequest {
    options: string[];
}

export interface RemovePollOptionsRequest {
    option_ids: string[];
}

export interface PostVoteRequest {
    value: boolean | null; // true = upvote, false = downvote, null = remove vote
}

export interface PollVoteRequest {
    option_id: string;
}

export interface RemoveImagesRequest {
    public_ids: string[];
}

export interface ReportPostRequest {
    reason: string;
    description?: string;
}

export interface GetPostsQuery {
    community_id?: string;
    author_id?: string;
    type?: string;
    sort?: string;
    time?: string;
    feed_type?: "home" | "popular" | "explore" | "all"; // Feed filtering
    search?: string; // Search query
    page?: number;
    limit?: number;
}

// --- Response DTOs ---

export interface Image {
    public_id: string;
    url: string;
    width?: number;
    height?: number;
}

export interface Video {
    public_id: string;
    url: string;
    duration?: number;
    thumbnail?: string;
}

export interface AuthorResponse {
    id: string;
    username: string;
    avatar?: Image;
}

export interface CommunityShortResponse {
    id: string;
    name: string;
}

export interface PollOptionResponse {
    id: string;
    text: string;
    votes: number;
    percentage: number;
}

export interface PollResponse {
    question: string;
    options: PollOptionResponse[];
    total_votes: number;
    user_vote_ids?: string[];
    expires_at?: string; // ISO 8601 format
    allow_multiple: boolean;
}

export interface PostContentResponse {
    text?: string;
    images?: Image[];
    videos?: Video[];
    poll?: PollResponse;
}

export interface VotesCountResponse {
    up: number;
    down: number;
    score: number;
}

export interface PostResponse {
    id: string;
    author: AuthorResponse;
    community: CommunityShortResponse;
    title: string;
    type: "text" | "poll" | "image" | "video";
    content: PostContentResponse;
    votes_count?: VotesCountResponse;
    user_vote?: string; // "up" or "down" or ""
    comments_count: number;
    created_at: string; // ISO 8601 format
    updated_at?: string; // ISO 8601 format
    tags?: string[];
}

export interface Pagination {
    page: number;
    page_size: number;
    total: number;
}

export interface PaginatedPostsResponse {
    posts: PostResponse[];
    pagination: Pagination;
}
