import apiClient from "@/lib/api-client";
import type { SalesRepKPIResponse, SalesManagerKPIResponse, SalesRepKPIParams, SalesManagerKPIParams } from "../types";
import type { ApiResponse } from "@/types/api";

export const kpiService = {
  async getSalesRepKPI(params: SalesRepKPIParams) {
    const response = await apiClient.get<ApiResponse<SalesRepKPIResponse>>("/kpi/sales-rep", { params });
    return response.data.data;
  },

  async getSalesManagerKPI(params: SalesManagerKPIParams) {
    const response = await apiClient.get<ApiResponse<SalesManagerKPIResponse>>("/kpi/sales-manager", { params });
    return response.data.data;
  },
};
