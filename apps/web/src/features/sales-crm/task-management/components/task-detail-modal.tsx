"use client";

import { Calendar, User, Building2, Contact, FileText, CheckCircle2, Edit, Trash2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Drawer } from "@/components/ui/drawer";
import { useTask, useCompleteTask, useDeleteTask, useUpdateTask } from "../hooks/useTasks";
import { toast } from "sonner";
import React, { useState } from "react";
import { DeleteDialog } from "@/components/ui/delete-dialog";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { TaskForm } from "./task-form";
import { useHasPermission } from "@/features/master-data/user-management/hooks/useHasPermission";
import type { Task } from "../types";
import type { UpdateTaskFormData } from "../schemas/task.schema";
import { useTranslations } from "next-intl";
import { ContactDetailModal } from "@/features/sales-crm/account-management/components/contact-detail-modal";
import { parseWallClockDateTime } from "@/lib/utils";

const statusVariantMap: Record<Task["status"], "default" | "secondary" | "outline" | "destructive"> = {
  pending: "outline",
  completed: "default",
};

const priorityVariantMap: Record<Task["priority"], "default" | "secondary" | "outline" | "destructive"> = {
  low: "outline",
  medium: "secondary",
  high: "default",
  urgent: "destructive",
};

interface TaskDetailModalProps {
  readonly taskId: string | null;
  readonly open: boolean;
  readonly onOpenChange: (open: boolean) => void;
  readonly onTaskUpdated?: () => void;
}

