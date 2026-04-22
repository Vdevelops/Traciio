"use client";

import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { productAnalyticsService } from "../services/productAnalyticsService";

export function useProductPerformance(
  productId: string,
  startDate?: string,
  endDate?: string
) {
  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ["product-analytics", "performance", productId, { start_date: startDate, end_date: endDate }],
    queryFn: () =>
      productAnalyticsService.getProductPerformance(productId, {
        start_date: startDate,
        end_date: endDate,
      }),
    enabled: !!productId,
  });

  return {
    performance: data?.data,
    filters: data?.meta?.filters,
    isLoading,
    error,
    refetch,
  };
}

export function useProductComparison(
  productIds: string[],
  startDate?: string,
  endDate?: string
) {
  const { data, isLoading, error, refetch } = useQuery({
    queryKey: [
      "product-analytics",
      "comparison",
      { product_ids: productIds, start_date: startDate, end_date: endDate },
    ],
    queryFn: () =>
      productAnalyticsService.getProductComparison({
        product_ids: productIds,
        start_date: startDate,
        end_date: endDate,
      }),
    enabled: productIds.length > 0,
  });

  return {
    comparison: data?.data?.products ?? [],
    filters: data?.meta?.filters,
    isLoading,
    error,
    refetch,
  };
}

export function useProductTrends(
  productId: string,
  startDate?: string,
  endDate?: string,
  groupBy?: "day" | "week" | "month" | "year"
) {
  const { data, isLoading, error, refetch } = useQuery({
    queryKey: [
      "product-analytics",
      "trends",
      productId,
      { start_date: startDate, end_date: endDate, group_by: groupBy },
    ],
    queryFn: () =>
      productAnalyticsService.getProductTrends(productId, {
        start_date: startDate,
        end_date: endDate,
        group_by: groupBy,
      }),
    enabled: !!productId,
  });

  return {
    trends: data?.data,
    filters: data?.meta?.filters,
    isLoading,
    error,
    refetch,
  };
}

export function useTopProducts() {
  const { data, isLoading, error } = useQuery({
    queryKey: ["product-analytics", "top-products"],
    queryFn: () =>
      productAnalyticsService.getProductsList({
        period: "month", // Default to month like leaderboard
        sort_by: "total_sold",
        order: "desc",
        limit: 3,
      }),
  });

  // Calculate rank based on index since it's the top list
  const topProducts = (data?.data ?? []).map((product: any, index: number) => ({
    ...product,
    rank: index + 1,
  }));

  return {
    products: topProducts,
    isLoading,
    error,
  };
}

export function useProductStats(
  period: "all" | "day" | "week" | "month" | "year" = "month",
  startDate?: string,
  endDate?: string
) {
  const { data, isLoading, error } = useQuery({
    queryKey: ["product-analytics", "product-stats", { period, startDate, endDate, limit: 100 }],
    queryFn: () =>
      productAnalyticsService.getProductsList({
        period,
        start_date: startDate,
        end_date: endDate,
        sort_by: "total_sold",
        order: "desc",
        limit: 100, // Fetch top 100 to calculate stats approximates
      }),
  });

  return {
    products: data?.data ?? [],
    isLoading,
    error,
  };
}

export function useProductsList(
  period: "all" | "day" | "week" | "month" | "year" = "month",
  startDate?: string,
  endDate?: string,
  search?: string
) {
  const [sortBy, setSortBy] = useState<"total_sold" | "revenue" | "profit" | "name">("total_sold");
  const [orderBy, setOrderBy] = useState<"asc" | "desc">("desc");
  const [page, setPage] = useState(1);
  const [perPage, setPerPage] = useState(10); // Default to 10 for table

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ["product-analytics", "products-list", { period, startDate, endDate, search, sort_by: sortBy, order: orderBy, page, per_page: perPage }],
    queryFn: () =>
      productAnalyticsService.getProductsList({
        period,
        start_date: startDate,
        end_date: endDate,
        search,
        sort_by: sortBy,
        order: orderBy,
        page,
        per_page: perPage,
      }),
  });

  return {
    products: data?.data ?? [],
    // Safely access pagination from meta, handling potential undefined structure
    pagination: data?.meta?.pagination ?? {
      page: 1,
      per_page: perPage,
      total: 0,
      total_pages: 0,
      has_next: false,
      has_prev: false,
    },
    filters: data?.meta?.filters,
    isLoading,
    error,
    refetch,
    sortBy,
    setSortBy,
    orderBy,
    setOrderBy,
    page,
    setPage,
    perPage,
    setPerPage,
  };
}

/**
 * Hook for fetching monthly sales data for the chart
 */
export function useMonthlySales(startDate?: string, endDate?: string) {
  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ["product-analytics", "monthly-sales", { startDate, endDate }],
    queryFn: () => productAnalyticsService.getMonthlySales({ start_date: startDate, end_date: endDate }),
  });

  return {
    monthlySales: data?.data,
    filters: data?.meta?.filters,
    isLoading,
    error,
    refetch,
  };
}

/**
 * Hook for fetching monthly sales data for a specific product
 */
export function useProductMonthlySales(productId: string, startDate?: string, endDate?: string) {
  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ["product-analytics", "product-monthly-sales", productId, { startDate, endDate }],
    queryFn: () => productAnalyticsService.getProductMonthlySales(productId, { start_date: startDate, end_date: endDate }),
    enabled: !!productId,
  });

  return {
    monthlySales: data?.data,
    filters: data?.meta?.filters,
    isLoading,
    error,
    refetch,
  };
}

export function useUserProductSales(
  userId: string,
  options?: {
    startDate?: string;
    endDate?: string;
  }
) {
  const [sortBy, setSortBy] = useState<"total_sold" | "revenue" | "profit" | "name">("total_sold");
  const [orderBy, setOrderBy] = useState<"asc" | "desc">("desc");
  const [page, setPage] = useState(1);
  const [perPage, setPerPage] = useState(20);

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: [
      "product-analytics",
      "user-products",
      userId,
      { start_date: options?.startDate, end_date: options?.endDate, sort_by: sortBy, order: orderBy, page, per_page: perPage },
    ],
    queryFn: () =>
      productAnalyticsService.getUserProductSales(userId, {
        start_date: options?.startDate,
        end_date: options?.endDate,
        sort_by: sortBy,
        order: orderBy,
        page,
        per_page: perPage,
      }),
    enabled: !!userId,
    staleTime: 30000, // 30 seconds
  });

  return {
    products: data?.data ?? [],
    pagination: data?.meta?.pagination,
    filters: data?.meta?.filters,
    isLoading,
    error,
    refetch,
    sortBy,
    setSortBy,
    orderBy,
    setOrderBy,
    page,
    setPage,
    perPage,
    setPerPage,
  };
}
