"use client";

import { Calendar as CalendarIcon, Clock } from "lucide-react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { useSchedule } from "../hooks/useSchedules";
import { useTranslations } from "next-intl";
import { format } from "date-fns";
import { Skeleton } from "@/components/ui/skeleton";
import { Edit, Trash2 } from "lucide-react";

interface ScheduleDetailModalProps {
  readonly scheduleId: string;
  readonly open: boolean;
  readonly onOpenChange: (open: boolean) => void;
  readonly onEdit?: (id: string) => void;
  readonly onDelete?: (id: string) => Promise<void>;
  readonly hasEditPermission?: boolean;
  readonly hasDeletePermission?: boolean;
}

export function ScheduleDetailModal({
  scheduleId,
  open,
  onOpenChange,
  onEdit,
  onDelete,
  hasEditPermission = false,
  hasDeletePermission = false,
}: ScheduleDetailModalProps) {
  const t = useTranslations("scheduleManagement.detail");
  const { data, isLoading } = useSchedule(scheduleId);

  const schedule = data?.data;

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
