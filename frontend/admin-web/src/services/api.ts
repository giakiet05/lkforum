import { tokenManager } from "../auth/token";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://localhost:8081";

let isRefreshing = false;
let refreshSubscribers: ((token: string) => void)[] = [];

function onAccessTokenRefreshed(token: string) {
  refreshSubscribers.forEach((callback) => callback(token));
  refreshSubscribers = [];
}

function addRefreshSubscriber(callback: (token: string) => void) {
  refreshSubscribers.push(callback);
}

async function refreshAccessToken(): Promise<string> {
  console.log("[API] Attempting to refresh access token");
  const refreshToken = tokenManager.getRefreshToken();
  if (!refreshToken) {
    console.log("[API] No refresh token available, redirecting to login");
    throw new Error("No refresh token available");
  }

  const response = await fetch(`${API_BASE_URL}/api/admin/auth/refresh`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ refresh_token: refreshToken }),
  });

  if (!response.ok) {
    console.log("[API] Refresh token failed, clearing tokens and redirecting");
    tokenManager.clearTokens();
    window.location.href = "/#/login";
    throw new Error("Failed to refresh token");
  }

  const data = await response.json();
  const newAccessToken = data.data.access_token;
  tokenManager.setAccessToken(newAccessToken);
  return newAccessToken;
}

export async function authenticatedFetch(
  url: string,
  options: RequestInit = {}
): Promise<Response> {
  console.log(`[API] authenticatedFetch: ${options.method || 'GET'} ${url}`);
  let accessToken = tokenManager.getAccessToken();

  if (!accessToken) {
    throw new Error("No access token available");
  }

  const makeRequest = async (token: string) => {
    return fetch(`${API_BASE_URL}${url}`, {
      ...options,
      headers: {
        ...options.headers,
        Authorization: `Bearer ${token}`,
      },
    });
  };

  let response = await makeRequest(accessToken);

  if (response.status === 401) {
    if (!isRefreshing) {
      isRefreshing = true;
      try {
        const newToken = await refreshAccessToken();
        isRefreshing = false;
        onAccessTokenRefreshed(newToken);
        response = await makeRequest(newToken);
      } catch (error) {
        isRefreshing = false;
        throw error;
      }
    } else {
      const newToken = await new Promise<string>((resolve) => {
        addRefreshSubscriber((token: string) => {
          resolve(token);
        });
      });
      response = await makeRequest(newToken);
    }
  }

  return response;
}

export async function publicFetch(
  url: string,
  options: RequestInit = {}
): Promise<Response> {
  return fetch(`${API_BASE_URL}${url}`, options);
}

export async function handleApiResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    const error = await response.json().catch(() => ({
      message: "An error occurred",
    }));
    throw new Error(error.message || "API request failed");
  }

  const data = await response.json();
  return data.data;
}
