"use client";

import { useQuery, useMutation, useQueryClient, useInfiniteQuery } from "@tanstack/react-query";
import { monthlyTargetService } from "../services/monthlyTargetService";
import type {
  BulkCreateMonthlyTargetFormData,
  UpdateMonthlyTargetFormData,
  CreateGroupTargetWithUserAssignmentFormData,
  BulkSetTargetFormData,
} from "../schemas/monthly-target.schema";

export function useMonthlyTargets(
  params?: {
    page?: number;
    per_page?: number;
    group_id?: string;
    user_id?: string;
    year?: number;
    month?: number;
    search?: string;
    manager_id?: string;
    scope?: "all" | "user" | "group" | "brick";
  },
  options?: { enabled?: boolean }
) {
  return useQuery({
    queryKey: ["monthly-targets", params],
    queryFn: () => monthlyTargetService.list(params),
    ...options,
  });
}

export function useInfiniteMonthlyTargets(
  params?: {
    per_page?: number;
    group_id?: string;
    user_id?: string;
    year?: number;
    month?: number;
    search?: string;
    manager_id?: string;
    scope?: "all" | "user" | "group" | "brick";
  },
  options?: { enabled?: boolean }
) {
  return useInfiniteQuery({
    queryKey: ["monthly-targets", "infinite", params],
    queryFn: ({ pageParam = 1 }) => 
      monthlyTargetService.list({ ...params, page: pageParam }),
    getNextPageParam: (lastPage) => {
      if (lastPage.meta.pagination.has_next) {
        return lastPage.meta.pagination.page + 1;
      }
      return undefined;
    },
    initialPageParam: 1,
    ...options,
  });
}

export function useMonthlyTarget(id: string) {
  return useQuery({
    queryKey: ["monthly-target", id],
    queryFn: () => monthlyTargetService.getById(id),
    enabled: !!id,
  });
}

export function useUserEffectiveTarget(params: {
  user_id: string;
  year: number;
  month: number;
}) {
  return useQuery({
    queryKey: ["monthly-target", "user-effective", params],
    queryFn: () => monthlyTargetService.getUserEffectiveTarget(params),
    enabled: !!params.user_id && !!params.year && !!params.month,
  });
}

export function useCreateMonthlyTarget() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: BulkCreateMonthlyTargetFormData) =>
      monthlyTargetService.bulkCreate(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["monthly-targets"] });
    },
  });
}

export function useUpdateMonthlyTarget() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      id,
      data,
    }: {
      id: string;
      data: UpdateMonthlyTargetFormData;
    }) => monthlyTargetService.update(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["monthly-targets"] });
      queryClient.invalidateQueries({
        queryKey: ["monthly-target", variables.id],
      });
    },
  });
}

export function useDeleteMonthlyTarget() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => monthlyTargetService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["monthly-targets"] });
    },
  });
}

export function useCreateGroupTargetWithUserAssignment() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateGroupTargetWithUserAssignmentFormData) =>
      monthlyTargetService.createGroupTargetWithUserAssignment(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["monthly-targets"] });
    },
  });
}


export function useBulkSetMonthlyTarget() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: BulkSetTargetFormData) =>
      monthlyTargetService.bulkSetTarget(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["monthly-targets"] });
    },
  });
}

