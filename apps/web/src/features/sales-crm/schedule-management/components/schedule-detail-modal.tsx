"use client";

import { Calendar as CalendarIcon, Clock, ExternalLink, CheckCircle2, XCircle, Link as LinkIcon } from "lucide-react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { useSchedule, useSyncToGoogleCalendar, useUnsyncFromGoogleCalendar } from "../hooks/useSchedules";
import { useGoogleCalendarStatus, useConnectGoogleCalendar } from "@/features/profile/hooks/useGoogleCalendar";
import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import { format } from "date-fns";
import { Skeleton } from "@/components/ui/skeleton";
import { Edit, Trash2 } from "lucide-react";
import { getGoogleCalendarEventURL } from "../utils/googleCalendar";
import { Alert, AlertDescription } from "@/components/ui/alert";

interface ScheduleDetailModalProps {
  readonly scheduleId: string;
  readonly open: boolean;
  readonly onOpenChange: (open: boolean) => void;
  readonly onEdit?: (id: string) => void;
  readonly onDelete?: (id: string) => Promise<void>;
  readonly onSyncToGoogleCalendar?: (id: string) => Promise<void>;
  readonly onUnsyncFromGoogleCalendar?: (id: string) => Promise<void>;
  readonly hasEditPermission?: boolean;
  readonly hasDeletePermission?: boolean;
}

