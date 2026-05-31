"use client";

import { useMemo, useState } from "react";
import { CalendarDays, Eye, MapPin, Building2 } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import { useVisitReports } from "../hooks/useVisitReports";
import type { VisitReport } from "../types";
import type { SalesRepGroup } from "./visit-report-team-overview";
import { useTranslations } from "next-intl";

const PER_PAGE = 10;
const DIALOG_SKELETONS = ["d1", "d2", "d3", "d4", "d5"];

const STATUS_CONFIG: Record<
  string,
  {
    variant: "default" | "secondary" | "destructive" | "outline";
  }
> = {
  draft: { variant: "outline" },
  submitted: { variant: "secondary" },
  approved: { variant: "default" },
  rejected: { variant: "destructive" },
};

function getDicebearSrc(seed: string, avatarUrl?: string): string {
  if (avatarUrl && avatarUrl.trim() !== "") return avatarUrl;
  return `https://api.dicebear.com/7.x/lorelei/svg?seed=${encodeURIComponent(seed)}`;
}

function formatDate(dateString: string): string {
  return new Date(dateString).toLocaleDateString("id-ID", {
    weekday: "short",
    year: "numeric",
    month: "short",
    day: "numeric",
  });
}

function buildPageNumbers(current: number, total: number): (number | "...")[] {
  if (total <= 7) return Array.from({ length: total }, (_, i) => i + 1);
  const pages: (number | "...")[] = [1];
  if (current > 3) pages.push("...");
  for (
    let i = Math.max(2, current - 1);
    i <= Math.min(total - 1, current + 1);
    i++
  ) {
    pages.push(i);
  }
  if (current < total - 2) pages.push("...");
  pages.push(total);
  return pages;
}

interface ReportRowProps {
  readonly report: VisitReport;
  readonly onView: (id: string) => void;
}

