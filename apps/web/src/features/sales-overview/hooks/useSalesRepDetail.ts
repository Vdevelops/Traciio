"use client";

import { useQuery } from "@tanstack/react-query";
import { salesOverviewService } from "../services/salesOverviewService";
import type { GetSalesRepDetailRequest } from "../types";

export function useSalesRepDetail(userId: string, params?: GetSalesRepDetailRequest) {
  // Normalize params for query key (only include defined values)
  const normalizedParams = params
    ? {
        start_date: params.start_date,
        end_date: params.end_date,
        // Don't include period if we're using date range
        ...(params.period && !params.start_date && !params.end_date ? { period: params.period } : {}),
      }
    : undefined;

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ["sales-overview", "sales-rep", userId, normalizedParams],
    queryFn: async () => {
      return await salesOverviewService.getSalesRepDetail(userId, params);
    },
    enabled: !!userId,
    staleTime: 30000, // 30 seconds
  });

  return {
    detail: data?.data, // Direct access to data object
    isLoading,
    error,
    refetch,
  };
}

