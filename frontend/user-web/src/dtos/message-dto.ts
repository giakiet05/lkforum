// --- Request DTOs ---

export type MessageType = "text" | "image" | "video" | "file";

export interface CreateMessageRequest {
    channel_id: string;
    sender_id: string;
    type: MessageType;
    content: string;
}

export interface GetMessageFilterQuery {
    channel_id: string;
    sender_id?: string;
    search_content?: string;
    is_read?: boolean;
    is_send?: boolean;
    is_media?: boolean;
    page?: number;
    page_size?: number;
}

// --- Response DTOs ---

export interface MessageResponse {
    id: string;
    channel_id: string;
    sender_id: string;
    sender_username: string;
    type: MessageType;
    content: string;
    created_at: string; // ISO 8601 format
    is_read: boolean;
}
