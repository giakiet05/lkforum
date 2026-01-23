import { authenticatedFetch, handleApiResponse } from "./api";
import type { NotificationResponse, PaginatedNotificationsResponse } from "../dtos/notification-dto";

/**
 * Get user notifications
 */
export async function getNotifications(params?: {
    page?: number;
    pageSize?: number;
}): Promise<PaginatedNotificationsResponse> {
    const queryParams = new URLSearchParams();
    if (params) {
        Object.entries(params).forEach(([key, value]) => {
            if (value !== undefined && value !== null) {
                queryParams.append(key, String(value));
            }
        });
    }
    
    const res = await authenticatedFetch(`/api/notifications?${queryParams.toString()}`, {
        method: "GET",
    });
    
    return await handleApiResponse(res);
}

/**
 * Mark notification as read
 */
export async function markNotificationAsRead(notificationId: string): Promise<void> {
    const res = await authenticatedFetch(`/api/notifications/${notificationId}/read`, {
        method: "PUT",
    });
    
    return await handleApiResponse(res);
}

/**
 * Delete notification
 */
export async function deleteNotification(notificationId: string): Promise<void> {
    const res = await authenticatedFetch(`/api/notifications/${notificationId}`, {
        method: "DELETE",
    });
    
    return await handleApiResponse(res);
}

/**
 * Mark all notifications as read
 */
export async function markAllAsRead(): Promise<void> {
    console.log("🔔 [notification-service] Calling PUT /api/notifications/read-all");
    const res = await authenticatedFetch(`/api/notifications/read-all`, {
        method: "PUT",
    });
    console.log("🔔 [notification-service] Response status:", res.status);
    const result = await handleApiResponse(res);
    console.log("🔔 [notification-service] Response data:", result);
    return result;
}
