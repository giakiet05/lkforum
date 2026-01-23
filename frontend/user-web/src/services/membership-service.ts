import type {
    CreateMembershipRequest,
    DeleteMembershipRequest,
    MembershipResponse,
    PaginatedMembershipsResponse
} from "../dtos/membership-dto";
import { authenticatedFetch, handleApiResponse } from "./api";

/**
 * Create a new membership (join community)
 */
export async function createMembership(
    userId: string,
    communityId: string
): Promise<MembershipResponse> {
    const reqBody: CreateMembershipRequest = {
        user_id: userId,
        community_id: communityId
    };

    const res = await authenticatedFetch("/api/memberships", {
        method: "POST",
        body: JSON.stringify(reqBody)
    });

    return await handleApiResponse(res);
}

/**
 * Delete membership (leave community)
 */
export async function deleteMembership(
    userId: string,
    communityId: string
): Promise<void> {
    const reqBody: DeleteMembershipRequest = {
        user_id: userId,
        community_id: communityId
    };

    const res = await authenticatedFetch("/api/memberships", {
        method: "DELETE",
        body: JSON.stringify(reqBody)
    });

    await handleApiResponse(res);
}

/**
 * Get all memberships for a user
 */
export async function getMembershipsByUserId(
    userId: string
): Promise<MembershipResponse[]> {
    const res = await authenticatedFetch(`/api/memberships/user/${userId}`, {
        method: "GET"
    });

    return await handleApiResponse(res);
}

/**
 * Get all memberships for a community (with user info)
 */
export async function getMembershipsByCommunityId(
    communityId: string,
    page: number = 1,
    pageSize: number = 20
): Promise<MembershipResponse[]> {
    const params = new URLSearchParams({
        page: page.toString(),
        page_size: pageSize.toString()
    });

    const res = await authenticatedFetch(
        `/api/memberships/community/${communityId}?${params.toString()}`,
        { method: "GET" }
    );

    const data: PaginatedMembershipsResponse = await handleApiResponse(res);
    return data.memberships || [];
}

/**
 * Check if user is member of a community
 */
export async function checkMembership(
    userId: string,
    communityId: string
): Promise<boolean> {
    try {
        const memberships = await getMembershipsByUserId(userId);
        return memberships.some(m => m.community_id === communityId);
    } catch (error) {
        console.error("Failed to check membership:", error);
        return false;
    }
}

/**
 * Kick a member from community (moderator/creator only)
 */
export async function kickMember(
    communityId: string,
    userId: string
): Promise<void> {
    const res = await authenticatedFetch(`/api/memberships/kick/${communityId}/${userId}`, {
        method: "DELETE"
    });

    await handleApiResponse(res);
}

/**
 * Get pending membership requests for a community (moderator/creator only)
 */
export async function getPendingMembers(
    communityId: string,
    page: number = 1,
    pageSize: number = 20
): Promise<PaginatedMembershipsResponse> {
    const params = new URLSearchParams({
        page: page.toString(),
        page_size: pageSize.toString()
    });

    const res = await authenticatedFetch(
        `/api/memberships/community/${communityId}/pending?${params.toString()}`,
        { method: "GET" }
    );

    return await handleApiResponse(res);
}

/**
 * Approve a pending membership request (moderator/creator only)
 */
export async function approveMember(
    communityId: string,
    membershipId: string
): Promise<void> {
    const res = await authenticatedFetch(
        `/api/memberships/community/${communityId}/approve/${membershipId}`,
        { method: "POST" }
    );

    await handleApiResponse(res);
}

/**
 * Reject a pending membership request (moderator/creator only)
 */
export async function rejectMember(
    communityId: string,
    membershipId: string
): Promise<void> {
    const res = await authenticatedFetch(
        `/api/memberships/community/${communityId}/reject/${membershipId}`,
        { method: "POST" }
    );

    await handleApiResponse(res);
}

/**
 * Get approved members for a community (moderator/creator only)
 */
export async function getApprovedMembers(
    communityId: string,
    page: number = 1,
    pageSize: number = 20
): Promise<PaginatedMembershipsResponse> {
    const params = new URLSearchParams({
        page: page.toString(),
        page_size: pageSize.toString()
    });

    const res = await authenticatedFetch(
        `/api/memberships/community/${communityId}/approved?${params.toString()}`,
        { method: "GET" }
    );

    return await handleApiResponse(res);
}
