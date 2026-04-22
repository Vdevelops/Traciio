"use client";

import { useParams, useRouter } from "next/navigation";
import { useTranslations } from "next-intl";
import { AuthGuard } from "@/features/auth/components/auth-guard";
import { PermissionGuard } from "@/features/auth/components/permission-guard";
import { PageMotion } from "@/components/motion";
import { PageDetailLayout } from "@/components/layouts/page-detail-layout";
import {
  useTask,
  useUpdateTask,
  useDeleteTask,
  useCompleteTask,
} from "@/features/sales-crm/task-management/hooks/useTasks";
import dynamic from "next/dynamic";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Pencil, Trash2, CheckCircle2, Plus } from "lucide-react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "@/components/ui/dialog";
import { Skeleton } from "@/components/ui/skeleton";
import { toast } from "sonner";
import { useState, Suspense } from "react";

const TaskForm = dynamic(
  () =>
    import(
      "@/features/sales-crm/task-management/components/task-form"
    ).then((mod) => ({ default: mod.TaskForm })),
  {
    loading: () => <Skeleton className="h-[400px] w-full" />,
    ssr: false,
  },
);

const ReminderSettings = dynamic(
  () =>
    import(
      "@/features/sales-crm/task-management/components/reminder-settings"
    ).then((mod) => ({ default: mod.ReminderSettings })),
  {
    loading: () => <Skeleton className="h-[200px] w-full" />,
    ssr: false,
  },
);

