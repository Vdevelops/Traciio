// Schedule Management Types - User-specific schedules connected to tasks

import type { Task } from "../../task-management/types";

export type ScheduleStatus = "pending" | "confirmed" | "completed" | "cancelled";

export type GoogleCalendarSyncStatus = "not_synced" | "synced" | "sync_failed";

export interface ScheduleTaskRef {
  id: string;
  title: string;
  due_date: string | null;
  assigned_to: string;
  assigned_user?: {
    id: string;
    name: string;
    email: string;
  };
}

export interface Schedule {
  id: string;
  task_id: string;
  task?: ScheduleTaskRef;
  user_id: string; // User who owns this schedule (from task.assigned_to)
  title: string;
  description: string | null;
  scheduled_at: string; // When the schedule/reminder should occur
  status: ScheduleStatus;
  google_calendar_event_id: string | null; // Google Calendar event ID if synced
  google_calendar_sync_status: GoogleCalendarSyncStatus;
  google_calendar_synced_at: string | null;
  google_calendar_event_link: string | null; // Direct URL to view in Google Calendar
  reminder_minutes_before: number | null; // Minutes before task due_date to remind
  created_by: string;
  created_at: string;
  updated_at: string;
}

export interface ListSchedulesResponse {
  success: boolean;
  data: Schedule[];
  meta: {
    pagination: {
      page: number;
      per_page: number;
      total: number;
      total_pages: number;
      has_next: boolean;
      has_prev: boolean;
    };
    filters?: Record<string, unknown>;
  };
  timestamp: string;
  request_id: string;
}

export interface ScheduleResponse {
  success: boolean;
  data: Schedule;
  timestamp: string;
  request_id: string;
}

export interface ScheduleListParams {
  page?: number;
  per_page?: number;
  search?: string;
  status?: ScheduleStatus;
  task_id?: string;
  user_id?: string;
  scheduled_at_from?: string;
  scheduled_at_to?: string;
  google_calendar_sync_status?: GoogleCalendarSyncStatus;
}

export interface GoogleCalendarEvent {
  id: string;
  summary: string;
  description: string | null;
  start: {
    dateTime: string;
    timeZone: string;
  };
  end: {
    dateTime: string;
    timeZone: string;
  };
  reminders: {
    useDefault: boolean;
    overrides: Array<{
      method: "email" | "popup";
      minutes: number;
    }>;
  };
}

export interface GoogleCalendarSyncResponse {
  success: boolean;
  data: {
    schedule_id: string;
    google_calendar_event_id: string;
    event_url: string;
    sync_status: GoogleCalendarSyncStatus;
  };
  timestamp: string;
  request_id: string;
}

