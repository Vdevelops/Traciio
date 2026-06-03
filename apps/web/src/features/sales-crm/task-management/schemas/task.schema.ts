import { z } from "zod";
import type { TaskPriority, TaskStatus, TaskType } from "../types";

export const taskTypeValues: TaskType[] = ["general", "call", "email", "meeting", "follow_up"];

export const taskStatusValues: TaskStatus[] = ["pending", "completed"];

export const taskPriorityValues: TaskPriority[] = ["low", "medium", "high", "urgent"];

export const createTaskSchema = z.object({
  title: z
    .string()
    .min(1, "Title is required")
    .min(3, "Title must be at least 3 characters")
    .max(255, "Title must be at most 255 characters"),
  description: z.string().optional(),
  type: z.enum(taskTypeValues).default("general"),
  priority: z.enum(taskPriorityValues).default("medium"),
  due_date: z
    .date()
    .nullable()
    .optional()
    .or(z.null())
    .refine(
      (value) => {
        if (!value) return true;
        return value instanceof Date && !Number.isNaN(value.getTime());
      },
      { message: "Invalid due date" }
    ),
  due_time: z
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
  assigned_to: z.string().uuid("Invalid user ID").optional().or(z.literal("")),
  lead_id: z.string().uuid("Invalid lead ID").optional().or(z.literal("")),
  account_id: z.string().uuid("Invalid account ID").optional().or(z.literal("")),
  contact_id: z.string().uuid("Invalid contact ID").optional().or(z.literal("")),
  deal_id: z.string().uuid("Invalid deal ID").optional().or(z.literal("")),
  sync_to_google_calendar: z.boolean().default(false),
});

export const updateTaskSchema = z.object({
  title: z
    .string()
    .min(3, "Title must be at least 3 characters")
    .max(255, "Title must be at most 255 characters")
    .optional(),
  description: z.string().optional(),
  type: z.enum(taskTypeValues).optional(),
  status: z.enum(taskStatusValues).optional(),
  priority: z.enum(taskPriorityValues).optional(),
  due_date: z
    .date()
    .nullable()
    .optional()
    .or(z.null())
    .refine(
      (value) => {
        if (!value) return true;
        return value instanceof Date && !Number.isNaN(value.getTime());
      },
      { message: "Invalid due date" }
    ),
  due_time: z
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
  assigned_to: z.string().uuid("Invalid user ID").optional().or(z.literal("")),
  lead_id: z.string().uuid("Invalid lead ID").optional().or(z.literal("")),
  account_id: z.string().uuid("Invalid account ID").optional().or(z.literal("")),
  contact_id: z.string().uuid("Invalid contact ID").optional().or(z.literal("")),
  deal_id: z.string().uuid("Invalid deal ID").optional().or(z.literal("")),
  sync_to_google_calendar: z.boolean().default(false),
});

export const assignTaskSchema = z.object({
  assigned_to: z
    .string()
    .min(1, "Assignee is required")
    .uuid("Invalid user ID"),
});

export type CreateTaskFormData = z.infer<typeof createTaskSchema>;
export type UpdateTaskFormData = z.infer<typeof updateTaskSchema>;
export type AssignTaskFormData = z.infer<typeof assignTaskSchema>;
