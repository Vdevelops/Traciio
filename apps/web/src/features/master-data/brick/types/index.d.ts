import type { User } from "@/features/master-data/user-management/types";

export interface Brick {
  id: string;
  name: string;
  code: string;
  description?: string;
  province: string;
  regency: string;
  district?: string;
  manager_id?: string;
  manager?: User;
  status: "active" | "inactive";
  sales_count?: number;
  created_at: string;
  updated_at: string;
}

export interface ListBricksResponse {
  success: boolean;
  data: Brick[];
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

export interface BrickResponse {
  success: boolean;
  data: Brick;
  timestamp: string;
  request_id: string;
}

export interface BrickTargetDistribution {
  id: string;
  brick_id: string;
  brick?: Brick;
  brick_target_id: string;
  brick_target?: {
    id: string;
    year: number;
    month: number;
    target_amount: number;
  };
  sales_user_id: string;
  sales_user?: User;
  distributed_amount: number;
  distributed_amount_formatted?: string;
  distributed_by: string;
  distributed_by_user?: User;
  distributed_at: string;
  created_at: string;
  updated_at: string;
}

export interface BrickTargetWithDistributions {
  brick_id: string;
  brick?: Brick;
  year: number;
  month: number;
  target: {
    id: string;
    target_amount: number;
    target_amount_formatted?: string;
  } | null;
  distributions: BrickTargetDistribution[];
  total_distributed: number;
  total_distributed_formatted?: string;
  remaining_amount: number;
  remaining_amount_formatted?: string;
}

