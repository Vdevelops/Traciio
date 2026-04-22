"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useAuthStore } from "@/features/auth/stores/useAuthStore";
import { profileService } from "../services/profileService";
import type { UpdateProfileFormData, ChangePasswordFormData } from "../schemas/profile.schema";
import { toast } from "sonner";

export function useProfile(dateRange?: { start_date?: string; end_date?: string }) {
  return useQuery({
    queryKey: ["profile", "me", dateRange?.start_date, dateRange?.end_date],
    queryFn: () => profileService.getProfile(dateRange),
  });
}

export function useUpdateProfile() {
  const queryClient = useQueryClient();
  const { user, setUser } = useAuthStore();

  return useMutation({
    mutationFn: (data: UpdateProfileFormData) => {
      if (!user?.id) throw new Error("User not authenticated");
      return profileService.updateProfile(user.id, data);
    },
    onSuccess: (response) => {
      if (response.success && response.data && user) {
        // Update auth store with new user data
        setUser({
          ...user,
          name: response.data.name,
          avatar_url: response.data.avatar_url,
        });
        queryClient.invalidateQueries({ queryKey: ["profile", user.id] });
        toast.success("Profile updated successfully");
      }
    },
    onError: (error: unknown) => {
      const err = error as { response?: { data?: { error?: { message?: string } } } };
      const message = err.response?.data?.error?.message || "Failed to update profile";
      toast.error(message);
    },
  });
}

export function useChangePassword() {
  const { user } = useAuthStore();

  return useMutation({
    mutationFn: (data: ChangePasswordFormData) => {
      if (!user?.id) throw new Error("User not authenticated");
      return profileService.changePassword(user.id, data);
    },
    onSuccess: () => {
      toast.success("Password changed successfully");
    },
    onError: (error: unknown) => {
      const err = error as { response?: { data?: { error?: { message?: string } } } };
      const message = err.response?.data?.error?.message || "Failed to change password";
      toast.error(message);
    },
  });
}

