import { z } from "zod";

export const createMonthlyTargetSchema = z
  .object({
    group_id: z.string().uuid().optional(),
    user_id: z.string().uuid().optional(),
    brick_id: z.string().uuid().optional(),
    year: z
      .number()
      .int()
      .min(2000, "Year must be at least 2000")
      .max(2100, "Year must be at most 2100"),
    month: z
      .number()
      .int()
      .min(1, "Month must be between 1 and 12")
      .max(12, "Month must be between 1 and 12"),
    target_amount: z
      .number()
      .int()
      .min(0, "Target amount must be non-negative"),
  })
  .refine(
    (data) => data.group_id || data.user_id || data.brick_id,
    {
      message: "Either group_id, user_id, or brick_id must be provided",
      path: ["group_id"],
    }
  )
  .refine(
    (data) => {
      const count = [data.group_id, data.user_id, data.brick_id].filter(Boolean).length;
      return count === 1;
    },
    {
      message: "Only one of group_id, user_id, or brick_id can be provided",
      path: ["group_id"],
    }
  );

export const updateMonthlyTargetSchema = z.object({
  year: z
    .number()
    .int()
    .min(2000, "Year must be at least 2000")
    .max(2100, "Year must be at most 2100")
    .optional(),
  month: z
    .number()
    .int()
    .min(1, "Month must be between 1 and 12")
    .max(12, "Month must be between 1 and 12")
    .optional(),
  target_amount: z
    .number()
    .int()
    .min(0, "Target amount must be non-negative")
    .optional(),
});

export const bulkCreateMonthlyTargetSchema = z
  .object({
    group_ids: z.array(z.string().uuid()).optional(),
    user_ids: z.array(z.string().uuid()).optional(),
    brick_ids: z.array(z.string().uuid()).optional(),
    year: z
      .number()
      .int()
      .min(2000, "Year must be at least 2000")
      .max(2100, "Year must be at most 2100"),
    month: z
      .number()
      .int()
      .min(1, "Month must be between 1 and 12")
      .max(12, "Month must be between 1 and 12"),
    target_amount: z
      .number()
      .int()
      .min(0, "Target amount must be non-negative"),
  })
  .refine(
    (data) => 
      (data.group_ids && data.group_ids.length > 0) || 
      (data.user_ids && data.user_ids.length > 0) ||
      (data.brick_ids && data.brick_ids.length > 0),
    {
      message: "Either group_ids, user_ids, or brick_ids must be provided with at least one item",
      path: ["group_ids"],
    }
  )
  .refine(
    (data) => {
      const count = [
        data.group_ids && data.group_ids.length > 0,
        data.user_ids && data.user_ids.length > 0,
        data.brick_ids && data.brick_ids.length > 0,
      ].filter(Boolean).length;
      return count === 1;
    },
    {
      message: "Only one of group_ids, user_ids, or brick_ids can be provided",
      path: ["group_ids"],
    }
  );

export const createGroupTargetWithUserAssignmentSchema = z.object({
  group_id: z.string().uuid("Invalid group ID"),
  year: z
    .number()
    .int()
    .min(2000, "Year must be at least 2000")
    .max(2100, "Year must be at most 2100"),
  month: z
    .number()
    .int()
    .min(1, "Month must be between 1 and 12")
    .max(12, "Month must be between 1 and 12"),
  target_amount: z
    .number()
    .int()
    .min(0, "Target amount must be non-negative"),
});

export const bulkSetTargetSchema = z
  .object({
    group_id: z.string().uuid().optional(),
    user_id: z.string().uuid().optional(),
    brick_id: z.string().uuid().optional(),
    year: z
      .number()
      .int()
      .min(2000, "Year must be at least 2000")
      .max(2100, "Year must be at most 2100"),
    start_month: z
      .number()
      .int()
      .min(1, "Month must be between 1 and 12")
      .max(12, "Month must be between 1 and 12"),
    end_month: z
      .number()
      .int()
      .min(1, "Month must be between 1 and 12")
      .max(12, "Month must be between 1 and 12"),
    target_amount: z
      .number()
      .int()
      .min(0, "Target amount must be non-negative"),
  })
  .refine(
    (data) => data.group_id || data.user_id || data.brick_id,
    {
      message: "Either group_id, user_id, or brick_id must be provided",
      path: ["group_id"],
    }
  )
  .refine(
    (data) => {
      const count = [data.group_id, data.user_id, data.brick_id].filter(Boolean).length;
      return count === 1;
    },
    {
      message: "Only one of group_id, user_id, or brick_id can be provided",
      path: ["group_id"],
    }
  )
  .refine((data) => data.start_month <= data.end_month, {
    message: "Start month cannot be greater than end month",
    path: ["start_month"],
  });

export type CreateMonthlyTargetFormData = z.infer<
  typeof createMonthlyTargetSchema
>;
export type BulkCreateMonthlyTargetFormData = z.infer<
  typeof bulkCreateMonthlyTargetSchema
>;
export type UpdateMonthlyTargetFormData = z.infer<
  typeof updateMonthlyTargetSchema
>;
export type CreateGroupTargetWithUserAssignmentFormData = z.infer<
  typeof createGroupTargetWithUserAssignmentSchema
>;
export type BulkSetTargetFormData = z.infer<typeof bulkSetTargetSchema>;

