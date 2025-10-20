import type { User } from "../models/user";

export interface RegisterDto {
  username: string;
  password: string;
}

export interface LoginRequest {
  identifier: string;
  password: string;
}

export interface RefreshTokenRequest {
  refresh_token: string;
}

export interface RegisterResponse extends LoginResponse {
  user: User;
  access_token: string; 
  refresh_token: string;
}

export interface LoginResponse {
  user: User;
  access_token: string;
  refresh_token: string;
}

export interface RefreshTokenResponse {
  access_token: string;
  refresh_token: string;
}
