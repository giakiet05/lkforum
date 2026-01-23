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

export interface MembershipUserInfo {
    id: string;
    username: string;
    avatar?: {
        url: string;
        public_id?: string;
    };
}

export interface MembershipResponse {
    id: string;
    user_id: string;
    community_id: string;
    status?: 'pending' | 'approved' | 'rejected';
    created_at?: string;
    user?: MembershipUserInfo;
}

export interface PaginatedMembershipsResponse {
    memberships: MembershipResponse[];
    pagination: {
        total: number;
        page: number;
    };
}
