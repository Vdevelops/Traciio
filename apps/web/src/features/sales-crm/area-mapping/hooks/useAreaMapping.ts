import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { areaMappingService } from "../services/areaMappingService";
import type {
  CreateTerritoryRequest,
  UpdateTerritoryRequest,
  CreateAreaCaptureRequest,
} from "../types";
import { toast } from "sonner";

// Query Keys
export const areaMappingKeys = {
  all: ["area-mapping"] as const,
  territories: () => [...areaMappingKeys.all, "territories"] as const,
  territory: (id: string) => [...areaMappingKeys.territories(), id] as const,
  captures: () => [...areaMappingKeys.all, "captures"] as const,
  coverage: () => [...areaMappingKeys.all, "coverage"] as const,
  heatmap: (params?: Record<string, any>) =>
    [...areaMappingKeys.all, "heatmap", params] as const,
};

// Territory Hooks
export function useTerritories(params?: {
  page?: number;
  page_size?: number;
  search?: string;
  assigned_to?: string;
}) {
  return useQuery({
    queryKey: [...areaMappingKeys.territories(), params],
    queryFn: () => areaMappingService.listTerritories(params),
  });
}

export function useTerritory(id: string) {
  return useQuery({
    queryKey: areaMappingKeys.territory(id),
    queryFn: () => areaMappingService.getTerritory(id),
    enabled: !!id,
  });
}

export function useCreateTerritory() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateTerritoryRequest) =>
      areaMappingService.createTerritory(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: areaMappingKeys.territories() });
      toast.success("Territory created successfully");
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.message || "Failed to create territory");
    },
  });
}

export function useUpdateTerritory() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateTerritoryRequest }) =>
      areaMappingService.updateTerritory(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: areaMappingKeys.territories() });
      queryClient.invalidateQueries({ queryKey: areaMappingKeys.territory(variables.id) });
      toast.success("Territory updated successfully");
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.message || "Failed to update territory");
    },
  });
}

export function useDeleteTerritory() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => areaMappingService.deleteTerritory(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: areaMappingKeys.territories() });
      toast.success("Territory deleted successfully");
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.message || "Failed to delete territory");
    },
  });
}

// Area Capture Hooks
export function useAreaCaptures(params?: {
  page?: number;
  per_page?: number;
  visit_report_id?: string;
  capture_type?: string;
  captured_after?: string;
  captured_before?: string;
}) {
  return useQuery({
    queryKey: [...areaMappingKeys.captures(), params],
    queryFn: () => areaMappingService.listCaptures(params),
  });
}

export function useCreateAreaCapture() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateAreaCaptureRequest) =>
      areaMappingService.createCapture(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: areaMappingKeys.captures() });
      toast.success("Location captured successfully");
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.message || "Failed to capture location");
    },
  });
}


// Coverage Analysis Hooks
export function useCoverageAnalysis(params?: {
  territory_id: string;
  start_date: string;
  end_date: string;
}) {
  return useQuery({
    queryKey: [...areaMappingKeys.coverage(), params],
    queryFn: () => {
      if (!params || !params.territory_id || !params.start_date || !params.end_date) {
        throw new Error("Missing required parameters");
      }
      return areaMappingService.getCoverageAnalysis(params);
    },
    enabled: !!params?.territory_id && !!params?.start_date && !!params?.end_date,
  });
}

// Heatmap Hook
export function useHeatmapData(params?: {
  territory_id?: string;
  capture_type?: string;
  start_date?: string;
  end_date?: string;
}) {
  return useQuery({
    queryKey: areaMappingKeys.heatmap(params),
    queryFn: () => areaMappingService.getHeatmapData(params),
  });
}
