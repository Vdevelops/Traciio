"use client";

import { useState, useCallback, useMemo } from "react";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useStages } from "./useStages";
import { useMoveDeal, useUpdateDeal, dealKeys } from "./useDeals";
import * as dealService from "../services/dealService";
import type { Deal, PipelineStage, DealFilters } from "../types";
import type { DealUpdateData } from "../schemas/deal.schema";
import { toast } from "sonner";

interface UseProgressiveKanbanBoardParams {
  readonly filters?: DealFilters;
  readonly pageSize?: number;
}

// Cached data structure returned by the initial query
interface KanbanInitialData {
  deals: Record<string, Deal[]>;
  pages: Record<string, number>;
  hasMore: Record<string, boolean>;
}

export function useProgressiveKanbanBoard(
  params: UseProgressiveKanbanBoardParams = {}
) {
  const { filters, pageSize = 20 } = params;
  const queryClient = useQueryClient();

  // Fetch stages first (this is fast)
  const { data: pipelines, isLoading: stagesLoading } = useStages();
  const pipelinesArray = Array.isArray(pipelines) ? pipelines : [];

  // Local state only for pagination beyond the initial load
  const [extraDeals, setExtraDeals] = useState<Record<string, Deal[]>>({});
  const [currentPages, setCurrentPages] = useState<Record<string, number>>({});
  const [paginationHasMore, setPaginationHasMore] = useState<Record<string, boolean>>({});
  const [loadingMore, setLoadingMore] = useState<Record<string, boolean>>({});

  // Initial load - return actual data so TanStack Query caches it properly
  const { data: initialData, isLoading: initialLoading } = useQuery<KanbanInitialData>({
    queryKey: [...dealKeys.byStage(filters), "paginated-initial"],
    queryFn: async () => {
      const promises = pipelinesArray.map(async (stage) => {
        const response = await dealService.getDeals(
          { ...filters, stage_id: stage.id },
          1,
          pageSize
        );

        const currentPage = response.meta?.pagination?.page ?? 0;
        const totalPages = response.meta?.pagination?.total_pages ?? 0;
        const hasMore = currentPage < totalPages;

        return {
          stageId: stage.id,
          deals: response.data || [],
          hasMore,
        };
      });

      const results = await Promise.all(promises);

      const deals: Record<string, Deal[]> = {};
      const pages: Record<string, number> = {};
      const hasMore: Record<string, boolean> = {};

      results.forEach((result) => {
        deals[result.stageId] = result.deals;
        pages[result.stageId] = 1;
        hasMore[result.stageId] = result.hasMore;
      });

      // Reset pagination state when initial data is refetched
      setExtraDeals({});
      setCurrentPages(pages);
      setPaginationHasMore(hasMore);

      return { deals, pages, hasMore };
    },
    enabled: pipelinesArray.length > 0,
    staleTime: 30000,
  });

  // Merge initial cached data with extra paginated data
  const mergedHasMore = useMemo(() => {
    if (!initialData) return paginationHasMore;
    return { ...initialData.hasMore, ...paginationHasMore };
  }, [initialData, paginationHasMore]);

  const mergedPages = useMemo(() => {
    if (!initialData) return currentPages;
    return { ...initialData.pages, ...currentPages };
  }, [initialData, currentPages]);

  // Function to load more deals for a specific stage
  const fetchNextPageForStage = useCallback(
    async (stageId: string) => {
      if (loadingMore[stageId] || !mergedHasMore[stageId]) {
        return;
      }

      const page = mergedPages[stageId] || 1;
      const nextPage = page + 1;

      setLoadingMore((prev) => ({ ...prev, [stageId]: true }));

      try {
        const response = await dealService.getDeals(
          { ...filters, stage_id: stageId },
          nextPage,
          pageSize
        );

        const newDeals = response.data || [];
        const currentPageNum = response.meta?.pagination?.page ?? 0;
        const totalPages = response.meta?.pagination?.total_pages ?? 0;
        const hasMore = currentPageNum < totalPages;

        // Append new deals, filtering out duplicates
        setExtraDeals((prev) => {
          const existing = prev[stageId] || [];
          const existingIds = new Set([
            ...existing.map((d) => d.id),
            ...(initialData?.deals[stageId] || []).map((d) => d.id),
          ]);
          const uniqueNewDeals = newDeals.filter((d) => !existingIds.has(d.id));

          return {
            ...prev,
            [stageId]: [...existing, ...uniqueNewDeals],
          };
        });

        setCurrentPages((prev) => ({ ...prev, [stageId]: nextPage }));
        setPaginationHasMore((prev) => ({ ...prev, [stageId]: hasMore }));
      } catch (error) {
        console.error(`[Pagination] Error loading more deals for stage ${stageId}:`, error);
      } finally {
        setLoadingMore((prev) => ({ ...prev, [stageId]: false }));
      }
    },
    [loadingMore, mergedHasMore, mergedPages, filters, pageSize, initialData]
  );

  // Derive deals from cached query data + extra paginated deals
  const dealsByStage = useMemo(() => {
    const result: Record<string, Deal[]> = {};
    pipelinesArray.forEach((stage) => {
      const base = initialData?.deals[stage.id] || [];
      const extra = extraDeals[stage.id] || [];
      result[stage.id] = [...base, ...extra];
    });
    return result;
  }, [pipelinesArray, initialData, extraDeals]);

  // Create loading state per stage (initial load only)
  const loadingByStage = useMemo(() => {
    const result: Record<string, boolean> = {};
    pipelinesArray.forEach((stage) => {
      result[stage.id] = initialLoading;
    });
    return result;
  }, [pipelinesArray, initialLoading]);

  // Has next page per stage
  const hasNextPageByStage = useMemo(() => {
    const result: Record<string, boolean> = {};
    pipelinesArray.forEach((stage) => {
      result[stage.id] = mergedHasMore[stage.id] ?? false;
    });
    return result;
  }, [pipelinesArray, mergedHasMore]);

  // Is fetching next page per stage
  const isFetchingNextPageByStage = useMemo(() => {
    const result: Record<string, boolean> = {};
    pipelinesArray.forEach((stage) => {
      result[stage.id] = loadingMore[stage.id] ?? false;
    });
    return result;
  }, [pipelinesArray, loadingMore]);

  const moveDeal = useMoveDeal();
  const updateDeal = useUpdateDeal();

  const [editingDeal, setEditingDeal] = useState<Deal | null>(null);
  const [draggedDeal, setDraggedDeal] = useState<Deal | null>(null);

  const handleDragStart = useCallback((e: React.DragEvent, deal: Deal) => {
    setDraggedDeal(deal);
    e.dataTransfer.effectAllowed = "move";
    e.dataTransfer.setData("text/plain", deal.id);
  }, []);

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.dataTransfer.dropEffect = "move";
  }, []);

  const moveDealToStage = useCallback(
    async (dealId: string, stageId: string, reason?: string) => {
      await moveDeal.mutateAsync({
        deal_id: dealId,
        stage_id: stageId,
        order: 0,
        reason,
      });

      setExtraDeals({});
      queryClient.invalidateQueries({ queryKey: dealKeys.byStage(filters) });
      setDraggedDeal(null);
    },
    [filters, moveDeal, queryClient]
  );

  const handleDrop = useCallback(
    async (e: React.DragEvent, targetStage: PipelineStage) => {
      e.preventDefault();

      if (!draggedDeal || draggedDeal.stage_id === targetStage.id) {
        setDraggedDeal(null);
        return;
      }

      try {
        await moveDealToStage(draggedDeal.id, targetStage.id);
      } catch (error) {
        // Error handled by mutation, but we should show a toast
        toast.error(
          error instanceof Error
            ? error.message
            : "Failed to move deal. The action was reverted."
        );
      }

      setDraggedDeal(null);
    },
    [draggedDeal, moveDealToStage]
  );

  const handleUpdateDeal = useCallback(
    async (data: DealUpdateData) => {
      if (!editingDeal) return;

      await updateDeal.mutateAsync({
        id: editingDeal.id,
        data,
      });

      setEditingDeal(null);
    },
    [editingDeal, updateDeal]
  );

  const openEditDialog = useCallback((deal: Deal) => {
    setEditingDeal(deal);
  }, []);

  const closeEditDialog = useCallback(() => {
    setEditingDeal(null);
  }, []);

  return {
    pipelines: pipelinesArray,
    dealsByStage,
    loadingByStage,
    hasNextPageByStage,
    isFetchingNextPageByStage,
    fetchNextPageForStage,
    isLoading: stagesLoading || initialLoading,
    editingDeal,
    draggedDeal,
    handleDragStart,
    handleDragOver,
    handleDrop,
    moveDealToStage,
    clearDraggedDeal: () => setDraggedDeal(null),
    handleUpdateDeal,
    openEditDialog,
    closeEditDialog,
    isUpdating: updateDeal.isPending || moveDeal.isPending,
  };
}
