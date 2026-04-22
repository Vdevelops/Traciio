"use client";

import { useQuery } from "@tanstack/react-query";
import { dashboardService } from "../services/dashboardService";
import type { DashboardRequestParams } from "../types";

export function useDashboardOverview(params?: DashboardRequestParams, options?: { enabled?: boolean }) {
  return useQuery({
    queryKey: ["dashboard", "overview", params],
    queryFn: () => dashboardService.getOverview(params),
    enabled: options?.enabled,
  });
}

export function useVisitStatistics(params?: DashboardRequestParams, options?: { enabled?: boolean }) {
  return useQuery({
    queryKey: ["dashboard", "visits", params],
    queryFn: () => dashboardService.getVisitStatistics(params),
    enabled: options?.enabled,
  });
}

export function usePipelineSummary(params?: DashboardRequestParams, options?: { enabled?: boolean }) {
  return useQuery({
    queryKey: ["dashboard", "pipeline", params],
    queryFn: () => dashboardService.getPipelineSummary(params),
    enabled: options?.enabled,
  });
}

export function useTopAccounts(params?: DashboardRequestParams, options?: { enabled?: boolean }) {
  return useQuery({
    queryKey: ["dashboard", "top-accounts", params],
    queryFn: () => dashboardService.getTopAccounts(params),
    enabled: options?.enabled,
  });
}

export function useTopSalesRep(params?: DashboardRequestParams, options?: { enabled?: boolean }) {
  return useQuery({
    queryKey: ["dashboard", "top-sales-rep", params],
    queryFn: () => dashboardService.getTopSalesRep(params),
    enabled: options?.enabled,
  });
}

export function useRecentActivities(params?: DashboardRequestParams, options?: { enabled?: boolean }) {
  return useQuery({
    queryKey: ["dashboard", "recent-activities", params],
    queryFn: () => dashboardService.getRecentActivities(params),
    enabled: options?.enabled,
  });
}

export function useActivityTrends(params?: DashboardRequestParams, options?: { enabled?: boolean }) {
  return useQuery({
    queryKey: ["dashboard", "activity-trends", params],
    queryFn: () => dashboardService.getActivityTrends(params),
    enabled: options?.enabled,
  });
}

// ============================================================================
// Super Admin Dashboard Hooks
// ============================================================================

export function useSuperAdminUsersByRole(params?: DashboardRequestParams, options?: { enabled?: boolean }) {
  return useQuery({
    queryKey: ["dashboard", "super-admin", "users-by-role", params],
    queryFn: () => dashboardService.getSuperAdminUsersByRole(params),
    enabled: options?.enabled,
  });
}

export function useSuperAdminSystemActivity(params?: DashboardRequestParams, options?: { enabled?: boolean }) {
  return useQuery({
    queryKey: ["dashboard", "super-admin", "system-activity", params],
    queryFn: () => dashboardService.getSuperAdminSystemActivity(params),
    enabled: options?.enabled,
  });
}

export function useSuperAdminAIUsage(params?: DashboardRequestParams, options?: { enabled?: boolean }) {
  return useQuery({
    queryKey: ["dashboard", "super-admin", "ai-usage", params],
    queryFn: () => dashboardService.getSuperAdminAIUsage(params),
    enabled: options?.enabled,
  });
}

export function useSuperAdminDataGrowth(params?: DashboardRequestParams, options?: { enabled?: boolean }) {
  return useQuery({
    queryKey: ["dashboard", "super-admin", "data-growth", params],
    queryFn: () => dashboardService.getSuperAdminDataGrowth(params),
    enabled: options?.enabled,
  });
}

export function useSuperAdminErrorSummary(params?: DashboardRequestParams, options?: { enabled?: boolean }) {
  return useQuery({
    queryKey: ["dashboard", "super-admin", "error-summary", params],
    queryFn: () => dashboardService.getSuperAdminErrorSummary(params),
    enabled: options?.enabled,
  });
}

// ============================================================================
// Admin Dashboard Hooks
// ============================================================================

export function useAdminTotalLeads(params?: DashboardRequestParams, options?: { enabled?: boolean }) {
  return useQuery({
    queryKey: ["dashboard", "admin", "total-leads", params],
    queryFn: () => dashboardService.getAdminTotalLeads(params),
    enabled: options?.enabled,
  });
}

