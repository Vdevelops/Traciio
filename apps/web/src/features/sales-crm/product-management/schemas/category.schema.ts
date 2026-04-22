import { z } from "zod";

export const createCategorySchema = z.object({
  name: z.string().min(3, "Name must be at least 3 characters").max(150, "Name must be at most 150 characters"),
  slug: z.string().max(150, "Slug must be at most 150 characters").optional(),
  description: z.string().optional(),
  status: z.enum(["active", "inactive"]).optional().default("active"),
});

export const updateCategorySchema = z.object({
  name: z.string().min(3, "Name must be at least 3 characters").max(150, "Name must be at most 150 characters").optional(),
  slug: z.string().max(150, "Slug must be at most 150 characters").optional(),
  description: z.string().optional(),
  status: z.enum(["active", "inactive"]).optional(),
});

export type CreateCategoryFormData = z.infer<typeof createCategorySchema>;
export type UpdateCategoryFormData = z.infer<typeof updateCategorySchema>;

