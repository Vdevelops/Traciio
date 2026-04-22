import apiClient from "@/lib/api-client";
import type {
  DashboardOverviewResponse,
  VisitStatisticsResponse,
  ActivityTrendsResponse,
  PipelineSummaryResponse,
  TopAccountsResponse,
  TopSalesRepResponse,
  RecentActivitiesResponse,
  DashboardRequestParams,
  SuperAdminUsersByRoleResponse,
  SuperAdminSystemActivityResponse,
  SuperAdminAIUsageResponse,
  SuperAdminDataGrowthResponse,
  SuperAdminErrorSummaryResponse,
  AdminTotalLeadsResponse,
  AdminPipelineValueResponse,
  AdminPendingApprovalsResponse,
  AdminTaskOverdueResponse,
  SalesManagerPipelineFunnelResponse,
  SalesManagerTargetVsActualResponse,
  SalesManagerVisitCompletionResponse,
  SalesManagerDealsAtRiskResponse,
  SalesManagerTeamDraftApprovalsResponse,
  SalesTodayTasksResponse,
  SalesAssignedLeadsResponse,
  SalesUpcomingVisitsResponse,
  SalesRemindersResponse,
  AnalystRevenueTrendResponse,
  AnalystConversionRateResponse,
  AnalystSalesVelocityResponse,
  AnalystAIInsightsResponse,
} from "../types";

