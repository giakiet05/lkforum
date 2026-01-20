import type {
    LoginRequest,
    AuthResponse,
    CompleteRegistrationRequest,
    VerifyEmailRequest,
    sendEmailVerificationRequest,
    LogoutRequest,
    CompleteGoogleSetupRequest
} from "../dtos/auth-dto";
import {setAuth, clearAuth} from "../stores/auth-store";
import {publicFetch, handleApiResponse, API_BASE_URL} from "./api";
import {ACCESS_TOKEN_KEY, REFRESH_TOKEN_KEY, USER_KEY} from "../constants/auth-constants";
import {getValidAccessToken, isTokenExpired} from "../auth/token";

// --- Authentication Flows ---

/**
 * Validate auth state khi app khởi động
 * Tự động refresh token nếu cần, hoặc clear auth nếu hết hạn
 */
export async function validateAuth(): Promise<void> {
    const accessToken = localStorage.getItem(ACCESS_TOKEN_KEY);
    const refreshToken = localStorage.getItem(REFRESH_TOKEN_KEY);

    // Không có token nào → chắc chắn chưa login
    if (!accessToken && !refreshToken) {
        clearAuth();
        return;
    }

    // Có access token và còn hạn → OK
    if (accessToken && !isTokenExpired(accessToken)) {
        // User vẫn đang login, authStore đã đúng
        return;
    }

    // Access token hết hạn, thử refresh
    if (refreshToken) {
        const newToken = await getValidAccessToken();
        if (newToken) {
            // Refresh thành công → update authStore
            const user = localStorage.getItem(USER_KEY) ? JSON.parse(localStorage.getItem(USER_KEY)!) : null;
            if (user) {
                setAuth(user);
            }
            return;
        }
    }

    // Refresh token cũng hết hạn hoặc không có → logout
    clearAuth();
    localStorage.removeItem(ACCESS_TOKEN_KEY);
    localStorage.removeItem(REFRESH_TOKEN_KEY);
    localStorage.removeItem(USER_KEY);
}

export async function sendVerificationEmail(email: string): Promise<void> {
    const reqBody: sendEmailVerificationRequest = {email};
    const res = await publicFetch(`/api/auth/local/send-verification`, {
        method: "POST",
        body: JSON.stringify(reqBody),
    });
    // We only care about success/failure, not the data
    await handleApiResponse(res);
}

export async function verifyEmail(email: string, otp: string): Promise<string> {
    const reqBody: VerifyEmailRequest = {email, otp};
    const res = await publicFetch(`/api/auth/local/verify-email`, {
        method: "POST",
        body: JSON.stringify(reqBody),
    });
    const data = await handleApiResponse(res);
    // Assuming the response data is { verification_token: "..." }
    return data.verification_token;
}

export async function completeRegistration(req: CompleteRegistrationRequest): Promise<AuthResponse> {
    const res = await publicFetch(`/api/auth/local/complete-registration`, {
        method: "POST",
        body: JSON.stringify(req),
    });

    const data: AuthResponse = await handleApiResponse(res);
    localStorage.setItem(ACCESS_TOKEN_KEY, data.access_token);
    localStorage.setItem(REFRESH_TOKEN_KEY, data.refresh_token);
    localStorage.setItem(USER_KEY, JSON.stringify(data.user));
    setAuth(data.user);
    return data;
}

export async function login(credentials: LoginRequest): Promise<AuthResponse> {
    const res = await publicFetch(`/api/auth/local/login`, {
        method: "POST",
        body: JSON.stringify(credentials),
    });

    const data: AuthResponse = await handleApiResponse(res);
    localStorage.setItem(ACCESS_TOKEN_KEY, data.access_token);
    localStorage.setItem(REFRESH_TOKEN_KEY, data.refresh_token);
    localStorage.setItem(USER_KEY, JSON.stringify(data.user));
    setAuth(data.user);
    return data;
}

export async function logout() {
    const accessToken = localStorage.getItem(ACCESS_TOKEN_KEY);
    const refreshToken = localStorage.getItem(REFRESH_TOKEN_KEY);

    // Clear local state immediately for instant UI update
    localStorage.removeItem(ACCESS_TOKEN_KEY);
    localStorage.removeItem(REFRESH_TOKEN_KEY);
    localStorage.removeItem(USER_KEY);
    clearAuth();

    // Only call logout API if we have valid tokens (in background)
    if (accessToken && refreshToken) {
        const req: LogoutRequest = {
            access_token: accessToken,
            refresh_token: refreshToken
        }

        try {
            const res = await publicFetch("/api/auth/logout", {
                method: "POST",
                body: JSON.stringify(req)
            })
            await handleApiResponse(res);
        } catch (error) {
            // Already cleared local auth above
            console.error("Logout API failed:", error);
        }
    }
}

export async function completeGoogleSetup(setupToken: string, username: string): Promise<AuthResponse> {
    const req: CompleteGoogleSetupRequest = {
        setup_token: setupToken,
        username
    };

    const res = await publicFetch(`/api/auth/google/complete-setup`, {
        method: "POST",
        body: JSON.stringify(req),
    });

    const data: AuthResponse = await handleApiResponse(res);
    localStorage.setItem(ACCESS_TOKEN_KEY, data.access_token);
    localStorage.setItem(REFRESH_TOKEN_KEY, data.refresh_token);
    localStorage.setItem(USER_KEY, JSON.stringify(data.user));
    setAuth(data.user);
    return data;
}

// --- Forgot Password Flow ---

/**
 * Gửi OTP để reset password
 */
export async function forgotPassword(email: string): Promise<void> {
    const res = await publicFetch(`${API_BASE_URL}/api/auth/local/forgot-password`, {
        method: "POST",
        body: JSON.stringify({ email }),
    });

    await handleApiResponse(res);
}

/**
 * Xác thực OTP và nhận reset token
 */
export async function verifyResetOTP(email: string, otp: string): Promise<{ reset_token: string }> {
    const res = await publicFetch(`${API_BASE_URL}/api/auth/local/verify-reset-otp`, {
        method: "POST",
        body: JSON.stringify({ email, otp }),
    });

    const data = await handleApiResponse(res);
    return data;
}

/**
 * Reset password với reset token
 */
export async function resetPassword(resetToken: string, newPassword: string): Promise<void> {
    const res = await publicFetch(`${API_BASE_URL}/api/auth/local/reset-password`, {
        method: "POST",
        body: JSON.stringify({ reset_token: resetToken, new_password: newPassword }),
    });

    await handleApiResponse(res);
}

