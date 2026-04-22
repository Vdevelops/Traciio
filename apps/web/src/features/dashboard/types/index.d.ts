import type { routing } from "@/i18n/routing";

export type Locale = (typeof routing)["locales"][number];

declare global {
  interface IntlMessages {
    common: Record<string, string>;
    nav: Record<string, string>;
    dashboard: Record<string, string>;
    dashboardOverview: Record<string, string>;
    activityTrends: Record<string, string>;
    visitStatistics: Record<string, string>;
    pipelineSummary: Record<string, string>;
    recentActivities: Record<string, string>;
    topAccounts: Record<string, string>;
    topSalesReps: Record<string, string>;
    sidebar: Record<string, string>;
    notFound: Record<string, string>;
  }
}

export interface DashboardOverview {
  period: {
    type: string;
    start: string;
    end: string;
  };
  visit_stats?: {
    total: number;
    completed: number;
    pending: number;
    approved: number;
    rejected: number;
    change_percent: number;
  };
  account_stats?: {
    total: number;
    active: number;
    inactive: number;
    change_percent: number;
  };
  activity_stats?: {
    total: number;
    visits: number;
    calls: number;
    emails: number;
    change_percent: number;
  };
  target?: {
    target_amount: number;
    target_amount_formatted: string;
    achieved_amount: number;
    achieved_amount_formatted: string;
    progress_percent: number;
    change_percent: number;
  };
  deals?: {
    total_deals: number;
    open_deals: number;
    won_deals: number;
    lost_deals: number;
    total_value: number;
    total_value_formatted: string;
    change_percent: number;
  };
  revenue?: {
    total_revenue: number;
    total_revenue_formatted: string;
    change_percent: number;
  };
  leads_by_source?: {
    total: number;
    by_source: Array<{
      source: string;
      count: number;
    }>;
  };
  upcoming_tasks?: Array<{
    id: string;
    title: string;
    priority: string;
    status: string;
    due_date?: string | null;
  }>;
  pipeline_stages?: Array<{
    stage_id: string;
    stage_name: string;
    stage_code: string;
    deal_count: number;
    percentage: number;
  }>;
  lead_stats?: {
    total: number;
    new: number;
    contacted: number;
    qualified: number;
    converted: number;
    lost: number;
    change_percent: number;
  };
}

export interface VisitStatistics {
  period: {
    start: string;
    end: string;
  };
  total: number;
  completed: number;
  pending: number;
  approved: number;
  rejected: number;
  by_status: Record<string, number>;
  by_date: Array<{
    date: string;
    count: number;
    completed: number;
    approved: number;
    pending: number;
    rejected: number;
  }>;
  change_percent: number;
}

export interface PipelineSummaryStage {
  stage_id: string;
  stage_name: string;
  stage_code: string;
  stage_color: string;
  deal_count: number;
  total_value: number;
  total_value_formatted: string;
  percentage: number;
}

export interface PipelineSummary {
  total_deals: number;
  total_value: number;
  won_deals: number;
  lost_deals: number;
  open_deals: number;
  by_stage: PipelineSummaryStage[];
}

export interface TopAccount {
  account: {
    id: string;
    name: string;
  };
  visit_count: number;
  activity_count: number;
  last_visit_date?: string;
}

export interface TopSalesRep {
  sales_rep: {
    id: string;
    name: string;
    email: string;
  };
  visit_count: number;
  account_count: number;
  activity_count: number;
  deals_closed: number;
  actual_revenue: number;
  actual_revenue_formatted: string;
  target_amount: number;
  target_amount_formatted: string;
  target_achievement_percent: number;
}

export interface RecentActivity {
  id: string;
  type: string;
  description: string;
  account?: {
    id: string;
    name: string;
  };
  contact?: {
    id: string;
    name: string;
  };
  user: {
    id: string;
    name: string;
  };
  timestamp: string;
}

export interface DashboardOverviewResponse {
  success: boolean;
  data: DashboardOverview;
  timestamp: string;
  request_id: string;
}

export interface VisitStatisticsResponse {
  success: boolean;
  data: VisitStatistics;
  timestamp: string;
  request_id: string;
}

