"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { industryService } from "../services/industryService";
import type { CreateIndustryRequest, UpdateIndustryRequest } from "../types/industry";
import { toast } from "sonner";

export function useIndustries(params?: {
  page?: number;
  per_page?: number;
  search?: string;
  is_active?: boolean;
  sort_by?: string;
  sort_order?: "asc" | "desc";
}) {
  return useQuery({
    queryKey: ["industries", params],
    queryFn: () => industryService.list(params),
  });
}

export function useAllIndustries() {
  return useQuery({
    queryKey: ["industries", "all"],
    queryFn: () => industryService.listAll(),
  });
}

export function useIndustry(id: string) {
  return useQuery({
    queryKey: ["industries", id],
    queryFn: () => industryService.getById(id),
    enabled: !!id,
  });
}

export function useCreateIndustry() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateIndustryRequest) => industryService.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["industries"] });
      toast.success("Industry created successfully");
    },
    onError: () => {
      toast.error("Failed to create industry");
    },
  });
}

export function useUpdateIndustry() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateIndustryRequest }) =>
      industryService.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["industries"] });
      toast.success("Industry updated successfully");
    },
    onError: () => {
      toast.error("Failed to update industry");
    },
  });
}

export function useDeleteIndustry() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => industryService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["industries"] });
      queryClient.invalidateQueries({ queryKey: ["leads", "form-data"] }); // Invalidate form data cache
      toast.success("Industry deleted successfully");
    },
    onError: () => {
      toast.error("Failed to delete industry");
    },
  });
}

