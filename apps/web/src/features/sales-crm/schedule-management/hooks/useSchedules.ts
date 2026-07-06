"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { scheduleService } from "../services/scheduleService";
import type {
  CreateScheduleFormData,
  UpdateScheduleFormData,
} from "../schemas/schedule.schema";
import type { ScheduleListParams } from "../types";

// Helper function to extract error message from API response
function extractErrorMessage(error: unknown, defaultMessage: string): string {
  if (!error || typeof error !== "object") {
    return defaultMessage;
  }

  // Check for axios error with response
  if ("response" in error) {
    const axiosError = error as {
      response?: { data?: { error?: { message?: string }; message?: string } };
      message?: string;
    };

    const responseData = axiosError.response?.data;
    if (responseData && typeof responseData === "object") {
      // API response structure: { error: { message: string } }
      return responseData.error?.message || responseData.message || defaultMessage;
    }

    // Fallback to axios error message
    return axiosError.message || defaultMessage;
  }

  // Check for error with message property
  if ("message" in error) {
    return (error as { message: string }).message;
  }

  return defaultMessage;
}

export function useSchedules(params?: ScheduleListParams) {
  return useQuery({
    queryKey: ["schedules", params],
    queryFn: () => scheduleService.list(params),
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

export function useSchedule(id: string) {
  return useQuery({
    queryKey: ["schedules", id],
    queryFn: () => scheduleService.getById(id),
    enabled: !!id,
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

export function useCreateSchedule() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateScheduleFormData) => scheduleService.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["schedules"] });
      toast.success("Schedule created successfully");
    },
    onError: (error: unknown) => {
      const message =
        error instanceof Error ? error.message : "Failed to create schedule";
      toast.error(message);
    },
  });
}

export function useUpdateSchedule() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateScheduleFormData }) =>
      scheduleService.update(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["schedules"] });
      queryClient.invalidateQueries({ queryKey: ["schedules", variables.id] });
      toast.success("Schedule updated successfully");
    },
    onError: (error: unknown) => {
      const message =
        error instanceof Error ? error.message : "Failed to update schedule";
      toast.error(message);
    },
  });
}

export function useDeleteSchedule() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => scheduleService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["schedules"] });
      toast.success("Schedule deleted successfully");
    },
    onError: (error: unknown) => {
      const message =
        error instanceof Error ? error.message : "Failed to delete schedule";
      toast.error(message);
    },
  });
}