export function ScheduleDetailModal({
  scheduleId,
  open,
  onOpenChange,
  onEdit,
  onDelete,
  onSyncToGoogleCalendar,
  onUnsyncFromGoogleCalendar,
  hasEditPermission = false,
  hasDeletePermission = false,
}: ScheduleDetailModalProps) {
  const t = useTranslations("scheduleManagement.detail");
  const router = useRouter();
  const { data, isLoading } = useSchedule(scheduleId);
  const { data: googleCalendarStatus } = useGoogleCalendarStatus();
  const connectGoogleCalendar = useConnectGoogleCalendar();
  
  // Check if Google Calendar is connected
  const isGoogleCalendarConnected = googleCalendarStatus?.data?.connected ?? false;
  
  // Use provided handlers or fallback to hooks
  const syncToGoogleCalendarHook = useSyncToGoogleCalendar();
  const unsyncFromGoogleCalendarHook = useUnsyncFromGoogleCalendar();
  
  const syncToGoogleCalendar = onSyncToGoogleCalendar 
    ? { mutate: onSyncToGoogleCalendar, isPending: false }
    : syncToGoogleCalendarHook;
  const unsyncFromGoogleCalendar = onUnsyncFromGoogleCalendar
    ? { mutate: onUnsyncFromGoogleCalendar, isPending: false }
    : unsyncFromGoogleCalendarHook;

  const schedule = data?.data;

  const handleConnectGoogleCalendar = () => {
    connectGoogleCalendar.mutate();
  };

  const handleGoToProfile = () => {
    onOpenChange(false);
    router.push("/profile");
  };

  if (isLoading) {
  return (
      <Dialog open={open} onOpenChange={onOpenChange}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>{t("title")}</DialogTitle>
          </DialogHeader>
                <div className="space-y-4">
                  <Skeleton className="h-4 w-full" />
                  <Skeleton className="h-4 w-3/4" />
            <Skeleton className="h-4 w-1/2" />
          </div>
        </DialogContent>
      </Dialog>
    );
  }

  if (!schedule) {
    return null;
  }

  const isSynced = schedule.google_calendar_sync_status === "synced";

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl">
        <DialogHeader>
          <DialogTitle>{t("title")}</DialogTitle>
        </DialogHeader>
        <div className="space-y-6">
          {/* Title and Status */}
          <div>
            <h3 className="text-lg font-medium">{schedule.title}</h3>
            <Badge 
              variant={
                schedule.status === "completed" 
                  ? "default" 
                  : schedule.status === "confirmed" 
                  ? "success" 
                  : schedule.status === "cancelled"
                  ? "destructive"
                  : "outline"
              } 
              className="mt-2"
            >
              {t(`status.${schedule.status}`)}
            </Badge>
          </div>

          {/* Description */}
          {schedule.description && (
            <div>
              <h4 className="text-sm font-medium mb-2">{t("description")}</h4>
              <p className="text-sm text-muted-foreground">{schedule.description}</p>
          </div>
        )}

          {/* Task Info */}
          {schedule.task && (
            <div>
              <h4 className="text-sm font-medium mb-2">{t("task")}</h4>
              <div className="rounded-md bg-muted p-3">
                <p className="font-medium">{schedule.task.title}</p>
                {schedule.task.assigned_user && (
                  <p className="text-sm text-muted-foreground mt-1">
                    {t("assignedTo")}: {schedule.task.assigned_user.name}
                  </p>
                )}
                {schedule.task.due_date && (
                  <p className="text-sm text-muted-foreground">
                    {t("dueDate")}: {format(new Date(schedule.task.due_date), "MMM dd, yyyy HH:mm")}
                  </p>
                )}
              </div>
            </div>
          )}

          {/* Scheduled At */}
                <div>
            <h4 className="text-sm font-medium mb-2 flex items-center gap-2">
              <CalendarIcon className="h-4 w-4" />
              {t("scheduledAt")}
            </h4>
            <p className="text-sm">
              {format(new Date(schedule.scheduled_at), "EEEE, MMMM dd, yyyy 'at' HH:mm")}
                    </p>
                </div>

          {/* Reminder */}
          {schedule.reminder_minutes_before !== null && schedule.reminder_minutes_before !== undefined && (
                  <div>
              <h4 className="text-sm font-medium mb-2 flex items-center gap-2">
                <Clock className="h-4 w-4" />
                {t("reminder")}
              </h4>
              <p className="text-sm">
                {t("reminderMinutesBefore", { minutes: schedule.reminder_minutes_before })}
              </p>
                  </div>
                )}

          {/* Google Calendar Sync */}
          <div>
            <h4 className="text-sm font-medium mb-2">{t("googleCalendar")}</h4>
            {!isGoogleCalendarConnected ? (
              <Alert className="border-yellow-200 bg-yellow-50">
                <XCircle className="h-4 w-4 text-yellow-600" />
                <AlertDescription className="text-sm text-yellow-800">
                  <div className="space-y-2">
                    <p>{t("googleCalendarNotConnected")}</p>
                    <div className="flex gap-2">
                      <Button
                        size="sm"
                        variant="outline"
                        onClick={handleConnectGoogleCalendar}
                        disabled={connectGoogleCalendar.isPending}
                        className="cursor-pointer"
                      >
                        <LinkIcon className="h-3.5 w-3.5 mr-1.5" />
                        {connectGoogleCalendar.isPending
                          ? t("connecting")
                          : t("connect")}
                      </Button>
                    </div>
                  </div>
                </AlertDescription>
              </Alert>
            ) : isSynced ? (
              <div className="rounded-lg border border-green-200 bg-green-50 p-3">
                <div className="flex items-center justify-between gap-2">
                  <div className="flex items-center gap-2">
                    <CheckCircle2 className="h-4 w-4 text-green-600 shrink-0" />
                    <span className="text-sm font-medium text-green-700">{t("synced")}</span>
                  </div>
                  <div className="flex items-center gap-2">
                    {(schedule.google_calendar_event_link || schedule.google_calendar_event_id) && (
                      <a
                        href={schedule.google_calendar_event_link ?? getGoogleCalendarEventURL(schedule.google_calendar_event_id!)}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="inline-flex items-center gap-1 text-xs text-blue-600 hover:text-blue-700 hover:underline cursor-pointer"
                      >
                        <ExternalLink className="h-3 w-3" />
                        {t("view")}
                      </a>
                    )}
                    <Button
                      size="sm"
                      variant="ghost"
                      onClick={() => unsyncFromGoogleCalendar.mutate(schedule.id)}
                      disabled={unsyncFromGoogleCalendar.isPending}
                      className="h-7 px-2 text-xs cursor-pointer text-muted-foreground hover:text-destructive"
                    >
                      {t("unsync")}
                    </Button>
                  </div>
                </div>
                {schedule.google_calendar_synced_at && (
                  <p className="text-xs text-muted-foreground mt-1">
                    {t("syncedAt")}: {format(new Date(schedule.google_calendar_synced_at), "MMM dd, yyyy HH:mm")}
                  </p>
                )}
              </div>
            ) : (
              <div className="rounded-lg border border-gray-200 bg-gray-50 p-3">
                <div className="flex items-center justify-between gap-2">
                  <div className="flex items-center gap-2">
                    <XCircle className="h-4 w-4 text-muted-foreground shrink-0" />
                    <span className="text-sm text-muted-foreground">{t("notSynced")}</span>
                  </div>
                  <Button
                    size="sm"
                    variant="outline"
                    onClick={() => syncToGoogleCalendar.mutate(schedule.id)}
                    disabled={syncToGoogleCalendar.isPending}
                    className="h-7 px-2 text-xs cursor-pointer"
                  >
                    {t("sync")}
                  </Button>
                </div>
              </div>
            )}
          </div>

          {/* Actions */}
          <div className="flex justify-end gap-2 pt-4 border-t">
            {hasEditPermission && onEdit && (
              <Button
                variant="outline"
                      onClick={() => {
                  onEdit(schedule.id);
                  onOpenChange(false);
                      }}
                className="cursor-pointer"
              >
                <Edit className="mr-2 h-4 w-4" />
                {t("edit")}
              </Button>
            )}
            {hasDeletePermission && onDelete && (
              <Button
                variant="destructive"
                onClick={async () => {
                  await onDelete(schedule.id);
                  onOpenChange(false);
                }}
                className="cursor-pointer"
              >
                <Trash2 className="mr-2 h-4 w-4" />
                {t("delete")}
              </Button>
            )}
            <Button variant="outline" onClick={() => onOpenChange(false)} className="cursor-pointer">
              {t("close")}
            </Button>
          </div>
        </div>
          </DialogContent>
        </Dialog>
  );
}

