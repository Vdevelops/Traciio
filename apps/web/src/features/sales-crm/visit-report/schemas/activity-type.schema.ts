import { z } from "zod";

export const createActivityTypeSchema = z.object({
  name: z.string().min(2, "Name must be at least 2 characters").max(100, "Name must not exceed 100 characters"),
  code: z.string().min(2, "Code must be at least 2 characters").max(50, "Code must not exceed 50 characters"),
  description: z.string().optional(),
  icon: z.string().max(50, "Icon must not exceed 50 characters").optional(),
  badge_color: z.enum(["default", "secondary", "destructive", "outline", "success", "warning", "active"]).optional(),
  status: z.enum(["active", "inactive"]).optional(),
  order: z.number().int().min(0).optional(),
});

export const updateActivityTypeSchema = z.object({
  name: z.string().min(2, "Name must be at least 2 characters").max(100, "Name must not exceed 100 characters").optional(),
  code: z.string().min(2, "Code must be at least 2 characters").max(50, "Code must not exceed 50 characters").optional(),
  description: z.string().optional(),
  icon: z.string().max(50, "Icon must not exceed 50 characters").optional(),
  badge_color: z.enum(["default", "secondary", "destructive", "outline", "success", "warning", "active"]).optional(),
  status: z.enum(["active", "inactive"]).optional(),
  order: z.number().int().min(0).optional(),
});

export type CreateActivityTypeFormData = z.infer<typeof createActivityTypeSchema>;
export type UpdateActivityTypeFormData = z.infer<typeof updateActivityTypeSchema>;

