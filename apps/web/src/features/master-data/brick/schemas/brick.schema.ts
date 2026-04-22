import { z } from "zod";

export const createBrickSchema = z.object({
  name: z.string().min(3, "Name must be at least 3 characters"),
  code: z
    .string()
    .regex(/^[a-zA-Z0-9_-]*$/, "Code can only contain alphanumeric characters, dashes, and underscores")
    .optional()
    .or(z.literal("")),
  description: z.string().optional(),
  province: z.string().min(1, "Province is required"),
  regency: z.string().min(1, "Regency is required"),
  manager_id: z.string().uuid().optional(),
  status: z.enum(["active", "inactive"]).optional().default("active"),
});

export const updateBrickSchema = z.object({
  name: z.string().min(3, "Name must be at least 3 characters").optional(),
  code: z
    .string()
    .min(3, "Code must be at least 3 characters")
    .regex(/^[a-zA-Z0-9_-]+$/, "Code can only contain alphanumeric characters, dashes, and underscores")
    .optional(),
  description: z.string().optional(),
  province: z.string().min(1, "Province is required").optional(),
  regency: z.string().min(1, "Regency is required").optional(),
  manager_id: z.string().uuid().optional(),
  status: z.enum(["active", "inactive"]).optional(),
});

export type CreateBrickFormData = z.infer<typeof createBrickSchema>;
export type UpdateBrickFormData = z.infer<typeof updateBrickSchema>;

export const createBrickTargetDistributionSchema = z.object({
  distributions: z.array(
    z.object({
      sales_user_id: z.string().uuid("Sales user ID must be a valid UUID"),
      distributed_amount: z.number().int().min(0, "Distributed amount must be non-negative"),
    })
  ).min(1, "At least one distribution is required"),
});

export const updateBrickTargetDistributionSchema = z.object({
  distributed_amount: z.number().int().min(0, "Distributed amount must be non-negative"),
});

export type CreateBrickTargetDistributionFormData = z.infer<typeof createBrickTargetDistributionSchema>;
export type UpdateBrickTargetDistributionFormData = z.infer<typeof updateBrickTargetDistributionSchema>;

