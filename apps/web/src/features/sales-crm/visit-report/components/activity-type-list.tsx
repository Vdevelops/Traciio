"use client";

import * as React from "react";
import { useState } from "react";
import { Edit2, Trash2, Plus, ArrowUpDown, Circle } from "lucide-react";
import { StatusSwitch } from "@/components/ui/status-switch";
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
import { Skeleton } from "@/components/ui/skeleton";
import { Badge } from "@/components/ui/badge";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { DeleteDialog } from "@/components/ui/delete-dialog";
import { toast } from "sonner";

import {
  useActivityTypes,
  useCreateActivityType,
  useUpdateActivityType,
  useDeleteActivityType,
} from "../hooks/useActivityTypes";
import { ActivityTypeForm } from "./activity-type-form";
import { renderIcon } from "../lib/icon-utils";
import type { ActivityType } from "../types/activity-type";
import type { ActivityTypeFormData } from "../schemas/activity-type.schema";

export function ActivityTypeList() {
  const t = useTranslations("visitReportActivityType.list");
  const { data: activityTypesData, isLoading } = useActivityTypes();

  const createMutation = useCreateActivityType();
  const updateMutation = useUpdateActivityType();
  const deleteMutation = useDeleteActivityType();

  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [editingActivityType, setEditingActivityType] = useState<ActivityType | null>(null);
  const [deletingId, setDeletingId] = useState<string | null>(null);

  const activityTypes = (activityTypesData?.data ?? []) as ActivityType[];

  const handleCreateSubmit = async (data: ActivityTypeFormData) => {
    try {
      await createMutation.mutateAsync(data);
      toast.success("Activity type created successfully");
      setIsCreateOpen(false);
    } catch (error: any) {
      toast.error(error?.response?.data?.error?.message || "Failed to create activity type");
    }
  };

  const handleEditSubmit = async (data: ActivityTypeFormData) => {
    if (!editingActivityType) return;
    try {
      await updateMutation.mutateAsync({
        id: editingActivityType.id,
        data,
      });
      toast.success("Activity type updated successfully");
      setEditingActivityType(null);
    } catch (error: any) {
      toast.error(error?.response?.data?.error?.message || "Failed to update activity type");
    }
  };

  const handleDeleteConfirm = async () => {
    if (!deletingId) return;
    try {
      await deleteMutation.mutateAsync(deletingId);
      toast.success("Activity type deleted successfully");
      setDeletingId(null);
    } catch (error: any) {
      toast.error(error?.response?.data?.error?.message || "Failed to delete activity type");
    }
  };

  return (
    <div className="space-y-4">
      {/* Header action */}
      <div className="flex justify-end">
        <Button onClick={() => setIsCreateOpen(true)} size="sm">
          <Plus className="h-4 w-4 mr-2" />
          {t("addType")}
        </Button>
      </div>

      {/* Table */}
      <div className="border rounded-lg overflow-x-auto bg-card">
        {isLoading ? (
          <div className="p-4 space-y-3">
            {Array.from({ length: 4 }).map((_, i) => (
              <Skeleton key={`skeleton-${i}`} className="h-10 w-full" />
            ))}
          </div>
        ) : (
          <Table className="min-w-max">
            <TableHeader>
              <TableRow>
                <TableHead className="w-[80px]">{t("order")}</TableHead>
                <TableHead className="w-[180px]">{t("name")}</TableHead>
                <TableHead className="w-[120px]">{t("code")}</TableHead>
                <TableHead className="w-[80px]">{t("icon")}</TableHead>
                <TableHead className="w-[120px]">{t("badgeColor")}</TableHead>
                <TableHead>{t("description")}</TableHead>
                <TableHead className="w-[100px]">{t("status")}</TableHead>
                <TableHead className="w-[100px] text-right">{t("actions")}</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {activityTypes.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={8} className="text-center text-muted-foreground py-8">
                    {t("empty")}
                  </TableCell>
                </TableRow>
              ) : (
                activityTypes
                  .sort((a, b) => a.order - b.order)
                  .map((type) => (
                    <TableRow key={type.id} className="hover:bg-muted/50">
                      <TableCell className="font-mono text-xs text-muted-foreground">
                        {type.order}
                      </TableCell>
                      <TableCell className="font-semibold">{type.name}</TableCell>
                      <TableCell>
                        <code className="text-xs bg-muted px-1.5 py-0.5 rounded text-primary">
                          {type.code}
                        </code>
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center justify-center bg-muted/40 w-8 h-8 rounded border">
                          {type.icon ? (
                            renderIcon(type.icon, "h-4 w-4 text-muted-foreground")
                          ) : (
                            <Circle className="h-3.5 w-3.5 text-muted-foreground" />
                          )}
                        </div>
                      </TableCell>
                      <TableCell>
                        <Badge variant={type.badge_color} className="text-[10px] font-medium">
                          {type.badge_color}
                        </Badge>
                      </TableCell>
                      <TableCell className="text-muted-foreground max-w-xs truncate">
                        {type.description || "-"}
                      </TableCell>
                      <TableCell>
                        <StatusSwitch
                          checked={type.status === "active"}
                          onCheckedChange={async (checked) => {
                            try {
                              await updateMutation.mutateAsync({
                                id: type.id,
                                data: { status: checked ? "active" : "inactive" },
                              });
                              toast.success(`Activity type marked as ${checked ? "active" : "inactive"}`);
                            } catch (err: any) {
                              toast.error(err?.response?.data?.error?.message || "Failed to toggle status");
                            }
                          }}
                        />
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center justify-end gap-1">
                          <Button
                            variant="ghost"
                            size="icon-sm"
                            onClick={() => setEditingActivityType(type)}
                            className="h-8 w-8"
                          >
                            <Edit2 className="h-3.5 w-3.5" />
                          </Button>
                          <Button
                            variant="ghost"
                            size="icon-sm"
                            onClick={() => setDeletingId(type.id)}
                            className="h-8 w-8 text-destructive hover:text-destructive"
                            disabled={type.activity_count ? type.activity_count > 0 : false}
                          >
                            <Trash2 className="h-3.5 w-3.5" />
                          </Button>
                        </div>
                      </TableCell>
                    </TableRow>
                  ))
              )}
            </TableBody>
          </Table>
        )}
      </div>

      {/* Create Dialog */}
      <Dialog open={isCreateOpen} onOpenChange={setIsCreateOpen}>
        <DialogContent className="sm:max-w-[500px]">
          <DialogHeader>
            <DialogTitle>{t("createTitle")}</DialogTitle>
          </DialogHeader>
          <ActivityTypeForm
            onSubmit={handleCreateSubmit}
            onCancel={() => setIsCreateOpen(false)}
            isLoading={createMutation.isPending}
          />
        </DialogContent>
      </Dialog>

      {/* Edit Dialog */}
      <Dialog
        open={!!editingActivityType}
        onOpenChange={(open) => !open && setEditingActivityType(null)}
      >
        <DialogContent className="sm:max-w-[500px]">
          <DialogHeader>
            <DialogTitle>{t("editTitle")}</DialogTitle>
          </DialogHeader>
          {editingActivityType && (
            <ActivityTypeForm
              activityType={editingActivityType}
              onSubmit={handleEditSubmit}
              onCancel={() => setEditingActivityType(null)}
              isLoading={updateMutation.isPending}
            />
          )}
        </DialogContent>
      </Dialog>

      {/* Delete Dialog */}
      <DeleteDialog
        open={!!deletingId}
        onOpenChange={(open) => !open && setDeletingId(null)}
        onConfirm={handleDeleteConfirm}
        title={t("deleteTitle")}
        description={t("deleteDescription")}
        itemName="activity type"
        isLoading={deleteMutation.isPending}
      />
    </div>
  );
}
