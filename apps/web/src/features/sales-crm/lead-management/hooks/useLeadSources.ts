"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { leadSourceService } from "../services/leadSourceService";
import type { CreateLeadSourceRequest, UpdateLeadSourceRequest } from "../types/lead-source";
import { toast } from "sonner";

export function useLeadSources(params?: {
  page?: number;
  per_page?: number;
  search?: string;
  is_active?: boolean;
  sort_by?: string;
  sort_order?: "asc" | "desc";
}) {
  return useQuery({
    queryKey: ["lead-sources", params],
    queryFn: () => leadSourceService.list(params),
  });
}

export function useAllLeadSources() {
  return useQuery({
    queryKey: ["lead-sources", "all"],
    queryFn: () => leadSourceService.listAll(),
    // Master data rarely changes — cache aggressively to avoid redundant requests
    staleTime: 10 * 60 * 1000,
    gcTime: 30 * 60 * 1000,
    refetchOnWindowFocus: false,
  });
}

export function useLeadSource(id: string) {
  return useQuery({
    queryKey: ["lead-sources", id],
    queryFn: () => leadSourceService.getById(id),
    enabled: !!id,
  });
}

export function useCreateLeadSource() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateLeadSourceRequest) => leadSourceService.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["lead-sources"] });
      queryClient.invalidateQueries({ queryKey: ["leads", "form-data"] }); // Invalidate form data cache
      toast.success("Lead source created successfully");
    },
    onError: () => {
      toast.error("Failed to create lead source");
    },
  });
}

export function useUpdateLeadSource() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateLeadSourceRequest }) =>
      leadSourceService.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["lead-sources"] });
      queryClient.invalidateQueries({ queryKey: ["leads", "form-data"] }); // Invalidate form data cache
      toast.success("Lead source updated successfully");
    },
    onError: () => {
      toast.error("Failed to update lead source");
    },
  });
}

export function useDeleteLeadSource() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => leadSourceService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["lead-sources"] });
      queryClient.invalidateQueries({ queryKey: ["leads", "form-data"] }); // Invalidate form data cache
      toast.success("Lead source deleted successfully");
    },
    onError: () => {
      toast.error("Failed to delete lead source");
    },
  });
}

