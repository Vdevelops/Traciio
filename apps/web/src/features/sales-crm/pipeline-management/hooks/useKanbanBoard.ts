"use client";

import { useState, useCallback } from "react";
import { useStages } from "./useStages";
import { useDealsByStage, useMoveDeal, useUpdateDeal } from "./useDeals";
import type { Deal, PipelineStage, DealFilters } from "../types";
import type { DealUpdateData } from "../schemas/deal.schema";

interface UseKanbanBoardParams {
  readonly filters?: DealFilters;
}

export function useKanbanBoard(params: UseKanbanBoardParams = {}) {
  const { filters } = params;
  const { data: pipelines, isLoading: stagesLoading } = useStages();
  const { data: dealsByStage, isLoading: dealsLoading } = useDealsByStage(filters);
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

  const handleDrop = useCallback(
    async (e: React.DragEvent, targetStage: PipelineStage) => {
      e.preventDefault();
      
      if (!draggedDeal || draggedDeal.stage_id === targetStage.id) {
        setDraggedDeal(null);
        return;
      }

      try {
        await moveDeal.mutateAsync({
          deal_id: draggedDeal.id,
          stage_id: targetStage.id,
          order: 0,
        });
      } catch (error) {
        // Error handled by mutation
      }
      
      setDraggedDeal(null);
    },
    [draggedDeal, moveDeal]
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

  // Ensure pipelines is an array
  const pipelinesArray = Array.isArray(pipelines) ? pipelines : [];

  // Ensure dealsByStage is an object
  const dealsByStageObject = dealsByStage && typeof dealsByStage === "object" && !Array.isArray(dealsByStage)
    ? dealsByStage
    : {};

  return {
    pipelines: pipelinesArray,
    dealsByStage: dealsByStageObject,
    isLoading: stagesLoading || dealsLoading,
    editingDeal,
    draggedDeal,
    handleDragStart,
    handleDragOver,
    handleDrop,
    handleUpdateDeal,
    openEditDialog,
    closeEditDialog,
    isUpdating: updateDeal.isPending || moveDeal.isPending,
  };
}
