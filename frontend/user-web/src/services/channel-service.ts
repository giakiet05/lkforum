import type {
    CreateChannelRequest,
    ChannelResponse,
    GetChannelByUserIDQuery,
    UpdateChannelRequest,
    PaginatedChannelsResponse
} from "../dtos/channel-dto";
import { authenticatedFetch, handleApiResponse } from "./api";
import { USER_KEY } from "../constants/auth-constants";

/**
 * Create a new channel (1-1 conversation)
 */
export async function createChannel(
    member2_id: string,
    member2_username: string,
    member2_avatar: string
): Promise<ChannelResponse> {
    // Get current user info from localStorage
    const userStr = localStorage.getItem(USER_KEY);
    if (!userStr) {
        throw new Error("User not authenticated");
    }
    const currentUser = JSON.parse(userStr);

    const reqBody: CreateChannelRequest = {
        member_1: currentUser.id,
        member_1_username: currentUser.username,
        member_1_avatar: currentUser.avatar_url || "",
        member_2: member2_id,
        member_2_username: member2_username,
        member_2_avatar: member2_avatar
    };

    const res = await authenticatedFetch("/api/channels", {
        method: "POST",
        body: JSON.stringify(reqBody)
    });

    return await handleApiResponse(res);
}

/**
 * Get channel by ID
 */
export async function getChannelById(channelId: string): Promise<ChannelResponse> {
    const res = await authenticatedFetch(`/api/channels/${channelId}`, {
        method: "GET"
    });

    return await handleApiResponse(res);
}

/**
 * Get all channels for current user
 */
export async function getChannelsByUser(
    userId: string,
    page: number = 1,
    pageSize: number = 20
): Promise<PaginatedChannelsResponse> {
    const params = new URLSearchParams({
        user_id: userId,
        page: page.toString(),
        page_size: pageSize.toString()
    });

    const res = await authenticatedFetch(`/api/channels/user?${params.toString()}`, {
        method: "GET"
    });

    return await handleApiResponse(res);
}

/**
 * Get channel between two users
 */
export async function getChannelBetweenUsers(
    user1Id: string,
    user2Id: string
): Promise<ChannelResponse | null> {
    const res = await authenticatedFetch(`/api/channels/between/${user1Id}/${user2Id}`, {
        method: "GET"
    });

    try {
        return await handleApiResponse(res);
    } catch (error) {
        // If channel doesn't exist, return null
        return null;
    }
}

/**
 * Update channel settings (nickname, notifications, etc.)
 */
export async function updateChannel(request: UpdateChannelRequest): Promise<ChannelResponse> {
    const res = await authenticatedFetch("/api/channels", {
        method: "PUT",
        body: JSON.stringify(request)
    });

    return await handleApiResponse(res);
}

/**
 * Delete channel by ID
 */
export async function deleteChannel(channelId: string): Promise<void> {
    const res = await authenticatedFetch(`/api/channels/${channelId}`, {
        method: "DELETE"
    });

    await handleApiResponse(res);
}
