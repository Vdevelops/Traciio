import { useQuery, useMutation, useQueryClient, type QueryClient } from "@tanstack/react-query";
import type { Deal, DealFilters } from "../types";
import type { DealFormData, DealUpdateData, DealMoveData } from "../schemas/deal.schema";
import * as dealService from "../services/dealService";

type DealDetailResponse = Awaited<ReturnType<typeof dealService.getDeal>>;
type DealListResponse = Awaited<ReturnType<typeof dealService.getDeals>>;
type DealsByStageResponse = Awaited<ReturnType<typeof dealService.getDealsByStage>>;

export const dealKeys = {
  all: ["deals"] as const,
  lists: () => [...dealKeys.all, "list"] as const,
  list: (filters?: DealFilters, page?: number, limit?: number) => 
    [...dealKeys.lists(), { filters, page, limit }] as const,
  byStage: (filters?: DealFilters) => [...dealKeys.all, "by-stage", filters] as const,
  details: () => [...dealKeys.all, "detail"] as const,
  detail: (id: string) => [...dealKeys.details(), id] as const,
};

function isDealListQueryKey(queryKey: readonly unknown[]) {
  return queryKey[0] === "deals" && queryKey[1] === "list";
}

function isDealsByStageQueryKey(queryKey: readonly unknown[]) {
  return queryKey[0] === "deals" && queryKey[1] === "by-stage";
}

function invalidateDealDerivedQueries(queryClient: QueryClient) {
  return Promise.all([
    queryClient.invalidateQueries({ queryKey: ["product-analytics"] }),
    queryClient.invalidateQueries({ queryKey: ["sales-overview"] }),
    queryClient.invalidateQueries({ queryKey: ["dashboard"] }),
  ]);
}

function replaceDealInList(current: DealListResponse | undefined, updatedDeal: Deal) {
  if (!current?.data) return current;

  return {
    ...current,
    data: current.data.map((deal) => (deal.id === updatedDeal.id ? { ...deal, ...updatedDeal } : deal)),
  };
}

function replaceDealInStageGroups(current: DealsByStageResponse | undefined, updatedDeal: Deal) {
  if (!current?.data) return current;

  const wasPresent = Object.values(current.data).some((deals) =>
    deals.some((deal) => deal.id === updatedDeal.id)
  );
  if (!wasPresent) return current;

  const nextGroups = Object.entries(current.data).reduce<Record<string, Deal[]>>((groups, [stageId, deals]) => {
    groups[stageId] = deals.filter((deal) => deal.id !== updatedDeal.id);
    return groups;
  }, {});

  nextGroups[updatedDeal.stage_id] = [{ ...updatedDeal }, ...(nextGroups[updatedDeal.stage_id] ?? [])];

  return {
    ...current,
    data: nextGroups,
  };
}

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
      queryClient.invalidateQueries({ queryKey: ["leads"] });
      queryClient.invalidateQueries({ queryKey: ["leads", "analytics"] });
      queryClient.invalidateQueries({ queryKey: ["accounts"] });
      queryClient.invalidateQueries({ queryKey: ["contacts"] });
      queryClient.invalidateQueries({ queryKey: ["tasks"] });
      void invalidateDealDerivedQueries(queryClient);
    },
  });
}

export function useUpdateDeal() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: DealUpdateData }) => 
      dealService.updateDeal(id, data),
    onSuccess: (response, variables) => {
      const updatedDeal = response.data;

      queryClient.setQueryData<DealDetailResponse>(dealKeys.detail(variables.id), response);
      queryClient.setQueriesData<DealListResponse>(
        { predicate: (query) => isDealListQueryKey(query.queryKey) },
        (current) => replaceDealInList(current, updatedDeal)
      );
      queryClient.setQueriesData<DealsByStageResponse>(
        { predicate: (query) => isDealsByStageQueryKey(query.queryKey) },
        (current) => replaceDealInStageGroups(current, updatedDeal)
      );

      queryClient.invalidateQueries({ queryKey: dealKeys.all });
      void invalidateDealDerivedQueries(queryClient);
    },
  });
}

export function useMoveDeal() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: DealMoveData) => dealService.moveDeal(data),
    onSuccess: async (response, variables) => {
      const updatedDeal = response.data;

      queryClient.setQueryData<DealDetailResponse>(dealKeys.detail(variables.deal_id), response);
      queryClient.setQueriesData<DealListResponse>(
        { predicate: (query) => isDealListQueryKey(query.queryKey) },
        (current) => replaceDealInList(current, updatedDeal)
      );
      queryClient.setQueriesData<DealsByStageResponse>(
        { predicate: (query) => isDealsByStageQueryKey(query.queryKey) },
        (current) => replaceDealInStageGroups(current, updatedDeal)
      );

      await queryClient.invalidateQueries({ queryKey: dealKeys.all });
      await Promise.all([
        invalidateDealDerivedQueries(queryClient),
        queryClient.refetchQueries({
          queryKey: dealKeys.detail(variables.deal_id),
          type: "active",
        }),
        queryClient.refetchQueries({
          predicate: (query) => isDealListQueryKey(query.queryKey),
          type: "active",
        }),
        queryClient.refetchQueries({
          predicate: (query) => isDealsByStageQueryKey(query.queryKey),
          type: "active",
        }),
      ]);
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
      void invalidateDealDerivedQueries(queryClient);
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
