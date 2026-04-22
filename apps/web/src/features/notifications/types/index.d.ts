export type NotificationType = "reminder" | "task" | "deal" | "activity";

export interface Notification {
  id: string;
  user_id: string;
  title: string;
  message: string;
  type: NotificationType;
  is_read: boolean;
  read_at: string | null;
  data: string; // JSON string
  created_at: string;
  updated_at: string;
}

export interface NotificationData {
  reminder_id?: string;
  task_id?: string;
  task_title?: string;
  remind_at?: string;
  [key: string]: unknown;
}

export interface ListNotificationsParams {
  page?: number;
  per_page?: number;
  type?: NotificationType;
  is_read?: boolean;
}

export interface ListNotificationsResponse {
  success: boolean;
  data: Notification[];
  meta: {
    pagination: {
      page: number;
      per_page: number;
      total: number;
      total_pages: number;
      has_next: boolean;
      has_prev: boolean;
    };
    filters?: Record<string, unknown>;
  };
  timestamp: string;
  request_id: string;
}

export interface NotificationResponse {
  success: boolean;
  data: Notification;
  timestamp: string;
  request_id: string;
}

export interface UnreadCountResponse {
  success: boolean;
  data: {
    unread_count: number;
  };
  timestamp: string;
  request_id: string;
}

// WebSocket message types
export interface WebSocketMessage {
  type: "notification.created" | "notification.updated" | "notification.deleted" | "permissions.updated";
  data: Notification | { user_id: string; notification_id: string } | { user_id: string; role_id: string };
}

