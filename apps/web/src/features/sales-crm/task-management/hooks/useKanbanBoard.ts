"use client";

import { useState, useMemo } from "react";
import { toast } from "sonner";
import { useTranslations } from "next-intl";
import { useTasks, useUpdateTask } from "./useTasks";
import type { Task, TaskStatus } from "../types";
import type { TaskListParams } from "../types";

const BOARD_STATUSES: TaskStatus[] = ["pending", "completed"];

interface UseKanbanBoardParams {
  readonly search?: string;
  readonly status?: string;
  readonly priority?: string;
  readonly type?: string;
  readonly assignedTo?: string;
  readonly accountId?: string;
  readonly startDueDate?: string;
  readonly endDueDate?: string;
}

export function useKanbanBoard(params: UseKanbanBoardParams = {}) {
  const t = useTranslations("taskManagement.board");
  
  const taskParams: TaskListParams = {
    page: 1,
    per_page: 100,
    search: params.search || undefined,
    status: params.status ? (params.status as TaskStatus) : undefined,
    priority: params.priority ? (params.priority as TaskListParams["priority"]) : undefined,
    type: params.type ? (params.type as TaskListParams["type"]) : undefined,
    assigned_to: params.assignedTo || undefined,
    account_id: params.accountId || undefined,
    due_date_from: params.startDueDate || undefined,
    due_date_to: params.endDueDate || undefined,
  };

  const { data: tasksData, isLoading } = useTasks(taskParams);
  const updateTask = useUpdateTask();

  const [draggedTask, setDraggedTask] = useState<Task | null>(null);

  const tasks = tasksData?.data ?? [];

  // Group tasks by status
  const tasksByStatus = useMemo(() => {
    const grouped: Record<TaskStatus, Task[]> = {
      pending: [],
      completed: [],
    };

    tasks.forEach((task) => {
      if (BOARD_STATUSES.includes(task.status)) {
        grouped[task.status].push(task);
      }
    });

    return grouped;
  }, [tasks]);

  const handleDragStart = (task: Task) => {
    setDraggedTask(task);
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    e.dataTransfer.dropEffect = "move";
  };

  const handleDrop = async (e: React.DragEvent, targetStatus: TaskStatus) => {
    e.preventDefault();

    if (!draggedTask || draggedTask.status === targetStatus) {
      setDraggedTask(null);
      return;
    }

    try {
      await updateTask.mutateAsync({
        id: draggedTask.id,
        data: { status: targetStatus, sync_to_google_calendar: false },
      });

      const statusLabels: Record<TaskStatus, string> = {
        pending: t("statusPending"),
        completed: t("statusCompleted"),
      };

      toast.success(
        t("toastTaskMoved", {
          status: statusLabels[targetStatus],
        }),
      );
    } catch (error) {
      // Error already handled in api-client interceptor
    } finally {
      setDraggedTask(null);
    }
  };

  return {
    // Data
    tasksByStatus,
    isLoading,
    boardStatuses: BOARD_STATUSES,

    // Actions
    handleDragStart,
    handleDragOver,
    handleDrop,

    // Mutations
    isUpdating: updateTask.isPending,
  };
}

