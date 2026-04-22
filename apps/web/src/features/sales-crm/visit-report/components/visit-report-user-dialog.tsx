/* eslint-disable @next/next/no-img-element */
"use client";

import { useState, useMemo } from "react";
import {
  CheckCircle2,
  Clock,
  FileText,
  XCircle,
  CalendarDays,
  Send,
  Edit,
  Trash2,
  Eye,
  MoreVertical,
  Bell,
  MapPin,
  Building2,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { DeleteDialog } from "@/components/ui/delete-dialog";
import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import {
  useVisitReports,
  useApproveVisitReport,
  useRejectVisitReport,
  useDeleteVisitReport,
  useUpdateVisitReport,
} from "../hooks/useVisitReports";
import { useHasPermission } from "@/features/master-data/user-management/hooks/useHasPermission";
import type { VisitReport } from "../types";
import type { SalesRepGroup } from "./visit-report-team-overview";
import { useTranslations } from "next-intl";

// ─── Constants ────────────────────────────────────────────────────────────────

const PER_PAGE = 10;
const DIALOG_SKELETONS = ["d1", "d2", "d3", "d4", "d5"];

const STATUS_CONFIG: Record<
  string,
  {
    variant: "default" | "secondary" | "destructive" | "outline";
    icon: typeof CheckCircle2;
  }
> = {
  draft: { variant: "outline", icon: FileText },
  submitted: { variant: "secondary", icon: Clock },
  approved: { variant: "default", icon: CheckCircle2 },
  rejected: { variant: "destructive", icon: XCircle },
};

// ─── Helpers ──────────────────────────────────────────────────────────────────

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

// ─── ReportRow ────────────────────────────────────────────────────────────────

interface ReportRowProps {
  readonly report: VisitReport;
  readonly onView: (id: string) => void;
  readonly onEdit?: (id: string) => void;
  readonly onDelete?: (id: string) => void;
  readonly onApprove?: (id: string) => void;
  readonly onReject?: (id: string) => void;
  readonly onSubmit: (id: string) => void;
  readonly isApproving: boolean;
  readonly isRejecting: boolean;
  readonly rejectingId: string | null;
  readonly isUpdating: boolean;
}

function ReportRow({
  report,
  onView,
  onEdit,
  onDelete,
  onApprove,
  onReject,
  onSubmit,
  isApproving,
  isRejecting,
  rejectingId,
  isUpdating,
}: ReportRowProps) {
  const t = useTranslations("visitReportTeamOverview");
  const cfg = STATUS_CONFIG[report.status] ?? STATUS_CONFIG.draft;
  const StatusIcon = cfg.icon;

  const canApprove = !!onApprove && report.status === "submitted";
  const canReject = !!onReject && report.status === "submitted";
  const canEdit =
    !!onEdit && (report.status === "draft" || report.status === "submitted");
  const canSubmit = report.status === "draft";
  const hasActions =
    canApprove || canReject || canEdit || canSubmit || !!onDelete;

  const primaryLabel = report.account?.name ?? report.deal?.title ?? null;

  return (
    <Card
      className="overflow-hidden transition-colors hover:bg-muted/30 cursor-pointer"
      onClick={() => onView(report.id)}
    >
      <CardContent className="px-3 py-2.5">
        {/* Row 1: date + status + actions */}
        <div className="flex items-center justify-between gap-2">
          <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
            <CalendarDays className="h-3 w-3 shrink-0" />
            <span>{formatDate(report.visit_date)}</span>
          </div>

          <div className="flex items-center gap-1">
            <Badge variant={cfg.variant} className="text-xs gap-1 py-0">
              <StatusIcon className="h-2.5 w-2.5" />
              {t(`status.${report.status}` as "status.draft")}
            </Badge>

            {hasActions && (
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-7 w-7"
                    onClick={(e) => e.stopPropagation()}
                  >
                    <MoreVertical className="h-3.5 w-3.5" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" className="w-44 z-200">
                  <DropdownMenuItem
                    onClick={(e) => {
                      e.stopPropagation();
                      onView(report.id);
                    }}
                  >
                    <Eye className="h-4 w-4 mr-2" />
                    {t("actions.view")}
                  </DropdownMenuItem>

                  {(canApprove || canReject) && <DropdownMenuSeparator />}

                  {canApprove && (
                    <DropdownMenuItem
                      onClick={(e) => {
                        e.stopPropagation();
                        onApprove?.(report.id);
                      }}
                      disabled={isApproving}
                      className="hover:bg-transparent"
                      style={{ color: "var(--color-success)" }}
                    >
                      <CheckCircle2 className="h-4 w-4 mr-2" />
                      {t("actions.approve")}
                    </DropdownMenuItem>
                  )}

                  {canReject && (
                    <DropdownMenuItem
                      onClick={(e) => {
                        e.stopPropagation();
                        onReject?.(report.id);
                      }}
                      disabled={isRejecting && rejectingId === report.id}
                      variant="destructive"
                    >
                      <XCircle className="h-4 w-4 mr-2" />
                      {t("actions.reject")}
                    </DropdownMenuItem>
                  )}

                  {(canEdit || canSubmit || !!onDelete) &&
                    (canApprove || canReject) && <DropdownMenuSeparator />}

                  {canSubmit && (
                    <DropdownMenuItem
                      onClick={(e) => {
                        e.stopPropagation();
                        onSubmit(report.id);
                      }}
                      disabled={isUpdating}
                      className="text-primary"
                    >
                      <Send className="h-4 w-4 mr-2" />
                      {t("actions.submit")}
                    </DropdownMenuItem>
                  )}

                  {canEdit && (
                    <DropdownMenuItem
                      onClick={(e) => {
                        e.stopPropagation();
                        onEdit?.(report.id);
                      }}
                    >
                      <Edit className="h-4 w-4 mr-2" />
                      {t("actions.edit")}
                    </DropdownMenuItem>
                  )}

                  {onDelete && (
                    <DropdownMenuItem
                      onClick={(e) => {
                        e.stopPropagation();
                        onDelete(report.id);
                      }}
                      variant="destructive"
                    >
                      <Trash2 className="h-4 w-4 mr-2" />
                      {t("actions.delete")}
                    </DropdownMenuItem>
                  )}
                </DropdownMenuContent>
              </DropdownMenu>
            )}
          </div>
        </div>

        {/* Row 2: account / deal name */}
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

        {/* Row 3: purpose + check-in time */}
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

        {/* Row 4: inline quick-action buttons */}
        {(canApprove || canReject || canSubmit) && (
          <div className="flex items-center gap-1.5 mt-2 pt-2 border-t">
            {canApprove && (
              <Button
                size="sm"
                variant="outline"
                className="h-6 text-xs px-2 rounded-md"
                style={{ backgroundColor: "var(--color-success)", color: "var(--color-success-foreground)" }}
                disabled={isApproving}
                onClick={(e) => {
                  e.stopPropagation();
                  onApprove?.(report.id);
                }}
              >
                <CheckCircle2 className="h-3 w-3 mr-1" />
                {isApproving ? t("actions.approving") : t("actions.approve")}
              </Button>
            )}
            {canReject && (
              <Button
                size="sm"
                variant="outline"
                className="h-6 text-xs px-2 text-destructive border-destructive/30 hover:bg-destructive/10"
                disabled={isRejecting && rejectingId === report.id}
                onClick={(e) => {
                  e.stopPropagation();
                  onReject?.(report.id);
                }}
              >
                <XCircle className="h-3 w-3 mr-1" />
                {t("actions.reject")}
              </Button>
            )}
            {canSubmit && (
              <Button
                size="sm"
                variant="outline"
                className="h-6 text-xs px-2 text-primary border-primary/30 hover:bg-primary/10"
                disabled={isUpdating}
                onClick={(e) => {
                  e.stopPropagation();
                  onSubmit(report.id);
                }}
              >
                <Send className="h-3 w-3 mr-1" />
                {isUpdating
                  ? t("actions.submitting")
                  : t("actions.submitForApproval")}
              </Button>
            )}
          </div>
        )}
      </CardContent>
    </Card>
  );
}

// ─── UserVisitReportsDialog ───────────────────────────────────────────────────

export interface UserVisitReportsDialogProps {
  readonly open: boolean;
  readonly onOpenChange: (open: boolean) => void;
  readonly rep: SalesRepGroup;
  /** Called when the user wants to view a specific report — opens the detail drawer at the parent level. */
  readonly onViewReport: (id: string) => void;
  /** When true, a top-level drawer is open and the dialog should not trap focus. */
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
  const [deletingId, setDeletingId] = useState<string | null>(null);
  const [rejectingId, setRejectingId] = useState<string | null>(null);
  const [rejectReason, setRejectReason] = useState("");

  const hasEditPermission = useHasPermission("visit-reports.edit");
  const hasDeletePermission = useHasPermission("visit-reports.delete");
  const hasApprovePermission = useHasPermission("visit-reports.approve");
  const hasRejectPermission = useHasPermission("visit-reports.reject");

  const approveVisitReport = useApproveVisitReport();
  const rejectVisitReport = useRejectVisitReport();
  const deleteVisitReport = useDeleteVisitReport();
  const updateVisitReport = useUpdateVisitReport();

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

  const handleApprove = (id: string) => approveVisitReport.mutate(id);

  const handleSubmit = (id: string) =>
    updateVisitReport.mutate({ id, data: { status: "submitted" } });

  const handleDeleteConfirm = () => {
    if (!deletingId) return;
    deleteVisitReport.mutate(deletingId, {
      onSuccess: () => setDeletingId(null),
    });
  };

  const handleRejectConfirm = () => {
    if (!rejectingId || !rejectReason.trim()) return;
    rejectVisitReport.mutate(
      { id: rejectingId, data: { reason: rejectReason } },
      {
        onSuccess: () => {
          setRejectingId(null);
          setRejectReason("");
        },
      }
    );
  };

  const handleRejectDialogClose = (isOpen: boolean) => {
    if (!isOpen) {
      setRejectingId(null);
      setRejectReason("");
    }
  };

  return (
    <>
      {/* Layer 2: list of visit reports for the selected sales rep */}
      <Dialog
        open={open}
        onOpenChange={(isOpen) => {
          // Prevent the dialog from closing while the top-level drawer is open.
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
            {rep.submitted > 0 && (
              <DialogDescription className="flex items-center gap-1.5 text-amber-600 dark:text-amber-400">
                <Bell className="h-3.5 w-3.5 shrink-0" />
                {rep.submitted > 1
                  ? t("dialog.waitingApprovalPlural", { count: rep.submitted })
                  : t("dialog.waitingApproval", { count: rep.submitted })}
              </DialogDescription>
            )}
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
                    onEdit={
                      hasEditPermission ? onViewReport : undefined
                    }
                    onDelete={hasDeletePermission ? setDeletingId : undefined}
                    onApprove={
                      hasApprovePermission ? handleApprove : undefined
                    }
                    onReject={
                      hasRejectPermission ? setRejectingId : undefined
                    }
                    onSubmit={handleSubmit}
                    isApproving={approveVisitReport.isPending}
                    isRejecting={rejectVisitReport.isPending}
                    rejectingId={rejectingId}
                    isUpdating={updateVisitReport.isPending}
                  />
                ))}
              </div>
            )}
          </div>

          {/* Pagination */}
          {totalPages > 1 && (
            <div className="shrink-0 border-t px-4 py-3">
              <Pagination>
                <PaginationContent>
                  <PaginationItem>
                    <PaginationPrevious
                      onClick={() => {
                        if (page > 1) setPage(page - 1);
                      }}
                      disabled={page === 1}
                    />
                  </PaginationItem>
                  {pageNumbers.map((p, idx) => {
                    if (p === "...") {
                      const nextPage = pageNumbers[idx + 1];
                      return (
                        <PaginationItem key={`ellipsis-before-${nextPage}`}>
                          <PaginationEllipsis />
                        </PaginationItem>
                      );
                    }
                    return (
                      <PaginationItem key={p}>
                        <PaginationLink
                          isActive={p === page}
                          onClick={() => setPage(p)}
                        >
                          {p}
                        </PaginationLink>
                      </PaginationItem>
                    );
                  })}
                  <PaginationItem>
                    <PaginationNext
                      onClick={() => {
                        if (page < totalPages) setPage(page + 1);
                      }}
                      disabled={page === totalPages}
                    />
                  </PaginationItem>
                </PaginationContent>
              </Pagination>
            </div>
          )}
        </DialogContent>
      </Dialog>

      {/* Delete confirmation dialog */}
      <DeleteDialog
        open={!!deletingId}
        onOpenChange={(isOpen) => {
          if (!isOpen) setDeletingId(null);
        }}
        onConfirm={handleDeleteConfirm}
        title={t("deleteDialog.title")}
        description={t("deleteDialog.description")}
        itemName={t("deleteDialog.itemName")}
        isLoading={deleteVisitReport.isPending}
      />

      {/* Reject reason dialog */}
      <Dialog open={!!rejectingId} onOpenChange={handleRejectDialogClose}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{t("rejectDialog.title")}</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div>
              <Label className="text-sm font-medium mb-2 block">
                {t("rejectDialog.reasonLabel")} *
              </Label>
              <Textarea
                value={rejectReason}
                onChange={(e) => setRejectReason(e.target.value)}
                placeholder={t("rejectDialog.reasonPlaceholder")}
                className="min-h-[100px]"
                rows={4}
              />
            </div>
            <div className="flex justify-end gap-2">
              <Button
                variant="outline"
                onClick={() => handleRejectDialogClose(false)}
                disabled={rejectVisitReport.isPending}
              >
                {t("rejectDialog.cancel")}
              </Button>
              <Button
                onClick={handleRejectConfirm}
                disabled={
                  rejectVisitReport.isPending || !rejectReason.trim()
                }
                variant="destructive"
              >
                {rejectVisitReport.isPending
                  ? t("rejectDialog.rejecting")
                  : t("rejectDialog.reject")}
              </Button>
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </>
  );
}
