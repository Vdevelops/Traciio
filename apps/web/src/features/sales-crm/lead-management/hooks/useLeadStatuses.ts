"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { leadStatusService } from "../services/leadStatusService";
import type { CreateLeadStatusRequest, UpdateLeadStatusRequest } from "../types/lead-status";
import { toast } from "sonner";

export function useLeadStatuses(params?: {
  page?: number;
  per_page?: number;
  search?: string;
  is_active?: boolean;
  sort_by?: string;
  sort_order?: "asc" | "desc";
}) {
  return useQuery({
    queryKey: ["lead-statuses", params],
    queryFn: () => leadStatusService.list(params),
  });
}

export function useAllLeadStatuses() {
  return useQuery({
    queryKey: ["lead-statuses", "all"],
    queryFn: () => leadStatusService.listAll(),
  });
}

export function useLeadStatus(id: string) {
  return useQuery({
    queryKey: ["lead-statuses", id],
    queryFn: () => leadStatusService.getById(id),
    enabled: !!id,
  });
}

export function useCreateLeadStatus() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateLeadStatusRequest) => leadStatusService.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["lead-statuses"] });
      toast.success("Lead status created successfully");
    },
    onError: () => {
      toast.error("Failed to create lead status");
    },
  });
}

export function useUpdateLeadStatus() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateLeadStatusRequest }) =>
      leadStatusService.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["lead-statuses"] });
      toast.success("Lead status updated successfully");
    },
    onError: () => {
      toast.error("Failed to update lead status");
    },
  });
}

export function useDeleteLeadStatus() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => leadStatusService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["lead-statuses"] });
      toast.success("Lead status deleted successfully");
    },
    onError: () => {
      toast.error("Failed to delete lead status");
    },
  });
}

export function useSetDefaultLeadStatus() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => leadStatusService.setDefault(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["lead-statuses"] });
      toast.success("Default lead status updated successfully");
    },
    onError: () => {
      toast.error("Failed to set default lead status");
    },
  });
}
