import { getValidAccessToken, logout } from "./auth-service";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "";

/**
 * Centralized API request function for authenticated calls.
 * Handles token refresh, error management, and header setup.
 * @param path - API endpoint path (relative or absolute)
 * @param options - Fetch API options
 * @returns Response data or throws error
 */
export async function apiFetch(
  path: string,
  options: RequestInit = {}
): Promise<any> {
  // Get a valid access token (refresh if needed)
  const accessToken = await getValidAccessToken();
  if (!accessToken) {
    logout();
    throw new Error("Not authenticated");
  }

  // Check if body is FormData for correct Content-Type
  const isFormData = options.body instanceof FormData;

  // Normalize headers to a plain object
  const headers: Record<string, string> = {
    ...((options.headers as Record<string, string>) || {}),
    Authorization: `Bearer ${accessToken}`,
    ...(!isFormData ? { "Content-Type": "application/json" } : {}),
  };
  options.headers = headers;

  // Build full URL
  const url =
    path.startsWith("http://") || path.startsWith("https://")
      ? path
      : API_BASE_URL + path;

  // Make the request
  const res = await fetch(url, options);

  // Handle unauthorized (token expired)
  if (res.status === 401) {
    logout();
    throw new Error("Token expired");
  }

  // Handle other errors
  if (!res.ok) {
    let errObj: any = {};
    try {
      errObj = await res.json();
    } catch (e) {
      try {
        const text = await res.text();
        errObj = { error: text || `HTTP ${res.status}` };
      } catch {
        errObj = { error: `HTTP ${res.status}` };
      }
    }
    throw errObj.error || "Unknown error";
  }

  // Return parsed response data
  try {
    return await res.json();
  } catch {
    return res;
  }
}
