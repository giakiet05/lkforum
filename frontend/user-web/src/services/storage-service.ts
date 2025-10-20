import type { User } from "../models/user";

export const TOKEN_KEY = "access_token";
export const REFRESH_KEY = "refresh_token";
export const USER_KEY = "user";

// Store tokens to localStorage
export function setAccessToken(token: string): void {
  localStorage.setItem(TOKEN_KEY, token);
}

// Get access token from localStorage
export function getAccessToken(): string | null {
  return localStorage.getItem(TOKEN_KEY);
}

// Clear access token from localStorage
export function clearAccessToken(): void {
  localStorage.removeItem(TOKEN_KEY);
}

// Store refresh token to localStorage
export function setRefreshToken(refreshToken: string): void {
  localStorage.setItem(REFRESH_KEY, refreshToken);
}

// Get refresh token from localStorage
export function getRefreshToken(): string | null {
  return localStorage.getItem(REFRESH_KEY);
}

// Clear refresh token from localStorage
export function clearRefreshToken(): void {
  localStorage.removeItem(REFRESH_KEY);
}

// Store user info to localStorage
export function setUser(user: User): void {
  localStorage.setItem(USER_KEY, JSON.stringify(user));
}

// Get user info from localStorage
export function getUser(): User | null {
  try {
    const user = localStorage.getItem(USER_KEY);
    return user ? (JSON.parse(user) as User) : null;
  } catch {
    return null;
  }
}

// Clear user info from localStorage
export function clearUser(): void {
  localStorage.removeItem(USER_KEY);
}