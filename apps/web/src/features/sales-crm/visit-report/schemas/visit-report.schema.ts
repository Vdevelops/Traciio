import { z } from "zod";

const locationSchema = z.object({
  latitude: z.number().min(-90).max(90),
  longitude: z.number().min(-180).max(180),
  address: z.string().optional(),
});

export const createVisitReportSchema = z.object({
  account_id: z.string().uuid("Invalid account ID").optional(),
  contact_id: z.string().uuid("Invalid contact ID").optional(),
  deal_id: z.string().uuid("Invalid deal ID").optional(),
  lead_id: z.string().uuid("Invalid lead ID").optional(),
  visit_date: z.string().regex(/^\d{4}-\d{2}-\d{2} \d{2}:\d{2}$/, "Invalid date format (YYYY-MM-DD HH:mm)"),
  purpose: z.string().min(3, "Purpose must be at least 3 characters"),
  notes: z.string().optional(),
  check_in_location: locationSchema.optional(),
  check_out_location: locationSchema.optional(),
  photos: z.array(z.string().url("Invalid photo URL")).optional(),
});

export const updateVisitReportSchema = z.object({
  account_id: z.string().uuid("Invalid account ID").optional(),
  contact_id: z.string().uuid("Invalid contact ID").optional(),
  deal_id: z.string().uuid("Invalid deal ID").optional(),
  lead_id: z.string().uuid("Invalid lead ID").optional(),
  visit_date: z.string().regex(/^\d{4}-\d{2}-\d{2} \d{2}:\d{2}$/, "Invalid date format (YYYY-MM-DD HH:mm)").optional(),
  purpose: z.string().min(3, "Purpose must be at least 3 characters").optional(),
  notes: z.string().optional(),
  check_in_location: locationSchema.optional(),
  check_out_location: locationSchema.optional(),
  photos: z.array(z.string().url("Invalid photo URL")).optional(),
  status: z.enum(["draft", "submitted"]).optional(),
});

export const checkInSchema = z.object({
  location: locationSchema,
});

export const checkOutSchema = z.object({
  location: locationSchema,
});

export const rejectSchema = z.object({
  reason: z.string().min(3, "Reason must be at least 3 characters"),
});

export const uploadPhotoSchema = z.object({
  photo_url: z.string().url("Invalid photo URL"),
});

export const submitVisitReportSchema = z.object({
  outcome: z.enum(["positive", "neutral", "negative", "very_positive"]),
  next_steps: z.string().min(3, "Next steps must be at least 3 characters").optional(),
});

export type CreateVisitReportFormData = z.infer<typeof createVisitReportSchema>;
export type UpdateVisitReportFormData = z.infer<typeof updateVisitReportSchema>;
export type CheckInFormData = z.infer<typeof checkInSchema>;
export type CheckOutFormData = z.infer<typeof checkOutSchema>;
export type RejectFormData = z.infer<typeof rejectSchema>;
export type UploadPhotoFormData = z.infer<typeof uploadPhotoSchema>;
export type SubmitVisitReportFormData = z.infer<typeof submitVisitReportSchema>;

