export interface Location {
  latitude: number;
  longitude: number;
  address?: string;
}

export interface VisitReport {
  id: string;
  account_id?: string;
  contact_id?: string;
  deal_id?: string;
  lead_id?: string;
  sales_rep_id: string;
  visit_date: string;
  check_in_time?: string;
  check_out_time?: string;
  check_in_location?: Location;
  check_out_location?: Location;
  purpose: string;
  notes: string;
  outcome?: "positive" | "neutral" | "negative" | "very_positive" | string;
  next_steps?: string;
  metadata?: Record<string, unknown>;
  photos?: string[];
  status: "pending" | "completed" | "draft" | "submitted" | "approved" | "rejected";
  approved_by?: string;
  approved_at?: string;
  rejection_reason?: string;
  created_at: string;
  updated_at: string;
  account?: {
    id: string;
    name: string;
  };
  contact?: {
    id: string;
    name: string;
  };
  deal?: {
    id: string;
    title: string;
    status?: string;
  };
  lead?: {
    id: string;
    first_name?: string;
    last_name?: string;
    company_name?: string;
    status?: string;
    lead_status?: string;
  };
  sales_rep?: {
    id: string;
    name: string;
  };
}

export interface ListVisitReportsResponse {
  success: boolean;
  data: VisitReport[];
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

export interface VisitReportResponse {
  success: boolean;
  data: VisitReport;
  timestamp: string;
  request_id: string;
}

export interface CreateVisitReportFormData {
  account_id?: string;
  contact_id?: string;
  deal_id?: string;
  lead_id?: string;
  visit_date: string;
  purpose: string;
  notes?: string;
  check_in_location?: Location;
  check_out_location?: Location;
  photos?: string[];
  metadata?: Record<string, unknown>;
}

export interface UpdateVisitReportFormData {
  account_id?: string;
  contact_id?: string;
  deal_id?: string;
  lead_id?: string;
  visit_date?: string;
  purpose?: string;
  notes?: string;
  check_in_location?: Location;
  check_out_location?: Location;
  photos?: string[];
  metadata?: Record<string, unknown>;
}

export interface CheckInFormData {
  location: Location;
}

export interface CheckOutFormData {
  location: Location;
}

export interface RejectFormData {
  reason: string;
}

export interface UploadPhotoFormData {
  photo_url: string;
}

export interface SubmitVisitReportFormData {
  outcome: "positive" | "neutral" | "negative" | "very_positive";
  next_steps?: string;
}
