import apiClient from "@/lib/api-client";
import type { SalesRepKPIResponse, SalesManagerKPIResponse, SalesRepKPIParams, SalesManagerKPIParams } from "../types";

export const kpiService = {
  async getSalesRepKPI(params: SalesRepKPIParams) {
    const response = await apiClient.get<SalesRepKPIResponse>("/kpi/sales-rep", { params });
    return response.data;
  },

  async getSalesManagerKPI(params: SalesManagerKPIParams) {
    const response = await apiClient.get<SalesManagerKPIResponse>("/kpi/sales-manager", { params });
    return response.data;
  },
};
