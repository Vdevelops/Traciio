import apiClient from "@/lib/api-client";
import type { ProfileResponse, UserResponse } from "../types";
import type { UpdateProfileFormData, ChangePasswordFormData } from "../schemas/profile.schema";

export const profileService = {
  async getProfile(params?: { start_date?: string; end_date?: string }): Promise<ProfileResponse> {
    // Use new user-scoped endpoint (extracts user ID from JWT)
    // Pass optional date range parameters for filtering
    const response = await apiClient.get<ProfileResponse>(`/users/me/settings-summary`, { params });
    return response.data;
  },

  async updateProfile(userId: string, data: UpdateProfileFormData): Promise<UserResponse> {
    const response = await apiClient.put<UserResponse>(`/users/${userId}/profile`, data);
    return response.data;
  },

  async changePassword(userId: string, data: ChangePasswordFormData): Promise<void> {
    await apiClient.put(`/users/${userId}/password`, {
      current_password: data.current_password,
      password: data.password,
      confirm_password: data.confirm_password,
    });
  },
};

