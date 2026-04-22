import type { Group } from "@/features/master-data/group/types";
import type { User } from "@/features/master-data/user-management/types";
import type { Brick } from "@/features/master-data/brick/types";

export interface MonthlyTarget {
  id: string;
  group_id?: string;
  group?: Group;
  user_id?: string;
  user?: User;
  brick_id?: string;
  brick?: Brick;
  year: number;
  month: number;
  target_amount: number;
  created_at: string;
  updated_at: string;
}

export interface ListMonthlyTargetsResponse {
  success: boolean;
  data: MonthlyTarget[];
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
    additional?: {
      total_target_amount?: number;
    };
  };
  timestamp: string;
  request_id: string;
}

export interface MonthlyTargetResponse {
  success: boolean;
  data: MonthlyTarget;
  timestamp: string;
  request_id: string;
}

