import { authenticatedFetch, publicFetch, handleApiResponse } from "./api";
import type { 
    CreateCommunityRequest, 
    UpdateCommunityRequest,
    CommunityResponse,
    PaginatedCommunitiesResponse,
    AddModeratorRequest,
    RemoveModeratorRequest,
    BanUserRequest,
    UnbanUserRequest,
    UserResponse
} from "../dtos/community-dto";

/**
 * Get communities by user ID (joined communities)
 */
export async function getCommunitiesByUserId(userId: string): Promise<CommunityResponse[]> {
    const res = await authenticatedFetch(`/api/communities/user/${userId}`, {
        method: "GET",
    });
    
    return await handleApiResponse(res);
}

/**
 * Get all communities with optional filters
 */
export async function getCommunities(params?: {
    name?: string;
    description?: string;
    create_from?: string;
    page?: number;
    limit?: number;
}): Promise<PaginatedCommunitiesResponse> {
    const queryParams = new URLSearchParams();
    if (params) {
        Object.entries(params).forEach(([key, value]) => {
            if (value !== undefined && value !== null) {
                queryParams.append(key, String(value));
            }
        });
    }
    
    const res = await publicFetch(`/api/communities/filter?${queryParams.toString()}`, {
        method: "GET",
    });
    
    return await handleApiResponse(res);
}

/**
 * Get a single community by ID
 */
export async function getCommunityById(communityId: string): Promise<CommunityResponse> {
    const res = await publicFetch(`/api/communities/${communityId}`, {
        method: "GET",
    });
    
    return await handleApiResponse(res);
}

/**
 * Get a single community by name
 */
export async function getCommunityByName(name: string): Promise<CommunityResponse> {
    const res = await publicFetch(`/api/communities/name/${name}`, {
        method: "GET",
    });
    
    return await handleApiResponse(res);
}

/**
 * Create a new community (requires authentication)
 */
export async function createCommunity(data: CreateCommunityRequest): Promise<CommunityResponse> {
    const res = await authenticatedFetch("/api/communities", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
    });
    
    return await handleApiResponse(res);
}

/**
 * Update an existing community (requires authentication)
 */
export async function updateCommunity(data: UpdateCommunityRequest): Promise<CommunityResponse> {
    const res = await authenticatedFetch("/api/communities", {
        method: "PUT",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
    });
    
    return await handleApiResponse(res);
}

/**
 * Add moderators to a community (requires authentication)
 */
export async function addModerators(data: AddModeratorRequest): Promise<CommunityResponse> {
    const res = await authenticatedFetch("/api/communities/add_moderator", {
        method: "PUT",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
    });
    
    return await handleApiResponse(res);
}

/**
 * Remove moderators from a community (requires authentication)
 */
export async function removeModerators(data: RemoveModeratorRequest): Promise<CommunityResponse> {
    const res = await authenticatedFetch("/api/communities/remove_moderator", {
        method: "PUT",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
    });
    
    return await handleApiResponse(res);
}

/**
 * Activate moderator status (accept moderator invitation)
 */
export async function activateModerator(communityId: string): Promise<void> {
    const res = await authenticatedFetch(`/api/communities/activate_moderator/${communityId}`, {
        method: "PUT",
    });
    
    return await handleApiResponse(res);
}

/**
 * Delete a community by ID (requires authentication)
 */
export async function deleteCommunity(communityId: string): Promise<void> {
    const res = await authenticatedFetch(`/api/communities/${communityId}`, {
        method: "DELETE",
    });
    
    await handleApiResponse(res);
}

/**
 * Ban or mute a user from a community (requires authentication)
 */
export async function banUser(data: BanUserRequest): Promise<void> {
    const res = await authenticatedFetch("/api/communities/ban/user", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
    });
    
    await handleApiResponse(res);
}

/**
 * Unban a user from a community (requires authentication)
 */
export async function unbanUser(data: UnbanUserRequest): Promise<void> {
    const res = await authenticatedFetch("/api/communities/unban/user", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
    });
    
    await handleApiResponse(res);
}

/**
 * Unmute a user from a community (requires authentication)
 */
export async function unmuteUser(data: UnbanUserRequest): Promise<void> {
    const res = await authenticatedFetch("/api/communities/unmute/user", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
    });
    
    await handleApiResponse(res);
}

/**
 * Get banned or muted users in a community (requires authentication)
 */
export async function getBannedUsers(communityId: string, banType: "banned" | "muted"): Promise<UserResponse[]> {
    const res = await authenticatedFetch(`/api/communities/banned_user?community_id=${communityId}&ban_type=${banType}`, {
        method: "GET",
    });
    
    return await handleApiResponse(res);
}

/**
 * Get pending posts in a community (requires moderator authentication)
 */
export async function getPendingPosts(communityId: string, page: number = 1, pageSize: number = 10): Promise<any> {
    const res = await authenticatedFetch(`/api/communities/${communityId}/posts/pending?page=${page}&page_size=${pageSize}`, {
        method: "GET",
    });
    
    return await handleApiResponse(res);
}

/**
 * Get edited posts in a community (requires moderator authentication)
 */
export async function getEditedPosts(communityId: string, page: number = 1, pageSize: number = 10): Promise<any> {
    const res = await authenticatedFetch(`/api/communities/${communityId}/posts/edited?page=${page}&page_size=${pageSize}`, {
        method: "GET",
    });
    
    return await handleApiResponse(res);
}

/**
 * Moderate a post (approve or reject) (requires moderator authentication)
 */
export async function moderatePost(communityId: string, postId: string, approve: boolean, reason?: string): Promise<void> {
    const res = await authenticatedFetch(`/api/communities/${communityId}/posts/${postId}/moderate`, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ approve, reason }),
    });
    
    await handleApiResponse(res);
}
