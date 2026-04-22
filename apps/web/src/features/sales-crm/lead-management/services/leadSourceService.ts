import apiClient from "@/lib/api-client";
import type {
  LeadSource,
  LeadSourceResponse,
  ListLeadSourcesResponse,
  CreateLeadSourceRequest,
  UpdateLeadSourceRequest,
} from "../types/lead-source";

export const leadSourceService = {
  async list(params?: {
    page?: number;
    per_page?: number;
    search?: string;
    is_active?: boolean;
    sort_by?: string;
    sort_order?: "asc" | "desc";
  }): Promise<ListLeadSourcesResponse> {
    const response = await apiClient.get<ListLeadSourcesResponse>("/lead-sources", { params });
    return response.data;
  },

  async listAll(): Promise<ListLeadSourcesResponse> {
    const response = await apiClient.get<ListLeadSourcesResponse>("/lead-sources/all");
    return response.data;
  },

  async getById(id: string): Promise<LeadSourceResponse> {
    const response = await apiClient.get<LeadSourceResponse>(`/lead-sources/${id}`);
    return response.data;
  },

  async create(data: CreateLeadSourceRequest): Promise<LeadSourceResponse> {
    const response = await apiClient.post<LeadSourceResponse>("/lead-sources", data);
    return response.data;
  },

  async update(id: string, data: UpdateLeadSourceRequest): Promise<LeadSourceResponse> {
    const response = await apiClient.put<LeadSourceResponse>(`/lead-sources/${id}`, data);
    return response.data;
  },

  async delete(id: string): Promise<{ success: boolean }> {
    const response = await apiClient.delete<{ success: boolean }>(`/lead-sources/${id}`);
    return response.data;
  },
};

