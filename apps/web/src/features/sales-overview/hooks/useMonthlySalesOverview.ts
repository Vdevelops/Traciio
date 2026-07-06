"use client";

import { useQuery } from "@tanstack/react-query";
import { salesOverviewService } from "../services/salesOverviewService";

export function useMonthlySalesOverview(
  startDate?: string,
  endDate?: string,
  trendMode:
    | "monthly"
    | "mom"
    | "rolling_30d"
    | "rolling_90d"
    | "qoq" = "monthly",
) {
  const { data, isLoading, error, refetch } = useQuery({
    queryKey: [
      "sales-overview",
      "monthly",
      { start_date: startDate, end_date: endDate, trend_mode: trendMode },
    ],
    queryFn: async () => {
      return await salesOverviewService.getMonthlySalesOverview(
        startDate,
        endDate,
        trendMode,
      );
    },
    staleTime: 30000,
  });

  return {
    monthlyData: data?.data,
    isLoading,
    error,
    refetch,
  };
}
