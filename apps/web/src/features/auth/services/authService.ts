import apiClient from "@/lib/api-client";
import type { LoginRequest, LoginResponse } from "../types";

export const authService = {
  async login(credentials: LoginRequest): Promise<LoginResponse> {
    const response = await apiClient.post<LoginResponse>("/auth/login", credentials);
    return response.data;
  },

  async refreshToken(refreshToken: string): Promise<LoginResponse> {
    const response = await apiClient.post<LoginResponse>("/auth/refresh", {
      refresh_token: refreshToken,
    });
    return response.data;
  },

  async logout(): Promise<void> {
    await apiClient.post("/auth/logout");
  },
};

