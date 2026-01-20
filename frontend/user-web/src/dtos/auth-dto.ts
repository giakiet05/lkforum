import type { UserResponse } from "./user-dto";

export interface sendEmailVerificationRequest {
    email: string;
}

export interface VerifyEmailRequest {
    email: string;
    otp: string;
}

export interface CompleteRegistrationRequest {
  verification_token: string;
  username: string;
  password: string;
}

export interface ResendOTPRequest {
    email: string;
}

export interface LoginRequest {
  identifier: string;
  password: string;
}

export interface RefreshTokenRequest {
  refresh_token: string;
}

export interface LogoutRequest {
    access_token: string;
    refresh_token: string;
}

export interface CompleteGoogleSetupRequest {
    setup_token: string;
    username: string;
}


export interface AuthResponse {
  user: UserResponse;
  access_token: string;
  refresh_token: string;
}

export interface RefreshTokenResponse {
  access_token: string;
  refresh_token: string;
}
