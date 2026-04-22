import apiClient from "@/lib/api-client";
import type { PipelineSummary } from "../types";

interface ApiResponse<T> {
  success: boolean;
  data: T;
  error?: {
    code: string;
    message: string;
  };
}

export async function getPipelineSummary(
  dateFrom?: string,
  dateTo?: string
): Promise<ApiResponse<PipelineSummary>> {
  const params = new URLSearchParams();
  if (dateFrom) params.append("date_from", dateFrom);
  if (dateTo) params.append("date_to", dateTo);
  const query = params.toString() ? `?${params.toString()}` : "";
  const response = await apiClient.get<ApiResponse<PipelineSummary>>(`/pipeline/summary${query}`);
  return response.data;
}
