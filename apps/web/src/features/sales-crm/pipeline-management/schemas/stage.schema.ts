import { z } from "zod";

/**
 * Pipeline Stage Form Schema
 * Used for creating and updating stages
 */
export const stageSchema = z.object({
  name: z
    .string()
    .min(1, "Name is required")
    .max(100, "Name must be less than 100 characters"),
  code: z
    .string()
    .min(1, "Code is required")
    .max(50, "Code must be less than 50 characters")
    .regex(/^[a-z0-9_]+$/, "Code must contain only lowercase letters, numbers, and underscores"),
  description: z.string().optional(),
  order: z
    .number()
    .int("Order must be an integer")
    .min(0, "Order must be greater than or equal to 0"),
  color: z
    .string()
    .regex(/^#[0-9A-F]{6}$/i, "Color must be a valid hex color")
    .optional(),
  is_active: z.boolean().default(true),
  is_won: z.boolean().default(false),
  is_lost: z.boolean().default(false),
  probability: z
    .number()
    .min(0, "Probability must be between 0 and 100")
    .max(100, "Probability must be between 0 and 100")
    .optional(),
  requirements: z.string().optional(),
});

export type StageFormData = z.infer<typeof stageSchema>;

/**
 * Stage Update Schema (all fields optional)
 */
export const stageUpdateSchema = stageSchema.partial();

export type StageUpdateData = z.infer<typeof stageUpdateSchema>;

/**
 * Stage Reorder Schema
 */
export const stageReorderSchema = z.object({
  stage_id: z.string().uuid(),
  new_order: z.number().int().min(0),
});

export type StageReorderData = z.infer<typeof stageReorderSchema>;
