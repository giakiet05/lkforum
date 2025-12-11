import type { Pagination } from './pagination-dto';

// --- Response DTOs ---

export type NotificationType =
    | "comment"
    | "upvote"
    | "mention"
    | "message"
    | "system";

export interface NotificationResponse {
    id: string;
    type: NotificationType;
    message: string;
    link: string;
    is_read: boolean;
    created_at: string; // ISO 8601 format
    metadata?: {
        community_id?: string;
        [key: string]: any;
    };
}

export interface PaginatedNotificationsResponse {
    notifications: NotificationResponse[];
    pagination: Pagination;
}
