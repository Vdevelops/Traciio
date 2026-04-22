import { z } from "zod";

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

export type CreateGroupTargetWithUserAssignmentFormData = z.infer<
  typeof createGroupTargetWithUserAssignmentSchema
>;
