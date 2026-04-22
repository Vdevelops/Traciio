import apiClient from "@/lib/api-client";
import type { Account, ListAccountsResponse, AccountResponse } from "../types";
import type { CreateAccountFormData, UpdateAccountFormData } from "../schemas/account.schema";

export const accountService = {
  async list(params?: {
    page?: number;
    per_page?: number;
    search?: string;
    status?: string;
    category_id?: string;
    assigned_to?: string;
  }): Promise<ListAccountsResponse> {
    const response = await apiClient.get<ListAccountsResponse>("/accounts", { params });
    return response.data;
  },

  async getById(id: string): Promise<AccountResponse> {
    const response = await apiClient.get<AccountResponse>(`/accounts/${id}`);
    return response.data;
  },

  async create(data: CreateAccountFormData): Promise<AccountResponse> {
    const response = await apiClient.post<AccountResponse>("/accounts", data);
    return response.data;
  },

  async update(id: string, data: UpdateAccountFormData): Promise<AccountResponse> {
    const response = await apiClient.put<AccountResponse>(`/accounts/${id}`, data);
    return response.data;
  },

  async delete(id: string): Promise<void> {
    await apiClient.delete(`/accounts/${id}`);
  },

  async listForMap(params?: {
    status?: string;
  }): Promise<ListAccountsResponse> {
    const response = await apiClient.get<ListAccountsResponse>("/accounts/map", { params });
    return response.data;
  },

  async listByBBox(params: {
    north: number;
    south: number;
    east: number;
    west: number;
    search?: string;
    status?: string;
    category_id?: string;
    limit?: number;
  }): Promise<ListAccountsResponse> {
    const response = await apiClient.get<ListAccountsResponse>("/accounts/map/bbox", { params });
    return response.data;
  },
};

