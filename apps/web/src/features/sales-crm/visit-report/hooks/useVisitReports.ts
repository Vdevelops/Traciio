"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { visitReportService } from "../services/visitReportService";
import { activityService } from "../services/activityService";
import type {
  CreateVisitReportFormData,
  UpdateVisitReportFormData,
  CheckInFormData,
  CheckOutFormData,
  RejectFormData,
  UploadPhotoFormData,
  SubmitVisitReportFormData,
} from "../schemas/visit-report.schema";

export function useVisitReports(params?: {
  page?: number;
  per_page?: number;
  search?: string;
  status?: string;
  account_id?: string;
  deal_id?: string;
  sales_rep_id?: string;
  start_date?: string;
  end_date?: string;
}) {
  return useQuery({
    queryKey: ["visit-reports", params],
    queryFn: () => visitReportService.list(params),
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

export function useVisitReport(id: string) {
  return useQuery({
    queryKey: ["visit-reports", id],
    queryFn: () => visitReportService.getById(id),
    enabled: !!id,
  });
}

export function useCreateVisitReport() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateVisitReportFormData) => visitReportService.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["visit-reports"] });
      queryClient.invalidateQueries({ queryKey: ["activities"] });
    },
  });
}

export function useUpdateVisitReport() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateVisitReportFormData }) =>
      visitReportService.update(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["visit-reports"] });
      queryClient.invalidateQueries({ queryKey: ["visit-reports", variables.id] });
      queryClient.invalidateQueries({ queryKey: ["activities"] });
    },
  });
}

export function useDeleteVisitReport() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => visitReportService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["visit-reports"] });
      queryClient.invalidateQueries({ queryKey: ["activities"] });
    },
  });
}

export function useCheckIn() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({
      id,
      data,
      options,
    }: {
      id: string;
      data: CheckInFormData;
      options?: {
        photo?: File;
        deviceGPS?: {
          latitude: number;
          longitude: number;
          accuracy?: number;
          timestamp?: number;
        };
        photoGPS?: {
          latitude: number;
          longitude: number;
          timestamp?: number;
        };
      };
    }) => {
      try {
        const result = await visitReportService.checkIn(id, data, options);
        return result;
      } catch (error) {
        // Re-throw with more context if it's an axios error
        if (error && typeof error === "object" && "response" in error) {
          const axiosError = error as { 
            response?: { 
              data?: { 
                error?: { message?: string; code?: string };
                message?: string;
              };
              status?: number;
              statusText?: string;
            };
            message?: string;
          };
          
          const errorMessage = axiosError.response?.data?.error?.message 
            || axiosError.response?.data?.message
            || axiosError.message
            || "Failed to check in";
          const errorCode = axiosError.response?.data?.error?.code || "UNKNOWN_ERROR";
          
          // Create a new Error with the message
          const enhancedError = new Error(`${errorCode}: ${errorMessage}`);
          // Preserve original error for debugging
          (enhancedError as unknown as { originalError: unknown }).originalError = error;
          throw enhancedError;
        }
        
        // If it's already an Error, re-throw it
        if (error instanceof Error) {
          throw error;
        }
        
        // Otherwise, wrap it in an Error
        throw new Error(`Check-in failed: ${String(error)}`);
      }
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["visit-reports"] });
      queryClient.invalidateQueries({ queryKey: ["visit-reports", variables.id] });
      queryClient.invalidateQueries({ queryKey: ["activities"] });
    },
    onError: (error) => {
    },
  });
}

export function useCheckOut() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: CheckOutFormData }) =>
      visitReportService.checkOut(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["visit-reports"] });
      queryClient.invalidateQueries({ queryKey: ["visit-reports", variables.id] });
      queryClient.invalidateQueries({ queryKey: ["activities"] });
    },
  });
}

export function useSubmitVisitReport() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: SubmitVisitReportFormData }) =>
      visitReportService.submit(id, data),
    onSuccess: (_, variables) => {
      // Invalidate visit reports
      queryClient.invalidateQueries({ queryKey: ["visit-reports"] });
      queryClient.invalidateQueries({ queryKey: ["visit-reports", variables.id] });
      // Invalidate related data that might be affected by auto-triggers
      queryClient.invalidateQueries({ queryKey: ["leads"] }); // Lead status might change
      queryClient.invalidateQueries({ queryKey: ["tasks"] }); // Auto-tasks created
      queryClient.invalidateQueries({ queryKey: ["activities"] });
      queryClient.invalidateQueries({ queryKey: ["notifications"] }); // Manager notified
    },
  });
}

export function useApproveVisitReport() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => visitReportService.approve(id),
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: ["visit-reports"] });
      queryClient.invalidateQueries({ queryKey: ["visit-reports", id] });
      queryClient.invalidateQueries({ queryKey: ["activities"] });
    },
  });
}

export function useRejectVisitReport() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: RejectFormData }) =>
      visitReportService.reject(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["visit-reports"] });
      queryClient.invalidateQueries({ queryKey: ["visit-reports", variables.id] });
      queryClient.invalidateQueries({ queryKey: ["activities"] });
    },
  });
}

export function useUploadPhoto() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, file }: { id: string; file: File }) =>
      visitReportService.uploadPhoto(id, file),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["visit-reports"] });
      queryClient.invalidateQueries({ queryKey: ["visit-reports", variables.id] });
    },
  });
}

// Activity hooks
export function useActivities(params?: {
  page?: number;
  per_page?: number;
  type?: string;
  account_id?: string;
  contact_id?: string;
  deal_id?: string;
  user_id?: string;
  start_date?: string;
  end_date?: string;
}) {
  return useQuery({
    queryKey: ["activities", params],
    queryFn: () => activityService.list(params),
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

export function useActivityTimeline(params?: {
  account_id?: string;
  contact_id?: string;
  deal_id?: string;
  user_id?: string;
  start_date?: string;
  end_date?: string;
  limit?: number;
}) {
  return useQuery({
    queryKey: ["activities", "timeline", params],
    queryFn: () => activityService.getTimeline(params),
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

export function useCreateActivity() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: {
      activity_type_id: string;
      account_id?: string;
      contact_id?: string;
      deal_id?: string;
      description: string;
      timestamp: string;
      metadata?: Record<string, unknown>;
    }) => activityService.create(data),
    onSuccess: (_, variables) => {
      // Invalidate all activities queries using predicate for more comprehensive matching
      queryClient.invalidateQueries({
        predicate: (query) => {
          const key = query.queryKey;
          // Match any query that starts with ["activities"]
          if (key[0] === "activities") {
            return true;
          }
          // Match deal-specific activity queries
          if (variables.deal_id && key[0] === "deals" && key[1] === variables.deal_id) {
            return true;
          }
          return false;
        },
      });
    },
  });
}

