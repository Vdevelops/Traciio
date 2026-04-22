import { z } from "zod";

export const createRoleSchema = z.object({
  name: z.string().min(3, "Name must be at least 3 characters"),
  code: z.string().min(3, "Code must be at least 3 characters"),
  description: z.string().optional(),
  status: z.enum(["active", "inactive"]).optional().default("active"),
  mobile_access: z.boolean().optional().default(false),
});

export const updateRoleSchema = z.object({
  name: z.string().min(3, "Name must be at least 3 characters").optional(),
  code: z.string().min(3, "Code must be at least 3 characters").optional(),
  description: z.string().optional(),
  status: z.enum(["active", "inactive"]).optional(),
  mobile_access: z.boolean().optional(),
});

export const assignPermissionsSchema = z.object({
  permission_ids: z.array(z.string().uuid("Invalid permission ID")).min(1, "At least one permission is required"),
});

export type CreateRoleFormData = z.infer<typeof createRoleSchema>;
export type UpdateRoleFormData = z.infer<typeof updateRoleSchema>;
export type AssignPermissionsFormData = z.infer<typeof assignPermissionsSchema>;

