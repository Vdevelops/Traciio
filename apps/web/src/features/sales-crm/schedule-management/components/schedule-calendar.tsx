"use client";

import { useState, useMemo } from "react";
import { Calendar, dateFnsLocalizer, type View, type Event, type ToolbarProps } from "react-big-calendar";
import {
  endOfDay,
  endOfMonth,
  endOfWeek,
  format,
  getDay,
  parse,
  startOfDay,
  startOfMonth,
  startOfWeek,
} from "date-fns";
import type { Locale } from "date-fns";
import { enUS, id as idLocale } from "date-fns/locale";
import {
  Plus,
  Calendar as CalendarIcon,
  ChevronLeft,
  ChevronRight,
  RefreshCw,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Drawer } from "@/components/ui/drawer";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { useSchedules, useDeleteSchedule, useSyncToGoogleCalendar, useUnsyncFromGoogleCalendar, useUpdateSchedule, useCreateSchedule } from "../hooks/useSchedules";
import type { Schedule } from "../types";
import { ScheduleForm } from "./schedule-form";
import { ScheduleDetailModal } from "./schedule-detail-modal";
import { useHasPermission } from "@/features/master-data/user-management/hooks/useHasPermission";
import { useLocale, useTranslations } from "next-intl";
import type { CreateScheduleFormData, UpdateScheduleFormData } from "../schemas/schedule.schema";

const locales = {
  en: enUS,
  id: idLocale,
};

const localizer = dateFnsLocalizer({
  format,
  parse,
  startOfWeek: (date: Date, options?: { locale?: Locale }) =>
    startOfWeek(date, { ...options, weekStartsOn: 1 }),
  getDay,
  locales,
});

