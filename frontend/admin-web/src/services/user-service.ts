import { authenticatedFetch, publicFetch, handleApiResponse } from "./api";
import type { PaginatedUsersResponse } from "../dtos/user-dto";

// Note: Backend admin routes are not registered yet
// Using public communities API for now as placeholder
export async function getUsers(params?: {
  page?: number;
  limit?: number;
}): Promise<PaginatedUsersResponse> {
  // TODO: Backend needs to register admin routes in init.go
  // For now, return empty data since /api/admin/users doesn't exist
  return {
    users: [],
    total: 0,
    page: params?.page || 1,
    page_size: params?.limit || 10,
  };
}

export async function banUser(userId: string): Promise<void> {
  const res = await authenticatedFetch(`/api/admin/users/${userId}/ban`, {
    method: "POST",
  });

  await handleApiResponse(res);
}

export async function unbanUser(userId: string): Promise<void> {
  const res = await authenticatedFetch(`/api/admin/users/${userId}/unban`, {
    method: "POST",
  });

  await handleApiResponse(res);
}
