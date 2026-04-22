export interface ActivityType {
  id: string;
  name: string;
  code: string;
  description?: string;
  icon?: string;
  badge_color: "default" | "secondary" | "destructive" | "outline" | "success" | "warning" | "active";
  status: "active" | "inactive";
  order: number;
  created_at: string;
  updated_at: string;
  activity_count?: number;
}

export interface ActivityTypeResponse {
  success: boolean;
  data: ActivityType;
}

export interface ListActivityTypesResponse {
  success: boolean;
  data: ActivityType[];
}

