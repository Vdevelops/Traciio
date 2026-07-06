import { z } from "zod";

export const activityTypeSchema = z.object({
  name: z.string().min(1, "Name is required").max(100, "Name must be at most 100 characters"),
  code: z.string().min(1, "Code is required").max(50, "Code must be at most 50 characters").regex(/^[a-z0-9_-]+$/, "Code must contain only lowercase letters, numbers, hyphens, and underscores"),
  description: z.string().optional(),
  icon: z.string().optional(),
  badge_color: z.enum(["default", "secondary", "destructive", "outline", "success", "warning", "active"]),
  status: z.enum(["active", "inactive"]),
  order: z.number().int().nonnegative("Order must be a non-negative integer"),
});

export type ActivityTypeFormData = z.infer<typeof activityTypeSchema>;