export interface PipelineSummaryResponse {
  success: boolean;
  data: PipelineSummary;
  timestamp: string;
  request_id: string;
}

export interface TopAccountsResponse {
  success: boolean;
  data: TopAccount[];
  timestamp: string;
  request_id: string;
}

export interface TopSalesRepResponse {
  success: boolean;
  data: TopSalesRep[];
  timestamp: string;
  request_id: string;
}

export interface RecentActivitiesResponse {
  success: boolean;
  data: RecentActivity[];
  timestamp: string;
  request_id: string;
}

export interface ActivityTrends {
  period: {
    start: string;
    end: string;
  };
  by_date: Array<{
    date: string;
    visits: number;
    calls: number;
    emails: number;
    total: number;
  }>;
}

export interface ActivityTrendsResponse {
  success: boolean;
  data: ActivityTrends;
  timestamp: string;
  request_id: string;
}

export interface DashboardRequestParams {
  start_date?: string;
  end_date?: string;
  period?: "today" | "week" | "month" | "year";
  limit?: number;
  offset?: number;
  days?: number;
  status?: string;
  team_id?: string;
}

// ============================================================================
// Super Admin Dashboard Types
// ============================================================================

export interface SuperAdminUsersByRole {
  users_by_role: Array<{
    role_code: string;
    role_name: string;
    total_users: number;
    active_users: number;
    inactive_users: number;
  }>;
  total_users: number;
  total_active: number;
  total_inactive: number;
}

export interface SuperAdminSystemActivity {
  activities: Array<{
    id: string;
    type: string;
    description: string;
    user_id?: string;
    user_name?: string;
    created_at: string;
  }>;
  total: number;
  recent_errors: number;
  recent_warnings: number;
}

export interface SuperAdminAIUsage {
  total_requests: number;
  requests_today: number;
  requests_this_week: number;
  requests_this_month: number;
  estimated_cost: {
    today: number;
    week: number;
    month: number;
    currency: string;
  };
  success_rate: number;
  fallback_rate: number;
  average_response_time: number;
}

export interface SuperAdminDataGrowth {
  accounts: {
    total: number;
    growth_percent: number;
    growth_count: number;
    period: string;
  };
  leads: {
    total: number;
    growth_percent: number;
    growth_count: number;
    period: string;
  };
  deals: {
    total: number;
    growth_percent: number;
    growth_count: number;
    period: string;
  };
}

export interface SuperAdminErrorSummary {
  total_errors: number;
  errors_today: number;
  errors_this_week: number;
  errors_this_month: number;
  failed_processes: Array<{
    process_name: string;
    failure_count: number;
    last_failure: string;
  }>;
  error_types: Array<{
    type: string;
    count: number;
  }>;
}

// ============================================================================
// Admin Dashboard Types
// ============================================================================

export interface AdminTotalLeads {
  today: {
    total: number;
    new: number;
    contacted: number;
    qualified: number;
  };
  this_month: {
    total: number;
    new: number;
    contacted: number;
    qualified: number;
    converted: number;
  };
  change_percent: number;
}

export interface AdminPipelineValue {
  total_value: number;
  total_value_formatted: string;
  open_deals_value: number;
  won_deals_value: number;
  lost_deals_value: number;
  change_percent: number;
}

export interface AdminPendingApprovals {
  total: number;
  visit_reports: number;
  expense_reports: number;
  other: number;
  items: Array<{
    id: string;
    type: string;
    title: string;
    submitted_by: string;
    submitted_at: string;
    priority: string;
  }>;
}

export interface AdminTaskOverdue {
  total_overdue: number;
  critical_overdue: number;
  tasks: Array<{
    id: string;
    title: string;
    assigned_to: string;
    due_date: string;
    days_overdue: number;
    priority: string;
  }>;
}

// ============================================================================
// Sales Manager Dashboard Types
// ============================================================================

export interface SalesManagerPipelineFunnel {
  funnel: Array<{
    stage: string;
    count: number;
    percentage: number;
  }>;
  conversion_rate: number;
}

