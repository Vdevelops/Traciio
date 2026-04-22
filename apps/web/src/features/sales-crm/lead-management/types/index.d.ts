export interface Lead {
  id: string;
  first_name: string;
  last_name: string;
  company_name: string;
  email: string;
  phone: string;
  job_title: string;
  industry: string;
  lead_source: string;
  lead_status: "new" | "contacted" | "qualified" | "disqualified";
  lead_score?: number;
  assigned_to: string;
  assigned_user?: {
    id: string;
    name: string;
    email: string;
    avatar_url?: string;
  };
  account_id?: string;
  account?: {
    id: string;
    name: string;
  };
  contact_id?: string;
  contact?: {
    id: string;
    name: string;
    email: string;
    phone: string;
  };
  converted_at?: string;
  converted_by?: string;
  notes: string;
  address: string;
  city: string;
  province: string;
  postal_code: string;
  country: string;
  website: string;
  // BANT Qualification fields - REMOVED
  probability?: number;
  estimated_value?: number;
  // Optional BANT fields (used by UI components)
  budget_confirmed?: boolean;
  budget_amount?: number;
  authority_confirmed?: boolean;
  authority_person?: string;
  need_confirmed?: boolean;
  need_description?: string;
  timeline_confirmed?: boolean;
  expected_close_date?: string;
  // Optional opportunity relation (used by dashboard)
  opportunity?: {
    id?: string;
    value?: number;
    value_formatted?: string;
    stage_id?: string;
    status?: "open" | "won" | "lost";
  };
  created_by: string;
  created_at: string;
  updated_at: string;
}

export interface ListLeadsResponse {
  success: boolean;
  data: Lead[];
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
    sort?: {
      field: string;
      order: "asc" | "desc";
    };
  };
  timestamp: string;
  request_id: string;
}

export interface LeadResponse {
  success: boolean;
  data: Lead;
  timestamp: string;
  request_id: string;
}

export interface LeadFormDataResponse {
  success: boolean;
  data: {
    lead_sources: Array<{ value: string; label: string }>;
    lead_statuses: Array<{ value: string; label: string }>;
    users: Array<{ id: string; name: string; email: string }>;
    industries: string[];
    provinces: string[];
    defaults: {
      country: string;
      lead_status: string;
    };
  };
  timestamp: string;
  request_id: string;
}

export interface LeadAnalyticsResponse {
  success: boolean;
  data: {
    total_leads: number;
    by_status: Array<{
      status: string;
      count: number;
      percentage: number;
    }>;
    by_source: Array<{
      source: string;
      count: number;
      percentage: number;
    }>;
    conversion_rate: number;
    average_score: number;
    trends?: Array<{
      date: string;
      leads: number;
      converted: number;
    }>;
  };
  timestamp: string;
  request_id: string;
}

export interface ConvertLeadResponse {
  success: boolean;
  data: {
    lead: Lead;
    opportunity: unknown; // DealResponse
    account?: unknown; // AccountResponse
    contact?: unknown; // ContactResponse
  };
  timestamp: string;
  request_id: string;
}

export interface CreateAccountFromLeadResponse {
  success: boolean;
  data: {
    lead: Lead;
    account: unknown; // AccountResponse
    contact?: unknown; // ContactResponse
  };
  timestamp: string;
  request_id: string;
}

