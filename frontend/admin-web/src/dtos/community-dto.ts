export interface CommunityResponse {
  id: string;
  name: string;
  description?: string;
  avatar?: string;
  member_count: number;
  created_at: string;
  is_banned?: boolean;
}

export interface PaginatedCommunitiesResponse {
  communities: CommunityResponse[];
  total: number;
  page: number;
  page_size: number;
}
