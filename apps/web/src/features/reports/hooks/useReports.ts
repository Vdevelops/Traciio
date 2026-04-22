"use client";

import { useQuery } from "@tanstack/react-query";
import { reportService } from "../services/reportService";
import type { ReportRequestParams } from "../types";

export function useVisitReportReport(params?: ReportRequestParams) {
  return useQuery({
    queryKey: ["reports", "visit-reports", params],
    queryFn: () => reportService.getVisitReportReport(params),
  });
}

export function usePipelineReport(params?: ReportRequestParams) {
  return useQuery({
    queryKey: ["reports", "pipeline", params],
    queryFn: () => reportService.getPipelineReport(params),
  });
}

export function useSalesPerformanceReport(params?: ReportRequestParams) {
  return useQuery({
    queryKey: ["reports", "sales-performance", params],
    queryFn: () => reportService.getSalesPerformanceReport(params),
  });
}

export function useAccountActivityReport(params?: ReportRequestParams) {
  return useQuery({
    queryKey: ["reports", "account-activity", params],
    queryFn: () => reportService.getAccountActivityReport(params),
    enabled: !!params?.account_id, // Only fetch if account_id is provided
  });
}

