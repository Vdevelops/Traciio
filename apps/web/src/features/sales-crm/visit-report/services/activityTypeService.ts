import apiClient from "@/lib/api-client";
import type { ActivityType, ListActivityTypesResponse, ActivityTypeResponse } from "../types/activity-type";

export const activityTypeService = {
  async list(params?: {
    status?: string;
  }): Promise<ListActivityTypesResponse> {
    const response = await apiClient.get<ListActivityTypesResponse>("/visit-reports/activity-types", { params });
    return response.data;
  },

  async getById(id: string): Promise<ActivityTypeResponse> {
    const response = await apiClient.get<ActivityTypeResponse>(`/visit-reports/activity-types/${id}`);
    return response.data;
  },

  async create(data: {
    name: string;
    code: string;
    description?: string;
    icon?: string;
    badge_color?: string;
    status?: string;
    order?: number;
  }): Promise<ActivityTypeResponse> {
    const response = await apiClient.post<ActivityTypeResponse>("/visit-reports/activity-types", data);
    return response.data;
  },

  async update(id: string, data: {
    name?: string;
    code?: string;
    description?: string;
    icon?: string;
    badge_color?: string;
    status?: string;
    order?: number;
  }): Promise<ActivityTypeResponse> {
    const response = await apiClient.put<ActivityTypeResponse>(`/visit-reports/activity-types/${id}`, data);
    return response.data;
  },

  async delete(id: string): Promise<void> {
    await apiClient.delete(`/visit-reports/activity-types/${id}`);
  },
};

