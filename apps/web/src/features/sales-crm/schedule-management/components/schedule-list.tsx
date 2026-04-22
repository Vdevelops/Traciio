"use client";

import { useState } from "react";
import {
  Plus,
  Search,
  Eye,
  Edit,
  Trash2,
  Calendar as CalendarIcon,
  CheckCircle2,
  XCircle,
  ExternalLink,
  RefreshCw,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { DataTable, type Column } from "@/components/ui/data-table";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { DeleteDialog } from "@/components/ui/delete-dialog";
import { useSchedules, useDeleteSchedule, useSyncToGoogleCalendar, useUnsyncFromGoogleCalendar } from "../hooks/useSchedules";
import type { Schedule } from "../types";
import { ScheduleForm } from "./schedule-form";
import { ScheduleDetailModal } from "./schedule-detail-modal";
import { useHasPermission } from "@/features/master-data/user-management/hooks/useHasPermission";
import { useTranslations } from "next-intl";
import type { CreateScheduleFormData, UpdateScheduleFormData } from "../schemas/schedule.schema";
import { useCreateSchedule, useUpdateSchedule } from "../hooks/useSchedules";
import { getGoogleCalendarEventURL } from "../utils/googleCalendar";

export function ScheduleList() {
  const t = useTranslations("scheduleManagement.list");
  
  // Permission checks
  const hasCreatePermission = useHasPermission("schedules.create");
  const hasEditPermission = useHasPermission("schedules.edit");
  const hasDeletePermission = useHasPermission("schedules.delete");
  
  const [page, setPage] = useState(1);
  const [perPage, setPerPage] = useState(10);
  const [search, setSearch] = useState("");
  const [status, setStatus] = useState<Schedule["status"] | "all">("all");
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [editingScheduleId, setEditingScheduleId] = useState<string | null>(null);
  const [deletingScheduleId, setDeletingScheduleId] = useState<string | null>(null);
  const [viewingScheduleId, setViewingScheduleId] = useState<string | null>(null);
  const [isDetailModalOpen, setIsDetailModalOpen] = useState(false);

  const { data, isLoading } = useSchedules({
    page,
    per_page: perPage,
    search: search || undefined,
    status: status !== "all" ? status : undefined,
  });

  const schedules = data?.data ?? [];
  const pagination = data?.meta?.pagination;

  const createSchedule = useCreateSchedule();
  const updateSchedule = useUpdateSchedule();
  const deleteSchedule = useDeleteSchedule();
  const syncToGoogleCalendar = useSyncToGoogleCalendar();
  const unsyncFromGoogleCalendar = useUnsyncFromGoogleCalendar();

  const handleCreate = async (data: CreateScheduleFormData | UpdateScheduleFormData) => {
    // Type guard to ensure it's CreateScheduleFormData
    if ("title" in data && data.title) {
      await createSchedule.mutateAsync(data as CreateScheduleFormData);
      setIsCreateDialogOpen(false);
    }
  };

  const handleUpdate = async (data: UpdateScheduleFormData) => {
    if (!editingScheduleId) return;
    await updateSchedule.mutateAsync({ id: editingScheduleId, data });
    setEditingScheduleId(null);
  };

  const handleDeleteClick = (id: string) => {
    setDeletingScheduleId(id);
  };

  const handleDeleteConfirm = async () => {
    if (!deletingScheduleId) return;
    await deleteSchedule.mutateAsync(deletingScheduleId);
    setDeletingScheduleId(null);
  };

  const handleViewSchedule = (id: string) => {
    setViewingScheduleId(id);
    setIsDetailModalOpen(true);
  };

  const handleSyncToGoogleCalendar = async (id: string): Promise<void> => {
    await syncToGoogleCalendar.mutateAsync(id);
  };

  const handleUnsyncFromGoogleCalendar = async (id: string): Promise<void> => {
    await unsyncFromGoogleCalendar.mutateAsync(id);
  };

  const statusVariantMap: Record<Schedule["status"], "default" | "secondary" | "outline" | "destructive" | "success"> = {
    pending: "outline",
    confirmed: "success",
    completed: "default",
    cancelled: "destructive",
  };

  const formatDate = (dateString: string | null) => {
    if (!dateString) return "-";
    const date = new Date(dateString);
    if (Number.isNaN(date.getTime())) return "-";
    return date.toLocaleString("id-ID", {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const editingSchedule = editingScheduleId
    ? schedules.find((s) => s.id === editingScheduleId)
    : null;

  const columns: Column<Schedule>[] = [
    {
      id: "title",
      header: t("columns.title"),
      accessor: (row) => (
        <div className="flex flex-col">
          <span className="font-medium">{row.title}</span>
          {row.task && (
            <span className="text-xs text-muted-foreground">
              {t("task")}: {row.task.title}
            </span>
          )}
        </div>
      ),
    },
    {
      id: "scheduled_at",
      header: t("columns.scheduledAt"),
      accessor: (row) => (
        <div className="flex items-center gap-2">
          <CalendarIcon className="h-4 w-4 text-muted-foreground" />
          <span>{formatDate(row.scheduled_at)}</span>
        </div>
      ),
    },
    {
      id: "status",
      header: t("columns.status"),
      accessor: (row) => (
        <Badge variant={statusVariantMap[row.status]}>
          {t(`status.${row.status}`)}
        </Badge>
      ),
    },
    {
      id: "google_calendar",
      header: t("columns.googleCalendar"),
      accessor: (row) => {
        const isSynced = row.google_calendar_sync_status === "synced";
        return (
          <div className="flex items-center gap-2">
            {isSynced ? (
              <>
                <CheckCircle2 className="h-4 w-4 text-green-600 flex-shrink-0" />
                <span className="text-sm text-green-600">{t("synced")}</span>
                {row.google_calendar_event_id && (
                  <a
                    href={getGoogleCalendarEventURL(row.google_calendar_event_id)}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="cursor-pointer text-blue-600 hover:text-blue-700 hover:underline flex items-center gap-1"
                    title={t("viewInCalendar")}
                  >
                    <ExternalLink className="h-3.5 w-3.5" />
                  </a>
                )}
              </>
            ) : (
              <>
                <XCircle className="h-4 w-4 text-muted-foreground flex-shrink-0" />
                <span className="text-sm text-muted-foreground">{t("notSynced")}</span>
              </>
            )}
          </div>
        );
      },
    },
    {
      id: "actions",
      header: t("columns.actions"),
      accessor: (row) => (
        <div className="flex items-center gap-2">
          <Button
            variant="ghost"
            size="sm"
            onClick={() => handleViewSchedule(row.id)}
            className="cursor-pointer"
          >
            <Eye className="h-4 w-4" />
          </Button>
          {hasEditPermission && (
            <Button
              variant="ghost"
              size="sm"
              onClick={() => setEditingScheduleId(row.id)}
              className="cursor-pointer"
            >
              <Edit className="h-4 w-4" />
            </Button>
          )}
          {row.google_calendar_sync_status === "synced" ? (
            <Button
              variant="ghost"
              size="sm"
              onClick={() => handleUnsyncFromGoogleCalendar(row.id)}
              disabled={unsyncFromGoogleCalendar.isPending}
              className="cursor-pointer"
            >
              <RefreshCw className="h-4 w-4" />
            </Button>
          ) : (
            <Button
              variant="ghost"
              size="sm"
              onClick={() => handleSyncToGoogleCalendar(row.id)}
              disabled={syncToGoogleCalendar.isPending}
              className="cursor-pointer"
            >
              <CalendarIcon className="h-4 w-4" />
            </Button>
          )}
          {hasDeletePermission && (
            <Button
              variant="ghost"
              size="sm"
              onClick={() => handleDeleteClick(row.id)}
              className="cursor-pointer text-destructive hover:text-destructive"
            >
              <Trash2 className="h-4 w-4" />
            </Button>
          )}
        </div>
      ),
    },
  ];

  return (
    <div className="space-y-4">
      {/* Filters */}
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex flex-col gap-3 sm:flex-row sm:flex-1 sm:items-center">
          <div className="relative flex-1 min-w-0">
            <Search className="absolute left-2 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground pointer-events-none" />
            <Input
              placeholder={t("searchPlaceholder")}
              value={search}
              onChange={(e) => {
                setSearch(e.target.value);
                setPage(1);
              }}
              className="pl-8 w-full"
            />
          </div>
          <Select value={status} onValueChange={(value) => setStatus(value as Schedule["status"] | "all")}>
            <SelectTrigger className="w-full sm:w-[180px] cursor-pointer">
              <SelectValue placeholder={t("filterStatus")} />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all" className="cursor-pointer">{t("allStatuses")}</SelectItem>
              <SelectItem value="pending" className="cursor-pointer">{t("status.pending")}</SelectItem>
              <SelectItem value="confirmed" className="cursor-pointer">{t("status.confirmed")}</SelectItem>
              <SelectItem value="completed" className="cursor-pointer">{t("status.completed")}</SelectItem>
              <SelectItem value="cancelled" className="cursor-pointer">{t("status.cancelled")}</SelectItem>
            </SelectContent>
          </Select>
        </div>
        {hasCreatePermission && (
          <Button 
            onClick={() => setIsCreateDialogOpen(true)} 
            className="cursor-pointer w-full sm:w-auto"
          >
            <Plus className="mr-2 h-4 w-4" />
            <span className="hidden sm:inline">{t("createSchedule")}</span>
            <span className="sm:hidden">{t("createSchedule")}</span>
          </Button>
        )}
      </div>

      {/* Table */}
      <DataTable
        columns={columns}
        data={schedules}
        isLoading={isLoading}
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
      />

      {/* Create Dialog */}
      <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>{t("createSchedule")}</DialogTitle>
          </DialogHeader>
        <ScheduleForm
            onSubmit={handleCreate}
          onCancel={() => setIsCreateDialogOpen(false)}
          isLoading={createSchedule.isPending}
        />
        </DialogContent>
      </Dialog>

      {/* Edit Dialog */}
      {editingSchedule && (
        <Dialog open={!!editingScheduleId} onOpenChange={(open) => !open && setEditingScheduleId(null)}>
          <DialogContent className="max-w-2xl">
            <DialogHeader>
              <DialogTitle>{t("editSchedule")}</DialogTitle>
            </DialogHeader>
            <ScheduleForm
              schedule={editingSchedule}
              onSubmit={handleUpdate}
              onCancel={() => setEditingScheduleId(null)}
              isLoading={updateSchedule.isPending}
            />
          </DialogContent>
        </Dialog>
      )}

      {/* Detail Modal */}
      {viewingScheduleId && (
        <ScheduleDetailModal
          scheduleId={viewingScheduleId}
          open={isDetailModalOpen}
          onOpenChange={setIsDetailModalOpen}
        />
      )}

      {/* Delete Dialog */}
      <DeleteDialog
        open={!!deletingScheduleId}
        onOpenChange={(open) => !open && setDeletingScheduleId(null)}
        onConfirm={handleDeleteConfirm}
        isLoading={deleteSchedule.isPending}
        title={t("deleteSchedule")}
        description={t("deleteScheduleDescription")}
        />
    </div>
  );
}

