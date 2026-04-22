import apiClient from "@/lib/api-client";
import type {
  LeadStatus,
  LeadStatusResponse,
  ListLeadStatusesResponse,
  CreateLeadStatusRequest,
  UpdateLeadStatusRequest,
} from "../types/lead-status";

export const leadStatusService = {
  async list(params?: {
    page?: number;
    per_page?: number;
    search?: string;
    is_active?: boolean;
    sort_by?: string;
    sort_order?: "asc" | "desc";
  }): Promise<ListLeadStatusesResponse> {
    const response = await apiClient.get<ListLeadStatusesResponse>("/lead-statuses", { params });
    return response.data;
  },

  async listAll(): Promise<ListLeadStatusesResponse> {
    const response = await apiClient.get<ListLeadStatusesResponse>("/lead-statuses/all");
    return response.data;
  },

  async getById(id: string): Promise<LeadStatusResponse> {
    const response = await apiClient.get<LeadStatusResponse>(`/lead-statuses/${id}`);
    return response.data;
  },

  async create(data: CreateLeadStatusRequest): Promise<LeadStatusResponse> {
    const response = await apiClient.post<LeadStatusResponse>("/lead-statuses", data);
    return response.data;
  },

  async update(id: string, data: UpdateLeadStatusRequest): Promise<LeadStatusResponse> {
    const response = await apiClient.put<LeadStatusResponse>(`/lead-statuses/${id}`, data);
    return response.data;
  },

  async delete(id: string): Promise<{ success: boolean }> {
    const response = await apiClient.delete<{ success: boolean }>(`/lead-statuses/${id}`);
    return response.data;
  },

  async setDefault(id: string): Promise<{ success: boolean }> {
    const response = await apiClient.patch<{ success: boolean }>(`/lead-statuses/${id}/set-default`);
    return response.data;
  },
};
