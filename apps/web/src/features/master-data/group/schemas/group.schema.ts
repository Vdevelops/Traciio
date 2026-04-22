import { z } from "zod";

export const createGroupSchema = z.object({
  name: z.string().min(3, "Name must be at least 3 characters"),
  code: z.string().min(3, "Code must be at least 3 characters"),
  description: z.string().optional(),
  status: z.enum(["active", "inactive"]).optional().default("active"),
});

export const updateGroupSchema = z.object({
  name: z.string().min(3, "Name must be at least 3 characters").optional(),
  code: z.string().min(3, "Code must be at least 3 characters").optional(),
  description: z.string().optional(),
  status: z.enum(["active", "inactive"]).optional(),
});

export type CreateGroupFormData = z.infer<typeof createGroupSchema>;
export type UpdateGroupFormData = z.infer<typeof updateGroupSchema>;

