"use client";

import { useState } from "react";
import { Plus } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { DealCard } from "./deal-card";
import type { Deal, DealFilters } from "../types";
import type { UpdateDealFormData } from "../schemas/deal.schema";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { DealForm } from "./deal-form";
import { CreateDealDialog } from "./create-deal-dialog";
import { MoveStageModal } from "./move-stage-modal";
import { PipelineFilters } from "./pipeline-filters";
import { StageScrollLoader } from "./stage-scroll-loader";
import { useProgressiveKanbanBoard } from "../hooks/useProgressiveKanbanBoard";
import { useTranslations } from "next-intl";
import { useHasPermission } from "@/features/auth/providers/permissions-provider";

interface KanbanBoardProps {
  readonly onDealClick?: (deal: Deal) => void;
}

export function KanbanBoard({ onDealClick }: KanbanBoardProps) {
  const t = useTranslations("pipelineManagement.kanban");
  const [filters, setFilters] = useState<DealFilters>({});
  
  const {
    pipelines,
    dealsByStage,
    loadingByStage,
    hasNextPageByStage,
    isFetchingNextPageByStage,
    fetchNextPageForStage,
    isLoading,
    editingDeal,
    handleDragStart,
    handleDragOver,
    handleDrop,
    clearDraggedDeal,
    draggedDeal,
    handleUpdateDeal,
    openEditDialog,
    closeEditDialog,
    isUpdating,
  } = useProgressiveKanbanBoard({ filters });
  const handleResetFilters = () => {
    setFilters({});
  };

  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [pendingStageMove, setPendingStageMove] = useState<{
    dealId: string;
    currentStageId: string;
    stageId: string;
  } | null>(null);
  const hasCreatePermission = useHasPermission("pipeline.opportunity-create");

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div>
          <div className="h-8 bg-muted animate-pulse rounded w-64 mb-2" />
          <div className="h-4 bg-muted animate-pulse rounded w-96" />
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
          {Array.from({ length: 4 }, (_, i) => (
            <Card key={`skeleton-${i}`} className="p-4">
              <div className="space-y-4">
                <div className="h-12 bg-muted animate-pulse rounded" />
                <div className="space-y-2">
                  {Array.from({ length: 2 }, (_, j) => (
                    <div key={`skeleton-inner-${i}-${j}`} className="h-32 bg-muted animate-pulse rounded" />
                  ))}
                </div>
              </div>
            </Card>
          ))}
        </div>
      </div>
    );
  }

  // Show empty state if no pipelines
  if (!pipelines || pipelines.length === 0) {
    return (
      <div className="space-y-6">
        {hasCreatePermission && (
          <div className="flex justify-end">
            <Button onClick={() => setIsCreateDialogOpen(true)}>
              <Plus className="h-4 w-4 mr-2" />
              {t("addOpportunity") || "Add Opportunity"}
            </Button>
          </div>
        )}
        <div className="text-center py-12 border border-dashed rounded-xl bg-muted/20">
          <p className="text-sm text-muted-foreground">{t("noStages") || "No pipeline stages found. Please create stages first."}</p>
        </div>
      </div>
    );
  }

  const handleCreateSuccess = () => {
    setIsCreateDialogOpen(false);
  };

  const handleDealClick = (deal: Deal) => {
    if (onDealClick) {
      onDealClick(deal);
    } else {
      openEditDialog(deal);
    }
  };

  const renderDeal = (deal: Deal) => (
    <div
      key={deal.id}
      draggable
      onDragStart={(e) => handleDragStart(e, deal)}
      className="cursor-grab active:cursor-grabbing outline-hidden"
      aria-label={t("dragToMove", { title: deal.title || 'deal' })}
      role="button"
      tabIndex={0}
      onKeyDown={(e) => {
        if (e.key === 'Enter' || e.key === ' ') {
          e.preventDefault();
          handleDealClick(deal);
        }
      }}
    >
      <DealCard
        deal={deal}
        onClick={() => handleDealClick(deal)}
      />
    </div>
  );

  return (
    <div className="space-y-6">
      {/* Filters with Add Button */}
      <PipelineFilters
        filters={filters}
        onFiltersChange={setFilters}
        onReset={handleResetFilters}
        onAdd={hasCreatePermission ? () => setIsCreateDialogOpen(true) : undefined}
        addLabel={t("addOpportunity") || "Add Opportunity"}
      />
      <div className="flex overflow-x-auto overflow-y-hidden pb-6 gap-6 -mx-6 px-6 scrollbar-thin scrollbar-thumb-muted-foreground/20">
        {pipelines.map((stage) => {
          const stageDeals = dealsByStage?.[stage.id] || [];
          
          return (
            <section
              key={stage.id}
              className="shrink-0 w-80 h-full flex flex-col"
              onDragOver={handleDragOver}
              onDrop={(e) => {
                if (draggedDeal && draggedDeal.stage_id !== stage.id && (stage.is_won || stage.is_lost)) {
                  e.preventDefault();
                  setPendingStageMove({
                    dealId: draggedDeal.id,
                    currentStageId: draggedDeal.stage_id,
                    stageId: stage.id,
                  });
                  return;
                }
                void handleDrop(e, stage);
              }}
              aria-label={`Drop zone for ${stage.name} stage`}
            >
                <div className="flex items-center gap-2.5 mb-4 shrink-0 pb-3 px-1 border-b border-border/40">
                <div
                  className="w-2.5 h-2.5 rounded-full shrink-0 ring-2 ring-offset-2 ring-offset-background/50 shadow-sm"
                  style={{
                    backgroundColor: stage.color,
                  }}
                />
                <h3 className="font-semibold text-sm truncate flex-1 tracking-tight text-foreground/80">{stage.name}</h3>
                <Badge variant="secondary" className="shrink-0 text-[10px] font-bold h-5 px-1.5 bg-muted/40 border-none text-muted-foreground">
                  {stageDeals.length}
                </Badge>
              </div>

              <div className="space-y-3 min-h-[200px] flex-1 overflow-y-auto pr-1">
                {(() => {
                  if (loadingByStage[stage.id]) {
                    return (
                      <div className="space-y-3">
                        {Array.from({ length: 3 }, (_, i) => (
                          <div key={`stage-skeleton-${stage.id}-${i}`} className="bg-muted/20 rounded-lg p-3 animate-pulse">
                            <div className="h-4 bg-muted/40 rounded w-3/4 mb-2" />
                            <div className="h-3 bg-muted/30 rounded w-1/2" />
                          </div>
                        ))}
                      </div>
                    );
                  }

                  if (stageDeals.length === 0) {
                    return (
                      <div className="flex flex-col items-center justify-center py-12 text-center">
                        <div className="w-12 h-12 rounded-full bg-muted flex items-center justify-center mb-3">
                          <Plus className="h-5 w-5 text-muted-foreground" />
                        </div>
                        <p className="text-sm text-muted-foreground font-medium">
                          {t("noDeals")}
                        </p>
                        <p className="text-xs text-muted-foreground mt-1">
                          {t("noDealsHint")}
                        </p>
                      </div>
                    );
                  }

                  return (
                    <>
                      {stageDeals.map(renderDeal)}
                      
                      {/* Auto-load on scroll */}
                      <StageScrollLoader
                        onLoadMore={() => fetchNextPageForStage(stage.id)}
                        hasMore={hasNextPageByStage[stage.id] ?? false}
                        isLoading={isFetchingNextPageByStage[stage.id] ?? false}
                      />
                      
                      {/* Fallback "Load More" button */}
                      {hasNextPageByStage[stage.id] && !isFetchingNextPageByStage[stage.id] && (
                        <Button
                          variant="ghost"
                          size="sm"
                          className="w-full text-xs"
                          onClick={() => fetchNextPageForStage(stage.id)}
                        >
                          <Plus className="h-3 w-3 mr-1" />
                          Load More
                        </Button>
                      )}
                    </>
                  );
                })()}
              </div>
            </section>
          );
        })}
      </div>

      {/* Create Deal Dialog */}
      <CreateDealDialog
        open={isCreateDialogOpen}
        onOpenChange={setIsCreateDialogOpen}
        onSuccess={handleCreateSuccess}
      />

      {/* Edit Deal Dialog */}
      {editingDeal && (
        <Dialog open={!!editingDeal} onOpenChange={closeEditDialog}>
          <DialogContent className="sm:max-w-[600px]">
            <DialogHeader>
              <DialogTitle>{t("editDialogTitle")}</DialogTitle>
            </DialogHeader>
            <DealForm
              deal={editingDeal}
              onSubmit={(data: UpdateDealFormData) => handleUpdateDeal(data)}
              onCancel={closeEditDialog}
              isLoading={isUpdating}
            />
          </DialogContent>
        </Dialog>
      )}

      {pendingStageMove && (
        <MoveStageModal
          dealId={pendingStageMove.dealId}
          currentStageId={pendingStageMove.currentStageId}
          availableStages={pipelines}
          isOpen={!!pendingStageMove}
          initialStageId={pendingStageMove.stageId}
          onClose={() => {
            setPendingStageMove(null);
            clearDraggedDeal();
          }}
          onSuccess={() => {
            setPendingStageMove(null);
            clearDraggedDeal();
          }}
        />
      )}
    </div>
  );
}
