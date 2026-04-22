export interface Group {
  id: string;
  name: string;
  code: string;
  description?: string;
  status: "active" | "inactive";
  created_at: string;
  updated_at: string;
}

export interface ListGroupsResponse {
  success: boolean;
  data: Group[];
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

export interface GroupResponse {
  success: boolean;
  data: Group;
  timestamp: string;
  request_id: string;
}

