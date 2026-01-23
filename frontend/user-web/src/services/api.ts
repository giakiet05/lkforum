import { getValidAccessToken } from "../auth/token";
import type { ApiResponse } from "../dtos/response-dto";
import { ApiErrorCode } from "../errors/error-codes";
import { ApiError } from "../errors/api-error";


export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "";

/**
 * Handles the response from the Fetch API.
 * Checks for network errors and API-level errors (based on `success` field).
 * @param res - The Response object from a fetch call.
 * @returns The `data` field from the API response.
 * @throws {ApiError} If the response is not ok, `success` is false, or JSON parse fails.
 */
export async function handleApiResponse(res: Response): Promise<any> {
    try {
        const json: ApiResponse = await res.json();
        if (!res.ok || !json.success) {
            // Log full error response for debugging
            console.error("❌ API Error Response:", {
                status: res.status,
                statusText: res.statusText,
                success: json.success,
                message: json.message,
                error_code: json.error_code,
                data: json.data,
                fullResponse: json
            });
            // Use the message from the API, or a default one
            const message = json.message || `Request failed with status ${res.status}`;
            throw new ApiError(message, json.error_code as ApiErrorCode, json.data);
        }
        return json.data;
    } catch (error) {
        // JSON parse error hoặc ApiError đã throw ở trên
        if (error instanceof ApiError) {
            throw error; // Re-throw ApiError
        }
        // JSON parse error - likely 500 error with empty body
        console.error("Failed to parse API response:", error);
        const statusMessage = res.status === 500 
            ? "Server error occurred. Please try with smaller images or contact support."
            : "Phản hồi từ server không hợp lệ.";
        throw new ApiError(
            statusMessage,
            ApiErrorCode.INTERNAL_ERROR
        );
    }
}

/**
 * A wrapper for the standard Fetch API that includes base URL and default headers.
 * This function is for public endpoints that do not require authentication.
 * @param path - The API endpoint path (e.g., "/api/auth/login").
 * @param options - Standard RequestInit options.
 * @returns The raw Response object.
 * @throws {ApiError} If network error occurs (mất mạng, timeout, CORS, DNS fail).
 */
export async function publicFetch(path: string, options: RequestInit = {}): Promise<Response> {
    const url = path.startsWith("http") ? path : API_BASE_URL + path;

    const headers: Record<string, string> = {
        "Content-Type": "application/json",
        ...((options.headers as Record<string, string>) || {}),
    };

    try {
        return await fetch(url, { ...options, headers });
    } catch (error) {
        // Network errors: mất mạng, timeout, CORS, DNS fail
        console.error("Network error in publicFetch:", error);
        throw new ApiError(
            "Không thể kết nối đến server. Vui lòng kiểm tra kết nối mạng.",
            ApiErrorCode.SERVICE_UNAVAILABLE
        );
    }
}


/**
 * Centralized API request function for authenticated calls.
 * This function automatically adds the Authorization header and handles token refreshing.
 * @param path - API endpoint path.
 * @param options - Fetch API options.
 * @returns The raw Response object.
 * @throws {ApiError} If the request fails, token cannot be refreshed, or network error occurs.
 */

let isUnauthorizedEventDispatched = false;

export async function authenticatedFetch(path: string, options: RequestInit = {}): Promise<Response> {
    const url = path.startsWith("http") ? path : API_BASE_URL + path;
    const accessToken = await getValidAccessToken();
    
    console.log("🔑 Access token:", accessToken ? "✅ Found" : "❌ Missing");
    
    if (!accessToken) {
        // Dispatch event only once to avoid infinite loops
        if (!isUnauthorizedEventDispatched) {
            isUnauthorizedEventDispatched = true;
            window.dispatchEvent(new CustomEvent('auth:unauthorized'));
            // Reset after 2 seconds to allow future logouts
            setTimeout(() => { isUnauthorizedEventDispatched = false; }, 2000);
        }
        throw new ApiError("Phiên đăng nhập đã hết hạn. Vui lòng đăng nhập lại.", ApiErrorCode.FORBIDDEN);
    }

    const isFormData = options.body instanceof FormData;

    const headers: Record<string, string> = {
        ...((options.headers as Record<string, string>) || {}),
        Authorization: `Bearer ${accessToken}`,
    };

    // Do not set Content-Type for FormData, the browser does it.
    if (!isFormData) {
        headers["Content-Type"] = "application/json";
    }

    console.log("📤 Request headers:", { ...headers, Authorization: `Bearer ${accessToken.substring(0, 20)}...` });
    console.log("📤 Request URL:", url);
    console.log("📤 Request method:", options.method || "GET");
    if (!isFormData && options.body) {
        console.log("📤 Request body:", options.body);
    }

    try {
        return await fetch(url, { ...options, headers });
    } catch (error) {
        // Network errors: mất mạng, timeout, CORS, DNS fail
        console.error("Network error in authenticatedFetch:", error);
        throw new ApiError(
            "Không thể kết nối đến server. Vui lòng kiểm tra kết nối mạng.",
            ApiErrorCode.SERVICE_UNAVAILABLE
        );
    }
}

/**
 * A wrapper for endpoints that can work with or without authentication.
 * If token is available, it will be sent. Otherwise, request proceeds without auth.
 * Useful for endpoints like getPosts where logged-in users get personalized results.
 * @param path - API endpoint path.
 * @param options - Fetch API options.
 * @returns The raw Response object.
 */
export async function optionalAuthFetch(path: string, options: RequestInit = {}): Promise<Response> {
    const url = path.startsWith("http") ? path : API_BASE_URL + path;
    const accessToken = await getValidAccessToken();
    
    const headers: Record<string, string> = {
        "Content-Type": "application/json",
        ...((options.headers as Record<string, string>) || {}),
    };

    // Add auth header if token is available
    if (accessToken) {
        headers["Authorization"] = `Bearer ${accessToken}`;
    }

    try {
        return await fetch(url, { ...options, headers });
    } catch (error) {
        // Network errors: mất mạng, timeout, CORS, DNS fail
        console.error("Network error in optionalAuthFetch:", error);
        throw new ApiError(
            "Không thể kết nối đến server. Vui lòng kiểm tra kết nối mạng.",
            ApiErrorCode.SERVICE_UNAVAILABLE
        );
    }
}
