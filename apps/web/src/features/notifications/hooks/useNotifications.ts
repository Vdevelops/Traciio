"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { notificationService } from "../services/notificationService";
import type { ListNotificationsParams } from "../types";

// Helper to check if user has token (client-side only)
function hasToken(): boolean {
  if (typeof window === "undefined") return false;
  return !!localStorage.getItem("token");
}

export function useNotifications(params?: ListNotificationsParams) {
  return useQuery({
    queryKey: ["notifications", params],
    queryFn: () => notificationService.list(params),
    enabled: hasToken(), // Only fetch when authenticated
    retry: (failureCount, error) => {
      if (error && typeof error === "object" && "response" in error) {
        const axiosError = error as { response?: { status?: number } };
        // Don't retry on 401 (unauthorized) or 404
        if (axiosError.response?.status === 401 || axiosError.response?.status === 404) {
          return false;
        }
      }
      return failureCount < 1;
    },
  });
}

export function useNotificationCount() {
  return useQuery({
    queryKey: ["notifications", "unread-count"],
    queryFn: () => notificationService.getUnreadCount(),
    enabled: hasToken(), // Only fetch when authenticated
    refetchInterval: false, // Disable polling, rely on WS and window focus for best scalability
    retry: (failureCount, error) => {
      if (error && typeof error === "object" && "response" in error) {
        const axiosError = error as { response?: { status?: number } };
        // Don't retry on 401 (unauthorized)
        if (axiosError.response?.status === 401) {
          return false;
        }
      }
      return failureCount < 1;
    },
  });
}

export function useMarkAsRead() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => notificationService.markAsRead(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["notifications"] });
      queryClient.invalidateQueries({ queryKey: ["notifications", "unread-count"] });
    },
    onError: (error: unknown) => {
      const message =
        error instanceof Error ? error.message : "Failed to mark notification as read";
      toast.error(message);
    },
  });
}

export function useMarkAllAsRead() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => notificationService.markAllAsRead(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["notifications"] });
      queryClient.invalidateQueries({ queryKey: ["notifications", "unread-count"] });
      toast.success("All notifications marked as read");
    },
    onError: (error: unknown) => {
      const message =
        error instanceof Error ? error.message : "Failed to mark all notifications as read";
      toast.error(message);
    },
  });
}

export function useDeleteNotification() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => notificationService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["notifications"] });
      queryClient.invalidateQueries({ queryKey: ["notifications", "unread-count"] });
      toast.success("Notification deleted");
    },
    onError: (error: unknown) => {
      const message = error instanceof Error ? error.message : "Failed to delete notification";
      toast.error(message);
    },
  });
}

