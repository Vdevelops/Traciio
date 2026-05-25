// Task & Reminder Types for Sales CRM - Task Management

export type TaskType = "general" | "call" | "email" | "meeting" | "follow_up";

export type TaskStatus = "pending" | "in_progress" | "completed" | "cancelled";

export type TaskPriority = "low" | "medium" | "high" | "urgent";

export interface TaskAssignedUserRef {
  id: string;
  name: string;
  email: string;
  avatar_url?: string;
}

export interface TaskAccountRef {
  id: string;
  name: string;
}

export interface TaskContactRef {
  id: string;
  name: string;
  email: string;
  phone: string;
}

export interface TaskDealRef {
  id: string;
  title: string;
}

export interface Task {
  id: string;
  title: string;
  description: string;
  type: TaskType;
  status: TaskStatus;
  priority: TaskPriority;
  due_date: string | null;
  completed_at: string | null;
  assigned_to: string;
  assigned_user?: TaskAssignedUserRef;
  account_id?: string;
  account?: TaskAccountRef;
  contact_id?: string;
  contact?: TaskContactRef;
  deal_id?: string;
  deal?: TaskDealRef;

  // Lead integration
  lead_id?: string;
  lead?: LeadRef;

  // Source tracking
  task_source?: "manual" | "lead_tab" | "pipeline_tab" | "auto_generated";

  // Schedule unification
  scheduled_start_time?: string;
  scheduled_end_time?: string;
  scheduled_location?: string;
  is_schedule_task?: boolean;

  // Quick action
  quick_action_type?: "add_lead" | "convert_lead" | "add_visit" | "add_deal";
  quick_action_payload?: Record<string, unknown>;

  google_calendar_event_id?: string | null;
  google_calendar_sync_status?: string | null;
  google_calendar_synced_at?: string | null;
  created_by: string;
  created_at: string;
  updated_at: string;
}

export interface LeadRef {
  id: string;
  first_name: string;
  last_name: string;
  company_name: string;
  email: string;
  full_name: string;
}

// Request for creating lead from task
export interface AddLeadFromTaskRequest {
  first_name: string;
  last_name?: string;
  company_name?: string;
  email: string;
  phone?: string;
  lead_source: string;
  notes?: string;
  budget_amount?: number;
  expected_close_date?: string;
}

export interface AddLeadFromTaskResponse {
  task_id: string;
  lead_id: string;
  lead?: LeadRef;
  task_updated: boolean;
  message: string;
}

export interface ListTasksResponse {
  success: boolean;
  data: Task[];
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

export interface TaskResponse {
  success: boolean;
  data: Task;
  timestamp: string;
  request_id: string;
}

export type ReminderType = "in_app" | "email" | "sms";

export interface ReminderTaskRef {
  id: string;
  title: string;
}

export interface Reminder {
  id: string;
  task_id: string;
  task?: ReminderTaskRef;
  remind_at: string;
  reminder_type: ReminderType;
  is_sent: boolean;
  sent_at: string | null;
  message: string;
  created_by: string;
  created_at: string;
  updated_at: string;
}

export interface ListRemindersResponse {
  success: boolean;
  data: Reminder[];
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

export interface ReminderResponse {
  success: boolean;
  data: Reminder;
  timestamp: string;
  request_id: string;
}

export interface TaskListParams {
  page?: number;
  per_page?: number;
  search?: string;
  status?: TaskStatus;
  priority?: TaskPriority;
  type?: TaskType;
  assigned_to?: string;
  account_id?: string;
  contact_id?: string;
  deal_id?: string;
  lead_id?: string;
  due_date_from?: string;
  due_date_to?: string;
}

export interface ReminderListParams {
  page?: number;
  per_page?: number;
  task_id?: string;
  reminder_type?: ReminderType;
  is_sent?: boolean;
  remind_at_from?: string;
  remind_at_to?: string;
}


