"use client";

import { useState, useMemo, useEffect } from "react";
import { useQuery } from "@tanstack/react-query";
import { useDebounce } from "@/hooks/use-debounce";
import { salesOverviewService } from "../services/salesOverviewService";
import type { ListSalesPerformanceRequest } from "../types";

// Helper function to get default monthly date range (month to date)
const getDefaultMonthlyRange = () => {
  const today = new Date();
  today.setHours(23, 59, 59, 999);
  
  const startOfMonth = new Date(today);
  startOfMonth.setDate(1);
  startOfMonth.setHours(0, 0, 0, 0);
  
  // Set end date to end of month to include all potential data (especially for demos/seeding)
  const endOfMonth = new Date(today.getFullYear(), today.getMonth() + 1, 0);
  endOfMonth.setHours(23, 59, 59, 999);
  
  const startDateStr = `${startOfMonth.getFullYear()}-${String(startOfMonth.getMonth() + 1).padStart(2, "0")}-${String(startOfMonth.getDate()).padStart(2, "0")}`;
  const endDateStr = `${endOfMonth.getFullYear()}-${String(endOfMonth.getMonth() + 1).padStart(2, "0")}-${String(endOfMonth.getDate()).padStart(2, "0")}`;
  
  return { startDate: startDateStr, endDate: endDateStr };
};

export function useSalesPerformanceList(limit?: number, filters?: { startDate?: string; endDate?: string }) {
  // Set default to monthly (month to date)
  const defaultRange = useMemo(() => getDefaultMonthlyRange(), []);
  
  const [page, setPage] = useState(1);
  const [perPage, setPerPage] = useState(limit || 10);
  const [search, setSearch] = useState("");
  const debouncedSearch = useDebounce(search, 500);

  // Reset page to 1 when a new search is initiated. We update page immediately
  // when the search setter is called (see `handleSetSearch`) so no effect is required.
  const [internalStartDate, setInternalStartDate] = useState<string>(defaultRange.startDate);
  const [internalEndDate, setInternalEndDate] = useState<string>(defaultRange.endDate);

  const startDate = filters?.startDate ?? internalStartDate;
  const endDate = filters?.endDate ?? internalEndDate;
  const [sortBy, setSortBy] = useState<"revenue" | "deals" | "visits" | "tasks" | "name" | "target" | "achievement">("revenue");
  const [order, setOrder] = useState<"asc" | "desc">("desc");

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: [
      "sales-overview",
      "performance",
      "performance",
      { page, per_page: perPage, search: debouncedSearch, start_date: startDate, end_date: endDate, sort_by: sortBy, order },
    ],
    queryFn: async () => {
      const params: ListSalesPerformanceRequest = {
        page,
        per_page: perPage,
        order,
      };

      // Only include API-supported sort fields
      const apiAllowedSorts = ["revenue", "deals", "visits", "tasks", "name"] as const;
      if ((apiAllowedSorts as readonly string[]).includes(sortBy)) {
        const sortValue = sortBy as ListSalesPerformanceRequest["sort_by"];
        params.sort_by = sortValue;
      }

      // Only include search if not empty
      if (debouncedSearch && debouncedSearch.trim() !== "") {
        params.search = debouncedSearch.trim();
      }

      // Include date range if provided
      // If empty (All Time), the backend will handle showing all data
      if (startDate && startDate.trim() !== "") {
        params.start_date = startDate.trim();
      }
      if (endDate && endDate.trim() !== "") {
        params.end_date = endDate.trim();
      }

      return await salesOverviewService.listSalesPerformance(params);
    },
    staleTime: 30000, // 30 seconds
  });

  return {
    performanceList: Array.isArray(data?.data) ? data.data : [],
    pagination: data?.meta?.pagination,
    isLoading,
    error,
    refetch,
    page,
    setPage,
    perPage,
    setPerPage,
    search,
    setSearch: (s: string) => {
      setSearch(s);
      setPage(1);
    },
    startDate,
    setStartDate: setInternalStartDate,
    endDate,
    setEndDate: setInternalEndDate,
    sortBy,
    setSortBy,
    order,
    setOrder,
  };
}

