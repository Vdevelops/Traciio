"use client";

import { useState, useMemo, useCallback } from "react";
import { Calendar, momentLocalizer, type View, type Event, type ToolbarProps } from "react-big-calendar";
import moment from "moment";
import {
  Calendar as CalendarIcon,
  ChevronLeft,
  ChevronRight,
  RefreshCw,
  MapPin,
  Clock,
  Building2,
  Search,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Drawer } from "@/components/ui/drawer";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Input } from "@/components/ui/input";
import { Skeleton } from "@/components/ui/skeleton";
import { useVisitReports } from "../hooks/useVisitReports";
import type { VisitReport } from "../types";
import { VisitReportDetailModal } from "./visit-report-detail-modal";
import { useAccounts } from "../../account-management/hooks/useAccounts";
import { useTranslations } from "next-intl";
import { useDebounce } from "@/hooks/use-debounce";

const localizer = momentLocalizer(moment);

const statusColors: Record<string, { bg: string; text: string; badge: "default" | "secondary" | "destructive" | "outline" }> = {
  draft: { bg: "var(--muted)", text: "var(--muted-foreground)", badge: "outline" },
  submitted: { bg: "var(--secondary)", text: "var(--secondary-foreground)", badge: "secondary" },
  approved: { bg: "var(--primary)", text: "var(--primary-foreground)", badge: "default" },
  rejected: { bg: "var(--destructive)", text: "var(--destructive-foreground)", badge: "destructive" },
};

interface VisitReportCalendarProps {
  readonly accountId?: string;
  readonly dealId?: string;
}