export interface SalesManagerTargetVsActual {
  period: string;
  target: {
    revenue: number;
    deals: number;
    visits: number;
  };
  actual: {
    revenue: number;
    deals: number;
    visits: number;
  };
  achievement: {
    revenue_percent: number;
    deals_percent: number;
    visits_percent: number;
  };
  gap: {
    revenue: number;
    deals: number;
    visits: number;
  };
}

export interface SalesManagerVisitCompletion {
  total_scheduled: number;
  completed: number;
  pending: number;
  missed: number;
  completion_rate: number;
  by_sales_rep: Array<{
    sales_rep_id: string;
    sales_rep_name: string;
    scheduled: number;
    completed: number;
    completion_rate: number;
  }>;
}

export interface SalesManagerDealsAtRisk {
  total_at_risk: number;
  deals: Array<{
    id: string;
    name: string;
    value: number;
    stage: string;
    days_without_activity: number;
    risk_reason: string;
    assigned_to: string;
    last_activity: string;
  }>;
}

// ============================================================================
// Sales Manager Team Draft Approvals Types
// Approval workflow: draft → submitted → approved / rejected
// ============================================================================

export interface DraftVisitItem {
  id: string;
  purpose: string;
  status: string;
  visit_date: string;
  assigned_to: string;
  created_at: string;
}

export interface DraftTaskItem {
  id: string;
  title: string;
  priority: string;
  due_date?: string;
  assigned_to: string;
  created_at: string;
}

export interface DraftScheduleItem {
  id: string;
  title: string;
  scheduled_at: string;
  assigned_to: string;
  created_at: string;
}

export interface DraftLeadItem {
  id: string;
  name: string;
  company: string;
  status: string;
  assigned_to: string;
  created_at: string;
}

export interface DraftPipelineItem {
  id: string;
  name: string;
  stage: string;
  value: number;
  assigned_to: string;
  created_at: string;
}

export interface SalesManagerTeamDraftApprovals {
  total: number;
  visits: {
    total: number;
    items: DraftVisitItem[];
  };
  tasks: {
    total: number;
    items: DraftTaskItem[];
  };
  schedules: {
    total: number;
    items: DraftScheduleItem[];
  };
  leads: {
    total: number;
    items: DraftLeadItem[];
  };
  pipeline: {
    total: number;
    items: DraftPipelineItem[];
  };
}

// ============================================================================
// Sales Dashboard Types
// ============================================================================

export interface SalesTodayTasks {
  total: number;
  completed: number;
  pending: number;
  overdue: number;
  tasks: Array<{
    id: string;
    title: string;
    due_date: string;
    due_time?: string;
    status: string;
    priority: string;
    related_to?: {
      type: string;
      id: string;
      name: string;
    };
  }>;
}

export interface SalesAssignedLeads {
  total: number;
  new: number;
  contacted: number;
  qualified: number;
  converted: number;
  leads: Array<{
    id: string;
    name: string;
    company: string;
    status: string;
    assigned_date: string;
    last_contact?: string;
  }>;
}

export interface SalesUpcomingVisits {
  total: number;
  today: number;
  this_week: number;
  next_week: number;
  visits: Array<{
    id: string;
    account_name: string;
    scheduled_date: string;
    scheduled_time?: string;
    purpose: string;
    status: string;
  }>;
}

export interface SalesReminders {
  total: number;
  unread: number;
  reminders: Array<{
    id: string;
    type: string;
    title: string;
    message: string;
    related_to?: {
      type: string;
      id: string;
    };
    created_at: string;
    read: boolean;
  }>;
}

// ============================================================================
// Analyst Dashboard Types
// ============================================================================

export interface AnalystRevenueTrend {
  period: string;
  total_revenue: number;
  trend: Array<{
    date: string;
    revenue: number;
    deals: number;
  }>;
  growth_percent: number;
  average_daily: number;
}

export interface AnalystConversionRate {
  period: string;
  total_leads: number;
  converted_leads: number;
  conversion_rate: number;
  by_source: Array<{
    source: string;
    leads: number;
    converted: number;
    conversion_rate: number;
  }>;
  trend: Array<{
    date: string;
    conversion_rate: number;
  }>;
}

