import { publicFetch, handleApiResponse } from "./api";
import type { LoginRequest, LoginResponse } from "../dtos/auth-dto";
import { tokenManager } from "../auth/token";

export async function login(credentials: LoginRequest): Promise<LoginResponse> {
  const res = await publicFetch("/api/admin/auth/login", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(credentials),
  });

  const data = await handleApiResponse<LoginResponse>(res);
  tokenManager.setAccessToken(data.access_token);
  tokenManager.setRefreshToken(data.refresh_token);
  return data;
}

export async function logout(): Promise<void> {
  tokenManager.clearTokens();
}
