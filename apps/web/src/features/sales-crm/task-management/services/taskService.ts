import apiClient from "@/lib/api-client";
import type {
  ListTasksResponse,
  TaskResponse,
  TaskListParams,
  ListRemindersResponse,
  ReminderResponse,
  ReminderListParams,
} from "../types";
import type {
  CreateTaskFormData,
  UpdateTaskFormData,
  AssignTaskFormData,
} from "../schemas/task.schema";
import type { CreateReminderFormData, UpdateReminderFormData } from "../schemas/reminder.schema";

export const taskService = {
  async list(params?: TaskListParams): Promise<ListTasksResponse> {
    const response = await apiClient.get<ListTasksResponse>("/tasks", { params });
    return response.data;
  },

  async getById(id: string): Promise<TaskResponse> {
    const response = await apiClient.get<TaskResponse>(`/tasks/${id}`);
    return response.data;
  },

  async create(data: CreateTaskFormData | Record<string, unknown>): Promise<TaskResponse> {
    const dataRecord = data as Record<string, unknown>;

    const payload: Record<string, unknown> = {
      title: dataRecord.title as string,
      description: dataRecord.description as string | undefined,
      type: dataRecord.type as string | undefined,
      priority: dataRecord.priority as string | undefined,
    };

    const dueDateValue = dataRecord.due_date;
    if (dueDateValue) {
      if (typeof dueDateValue === "string") {
        payload.due_date = dueDateValue;
      } else if (dueDateValue instanceof Date) {
        if (!Number.isNaN(dueDateValue.getTime())) {
          payload.due_date = dueDateValue.toISOString();
        }
      }
    }

    const assignedToValue = dataRecord.assigned_to;
    if (assignedToValue && typeof assignedToValue === "string" && assignedToValue.trim() !== "") {
      payload.assigned_to = assignedToValue;
    }

    if (dataRecord.account_id && typeof dataRecord.account_id === "string" && dataRecord.account_id.trim() !== "") {
      payload.account_id = dataRecord.account_id;
    }

    if (dataRecord.contact_id && typeof dataRecord.contact_id === "string" && dataRecord.contact_id.trim() !== "") {
      payload.contact_id = dataRecord.contact_id;
    }

    if (dataRecord.deal_id && typeof dataRecord.deal_id === "string" && dataRecord.deal_id.trim() !== "") {
      payload.deal_id = dataRecord.deal_id;
    }

    const response = await apiClient.post<TaskResponse>("/tasks", payload);
    return response.data;
  },

  async update(id: string, data: UpdateTaskFormData): Promise<TaskResponse> {
    const payload: Record<string, unknown> = {};

    if (data.title !== undefined) {
      payload.title = data.title;
    }

    if (data.description !== undefined) {
      payload.description = data.description;
    }

    if (data.type !== undefined) {
      payload.type = data.type;
    }

    if (data.status !== undefined) {
      payload.status = data.status;
    }

    if (data.priority !== undefined) {
      payload.priority = data.priority;
    }

    if (data.due_date !== undefined) {
      if (data.due_date) {
        if (typeof data.due_date === "string") {
          // Already ISO string from form submit handler
          payload.due_date = data.due_date;
        } else if (data.due_date instanceof Date) {
          // Date object - convert to ISO string
          if (!Number.isNaN(data.due_date.getTime())) {
            payload.due_date = data.due_date.toISOString();
          }
        }
      } else {
        payload.due_date = null;
      }
    }

    if (data.assigned_to !== undefined) {
      payload.assigned_to = data.assigned_to || null;
    }

    if (data.account_id !== undefined) {
      payload.account_id = data.account_id || null;
    }

    if (data.contact_id !== undefined) {
      payload.contact_id = data.contact_id || null;
    }

    if (data.deal_id !== undefined) {
      payload.deal_id = data.deal_id || null;
    }

    const response = await apiClient.put<TaskResponse>(`/tasks/${id}`, payload);
    return response.data;
  },

  async delete(id: string): Promise<void> {
    await apiClient.delete(`/tasks/${id}`);
  },

  async assign(id: string, data: AssignTaskFormData): Promise<TaskResponse> {
    const response = await apiClient.post<TaskResponse>(`/tasks/${id}/assign`, data);
    return response.data;
  },

  async complete(id: string): Promise<TaskResponse> {
    const response = await apiClient.post<TaskResponse>(`/tasks/${id}/complete`);
    return response.data;
  },

  // Add lead from task (quick action)
  async addLeadFromTask(
    taskId: string,
    data: import("../types").AddLeadFromTaskRequest,
  ): Promise<import("../types").AddLeadFromTaskResponse> {
    const response = await apiClient.post<{ data: import("../types").AddLeadFromTaskResponse }>(
      `/tasks/${taskId}/add-lead`,
      data,
    );
    return response.data.data;
  },

  // Get scheduled tasks (unified schedule view)
  async getScheduledTasks(
    startDate?: string,
    endDate?: string,
    page = 1,
    perPage = 20,
  ): Promise<ListTasksResponse> {
    const params: Record<string, unknown> = { page, per_page: perPage };
    if (startDate) params.start_date = startDate;
    if (endDate) params.end_date = endDate;

    const response = await apiClient.get<ListTasksResponse>(
      "/tasks/schedule",
      { params },
    );
    return response.data;
  },

  // Get tasks by lead
  async getTasksByLead(
    leadId: string,
    page = 1,
    perPage = 20,
  ): Promise<ListTasksResponse> {
    const response = await apiClient.get<ListTasksResponse>(
      `/leads/${leadId}/tasks`,
      { params: { page, per_page: perPage } },
    );
    return response.data;
  },

  // Get tasks by deal
  async getTasksByDeal(
    dealId: string,
    page = 1,
    perPage = 20,
  ): Promise<ListTasksResponse> {
    const response = await apiClient.get<ListTasksResponse>(
      `/pipeline/deals/${dealId}/tasks`,
      { params: { page, per_page: perPage } },
    );
    return response.data;
  },
};

