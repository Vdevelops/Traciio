export interface LeadSource {
  id: string;
  name: string;
  code: string;
  description?: string;
  order: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  lead_count?: number;
}

export interface ListLeadSourcesResponse {
  success: boolean;
  data: LeadSource[];
  meta?: {
    pagination?: {
      current_page: number;
      per_page: number;
      total: number;
      total_pages: number;
    };
  };
}

export interface LeadSourceResponse {
  success: boolean;
  data: LeadSource;
}

export interface CreateLeadSourceRequest {
  name: string;
  code: string;
  description?: string;
  order?: number;
  is_active?: boolean;
}

export interface UpdateLeadSourceRequest {
  name?: string;
  code?: string;
  description?: string;
  order?: number;
  is_active?: boolean;
}

