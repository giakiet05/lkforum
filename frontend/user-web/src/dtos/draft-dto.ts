import type { PostType } from "./post-dto";

export interface CreateDraftRequest {
    community_id?: string;
    type?: PostType;
    title?: string;
    text?: string;
    tags?: string[];
    images?: Image[];
    videos?: Video[];
    poll?: Poll;
}

export interface UpdateDraftRequest {
    community_id?: string;
    type?: PostType;
    title?: string;
    text?: string;
    tags?: string[];
    images?: Image[];
    videos?: Video[];
    poll?: Poll;
}

export interface Image {
    url: string;
    width?: number;
    height?: number;
}

export interface Video {
    url: string;
    thumbnail_url?: string;
    duration?: number;
}

export interface PollOption {
    id?: string;
    text: string;
    votes?: number;
}

export interface Poll {
    question: string;
    options: PollOption[];
    expires_at?: string;
    allow_multiple?: boolean;
}

export interface DraftSummaryResponse {
    id: string;
    title?: string;
    updated_at: string;
}

export interface DraftResponse {
    id: string;
    author_id: string;
    community_id?: string;
    type?: PostType;
    title?: string;
    content?: PostContent;
    tags?: string[];
    created_at: string;
    updated_at: string;
}

export interface PostContent {
    text?: string;
    images?: Image[];
    videos?: Video[];
    poll?: Poll;
}

export interface PaginatedDraftsResponse {
    drafts: DraftSummaryResponse[];
    pagination: {
        page: number;
        page_size: number;
        total: number;
        total_pages?: number;
    };
}