export const dashboardService = {
  async getOverview(params?: DashboardRequestParams): Promise<DashboardOverviewResponse> {
    const response = await apiClient.get<DashboardOverviewResponse>("/dashboard/overview", {
      params,
    });
    return response.data;
  },

  async getVisitStatistics(params?: DashboardRequestParams): Promise<VisitStatisticsResponse> {
    const response = await apiClient.get<VisitStatisticsResponse>("/dashboard/visits", {
      params,
    });
    return response.data;
  },

  async getPipelineSummary(params?: DashboardRequestParams): Promise<PipelineSummaryResponse> {
    const response = await apiClient.get<PipelineSummaryResponse>("/dashboard/pipeline", {
      params,
    });
    return response.data;
  },

  async getTopAccounts(params?: DashboardRequestParams): Promise<TopAccountsResponse> {
    const response = await apiClient.get<TopAccountsResponse>("/dashboard/top-accounts", {
      params,
    });
    return response.data;
  },

  async getTopSalesRep(params?: DashboardRequestParams): Promise<TopSalesRepResponse> {
    const response = await apiClient.get<TopSalesRepResponse>("/dashboard/top-sales-rep", {
      params,
    });
    return response.data;
  },

  async getRecentActivities(params?: DashboardRequestParams): Promise<RecentActivitiesResponse> {
    const response = await apiClient.get<RecentActivitiesResponse>("/dashboard/recent-activities", {
      params,
    });
    return response.data;
  },

  async getActivityTrends(params?: DashboardRequestParams): Promise<ActivityTrendsResponse> {
    const response = await apiClient.get<ActivityTrendsResponse>("/dashboard/activity-trends", {
      params,
    });
    return response.data;
  },

  // ============================================================================
  // Super Admin Dashboard Methods
  // ============================================================================

  async getSuperAdminUsersByRole(
    params?: DashboardRequestParams
  ): Promise<SuperAdminUsersByRoleResponse> {
    const response = await apiClient.get<SuperAdminUsersByRoleResponse>(
      "/dashboard/super-admin/users-by-role",
      { params }
    );
    return response.data;
  },

  async getSuperAdminSystemActivity(
    params?: DashboardRequestParams
  ): Promise<SuperAdminSystemActivityResponse> {
    const response = await apiClient.get<SuperAdminSystemActivityResponse>(
      "/dashboard/super-admin/system-activity",
      { params }
    );
    return response.data;
  },

  async getSuperAdminAIUsage(
    params?: DashboardRequestParams
  ): Promise<SuperAdminAIUsageResponse> {
    const response = await apiClient.get<SuperAdminAIUsageResponse>(
      "/dashboard/super-admin/ai-usage",
      { params }
    );
    return response.data;
  },

  async getSuperAdminDataGrowth(
    params?: DashboardRequestParams
  ): Promise<SuperAdminDataGrowthResponse> {
    const response = await apiClient.get<SuperAdminDataGrowthResponse>(
      "/dashboard/super-admin/data-growth",
      { params }
    );
    return response.data;
  },

  async getSuperAdminErrorSummary(
    params?: DashboardRequestParams
  ): Promise<SuperAdminErrorSummaryResponse> {
    const response = await apiClient.get<SuperAdminErrorSummaryResponse>(
      "/dashboard/super-admin/error-summary",
      { params }
    );
    return response.data;
  },

  // ============================================================================
  // Admin Dashboard Methods
  // ============================================================================

  async getAdminTotalLeads(
    params?: DashboardRequestParams
  ): Promise<AdminTotalLeadsResponse> {
    const response = await apiClient.get<AdminTotalLeadsResponse>(
      "/dashboard/admin/total-leads",
      { params }
    );
    return response.data;
  },

  async getAdminPipelineValue(
    params?: DashboardRequestParams
  ): Promise<AdminPipelineValueResponse> {
    const response = await apiClient.get<AdminPipelineValueResponse>(
      "/dashboard/admin/pipeline-value",
      { params }
    );
    return response.data;
  },

  async getAdminPendingApprovals(
    params?: DashboardRequestParams
  ): Promise<AdminPendingApprovalsResponse> {
    const response = await apiClient.get<AdminPendingApprovalsResponse>(
      "/dashboard/admin/pending-approvals",
      { params }
    );
    return response.data;
  },

  async getAdminTaskOverdue(
    params?: DashboardRequestParams
  ): Promise<AdminTaskOverdueResponse> {
    const response = await apiClient.get<AdminTaskOverdueResponse>(
      "/dashboard/admin/task-overdue",
      { params }
    );
    return response.data;
  },

  // ============================================================================
  // Sales Manager Dashboard Methods
  // ============================================================================

  async getSalesManagerPipelineFunnel(
    params?: DashboardRequestParams
  ): Promise<SalesManagerPipelineFunnelResponse> {
    const response = await apiClient.get<SalesManagerPipelineFunnelResponse>(
      "/dashboard/sales-manager/pipeline-funnel",
      { params }
    );
    return response.data;
  },

  async getSalesManagerTargetVsActual(
    params?: DashboardRequestParams
  ): Promise<SalesManagerTargetVsActualResponse> {
    const response = await apiClient.get<SalesManagerTargetVsActualResponse>(
      "/dashboard/sales-manager/target-vs-actual",
      { params }
    );
    return response.data;
  },

  async getSalesManagerVisitCompletion(
    params?: DashboardRequestParams
  ): Promise<SalesManagerVisitCompletionResponse> {
    const response = await apiClient.get<SalesManagerVisitCompletionResponse>(
      "/dashboard/sales-manager/visit-completion",
      { params }
    );
    return response.data;
  },

  async getSalesManagerDealsAtRisk(
    params?: DashboardRequestParams
  ): Promise<SalesManagerDealsAtRiskResponse> {
    const response = await apiClient.get<SalesManagerDealsAtRiskResponse>(
      "/dashboard/sales-manager/deals-at-risk",
      { params }
    );
    return response.data;
  },

  async getSalesManagerTeamDraftApprovals(): Promise<SalesManagerTeamDraftApprovalsResponse> {
    const response = await apiClient.get<SalesManagerTeamDraftApprovalsResponse>(
      "/dashboard/sales-manager/team-draft-approvals"
    );
    return response.data;
  },

  // ============================================================================
  // Sales Dashboard Methods
  // ============================================================================

  async getSalesTodayTasks(
    params?: DashboardRequestParams
  ): Promise<SalesTodayTasksResponse> {
    const response = await apiClient.get<SalesTodayTasksResponse>(
      "/dashboard/sales/today-tasks",
      { params }
    );
    return response.data;
  },

  async getSalesAssignedLeads(
    params?: DashboardRequestParams
  ): Promise<SalesAssignedLeadsResponse> {
    const response = await apiClient.get<SalesAssignedLeadsResponse>(
      "/dashboard/sales/assigned-leads",
      { params }
    );
    return response.data;
  },

  async getSalesUpcomingVisits(
    params?: DashboardRequestParams
  ): Promise<SalesUpcomingVisitsResponse> {
    const response = await apiClient.get<SalesUpcomingVisitsResponse>(
      "/dashboard/sales/upcoming-visits",
      { params }
    );
    return response.data;
  },

  async getSalesReminders(
    params?: DashboardRequestParams
  ): Promise<SalesRemindersResponse> {
    const response = await apiClient.get<SalesRemindersResponse>(
      "/dashboard/sales/reminders",
      { params }
    );
    return response.data;
  },

  // ============================================================================
  // Analyst Dashboard Methods
  // ============================================================================

  async getAnalystRevenueTrend(
    params?: DashboardRequestParams
  ): Promise<AnalystRevenueTrendResponse> {
    const response = await apiClient.get<AnalystRevenueTrendResponse>(
      "/dashboard/analyst/revenue-trend",
      { params }
    );
    return response.data;
  },

  async getAnalystConversionRate(
    params?: DashboardRequestParams
  ): Promise<AnalystConversionRateResponse> {
    const response = await apiClient.get<AnalystConversionRateResponse>(
      "/dashboard/analyst/conversion-rate",
      { params }
    );
    return response.data;
  },

  async getAnalystSalesVelocity(
    params?: DashboardRequestParams
  ): Promise<AnalystSalesVelocityResponse> {
    const response = await apiClient.get<AnalystSalesVelocityResponse>(
      "/dashboard/analyst/sales-velocity",
      { params }
    );
    return response.data;
  },

  async getAnalystAIInsights(
    params?: DashboardRequestParams
  ): Promise<AnalystAIInsightsResponse> {
    const response = await apiClient.get<AnalystAIInsightsResponse>(
      "/dashboard/analyst/ai-insights",
      { params }
    );
    return response.data;
  },

};

