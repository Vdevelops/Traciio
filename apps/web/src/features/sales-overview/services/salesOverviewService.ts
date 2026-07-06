import apiClient from "@/lib/api-client";
import type {
  GetSalesPerformanceDetailRequest,
  GetSalesRepDetailRequest,
  ListSalesPerformanceRequest,
  GetSalesRepCheckInLocationsRequest,
  ListProspectOutcomesRequest,
  ListProspectOutcomesResponse,
  MonthlySalesOverviewResponse,
  FunnelDiagnosticsResponse,
  GetFunnelDiagnosticsRequest,
  ListSalesPerformanceResponse,
} from "../types";

export const salesOverviewService = {
  /**
   * Get sales performance detail for a user
   */
  async getSalesPerformanceDetail(
    userId: string,
    params?: GetSalesPerformanceDetailRequest,
  ) {
    const response = await apiClient.get(
      `/sales-overview/performance/${userId}`,
      {
        params,
      },
    );
    return response.data;
  },

  /**
   * Get comprehensive sales rep detail
   */
  async getSalesRepDetail(userId: string, params?: GetSalesRepDetailRequest) {
    const response = await apiClient.get(
      `/sales-overview/sales-rep/${userId}`,
      {
        params,
      },
    );
    return response.data;
  },

  /**
   * List all sales performance
   */
  async listSalesPerformance(
    params?: ListSalesPerformanceRequest,
  ): Promise<ListSalesPerformanceResponse> {
    const response = await apiClient.get<ListSalesPerformanceResponse>(
      "/sales-overview/performance",
      {
        params,
      },
    );
    return response.data;
  },

  /**
   * List prospect outcomes across sales reps
   */
  async listProspectOutcomes(
    params?: ListProspectOutcomesRequest,
  ): Promise<ListProspectOutcomesResponse> {
    const response = await apiClient.get<ListProspectOutcomesResponse>(
      "/sales-overview/prospect-outcomes",
      {
        params,
      },
    );
    return response.data;
  },

  /**
   * Get monthly sales overview
   */
  async getMonthlySalesOverview(
    startDate?: string,
    endDate?: string,
    trendMode:
      | "monthly"
      | "mom"
      | "rolling_30d"
      | "rolling_90d"
      | "qoq" = "monthly",
  ) {
    const params: Record<string, string> = {};
    if (startDate) params.start_date = startDate;
    if (endDate) params.end_date = endDate;
    params.trend_mode = trendMode;

    const response = await apiClient.get<MonthlySalesOverviewResponse>(
      "/sales-overview/monthly-overview",
      { params },
    );
    return response.data;
  },

  async getFunnelDiagnostics(
    params?: GetFunnelDiagnosticsRequest,
  ): Promise<FunnelDiagnosticsResponse> {
    const response = await apiClient.get<FunnelDiagnosticsResponse>(
      "/sales-overview/funnel-diagnostics",
      { params },
    );
    return response.data;
  },

  /**
   * Get sales rep check-in locations
   */
  async getSalesRepCheckInLocations(
    userId: string,
    params?: GetSalesRepCheckInLocationsRequest,
  ) {
    const response = await apiClient.get(
      `/sales-overview/sales-rep/${userId}/check-in-locations`,
      {
        params,
      },
    );
    return response.data;
  },
};
