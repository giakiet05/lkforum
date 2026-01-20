import {writable} from "svelte/store";
import type {UserResponse} from "../dtos/user-dto";
import {ACCESS_TOKEN_KEY, REFRESH_TOKEN_KEY, USER_KEY} from "../constants/auth-constants";
import {isTokenExpired} from "../auth/token";

interface AuthState {
    user: UserResponse | null;
    isAuthenticated: boolean;
}

export function getInitialAuthState(): AuthState {
    const access_token = localStorage.getItem(ACCESS_TOKEN_KEY);
    const user: UserResponse | null = localStorage.getItem(USER_KEY) ? JSON.parse(localStorage.getItem(USER_KEY)!) : null;
    const isAuthenticated = !!access_token && !isTokenExpired(access_token);
    return {
        user,
        isAuthenticated,
    };
}

export const authStore = writable<AuthState>(getInitialAuthState());

export function setAuth(user: UserResponse) {
    authStore.set({user, isAuthenticated: true});
}

export function clearAuth() {
    authStore.set({user: null, isAuthenticated: false});
}
