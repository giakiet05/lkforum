import { authenticatedFetch, handleApiResponse } from "./api";
import type { PaginatedCommunitiesResponse } from "../dtos/community-dto";

export async function getCommunities(params?: {
  name?: string;
  page?: number;
  limit?: number;
  is_banned?: boolean;
}): Promise<PaginatedCommunitiesResponse> {
  const queryParams = new URLSearchParams();
  if (params?.name) queryParams.append("name", params.name);
  if (params?.page) queryParams.append("page", String(params.page));
  if (params?.limit) queryParams.append("limit", String(params.limit));
  if (params?.is_banned !== undefined)
    queryParams.append("is_banned", String(params.is_banned));

  const res = await authenticatedFetch(
    `/api/admin/communities?${queryParams.toString()}`
  );

  return await handleApiResponse<PaginatedCommunitiesResponse>(res);
}

export async function banCommunity(
  communityId: string,
  reason: string
): Promise<void> {
  const res = await authenticatedFetch(
    `/api/admin/communities/${communityId}/ban`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ reason }),
    }
  );

  await handleApiResponse(res);
}

export async function unbanCommunity(communityId: string): Promise<void> {
  const res = await authenticatedFetch(
    `/api/admin/communities/${communityId}/unban`,
    {
      method: "POST",
    }
  );

  await handleApiResponse(res);
}
