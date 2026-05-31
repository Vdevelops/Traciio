/**
 * Pipeline Management Type Definitions
 * Domain types for deals, stages, forecasts, and summaries
 */

// Pipeline Stage Types
export interface PipelineStage {
  readonly id: string;
  readonly name: string;
  readonly code: string;
  readonly description?: string;
  readonly order: number;
  readonly color?: string;
  readonly is_active: boolean;
  readonly is_won: boolean;
  readonly is_lost: boolean;
  readonly probability?: number;
  readonly requirements?: string;
  readonly created_at: string;
  readonly updated_at: string;
}

// Deal Reference Types
export interface AccountRef {
  readonly id: string;
  readonly name: string;
}

export interface ContactRef {
  readonly id: string;
  readonly name: string;
  readonly email?: string;
  readonly phone?: string;
}

export interface UserRef {
  readonly id: string;
  readonly name: string;
  readonly email?: string;
  readonly avatar_url?: string;
}

// Deal Product Item Types
export interface DealProductItem {
  readonly id?: string;
  readonly deal_id?: string;
  readonly product_id: string;
  readonly product_name?: string;
  readonly product_sku?: string;
  readonly unit_price: number;
  readonly unit_price_formatted?: string;
  readonly quantity: number;
  readonly discount_amount?: number;
  readonly discount_amount_formatted?: string;
  readonly subtotal?: number;
  readonly subtotal_formatted?: string;
  readonly notes?: string;
}

// Deal Types
export interface Deal {
  readonly id: string;
  readonly title: string;
  readonly description?: string;
  readonly account_id: string;
  readonly account?: AccountRef;
  readonly contact_id?: string;
  readonly contact?: ContactRef;
  readonly stage_id: string;
  readonly stage?: PipelineStage;
  readonly value: number;
  readonly value_formatted?: string;
  readonly probability: number;
  readonly expected_close_date?: string;
  readonly actual_close_date?: string;
  readonly assigned_to?: string;
  readonly assigned_user?: UserRef;
  readonly lead_id?: string;
  readonly status: "open" | "won" | "lost";
  readonly source?: string;
  readonly budget_confirmed?: boolean;
  readonly authority_confirmed?: boolean;
  readonly need_confirmed?: boolean;
  readonly timeline_confirmed?: boolean;
  readonly qualification_snapshot?: {
    readonly budget_target_amount?: number;
    readonly budget_target_currency?: string;
    readonly budget_confirmed?: boolean;
    readonly budget_notes?: string;
    readonly authority_target_person?: string;
    readonly authority_target_role?: string;
    readonly authority_confirmed?: boolean;
    readonly authority_notes?: string;
    readonly need_target_products?: ReadonlyArray<{
      readonly product_id?: string;
      readonly product_name?: string;
    }>;
    readonly need_priority_level?: string;
    readonly need_confirmed?: boolean;
    readonly need_notes?: string;
    readonly timeline_target_date?: string;
    readonly timeline_flexibility?: string;
    readonly timeline_confirmed?: boolean;
    readonly timeline_notes?: string;
    readonly qualification_score?: number;
    readonly qualification_status?: string;
  };
  readonly notes?: string;
  readonly created_by?: string;
  readonly created_at: string;
  readonly updated_at: string;
  readonly deleted_at?: string;

  // New: deal line items
  readonly product_items?: DealProductItem[];
}

// Pipeline Summary Types
export interface StageSummary {
  readonly stage_id: string;
  readonly stage_name: string;
  readonly stage_code: string;
  readonly deal_count: number;
  readonly total_value: number;
  readonly total_value_formatted: string;
}

export interface PipelineSummary {
  readonly total_deals: number;
  readonly open_deals: number;
  readonly won_deals: number;
  readonly lost_deals: number;
  readonly total_value: number;
  readonly total_value_formatted: string;
  readonly won_value?: number;
  readonly won_value_formatted?: string;
  readonly lost_value?: number;
  readonly lost_value_formatted?: string;
  readonly open_value?: number;
  readonly open_value_formatted?: string;
  readonly average_deal_size: number;
  readonly average_deal_size_formatted: string;
  readonly conversion_rate: number;
  readonly by_stage?: StageSummary[];
  readonly stages?: ReadonlyArray<{
    readonly stage_id: string;
    readonly stage_name: string;
    readonly deal_count: number;
    readonly total_value: number;
    readonly total_value_formatted: string;
  }>;
}

