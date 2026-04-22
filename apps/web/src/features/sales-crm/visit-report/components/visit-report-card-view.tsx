"use client";

import { useState, useMemo } from "react";
import {
  CheckCircle2,
  Clock,
  FileText,
  XCircle,
  ChevronRight,
  CalendarDays,
  Send,
  Edit,
  Trash2,
  Eye,
  MoreVertical,
  Users,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
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
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
  PaginationEllipsis,
} from "@/components/ui/pagination";
import {
  useVisitReports,
  useApproveVisitReport,
  useRejectVisitReport,
  useDeleteVisitReport,
  useUpdateVisitReport,
} from "../hooks/useVisitReports";
import { VisitReportDetailModal } from "./visit-report-detail-modal";
import { useHasPermission } from "@/features/master-data/user-management/hooks/useHasPermission";
import { useUsers } from "@/features/master-data/user-management/hooks/useUsers";
import type { VisitReport } from "../types";
import { useTranslations } from "next-intl";

// ─── Status config ────────────────────────────────────────────────────────────

const STATUS_CONFIG: Record<
  string,
  {
    variant: "default" | "secondary" | "destructive" | "outline";
    icon: typeof CheckCircle2;
    label: string;
  }
> = {
  draft: { variant: "outline", icon: FileText, label: "Draft" },
  submitted: { variant: "secondary", icon: Clock, label: "Submitted" },
  approved: { variant: "default", icon: CheckCircle2, label: "Approved" },
  rejected: { variant: "destructive", icon: XCircle, label: "Rejected" },
};

// ─── Helpers ──────────────────────────────────────────────────────────────────

function getAvatarSrc(avatarUrl?: string, seed?: string): string {
  if (avatarUrl && avatarUrl.trim() !== "") return avatarUrl;
  const s = seed ?? "";
  return `https://api.dicebear.com/7.x/lorelei/svg?seed=${encodeURIComponent(s)}`;
}

function getInitials(name: string): string {
  return name
    .split(" ")
    .slice(0, 2)
    .map((n) => n[0])
    .join("")
    .toUpperCase();
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
  for (let i = Math.max(2, current - 1); i <= Math.min(total - 1, current + 1); i++) {
    pages.push(i);
  }
  if (current < total - 2) pages.push("...");
  pages.push(total);
  return pages;
}

// ─── Types ────────────────────────────────────────────────────────────────────

interface SalesRepGroup {
  id: string;
  name: string;
  email: string;
  avatarUrl?: string;
  total: number;
  approved: number;
  submitted: number;
  draft: number;
  rejected: number;
}

// ─── UserCard ─────────────────────────────────────────────────────────────────

interface UserCardProps {
  readonly rep: SalesRepGroup;
  readonly onClick: (id: string, name: string, avatarUrl?: string) => void;
}

