import type { Image } from './post-dto';

// --- Request DTOs ---

export interface CreateCommentRequest {
    post_id: string;
    parent_id?: string;
    content: string;
}

export interface GetCommentsFilterQuery {
    post_id?: string;
    parent_id?: string;
    user_id?: string;
    content?: string;
    page?: number;
    page_size?: number;
}

export interface GetCommentByPostIDQuery {
    post_id: string;
    depth?: number;
    children_page_size?: number;
    page?: number;
    page_size?: number;
}

// --- Response DTOs ---

export interface CommentAuthor {
    id: string;
    username: string;
    avatar?: Image;
}

export interface CommentResponse {
    id: string;
    author: CommentAuthor;
    post_id: string;
    parent_id?: string;
    children: CommentResponse[];
    content: string;
    created_at: string; // Formatted date string
    is_deleted: boolean;
}
