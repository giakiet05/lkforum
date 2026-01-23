import type { Pagination } from './pagination-dto';

// --- Response DTOs ---

// Backend notification types (from model/notification.go)
export type NotificationType =
    | "comment"
    | "like"         // Backend uses "like" for upvotes
    | "follow"
    | "mention"
    | "new_message"  // Backend uses "new_message" for messages
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
