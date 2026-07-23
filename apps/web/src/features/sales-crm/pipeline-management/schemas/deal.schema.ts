import { z } from "zod";

const dealProductItemSchema = z.object({
  product_id: z.string().uuid("Invalid product ID"),
  quantity: z.number().int().min(1, "Quantity must be at least 1"),
  // Money in Rupiah on FE (will be converted to sen on submit)
  unit_price: z.number().min(0, "Unit price must be greater than or equal to 0"),
  discount_amount: z.number().min(0, "Discount must be greater than or equal to 0").optional(),
  notes: z.string().max(500).optional(),
});

export const dealSchema = z.object({
  title: z
    .string()
    .min(3, "Title must be at least 3 characters")
    .max(255, "Title must be less than 255 characters"),
  description: z.string().optional(),
  account_id: z.string().uuid("Account is required"),
  contact_id: z.string().uuid("Contact is required"),
  stage_id: z.string().uuid("Stage is required"),
  // Value is auto-calculated from product_items (in Rupiah on FE)
  value: z
    .number()
    .min(0, "Value must be greater than or equal to 0")
    .default(0),
  probability: z
    .number()
    .min(0, "Probability must be between 0 and 100")
    .max(100, "Probability must be between 0 and 100")
    .default(0),
  expected_close_date: z.string().optional().or(z.literal("")),
  assigned_to: z.string().uuid("Invalid assigned user ID").optional(),
  lead_id: z.string().uuid("Invalid lead ID").optional().or(z.literal("")),
  source: z.string().max(100).optional().or(z.literal("")),
  close_reason: z.string().max(500, "Close reason must be at most 500 characters").optional().or(z.literal("")),
  notes: z.string().optional(),
  budget_confirmed: z.boolean().optional(),
  authority_confirmed: z.boolean().optional(),
  need_confirmed: z.boolean().optional(),
  timeline_confirmed: z.boolean().optional(),
  qualification_snapshot: z.record(z.string(), z.unknown()).optional(),

  // New: product line items (optional)
  product_items: z.array(dealProductItemSchema).optional().default([]),
});

export type DealFormData = z.infer<typeof dealSchema>;

// Alias for create schema
export const createDealSchema = dealSchema;
export type CreateDealFormData = z.infer<typeof createDealSchema>;

// Update schema (all fields optional)
export const dealUpdateSchema = dealSchema.partial();
export const updateDealSchema = dealUpdateSchema;

export type DealUpdateData = z.infer<typeof dealUpdateSchema>;
export type UpdateDealFormData = DealUpdateData;

// Move stage schema (used by MoveStageModal)
export const moveStageSchema = z.object({
  to_stage_id: z.string().uuid("Invalid stage ID"),
  reason: z.string().max(255).optional().or(z.literal("")),
  notes: z.string().max(500).optional().or(z.literal("")),
  product_items: z.array(dealProductItemSchema).optional(),
});

export type MoveStageFormData = z.infer<typeof moveStageSchema>;

export const dealMoveSchema = z.object({
  deal_id: z.string().uuid(),
  stage_id: z.string().uuid(),
  order: z.number().optional(),
  reason: z.string().max(500).optional().or(z.literal("")),
  product_items: z.array(dealProductItemSchema).optional(),
});

export type DealMoveData = z.infer<typeof dealMoveSchema>;
