import apiClient from "@/lib/api-client";
import type { RevenueForecast } from "../types";

interface ApiResponse<T> {
  success: boolean;
  data: T;
  error?: {
    code: string;
    message: string;
  };
}

export async function getRevenueForecast(
  period: "month" | "quarter" | "year" = "month"
): Promise<ApiResponse<RevenueForecast>> {
  const response = await apiClient.get<ApiResponse<RevenueForecast>>(`/pipelines/forecast`, {
    params: { period },
  });
  return response.data;
}
