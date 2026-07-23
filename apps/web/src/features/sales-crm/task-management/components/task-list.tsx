"use client";

import { useState } from "react";
import { Plus, Search, Eye, Edit, Trash2, CheckCircle2, Calendar as CalendarIcon } from "lucide-react";
import { Avatar, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { DataTable, type Column } from "@/components/ui/data-table";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { DeleteDialog } from "@/components/ui/delete-dialog";
import { useTaskList } from "../hooks/useTaskList";
import type { Task } from "../types";
import { TaskForm } from "./task-form";
import { TaskDetailModal } from "./task-detail-modal";
import { useUsers } from "@/features/master-data/user-management/hooks/useUsers";
import { useAccounts } from "@/features/sales-crm/account-management/hooks/useAccounts";
import { useHasPermission } from "@/features/master-data/user-management/hooks/useHasPermission";
import { useTranslations } from "next-intl";
import type { CreateTaskFormData, UpdateTaskFormData } from "../schemas/task.schema";
import { parseWallClockDateTime } from "@/lib/utils";

interface TaskListProps {
  readonly onTaskClick?: (task: Task) => void;
}

export function TaskList({ onTaskClick }: TaskListProps) {
  const t = useTranslations("taskManagement.list");
  
  // Permission checks
  const hasCreatePermission = useHasPermission("tasks.create");
  const hasEditPermission = useHasPermission("tasks.edit");
  const hasDeletePermission = useHasPermission("tasks.delete");

  const {
    page,
    setPage,
    perPage,
    setPerPage,
    search,
    setSearch,
    status,
    setStatus,
    priority,
    setPriority,
    type,
    setType,
    assignedTo,
    setAssignedTo,
    accountId,
    setAccountId,
    startDueDate,
    setStartDueDate,
    endDueDate,
    setEndDueDate,
    isCreateDialogOpen,
    setIsCreateDialogOpen,
    editingTaskId,
    setEditingTaskId,
    deletingTaskId,
    setDeletingTaskId,
    tasks,
    pagination,
    editingTaskData,
    isLoading,
    handleCreate,
    handleUpdate,
    handleComplete,
    handleDeleteClick,
    handleDeleteConfirm,
    createTask,
    updateTask,
    deleteTask,
  } = useTaskList();

  const { data: usersData } = useUsers({ status: "active", per_page: 100 });
  const users = usersData?.data ?? [];

  useAccounts({ status: "active", per_page: 100 });

  const handleViewTask = (taskId: string) => {
    onTaskClick?.({ id: taskId } as Task);
  };

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

  const formatDate = (dateString: string | null) => {
    if (!dateString) return "-";
    const date = parseWallClockDateTime(dateString);
    if (!date) return "-";
    return date.toLocaleDateString("id-ID", {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
  };

  const columns: Column<Task>[] = [
    {
      id: "title",
      header: t("table.columnTask"),
      accessor: (row) => (
        <button
          onClick={() => handleViewTask(row.id)}
          className="font-medium text-primary hover:underline text-left cursor-pointer"
        >
          {row.title}
        </button>
      ),
      className: "w-[250px]",
    },
    {
      id: "contact",
      header: t("table.columnContact"),
      accessor: (row) =>
        row.contact ? (
          <button
            type="button"
            className="text-sm text-primary hover:underline cursor-pointer"
          >
            {row.contact?.name || "-"}
          </button>
        ) : (
          <span className="text-sm text-muted-foreground">-</span>
        ),
      className: "w-[160px]",
    },
    {
      id: "type",
      header: t("table.columnType"),
      accessor: (row) => (
        <Badge variant="outline" className="font-normal capitalize">
          {row.type.replace("_", " ")}
        </Badge>
      ),
      className: "w-[100px]",
    },
    {
      id: "status",
      header: t("table.columnStatus"),
      accessor: (row) => (
        <Badge variant={statusVariantMap[row.status]} className="font-normal capitalize">
          {row.status.replace("_", " ")}
        </Badge>
      ),
      className: "w-[120px]",
    },
    {
      id: "priority",
      header: t("table.columnPriority"),
      accessor: (row) => (
        <Badge variant={priorityVariantMap[row.priority]} className="font-normal capitalize">
          {row.priority}
        </Badge>
      ),
      className: "w-[100px]",
    },
    {
      id: "assigned_to",
      header: t("table.columnAssignee"),
      accessor: (row) => (
        <div className="flex items-center gap-2 text-sm text-muted-foreground">
          {row.assigned_user ? (
            <div className="flex items-center gap-2">
              {row.assigned_user.avatar_url && (
                <Avatar className="h-6 w-6 border border-border">
                  <AvatarImage src={row.assigned_user.avatar_url} alt={row.assigned_user.name} />
                </Avatar>
              )}
              <span className="truncate">{row.assigned_user.name}</span>
            </div>
          ) : (
            <span>-</span>
          )}
        </div>
      ),
      className: "w-[150px]",
    },
    {
      id: "due_date",
      header: t("table.columnDueDate"),
      accessor: (row) => (
        <div className="flex items-center gap-2 text-sm text-muted-foreground">
          {row.due_date ? (
            <>
              <CalendarIcon className="h-3.5 w-3.5" />
              <span>{formatDate(row.due_date)}</span>
            </>
          ) : (
            <span>-</span>
          )}
        </div>
      ),
      className: "w-[140px]",
    },
    {
      id: "actions",
      header: t("table.columnActions"),
      accessor: (row) => (
        <div className="flex items-center justify-end gap-1">
          <Button
            variant="ghost"
            size="icon-sm"
            className="h-8 w-8 cursor-pointer"
            title="View Details"
            onClick={() => handleViewTask(row.id)}
          >
            <Eye className="h-3.5 w-3.5" />
          </Button>
          {row.status !== "completed" && (
            <Button
              variant="ghost"
              size="icon-sm"
              className="h-8 w-8 text-primary hover:text-primary hover:bg-primary/10 cursor-pointer"
              onClick={() => handleComplete(row.id)}
              title={t("buttons.markComplete")}
            >
              <CheckCircle2 className="h-3.5 w-3.5" />
            </Button>
          )}
          {hasEditPermission && (
            <Button
              variant="ghost"
              size="icon-sm"
              onClick={() => setEditingTaskId(row.id)}
              className="h-8 w-8 cursor-pointer"
              title="Edit"
            >
              <Edit className="h-3.5 w-3.5" />
            </Button>
          )}
          {hasDeletePermission && (
            <Button
              variant="ghost"
              size="icon-sm"
              onClick={() => handleDeleteClick(row.id)}
              className="h-8 w-8 text-destructive hover:text-destructive cursor-pointer"
              title="Delete"
            >
              <Trash2 className="h-3.5 w-3.5" />
            </Button>
          )}
        </div>
      ),
      className: "w-[160px] text-right",
    },
  ];

  return (
    <div className="space-y-4">
      {/* Header with Actions */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3 flex-1 flex-wrap">
          <div className="relative flex-1 max-w-sm">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder={t("searchPlaceholder")}
              value={search}
              onChange={(event) => setSearch(event.target.value)}
              className="pl-10 h-9"
            />
          </div>

          <Select
            value={status || "all"}
            onValueChange={(value) => setStatus(value === "all" ? "" : value)}
          >
            <SelectTrigger className="w-[140px] h-9 cursor-pointer">
              <SelectValue placeholder={t("filters.statusPlaceholder")} />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all" className="cursor-pointer">{t("filters.statusAll")}</SelectItem>
              <SelectItem value="pending" className="cursor-pointer">{t("filters.statusPending")}</SelectItem>
              <SelectItem value="completed" className="cursor-pointer">{t("filters.statusCompleted")}</SelectItem>
            </SelectContent>
          </Select>

          <Select
            value={priority || "all"}
            onValueChange={(value) => setPriority(value === "all" ? "" : value)}
          >
            <SelectTrigger className="w-[140px] h-9 cursor-pointer">
              <SelectValue placeholder={t("filters.priorityPlaceholder")} />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all" className="cursor-pointer">{t("filters.priorityAll")}</SelectItem>
              <SelectItem value="low" className="cursor-pointer">{t("filters.priorityLow")}</SelectItem>
              <SelectItem value="medium" className="cursor-pointer">{t("filters.priorityMedium")}</SelectItem>
              <SelectItem value="high" className="cursor-pointer">{t("filters.priorityHigh")}</SelectItem>
              <SelectItem value="urgent" className="cursor-pointer">{t("filters.priorityUrgent")}</SelectItem>
            </SelectContent>
          </Select>

          <Select value={type || "all"} onValueChange={(value) => setType(value === "all" ? "" : value)}>
            <SelectTrigger className="w-[140px] h-9 cursor-pointer">
              <SelectValue placeholder={t("filters.typePlaceholder")} />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all" className="cursor-pointer">{t("filters.typeAll")}</SelectItem>
              <SelectItem value="general" className="cursor-pointer">{t("filters.typeGeneral")}</SelectItem>
              <SelectItem value="call" className="cursor-pointer">{t("filters.typeCall")}</SelectItem>
              <SelectItem value="email" className="cursor-pointer">{t("filters.typeEmail")}</SelectItem>
              <SelectItem value="meeting" className="cursor-pointer">{t("filters.typeMeeting")}</SelectItem>
              <SelectItem value="follow_up" className="cursor-pointer">{t("filters.typeFollowUp")}</SelectItem>
            </SelectContent>
          </Select>

          <Select
            value={assignedTo || "all"}
            onValueChange={(value) => setAssignedTo(value === "all" ? "" : value)}
          >
            <SelectTrigger className="w-40 h-9 cursor-pointer">
              <SelectValue placeholder={t("filters.assigneePlaceholder")} />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all" className="cursor-pointer">{t("filters.assigneeAll")}</SelectItem>
              {users.map((user) => (
                <SelectItem key={user.id} value={user.id} className="cursor-pointer">
                  {user.name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        {hasCreatePermission && (
          <Button type="button" onClick={() => setIsCreateDialogOpen(true)} size="sm" className="cursor-pointer">
            <Plus className="h-4 w-4 mr-2" />
            {t("buttons.addTask")}
          </Button>
        )}
      </div>

      <DataTable
        columns={columns}
        data={tasks}
        isLoading={isLoading}
        emptyMessage={t("table.empty")}
        pagination={
          pagination
            ? {
                page: pagination.page,
                per_page: pagination.per_page,
                total: pagination.total,
                total_pages: pagination.total_pages,
                has_next: pagination.has_next,
                has_prev: pagination.has_prev,
              }
            : undefined
        }
        onPageChange={setPage}
        onPerPageChange={setPerPage}
        itemName={t("table.itemName")}
        perPageOptions={[10, 20, 50, 100]}
        onResetFilters={() => {
          setSearch("");
          setStatus("");
          setPriority("");
          setType("");
          setAssignedTo("");
          setAccountId("");
          setStartDueDate("");
          setEndDueDate("");
          setPage(1);
        }}
      />

      <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
        <DialogContent className="sm:max-w-[600px]">
          <DialogHeader>
            <DialogTitle>{t("buttons.createTitle")}</DialogTitle>
          </DialogHeader>
          <TaskForm
            onSubmit={async (data) => {
              await handleCreate(data as CreateTaskFormData);
            }}
            onCancel={() => setIsCreateDialogOpen(false)}
            isLoading={createTask.isPending}
          />
        </DialogContent>
      </Dialog>

      {editingTaskId && editingTaskData?.data && (
        <Dialog open={!!editingTaskId} onOpenChange={(open) => !open && setEditingTaskId(null)}>
          <DialogContent className="sm:max-w-[600px]">
            <DialogHeader>
              <DialogTitle>{t("buttons.editTitle")}</DialogTitle>
            </DialogHeader>
            <TaskForm
              task={editingTaskData.data}
              onSubmit={async (data) => {
                await handleUpdate(data as UpdateTaskFormData);
              }}
              onCancel={() => setEditingTaskId(null)}
              isLoading={updateTask.isPending}
            />
          </DialogContent>
        </Dialog>
      )}

      <DeleteDialog
        open={!!deletingTaskId}
        onOpenChange={(open) => {
          if (!open) {
            setDeletingTaskId(null);
          }
        }}
        onConfirm={handleDeleteConfirm}
        title={t("deleteDialog.title")}
        description={t("deleteDialog.description")}
        itemName={t("table.itemName")}
        isLoading={deleteTask.isPending}
      />
    </div>
  );
}
