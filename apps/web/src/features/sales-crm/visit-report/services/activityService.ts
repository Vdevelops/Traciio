import apiClient from "@/lib/api-client";
import type {
  Activity,
  ListActivitiesResponse,
  ActivityResponse,
  ActivityTimelineResponse,
} from "../types/activity";

export const activityService = {
  async list(params?: {
    page?: number;
    per_page?: number;
    type?: string;
    account_id?: string;
    contact_id?: string;
    deal_id?: string;
    lead_id?: string;
    user_id?: string;
    start_date?: string;
    end_date?: string;
  }): Promise<ListActivitiesResponse> {
    const response = await apiClient.get<ListActivitiesResponse>("/activities", { params });
    return response.data;
  },

  async getById(id: string): Promise<ActivityResponse> {
    const response = await apiClient.get<ActivityResponse>(`/activities/${id}`);
    return response.data;
  },

  async getTimeline(params?: {
    account_id?: string;
    contact_id?: string;
    deal_id?: string;
    lead_id?: string;
    user_id?: string;
    start_date?: string;
    end_date?: string;
    limit?: number;
  }): Promise<ActivityTimelineResponse> {
    const response = await apiClient.get<ActivityTimelineResponse>("/activities/timeline", { params });
    return response.data;
  },

  async create(data: {
    activity_type_id: string;
    account_id?: string;
    contact_id?: string;
    deal_id?: string;
    lead_id?: string;
    description: string;
    timestamp: string;
    metadata?: Record<string, unknown>;
  }): Promise<ActivityResponse> {
    const response = await apiClient.post<ActivityResponse>("/activities", data);
    return response.data;
  },
};