function TaskDetailPageContent() {
  const params = useParams();
  const router = useRouter();
  const taskId = params.id as string;
  const t = useTranslations("tasks.detail");
  const tCommon = useTranslations("common");

  const { data, isLoading } = useTask(taskId);
  const updateTask = useUpdateTask();
  const deleteTask = useDeleteTask();
  const completeTask = useCompleteTask();

  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);

  const task = data?.data;

  const handleUpdate = async (
    formData: Parameters<typeof updateTask.mutateAsync>[0]["data"],
  ) => {
    try {
      await updateTask.mutateAsync({ id: taskId, data: formData });
      setIsEditDialogOpen(false);
      toast.success(t("toast.updated"));
    } catch {
      // Error already handled in api-client interceptor
    }
  };

  const handleDelete = async () => {
    try {
      await deleteTask.mutateAsync(taskId);
      setIsDeleteDialogOpen(false);
      toast.success(t("toast.deleted"));
      router.push("/tasks");
    } catch {
      // Error already handled
    }
  };

  const handleCompleteTask = async () => {
    try {
      await completeTask.mutateAsync(taskId);
      toast.success(t("toast.completed"));
    } catch {
      // Error already handled
    }
  };

  const statusColor = (status: string) => {
    switch (status) {
      case "completed":
        return "default";
      case "in_progress":
        return "secondary";
      case "cancelled":
        return "destructive";
      default:
        return "outline";
    }
  };

  const priorityColor = (priority: string) => {
    switch (priority) {
      case "urgent":
        return "destructive";
      case "high":
        return "default";
      case "medium":
        return "secondary";
      default:
        return "outline";
    }
  };

  return (
    <PageMotion className="p-2 sm:p-4">
      <PageDetailLayout
        title={
          isLoading ? (
            <Skeleton className="w-64 h-8" />
          ) : task ? (
            task.title
          ) : (
            "Task Not Found"
          )
        }
        subtitle={
          task ? (
            <div className="flex flex-wrap items-center gap-2 mt-1">
              <Badge
                variant={statusColor(task.status ?? "pending")}
                className="text-xs"
              >
                {(task.status ?? "pending").replace("_", " ").toUpperCase()}
              </Badge>
              <Badge
                variant={priorityColor(task.priority ?? "medium")}
                className="text-xs"
              >
                {(task.priority ?? "medium").toUpperCase()}
              </Badge>
              <Badge variant="outline" className="text-xs">
                {(task.type ?? "general").replace("_", " ")}
              </Badge>
              {task.due_date && (
                <span className="text-xs text-muted-foreground">
                  Due: {formatSafeDate(task.due_date)}
                </span>
              )}
            </div>
          ) : undefined
        }
        backHref="/tasks"
        actions={
          task ? (
            <>
              {task.status !== "completed" &&
                task.status !== "cancelled" && (
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    onClick={handleCompleteTask}
                    disabled={completeTask.isPending}
                    className="cursor-pointer"
                  >
                    <CheckCircle2 className="h-4 w-4 mr-2" />
                    {completeTask.isPending
                      ? t("actions.markCompletedLoading")
                      : t("actions.markCompleted")}
                  </Button>
                )}
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={() => setIsEditDialogOpen(true)}
                className="cursor-pointer"
              >
                <Pencil className="h-4 w-4 mr-2" />
                {tCommon("edit")}
              </Button>
              <Button
                type="button"
                variant="destructive"
                size="sm"
                onClick={() => setIsDeleteDialogOpen(true)}
                className="cursor-pointer"
              >
                <Trash2 className="h-4 w-4 mr-2" />
                {tCommon("delete")}
              </Button>
            </>
          ) : undefined
        }
      >
        {/* Task detail content */}
        {task && (
          <div className="space-y-4">
            {/* Main Info Card */}
            <Card>
              <CardContent className="pt-6 space-y-6">
                {task.description && (
                  <div>
                    <h3 className="font-medium mb-2 text-sm text-muted-foreground">
                      {t("sections.description")}
                    </h3>
                    <p className="text-sm whitespace-pre-wrap">
                      {task.description}
                    </p>
                  </div>
                )}

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {task.account && (
                    <DetailItem
                      label={t("sections.account")}
                      value={task.account.name}
                    />
                  )}
                  {task.contact && (
                    <DetailItem
                      label={t("sections.contact")}
                      value={task.contact.name}
                    />
                  )}
                  {task.assigned_user && (
                    <DetailItem
                      label={t("sections.assignedTo")}
                      value={task.assigned_user.name}
                    />
                  )}
                  {task.deal && (
                    <DetailItem
                      label="Related Deal"
                      value={task.deal.title ?? "—"}
                    />
                  )}
                  {task.lead && (
                    <DetailItem
                      label="Related Lead"
                      value={`${task.lead.first_name ?? ""} ${task.lead.last_name ?? ""}`.trim() || "—"}
                    />
                  )}
                </div>

                {/* Schedule Info */}
                {task.is_schedule_task && (
                  <div className="pt-4 border-t">
                    <h3 className="font-medium mb-2 text-sm text-muted-foreground">
                      Schedule
                    </h3>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <DetailItem
                        label="Start"
                        value={formatSafeDate(task.scheduled_start_time, true)}
                      />
                      <DetailItem
                        label="End"
                        value={formatSafeDate(task.scheduled_end_time, true)}
                      />
                    </div>
                  </div>
                )}

                <div className="text-xs text-muted-foreground pt-4 border-t space-y-1">
                  <p>
                    {t("sections.createdAt")}{" "}
                    {formatSafeDate(task.created_at, true)}
                  </p>
                  <p>
                    {t("sections.updatedAt")}{" "}
                    {formatSafeDate(task.updated_at, true)}
                  </p>
                </div>
              </CardContent>
            </Card>

            {/* Quick Action: Create & Link Lead */}
            {!task.lead_id && (
              <Card>
                <CardHeader>
                  <CardTitle className="text-base">Quick Actions</CardTitle>
                </CardHeader>
                <CardContent>
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    className="cursor-pointer"
                    onClick={() => {
                      toast.info(
                        "Create & Link Lead feature coming soon",
                      );
                    }}
                  >
                    <Plus className="h-4 w-4 mr-2" />
                    Create & Link Lead
                  </Button>
                </CardContent>
              </Card>
            )}

            {/* Reminders */}
            <Card>
              <CardHeader>
                <CardTitle className="text-base">Reminders</CardTitle>
              </CardHeader>
              <CardContent>
                <Suspense
                  fallback={<Skeleton className="h-[200px] w-full" />}
                >
                  <ReminderSettings taskId={taskId} />
                </Suspense>
              </CardContent>
            </Card>
          </div>
        )}
      </PageDetailLayout>

      {/* Edit Dialog */}
      <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
        <DialogContent className="sm:max-w-[600px] max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>{tCommon("edit")}</DialogTitle>
          </DialogHeader>
          <Suspense fallback={<Skeleton className="h-[400px] w-full" />}>
            {task && (
              <TaskForm
                task={task}
                onSubmit={handleUpdate}
                onCancel={() => setIsEditDialogOpen(false)}
                isLoading={updateTask.isPending}
              />
            )}
          </Suspense>
        </DialogContent>
      </Dialog>

      {/* Delete Dialog */}
      <Dialog
        open={isDeleteDialogOpen}
        onOpenChange={setIsDeleteDialogOpen}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{t("deleteDialog.title")}</DialogTitle>
            <DialogDescription>
              {t("deleteDialog.description")}
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => setIsDeleteDialogOpen(false)}
              className="cursor-pointer"
            >
              {t("deleteDialog.cancel")}
            </Button>
            <Button
              type="button"
              variant="destructive"
              onClick={handleDelete}
              disabled={deleteTask.isPending}
              className="cursor-pointer"
            >
              {deleteTask.isPending
                ? t("deleteDialog.confirmLoading")
                : t("deleteDialog.confirm")}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </PageMotion>
  );
}

function DetailItem({
  label,
  value,
}: {
  readonly label: string;
  readonly value: string;
}) {
  return (
    <div>
      <p className="text-xs text-muted-foreground">{label}</p>
      <p className="text-sm">{value}</p>
    </div>
  );
}

function formatSafeDate(
  value?: string | null,
  includeTime = false,
): string {
  if (!value) return "—";
  const date = new Date(value);
  if (isNaN(date.getTime())) return "—";
  if (includeTime) {
    return date.toLocaleString("id-ID", {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  }
  return date.toLocaleDateString("id-ID", {
    year: "numeric",
    month: "long",
    day: "numeric",
  });
}

export default function TaskDetailPageClient() {
  return (
    <AuthGuard>
      <PermissionGuard requiredPermission="tasks.view">
        <TaskDetailPageContent />
      </PermissionGuard>
    </AuthGuard>
  );
}
