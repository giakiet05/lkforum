import { writable } from "svelte/store";
import type { User } from "../models/user";
import { getAccessToken, getUser, isTokenExpired } from "../services/auth-service";

interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
}

function getInitialAuthState(): AuthState {
  const token = getAccessToken();
  const user = getUser();
  const isAuthenticated = !!token && !isTokenExpired(token);
  return {
    user,
    token,
    isAuthenticated,
  };
}

export const authStore = writable<AuthState>(getInitialAuthState());

export function setAuth(user: User, token: string) {
  authStore.set({ user, token, isAuthenticated: true });
}
export function clearAuth() {
  authStore.set({ user: null, token: null, isAuthenticated: false });
}