export function ScheduleCalendar() {
  const t = useTranslations("scheduleManagement.calendar");
  const locale = useLocale();
  const calendarCulture = locale.startsWith("id") ? "id" : "en";
  const calendarLocale = calendarCulture === "id" ? idLocale : enUS;

  // Permission checks
  const hasCreatePermission = useHasPermission("schedules.create");
  const hasEditPermission = useHasPermission("schedules.edit");
  const hasDeletePermission = useHasPermission("schedules.delete");

  const [currentDate, setCurrentDate] = useState(new Date());
  const [currentView, setCurrentView] = useState<View>("month");
  const [status, setStatus] = useState<Schedule["status"] | "all">("all");
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [editingScheduleId, setEditingScheduleId] = useState<string | null>(null);
  const [viewingScheduleId, setViewingScheduleId] = useState<string | null>(null);
  const [isDetailModalOpen, setIsDetailModalOpen] = useState(false);
  const [selectedSlot, setSelectedSlot] = useState<{ start: Date; end: Date } | null>(null);
  const [isDateDrawerOpen, setIsDateDrawerOpen] = useState(false);
  const [selectedDate, setSelectedDate] = useState<Date | null>(null);

  // Calculate date range based on current view
  const dateRange = useMemo(() => {
    let viewType: "month" | "week" | "day";
    if (currentView === "month") {
      viewType = "month";
    } else if (currentView === "week") {
      viewType = "week";
    } else {
      viewType = "day";
    }
    const start =
      viewType === "month"
        ? startOfMonth(currentDate)
        : viewType === "week"
          ? startOfWeek(currentDate, { weekStartsOn: 1 })
          : startOfDay(currentDate);
    const end =
      viewType === "month"
        ? endOfMonth(currentDate)
        : viewType === "week"
          ? endOfWeek(currentDate, { weekStartsOn: 1 })
          : endOfDay(currentDate);
    return { start, end };
  }, [currentDate, currentView]);

  const scheduleStatus = status === "all" ? undefined : status;
  const { data, isLoading, refetch } = useSchedules({
    scheduled_at_from: dateRange.start.toISOString(),
    scheduled_at_to: dateRange.end.toISOString(),
    status: scheduleStatus,
  });

  const schedules = useMemo(() => {
    return data?.data ?? [];
  }, [data?.data]);

  const createSchedule = useCreateSchedule();
  const updateSchedule = useUpdateSchedule();
  const deleteSchedule = useDeleteSchedule();
  const syncToGoogleCalendar = useSyncToGoogleCalendar();
  const unsyncFromGoogleCalendar = useUnsyncFromGoogleCalendar();

  // Convert schedules to calendar events
  const events: Event[] = useMemo(() => {
    return schedules.map((schedule) => {
      const start = new Date(schedule.scheduled_at);
      // Default duration: 1 hour, or use reminder_minutes_before if available
      const duration = schedule.reminder_minutes_before ? schedule.reminder_minutes_before * 60 * 1000 : 60 * 60 * 1000;
      const end = new Date(start.getTime() + duration);

      return {
        id: schedule.id,
        title: schedule.title,
        start,
        end,
        resource: schedule,
      };
    });
  }, [schedules]);

  const handleUpdate = async (data: UpdateScheduleFormData) => {
    if (editingScheduleId) {
      await updateSchedule.mutateAsync({ id: editingScheduleId, data });
      setEditingScheduleId(null);
    }
  };

  const handleCreate = async (data: CreateScheduleFormData | UpdateScheduleFormData) => {
    // Type guard to ensure it's CreateScheduleFormData
    if ("title" in data && data.title && "scheduled_at" in data && data.scheduled_at) {
      await createSchedule.mutateAsync(data as CreateScheduleFormData);
      setIsCreateDialogOpen(false);
      setSelectedSlot(null);
      await refetch();
    }
  };

  const handleSelectSlot = (slotInfo: { start: Date; end: Date }) => {
    // Show drawer with schedules for the selected date
    setSelectedDate(slotInfo.start);
    setIsDateDrawerOpen(true);
  };

  const handleSelectEvent = (event: Event) => {
    const schedule = event.resource as Schedule;
    setViewingScheduleId(schedule.id);
    setIsDetailModalOpen(true);
  };


  const handleEditSchedule = (id: string) => {
    if (hasEditPermission) {
      setEditingScheduleId(id);
    }
  };

  const handleDeleteSchedule = async (id: string) => {
    if (hasDeletePermission) {
      await deleteSchedule.mutateAsync(id);
    }
  };

  const handleSyncToGoogleCalendar = async (id: string): Promise<void> => {
    await syncToGoogleCalendar.mutateAsync(id);
  };

  const handleUnsyncFromGoogleCalendar = async (id: string): Promise<void> => {
    await unsyncFromGoogleCalendar.mutateAsync(id);
  };

  // Event style getter - using theme colors from globals.css
  const eventStyleGetter = (event: Event) => {
    const schedule = event.resource as Schedule;
    // Get computed CSS variables from document root for accurate color values
    // This ensures colors match the current theme (light/dark mode)
    const root = document.documentElement;
    const getComputedColor = (varName: string, fallback: string) => {
      const value = getComputedStyle(root).getPropertyValue(varName).trim();
      return value || fallback;
    };

    let backgroundColor: string;
    let textColor: string;

    if (schedule.status === "completed") {
      // Completed: use muted-foreground (gray) with white text for contrast
      backgroundColor = getComputedColor("--muted-foreground", "oklch(0.45 0.02 240)");
      textColor = getComputedColor("--primary-foreground", "oklch(1.0 0 0)");
    } else if (schedule.status === "cancelled") {
      // Cancelled: use destructive color (red)
      backgroundColor = getComputedColor("--destructive", "oklch(0.5386 0.1937 26.7249)");
      textColor = getComputedColor("--destructive-foreground", "oklch(1.0000 0 0)");
    } else {
      // pending, confirmed, or default - use primary color (orange #F39200)
      backgroundColor = getComputedColor("--primary", "oklch(0.73 0.19 55)");
      textColor = getComputedColor("--primary-foreground", "oklch(1.0 0 0)");
    }

    return {
      style: {
        backgroundColor,
        color: textColor,
        borderRadius: "4px",
        border: "none",
        padding: "2px 4px",
        fontSize: "12px",
        cursor: "pointer",
        fontWeight: "500",
      },
    };
  };

  // Get schedules for selected date
  const schedulesForSelectedDate = useMemo(() => {
    if (selectedDate) {
      return schedules.filter((schedule) => {
        return isSameCalendarDay(selectedDate, new Date(schedule.scheduled_at));
      });
    }
    return [];
  }, [selectedDate, schedules]);

  // Custom toolbar component
  const CustomToolbar = (toolbar: ToolbarProps<Event, object>) => {
    return (
      <div className="mb-4 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex items-center gap-2 flex-wrap">
          <Button
            variant="outline"
            size="sm"
            onClick={() => toolbar.onNavigate("PREV")}
            className="cursor-pointer"
          >
            <ChevronLeft className="h-4 w-4" />
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={() => toolbar.onNavigate("TODAY")}
            className="cursor-pointer"
          >
            {t("today")}
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={() => toolbar.onNavigate("NEXT")}
            className="cursor-pointer"
          >
            <ChevronRight className="h-4 w-4" />
          </Button>
          <div className="ml-0 sm:ml-4 text-base sm:text-lg font-medium truncate">{toolbar.label}</div>
        </div>
        <div className="flex items-center gap-2 flex-wrap">
          <Select value={status} onValueChange={(value) => setStatus(value as Schedule["status"] | "all")}>
            <SelectTrigger className="w-full sm:w-[140px] cursor-pointer">
              <SelectValue placeholder={t("allStatuses")} />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all" className="cursor-pointer">{t("allStatuses")}</SelectItem>
              <SelectItem value="pending" className="cursor-pointer">{t("pending")}</SelectItem>
              <SelectItem value="confirmed" className="cursor-pointer">{t("confirmed")}</SelectItem>
              <SelectItem value="completed" className="cursor-pointer">{t("completed")}</SelectItem>
              <SelectItem value="cancelled" className="cursor-pointer">{t("cancelled")}</SelectItem>
            </SelectContent>
          </Select>
          <Button
            variant="outline"
            size="sm"
            onClick={() => refetch()}
            className="cursor-pointer"
            title={t("refresh")}
          >
            <RefreshCw className="h-4 w-4" />
          </Button>
          {hasCreatePermission && (
            <Button
              onClick={() => setIsCreateDialogOpen(true)}
              className="cursor-pointer w-full sm:w-auto"
            >
              <Plus className="mr-2 h-4 w-4" />
              {t("createSchedule")}
            </Button>
          )}
        </div>
      </div>
    );
  };

  if (isLoading) {
    return (
      <div className="flex h-[600px] items-center justify-center">
        <div className="text-sm text-muted-foreground">{t("loading")}</div>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Status Legend */}
      <div className="flex flex-wrap items-center gap-4 py-2">
        <div className="text-sm font-medium">{t("legend")}:</div>
        <div className="flex items-center gap-2">
          <div className="h-3 w-3 rounded" style={{ backgroundColor: "var(--primary)" }} />
          <span className="text-sm">{t("pending")}</span>
        </div>
        <div className="flex items-center gap-2">
          <div className="h-3 w-3 rounded" style={{ backgroundColor: "var(--success)" }} />
          <span className="text-sm">{t("confirmed")}</span>
        </div>
        <div className="flex items-center gap-2">
          <div className="h-3 w-3 rounded" style={{ backgroundColor: "var(--muted-foreground)" }} />
          <span className="text-sm">{t("completed")}</span>
        </div>
        <div className="flex items-center gap-2">
          <div className="h-3 w-3 rounded" style={{ backgroundColor: "var(--destructive)" }} />
          <span className="text-sm">{t("cancelled")}</span>
        </div>
        {schedules.some((s) => s.google_calendar_sync_status === "synced") && (
          <div className="ml-4 flex items-center gap-2">
            <Badge variant="outline" className="text-xs">
              <CalendarIcon className="mr-1 h-3 w-3" />
              {t("syncedToGoogle")}
            </Badge>
          </div>
        )}
      </div>

      {/* Calendar */}
      <div>
        <Calendar
          localizer={localizer}
          culture={calendarCulture}
          events={events}
          startAccessor="start"
          endAccessor="end"
          style={{ height: 600 }}
          view={currentView}
          onView={setCurrentView}
          date={currentDate}
          onNavigate={setCurrentDate}
          onSelectSlot={handleSelectSlot}
          onSelectEvent={handleSelectEvent}
          selectable={hasCreatePermission}
          eventPropGetter={eventStyleGetter}
          components={{
            toolbar: CustomToolbar,
          }}
          messages={{
            next: t("next"),
            previous: t("previous"),
            today: t("today"),
            month: t("month"),
            week: t("week"),
            day: t("day"),
            agenda: t("agenda"),
            date: t("date"),
            time: t("time"),
            event: t("event"),
            noEventsInRange: t("noEventsInRange"),
          }}
        />
      </div>

      {/* Create Schedule Dialog */}
      <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>{t("createSchedule")}</DialogTitle>
          </DialogHeader>
          <ScheduleForm
            onSubmit={handleCreate}
            onCancel={() => {
              setIsCreateDialogOpen(false);
              setSelectedSlot(null);
            }}
            defaultScheduledAt={selectedSlot?.start.toISOString()}
            isLoading={createSchedule.isPending}
          />
        </DialogContent>
      </Dialog>

      {/* Edit Schedule Dialog */}
      {editingScheduleId && (
        <Dialog open={!!editingScheduleId} onOpenChange={(open) => !open && setEditingScheduleId(null)}>
          <DialogContent className="max-w-2xl">
            <DialogHeader>
              <DialogTitle>{t("editSchedule")}</DialogTitle>
            </DialogHeader>
            <ScheduleForm
              scheduleId={editingScheduleId}
              onSubmit={handleUpdate}
              onCancel={() => setEditingScheduleId(null)}
              isLoading={updateSchedule.isPending}
            />
          </DialogContent>
        </Dialog>
      )}

      {/* Schedule Detail Modal */}
      {viewingScheduleId && (
        <ScheduleDetailModal
          scheduleId={viewingScheduleId}
          open={isDetailModalOpen}
          onOpenChange={setIsDetailModalOpen}
          onEdit={handleEditSchedule}
          onDelete={handleDeleteSchedule}
          onSyncToGoogleCalendar={handleSyncToGoogleCalendar}
          onUnsyncFromGoogleCalendar={handleUnsyncFromGoogleCalendar}
          hasEditPermission={hasEditPermission}
          hasDeletePermission={hasDeletePermission}
        />
      )}

      {/* Date Drawer - Shows schedules for selected date */}
      <Drawer
        open={isDateDrawerOpen}
        onOpenChange={setIsDateDrawerOpen}
        side="bottom"
        title={selectedDate ? format(selectedDate, "EEEE, MMMM dd, yyyy", { locale: calendarLocale }) : t("schedules")}
        description={
          schedulesForSelectedDate.length > 0
            ? t("schedulesCount", { count: schedulesForSelectedDate.length })
            : t("noSchedulesForDate")
        }
        className="max-h-[85vh]"
      >
        {schedulesForSelectedDate.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-16 text-center">
            <div className="rounded-full bg-muted p-6 mb-4">
              <CalendarIcon className="h-8 w-8 text-muted-foreground" />
            </div>
            <p className="text-sm font-medium mb-1">{t("noSchedulesForDate")}</p>
            <p className="text-xs text-muted-foreground mb-6">Click create to add a new schedule</p>
            {hasCreatePermission && (
              <Button
                onClick={() => {
                  if (selectedDate) {
                    setSelectedSlot({ start: selectedDate, end: new Date(selectedDate.getTime() + 60 * 60 * 1000) });
                    setIsDateDrawerOpen(false);
                    setIsCreateDialogOpen(true);
                  }
                }}
                className="cursor-pointer"
              >
                <Plus className="mr-2 h-4 w-4" />
                {t("createSchedule")}
              </Button>
            )}
          </div>
        ) : (
          <div className="space-y-2">
            {schedulesForSelectedDate.map((schedule) => {
              const scheduleStart = new Date(schedule.scheduled_at);
              const scheduleEnd = schedule.reminder_minutes_before
                ? new Date(scheduleStart.getTime() + schedule.reminder_minutes_before * 60 * 1000)
                : new Date(scheduleStart.getTime() + 60 * 60 * 1000);

              // Status colors using theme variables from globals.css
              let statusColor: string;
              let statusBgColor: string;
              switch (schedule.status) {
                case "pending":
                  statusColor = "bg-primary";
                  statusBgColor = "bg-accent/40";
                  break;
                case "confirmed":
                  statusColor = "bg-green-600";
                  statusBgColor = "bg-green-500/20";
                  break;
                case "completed":
                  statusColor = "bg-muted-foreground";
                  statusBgColor = "bg-muted";
                  break;
                case "cancelled":
                  statusColor = "bg-destructive";
                  statusBgColor = "bg-destructive/10";
                  break;
                default:
                  statusColor = "bg-primary";
                  statusBgColor = "bg-accent/30";
                  break;
              }

              return (
                <button
                  key={schedule.id}
                  type="button"
                  className="w-full text-left group"
                  onClick={() => {
                    setViewingScheduleId(schedule.id);
                    setIsDetailModalOpen(true);
                    setIsDateDrawerOpen(false);
                  }}
                >
                  <div className={`rounded-lg p-4 transition-all hover:shadow-sm border border-border/50 hover:border-border ${statusBgColor}`}>
                    <div className="flex items-start gap-3">
                      {/* Time indicator */}
                      <div className="flex flex-col items-center min-w-[60px] pt-1">
                        <div className={`h-2 w-2 rounded-full ${statusColor} mb-1`} />
                        <span className="text-xs font-medium text-foreground">
                          {format(scheduleStart, "HH:mm")}
                        </span>
                        <span className="text-xs text-muted-foreground">
                          {format(scheduleEnd, "HH:mm")}
                        </span>
                      </div>

                      {/* Content */}
                      <div className="flex-1 min-w-0">
                        <div className="flex items-start justify-between gap-2 mb-1">
                          <h3 className="font-medium text-sm leading-tight group-hover:text-primary transition-colors">
                            {schedule.title}
                          </h3>
                          <div className="flex items-center gap-1.5 shrink-0">
                            {schedule.google_calendar_sync_status === "synced" && (
                              <CalendarIcon className="h-3.5 w-3.5 text-muted-foreground" />
                            )}
                          </div>
                        </div>
                        {schedule.description && (
                          <p className="text-xs text-muted-foreground line-clamp-2 mb-2">
                            {schedule.description}
                          </p>
                        )}
                        <div className="flex items-center gap-3 flex-wrap">
                          <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-card border border-border">
                            {schedule.status}
                          </span>
                          {schedule.reminder_minutes_before && (
                            <span className="text-xs text-muted-foreground">
                              {t("reminder")} {schedule.reminder_minutes_before} {t("minutesBefore")}
                            </span>
                          )}
                        </div>
                      </div>
                    </div>
                  </div>
                </button>
              );
            })}
          </div>
        )}
      </Drawer>
    </div>
  );
}

function isSameCalendarDay(left: Date, right: Date): boolean {
  return (
    left.getFullYear() === right.getFullYear() &&
    left.getMonth() === right.getMonth() &&
    left.getDate() === right.getDate()
  );
}
