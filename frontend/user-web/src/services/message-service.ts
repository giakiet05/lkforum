import type {
    MessageResponse,
    GetMessageFilterQuery
} from "../dtos/message-dto";
import { authenticatedFetch, handleApiResponse } from "./api";

/**
 * Get message by ID
 */
export async function getMessageById(messageId: string): Promise<MessageResponse> {
    const res = await authenticatedFetch(`/api/messages/${messageId}`, {
        method: "GET"
    });

    return await handleApiResponse(res);
}

/**
 * Get messages with filters (channel_id, search, read status, etc.)
 */
export async function getMessages(query: GetMessageFilterQuery): Promise<MessageResponse[]> {
    const params = new URLSearchParams();
    
    // Add all query parameters
    if (query.channel_id) params.append("channel_id", query.channel_id);
    if (query.sender_id) params.append("sender_id", query.sender_id);
    if (query.search_content) params.append("search_content", query.search_content);
    if (query.is_read !== undefined) params.append("is_read", query.is_read.toString());
    if (query.is_send !== undefined) params.append("is_send", query.is_send.toString());
    if (query.is_media !== undefined) params.append("is_media", query.is_media.toString());
    if (query.page) params.append("page", query.page.toString());
    if (query.page_size) params.append("page_size", query.page_size.toString());

    const res = await authenticatedFetch(`/api/messages/filter?${params.toString()}`, {
        method: "GET"
    });

    const response = await handleApiResponse(res);
    
    // Backend returns { messages: [], pagination: {} }, we only need messages array
    return response.messages || [];
}

/**
 * Delete message by ID
 */
export async function deleteMessage(messageId: string): Promise<void> {
    const res = await authenticatedFetch(`/api/messages/${messageId}`, {
        method: "DELETE"
    });

    await handleApiResponse(res);
}

/**
 * Mark all messages in a channel as read
 */
export async function markChannelMessagesAsRead(channelId: string): Promise<void> {
    const res = await authenticatedFetch(`/api/messages/channels/${channelId}/read`, {
        method: "PUT"
    });

    await handleApiResponse(res);
}
