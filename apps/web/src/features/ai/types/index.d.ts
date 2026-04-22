export type InsightType = "visit_report" | "deal" | "contact" | "account";

export interface VisitReportInsight {
  summary: string;
  action_items: string[];
  sentiment: "positive" | "neutral" | "negative";
  key_points: string[];
  recommendations: string[];
}

export interface InsightResponse {
  type: InsightType;
  data: VisitReportInsight;
  tokens: number;
}

export interface AnalyzeVisitReportRequest {
  visit_report_id: string;
}

export interface AnalyzeVisitReportResponse {
  success: boolean;
  data: InsightResponse;
  timestamp: string;
  request_id: string;
}

export interface ChatMessage {
  role: "user" | "assistant";
  content: string;
}

export interface ChatRequest {
  message: string;
  context?: string;
  context_type?: "visit_report" | "deal" | "contact" | "account" | "lead";
  conversation_history?: ChatMessage[];
  model?: string;
  domain?: AIDomain;
}

export interface AISettingsResponse {
  id: string;
  enabled: boolean;
  provider: string;
  model: string;
  base_url?: string;
  api_key?: string;
  data_privacy: AIDataPrivacySettings;
  timezone: string;
  created_at: string;
  updated_at: string;
}

export interface ChatResponse {
  message: string;
  tokens: number;
}

export interface ChatAPIResponse {
  success: boolean;
  data: ChatResponse;
  timestamp: string;
  request_id: string;
}

/**
 * Data privacy settings controlling which CRM data AI can access.
 * Grouped by domain category to minimize token usage via modular context loading.
 */
export interface AIDataPrivacySettings {
  // Sales domain
  allow_leads: boolean;
  allow_deals: boolean;
  allow_visit_reports: boolean;
  allow_activities: boolean;
  allow_tasks: boolean;
  allow_schedule: boolean;
  allow_pipelines: boolean;

  // Customer domain
  allow_accounts: boolean;
  allow_contacts: boolean;

  // Inventory domain
  allow_products: boolean;

  // Analytics domain
  allow_sales_performance: boolean;
  allow_product_analysis: boolean;
  allow_reports: boolean;

  // Management domain
  allow_users: boolean;
  allow_roles: boolean;
  allow_groups: boolean;
  allow_brick_management: boolean;
  allow_target: boolean;

  // Route Optimization domain
  allow_route_optimization: boolean;
}

export interface AISettings {
  data_privacy: AIDataPrivacySettings;
  enabled: boolean;
}

/**
 * CRM domain categories for modular AI prompt routing.
 * Each domain loads only its relevant context to minimize token usage.
 */
export type AIDomain =
  | "route_optimization"
  | "sales"
  | "inventory"
  | "customers"
  | "analytics"
  | "management"
  | "general";

/**
 * Domain module definition for the modular prompt system.
 * Contains metadata and context-building instructions per domain.
 */
export interface AIDomainModule {
  id: AIDomain;
  label: string;
  description: string;
  /** Data privacy keys required for this domain to function */
  requiredPrivacyKeys: (keyof AIDataPrivacySettings)[];
  /** Keywords used for intent detection from user messages */
  intentKeywords: string[];
  /** Entity types this domain can perform CRUD on */
  supportedEntities: string[];
  /** Compact prompt fragment injected alongside core system prompt */
  promptFragment: string;
}

/**
 * AI Action Cards - structured actions embedded in AI responses.
 * Rendered as clickable cards below the message to navigate to pages or open entity details.
 */
export type AIActionType = "navigate" | "detail";

export type AIActionIcon =
  | "map"
  | "trending-up"
  | "package"
  | "users"
  | "bar-chart"
  | "settings"
  | "clipboard"
  | "calendar"
  | "target"
  | "file-text"
  | "user"
  | "building"
  | "phone";

export type AIActionEntity = "account" | "contact" | "lead" | "deal" | "visit" | "task";

export interface AIActionCard {
  type: AIActionType;
  label: string;
  description: string;
  /** Page URL for "navigate" type */
  url?: string;
  /** Entity type for "detail" type */
  entity?: AIActionEntity;
  /** Entity ID for "detail" type */
  entityId?: string;
  /** Icon identifier */
  icon?: AIActionIcon;
}

