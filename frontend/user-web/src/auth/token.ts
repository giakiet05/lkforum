import type {RefreshTokenRequest, RefreshTokenResponse} from "../dtos/auth-dto";
import {ACCESS_TOKEN_KEY, REFRESH_TOKEN_KEY} from "../constants/auth-constants";
import type {ApiResponse} from "../dtos/response-dto";

function decodeToken(token: string): any {
    try {
        const payload = token.split(".")[1];
        const decoded = atob(payload);
        return JSON.parse(decoded);
    } catch {
        return null;
    }
}

export function isTokenExpired(token: string): boolean {
    const decoded = decodeToken(token);
    if (!decoded || !decoded.exp) return true;
    const now = Math.floor(Date.now() / 1000); // seconds
    return decoded.exp < now;
}

export async function getValidAccessToken(): Promise<string | null> {
    const token = localStorage.getItem(ACCESS_TOKEN_KEY);
    console.log("🔍 Checking access token:", token ? "Found" : "Missing");
    
    if (token) {
        const decoded = decodeToken(token);
        console.log("🔐 Token payload:", decoded);
        const isExpired = isTokenExpired(token);
        console.log(`⏰ Token expired: ${isExpired}`, decoded?.exp ? `(expires at: ${new Date(decoded.exp * 1000).toISOString()})` : '');
        
        if (!isExpired) {
            console.log("✅ Access token is valid");
            return token;
        }
    }

    console.log("⚠️ Access token expired or missing, trying refresh...");
    const refreshToken = localStorage.getItem(REFRESH_TOKEN_KEY);
    if (!refreshToken) {
        console.log("❌ No refresh token found");
        return null;
    }

    const reqBody: RefreshTokenRequest = { refresh_token: refreshToken };

    try {
        const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "";
        const url = API_BASE_URL + "/api/auth/refresh";

        const res = await fetch(url, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(reqBody),
        });

        if (!res.ok) {
            console.error("Refresh token failed with status:", res.status);
            return null;
        }

        const apiRes: ApiResponse = await res.json();

        // Check API response success
        if (!apiRes.success) {
            console.error("Refresh token failed:", apiRes.message);
            return null;
        }

        const data: RefreshTokenResponse = apiRes.data;
        localStorage.setItem(ACCESS_TOKEN_KEY, data.access_token);
        localStorage.setItem(REFRESH_TOKEN_KEY, data.refresh_token);

        // Dispatch event để App.svelte có thể update authStore
        window.dispatchEvent(new CustomEvent('auth:refreshed'));

        return data.access_token;
    } catch (err) {
        // Network error hoặc JSON parse error
        console.error("Refresh token error:", err);
        return null;
    }
}