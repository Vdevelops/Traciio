import { z } from "zod";

// Safe URL validation that allows empty strings
const safeUrlSchema = z.string().refine(
  (val) => {
    if (!val || val === "") return true;
    try {
      const url = new URL(val);
      return url.protocol === "http:" || url.protocol === "https:";
    } catch {
      return false;
    }
  },
  { message: "Invalid image URL" }
);

export const createProductSchema = z.object({
  name: z.string().min(3, "Name must be at least 3 characters").max(200, "Name must be at most 200 characters"),
  sku: z.string().min(1, "SKU is required").max(100, "SKU must be at most 100 characters"),
  barcode: z.string().max(100, "Barcode must be at most 100 characters").optional(),
  price: z.number().min(0, "Price must be at least 0"),
  cost: z.number().min(0, "Cost must be at least 0").optional(),
  category_id: z.string().uuid("Invalid category ID"),
  description: z.string().optional(),
  status: z.enum(["active", "inactive"]).optional().default("active"),
  image_url: safeUrlSchema.optional(),
});

export const updateProductSchema = z.object({
  name: z.string().min(3, "Name must be at least 3 characters").max(200, "Name must be at most 200 characters").optional(),
  sku: z.string().min(1, "SKU is required").max(100, "SKU must be at most 100 characters").optional(),
  barcode: z.string().max(100, "Barcode must be at most 100 characters").optional(),
  price: z.number().min(0, "Price must be at least 0").optional(),
  cost: z.number().min(0, "Cost must be at least 0").optional(),
  category_id: z.string().uuid("Invalid category ID").optional(),
  description: z.string().optional(),
  status: z.enum(["active", "inactive"]).optional(),
  image_url: safeUrlSchema.optional(),
});

export type CreateProductFormData = z.infer<typeof createProductSchema>;
export type UpdateProductFormData = z.infer<typeof updateProductSchema>;

