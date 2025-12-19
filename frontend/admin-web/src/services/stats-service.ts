import { authenticatedFetch, handleApiResponse } from "./api";
import type {
  PlatformOverview,
  UserStats,
  ContentStats,
  GetUserStatsQuery,
  GetContentStatsQuery,
} from "../dtos/stats-dto";

export async function getPlatformOverview(): Promise<PlatformOverview> {
  const response = await authenticatedFetch("/api/admin/stats/overview");
  return handleApiResponse<{ data: PlatformOverview }>(response).then(
    (res) => res.data
  );
}

export async function getUserStats(
  query: GetUserStatsQuery = {}
): Promise<UserStats> {
  const params = new URLSearchParams();
  if (query.period) params.append("period", query.period);

  const url = `/api/admin/stats/users${params.toString() ? `?${params}` : ""}`;
  const response = await authenticatedFetch(url);
  return handleApiResponse<{ data: UserStats }>(response).then(
    (res) => res.data
  );
}

export async function getContentStats(
  query: GetContentStatsQuery = {}
): Promise<ContentStats> {
  const params = new URLSearchParams();
  if (query.period) params.append("period", query.period);

  const url = `/api/admin/stats/content${params.toString() ? `?${params}` : ""}`;
  const response = await authenticatedFetch(url);
  return handleApiResponse<{ data: ContentStats }>(response).then(
    (res) => res.data
  );
}
