"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { taskService, reminderService } from "../services/taskService";
import type {
  CreateTaskFormData,
  UpdateTaskFormData,
  AssignTaskFormData,
} from "../schemas/task.schema";
import type { CreateReminderFormData, UpdateReminderFormData } from "../schemas/reminder.schema";
import type { ReminderListParams, Task, TaskListParams, TaskResponse } from "../types";

function syncTaskIntoCaches(queryClient: ReturnType<typeof useQueryClient>, task: Task) {
  queryClient.setQueriesData(
    { queryKey: ["tasks"] },
    (existing: unknown) => {
      if (!existing || typeof existing !== "object") {
        return existing;
      }

      if ("data" in existing && Array.isArray((existing as { data?: unknown }).data)) {
        const typedExisting = existing as { data: Task[] };
        const alreadyExists = typedExisting.data.some((item) => item.id === task.id);
        return {
          ...typedExisting,
          data: alreadyExists
            ? typedExisting.data.map((item) => (item.id === task.id ? { ...item, ...task } : item))
            : [task, ...typedExisting.data],
        };
      }

      if ("success" in existing && "data" in existing) {
        const typedExisting = existing as TaskResponse;
        if (typedExisting.data?.id === task.id) {
          return {
            ...typedExisting,
            data: { ...typedExisting.data, ...task },
          };
        }
      }

      return existing;
    },
  );
}

export function useTasks(params?: TaskListParams) {
  return useQuery({
    queryKey: ["tasks", params],
    queryFn: () => taskService.list(params),
    retry: (failureCount, error) => {
      if (error && typeof error === "object" && "response" in error) {
        const axiosError = error as { response?: { status?: number } };
        if (axiosError.response?.status === 404) {
          return false;
        }
      }
      return failureCount < 1;
    },
  });
}

export function useTask(id: string) {
  return useQuery({
    queryKey: ["tasks", id],
    queryFn: () => taskService.getById(id),
    enabled: !!id,
  });
}

export function useCreateTask() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateTaskFormData) => taskService.create(data),
    onSuccess: (response) => {
      syncTaskIntoCaches(queryClient, response.data);
      queryClient.invalidateQueries({ queryKey: ["tasks"] });
    },
  });
}

export function useUpdateTask() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateTaskFormData }) =>
      taskService.update(id, data),
    onSuccess: (response, variables) => {
      syncTaskIntoCaches(queryClient, response.data);
      queryClient.invalidateQueries({ queryKey: ["tasks"] });
      queryClient.invalidateQueries({ queryKey: ["tasks", variables.id] });
    },
  });
}

export function useDeleteTask() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => taskService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["tasks"] });
    },
  });
}

export function useAssignTask() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: AssignTaskFormData }) =>
      taskService.assign(id, data),
    onSuccess: (response, variables) => {
      syncTaskIntoCaches(queryClient, response.data);
      queryClient.invalidateQueries({ queryKey: ["tasks"] });
      queryClient.invalidateQueries({ queryKey: ["tasks", variables.id] });
    },
  });
}

export function useCompleteTask() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => taskService.complete(id),
    onSuccess: (response, id) => {
      syncTaskIntoCaches(queryClient, response.data);
      queryClient.invalidateQueries({ queryKey: ["tasks"] });
      queryClient.invalidateQueries({ queryKey: ["tasks", id] });
    },
  });
}

export function useReminders(params?: ReminderListParams) {
  return useQuery({
    queryKey: ["reminders", params],
    queryFn: () => reminderService.list(params),
    retry: (failureCount, error) => {
      if (error && typeof error === "object" && "response" in error) {
        const axiosError = error as { response?: { status?: number } };
        if (axiosError.response?.status === 404) {
          return false;
        }
      }
      return failureCount < 1;
    },
  });
}

export function useReminder(id: string) {
  return useQuery({
    queryKey: ["reminders", id],
    queryFn: () => reminderService.getById(id),
    enabled: !!id,
  });
}

export function useCreateReminder() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateReminderFormData) => reminderService.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["reminders"] });
      queryClient.invalidateQueries({ queryKey: ["tasks"] });
    },
  });
}

export function useUpdateReminder() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateReminderFormData }) =>
      reminderService.update(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["reminders"] });
      queryClient.invalidateQueries({ queryKey: ["reminders", variables.id] });
      queryClient.invalidateQueries({ queryKey: ["tasks"] });
    },
  });
}

export function useDeleteReminder() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => reminderService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["reminders"] });
      queryClient.invalidateQueries({ queryKey: ["tasks"] });
    },
  });
}
