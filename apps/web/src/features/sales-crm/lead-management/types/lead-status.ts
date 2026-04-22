export interface LeadStatus {
  id: string;
  name: string;
  code: string;
  description?: string;
  score: number;
  color: string;
  order: number;
  is_active: boolean;
  is_default: boolean;
  is_converted: boolean;
  created_at: string;
  updated_at: string;
  updated_by?: string;
  lead_count?: number;
}

export interface LeadStatusResponse {
  success: boolean;
  data: LeadStatus;
}

export interface ListLeadStatusesResponse {
  success: boolean;
  data: LeadStatus[];
  meta?: {
    pagination?: {
      page: number;
      per_page: number;
      total: number;
      total_pages: number;
    };
  };
}

export interface CreateLeadStatusRequest {
  name: string;
  code: string;
  description?: string;
  score: number;
  color?: string;
  order?: number;
  is_active?: boolean;
  is_default?: boolean;
  is_converted?: boolean;
}

export interface UpdateLeadStatusRequest {
  name?: string;
  code?: string;
  description?: string;
  score?: number;
  color?: string;
  order?: number;
  is_active?: boolean;
  is_default?: boolean;
  is_converted?: boolean;
}
