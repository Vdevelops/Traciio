import apiClient from "@/lib/api-client";
import type { Deal, DealFilters, DealHistoryResponse } from "../types";
import type { DealFormData, DealUpdateData, DealMoveData } from "../schemas/deal.schema";
import type { ListVisitReportsResponse } from "@/features/sales-crm/visit-report/types";
import type { ListActivitiesResponse } from "@/features/sales-crm/visit-report/types/activity";

interface ApiResponse<T> {
  success: boolean;
  data: T;
  meta?: {
    pagination?: {
      page: number;
      per_page: number;
      total: number;
      total_pages: number;
      has_next: boolean;
      has_prev: boolean;
    };
    filters?: Record<string, unknown>;
  };
  error?: {
    code: string;
    message: string;
    details?: unknown;
    field_errors?: Array<{
      field: string;
      message: string;
    }>;
  };
}

// API returns data as array directly, not wrapped in PaginatedResponse
// The pagination info is in meta.pagination

function buildQueryString(filters?: DealFilters, page?: number, limit?: number): string {
  const params = new URLSearchParams();
  if (page) params.append("page", page.toString());
  if (limit) params.append("per_page", limit.toString());
  if (filters?.stage_id) params.append("stage_id", filters.stage_id);
  if (filters?.account_id) params.append("account_id", filters.account_id);
  if (filters?.assigned_to) params.append("assigned_to", filters.assigned_to);
  if (filters?.search) params.append("search", filters.search);
  if (filters?.min_value !== undefined) params.append("min_value", filters.min_value.toString());
  if (filters?.max_value !== undefined) params.append("max_value", filters.max_value.toString());
  if (filters?.date_from) params.append("date_from", filters.date_from);
  if (filters?.date_to) params.append("date_to", filters.date_to);
  
  // Add sorting parameters
  if (page) params.append("sort", "created_at");
  if (limit) params.append("order", "desc");
  
  const query = params.toString();
  return query ? `?${query}` : "";
}

export async function getDeals(
  filters?: DealFilters,
  page: number = 1,
  limit: number = 20,
  sort: string = "created_at",
  order: "asc" | "desc" = "desc"
): Promise<ApiResponse<Deal[]>> {
  const params = new URLSearchParams();
  if (page) params.append("page", page.toString());
  if (limit) params.append("per_page", limit.toString());
  if (sort) params.append("sort", sort);
  if (order) params.append("order", order);
  
  if (filters?.stage_id) params.append("stage_id", filters.stage_id);
  if (filters?.account_id) params.append("account_id", filters.account_id);
  if (filters?.assigned_to) params.append("assigned_to", filters.assigned_to);
  if (filters?.search) params.append("search", filters.search);
  if (filters?.min_value !== undefined) params.append("min_value", filters.min_value.toString());
  if (filters?.max_value !== undefined) params.append("max_value", filters.max_value.toString());
  if (filters?.date_from) params.append("date_from", filters.date_from);
  if (filters?.date_to) params.append("date_to", filters.date_to);

  const queryString = params.toString() ? `?${params.toString()}` : "";
  const response = await apiClient.get<ApiResponse<Deal[]>>(`/deals${queryString}`);
  return response.data;
}

export async function getDealsByStage(
  filters?: DealFilters
): Promise<ApiResponse<Record<string, Deal[]>>> {
  const queryString = buildQueryString(filters);
  const response = await apiClient.get<ApiResponse<Record<string, Deal[]>>>(`/deals/by-stage${queryString}`);
  return response.data;
}

export async function getDeal(id: string): Promise<ApiResponse<Deal>> {
  const response = await apiClient.get<ApiResponse<Deal>>(`/deals/${id}`);
  return response.data;
}

export async function createDeal(data: DealFormData): Promise<ApiResponse<Deal>> {
  const response = await apiClient.post<ApiResponse<Deal>>(`/deals`, data);
  return response.data;
}

export async function updateDeal(id: string, data: DealUpdateData): Promise<ApiResponse<Deal>> {
  const response = await apiClient.put<ApiResponse<Deal>>(`/deals/${id}`, data);
  return response.data;
}

export async function moveDeal(data: DealMoveData): Promise<ApiResponse<Deal>> {
  const response = await apiClient.patch<ApiResponse<Deal>>(`/deals/${data.deal_id}/move`, {
    stage_id: data.stage_id,
    order: data.order,
    reason: data.reason,
    product_items: data.product_items,
  });
  return response.data;
}

export async function deleteDeal(id: string): Promise<ApiResponse<{ message: string }>> {
  const response = await fetch(`/api/v1/deals/${id}`, {
    method: "DELETE",
    headers: { "Content-Type": "application/json" },
    credentials: "include",
  });
  if (!response.ok) {
    const error = await response.json();
    throw error;
  }
  return response.json();
}

export async function getDealVisitReports(
  dealId: string,
  page: number = 1,
  limit: number = 20
): Promise<ListVisitReportsResponse> {
  const params = new URLSearchParams();
  params.append("page", page.toString());
  params.append("per_page", limit.toString());
  const queryString = params.toString();
  const response = await apiClient.get<ListVisitReportsResponse>(
    `/deals/${dealId}/visit-reports?${queryString}`
  );
  return response.data;
}

export async function getDealActivities(
  dealId: string,
  page: number = 1,
  limit: number = 20
): Promise<ListActivitiesResponse> {
  const params = new URLSearchParams();
  params.append("page", page.toString());
  params.append("per_page", limit.toString());
  const queryString = params.toString();
  const response = await apiClient.get<ListActivitiesResponse>(
    `/deals/${dealId}/activities?${queryString}`
  );
  return response.data;
}

export async function getDealHistory(dealId: string): Promise<DealHistoryResponse> {
  const response = await apiClient.get<DealHistoryResponse>(`/deals/${dealId}/history`);
  return response.data;
}
