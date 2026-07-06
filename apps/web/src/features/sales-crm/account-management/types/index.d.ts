export interface Category {
  id: string;
  name: string;
  code: string;
  description?: string;
  badge_color: "default" | "secondary" | "outline" | "success" | "warning" | "active";
  status: "active" | "inactive";
  account_count?: number;
  created_at: string;
  updated_at: string;
}

export interface ContactRole {
  id: string;
  name: string;
  code: string;
  description?: string;
  badge_color: "default" | "secondary" | "outline" | "success" | "warning" | "active";
  status: "active" | "inactive";
  contact_count?: number;
  created_at: string;
  updated_at: string;
}

export interface Account {
  id: string;
  name: string;
  category_id: string;
  category?: Category;
  address?: string;
  city?: string;
  province?: string;
  phone?: string;
  email?: string;
  latitude?: number;
  longitude?: number;
  status: "active" | "inactive";
  assigned_to?: string;
  brick_id?: string;
  /** Pre-computed contact count from the API — avoids fetching contacts just for the count badge */
  contact_count?: number;
  created_at: string;
  updated_at: string;
}

export interface Contact {
  id: string;
  account_id: string;
  name: string;
  role_id: string;
  role?: ContactRole;
  phone?: string;
  email?: string;
  position?: string;
  notes?: string;
  created_at: string;
  updated_at: string;
}

export interface ListCategoriesResponse {
  success: boolean;
  data: Category[];
  timestamp: string;
  request_id: string;
}

export interface CategoryResponse {
  success: boolean;
  data: Category;
  timestamp: string;
  request_id: string;
}

export interface ListContactRolesResponse {
  success: boolean;
  data: ContactRole[];
  timestamp: string;
  request_id: string;
}

export interface ContactRoleResponse {
  success: boolean;
  data: ContactRole;
  timestamp: string;
  request_id: string;
}

export interface ListAccountsResponse {
  success: boolean;
  data: Account[];
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

export interface AccountResponse {
  success: boolean;
  data: Account;
  timestamp: string;
  request_id: string;
}

export interface ListContactsResponse {
  success: boolean;
  data: Contact[];
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

export interface ContactResponse {
  success: boolean;
  data: Contact;
  timestamp: string;
  request_id: string;
}
