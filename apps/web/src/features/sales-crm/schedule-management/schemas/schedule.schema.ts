import { z } from "zod";
import type { ScheduleStatus } from "../types";

export const scheduleStatusValues: ScheduleStatus[] = ["pending", "confirmed", "completed", "cancelled"];

export const createScheduleSchema = z.object({
  task_id: z
    .string()
    .uuid("Invalid task ID")
    .optional(),
  title: z
    .string()
    .min(1, "Title is required")
    .min(3, "Title must be at least 3 characters")
    .max(255, "Title must be at most 255 characters"),
  description: z.string().optional(),
  scheduled_at: z.string().optional(), // Will be set from scheduled_date + scheduled_time
  scheduled_date: z
    .date()
    .nullable()
    .optional()
    .or(z.null())
    .refine(
      (value) => {
        if (!value) return true;
        return value instanceof Date && !Number.isNaN(value.getTime());
      },
      { message: "Invalid scheduled date" }
    ),
  scheduled_time: z
    .string()
    .nullable()
    .optional()
    .or(z.null())
    .refine(
      (value) => {
        if (!value) return true;
        const timeRegex = /^([0-1]?\d|2[0-3]):[0-5]\d$/;
        return timeRegex.test(value);
      },
      { message: "Invalid time format (expected HH:mm)" }
    ),
  reminder_minutes_before: z
    .number()
    .int("Must be an integer")
    .min(0, "Must be 0 or positive")
    .max(10080, "Must be less than 7 days (10080 minutes)")
    .optional()
    .nullable(),
});

export const updateScheduleSchema = z.object({
  title: z
    .string()
    .min(3, "Title must be at least 3 characters")
    .max(255, "Title must be at most 255 characters")
    .optional(),
  description: z.string().optional(),
  scheduled_at: z.string().optional(), // Will be set from scheduled_date + scheduled_time
  scheduled_date: z
    .date()
    .nullable()
    .optional()
    .or(z.null())
    .refine(
      (value) => {
        if (!value) return true;
        return value instanceof Date && !Number.isNaN(value.getTime());
      },
      { message: "Invalid scheduled date" }
    ),
  scheduled_time: z
    .string()
    .nullable()
    .optional()
    .or(z.null())
    .refine(
      (value) => {
        if (!value) return true;
        const timeRegex = /^([0-1]?\d|2[0-3]):[0-5]\d$/;
        return timeRegex.test(value);
      },
      { message: "Invalid time format (expected HH:mm)" }
    ),
  status: z.enum(scheduleStatusValues).optional(),
  reminder_minutes_before: z
    .number()
    .int("Must be an integer")
    .min(0, "Must be 0 or positive")
    .max(10080, "Must be less than 7 days (10080 minutes)")
    .optional()
    .nullable(),
});

export type CreateScheduleFormData = z.infer<typeof createScheduleSchema>;
export type UpdateScheduleFormData = z.infer<typeof updateScheduleSchema>;
