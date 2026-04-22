export interface ProfileUpdateRequest {
  name?: string;
  password?: string;
  current_password?: string;
  avatar_url?: string;
}

export interface ProfileStats {
  visits: number;
  deals: number;
  tasks: number;
  // Extended stats from settings-summary endpoint
  total_revenue?: number;
  deals_won?: number;
  deals_lost?: number;
deals_open?: number;
  total_revenue_formatted?: string;
  // Additional metrics for detailed stats display
  conversion_rate?: number; // percentage
  average_deal_value?: number;
  average_deal_value_formatted?: string;
}


export interface ProfileActivity {
  id: string;
  title: string;
  description: string;
  type: string;
  date: string;
  download_url?: string;
}

export interface ProfileTransaction {
  id: string;
  product: string;
  status: "pending" | "paid" | "failed";
  date: string;
  amount: number;
}

export interface ProfileData {
  user: {
    id: string;
    email: string;
    name: string;
    avatar_url?: string;
    role_id: string;
    role?: {
      id: string;
      name: string;
      code: string;
    };
    status: "active" | "inactive";
    created_at: string;
    updated_at: string;
  };
  stats: ProfileStats;
  activities: ProfileActivity[];
  transactions: ProfileTransaction[];
}

export interface ProfileResponse {
  success: boolean;
  data: ProfileData;
  timestamp: string;
  request_id: string;
}

export interface UserResponse {
  success: boolean;
  data: {
    id: string;
    email: string;
    name: string;
    avatar_url?: string;
    role_id: string;
    role?: {
      id: string;
      name: string;
      code: string;
    };
    status: "active" | "inactive";
    created_at: string;
    updated_at: string;
  };
  timestamp: string;
  request_id: string;
}

