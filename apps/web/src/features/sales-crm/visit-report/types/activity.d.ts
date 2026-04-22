import type { ActivityType } from "./activity-type";

export interface Activity {
  id: string;
  type: "visit" | "call" | "email" | "task" | "deal"; // Kept for backward compatibility
  activity_type_id?: string;
  account_id?: string;
  contact_id?: string;
  deal_id?: string;
  lead_id?: string;
  user_id: string;
  description: string;
  timestamp: string;
  metadata?: Record<string, unknown>;
  created_at: string;
  updated_at: string;
  account?: {
    id: string;
    name: string;
  };
  contact?: {
    id: string;
    name: string;
  };
  deal?: {
    id: string;
    title: string;
  };
  user?: {
    id: string;
    name: string;
  };
  activity_type?: ActivityType;
}

export interface ListActivitiesResponse {
  success: boolean;
  data: Activity[];
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

export interface ActivityResponse {
  success: boolean;
  data: Activity;
  timestamp: string;
  request_id: string;
}

export interface ActivityTimelineResponse {
  success: boolean;
  data: Activity[];
  meta?: {
    filters?: Record<string, unknown>;
  };
  timestamp: string;
  request_id: string;
}

