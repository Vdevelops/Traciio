/**
 * Brick Analytics Hooks
 * React Query hooks for brick analytics data
 */

import { useQuery, UseQueryOptions } from "@tanstack/react-query";
import { brickAnalyticsService } from "../services/brickAnalyticsService";
import type {
  BrickPerformanceMetrics,
  BrickPerformanceResponse,
  BrickPerformanceListResponse,
} from "../types/analytics";

/**
 * Hook to get performance metrics for a single brick
 */
export function useBrickPerformance(
  brickId: string | null | undefined,
  periodStart?: string,
  periodEnd?: string,
  options?: Omit<
    UseQueryOptions<BrickPerformanceResponse, Error>,
    "queryKey" | "queryFn"
  >
) {
  return useQuery({
    queryKey: ["brick-performance", brickId, periodStart, periodEnd],
    queryFn: () => {
      if (!brickId) {
        throw new Error("Brick ID is required");
      }
      return brickAnalyticsService.getBrickPerformance(
        brickId,
        periodStart,
        periodEnd
      );
    },
    enabled: !!brickId,
    ...options,
  });
}

/**
 * Hook to get performance metrics for multiple bricks
 */
export function useBrickPerformanceList(
  brickIds: string[],
  periodStart?: string,
  periodEnd?: string,
  options?: Omit<
    UseQueryOptions<BrickPerformanceListResponse, Error>,
    "queryKey" | "queryFn"
  >
) {
  return useQuery({
    queryKey: ["brick-performance-list", brickIds, periodStart, periodEnd],
    queryFn: () =>
      brickAnalyticsService.listBrickPerformance(
        brickIds,
        periodStart,
        periodEnd
      ),
    enabled: brickIds.length > 0,
    ...options,
  });
}

/**
 * Helper hook to get current month period
 */
export function useCurrentMonthPeriod() {
  const now = new Date();
  const year = now.getFullYear();
  const month = now.getMonth();

  const periodStart = `${year}-${String(month + 1).padStart(2, "0")}-01`;
  const lastDay = new Date(year, month + 1, 0).getDate();
  const periodEnd = `${year}-${String(month + 1).padStart(2, "0")}-${String(
    lastDay
  ).padStart(2, "0")}`;

  return { periodStart, periodEnd };
}

