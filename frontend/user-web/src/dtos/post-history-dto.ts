// --- Request DTOs ---

export interface CreatePostHistoryRequest {
    user_id: string;
    post_id: string;
}

// --- Response DTOs ---

export interface PostHistoryResponse {
    id: string;
    post_id: string;
    user_id: string;
    viewed_at: string; // ISO 8601 format
}
