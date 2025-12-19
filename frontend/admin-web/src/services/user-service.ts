import { authenticatedFetch, handleApiResponse } from "./api";
import type { PaginatedUsersResponse } from "../dtos/user-dto";

export async function getUsers(params?: {
  page?: number;
  limit?: number;
  is_banned?: boolean;
  role?: string;
}): Promise<PaginatedUsersResponse> {
  const queryParams = new URLSearchParams();
  if (params?.page) queryParams.append("page", String(params.page));
  if (params?.limit) queryParams.append("limit", String(params.limit));
  if (params?.is_banned !== undefined)
    queryParams.append("is_banned", String(params.is_banned));
  if (params?.role) queryParams.append("role", params.role);

  const res = await authenticatedFetch(
    `/api/admin/users?${queryParams.toString()}`
  );
  return await handleApiResponse<PaginatedUsersResponse>(res);
}

export async function banUser(
  userId: string,
  reason?: string,
  banUntil?: string
): Promise<void> {
  const body: { reason?: string; ban_until?: string } = {};
  if (reason) body.reason = reason;
  if (banUntil) body.ban_until = banUntil;

  const res = await authenticatedFetch(`/api/admin/users/${userId}/ban`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(body),
  });

  await handleApiResponse(res);
}

export async function unbanUser(userId: string): Promise<void> {
  const res = await authenticatedFetch(`/api/admin/users/${userId}/unban`, {
    method: "POST",
  });

  await handleApiResponse(res);
}

export async function deleteUser(userId: string): Promise<void> {
  const res = await authenticatedFetch(`/api/admin/users/${userId}`, {
    method: "DELETE",
  });

  await handleApiResponse(res);
}

export async function restoreUser(userId: string): Promise<void> {
  const res = await authenticatedFetch(`/api/admin/users/${userId}/restore`, {
    method: "POST",
  });

  await handleApiResponse(res);
}
