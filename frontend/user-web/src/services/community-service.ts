import { authenticatedFetch, publicFetch, handleApiResponse } from "./api";
import type { 
    CreateCommunityRequest, 
    UpdateCommunityRequest,
    CommunityResponse,
    AddModeratorRequest,
    RemoveModeratorRequest
} from "../dtos/community-dto";

/**
 * Get all communities with optional filters
 */
export async function getCommunities(params?: {
    name?: string;
    description?: string;
    create_from?: string;
    page?: number;
    limit?: number;
}): Promise<CommunityResponse[]> {
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
 * Delete a community by ID (requires authentication)
 */
export async function deleteCommunity(communityId: string): Promise<void> {
    const res = await authenticatedFetch(`/api/communities/${communityId}`, {
        method: "DELETE",
    });
    
    await handleApiResponse(res);
}
