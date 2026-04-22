import apiClient from "@/lib/api-client";
import type {
  Lead,
  ListLeadsResponse,
  LeadResponse,
  LeadFormDataResponse,
  LeadAnalyticsResponse,
  ConvertLeadResponse,
  CreateAccountFromLeadResponse,
} from "../types";
import type { CreateLeadFormData, UpdateLeadFormData, ConvertLeadFormData } from "../schemas/lead.schema";
import type { ListVisitReportsResponse } from "@/features/sales-crm/visit-report/types";
import type { ListActivitiesResponse } from "@/features/sales-crm/visit-report/types/activity";

export const leadService = {
  async list(params?: {
    page?: number;
    per_page?: number;
    status?: string;
    source?: string;
    assigned_to?: string;
    search?: string;
    sort?: string;
    order?: "asc" | "desc";
  }): Promise<ListLeadsResponse> {
    const response = await apiClient.get<ListLeadsResponse>("/leads", { params });
    return response.data;
  },

  async getById(id: string): Promise<LeadResponse> {
    const response = await apiClient.get<LeadResponse>(`/leads/${id}`);
    return response.data;
  },

  async getFormData(): Promise<LeadFormDataResponse> {
    const response = await apiClient.get<LeadFormDataResponse>("/leads/form-data");
    return response.data;
  },

  async getAnalytics(params?: {
    start_date?: string;
    end_date?: string;
  }): Promise<LeadAnalyticsResponse> {
    const response = await apiClient.get<LeadAnalyticsResponse>("/leads/analytics", { params });
    return response.data;
  },

  async create(data: CreateLeadFormData): Promise<LeadResponse> {
    const response = await apiClient.post<LeadResponse>("/leads", data);
    return response.data;
  },

  async update(id: string, data: UpdateLeadFormData): Promise<LeadResponse> {
    const response = await apiClient.put<LeadResponse>(`/leads/${id}`, data);
    return response.data;
  },

  async delete(id: string): Promise<void> {
    await apiClient.delete(`/leads/${id}`);
  },

  async convert(id: string, data: ConvertLeadFormData): Promise<ConvertLeadResponse> {
    const response = await apiClient.post<ConvertLeadResponse>(`/leads/${id}/convert`, data);
    return response.data;
  },

  async getVisitReportsByLead(leadId: string, params?: {
    page?: number;
    per_page?: number;
    search?: string;
    status?: string;
    account_id?: string;
    sales_rep_id?: string;
    start_date?: string;
    end_date?: string;
  }): Promise<ListVisitReportsResponse> {
    const response = await apiClient.get<ListVisitReportsResponse>(`/leads/${leadId}/visit-reports`, { params });
    return response.data;
  },

  async getActivitiesByLead(leadId: string, params?: {
    page?: number;
    per_page?: number;
    type?: string;
    account_id?: string;
    contact_id?: string;
    user_id?: string;
    start_date?: string;
    end_date?: string;
  }): Promise<ListActivitiesResponse> {
    const response = await apiClient.get<ListActivitiesResponse>(`/leads/${leadId}/activities`, { params });
    return response.data;
  },

  async createAccountFromLead(leadId: string, data?: {
    category_id?: string;
    create_contact?: boolean;
  }): Promise<CreateAccountFromLeadResponse> {
    const response = await apiClient.post<CreateAccountFromLeadResponse>(`/leads/${leadId}/create-account`, data || {});
    return response.data;
  },
};

