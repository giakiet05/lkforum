import type {
    UserResponse,
    UserProfileUpdateRequest,
    ChangePasswordRequest,
} from "../dtos/user-dto";
import type {
    SettingsResponse,
    UpdateSettingsRequest,
} from "../dtos/user-settings-dto";
import { authenticatedFetch, publicFetch, handleApiResponse } from "./api";

// --- Profile Management ---

/**
 * Get current user's profile
 */
export async function getMyProfile(): Promise<UserResponse> {
    const res = await authenticatedFetch("/api/users/me", {
        method: "GET",
    });

    return await handleApiResponse(res);
}

/**
 * Get user profile by username (public)
 */
export async function getUserByUsername(username: string): Promise<UserResponse> {
    const res = await publicFetch(`/api/users/profile/${username}`, {
        method: "GET",
    });

    return await handleApiResponse(res);
}

/**
 * Update current user's profile
 */
export async function updateProfile(data: UserProfileUpdateRequest): Promise<UserResponse> {
    const res = await authenticatedFetch("/api/users/me/profile", {
        method: "PUT",
        body: JSON.stringify(data),
    });

    return await handleApiResponse(res);
}

/**
 * Change password
 */
export async function changePassword(data: ChangePasswordRequest): Promise<void> {
    const res = await authenticatedFetch("/api/users/me/password", {
        method: "PUT",
        body: JSON.stringify(data),
    });

    await handleApiResponse(res);
}

// --- Avatar & Cover Management ---

/**
 * Upload avatar
 */
export async function uploadAvatar(file: File): Promise<UserResponse> {
    const formData = new FormData();
    formData.append("avatar", file);

    const res = await authenticatedFetch("/api/users/me/avatar", {
        method: "POST",
        body: formData,
        // Don't set Content-Type - browser will set it with boundary for FormData
    });

    return await handleApiResponse(res);
}

/**
 * Delete avatar
 */
export async function deleteAvatar(): Promise<UserResponse> {
    const res = await authenticatedFetch("/api/users/me/avatar", {
        method: "DELETE",
    });

    return await handleApiResponse(res);
}

/**
 * Upload cover image
 */
export async function uploadCover(file: File): Promise<UserResponse> {
    const formData = new FormData();
    formData.append("cover", file);

    const res = await authenticatedFetch("/api/users/me/cover", {
        method: "POST",
        body: formData,
    });

    return await handleApiResponse(res);
}

/**
 * Delete cover image
 */
export async function deleteCover(): Promise<UserResponse> {
    const res = await authenticatedFetch("/api/users/me/cover", {
        method: "DELETE",
    });

    return await handleApiResponse(res);
}

// --- Settings Management ---

/**
 * Get user settings
 */
export async function getSettings(): Promise<SettingsResponse> {
    const res = await authenticatedFetch("/api/users/me/settings", {
        method: "GET",
    });

    return await handleApiResponse(res);
}

/**
 * Update user settings
 */
export async function updateSettings(data: UpdateSettingsRequest): Promise<SettingsResponse> {
    const res = await authenticatedFetch("/api/users/me/settings", {
        method: "PUT",
        body: JSON.stringify(data),
    });

    return await handleApiResponse(res);
}

// --- Metadata (Public endpoints for dropdowns) ---

/**
 * Get list of provinces/cities
 */
export async function getProvinces(): Promise<string[]> {
    const res = await publicFetch("/api/users/provinces", {
        method: "GET",
    });

    return await handleApiResponse(res);
}

/**
 * Get list of interests
 */
export async function getInterests(): Promise<string[]> {
    const res = await publicFetch("/api/users/interests", {
        method: "GET",
    });

    return await handleApiResponse(res);
}

/**
 * Get list of genders
 */
export async function getGenders(): Promise<string[]> {
    const res = await publicFetch("/api/users/genders", {
        method: "GET",
    });

    return await handleApiResponse(res);
}

// --- Search ---

export interface PaginatedUsersResponse {
    users: UserResponse[];
    pagination: {
        page: number;
        page_size: number;
        total_items: number;
        total_pages: number;
    };
}

/**
 * Search users (public endpoint)
 * @param username - Optional username filter
 * @param page - Page number (default: 1)
 * @param pageSize - Page size (default: 10)
 */
export async function searchUsers(username?: string, page = 1, pageSize = 5): Promise<PaginatedUsersResponse> {
    const params = new URLSearchParams({
        page: page.toString(),
        pageSize: pageSize.toString(),
    });

    // Add username filter if provided (backend now supports it!)
    if (username) {
        params.append('username', username);
    }
    
    const res = await publicFetch(`/api/users?${params.toString()}`, {
        method: "GET",
    });

    return await handleApiResponse(res);
}
