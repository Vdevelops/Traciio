"use client";

import { useState } from "react";
import { Edit, Trash2, Plus, ArrowUp, ArrowDown, GripVertical } from "lucide-react";
import { useTranslations } from "next-intl";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { DeleteDialog } from "@/components/ui/delete-dialog";
import { usePipelines } from "../hooks/usePipelines";
import { useCreateStage, useUpdateStage, useDeleteStage, useUpdateStagesOrder } from "../hooks/useStages";
import { StageForm } from "./stage-form";
import { toast } from "sonner";
import type { PipelineStage } from "../types";
import type { CreateStageFormData, UpdateStageFormData } from "../schemas/pipeline.schema";

export function StagesManagement() {
  const t = useTranslations("pipelineManagement.stages");
  const { data, isLoading, refetch } = usePipelines();
  const createStage = useCreateStage();
  const updateStage = useUpdateStage();
  const deleteStage = useDeleteStage();
  const updateStagesOrder = useUpdateStagesOrder();

  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [editingStageId, setEditingStageId] = useState<string | null>(null);
  const [deletingStageId, setDeletingStageId] = useState<string | null>(null);

  const stages = data?.data || [];
  const sortedStages = [...stages].sort((a, b) => a.order - b.order);
  const editingStage = editingStageId ? sortedStages.find((s) => s.id === editingStageId) : null;

  const handleCreate = async (formData: CreateStageFormData | UpdateStageFormData) => {
    try {
      if (!("name" in formData) || !formData.name || !("code" in formData) || !formData.code) {
        throw new Error("Name and code are required");
      }
      await createStage.mutateAsync(formData as CreateStageFormData);
      toast.success(t("toastCreated"));
      setIsCreateDialogOpen(false);
      refetch();
    } catch {
      // Error handled by interceptor
    }
  };

  const handleUpdate = async (formData: CreateStageFormData | UpdateStageFormData) => {
    if (!editingStageId) return;
    try {
      await updateStage.mutateAsync({
        id: editingStageId,
        data: formData as UpdateStageFormData,
      });
      toast.success(t("toastUpdated"));
      setEditingStageId(null);
      refetch();
    } catch {
      // Error handled by interceptor
    }
  };

  const handleDelete = async () => {
    if (!deletingStageId) return;
    try {
      await deleteStage.mutateAsync(deletingStageId);
      toast.success(t("toastDeleted"));
      setDeletingStageId(null);
      refetch();
    } catch {
      // Error handled by interceptor
    }
  };

  const handleMoveUp = async (stage: PipelineStage) => {
    const currentIndex = sortedStages.findIndex((s) => s.id === stage.id);
    if (currentIndex <= 0) return;

    const prevStage = sortedStages[currentIndex - 1];
    const newOrder = prevStage.order;
    const prevNewOrder = stage.order;

    try {
      await updateStagesOrder.mutateAsync([
        { stage_id: stage.id, new_order: newOrder },
        { stage_id: prevStage.id, new_order: prevNewOrder },
      ]);
      toast.success(t("toastOrderUpdated"));
      refetch();
    } catch {
      // Error handled by interceptor
    }
  };

  const handleMoveDown = async (stage: PipelineStage) => {
    const currentIndex = sortedStages.findIndex((s) => s.id === stage.id);
    if (currentIndex >= sortedStages.length - 1) return;

    const nextStage = sortedStages[currentIndex + 1];
    const newOrder = nextStage.order;
    const nextNewOrder = stage.order;

    try {
      await updateStagesOrder.mutateAsync([
        { stage_id: stage.id, new_order: newOrder },
        { stage_id: nextStage.id, new_order: nextNewOrder },
      ]);
      toast.success(t("toastOrderUpdated"));
      refetch();
    } catch {
      // Error handled by interceptor
    }
  };

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div>
          <Skeleton className="h-8 w-64 mb-2" />
          <Skeleton className="h-4 w-96" />
        </div>
        <div className="border border-border/50 rounded-xl p-6 bg-card/30">
          <div className="space-y-4">
            {Array.from({ length: 5 }).map((_, i) => (
              <Skeleton key={`skeleton-${i}`} className="h-16 w-full" />
            ))}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between pb-2">
        <div>
          <h3 className="text-xl font-semibold tracking-tight text-foreground/80">{t("title")}</h3>
          <p className="text-sm text-muted-foreground mt-1">{t("description")}</p>
        </div>
        <Button onClick={() => setIsCreateDialogOpen(true)} size="sm">
          <Plus className="h-4 w-4 mr-2" />
          {t("addStage")}
        </Button>
      </div>

      {sortedStages.length === 0 ? (
        <div className="py-12 text-center border border-dashed rounded-xl bg-muted/20">
          <p className="text-sm text-muted-foreground">{t("noStages")}</p>
        </div>
      ) : (
        <div className="border border-border/50 rounded-xl bg-card/30 overflow-hidden">
          <Table>
            <TableHeader className="bg-muted/50">
              <TableRow>
                <TableHead className="w-[50px]"></TableHead>
                <TableHead>{t("nameLabel")}</TableHead>
                <TableHead>{t("codeLabel")}</TableHead>
                <TableHead>{t("orderLabel")}</TableHead>
                <TableHead>{t("colorLabel")}</TableHead>
                <TableHead>{t("isActiveLabel")}</TableHead>
                <TableHead>{t("isWonLabel")}</TableHead>
                <TableHead>{t("isLostLabel")}</TableHead>
                <TableHead className="text-right">{t("actions") || "Actions"}</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {sortedStages.map((stage, index) => (
                <TableRow key={stage.id} className="hover:bg-muted/30 transition-colors">
                  <TableCell>
                    <div className="flex items-center gap-1">
                      <GripVertical className="h-4 w-4 text-muted-foreground cursor-grab" />
                      <Button
                        variant="ghost"
                        size="icon"
                        className="h-6 w-6"
                        onClick={() => handleMoveUp(stage)}
                        disabled={index === 0}
                      >
                        <ArrowUp className="h-3 w-3" />
                      </Button>
                      <Button
                        variant="ghost"
                        size="icon"
                        className="h-6 w-6"
                        onClick={() => handleMoveDown(stage)}
                        disabled={index === sortedStages.length - 1}
                      >
                        <ArrowDown className="h-3 w-3" />
                      </Button>
                    </div>
                  </TableCell>
                  <TableCell className="font-medium">{stage.name}</TableCell>
                  <TableCell>
                    <Badge variant="outline" className="font-mono text-[10px] tracking-wider uppercase bg-background/50">
                      {stage.code}
                    </Badge>
                  </TableCell>
                  <TableCell className="text-sm text-muted-foreground">{stage.order}</TableCell>
                  <TableCell>
                    <div className="flex items-center gap-2">
                      <div
                        className="w-3 h-3 rounded-full border border-border/50"
                        style={{ backgroundColor: stage.color }}
                      />
                      <span className="text-xs text-muted-foreground font-mono">{stage.color}</span>
                    </div>
                  </TableCell>
                  <TableCell>
                    <Badge variant={stage.is_active ? "active" : "secondary"} className="text-[10px]">
                      {stage.is_active ? t("active") || "Active" : t("inactive") || "Inactive"}
                    </Badge>
                  </TableCell>
                  <TableCell>
                    {stage.is_won ? (
                      <Badge variant="active" className="text-[10px]">{t("yes") || "Yes"}</Badge>
                    ) : (
                      <span className="text-xs text-muted-foreground">{t("no") || "No"}</span>
                    )}
                  </TableCell>
                  <TableCell>
                    {stage.is_lost ? (
                      <Badge variant="destructive" className="text-[10px]">{t("yes") || "Yes"}</Badge>
                    ) : (
                      <span className="text-xs text-muted-foreground">{t("no") || "No"}</span>
                    )}
                  </TableCell>
                  <TableCell className="text-right">
                    <div className="flex items-center justify-end gap-1">
                      <Button
                        variant="ghost"
                        size="icon"
                        className="h-8 w-8"
                        onClick={() => setEditingStageId(stage.id)}
                      >
                        <Edit className="h-4 w-4" />
                      </Button>
                      <Button
                        variant="ghost"
                        size="icon"
                        className="h-8 w-8 hover:text-destructive"
                        onClick={() => setDeletingStageId(stage.id)}
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      )}

      {/* Create Dialog */}
      <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
        <DialogContent className="sm:max-w-[600px]">
          <DialogHeader>
            <DialogTitle>{t("addStage")}</DialogTitle>
          </DialogHeader>
          <StageForm
            onSubmit={handleCreate}
            onCancel={() => setIsCreateDialogOpen(false)}
            isLoading={createStage.isPending}
          />
        </DialogContent>
      </Dialog>

      {/* Edit Dialog */}
      {editingStage && (
        <Dialog open={!!editingStage} onOpenChange={() => setEditingStageId(null)}>
          <DialogContent className="sm:max-w-[600px]">
            <DialogHeader>
              <DialogTitle>{t("editStage")}</DialogTitle>
            </DialogHeader>
            <StageForm
              stage={editingStage}
              onSubmit={handleUpdate}
              onCancel={() => setEditingStageId(null)}
              isLoading={updateStage.isPending}
            />
          </DialogContent>
        </Dialog>
      )}

      {/* Delete Dialog */}
      <DeleteDialog
        open={!!deletingStageId}
        onOpenChange={(open) => !open && setDeletingStageId(null)}
        onConfirm={handleDelete}
        title={t("deleteStage")}
        description={t("deleteConfirm")}
        isLoading={deleteStage.isPending}
      />
    </div>
  );
}

