import type { LoginRequest, LoginResponse, RegisterDto } from "../dtos/auth-dto";
import type {
    RefreshTokenRequest,
    RefreshTokenResponse,
} from "../dtos/auth-dto";
import type { User } from "../models/user";
import { setAuth, clearAuth } from "../stores/auth-store";
import {
    setAccessToken,
    getAccessToken,
    clearAccessToken,
    setRefreshToken,
    getRefreshToken,
    clearRefreshToken,
    setUser,
    getUser,
    clearUser
} from "./storage-service";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "";

// Check if user is logged in (has a valid access token)
export function isLoggedIn(): boolean {
  return !!getAccessToken();
}

// Decode JWT token payload
function decodeToken(token: string): any {
  try {
    const payload = token.split(".")[1];
    const decoded = atob(payload);
    return JSON.parse(decoded);
  } catch {
    return null;
  }
}

// Check if token is expired
export function isTokenExpired(token: string): boolean {
  const decoded = decodeToken(token);
  if (!decoded || !decoded.exp) return true;
  const now = Math.floor(Date.now() / 1000); // seconds
  return decoded.exp < now;
}

// Get a valid access token, refresh if expired, return null if cannot refresh
export async function getValidAccessToken(): Promise<string | null> {
  const token = getAccessToken();
  if (token && !isTokenExpired(token)) {
    return token;
  }

  // If token is expired or missing, try to refresh
  const refreshToken = getRefreshToken();
  if (!refreshToken) {
    return null;
  }

  const reqBody: RefreshTokenRequest = { refresh_token: refreshToken };

  try {
    const res = await fetch("/api/auth/refresh", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(reqBody),
    });

    if (!res.ok) return null;

    const data: RefreshTokenResponse = await res.json();
    setAccessToken(data.access_token); // Save new access token
    setRefreshToken(data.refresh_token); // Save new refresh token
    return data.access_token;
  } catch (err) {
    console.error("Refresh token error:", err);
    return null;
  }
}

export async function register(data: RegisterDto): Promise<LoginResponse> {
  const res = await fetch(`${API_BASE_URL}/api/auth/register`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });

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

  const response: LoginResponse = await res.json();
  setAccessToken(response.access_token);
  setRefreshToken(response.refresh_token);
  setUser(response.user);
  setAuth(response.user, response.access_token);
  return response;
}

export async function login(credentials: LoginRequest): Promise<LoginResponse> {
    console.log('Login request:', credentials);
    const res = await fetch(`${API_BASE_URL}/api/auth/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(credentials),
    });

    console.log('Login response status:', res.status);
    
    if (!res.ok) {
        let errObj: any = {};
        try {
            errObj = await res.json();
            console.log('Login error response:', errObj);
        } catch (e) {
            try {
                const text = await res.text();
                console.log('Login error text:', text);
                errObj = { error: text || `HTTP ${res.status}` };
            } catch {
                errObj = { error: `HTTP ${res.status}` };
            }
        }
        throw errObj.error || errObj.message || "Unknown error";
    }

    const response: LoginResponse = await res.json();
    console.log('Login success:', response);
    setAccessToken(response.access_token);
    setRefreshToken(response.refresh_token);
    setUser(response.user);
    setAuth(response.user, response.access_token);
    return response;
}

export function logout(): void {
  clearAccessToken();
  clearRefreshToken();
  clearUser();
  clearAuth();
  window.location.reload();
}
