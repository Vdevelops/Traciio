// Sales Overview types
export interface SalesPerformanceDetail {
  user_id: string;
  user?: {
    id: string;
    name: string;
    email: string;
    avatar_url?: string;
  };
  period_start: string;
  period_end: string;
  total_revenue: number; // in smallest currency unit (sen)
  total_revenue_formatted: string;
  won_deals: number;
  total_deals: number;
  lost_deals: number;
  open_deals: number;
  conversion_rate: number; // percentage
  average_deal_value: number;
  average_deal_value_formatted: string;
  visits_completed: number;
  tasks_completed: number;
  total_tasks: number;
  task_completion_rate: number; // percentage
  prospect_outcome?: ProspectOutcomeSummary;
}

export interface GetSalesPerformanceDetailRequest {
  period?: "today" | "week" | "month" | "year";
  start_date?: string; // YYYY-MM-DD
  end_date?: string; // YYYY-MM-DD
}

export interface PeriodComparison {
  revenue_change: number; // percentage
  deals_change: number; // percentage
  visits_change: number; // percentage
}

export interface SalesRepStatistics {
  total_revenue: number; // in smallest currency unit (sen)
  total_revenue_formatted: string;
  deals_closed: number;
  visits_completed: number;
  tasks_completed: number;
  conversion_rate: number; // percentage
  average_deal_value: number;
  average_deal_value_formatted: string;
  period_comparison?: PeriodComparison;
  prospect_outcome?: ProspectOutcomeSummary;
}

export interface SalesRepDetail {
  user_id: string;
  user?: {
    id: string;
    name: string;
    email: string;
    avatar_url?: string;
  };
  period_start?: string; // ISO 8601 date-time string
  period_end?: string; // ISO 8601 date-time string
  statistics?: SalesRepStatistics;
}

export interface GetSalesRepDetailRequest {
  period?: "today" | "week" | "month" | "year";
  start_date?: string; // YYYY-MM-DD
  end_date?: string; // YYYY-MM-DD
}

export interface SalesPerformanceListItem {
  user_id: string;
  user_name: string;
  user_email: string;
  avatar_url?: string;
  total_revenue: number;
  total_revenue_formatted: string;
  deals_closed: number;
  visits_completed: number;
  tasks_completed: number;
  conversion_rate: number;
  total_prospects?: number;
  won_prospects?: number;
  lost_prospects?: number;
  open_prospects?: number;
  prospect_conversion_rate?: number;
  top_won_reason?: string;
  top_lost_reason?: string;
  // Target fields
  target_amount?: number;
  target_amount_formatted?: string;
  target_achievement_percentage?: number;
  // Optional id for DataTable compatibility (mapped from user_id)
  id?: string;
}

export interface ProspectReasonBreakdown {
  reason: string;
  category?: string;
  count: number;
  percentage: number;
}

export interface ProspectOutcomeItem {
  id: string;
  type: "deal" | string;
  title: string;
  account_name?: string;
  status: "open" | "won" | "lost" | string;
  value: number;
  value_formatted: string;
  reason?: string;
  reason_category?: string;
  source?: string;
  created_at: string;
  closed_at?: string;
}

export interface ProspectOutcomeListItem extends ProspectOutcomeItem {
  sales_rep_id?: string;
  sales_rep_name?: string;
  sales_rep_email?: string;
  sales_rep_avatar_url?: string;
}

export interface ProspectOutcomeSummary {
  total_prospects: number;
  won_prospects: number;
  lost_prospects: number;
  open_prospects: number;
  prospect_conversion_rate: number;
  won_reasons: ProspectReasonBreakdown[];
  lost_reasons: ProspectReasonBreakdown[];
  recent_prospects?: ProspectOutcomeItem[];
}

export interface ListSalesPerformanceRequest {
  search?: string;
  start_date?: string; // YYYY-MM-DD
  end_date?: string; // YYYY-MM-DD
  page?: number;
  per_page?: number;
  sort_by?: "revenue" | "deals" | "visits" | "tasks" | "name";
  order?: "asc" | "desc";
}

export interface ListSalesPerformanceResponse {
  success: boolean;
  data: SalesPerformanceListItem[];
  meta: {
    pagination: {
      page: number;
      per_page: number;
      total: number;
      total_pages: number;
      has_next: boolean;
      has_prev: boolean;
    };
  };
}

export interface ListProspectOutcomesRequest {
  search?: string;
  start_date?: string; // YYYY-MM-DD
  end_date?: string; // YYYY-MM-DD
  page?: number;
  per_page?: number;
  sales_user_id?: string;
  status?: "open" | "won" | "lost";
}

export interface ListProspectOutcomesResponse {
  success: boolean;
  data: ProspectOutcomeListItem[];
  meta: {
    pagination: {
      page: number;
      per_page: number;
      total: number;
      total_pages: number;
      has_next: boolean;
      has_prev: boolean;
    };
  };
}

export interface MonthlySalesData {
  month: number;
  month_name: string;
  year: number;
  period_key?: string;
  period_label?: string;
  period_start: string;
  period_end: string;
  total_revenue: number;
  total_deals: number;
  total_visits: number;
  total_tasks: number;
  target_amount: number;
  change_rate: number;
}

export interface MonthlySalesOverviewData {
  trend_mode: "monthly" | "mom" | "rolling_30d" | "rolling_90d" | "qoq";
  monthly_data: MonthlySalesData[];
  total_revenue: number;
  total_deals: number;
  total_visits: number;
  total_tasks: number;
}

export interface MonthlySalesOverviewResponse {
  success: boolean;
  data: MonthlySalesOverviewData;
}

export interface Location {
  latitude: number;
  longitude: number;
  address?: string;
}

export interface AccountRef {
  id: string;
  name: string;
}

export interface SalesRepCheckInLocation {
  visit_number: number;
  visit_report_id: string;
  visit_date: string; // ISO 8601 date string
  check_in_time: string; // ISO 8601 date-time string
  location?: Location;
  account?: AccountRef;
  purpose: string;
}

export interface GetSalesRepCheckInLocationsRequest {
  start_date?: string; // YYYY-MM-DD
  end_date?: string; // YYYY-MM-DD
  page?: number;
  per_page?: number;
}

export interface SalesRepCheckInLocationsResponse {
  sales_rep?: {
    id: string;
    name: string;
    email: string;
    avatar_url?: string;
  };
  check_in_locations: SalesRepCheckInLocation[];
  total_visits: number; // int64 from backend
  period?: {
    start: string;
    end: string;
  };
}
