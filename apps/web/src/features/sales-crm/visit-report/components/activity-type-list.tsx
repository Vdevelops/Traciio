"use client";

import { Edit, Trash2, Plus } from "lucide-react";
import { StatusSwitch } from "@/components/ui/status-switch";
import { useTranslations } from "next-intl";
import { renderIcon } from "../lib/icon-utils";
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
import { useActivityTypeList } from "../hooks/useActivityTypeList";
import { ActivityTypeForm } from "./activity-type-form";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { DeleteDialog } from "@/components/ui/delete-dialog";
import { useHasPermission } from "@/features/master-data/user-management/hooks/useHasPermission";
import type { CreateActivityTypeFormData, UpdateActivityTypeFormData } from "../schemas/activity-type.schema";

export function ActivityTypeList() {
  const hasViewPermission = useHasPermission("visit-reports.activity-type");
  const hasCreatePermission = useHasPermission("visit-reports.activity-type");
  const hasEditPermission = useHasPermission("visit-reports.activity-type");
  const hasDeletePermission = useHasPermission("visit-reports.activity-type");

  const {
    editingType,
    setEditingType,
    deletingTypeId,
    setDeletingTypeId,
    isCreateDialogOpen,
    setIsCreateDialogOpen,
    types,
    typeForEdit,
    isLoading,
    handleCreate,
    handleUpdate,
    handleDelete,
    createType,
    updateType,
    deleteType,
  } = useActivityTypeList();
  const t = useTranslations("visitReportActivityType.list");

  if (!hasViewPermission) {
    return (
      <div className="text-center text-muted-foreground py-8">
        {t("noPermission")}
      </div>
    );
  }

  const handleDeleteClick = (id: string) => {
    setDeletingTypeId(id);
  };

  const handleDeleteConfirm = async () => {
    if (deletingTypeId) {
      await handleDelete(deletingTypeId);
      setDeletingTypeId(null);
    }
  };

  return (
    <div className="space-y-4">
      {/* Header with Actions */}
      <div className="flex items-center justify-end">
        {hasCreatePermission && (
          <Button onClick={() => setIsCreateDialogOpen(true)} size="sm">
            <Plus className="h-4 w-4 mr-2" />
            {t("addType")}
          </Button>
        )}
      </div>

      {/* Table */}
      <div className="border rounded-lg">
        {isLoading ? (
          <div className="p-4 space-y-3">
            {Array.from({ length: 5 }).map((_, i) => (
              <Skeleton key={`skeleton-${i}`} className="h-10 w-full" />
            ))}
          </div>
        ) : (
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="w-[150px]">{t("name")}</TableHead>
                <TableHead className="w-[120px]">{t("code")}</TableHead>
                <TableHead>{t("description")}</TableHead>
                <TableHead className="w-[100px]">{t("icon")}</TableHead>
                <TableHead className="w-[120px]">{t("badgeColor")}</TableHead>
                <TableHead className="w-[80px]">{t("order")}</TableHead>
                <TableHead className="w-[100px]">{t("status")}</TableHead>
                <TableHead className="w-[120px] text-right">{t("actions")}</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {types.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={8} className="text-center text-muted-foreground py-8">
                    {t("empty")}
                  </TableCell>
                </TableRow>
              ) : (
                types.map((type) => (
                  <TableRow key={type.id}>
                    <TableCell className="font-medium">{type.name}</TableCell>
                    <TableCell>
                      <code className="text-xs bg-muted px-2 py-1 rounded">{type.code}</code>
                    </TableCell>
                    <TableCell className="text-sm text-muted-foreground">
                      {type.description || "-"}
                    </TableCell>
                    <TableCell>
                      {type.icon ? (
                        <div className="flex items-center gap-2">
                          {renderIcon(type.icon, "h-4 w-4")}
                          <span className="text-xs text-muted-foreground">{type.icon}</span>
                        </div>
                      ) : (
                        <span className="text-muted-foreground">-</span>
                      )}
                    </TableCell>
                    <TableCell>
                      <Badge variant={type.badge_color as "default" | "secondary" | "destructive" | "outline"}>
                        {type.badge_color}
                      </Badge>
                    </TableCell>
                    <TableCell>{type.order}</TableCell>
                    <TableCell>
                      <StatusSwitch
                        checked={type.status === "active"}
                        onCheckedChange={(checked) => {
                          updateType.mutate({
                            id: type.id,
                            data: { status: checked ? "active" : "inactive" },
                          });
                        }}
                      />
                    </TableCell>
                    <TableCell className="text-right">
                      <div className="flex items-center justify-end gap-1">
                        {hasEditPermission && (
                          <Button
                            variant="ghost"
                            size="icon-sm"
                            className="h-8 w-8"
                            onClick={() => setEditingType(type.id)}
                            title={t("edit")}
                          >
                            <Edit className="h-3.5 w-3.5" />
                          </Button>
                        )}
                        {hasDeletePermission && (!type.activity_count || type.activity_count === 0) && (
                          <Button
                            variant="ghost"
                            size="icon-sm"
                            className="h-8 w-8 text-destructive hover:text-destructive"
                            onClick={() => handleDeleteClick(type.id)}
                            title={t("delete")}
                          >
                            <Trash2 className="h-3.5 w-3.5" />
                          </Button>
                        )}
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
      {hasCreatePermission && (
        <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
          <DialogContent className="max-w-2xl">
            <DialogHeader>
              <DialogTitle>{t("createTitle")}</DialogTitle>
            </DialogHeader>
            <ActivityTypeForm
              onSubmit={handleCreate}
              onCancel={() => setIsCreateDialogOpen(false)}
              isLoading={createType.isPending}
            />
          </DialogContent>
        </Dialog>
      )}

      {/* Edit Dialog */}
      {hasEditPermission && typeForEdit && (
        <Dialog open={!!editingType} onOpenChange={(open) => !open && setEditingType(null)}>
          <DialogContent className="max-w-2xl">
            <DialogHeader>
              <DialogTitle>{t("editTitle")}</DialogTitle>
            </DialogHeader>
            <ActivityTypeForm
              activityType={typeForEdit}
              onSubmit={handleUpdate}
              onCancel={() => setEditingType(null)}
              isLoading={updateType.isPending}
            />
          </DialogContent>
        </Dialog>
      )}

      {/* Delete Dialog */}
      {hasDeletePermission && (
        <DeleteDialog
          open={!!deletingTypeId}
          onOpenChange={(open) => !open && setDeletingTypeId(null)}
          onConfirm={handleDeleteConfirm}
          isLoading={deleteType.isPending}
          title={t("deleteTitle")}
          description={t("deleteDescription")}
        />
      )}
    </div>
  );
}

