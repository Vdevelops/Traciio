import { apiClient } from "@/lib/api-client";
import type {
  ListNotificationsParams,
  ListNotificationsResponse,
  NotificationResponse,
  UnreadCountResponse,
} from "../types";

export const notificationService = {
  async list(params?: ListNotificationsParams): Promise<ListNotificationsResponse> {
    const response = await apiClient.get<ListNotificationsResponse>("/notifications", {
      params,
    });
    return response.data;
  },

  async getUnreadCount(): Promise<number> {
    const response = await apiClient.get<UnreadCountResponse>("/notifications/unread-count");
    return response.data.data.unread_count;
  },

  async markAsRead(id: string): Promise<NotificationResponse> {
    const response = await apiClient.put<NotificationResponse>(`/notifications/${id}/read`);
    return response.data;
  },

  async markAllAsRead(): Promise<{ success: boolean; message: string }> {
    const response = await apiClient.put<{ success: boolean; message: string }>(
      "/notifications/read-all"
    );
    return response.data;
  },

  async delete(id: string): Promise<void> {
    await apiClient.delete(`/notifications/${id}`);
  },
};

