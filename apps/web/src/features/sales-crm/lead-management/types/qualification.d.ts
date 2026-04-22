export interface LeadQualificationChecklist {
  id: string;
  lead_id: string;
  budget_target_amount: number;
  budget_target_currency: string;
  budget_confirmed: boolean;
  budget_notes?: string;
  authority_target_role?: string;
  authority_target_person?: string;
  authority_confirmed: boolean;
  authority_notes?: string;
  need_target_products: ProductInterest[];
  need_priority_level: "low" | "medium" | "high" | "critical";
  need_confirmed: boolean;
  need_notes?: string;
  timeline_target_date?: string;
  timeline_flexibility: "fixed" | "flexible" | "urgent";
  timeline_confirmed: boolean;
  timeline_notes?: string;
  qualification_score: number;
  qualification_status:
  | "pending"
  | "unqualified"
  | "cold"
  | "warm"
  | "qualified";
  bant_progress: BANTProgress;
  created_at: string;
  updated_at: string;
}

export interface ProductInterest {
  product_id: string;
  product_name: string;
  category_id?: string;
  category_name?: string;
}

export interface BANTProgress {
  budget: BANTItemProgress;
  authority: BANTItemProgress;
  need: BANTItemProgress;
  timeline: BANTItemProgress;
}

export interface BANTItemProgress {
  completed: boolean;
  score: number;
  max_score: number;
}

// Update request
export interface UpdateLeadQualificationRequest {
  budget_target_amount?: number;
  budget_notes?: string;
  budget_confirmed?: boolean;

  authority_target_role?: string;
  authority_target_person?: string;
  authority_notes?: string;
  authority_confirmed?: boolean;

  need_target_products?: ProductInterest[];
  need_priority_level?: "low" | "medium" | "high" | "critical";
  need_notes?: string;
  need_confirmed?: boolean;

  timeline_target_date?: string;
  timeline_flexibility?: "fixed" | "flexible" | "urgent";
  timeline_notes?: string;
  timeline_confirmed?: boolean;
}
