import { z } from "zod";

export const createActivitySchema = z.object({
  activity_type_id: z.string().uuid("Invalid activity type ID"),
  account_id: z.string().uuid("Invalid account ID").optional(),
  contact_id: z.string().uuid("Invalid contact ID").optional(),
  deal_id: z.string().uuid("Invalid deal ID").optional(),
  lead_id: z.string().uuid("Invalid lead ID").optional(),
  description: z.string().min(3, "Description must be at least 3 characters"),
  timestamp: z.string().refine(
    (val) => {
      // Validate ISO 8601 datetime format (RFC3339)
      const date = new Date(val);
      return !isNaN(date.getTime()) && val.includes("T");
    },
    { message: "Invalid timestamp format (must be ISO 8601 datetime)" }
  ),
  metadata: z.record(z.string(), z.unknown()).optional(),
});

export type CreateActivityFormData = z.infer<typeof createActivitySchema>;

