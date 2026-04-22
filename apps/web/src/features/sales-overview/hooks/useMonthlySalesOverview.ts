"use client";

import { useQuery } from "@tanstack/react-query";
import { salesOverviewService } from "../services/salesOverviewService";

export function useMonthlySalesOverview(startDate?: string, endDate?: string) {
  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ["sales-overview", "monthly", { start_date: startDate, end_date: endDate }],
    queryFn: async () => {
      // Pass dates only if valid
      return await salesOverviewService.getMonthlySalesOverview(
        startDate, 
        endDate
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
