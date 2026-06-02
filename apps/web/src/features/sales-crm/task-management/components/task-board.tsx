"use client";

import { useState } from "react";
import { Columns2, Plus, Search, RotateCcw } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { DateRangePicker } from "@/components/ui/date-range-picker";
import type { DateRange } from "react-day-picker";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { DeleteDialog } from "@/components/ui/delete-dialog";
import { Badge } from "@/components/ui/badge";
import { Card } from "@/components/ui/card";
import { useTaskList } from "../hooks/useTaskList";
import { useKanbanBoard } from "../hooks/useKanbanBoard";
import type { TaskStatus } from "../types";
import { TaskForm } from "./task-form";
import { TaskDetailModal } from "./task-detail-modal";
import { useUsers } from "@/features/master-data/user-management/hooks/useUsers";
import { useAccounts } from "@/features/sales-crm/account-management/hooks/useAccounts";
import { TaskCard } from "./task-card";
import { useHasPermission } from "@/features/master-data/user-management/hooks/useHasPermission";
import { useTranslations } from "next-intl";
import { ContactDetailModal } from "@/features/sales-crm/account-management/components/contact-detail-modal";
import type { CreateTaskFormData, UpdateTaskFormData } from "../schemas/task.schema";

export function TaskBoard() {
  const tList = useTranslations("taskManagement.list");
  const tBoard = useTranslations("taskManagement.board");
  
  // Permission checks
  const hasCreatePermission = useHasPermission("tasks.create");
  const hasEditPermission = useHasPermission("tasks.edit");
  const hasDeletePermission = useHasPermission("tasks.delete");

  const {
    setPage,
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
    pagination,
    editingTaskData,
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

  const [viewingTaskId, setViewingTaskId] = useState<string | null>(null);
  const [isDetailModalOpen, setIsDetailModalOpen] = useState(false);
  const [viewingContactId, setViewingContactId] = useState<string | null>(null);
  const [isContactModalOpen, setIsContactModalOpen] = useState(false);

  const handleViewTask = (taskId: string) => {
    setViewingTaskId(taskId);
    setIsDetailModalOpen(true);
  };

  const handleViewContact = (contactId: string) => {
    setViewingContactId(contactId);
    setIsContactModalOpen(true);
  };

  const currentRange: DateRange | undefined =
    startDueDate && endDueDate
      ? {
          from: new Date(`${startDueDate}T00:00:00`),
          to: new Date(`${endDueDate}T00:00:00`),
        }
      : startDueDate
        ? {
            from: new Date(`${startDueDate}T00:00:00`),
            to: undefined,
          }
        : undefined;

  // Use kanban board hook for drag & drop
  const {
    tasksByStatus,
    isLoading: isKanbanLoading,
    boardStatuses,
    handleDragStart,
    handleDragOver,
    handleDrop,
  } = useKanbanBoard({
    search,
    priority,
    type,
    assignedTo,
    accountId,
    startDueDate,
    endDueDate,
  });

  const statusColors: Record<TaskStatus, string> = {
    pending: "#F59E0B", // amber-500
    completed: "#10B981", // emerald-500
  };

  const statusLabels: Record<TaskStatus, string> = {
    pending: tBoard("columns.todo"),
    completed: tBoard("columns.completed"),
  };

  return (
    <div className="space-y-4">
      {/* Header with Actions */}
      <div className="flex items-center justify-between gap-3">
        <div className="flex flex-1 items-center gap-3 flex-wrap">
          <div className="relative flex-1 min-w-[220px] max-w-sm">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder={tList("searchPlaceholder")}
              value={search}
              onChange={(event) => setSearch(event.target.value)}
              className="pl-10 h-9"
            />
          </div>

          <div className="flex items-center gap-2 flex-wrap">
            <Select
              value={status || "all"}
              onValueChange={(value) => setStatus(value === "all" ? "" : value)}
            >
              <SelectTrigger className="w-[120px] h-9 text-xs">
                <SelectValue placeholder={tList("filters.statusPlaceholder")} />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">{tList("filters.statusAll")}</SelectItem>
                <SelectItem value="pending">{tList("filters.statusPending")}</SelectItem>
                <SelectItem value="completed">{tList("filters.statusCompleted")}</SelectItem>
              </SelectContent>
            </Select>

            <Select
              value={priority || "all"}
              onValueChange={(value) => setPriority(value === "all" ? "" : value)}
            >
              <SelectTrigger className="w-[110px] h-9 text-xs">
                <SelectValue placeholder={tList("filters.priorityPlaceholder")} />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">{tList("filters.priorityAll")}</SelectItem>
                <SelectItem value="low">{tList("filters.priorityLow")}</SelectItem>
                <SelectItem value="medium">{tList("filters.priorityMedium")}</SelectItem>
                <SelectItem value="high">{tList("filters.priorityHigh")}</SelectItem>
                <SelectItem value="urgent">{tList("filters.priorityUrgent")}</SelectItem>
              </SelectContent>
            </Select>

            <Select
              value={type || "all"}
              onValueChange={(value) => setType(value === "all" ? "" : value)}
            >
              <SelectTrigger className="w-[120px] h-9 text-xs">
                <SelectValue placeholder={tList("filters.typePlaceholder")} />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">{tList("filters.typeAll")}</SelectItem>
                <SelectItem value="general">{tList("filters.typeGeneral")}</SelectItem>
                <SelectItem value="call">{tList("filters.typeCall")}</SelectItem>
                <SelectItem value="email">{tList("filters.typeEmail")}</SelectItem>
                <SelectItem value="meeting">{tList("filters.typeMeeting")}</SelectItem>
                <SelectItem value="follow_up">{tList("filters.typeFollowUp")}</SelectItem>
              </SelectContent>
            </Select>

            <Select
              value={assignedTo || "all"}
              onValueChange={(value) => setAssignedTo(value === "all" ? "" : value)}
            >
              <SelectTrigger className="w-[150px] h-9 text-xs">
                <SelectValue placeholder={tList("filters.assigneePlaceholder")} />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">{tList("filters.assigneeAll")}</SelectItem>
                {users.map((user) => (
                  <SelectItem key={user.id} value={user.id}>
                    {user.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>

            <DateRangePicker
              dateRange={currentRange}
              onDateChange={(range) => {
                if (range?.from) {
                  const fromDate = new Date(range.from);
                  fromDate.setHours(0, 0, 0, 0);
                  const fromStr = `${fromDate.getFullYear()}-${String(fromDate.getMonth() + 1).padStart(2, "0")}-${String(fromDate.getDate()).padStart(2, "0")}`;
                  setStartDueDate(fromStr);

                  if (range.to) {
                    const toDate = new Date(range.to);
                    toDate.setHours(0, 0, 0, 0);
                    const toStr = `${toDate.getFullYear()}-${String(toDate.getMonth() + 1).padStart(2, "0")}-${String(toDate.getDate()).padStart(2, "0")}`;
                    setEndDueDate(toStr);
                  } else {
                    setEndDueDate("");
                  }
                } else {
                  setStartDueDate("");
                  setEndDueDate("");
                }
              }}
            />
          </div>
        </div>

        <div className="flex items-center gap-2">
          {pagination && (
            <div className="hidden md:flex items-center text-xs text-muted-foreground mr-2">
              <Columns2 className="h-3.5 w-3.5 mr-1" />
              <span>
                {pagination.total} {tList("table.itemName")}
              </span>
            </div>
          )}
          <Button
            type="button"
            variant="ghost"
            size="icon-sm"
            className="h-8 w-8 text-muted-foreground hover:text-primary"
            title="Reset filter"
            onClick={() => {
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
          >
            <RotateCcw className="h-4 w-4" />
          </Button>
          {hasCreatePermission && (
            <Button type="button" onClick={() => setIsCreateDialogOpen(true)} size="sm">
              <Plus className="h-4 w-4 mr-2" />
              {tList("buttons.addTask")}
            </Button>
          )}
        </div>
      </div>

      {/* Kanban Board - consistent with pipeline */}
      {isKanbanLoading ? (
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
                      <div key={`skeleton-item-${i}-${j}`} className="h-32 bg-muted animate-pulse rounded" />
                    ))}
                  </div>
                </div>
              </Card>
            ))}
          </div>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
          {boardStatuses.map((boardStatus) => {
            const tasksForStatus = tasksByStatus[boardStatus] ?? [];

            return (
              <div
                key={boardStatus}
                className="bg-card border border-border rounded-lg p-4 h-full flex flex-col shadow-sm hover:shadow-md transition-shadow"
                onDragOver={handleDragOver}
                onDrop={(e) => handleDrop(e, boardStatus)}
              >
                <div className="flex items-center gap-2.5 mb-4 shrink-0 pb-3 border-b border-border">
                  <div
                    className="w-3 h-3 rounded-full shrink-0 ring-2 ring-offset-2 ring-offset-background"
                    style={{
                      backgroundColor: statusColors[boardStatus],
                    }}
                  />
                  <h3 className="font-medium text-base truncate flex-1">
                    {statusLabels[boardStatus]}
                  </h3>
                  <Badge variant="secondary" className="shrink-0 text-xs font-medium">
                    {tasksForStatus.length}
                  </Badge>
                </div>

                <div className="space-y-3 min-h-[200px] flex-1 overflow-y-auto pr-1">
                  {tasksForStatus.length === 0 ? (
                    <div className="flex flex-col items-center justify-center py-12 text-center">
                      <div className="w-12 h-12 rounded-full bg-muted flex items-center justify-center mb-3">
                        <Plus className="h-5 w-5 text-muted-foreground" />
                      </div>
                      <p className="text-sm text-muted-foreground font-medium">
                        {tBoard("noTasks")}
                      </p>
                      <p className="text-xs text-muted-foreground mt-1">
                        {tBoard("noTasksHint")}
                      </p>
                    </div>
                  ) : (
                    tasksForStatus.map((task) => (
                      <div
                        key={task.id}
                        draggable
                        onDragStart={() => handleDragStart(task)}
                        className="cursor-grab active:cursor-grabbing"
                      >
                        <TaskCard
                          task={task}
                          onClick={() => handleViewTask(task.id)}
                          onClickTitle={() => handleViewTask(task.id)}
                          onEdit={hasEditPermission ? () => setEditingTaskId(task.id) : undefined}
                          onDelete={hasDeletePermission ? () => handleDeleteClick(task.id) : undefined}
                          onComplete={() => handleComplete(task.id)}
                          onClickContact={task.contact ? handleViewContact : undefined}
                        />
                      </div>
                    ))
                  )}
                </div>
              </div>
            );
          })}
        </div>
      )}

      {/* Create Task Dialog */}
      <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
        <DialogContent className="sm:max-w-[600px]">
          <DialogHeader>
            <DialogTitle>{tList("buttons.createTitle")}</DialogTitle>
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

      {/* Edit Task Dialog */}
      {editingTaskId && editingTaskData?.data && (
        <Dialog open={!!editingTaskId} onOpenChange={(open) => !open && setEditingTaskId(null)}>
          <DialogContent className="sm:max-w-[600px]">
            <DialogHeader>
              <DialogTitle>{tList("buttons.editTitle")}</DialogTitle>
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

      {/* Task Detail Modal */}
      <TaskDetailModal
        taskId={viewingTaskId}
        open={isDetailModalOpen}
        onOpenChange={(open) => {
          setIsDetailModalOpen(open);
          if (!open) {
            setViewingTaskId(null);
          }
        }}
        onTaskUpdated={() => {
          // Refresh handled by query invalidation
        }}
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

      {/* Delete Dialog */}
      <DeleteDialog
        open={!!deletingTaskId}
        onOpenChange={(open) => {
          if (!open) {
            setDeletingTaskId(null);
          }
        }}
        onConfirm={handleDeleteConfirm}
        title={tList("deleteDialog.title")}
        description={tList("deleteDialog.description")}
        itemName={tList("table.itemName")}
        isLoading={deleteTask.isPending}
      />
    </div>
  );
}