function UserCard({ rep, onClick }: UserCardProps) {
  const src = getAvatarSrc(rep.avatarUrl, rep.email || rep.name);
  return (
    <Card
      className="cursor-pointer transition-all hover:shadow-sm hover:border-primary/30 active:scale-[0.99]"
      onClick={() => onClick(rep.id, rep.name, rep.avatarUrl)}
    >
      <CardContent className="flex items-center gap-4 px-4 py-3">
        {/* Avatar — same component pattern as header profile */}
        <Avatar className="h-10 w-10 shrink-0">
          <AvatarImage src={src} alt={rep.name} />
          <AvatarFallback className="text-sm font-semibold bg-primary/10 text-primary">
            {getInitials(rep.name)}
          </AvatarFallback>
        </Avatar>

        {/* Name + stats */}
        <div className="flex-1 min-w-0">
          <p className="text-sm font-semibold truncate">{rep.name}</p>
          <div className="flex items-center gap-2 mt-1 flex-wrap">
            <span className="text-xs text-muted-foreground">{rep.total} visits</span>
            {rep.approved > 0 && (
              <span className="inline-flex items-center gap-0.5 text-xs text-[color:var(--color-success)] dark:text-[color:var(--color-success-foreground)]">
                <CheckCircle2 className="h-3 w-3" />
                {rep.approved}
              </span>
            )}
            {rep.submitted > 0 && (
              <span className="inline-flex items-center gap-0.5 text-xs text-blue-600 dark:text-blue-400">
                <Clock className="h-3 w-3" />
                {rep.submitted}
              </span>
            )}
            {rep.draft > 0 && (
              <span className="inline-flex items-center gap-0.5 text-xs text-muted-foreground">
                <FileText className="h-3 w-3" />
                {rep.draft}
              </span>
            )}
            {rep.rejected > 0 && (
              <span className="inline-flex items-center gap-0.5 text-xs text-destructive">
                <XCircle className="h-3 w-3" />
                {rep.rejected}
              </span>
            )}
          </div>
        </div>

        <ChevronRight className="h-4 w-4 text-muted-foreground shrink-0" />
      </CardContent>
    </Card>
  );
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
  readonly approvingId: string | null;
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
  approvingId,
  rejectingId,
  isUpdating,
}: ReportRowProps) {
  const cfg = STATUS_CONFIG[report.status] ?? STATUS_CONFIG.draft;
  const StatusIcon = cfg.icon;

  const canApprove = !!onApprove && report.status === "submitted";
  const canReject = !!onReject && report.status === "submitted";
  const canEdit = !!onEdit && (report.status === "draft" || report.status === "submitted");
  const canSubmit = report.status === "draft";
  const hasActions = canApprove || canReject || canEdit || canSubmit || !!onDelete;

  return (
    <Card className="overflow-hidden">
      <CardContent className="p-3 space-y-2.5">
        {/* Top row: date + status + actions */}
        <div className="flex items-center justify-between gap-2">
          <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
            <CalendarDays className="h-3.5 w-3.5 shrink-0" />
            <span>{formatDate(report.visit_date)}</span>
          </div>

          <div className="flex items-center gap-1.5">
            <Badge variant={cfg.variant} className="text-xs gap-1">
              <StatusIcon className="h-3 w-3" />
              {cfg.label}
            </Badge>

            {hasActions && (
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="ghost" size="icon" className="h-7 w-7">
                    <MoreVertical className="h-3.5 w-3.5" />
                  </Button>
                </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-44">
            <DropdownMenuItem onClick={() => onView(report.id)}>
              <Eye className="h-4 w-4 mr-2" />
              View
            </DropdownMenuItem>
            {(canApprove || canReject) && <DropdownMenuSeparator />}
            {canApprove && (
              <DropdownMenuItem
                onClick={() => onApprove && onApprove(report.id)}
                disabled={isApproving && approvingId === report.id}
                className="text-[color:var(--color-success)] focus:text-[color:var(--color-success)] hover:bg-[color:var(--color-success)]/10 dark:focus:bg-[color:var(--color-success)]/90"
              >
                <CheckCircle2 className="h-4 w-4 mr-2" />
                Approve
              </DropdownMenuItem>
            )}
            {canReject && (
              <DropdownMenuItem
                onClick={() => onReject && onReject(report.id)}
                disabled={isRejecting && rejectingId === report.id}
                variant="destructive"
              >
                <XCircle className="h-4 w-4 mr-2" />
                Reject
              </DropdownMenuItem>
            )}
            {(canEdit || canSubmit || !!onDelete) && (canApprove || canReject) && (
              <DropdownMenuSeparator />
            )}
            {canSubmit && (
              <DropdownMenuItem
                onClick={() => onSubmit(report.id)}
                disabled={isUpdating}
                className="text-primary"
              >
                <Send className="h-4 w-4 mr-2" />
                Submit
              </DropdownMenuItem>
            )}
            {canEdit && (
              <DropdownMenuItem onClick={() => onEdit && onEdit(report.id)}>
                <Edit className="h-4 w-4 mr-2" />
                Edit
              </DropdownMenuItem>
            )}
            {onDelete && (
              <DropdownMenuItem
                onClick={() => onDelete(report.id)}
                variant="destructive"
              >
                <Trash2 className="h-4 w-4 mr-2" />
                Delete
              </DropdownMenuItem>
            )}
          </DropdownMenuContent>
              </DropdownMenu>
            )}
          </div>
        </div>

        {/* Account name — tappable */}
        <button
          onClick={() => onView(report.id)}
          className="w-full text-sm font-medium text-left truncate hover:text-primary transition-colors"
        >
          {report.account?.name ?? "—"}
        </button>
      </CardContent>
    </Card>
  );
}

// ─── UserVisitReportsDialog ───────────────────────────────────────────────────

const PER_PAGE = 10;

interface UserVisitReportsDialogProps {
  readonly open: boolean;
  readonly onOpenChange: (open: boolean) => void;
  readonly userId: string;
  readonly userName: string;
  readonly userAvatarUrl?: string;
  readonly userEmail?: string;
}

