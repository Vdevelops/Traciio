"use client";

import { useQuery } from "@tanstack/react-query";
import { kpiService } from "../services/kpiService";
import type { SalesRepKPIParams, SalesManagerKPIParams } from "../types";

export const kpiQueryKeys = {
  salesRep: (params: SalesRepKPIParams) => ["kpi", "sales-rep", params] as const,
  salesManager: (params: SalesManagerKPIParams) => ["kpi", "sales-manager", params] as const,
};

export function useSalesRepKPI(params: SalesRepKPIParams, options?: { enabled?: boolean }) {
  return useQuery({
    queryKey: kpiQueryKeys.salesRep(params),
    queryFn: () => kpiService.getSalesRepKPI(params),
    enabled: Boolean(params.startDate && params.endDate) && options?.enabled !== false,
    staleTime: 5 * 60 * 1000,
  });
}

export function useSalesManagerKPI(params: SalesManagerKPIParams, options?: { enabled?: boolean }) {
  return useQuery({
    queryKey: kpiQueryKeys.salesManager(params),
    queryFn: () => kpiService.getSalesManagerKPI(params),
    enabled: Boolean(params.startDate && params.endDate) && options?.enabled !== false,
    staleTime: 5 * 60 * 1000,
  });
}
