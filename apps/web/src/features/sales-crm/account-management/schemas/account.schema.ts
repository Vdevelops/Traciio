import { z } from "zod";

const latitudeSchema = z
  .number()
  .min(-90, "Latitude must be between -90 and 90")
  .max(90, "Latitude must be between -90 and 90")
  .optional();

const longitudeSchema = z
  .number()
  .min(-180, "Longitude must be between -180 and 180")
  .max(180, "Longitude must be between -180 and 180")
  .optional();

export const createAccountSchema = z.object({
  name: z.string().min(3, "Name must be at least 3 characters"),
  category_id: z
    .string()
    .min(1, "Category is required")
    .uuid("Invalid category ID"),
  address: z.string().optional(),
  city: z.string().optional(),
  province: z.string().optional(),
  phone: z.string().optional(),
  email: z.string().email("Invalid email format").optional().or(z.literal("")),
  latitude: latitudeSchema,
  longitude: longitudeSchema,
  status: z.enum(["active", "inactive"]).optional().default("active"),
  assigned_to: z.string().uuid("Invalid user ID").optional().or(z.literal("")),
});

export const updateAccountSchema = z.object({
  name: z.string().min(3, "Name must be at least 3 characters").optional(),
  category_id: z.string().uuid("Invalid category ID").optional(),
  address: z.string().optional(),
  city: z.string().optional(),
  province: z.string().optional(),
  phone: z.string().optional(),
  email: z.string().email("Invalid email format").optional().or(z.literal("")),
  latitude: latitudeSchema,
  longitude: longitudeSchema,
  status: z.enum(["active", "inactive"]).optional(),
  assigned_to: z.string().uuid("Invalid user ID").optional().or(z.literal("")),
});

export type CreateAccountFormData = z.infer<typeof createAccountSchema>;
export type UpdateAccountFormData = z.infer<typeof updateAccountSchema>;