function ReportRow({ report, onView }: ReportRowProps) {
  const t = useTranslations("visitReportTeamOverview");
  const cfg = STATUS_CONFIG[report.status] ?? STATUS_CONFIG.draft;
  const primaryLabel = report.account?.name ?? report.deal?.title ?? null;

  return (
    <Card
      className="overflow-hidden transition-colors hover:bg-muted/30 cursor-pointer"
      onClick={() => onView(report.id)}
    >
      <CardContent className="px-3 py-2.5">
        <div className="flex items-center justify-between gap-2">
          <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
            <CalendarDays className="h-3 w-3 shrink-0" />
            <span>{formatDate(report.visit_date)}</span>
          </div>

          <div className="flex items-center gap-1">
            <Badge variant={cfg.variant} className="text-xs gap-1 py-0">
              {t(`status.${report.status}` as "status.draft")}
            </Badge>
            <span className="rounded-md border border-border/70 p-1 text-muted-foreground">
              <Eye className="h-3.5 w-3.5" />
            </span>
          </div>
        </div>

        <div className="flex items-center gap-1.5 mt-1.5">
          <Building2 className="h-3 w-3 text-muted-foreground shrink-0" />
          <span className="text-sm font-medium text-foreground truncate">
            {primaryLabel ?? (
              <span className="text-muted-foreground italic text-xs">
                {t("row.noAccount")}
              </span>
            )}
          </span>
        </div>

        <div className="flex items-center justify-between gap-2 mt-0.5">
          <p className="text-xs text-muted-foreground truncate flex-1">
            {report.purpose}
          </p>
          {report.check_in_time && (
            <div className="flex items-center gap-1 text-xs text-muted-foreground shrink-0">
              <MapPin className="h-3 w-3" />
              <span>
                {new Date(report.check_in_time).toLocaleTimeString("id-ID", {
                  hour: "2-digit",
                  minute: "2-digit",
                })}
              </span>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}

export interface UserVisitReportsDialogProps {
  readonly open: boolean;
  readonly onOpenChange: (open: boolean) => void;
  readonly rep: SalesRepGroup;
  readonly onViewReport: (id: string) => void;
  readonly isDrawerOpen?: boolean;
}

export function UserVisitReportsDialog({
  open,
  onOpenChange,
  rep,
  onViewReport,
  isDrawerOpen = false,
}: UserVisitReportsDialogProps) {
  const t = useTranslations("visitReportTeamOverview");
  const [page, setPage] = useState(1);

  const { data, isLoading } = useVisitReports({
    sales_rep_id: rep.id,
    page,
    per_page: PER_PAGE,
  });

  const reports: VisitReport[] = data?.data ?? [];
  const totalPages = data?.meta?.pagination?.total_pages ?? 1;
  const pageNumbers = useMemo(
    () => buildPageNumbers(page, totalPages),
    [page, totalPages]
  );
  const avatarSrc = getDicebearSrc(rep.email || rep.name, rep.avatarUrl);

  return (
    <Dialog
      open={open}
      onOpenChange={(isOpen) => {
        if (!isOpen && isDrawerOpen) return;
        onOpenChange(isOpen);
      }}
      modal={!isDrawerOpen}
    >
      <DialogContent
        className="sm:max-w-[640px] p-0 gap-0"
        overlayClassName={isDrawerOpen ? "pointer-events-none" : undefined}
      >
        <DialogHeader className="shrink-0">
          <DialogTitle className="flex items-center gap-2 text-base">
            <img
              src={avatarSrc}
              alt={rep.name}
              className="h-8 w-8 rounded-full bg-primary/5 object-cover shrink-0"
            />
            {rep.name}
          </DialogTitle>
        </DialogHeader>

        <div className="flex-1 overflow-y-auto">
          {isLoading && (
            <div className="space-y-2 px-2 py-2">
              {DIALOG_SKELETONS.map((id) => (
                <Skeleton key={id} className="h-[72px] rounded-lg" />
              ))}
            </div>
          )}
          {!isLoading && reports.length === 0 && (
            <div className="flex flex-col items-center justify-center py-14 text-center">
              <CalendarDays className="h-10 w-10 text-muted-foreground/30 mb-3" />
              <p className="text-sm text-muted-foreground">
                {t("dialog.noReports")}
              </p>
            </div>
          )}
          {!isLoading && reports.length > 0 && (
            <div className="space-y-1.5 py-2">
              {reports.map((report) => (
                <ReportRow
                  key={report.id}
                  report={report}
                  onView={onViewReport}
                />
              ))}
            </div>
          )}
        </div>

        {totalPages > 1 && (
          <div className="border-t p-3">
            <Pagination>
              <PaginationContent>
                <PaginationItem>
                  <PaginationPrevious
                    onClick={(e) => {
                      e.preventDefault();
                      setPage((prev) => Math.max(1, prev - 1));
                    }}
                    className={page <= 1 ? "pointer-events-none opacity-50" : "cursor-pointer"}
                  />
                </PaginationItem>

                {pageNumbers.map((pageNumber, index) => (
                  <PaginationItem key={`${pageNumber}-${index}`}>
                    {pageNumber === "..." ? (
                      <PaginationEllipsis />
                    ) : (
                      <PaginationLink
                        isActive={pageNumber === page}
                        onClick={(e) => {
                          e.preventDefault();
                          setPage(pageNumber);
                        }}
                        className="cursor-pointer"
                      >
                        {pageNumber}
                      </PaginationLink>
                    )}
                  </PaginationItem>
                ))}

                <PaginationItem>
                  <PaginationNext
                    onClick={(e) => {
                      e.preventDefault();
                      setPage((prev) => Math.min(totalPages, prev + 1));
                    }}
                    className={page >= totalPages ? "pointer-events-none opacity-50" : "cursor-pointer"}
                  />
                </PaginationItem>
              </PaginationContent>
            </Pagination>
          </div>
        )}
      </DialogContent>
    </Dialog>
  );
}
