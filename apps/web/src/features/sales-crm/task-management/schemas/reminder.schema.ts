import { z } from "zod";
import type { ReminderType } from "../types";

export const reminderTypeValues: ReminderType[] = ["in_app", "email", "sms"];

export const createReminderSchema = z.object({
  task_id: z
    .string()
    .min(1, "Task is required")
    .uuid("Invalid task ID"),
  remind_at: z
    .string()
    .refine(
      (value) => {
        const date = new Date(value);
        return !Number.isNaN(date.getTime());
      },
      { message: "Invalid reminder time" }
    ),
  reminder_type: z.enum(reminderTypeValues).default("in_app"),
  message: z.string().optional(),
});

export const updateReminderSchema = z.object({
  remind_at: z
    .string()
    .optional()
    .refine(
      (value) => {
        if (!value) return true;
        const date = new Date(value);
        return !Number.isNaN(date.getTime());
      },
      { message: "Invalid reminder time" }
    ),
  reminder_type: z.enum(reminderTypeValues).optional(),
  message: z.string().optional(),
});

export type CreateReminderFormData = z.infer<typeof createReminderSchema>;
export type UpdateReminderFormData = z.infer<typeof updateReminderSchema>;


