import apiClient from "@/lib/api-client";
import type {
  ListSchedulesResponse,
  ScheduleResponse,
  ScheduleListParams,
} from "../types";
import type { CreateScheduleFormData, UpdateScheduleFormData } from "../schemas/schedule.schema";

export const scheduleService = {
  async list(params?: ScheduleListParams): Promise<ListSchedulesResponse> {
    const response = await apiClient.get<ListSchedulesResponse>("/schedules", { params });
    return response.data;
  },

  async getById(id: string): Promise<ScheduleResponse> {
    const response = await apiClient.get<ScheduleResponse>(`/schedules/${id}`);
    return response.data;
  },

  async create(data: CreateScheduleFormData): Promise<ScheduleResponse> {
    const payload: Record<string, unknown> = {
      title: data.title,
      scheduled_at: data.scheduled_at ? new Date(data.scheduled_at).toISOString() : undefined,
    };

    // Only include task_id if provided
    if (data.task_id) {
      payload.task_id = data.task_id;
    }

    if (data.description) {
      payload.description = data.description;
    }

    if (data.reminder_minutes_before !== undefined && data.reminder_minutes_before !== null) {
      payload.reminder_minutes_before = data.reminder_minutes_before;
    }

    const response = await apiClient.post<ScheduleResponse>("/schedules", payload);
    return response.data;
  },

  async update(id: string, data: UpdateScheduleFormData): Promise<ScheduleResponse> {
    const payload: Record<string, unknown> = {};

    if (data.title !== undefined) {
      payload.title = data.title;
    }

    if (data.description !== undefined) {
      payload.description = data.description;
    }

    if (data.scheduled_at !== undefined) {
      payload.scheduled_at = new Date(data.scheduled_at).toISOString();
    }

    if (data.status !== undefined) {
      payload.status = data.status;
    }

    if (data.reminder_minutes_before !== undefined) {
      payload.reminder_minutes_before = data.reminder_minutes_before;
    }

    const response = await apiClient.put<ScheduleResponse>(`/schedules/${id}`, payload);
    return response.data;
  },

  async delete(id: string): Promise<void> {
    await apiClient.delete(`/schedules/${id}`);
  },
};
