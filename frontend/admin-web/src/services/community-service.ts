import { authenticatedFetch, publicFetch, handleApiResponse } from "./api";
import type { PaginatedCommunitiesResponse } from "../dtos/community-dto";

export async function getCommunities(params?: {
  name?: string;
  page?: number;
  limit?: number;
}): Promise<PaginatedCommunitiesResponse> {
  const queryParams = new URLSearchParams();
  if (params?.name) queryParams.append("name", params.name);
  if (params?.page) queryParams.append("page", String(params.page));
  if (params?.limit) queryParams.append("limit", String(params.limit));

  const res = await publicFetch(`/api/communities/filter?${queryParams.toString()}`, {
    method: "GET",
  });

  return await handleApiResponse(res);
}

export async function deleteCommunity(communityId: string): Promise<void> {
  const res = await authenticatedFetch(`/api/communities/${communityId}`, {
    method: "DELETE",
  });

  await handleApiResponse(res);
}
