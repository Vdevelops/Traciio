import { z } from "zod";

export const createCategorySchema = z.object({
  name: z.string().min(3, "Name must be at least 3 characters"),
  code: z.string().min(3, "Code must be at least 3 characters"),
  description: z.string().optional(),
  badge_color: z.enum(["default", "secondary", "destructive", "outline", "success", "warning", "active"]).optional().default("outline"),
  status: z.enum(["active", "inactive"]).optional().default("active"),
});

export const updateCategorySchema = z.object({
  name: z.string().min(3, "Name must be at least 3 characters").optional(),
  code: z.string().min(3, "Code must be at least 3 characters").optional(),
  description: z.string().optional(),
  badge_color: z.enum(["default", "secondary", "destructive", "outline", "success", "warning", "active"]).optional(),
  status: z.enum(["active", "inactive"]).optional(),
});

export type CreateCategoryFormData = z.infer<typeof createCategorySchema>;
export type UpdateCategoryFormData = z.infer<typeof updateCategorySchema>;

