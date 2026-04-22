import apiClient from "@/lib/api-client";
import type {
  Industry,
  IndustryResponse,
  ListIndustriesResponse,
  CreateIndustryRequest,
  UpdateIndustryRequest,
} from "../types/industry";

export const industryService = {
  async list(params?: {
    page?: number;
    per_page?: number;
    search?: string;
    is_active?: boolean;
    sort_by?: string;
    sort_order?: "asc" | "desc";
  }): Promise<ListIndustriesResponse> {
    const response = await apiClient.get<ListIndustriesResponse>("/industries", { params });
    return response.data;
  },

  async listAll(): Promise<ListIndustriesResponse> {
    const response = await apiClient.get<ListIndustriesResponse>("/industries/all");
    return response.data;
  },

  async getById(id: string): Promise<IndustryResponse> {
    const response = await apiClient.get<IndustryResponse>(`/industries/${id}`);
    return response.data;
  },

  async create(data: CreateIndustryRequest): Promise<IndustryResponse> {
    const response = await apiClient.post<IndustryResponse>("/industries", data);
    return response.data;
  },

  async update(id: string, data: UpdateIndustryRequest): Promise<IndustryResponse> {
    const response = await apiClient.put<IndustryResponse>(`/industries/${id}`, data);
    return response.data;
  },

  async delete(id: string): Promise<{ success: boolean }> {
    const response = await apiClient.delete<{ success: boolean }>(`/industries/${id}`);
    return response.data;
  },
};