function UserVisitReportsDialog({
  open,
  onOpenChange,
  userId,
  userName,
  userAvatarUrl,
  userEmail,
}: UserVisitReportsDialogProps) {
  const [page, setPage] = useState(1);
  const [viewingId, setViewingId] = useState<string | null>(null);
  const [isDetailOpen, setIsDetailOpen] = useState(false);
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
    sales_rep_id: userId,
    page,
    per_page: PER_PAGE,
  });

  const reports: VisitReport[] = data?.data ?? [];
  const pagination = data?.meta?.pagination;
  const totalPages = pagination?.total_pages ?? 1;

  const pageNumbers = useMemo(
    () => buildPageNumbers(page, totalPages),
    [page, totalPages]
  );

  const handleApprove = (id: string) => {
    approveVisitReport.mutate(id);
  };

  const handleSubmit = (id: string) => {
    updateVisitReport.mutate({ id, data: { status: "submitted" } });
  };

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

  return (
    <>
      <Dialog open={open} onOpenChange={onOpenChange}>
        <DialogContent className="sm:max-w-[640px] max-h-[80vh] flex flex-col gap-0 p-0">
          <DialogHeader className="px-5 pt-5 pb-3 shrink-0">
            <DialogTitle className="flex items-center gap-2 text-base">
              <Avatar className="h-8 w-8">
                <AvatarImage
                  src={getAvatarSrc(userAvatarUrl, userEmail || userName)}
                  alt={userName}
                />
                <AvatarFallback className="text-xs font-semibold bg-primary/10 text-primary">
                  {getInitials(userName)}
                </AvatarFallback>
              </Avatar>
              {userName}
            </DialogTitle>
          </DialogHeader>

          <div className="flex-1 overflow-y-auto px-3 pb-2">
            {isLoading ? (
              <div className="space-y-2 px-2 py-2">
                {Array.from({ length: 5 }).map((_, i) => (
                  <Skeleton key={`sk-${i}`} className="h-10 rounded-lg" />
                ))}
              </div>
            ) : reports.length === 0 ? (
              <div className="flex flex-col items-center justify-center py-14 text-center">
                <CalendarDays className="h-10 w-10 text-muted-foreground/30 mb-3" />
                <p className="text-sm text-muted-foreground">No visit reports found</p>
              </div>
            ) : (
              <div className="space-y-0.5">
                {reports.map((report) => (
                  <ReportRow
                    key={report.id}
                    report={report}
                    onView={(id) => {
                      setViewingId(id);
                      setIsDetailOpen(true);
                    }}
                    onEdit={hasEditPermission ? (id) => { setViewingId(id); setIsDetailOpen(true); } : undefined}
                    onDelete={hasDeletePermission ? setDeletingId : undefined}
                    onApprove={hasApprovePermission ? handleApprove : undefined}
                    onReject={hasRejectPermission ? setRejectingId : undefined}
                    onSubmit={handleSubmit}
                    isApproving={approveVisitReport.isPending}
                    isRejecting={rejectVisitReport.isPending}
                    approvingId={null}
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
                      onClick={() => { if (page > 1) setPage(page - 1); }}
                      disabled={page === 1}
                    />
                  </PaginationItem>
                  {pageNumbers.map((p, idx) =>
                    p === "..." ? (
                      <PaginationItem key={`ellipsis-${idx}`}>
                        <PaginationEllipsis />
                      </PaginationItem>
                    ) : (
                      <PaginationItem key={p}>
                        <PaginationLink
                          isActive={p === page}
                          onClick={() => setPage(p as number)}
                        >
                          {p}
                        </PaginationLink>
                      </PaginationItem>
                    )
                  )}
                  <PaginationItem>
                    <PaginationNext
                      onClick={() => { if (page < totalPages) setPage(page + 1); }}
                      disabled={page === totalPages}
                    />
                  </PaginationItem>
                </PaginationContent>
              </Pagination>
            </div>
          )}
        </DialogContent>
      </Dialog>

      {/* Detail modal */}
      <VisitReportDetailModal
        visitReportId={viewingId}
        open={isDetailOpen}
        onOpenChange={(open) => {
          setIsDetailOpen(open);
          if (!open) setViewingId(null);
        }}
        onVisitReportUpdated={() => {}}
      />

      {/* Delete dialog */}
      <DeleteDialog
        open={!!deletingId}
        onOpenChange={(open) => { if (!open) setDeletingId(null); }}
        onConfirm={handleDeleteConfirm}
        title="Delete Visit Report"
        description="Are you sure you want to delete this visit report? This action cannot be undone."
        itemName="visit report"
        isLoading={deleteVisitReport.isPending}
      />

      {/* Reject dialog */}
      <Dialog
        open={!!rejectingId}
        onOpenChange={(open) => {
          if (!open) {
            setRejectingId(null);
            setRejectReason("");
          }
        }}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Reject Visit Report</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div>
              <Label className="text-sm font-medium mb-2 block">Reason *</Label>
              <Textarea
                value={rejectReason}
                onChange={(e) => setRejectReason(e.target.value)}
                placeholder="Enter rejection reason..."
                className="min-h-[100px]"
                rows={4}
              />
            </div>
            <div className="flex justify-end gap-2">
              <Button
                variant="outline"
                onClick={() => {
                  setRejectingId(null);
                  setRejectReason("");
                }}
                disabled={rejectVisitReport.isPending}
              >
                Cancel
              </Button>
              <Button
                onClick={handleRejectConfirm}
                disabled={rejectVisitReport.isPending || !rejectReason.trim()}
                variant="destructive"
              >
                {rejectVisitReport.isPending ? "Rejecting..." : "Reject"}
              </Button>
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </>
  );
}

// ─── Main export ──────────────────────────────────────────────────────────────

export function VisitReportCardView() {
  const t = useTranslations("visitReportManagement");

  const [selectedUser, setSelectedUser] = useState<{
    id: string;
    name: string;
    avatarUrl?: string;
    email?: string;
  } | null>(null);

  const { data, isLoading } = useVisitReports({ per_page: 100 });
  const allReports: VisitReport[] = data?.data ?? [];

  // Fetch all users to enrich sales rep cards with avatar and email
  const { data: usersData } = useUsers({ per_page: 100 });
  const userMap = useMemo(() => {
    const map = new Map<string, { avatar_url?: string; email?: string }>();
    for (const u of usersData?.data ?? []) {
      map.set(u.id, { avatar_url: u.avatar_url, email: u.email });
    }
    return map;
  }, [usersData?.data]);

  // Group reports by sales rep
  const repGroups = useMemo<SalesRepGroup[]>(() => {
    const map = new Map<string, SalesRepGroup>();

    for (const report of allReports) {
      const id = report.sales_rep?.id ?? report.sales_rep_id ?? "unknown";
      const name = report.sales_rep?.name ?? "Unknown";

      const existing = map.get(id);
      const userData = userMap.get(id);
      if (existing) {
        existing.total += 1;
        if (report.status === "approved") existing.approved += 1;
        else if (report.status === "submitted") existing.submitted += 1;
        else if (report.status === "draft") existing.draft += 1;
        else if (report.status === "rejected") existing.rejected += 1;
      } else {
        map.set(id, {
          id,
          name,
          email: userData?.email ?? "",
          avatarUrl: userData?.avatar_url,
          total: 1,
          approved: report.status === "approved" ? 1 : 0,
          submitted: report.status === "submitted" ? 1 : 0,
          draft: report.status === "draft" ? 1 : 0,
          rejected: report.status === "rejected" ? 1 : 0,
        });
      }
    }

    return Array.from(map.values()).sort((a, b) => b.total - a.total);
  }, [allReports, userMap]);

  return (
    <div className="space-y-3">
      {isLoading ? (
        <div className="space-y-2">
          {Array.from({ length: 5 }).map((_, i) => (
            <Skeleton key={`sk-${i}`} className="h-16 rounded-xl" />
          ))}
        </div>
      ) : repGroups.length === 0 ? (
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-16 text-center">
            <Users className="h-12 w-12 text-muted-foreground/30 mb-4" />
            <p className="text-sm text-muted-foreground">
              {t("empty") ?? "No team visit data available"}
            </p>
          </CardContent>
        </Card>
      ) : (
        <div className="grid grid-cols-1 gap-2">
          {repGroups.map((rep) => (
            <UserCard
              key={rep.id}
              rep={rep}
              onClick={(id, name, avatarUrl) =>
                setSelectedUser({ id, name, avatarUrl, email: rep.email })
              }
            />
          ))}
        </div>
      )}

      {selectedUser && (
        <UserVisitReportsDialog
          open={!!selectedUser}
          onOpenChange={(open) => { if (!open) setSelectedUser(null); }}
          userId={selectedUser.id}
          userName={selectedUser.name}
          userAvatarUrl={selectedUser.avatarUrl}
          userEmail={selectedUser.email}
        />
      )}
    </div>
  );
}