export const reminderService = {
  async list(params?: ReminderListParams): Promise<ListRemindersResponse> {
    const response = await apiClient.get<ListRemindersResponse>("/tasks/reminders", { params });
    return response.data;
  },

  async getById(id: string): Promise<ReminderResponse> {
    const response = await apiClient.get<ReminderResponse>(`/tasks/reminders/${id}`);
    return response.data;
  },

  async create(data: CreateReminderFormData): Promise<ReminderResponse> {
    const payload: Record<string, unknown> = {
      task_id: data.task_id,
      reminder_type: data.reminder_type,
      message: data.message,
    };

    const remindAt = new Date(data.remind_at);
    if (!Number.isNaN(remindAt.getTime())) {
      payload.remind_at = remindAt.toISOString();
    }

    const response = await apiClient.post<ReminderResponse>("/tasks/reminders", payload);
    return response.data;
  },

  async update(id: string, data: UpdateReminderFormData): Promise<ReminderResponse> {
    const payload: Record<string, unknown> = {};

    if (data.remind_at !== undefined) {
      if (data.remind_at) {
        const remindAt = new Date(data.remind_at);
        if (!Number.isNaN(remindAt.getTime())) {
          payload.remind_at = remindAt.toISOString();
        }
      } else {
        payload.remind_at = null;
      }
    }

    if (data.reminder_type !== undefined) {
      payload.reminder_type = data.reminder_type;
    }

    if (data.message !== undefined) {
      payload.message = data.message;
    }

    const response = await apiClient.put<ReminderResponse>(`/tasks/reminders/${id}`, payload);
    return response.data;
  },

  async delete(id: string): Promise<void> {
    await apiClient.delete(`/tasks/reminders/${id}`);
  },
};


