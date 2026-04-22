import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import type { DealFilters } from "../types";
import type { DealFormData, DealUpdateData, DealMoveData } from "../schemas/deal.schema";
import * as dealService from "../services/dealService";

export const dealKeys = {
  all: ["deals"] as const,
  lists: () => [...dealKeys.all, "list"] as const,
  list: (filters?: DealFilters, page?: number, limit?: number) => 
    [...dealKeys.lists(), { filters, page, limit }] as const,
  byStage: (filters?: DealFilters) => [...dealKeys.all, "by-stage", filters] as const,
  details: () => [...dealKeys.all, "detail"] as const,
  detail: (id: string) => [...dealKeys.details(), id] as const,
};

export function useDeals(
  filters?: DealFilters, 
  page: number = 1, 
  limit: number = 20,
  sort: string = "created_at",
  order: "asc" | "desc" = "desc"
) {
  return useQuery({
    queryKey: [...dealKeys.list(filters, page, limit), sort, order],
    queryFn: () => dealService.getDeals(filters, page, limit, sort, order),
    // Return full response to access meta.pagination
    select: (response) => response,
    staleTime: 30000,
  });
}

export function useDealsByStage(filters?: DealFilters) {
  return useQuery({
    queryKey: dealKeys.byStage(filters),
    queryFn: () => dealService.getDealsByStage(filters),
    select: (response) => response.data,
    staleTime: 30000,
  });
}

export function useDeal(id: string) {
  return useQuery({
    queryKey: dealKeys.detail(id),
    queryFn: () => dealService.getDeal(id),
    select: (response) => response.data,
    enabled: !!id,
  });
}

export function useCreateDeal() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: DealFormData) => dealService.createDeal(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: dealKeys.all });
    },
  });
}

export function useUpdateDeal() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: DealUpdateData }) => 
      dealService.updateDeal(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: dealKeys.all });
    },
  });
}

export function useMoveDeal() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: DealMoveData) => dealService.moveDeal(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: dealKeys.all });
    },
  });
}

// Alias for historical naming used by MoveStageModal.
// The API operation is the same as moveDeal: patch deal stage.
export const useMoveStage = useMoveDeal;

export function useDeleteDeal() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => dealService.deleteDeal(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: dealKeys.all });
    },
  });
}

// Hook for deal-related visit reports
export function useDealVisitReports(dealId: string, page: number = 1, limit: number = 20) {
  return useQuery({
    queryKey: [...dealKeys.detail(dealId), "visit-reports", page, limit],
    queryFn: () => dealService.getDealVisitReports(dealId, page, limit),
    select: (response) => response.data,
    enabled: !!dealId,
  });
}

export function useDealActivities(dealId: string, page: number = 1, limit: number = 20) {
  return useQuery({
    queryKey: [...dealKeys.detail(dealId), "activities", page, limit],
    queryFn: () => dealService.getDealActivities(dealId, page, limit),
    select: (response) => response.data,
    enabled: !!dealId,
  });
}

// Hook for deal history
export function useDealHistory(dealId: string) {
  return useQuery({
    queryKey: [...dealKeys.detail(dealId), "history"],
    queryFn: () => dealService.getDealHistory(dealId),
    select: (response) => response.data,
    enabled: !!dealId,
  });
}
