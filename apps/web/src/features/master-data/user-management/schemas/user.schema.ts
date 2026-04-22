import { z } from "zod";

export const createUserSchema = z.object({
  email: z.string().email("Invalid email format"),
  password: z.string().min(6, "Password must be at least 6 characters"),
  name: z.string().min(3, "Name must be at least 3 characters"),
  role_id: z.string().uuid("Invalid role ID"),
  group_id: z.string().uuid("Invalid group ID").optional().nullable(),
  status: z.enum(["active", "inactive"]).optional().default("active"),
});

export const updateUserSchema = z.object({
  email: z.string().email("Invalid email format").optional(),
  name: z.string().min(3, "Name must be at least 3 characters").optional(),
  role_id: z
    .string()
    .uuid("Invalid role ID")
    .optional()
    .or(z.literal("").transform(() => undefined)),
  brick_id: z
    .string()
    .uuid("Invalid brick ID")
    .optional()
    .nullable()
    .or(z.literal("").transform(() => null)),
  status: z.enum(["active", "inactive"]).optional(),
});

export type CreateUserFormData = z.infer<typeof createUserSchema>;
export type UpdateUserFormData = z.infer<typeof updateUserSchema>;