// Stage Performance Statistics
export interface StagePerformance {
  readonly stage_id: string;
  readonly stage_name: string;
  readonly stage_color?: string;
  readonly deal_count: number;
  readonly total_value: number;
  readonly total_value_formatted: string;
  readonly average_deal_size: number;
  readonly average_deal_size_formatted: string;
  readonly conversion_rate: number;
}

// Forecast Types
export interface ForecastPeriod {
  readonly type: "month" | "quarter" | "year";
  readonly start: string;
  readonly end: string;
}

export interface ForecastDeal {
  readonly id: string;
  readonly title: string;
  readonly account_id: string;
  readonly account_name: string;
  readonly contact_id?: string;
  readonly contact_name?: string;
  readonly stage_name: string;
  readonly value: number;
  readonly value_formatted: string;
  readonly probability: number;
  readonly expected_close_date: string;
  readonly weighted_value: number;
  readonly weighted_value_formatted: string;
}

export interface Forecast {
  readonly period: ForecastPeriod;
  readonly expected_revenue: number;
  readonly expected_revenue_formatted: string;
  readonly weighted_revenue: number;
  readonly weighted_revenue_formatted: string;
  readonly deals: ForecastDeal[];
}

export interface RevenueForecast {
  readonly period: "month" | "quarter" | "year";
  readonly expected_revenue: number;
  readonly expected_revenue_formatted: string;
  readonly weighted_revenue: number;
  readonly weighted_revenue_formatted: string;
  readonly deals_count: number;
  readonly deals?: ReadonlyArray<{
    readonly id: string;
    readonly title: string;
    readonly account_name?: string;
    readonly stage_name?: string;
    readonly value: number;
    readonly value_formatted: string;
    readonly probability: number;
    readonly weighted_value: number;
    readonly weighted_value_formatted: string;
    readonly expected_close_date?: string;
  }>;
}

// Deal Filter Options
export interface DealFilters {
  readonly stage_id?: string;
  readonly account_id?: string;
  readonly assigned_to?: string;
  readonly search?: string;
  readonly min_value?: number;
  readonly max_value?: number;
  readonly date_from?: string;
  readonly date_to?: string;
}

// Stage with statistics
export interface StageWithStats extends PipelineStage {
  readonly deals_count?: number;
  readonly total_value?: number;
  readonly total_value_formatted?: string;
}

// API Response Types
export interface ListPipelineStagesResponse {
  success: boolean;
  data: PipelineStage[];
  meta?: {
    filters?: Record<string, unknown>;
  };
  timestamp: string;
  request_id: string;
}

export interface PipelineStageResponse {
  success: boolean;
  data: PipelineStage;
  timestamp: string;
  request_id: string;
}

export interface ListDealsResponse {
  success: boolean;
  data: Deal[];
  meta: {
    pagination: {
      page: number;
      per_page: number;
      total: number;
      total_pages: number;
      has_next: boolean;
      has_prev: boolean;
    };
    filters?: Record<string, unknown>;
  };
  timestamp: string;
  request_id: string;
}

export interface DealResponse {
  success: boolean;
  data: Deal;
  timestamp: string;
  request_id: string;
}

export interface PipelineSummaryResponse {
  success: boolean;
  data: PipelineSummary;
  timestamp: string;
  request_id: string;
}

export interface ForecastResponse {
  success: boolean;
  data: Forecast;
  timestamp: string;
  request_id: string;
}

// Deal History Types
export interface DealHistory {
  readonly id: string;
  readonly deal_id: string;
  readonly from_stage_id?: string;
  readonly from_stage?: PipelineStage;
  readonly from_stage_name?: string;
  readonly to_stage_id: string;
  readonly to_stage?: PipelineStage;
  readonly to_stage_name: string;
  readonly from_probability: number;
  readonly to_probability: number;
  readonly days_in_prev_stage?: number;
  readonly changed_by: string;
  readonly changed_by_user?: UserRef;
  readonly changed_at: string;
  readonly reason?: string;
  readonly notes?: string;
  readonly created_at: string;
}

export interface DealHistoryResponse {
  success: boolean;
  data: DealHistory[];
  timestamp: string;
  request_id: string;
}
