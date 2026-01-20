// --- Request DTOs ---

export interface CreateMembershipRequest {
    user_id: string;
    community_id: string;
}

export interface DeleteMembershipRequest {
    user_id: string;
    community_id: string;
}

// --- Response DTOs ---

export interface MembershipResponse {
    id: string;
    user_id: string;
    community_id: string;
}
