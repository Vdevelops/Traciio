import { z } from "zod";

export const createContactSchema = z.object({
  account_id: z.string().uuid("Invalid account ID"),
  name: z.string().min(3, "Name must be at least 3 characters"),
  role_id: z
    .string()
    .min(1, "Role is required")
    .uuid("Invalid contact role ID"),
  phone: z.string().optional(),
  email: z.string().email("Invalid email format").optional().or(z.literal("")),
  position: z.string().optional(),
  notes: z.string().optional(),
});

export const updateContactSchema = z.object({
  account_id: z.string().uuid("Invalid account ID").optional(),
  name: z.string().min(3, "Name must be at least 3 characters").optional(),
  role_id: z.string().uuid("Invalid contact role ID").optional(),
  phone: z.string().optional(),
  email: z.string().email("Invalid email format").optional().or(z.literal("")),
  position: z.string().optional(),
  notes: z.string().optional(),
});

export type CreateContactFormData = z.infer<typeof createContactSchema>;
export type UpdateContactFormData = z.infer<typeof updateContactSchema>;