export function TaskDetailModal({
  taskId,
  open,
  onOpenChange,
  onTaskUpdated,
}: TaskDetailModalProps) {
  // Permission checks
  const hasEditPermission = useHasPermission("tasks.edit");
  const hasDeletePermission = useHasPermission("tasks.delete");
  const { data, isLoading, error } = useTask(taskId || "");
  const completeTask = useCompleteTask();
  const deleteTask = useDeleteTask();
  const updateTask = useUpdateTask();
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [viewingContactId, setViewingContactId] = useState<string | null>(null);
  const [isContactModalOpen, setIsContactModalOpen] = useState(false);

  const task = data?.data;
  const t = useTranslations("taskManagement.detail");

  const formatDateTime = (dateString?: string | null) => {
    if (!dateString) return "-";
    const date = new Date(dateString);
    return date.toLocaleString("id-ID", {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const formatDate = (dateString: string | null) => {
    if (!dateString) return "-";
    const date = parseWallClockDateTime(dateString);
    if (!date) return "-";
    return date.toLocaleDateString("id-ID", {
      year: "numeric",
      month: "long",
      day: "numeric",
    });
  };

  const handleComplete = async () => {
    if (!taskId) return;
    try {
      await completeTask.mutateAsync(taskId);
      toast.success(t("actions.completeSuccess"));
      onTaskUpdated?.();
    } catch {
      // Error handled by interceptor
    }
  };

  const handleUpdate = async (formData: UpdateTaskFormData) => {
    if (!taskId) return;
    try {
      await updateTask.mutateAsync({ id: taskId, data: formData });
      toast.success(t("actions.updateSuccess"));
      setIsEditDialogOpen(false);
      onTaskUpdated?.();
    } catch {
      // Error handled by interceptor
    }
  };

  const handleDelete = async () => {
    if (!taskId) return;
    try {
      await deleteTask.mutateAsync(taskId);
      toast.success(t("actions.deleteSuccess"));
      setIsDeleteDialogOpen(false);
      onOpenChange(false);
      onTaskUpdated?.();
    } catch {
      // Error handled by interceptor
    }
  };

  return (
    <>
      <Drawer
        open={open}
        onOpenChange={onOpenChange}
        title={t("drawerTitle")}
        side="right"
        className="max-w-2xl"
      >
        {isLoading && (
          <div className="space-y-6">
            <Skeleton className="h-8 w-48" />
            <Card>
              <CardHeader>
                <Skeleton className="h-6 w-32" />
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <Skeleton className="h-4 w-full" />
                  <Skeleton className="h-4 w-3/4" />
                </div>
              </CardContent>
            </Card>
          </div>
        )}

        {error && (
          <div className="text-center text-muted-foreground py-8">
            {t("loadError")}
          </div>
        )}

        {!isLoading && !error && task && (
          <div className="space-y-6">
            {/* Header */}
            <div className="flex items-center justify-between pb-4 border-b">
              <div className="flex items-center gap-3">
                <Badge variant={statusVariantMap[task.status]} className="font-normal">
                  {task.status.replace("_", " ")}
                </Badge>
                <Badge variant={priorityVariantMap[task.priority]} className="font-normal">
                  {task.priority}
                </Badge>
              </div>
              <div className="flex gap-2">
                {task.status !== "completed" && (
                  <Button
                    size="sm"
                    variant="outline"
                    onClick={handleComplete}
                    disabled={completeTask.isPending}
                    className="gap-2"
                  >
                    <CheckCircle2 className="h-4 w-4" />
                    {t("actions.markComplete")}
                  </Button>
                )}
                {hasEditPermission && (
                  <Button
                    size="sm"
                    variant="outline"
                    onClick={() => setIsEditDialogOpen(true)}
                    className="gap-2"
                  >
                    <Edit className="h-4 w-4" />
                    {t("actions.edit")}
                  </Button>
                )}
                {hasDeletePermission && (
                  <Button
                    size="sm"
                    variant="outline"
                    onClick={() => setIsDeleteDialogOpen(true)}
                    className="gap-2 text-destructive hover:text-destructive"
                  >
                    <Trash2 className="h-4 w-4" />
                    {t("actions.delete")}
                  </Button>
                )}
              </div>
            </div>

            {/* Task Information */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <FileText className="h-5 w-5" />
                  {t("sections.taskInformationTitle")}
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  <h3 className="text-lg font-medium">{task.title}</h3>
                  {task.description && (
                    <p className="text-sm text-muted-foreground mt-2 whitespace-pre-wrap">
                      {task.description}
                    </p>
                  )}
                </div>

                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <div className="text-sm text-muted-foreground mb-1">
                      {t("sections.typeLabel")}
                    </div>
                    <Badge variant="outline" className="font-normal">
                      {task.type.replace("_", " ")}
                    </Badge>
                  </div>
                  <div>
                    <div className="text-sm text-muted-foreground mb-1">
                      {t("sections.dueDateLabel")}
                    </div>
                    <div className="flex items-center gap-2">
                      <Calendar className="h-4 w-4 text-muted-foreground" />
                      <span className="text-sm font-medium">{formatDate(task.due_date)}</span>
                    </div>
                  </div>
                </div>

                {task.completed_at && (
                  <div>
                    <div className="text-sm text-muted-foreground mb-1">
                      {t("sections.completedAtLabel")}
                    </div>
                    <div className="flex items-center gap-2">
                      <CheckCircle2 className="h-4 w-4 text-primary" />
                      <span className="text-sm">{formatDateTime(task.completed_at)}</span>
                    </div>
                  </div>
                )}
              </CardContent>
            </Card>

            {/* Assignment & Relations */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <User className="h-5 w-5" />
                  {t("sections.assignmentTitle")}
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  {task.assigned_user && (
                    <div>
                      <div className="text-sm text-muted-foreground mb-1">
                        {t("sections.assignedToLabel")}
                      </div>
                      <div className="flex items-center gap-2">
                        <User className="h-4 w-4 text-muted-foreground" />
                        <span className="font-medium">{task.assigned_user?.name || t("sections.notAvailable")}</span>
                      </div>
                      {task.assigned_user?.email && (
                        <span className="text-xs text-muted-foreground">{task.assigned_user.email}</span>
                      )}
                    </div>
                  )}
                  {task.account && (
                    <div>
                      <div className="text-sm text-muted-foreground mb-1">
                        {t("sections.accountLabel")}
                      </div>
                      <div className="flex items-center gap-2">
                        <Building2 className="h-4 w-4 text-muted-foreground" />
                        <span className="font-medium">{task.account?.name || t("sections.notAvailable")}</span>
                      </div>
                    </div>
                  )}
                </div>

                {task.contact && (
                  <div>
                    <div className="text-sm text-muted-foreground mb-1">
                      {t("sections.contactLabel")}
                    </div>
                    <button
                      type="button"
                      className="flex items-center gap-2 text-left hover:text-primary"
                      onClick={() => {
                        if (task.contact?.id) {
                          setViewingContactId(task.contact.id);
                          setIsContactModalOpen(true);
                        }
                      }}
                    >
                      <Contact className="h-4 w-4 text-muted-foreground" />
                      <span className="font-medium">{task.contact?.name || t("sections.notAvailable")}</span>
                    </button>
                    {task.contact?.email && (
                      <span className="text-xs text-muted-foreground">{task.contact.email}</span>
                    )}
                  </div>
                )}
              </CardContent>
            </Card>

            {/* Metadata */}
            <Card>
              <CardHeader>
                <CardTitle>{t("sections.metadataTitle")}</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-2 gap-4 text-sm">
                  <div>
                    <div className="text-muted-foreground mb-1">
                      {t("sections.createdAtLabel")}
                    </div>
                    <div>{formatDateTime(task.created_at)}</div>
                  </div>
                  <div>
                    <div className="text-muted-foreground mb-1">
                      {t("sections.updatedAtLabel")}
                    </div>
                    <div>{formatDateTime(task.updated_at)}</div>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        )}
      </Drawer>

      {/* Edit Dialog */}
      {task && (
        <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
          <DialogContent className="sm:max-w-[600px]">
            <DialogHeader>
              <DialogTitle>{t("actions.editTitle")}</DialogTitle>
            </DialogHeader>
            <TaskForm
              task={task}
              onSubmit={handleUpdate}
              onCancel={() => setIsEditDialogOpen(false)}
              isLoading={updateTask.isPending}
            />
          </DialogContent>
        </Dialog>
      )}

      {/* Delete Dialog */}
      <DeleteDialog
        open={isDeleteDialogOpen}
        onOpenChange={setIsDeleteDialogOpen}
        onConfirm={handleDelete}
        title={t("deleteDialog.title")}
        description={t("deleteDialog.description")}
        itemName={t("deleteDialog.itemName")}
        isLoading={deleteTask.isPending}
      />

      {/* Contact Detail Modal */}
      <ContactDetailModal
        contactId={viewingContactId}
        open={isContactModalOpen}
        onOpenChange={(open) => {
          setIsContactModalOpen(open);
          if (!open) {
            setViewingContactId(null);
          }
        }}
      />
    </>
  );
}
