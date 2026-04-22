/**
 * Brick Analytics Types
 * Types for brick performance metrics and analytics
 */

export interface BrickPerformanceMetrics {
  brick_id: string;
  brick_name: string;
  brick_code: string;
  manager_name?: string;
  manager_id?: string;

  // Target & Achievement
  monthly_target: number;
  target_achieved: number;
  achievement_percentage: number;
  target_remaining: number;

  // Sales Team
  total_sales: number;
  active_sales: number;

  // Pipeline Metrics
  total_deals: number;
  open_deals: number;
  won_deals: number;
  lost_deals: number;
  total_deal_value: number;
  won_deal_value: number;
  win_rate: number;
  average_deal_size: number;

  // Visit Activity
  total_visits: number;
  visits_this_month: number;
  average_visits_per_sales: number;

  // Accounts
  total_accounts: number;
  active_accounts: number;
  new_accounts_this_month: number;

  // Revenue
  total_revenue: number;
  revenue_this_month: number;
  revenue_growth_percentage: number;
}

export interface BrickPerformanceResponse {
  success: boolean;
  data: BrickPerformanceMetrics;
  meta?: {
    period_start?: string;
    period_end?: string;
  };
}

export interface BrickPerformanceListResponse {
  success: boolean;
  data: BrickPerformanceMetrics[];
  meta?: {
    period_start?: string;
    period_end?: string;
  };
}