export function VisitReportCalendar({ accountId: propAccountId, dealId: propDealId }: VisitReportCalendarProps) {
  const t = useTranslations("visitReportCalendar");
  const tList = useTranslations("visitReportList");

  const [currentDate, setCurrentDate] = useState(new Date());
  const [currentView, setCurrentView] = useState<View>("month");
  const [status, setStatus] = useState<VisitReport["status"] | "all">("all");
  const [accountId, setAccountId] = useState<string>(propAccountId ?? "");
  const [search, setSearch] = useState("");
  const debouncedSearch = useDebounce(search, 300);
  const [viewingVisitReportId, setViewingVisitReportId] = useState<string | null>(null);
  const [isDetailModalOpen, setIsDetailModalOpen] = useState(false);
  const [isDateDrawerOpen, setIsDateDrawerOpen] = useState(false);
  const [selectedDate, setSelectedDate] = useState<Date | null>(null);

  // Calculate date range based on current view - fetch wider range for performance
  const dateRange = useMemo(() => {
    // Always fetch 3 months of data for smoother navigation
    const start = moment(currentDate).subtract(1, "month").startOf("month").toDate();
    const end = moment(currentDate).add(1, "month").endOf("month").toDate();
    return { start, end };
  }, [currentDate]);

  // Accounts for filter dropdown - memoized to prevent unnecessary re-renders
  const { data: accountsData } = useAccounts({ per_page: 100 });
  const accounts = useMemo(() => accountsData?.data || [], [accountsData?.data]);

  // Fetch visit reports with optimized query
  const visitReportStatus = status === "all" ? undefined : status;
  const { data, isLoading, refetch, isFetching } = useVisitReports({
    per_page: 100, // Fetch batch
    search: debouncedSearch || undefined,
    status: visitReportStatus,
    account_id: propAccountId || (accountId || undefined),
    deal_id: propDealId,
    start_date: dateRange.start.toISOString().split("T")[0],
    end_date: dateRange.end.toISOString().split("T")[0],
  });

  const visitReports = useMemo(() => {
    return data?.data ?? [];
  }, [data?.data]);

  // Convert visit reports to calendar events with memoization for performance
  const events: Event[] = useMemo(() => {
    return visitReports.map((report) => {
      const visitDate = new Date(report.visit_date);
      // Set proper start/end times based on check-in/check-out or default to full day
      let start: Date;
      let end: Date;

      if (report.check_in_time) {
        start = new Date(report.check_in_time);
        end = report.check_out_time 
          ? new Date(report.check_out_time) 
          : new Date(start.getTime() + 60 * 60 * 1000); // Default 1 hour
      } else {
        // All day event
        start = moment(visitDate).startOf("day").toDate();
        end = moment(visitDate).endOf("day").toDate();
      }

      return {
        id: report.id,
        title: report.account?.name || report.purpose,
        start,
        end,
        allDay: !report.check_in_time,
        resource: report,
      };
    });
  }, [visitReports]);

  const handleSelectSlot = useCallback((slotInfo: { start: Date; end: Date }) => {
    setSelectedDate(slotInfo.start);
    setIsDateDrawerOpen(true);
  }, []);

  const handleSelectEvent = useCallback((event: Event) => {
    const report = event.resource as VisitReport;
    setViewingVisitReportId(report.id);
    setIsDetailModalOpen(true);
  }, []);

  // Event style getter - using theme colors
  const eventStyleGetter = useCallback((event: Event) => {
    const report = event.resource as VisitReport;
    const root = document.documentElement;
    const getComputedColor = (varName: string, fallback: string) => {
      const value = getComputedStyle(root).getPropertyValue(varName).trim();
      return value || fallback;
    };

    let backgroundColor: string;
    let textColor: string;

    switch (report.status) {
      case "approved":
        backgroundColor = getComputedColor("--primary", "oklch(0.73 0.19 55)");
        textColor = getComputedColor("--primary-foreground", "oklch(1.0 0 0)");
        break;
      case "rejected":
        backgroundColor = getComputedColor("--destructive", "oklch(0.5386 0.1937 26.7249)");
        textColor = getComputedColor("--destructive-foreground", "oklch(1.0 0 0)");
        break;
      case "submitted":
        backgroundColor = getComputedColor("--secondary", "oklch(0.9 0.02 240)");
        textColor = getComputedColor("--secondary-foreground", "oklch(0.2 0 0)");
        break;
      default:
        backgroundColor = getComputedColor("--muted", "oklch(0.9 0.02 240)");
        textColor = getComputedColor("--muted-foreground", "oklch(0.45 0.02 240)");
    }

    return {
      style: {
        backgroundColor,
        color: textColor,
        borderRadius: "4px",
        border: "none",
        padding: "2px 6px",
        fontSize: "11px",
        cursor: "pointer",
        fontWeight: "500",
        overflow: "hidden",
        textOverflow: "ellipsis",
        whiteSpace: "nowrap" as const,
      },
    };
  }, []);

  // Get visit reports for selected date - memoized for performance
  const visitReportsForSelectedDate = useMemo(() => {
    if (selectedDate) {
      const selectedDateStr = moment(selectedDate).format("YYYY-MM-DD");
      return visitReports.filter((report) => {
        const reportDateStr = moment(report.visit_date).format("YYYY-MM-DD");
        return reportDateStr === selectedDateStr;
      });
    }
    return [];
  }, [selectedDate, visitReports]);

  // Custom toolbar component
  const CustomToolbar = useCallback((toolbar: ToolbarProps<Event, object>) => {
    return (
      <div className="mb-4 flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
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
          {/* Search */}
          <div className="relative flex-1 min-w-[180px] max-w-[220px]">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder={t("searchPlaceholder")}
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="pl-9 h-9"
            />
          </div>
          {/* Status filter */}
          <Select value={status} onValueChange={(value) => setStatus(value as VisitReport["status"] | "all")}>
            <SelectTrigger className="w-[130px] cursor-pointer h-9">
              <SelectValue placeholder={t("allStatuses")} />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all" className="cursor-pointer">{t("allStatuses")}</SelectItem>
              <SelectItem value="draft" className="cursor-pointer">{tList("filters.draft")}</SelectItem>
              <SelectItem value="submitted" className="cursor-pointer">{tList("filters.submitted")}</SelectItem>
              <SelectItem value="approved" className="cursor-pointer">{tList("filters.approved")}</SelectItem>
              <SelectItem value="rejected" className="cursor-pointer">{tList("filters.rejected")}</SelectItem>
            </SelectContent>
          </Select>
          {/* Account filter */}
          {!propAccountId && (
            <Select value={accountId || "all"} onValueChange={(value) => setAccountId(value === "all" ? "" : value)}>
              <SelectTrigger className="w-40 cursor-pointer h-9">
                <SelectValue placeholder={tList("filters.allAccounts")} />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all" className="cursor-pointer">{tList("filters.allAccounts")}</SelectItem>
                {accounts.map((account) => (
                  <SelectItem key={account.id} value={account.id} className="cursor-pointer">
                    {account.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          )}
          {/* Refresh */}
          <Button
            variant="outline"
            size="sm"
            onClick={() => refetch()}
            className="cursor-pointer h-9"
            title={t("refresh")}
            disabled={isFetching}
          >
            <RefreshCw className={`h-4 w-4 ${isFetching ? "animate-spin" : ""}`} />
          </Button>
        </div>
      </div>
    );
  }, [t, tList, search, status, accountId, propAccountId, accounts, refetch, isFetching]);

  if (isLoading) {
    return (
      <div className="space-y-4">
        <div className="flex justify-between items-center">
          <Skeleton className="h-8 w-48" />
          <div className="flex gap-2">
            <Skeleton className="h-9 w-32" />
            <Skeleton className="h-9 w-24" />
          </div>
        </div>
        <Skeleton className="h-[600px] w-full rounded-lg" />
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Status Legend */}
      <div className="flex flex-wrap items-center gap-4 py-2">
        <div className="text-sm font-medium">{t("legend")}:</div>
        <div className="flex items-center gap-2">
          <div className="h-3 w-3 rounded" style={{ backgroundColor: "var(--muted)" }} />
          <span className="text-sm">{t("draft")}</span>
        </div>
        <div className="flex items-center gap-2">
          <div className="h-3 w-3 rounded" style={{ backgroundColor: "var(--secondary)" }} />
          <span className="text-sm">{t("submitted")}</span>
        </div>
        <div className="flex items-center gap-2">
          <div className="h-3 w-3 rounded" style={{ backgroundColor: "var(--primary)" }} />
          <span className="text-sm">{t("approved")}</span>
        </div>
        <div className="flex items-center gap-2">
          <div className="h-3 w-3 rounded" style={{ backgroundColor: "var(--destructive)" }} />
          <span className="text-sm">{t("rejected")}</span>
        </div>
        <div className="ml-auto text-xs text-muted-foreground">
          {t("totalVisits", { count: visitReports.length })}
        </div>
      </div>

      {/* Calendar */}
      <div>
        <Calendar
          localizer={localizer}
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
          selectable
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
            showMore: (total: number) => t("showMore", { total }),
          }}
          popup
          popupOffset={{ x: 30, y: 20 }}
        />
      </div>

      {/* Visit Report Detail Modal */}
      {viewingVisitReportId && (
        <VisitReportDetailModal
          visitReportId={viewingVisitReportId}
          open={isDetailModalOpen}
          onOpenChange={(open) => {
            setIsDetailModalOpen(open);
            if (!open) {
              setViewingVisitReportId(null);
            }
          }}
          onVisitReportUpdated={() => {
            refetch();
          }}
        />
      )}

      {/* Date Drawer - Shows visit reports for selected date */}
      <Drawer
        open={isDateDrawerOpen}
        onOpenChange={setIsDateDrawerOpen}
        side="bottom"
        title={selectedDate ? moment(selectedDate).format("dddd, MMMM DD, YYYY") : t("visitReports")}
        description={
          visitReportsForSelectedDate.length > 0
            ? t("visitReportsCount", { count: visitReportsForSelectedDate.length })
            : t("noVisitReportsForDate")
        }
        className="max-h-[85vh]"
      >
        {visitReportsForSelectedDate.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-16 text-center">
            <div className="rounded-full bg-muted p-6 mb-4">
              <CalendarIcon className="h-8 w-8 text-muted-foreground" />
            </div>
            <p className="text-sm font-medium mb-1">{t("noVisitReportsForDate")}</p>
            <p className="text-xs text-muted-foreground mb-6">
              Visit pada tanggal ini akan muncul otomatis setelah sales membuat log visit dari lead atau deal.
            </p>
          </div>
        ) : (
          <div className="space-y-2 max-h-[60vh] overflow-y-auto">
            {visitReportsForSelectedDate.map((report) => {
              const statusStyle = statusColors[report.status] || statusColors.draft;

              return (
                <button
                  key={report.id}
                  type="button"
                  className="w-full text-left group"
                  onClick={() => {
                    setViewingVisitReportId(report.id);
                    setIsDetailModalOpen(true);
                    setIsDateDrawerOpen(false);
                  }}
                >
                  <div className="rounded-lg p-4 transition-all hover:shadow-sm border border-border/50 hover:border-border bg-card">
                    <div className="flex items-start gap-3">
                      {/* Time indicator */}
                      <div className="flex flex-col items-center min-w-[60px] pt-1">
                        {report.check_in_time ? (
                          <>
                            <Clock className="h-4 w-4 text-muted-foreground mb-1" />
                            <span className="text-xs font-medium text-foreground">
                              {moment(report.check_in_time).format("HH:mm")}
                            </span>
                            {report.check_out_time && (
                              <span className="text-xs text-muted-foreground">
                                {moment(report.check_out_time).format("HH:mm")}
                              </span>
                            )}
                          </>
                        ) : (
                          <>
                            <CalendarIcon className="h-4 w-4 text-muted-foreground mb-1" />
                            <span className="text-xs text-muted-foreground">{t("allDay")}</span>
                          </>
                        )}
                      </div>

                      {/* Content */}
                      <div className="flex-1 min-w-0">
                        <div className="flex items-start justify-between gap-2 mb-1">
                          <div className="flex items-center gap-2">
                            <Building2 className="h-4 w-4 text-muted-foreground shrink-0" />
                            <h3 className="font-medium text-sm leading-tight group-hover:text-primary transition-colors truncate">
                              {report.account?.name || t("noAccount")}
                            </h3>
                          </div>
                          <Badge variant={statusStyle.badge} className="shrink-0">
                            {report.status}
                          </Badge>
                        </div>
                        <p className="text-xs text-muted-foreground line-clamp-2 mb-2">
                          {report.purpose}
                        </p>
                        <div className="flex items-center gap-3 flex-wrap">
                          {report.check_in_location && (
                            <span className="inline-flex items-center gap-1 text-xs text-muted-foreground">
                              <MapPin className="h-3 w-3" />
                              {report.check_in_location.address || t("hasLocation")}
                            </span>
                          )}
                          {report.sales_rep && (
                            <span className="text-xs text-muted-foreground">
                              {report.sales_rep.name}
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
