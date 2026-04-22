"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { activityTypeService } from "../services/activityTypeService";
import type { ActivityType } from "../types/activity-type";

export function useActivityTypes(params?: {
  status?: string;
}) {
  return useQuery({
    queryKey: ["activity-types", params],
    queryFn: () => activityTypeService.list(params),
  });
}

export function useActivityType(id: string) {
  return useQuery({
    queryKey: ["activity-types", id],
    queryFn: () => activityTypeService.getById(id),
    enabled: !!id,
  });
}

export function useCreateActivityType() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: {
      name: string;
      code: string;
      description?: string;
      icon?: string;
      badge_color?: string;
      status?: string;
      order?: number;
    }) => activityTypeService.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["activity-types"] });
    },
  });
}

export function useUpdateActivityType() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: {
      id: string;
      data: {
        name?: string;
        code?: string;
        description?: string;
        icon?: string;
        badge_color?: string;
        status?: string;
        order?: number;
      };
    }) => activityTypeService.update(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["activity-types"] });
      queryClient.invalidateQueries({ queryKey: ["activity-types", variables.id] });
    },
  });
}

export function useDeleteActivityType() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => activityTypeService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["activity-types"] });
      queryClient.invalidateQueries({ queryKey: ["activities"] }); // Invalidate activities too since they depend on types
    },
  });
}