export interface AnalystSalesVelocity {
  period: string;
  average_sales_cycle_days: number;
  average_deal_value: number;
  sales_velocity: number;
  by_stage: Array<{
    stage: string;
    average_days: number;
  }>;
}

export interface AnalystAIInsights {
  insights: Array<{
    id: string;
    type: string;
    title: string;
    description: string;
    impact: string;
    generated_at: string;
  }>;
  total_insights: number;
}

// ============================================================================
// Response Wrappers
// ============================================================================

export interface SuperAdminUsersByRoleResponse {
  success: boolean;
  data: SuperAdminUsersByRole;
  timestamp: string;
  request_id: string;
}

export interface SuperAdminSystemActivityResponse {
  success: boolean;
  data: SuperAdminSystemActivity;
  timestamp: string;
  request_id: string;
}

export interface SuperAdminAIUsageResponse {
  success: boolean;
  data: SuperAdminAIUsage;
  timestamp: string;
  request_id: string;
}

export interface SuperAdminDataGrowthResponse {
  success: boolean;
  data: SuperAdminDataGrowth;
  timestamp: string;
  request_id: string;
}

export interface SuperAdminErrorSummaryResponse {
  success: boolean;
  data: SuperAdminErrorSummary;
  timestamp: string;
  request_id: string;
}

export interface AdminTotalLeadsResponse {
  success: boolean;
  data: AdminTotalLeads;
  timestamp: string;
  request_id: string;
}

export interface AdminPipelineValueResponse {
  success: boolean;
  data: AdminPipelineValue;
  timestamp: string;
  request_id: string;
}

export interface AdminPendingApprovalsResponse {
  success: boolean;
  data: AdminPendingApprovals;
  timestamp: string;
  request_id: string;
}

export interface AdminTaskOverdueResponse {
  success: boolean;
  data: AdminTaskOverdue;
  timestamp: string;
  request_id: string;
}

export interface SalesManagerPipelineFunnelResponse {
  success: boolean;
  data: SalesManagerPipelineFunnel;
  timestamp: string;
  request_id: string;
}

export interface SalesManagerTargetVsActualResponse {
  success: boolean;
  data: SalesManagerTargetVsActual;
  timestamp: string;
  request_id: string;
}

export interface SalesManagerVisitCompletionResponse {
  success: boolean;
  data: SalesManagerVisitCompletion;
  timestamp: string;
  request_id: string;
}

export interface SalesManagerDealsAtRiskResponse {
  success: boolean;
  data: SalesManagerDealsAtRisk;
  timestamp: string;
  request_id: string;
}

export interface SalesManagerTeamDraftApprovalsResponse {
  success: boolean;
  data: SalesManagerTeamDraftApprovals;
  timestamp: string;
  request_id: string;
}

export interface SalesTodayTasksResponse {
  success: boolean;
  data: SalesTodayTasks;
  timestamp: string;
  request_id: string;
}

export interface SalesAssignedLeadsResponse {
  success: boolean;
  data: SalesAssignedLeads;
  timestamp: string;
  request_id: string;
}

export interface SalesUpcomingVisitsResponse {
  success: boolean;
  data: SalesUpcomingVisits;
  timestamp: string;
  request_id: string;
}

export interface SalesRemindersResponse {
  success: boolean;
  data: SalesReminders;
  timestamp: string;
  request_id: string;
}

export interface AnalystRevenueTrendResponse {
  success: boolean;
  data: AnalystRevenueTrend;
  timestamp: string;
  request_id: string;
}

export interface AnalystConversionRateResponse {
  success: boolean;
  data: AnalystConversionRate;
  timestamp: string;
  request_id: string;
}

export interface AnalystSalesVelocityResponse {
  success: boolean;
  data: AnalystSalesVelocity;
  timestamp: string;
  request_id: string;
}

export interface AnalystAIInsightsResponse {
  success: boolean;
  data: AnalystAIInsights;
  timestamp: string;
  request_id: string;
}