export function useAdminPipelineValue(params?: DashboardRequestParams, options?: { enabled?: boolean }) {
  return useQuery({
    queryKey: ["dashboard", "admin", "pipeline-value", params],
    queryFn: () => dashboardService.getAdminPipelineValue(params),
    enabled: options?.enabled,
  });
}

export function useAdminPendingApprovals(params?: DashboardRequestParams, options?: { enabled?: boolean }) {
  return useQuery({
    queryKey: ["dashboard", "admin", "pending-approvals", params],
    queryFn: () => dashboardService.getAdminPendingApprovals(params),
    enabled: options?.enabled,
  });
}

export function useAdminTaskOverdue(params?: DashboardRequestParams, options?: { enabled?: boolean }) {
  return useQuery({
    queryKey: ["dashboard", "admin", "task-overdue", params],
    queryFn: () => dashboardService.getAdminTaskOverdue(params),
    enabled: options?.enabled,
  });
}

// ============================================================================
// Sales Manager Dashboard Hooks
// ============================================================================

export function useSalesManagerPipelineFunnel(params?: DashboardRequestParams) {
  return useQuery({
    queryKey: ["dashboard", "sales-manager", "pipeline-funnel", params],
    queryFn: () => dashboardService.getSalesManagerPipelineFunnel(params),
  });
}

export function useSalesManagerTargetVsActual(params?: DashboardRequestParams) {
  return useQuery({
    queryKey: ["dashboard", "sales-manager", "target-vs-actual", params],
    queryFn: () => dashboardService.getSalesManagerTargetVsActual(params),
  });
}

export function useSalesManagerVisitCompletion(params?: DashboardRequestParams) {
  return useQuery({
    queryKey: ["dashboard", "sales-manager", "visit-completion", params],
    queryFn: () => dashboardService.getSalesManagerVisitCompletion(params),
  });
}

export function useSalesManagerDealsAtRisk(params?: DashboardRequestParams) {
  return useQuery({
    queryKey: ["dashboard", "sales-manager", "deals-at-risk", params],
    queryFn: () => dashboardService.getSalesManagerDealsAtRisk(params),
  });
}

export function useSalesManagerTeamDraftApprovals() {
  return useQuery({
    queryKey: ["dashboard", "sales-manager", "team-draft-approvals"],
    queryFn: () => dashboardService.getSalesManagerTeamDraftApprovals(),
    staleTime: 1000 * 60 * 2, // 2-minute cache – approvals are time-sensitive
  });
}

// ============================================================================
// Sales Dashboard Hooks
// ============================================================================

export function useSalesTodayTasks(params?: DashboardRequestParams) {
  return useQuery({
    queryKey: ["dashboard", "sales", "today-tasks", params],
    queryFn: () => dashboardService.getSalesTodayTasks(params),
  });
}

export function useSalesAssignedLeads(params?: DashboardRequestParams) {
  return useQuery({
    queryKey: ["dashboard", "sales", "assigned-leads", params],
    queryFn: () => dashboardService.getSalesAssignedLeads(params),
  });
}

export function useSalesUpcomingVisits(params?: DashboardRequestParams) {
  return useQuery({
    queryKey: ["dashboard", "sales", "upcoming-visits", params],
    queryFn: () => dashboardService.getSalesUpcomingVisits(params),
  });
}

export function useSalesReminders(params?: DashboardRequestParams) {
  return useQuery({
    queryKey: ["dashboard", "sales", "reminders", params],
    queryFn: () => dashboardService.getSalesReminders(params),
  });
}

// ============================================================================
// Analyst Dashboard Hooks
// ============================================================================

export function useAnalystRevenueTrend(params?: DashboardRequestParams) {
  return useQuery({
    queryKey: ["dashboard", "analyst", "revenue-trend", params],
    queryFn: () => dashboardService.getAnalystRevenueTrend(params),
  });
}

export function useAnalystConversionRate(params?: DashboardRequestParams) {
  return useQuery({
    queryKey: ["dashboard", "analyst", "conversion-rate", params],
    queryFn: () => dashboardService.getAnalystConversionRate(params),
  });
}

export function useAnalystSalesVelocity(params?: DashboardRequestParams) {
  return useQuery({
    queryKey: ["dashboard", "analyst", "sales-velocity", params],
    queryFn: () => dashboardService.getAnalystSalesVelocity(params),
  });
}

export function useAnalystAIInsights(params?: DashboardRequestParams) {
  return useQuery({
    queryKey: ["dashboard", "analyst", "ai-insights", params],
    queryFn: () => dashboardService.getAnalystAIInsights(params),
  });
}


