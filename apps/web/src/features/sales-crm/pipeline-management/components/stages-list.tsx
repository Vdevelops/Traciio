"use client";

import { useState } from "react";
import { Edit, Trash2, Plus, ArrowUp, ArrowDown } from "lucide-react";
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

export function StagesList() {
  const t = useTranslations("pipelineManagement.stages");
  const { data, isLoading } = usePipelines();
  const createStage = useCreateStage();
  const updateStage = useUpdateStage();
  const deleteStage = useDeleteStage();
  const updateStagesOrder = useUpdateStagesOrder();

  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [editingStageId, setEditingStageId] = useState<string | null>(null);
  const [deletingStageId, setDeletingStageId] = useState<string | null>(null);

  const stages = data?.data || [];
  const editingStage = editingStageId ? stages.find((s) => s.id === editingStageId) : null;

  const handleCreate = async (formData: CreateStageFormData | UpdateStageFormData) => {
    try {
      // Type guard to ensure we have required fields for create
      if (!("name" in formData) || !formData.name || !("code" in formData) || !formData.code) {
        throw new Error("Name and code are required");
      }
      await createStage.mutateAsync(formData as CreateStageFormData);
      setIsCreateDialogOpen(false);
      toast.success(t("toastCreated"));
    } catch {
      // Error already handled in api-client interceptor
    }
  };

  const handleUpdate = async (formData: UpdateStageFormData) => {
    if (!editingStageId) return;
    try {
      await updateStage.mutateAsync({ id: editingStageId, data: formData });
      setEditingStageId(null);
      toast.success(t("toastUpdated"));
    } catch {
      // Error already handled
    }
  };

  const handleDeleteClick = (id: string) => {
    setDeletingStageId(id);
  };

  const handleDeleteConfirm = async () => {
    if (!deletingStageId) return;
    try {
      await deleteStage.mutateAsync(deletingStageId);
      toast.success(t("toastDeleted"));
      setDeletingStageId(null);
    } catch {
      // Error already handled
    }
  };

  const handleMoveUp = async (index: number) => {
    if (index === 0) return;
    const newStages = [...stages];
    [newStages[index - 1], newStages[index]] = [newStages[index], newStages[index - 1]];
    await handleReorder(newStages);
  };

  const handleMoveDown = async (index: number) => {
    if (index === stages.length - 1) return;
    const newStages = [...stages];
    [newStages[index], newStages[index + 1]] = [newStages[index + 1], newStages[index]];
    await handleReorder(newStages);
  };

  const handleReorder = async (reorderedStages: PipelineStage[]) => {
    try {
      await updateStagesOrder.mutateAsync(
        reorderedStages.map((stage, index) => ({
          stage_id: stage.id,
          new_order: index + 1,
        }))
      );
      toast.success(t("toastOrderUpdated"));
    } catch {
      // Error already handled
    }
  };

  return (
    <div className="space-y-4">
      {/* Header with Actions */}
      <div className="flex items-center justify-end">
        <Button onClick={() => setIsCreateDialogOpen(true)} size="sm">
          <Plus className="h-4 w-4 mr-2" />
          {t("addStage")}
        </Button>
      </div>

      {/* Table */}
      <div className="border rounded-lg">
        {isLoading ? (
          <div className="p-4 space-y-3">
            {Array.from({ length: 5 }).map((_, i) => (
              <Skeleton key={`skeleton-${i}`} className="h-10 w-full" />
            ))}
          </div>
        ) : stages.length === 0 ? (
          <div className="p-8 text-center text-muted-foreground">
            {t("noStages")}
          </div>
        ) : (
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="w-[50px]">Order</TableHead>
                <TableHead className="w-[200px]">{t("nameLabel")}</TableHead>
                <TableHead className="w-[150px]">{t("codeLabel")}</TableHead>
                <TableHead className="w-[100px]">Color</TableHead>
                <TableHead className="w-[100px]">Status</TableHead>
                <TableHead>{t("descriptionLabel")}</TableHead>
                <TableHead className="w-[150px] text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {stages.map((stage, index) => (
                <TableRow key={stage.id}>
                  <TableCell>
                    <div className="flex items-center gap-1">
                      <Button
                        variant="ghost"
                        size="icon-sm"
                        onClick={() => handleMoveUp(index)}
                        disabled={index === 0 || updateStagesOrder.isPending}
                        className="h-6 w-6"
                      >
                        <ArrowUp className="h-3 w-3" />
                      </Button>
                      <span className="text-sm font-medium w-8 text-center">
                        {stage.order}
                      </span>
                      <Button
                        variant="ghost"
                        size="icon-sm"
                        onClick={() => handleMoveDown(index)}
                        disabled={index === stages.length - 1 || updateStagesOrder.isPending}
                        className="h-6 w-6"
                      >
                        <ArrowDown className="h-3 w-3" />
                      </Button>
                    </div>
                  </TableCell>
                  <TableCell>
                    <div className="flex items-center gap-2">
                      <div
                        className="w-3 h-3 rounded-full"
                        style={{ backgroundColor: stage.color }}
                      />
                      <span className="font-medium">{stage.name}</span>
                    </div>
                  </TableCell>
                  <TableCell>
                    <code className="text-xs bg-muted px-2 py-1 rounded">
                      {stage.code}
                    </code>
                  </TableCell>
                  <TableCell>
                    <div className="flex items-center gap-2">
                      <div
                        className="w-6 h-6 rounded border"
                        style={{ backgroundColor: stage.color }}
                      />
                      <span className="text-xs text-muted-foreground">
                        {stage.color}
                      </span>
                    </div>
                  </TableCell>
                  <TableCell>
                    <div className="flex flex-col gap-1">
                      <Badge variant={stage.is_active ? "default" : "secondary"}>
                        {stage.is_active ? "Active" : "Inactive"}
                      </Badge>
                      {stage.is_won && (
                        <Badge variant="outline" className="text-green-600">
                          Won
                        </Badge>
                      )}
                      {stage.is_lost && (
                        <Badge variant="outline" className="text-red-600">
                          Lost
                        </Badge>
                      )}
                    </div>
                  </TableCell>
                  <TableCell>
                    <span className="text-sm text-muted-foreground">
                      {stage.description || "-"}
                    </span>
                  </TableCell>
                  <TableCell>
                    <div className="flex items-center justify-end gap-1">
                      <Button
                        variant="ghost"
                        size="icon-sm"
                        onClick={() => setEditingStageId(stage.id)}
                        className="h-8 w-8"
                      >
                        <Edit className="h-3.5 w-3.5" />
                      </Button>
                      <Button
                        variant="ghost"
                        size="icon-sm"
                        onClick={() => handleDeleteClick(stage.id)}
                        className="h-8 w-8 text-destructive hover:text-destructive"
                      >
                        <Trash2 className="h-3.5 w-3.5" />
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        )}
      </div>

      {/* Create Dialog */}
      <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
        <DialogContent className="sm:max-w-[500px]">
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
      {editingStageId && editingStage && (
        <Dialog open={!!editingStageId} onOpenChange={(open) => !open && setEditingStageId(null)}>
          <DialogContent className="sm:max-w-[500px]">
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
        onOpenChange={(open) => {
          if (!open) {
            setDeletingStageId(null);
          }
        }}
        onConfirm={handleDeleteConfirm}
        title={t("deleteStage")}
        description={t("deleteWarning")}
        itemName="stage"
        isLoading={deleteStage.isPending}
      />
    </div>
  );
}



