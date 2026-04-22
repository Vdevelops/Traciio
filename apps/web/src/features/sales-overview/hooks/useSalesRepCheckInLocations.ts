"use client";

import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { salesOverviewService } from "../services/salesOverviewService";
import type { GetSalesRepCheckInLocationsRequest } from "../types";

export function useSalesRepCheckInLocations(userId: string, initialParams?: GetSalesRepCheckInLocationsRequest) {
  const [page, setPage] = useState(initialParams?.page ?? 1);
  const [perPage, setPerPage] = useState(initialParams?.per_page ?? 50);

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ["sales-overview", "sales-rep", userId, "check-in-locations", { 
      ...initialParams, 
      page, 
      per_page: perPage 
    }],
    queryFn: async () => {
      return await salesOverviewService.getSalesRepCheckInLocations(userId, {
        ...initialParams,
        page,
        per_page: perPage,
      });
    },
    enabled: !!userId,
  });

  const locations = data?.data?.check_in_locations ?? [];

  return {
    locations,
    totalVisits: data?.data?.total_visits ?? 0,
    period: data?.data?.period,
    isLoading,
    error,
    refetch,
    page,
    setPage,
    perPage,
    setPerPage,
  };
}

