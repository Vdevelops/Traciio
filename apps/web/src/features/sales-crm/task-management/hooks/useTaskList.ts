"use client";

import { useState, useEffect } from "react";
import { toast } from "sonner";
import { useDebounce } from "@/hooks/use-debounce";
import {
  useTasks,
  useTask,
  useCreateTask,
  useUpdateTask,
  useDeleteTask,
  useAssignTask,
  useCompleteTask,
  useReminders,
  useCreateReminder,
  useUpdateReminder,
  useDeleteReminder,
} from "./useTasks";
import type {
  CreateTaskFormData,
  UpdateTaskFormData,
  AssignTaskFormData,
} from "../schemas/task.schema";
import type {
  CreateReminderFormData,
  UpdateReminderFormData,
} from "../schemas/reminder.schema";
import type { TaskListParams, ReminderListParams } from "../types";

export function useTaskList() {
  const [page, setPage] = useState(1);
  const [perPage, setPerPage] = useState(20);
  const [search, setSearch] = useState("");
  const debouncedSearch = useDebounce(search, 500);

  useEffect(() => {
    setPage(1);
  }, [debouncedSearch]);
  const [status, setStatus] = useState<string>("");
  const [priority, setPriority] = useState<string>("");
  const [type, setType] = useState<string>("");
  const [assignedTo, setAssignedTo] = useState<string>("");
  const [accountId, setAccountId] = useState<string>("");
  const [startDueDate, setStartDueDate] = useState<string>("");
  const [endDueDate, setEndDueDate] = useState<string>("");

  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [editingTaskId, setEditingTaskId] = useState<string | null>(null);
  const [deletingTaskId, setDeletingTaskId] = useState<string | null>(null);

  const taskParams: TaskListParams = {
    page,
    per_page: perPage,
    search: debouncedSearch,
    status: status ? (status as TaskListParams["status"]) : undefined,
    priority: priority ? (priority as TaskListParams["priority"]) : undefined,
    type: type ? (type as TaskListParams["type"]) : undefined,
    assigned_to: assignedTo || undefined,
    account_id: accountId || undefined,
    due_date_from: startDueDate || undefined,
    due_date_to: endDueDate || undefined,
  };

  const { data, isLoading } = useTasks(taskParams);
  const { data: editingTaskData } = useTask(editingTaskId || "");

  const createTask = useCreateTask();
  const updateTask = useUpdateTask();
  const deleteTask = useDeleteTask();
  const assignTask = useAssignTask();
  const completeTask = useCompleteTask();

  const tasks = data?.data ?? [];
  const pagination = data?.meta?.pagination;

  const handleCreate = async (formData: CreateTaskFormData) => {
    try {
      await createTask.mutateAsync(formData);
      setIsCreateDialogOpen(false);
      toast.success("Task created successfully");
    } catch {
      // Error already handled in api-client interceptor
    }
  };

  const handleUpdate = async (formData: UpdateTaskFormData) => {
    if (!editingTaskId) return;
    try {
      await updateTask.mutateAsync({ id: editingTaskId, data: formData });
      setEditingTaskId(null);
      toast.success("Task updated successfully");
    } catch {
      // Error already handled
    }
  };

  const handleAssign = async (taskId: string, data: AssignTaskFormData) => {
    try {
      await assignTask.mutateAsync({ id: taskId, data });
      toast.success("Task assigned successfully");
    } catch {
      // Error already handled
    }
  };

  const handleComplete = async (taskId: string) => {
    try {
      await completeTask.mutateAsync(taskId);
      toast.success("Task marked as completed");
    } catch {
      // Error already handled
    }
  };

  const handleUpdateStatus = async (taskId: string, newStatus: "in_progress" | "cancelled") => {
    try {
      await updateTask.mutateAsync({ id: taskId, data: { status: newStatus, sync_to_google_calendar: false } });
      const statusMessage = newStatus === "in_progress" ? "Task marked as in progress" : "Task cancelled";
      toast.success(statusMessage);
    } catch {
      // Error already handled
    }
  };

  const handleDeleteClick = (id: string) => {
    setDeletingTaskId(id);
  };

  const handleDeleteConfirm = async () => {
    if (!deletingTaskId) return;
    try {
      await deleteTask.mutateAsync(deletingTaskId);
      toast.success("Task deleted successfully");
      setDeletingTaskId(null);
    } catch {
      // Error already handled
    }
  };

  const handlePerPageChange = (newPerPage: number) => {
    setPerPage(newPerPage);
    setPage(1);
  };

  return {
    // State
    page,
    setPage,
    perPage,
    setPerPage: handlePerPageChange,
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
    // Data
    tasks,
    pagination,
    editingTaskData,
    isLoading,
    // Actions
    handleCreate,
    handleUpdate,
    handleAssign,
    handleComplete,
    handleUpdateStatus,
    handleDeleteClick,
    handleDeleteConfirm,
    // Mutations
    createTask,
    updateTask,
    deleteTask,
    assignTask,
    completeTask,
  };
}

export function useTaskReminders(taskId: string | null) {
  const reminderParams: ReminderListParams | undefined = taskId
    ? {
        task_id: taskId,
        page: 1,
        per_page: 20,
      }
    : undefined;

  const { data, isLoading } = useReminders(reminderParams);
  const reminders = data?.data ?? [];

  const createReminder = useCreateReminder();
  const updateReminder = useUpdateReminder();
  const deleteReminder = useDeleteReminder();

  const handleCreateReminder = async (formData: CreateReminderFormData) => {
    try {
      await createReminder.mutateAsync(formData);
      toast.success("Reminder created successfully");
    } catch {
      // Error already handled
    }
  };

  const handleUpdateReminder = async (id: string, formData: UpdateReminderFormData) => {
    try {
      await updateReminder.mutateAsync({ id, data: formData });
      toast.success("Reminder updated successfully");
    } catch {
      // Error already handled
    }
  };

  const handleDeleteReminder = async (id: string) => {
    try {
      await deleteReminder.mutateAsync(id);
      toast.success("Reminder deleted successfully");
    } catch {
      // Error already handled
    }
  };

  return {
    reminders,
    isLoading,
    handleCreateReminder,
    handleUpdateReminder,
    handleDeleteReminder,
    createReminder,
    updateReminder,
    deleteReminder,
  };
}


