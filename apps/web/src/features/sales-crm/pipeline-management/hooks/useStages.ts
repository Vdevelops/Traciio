import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import type { StageFormData, StageUpdateData, StageReorderData } from "../schemas/stage.schema";
import * as stageService from "../services/stageService";

export const stageKeys = {
  all: ["pipeline-stages"] as const,
  lists: () => [...stageKeys.all, "list"] as const,
  list: () => [...stageKeys.lists()] as const,
  withStats: () => [...stageKeys.all, "with-stats"] as const,
  details: () => [...stageKeys.all, "detail"] as const,
  detail: (id: string) => [...stageKeys.details(), id] as const,
};

export function useStages() {
  return useQuery({
    queryKey: stageKeys.list(),
    queryFn: async () => {
      const response = await stageService.getStages();
      // Ensure data is an array
      if (!response.data || !Array.isArray(response.data)) {
        return [];
      }
      return response.data;
    },
    staleTime: 60000,
  });
}

export function useStagesWithStats() {
  return useQuery({
    queryKey: stageKeys.withStats(),
    queryFn: () => stageService.getStagesWithStats(),
    select: (response) => response.data,
    staleTime: 30000,
  });
}

export function useStage(id: string) {
  return useQuery({
    queryKey: stageKeys.detail(id),
    queryFn: () => stageService.getStage(id),
    select: (response) => response.data,
    enabled: !!id,
  });
}

export function useCreateStage() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: StageFormData) => stageService.createStage(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: stageKeys.all });
    },
  });
}

export function useUpdateStage() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: StageUpdateData }) => 
      stageService.updateStage(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: stageKeys.detail(variables.id) });
      queryClient.invalidateQueries({ queryKey: stageKeys.all });
    },
  });
}

export function useReorderStages() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: StageReorderData[]) => stageService.reorderStages(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: stageKeys.all });
    },
  });
}

// Alias for useReorderStages
export const useUpdateStagesOrder = useReorderStages;

export function useDeleteStage() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => stageService.deleteStage(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: stageKeys.all });
    },
  });
}
