"use client";

import { useQuery } from "@tanstack/react-query";
import { reportService } from "../services/reportService";
import type { ReportRequestParams } from "../types";

interface ReportQueryOptions {
  enabled?: boolean;
}

export function useVisitReportReport(params?: ReportRequestParams, options?: ReportQueryOptions) {
  return useQuery({
    queryKey: ["reports", "visit-reports", params],
    queryFn: () => reportService.getVisitReportReport(params),
    enabled: options?.enabled,
  });
}

export function usePipelineReport(params?: ReportRequestParams, options?: ReportQueryOptions) {
  return useQuery({
    queryKey: ["reports", "pipeline", params],
    queryFn: () => reportService.getPipelineReport(params),
    enabled: options?.enabled,
  });
}

export function useSalesPerformanceReport(params?: ReportRequestParams, options?: ReportQueryOptions) {
  return useQuery({
    queryKey: ["reports", "sales-performance", params],
    queryFn: () => reportService.getSalesPerformanceReport(params),
    enabled: options?.enabled,
  });
}

export function useAccountActivityReport(params?: ReportRequestParams, options?: ReportQueryOptions) {
  return useQuery({
    queryKey: ["reports", "account-activity", params],
    queryFn: () => reportService.getAccountActivityReport(params),
    enabled: options?.enabled ?? !!params?.account_id, // Only fetch if account_id is provided
  });
}
