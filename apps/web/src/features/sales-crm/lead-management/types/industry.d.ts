export interface Industry {
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

export interface ListIndustriesResponse {
  success: boolean;
  data: Industry[];
  meta?: {
    pagination?: {
      current_page: number;
      per_page: number;
      total: number;
      total_pages: number;
    };
  };
}

export interface IndustryResponse {
  success: boolean;
  data: Industry;
}

export interface CreateIndustryRequest {
  name: string;
  code: string;
  description?: string;
  order?: number;
  is_active?: boolean;
}

export interface UpdateIndustryRequest {
  name?: string;
  code?: string;
  description?: string;
  order?: number;
  is_active?: boolean;
}

