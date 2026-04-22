"use client";

import { useState } from "react";
import { Edit, Trash2, Plus, Search, Eye, Calendar, MapPin, CheckCircle2, XCircle, MoreVertical, Send } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { DataTable, type Column } from "@/components/ui/data-table";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useVisitReportList } from "../hooks/useVisitReportList";
import { VisitReportForm } from "./visit-report-form";
import { VisitReportDetailModal } from "./visit-report-detail-modal";
import { DeleteDialog } from "@/components/ui/delete-dialog";
import { DateRangePicker } from "@/components/ui/date-range-picker";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { useAccounts } from "../../account-management/hooks/useAccounts";
import { useDeals } from "../../pipeline-management/hooks/useDeals";
import type { Deal } from "../../pipeline-management/types";
import { useHasPermission } from "@/features/master-data/user-management/hooks/useHasPermission";
import type { VisitReport } from "../types";
import type { CreateVisitReportFormData } from "../schemas/visit-report.schema";
import { useTranslations } from "next-intl";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";

const statusColors: Record<string, "default" | "secondary" | "destructive" | "outline"> = {
  draft: "outline",
  submitted: "secondary",
  approved: "default",
  rejected: "destructive",
};

export function VisitReportList() {
  const t = useTranslations("visitReportList");
  const tDetail = useTranslations("visitReportDetail");
  
  // Permission checks
  const hasCreatePermission = useHasPermission("visit-reports.create");
  const hasEditPermission = useHasPermission("visit-reports.edit");
  const hasDeletePermission = useHasPermission("visit-reports.delete");
  const hasApprovePermission = useHasPermission("visit-reports.approve");
  const hasRejectPermission = useHasPermission("visit-reports.reject");
  const {
    setPage,
    setPerPage,
    search,
    setSearch,
    status,
    setStatus,
    accountId,
    setAccountId,
    dealId,
    setDealId,
    startDate,
    setStartDate,
    endDate,
    setEndDate,
    isCreateDialogOpen,
    setIsCreateDialogOpen,
    editingVisitReport,
    setEditingVisitReport,
    visitReports,
    pagination,
    editingVisitReportData,
    isLoading,
    handleCreate,
    handleUpdate,
    handleDeleteClick,
    handleDeleteConfirm,
    deletingVisitReportId,
    setDeletingVisitReportId,
    deleteVisitReport,
    createVisitReport,
    updateVisitReport,
    handleApprove,
    handleRejectClick,
    handleSubmit,
    approvingVisitReportId,
    rejectingVisitReportId,
    setRejectingVisitReportId,
    rejectReason,
    setRejectReason,
    handleRejectConfirm,
    approveVisitReport,
    rejectVisitReport,
  } = useVisitReportList();

  const { data: accountsData } = useAccounts({ per_page: 100 });
  const accounts = accountsData?.data || [];
  const { data: dealsData } = useDeals(undefined, 1, 100);
  const deals = dealsData?.data || [];

  const [viewingVisitReportId, setViewingVisitReportId] = useState<string | null>(null);
  const [isDetailModalOpen, setIsDetailModalOpen] = useState(false);

  const handleViewVisitReport = (id: string) => {
    setViewingVisitReportId(id);
    setIsDetailModalOpen(true);
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString("id-ID", {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
  };

  const columns: Column<VisitReport>[] = [
    {
      id: "visit_date",
      header: t("table.visitDate"),
      accessor: (row) => (
        <div className="flex items-center gap-2">
          <Calendar className="h-4 w-4 text-muted-foreground" />
          <span>{formatDate(row.visit_date)}</span>
        </div>
      ),
      className: "w-[150px]",
    },
    {
      id: "account",
      header: t("table.account"),
      accessor: (row) => (
        <button
          onClick={() => handleViewVisitReport(row.id)}
          className="font-medium text-primary hover:underline text-left"
        >
          {row.account?.name ?? "N/A"}
        </button>
      ),
      className: "w-[200px]",
    },
    {
      id: "deal",
      header: t("table.deal"),
      accessor: (row) => (
        <span className="text-sm text-muted-foreground">
          {row.deal?.title ?? "-"}
        </span>
      ),
      className: "w-[180px]",
    },
    {
      id: "purpose",
      header: t("table.purpose"),
      accessor: (row) => (
        <span className="text-muted-foreground line-clamp-1">{row.purpose}</span>
      ),
    },
    {
      id: "status",
      header: t("table.status"),
      accessor: (row) => (
        <Badge variant={statusColors[row.status] || "outline"}>
          {row.status}
        </Badge>
      ),
      className: "w-[120px]",
    },
    {
      id: "check_in",
      header: t("table.checkIn"),
      accessor: (row) => (
        <div className="flex items-center gap-1 text-sm text-muted-foreground">
          {row.check_in_time ? (
            <>
              <MapPin className="h-3 w-3" />
              <span>{new Date(row.check_in_time).toLocaleTimeString("id-ID", { hour: "2-digit", minute: "2-digit" })}</span>
            </>
          ) : (
            <span>-</span>
          )}
        </div>
      ),
      className: "w-[120px]",
    },
    {
      id: "actions",
      header: t("table.actions"),
      accessor: (row) => {
        const hasApprovalActions = (hasApprovePermission || hasRejectPermission) && row.status === "submitted";
        const canEdit = hasEditPermission && (row.status === "draft" || row.status === "submitted");
        const canSubmit = row.status === "draft";
        const canDelete = hasDeletePermission;
        const hasAnyAction = hasApprovalActions || canEdit || canSubmit || canDelete;

        return (
          <div className="flex items-center justify-end gap-1">
            <Button
              variant="ghost"
              size="icon-sm"
              className="h-8 w-8"
              title={t("buttons.viewDetails")}
              onClick={() => handleViewVisitReport(row.id)}
            >
              <Eye className="h-3.5 w-3.5" />
            </Button>
            {hasAnyAction && (
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button
                    variant="ghost"
                    size="icon-sm"
                    className="h-8 w-8"
                    title="More actions"
                  >
                    <MoreVertical className="h-3.5 w-3.5" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" className="w-48">
                  {hasApprovalActions && (
                    <>
                      {hasApprovePermission && row.status === "submitted" && (
                        <DropdownMenuItem
                          onClick={() => handleApprove(row.id)}
                          disabled={approveVisitReport.isPending && approvingVisitReportId === row.id}
                          className="text-[color:var(--color-success)] focus:text-[color:var(--color-success)] hover:bg-[color:var(--color-success)]/10 dark:focus:bg-[color:var(--color-success)]/90"
                        >
                          <CheckCircle2 className="h-4 w-4 mr-2" />
                          {t("buttons.approve")}
                        </DropdownMenuItem>
                      )}
                      {hasRejectPermission && row.status === "submitted" && (
                        <DropdownMenuItem
                          onClick={() => handleRejectClick(row.id)}
                          disabled={rejectVisitReport.isPending && rejectingVisitReportId === row.id}
                          variant="destructive"
                        >
                          <XCircle className="h-4 w-4 mr-2" />
                          {t("buttons.reject")}
                        </DropdownMenuItem>
                      )}
                      {(canEdit || canSubmit || canDelete) && <DropdownMenuSeparator />}
                    </>
                  )}
                  {canSubmit && (
                    <DropdownMenuItem
                      onClick={() => handleSubmit(row.id)}
                      disabled={updateVisitReport.isPending}
                      className="text-primary focus:text-primary"
                    >
                      <Send className="h-4 w-4 mr-2" />
                      {t("buttons.submit")}
                    </DropdownMenuItem>
                  )}
                  {canEdit && (
                    <DropdownMenuItem
                      onClick={() => setEditingVisitReport(row.id)}
                    >
                      <Edit className="h-4 w-4 mr-2" />
                      {t("buttons.edit")}
                    </DropdownMenuItem>
                  )}
                  {canDelete && (
                    <DropdownMenuItem
                      onClick={() => handleDeleteClick(row.id)}
                      variant="destructive"
                    >
                      <Trash2 className="h-4 w-4 mr-2" />
                      {t("buttons.delete")}
                    </DropdownMenuItem>
                  )}
                </DropdownMenuContent>
              </DropdownMenu>
            )}
          </div>
        );
      },
      className: "w-[120px] text-right",
    },
  ];

  return (
    <div className="space-y-4">
      {/* Header with Actions */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3 flex-1 flex-wrap">
          <div className="relative flex-1 max-w-sm">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder={t("filters.searchPlaceholder")}
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="pl-10 h-9"
            />
          </div>
          <Select 
            value={status || "all"} 
            onValueChange={(value) => setStatus(value === "all" ? "" : value)}
          >
            <SelectTrigger className="w-[140px] h-9">
              <SelectValue placeholder={t("filters.allStatus")} />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">{t("filters.allStatus")}</SelectItem>
              <SelectItem value="draft">draft</SelectItem>
              <SelectItem value="submitted">submitted</SelectItem>
              <SelectItem value="approved">approved</SelectItem>
              <SelectItem value="rejected">rejected</SelectItem>
            </SelectContent>
          </Select>
          <Select 
            value={accountId || "all"} 
            onValueChange={(value) => setAccountId(value === "all" ? "" : value)}
          >
            <SelectTrigger className="w-[180px] h-9">
              <SelectValue placeholder={t("filters.allAccounts")} />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">{t("filters.allAccounts")}</SelectItem>
              {accounts.map((account) => (
                <SelectItem key={account.id} value={account.id}>
                  {account.name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          <Select 
            value={dealId || "all"} 
            onValueChange={(value) => setDealId(value === "all" ? "" : value)}
          >
            <SelectTrigger className="w-[180px] h-9">
              <SelectValue placeholder={t("filters.allDeals")} />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">{t("filters.allDeals")}</SelectItem>
              {deals.map((deal: Deal) => (
                <SelectItem key={deal.id} value={deal.id}>
                  {deal.title}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          <DateRangePicker
            dateRange={(() => {
              if (startDate && endDate) {
                const fromDate = new Date(startDate + "T00:00:00");
                fromDate.setHours(0, 0, 0, 0);
                const toDate = new Date(endDate + "T00:00:00");
                toDate.setHours(0, 0, 0, 0);
                return { from: fromDate, to: toDate };
              }
              if (startDate) {
                const fromDate = new Date(startDate + "T00:00:00");
                fromDate.setHours(0, 0, 0, 0);
                return { from: fromDate, to: undefined };
              }
              return undefined;
            })()}
            onDateChange={(range) => {
              if (range?.from) {
                const fromDate = new Date(range.from);
                fromDate.setHours(0, 0, 0, 0);
                const fromStr = `${fromDate.getFullYear()}-${String(fromDate.getMonth() + 1).padStart(2, "0")}-${String(fromDate.getDate()).padStart(2, "0")}`;
                setStartDate(fromStr);
                
                if (range.to) {
                  const toDate = new Date(range.to);
                  toDate.setHours(0, 0, 0, 0);
                  const toStr = `${toDate.getFullYear()}-${String(toDate.getMonth() + 1).padStart(2, "0")}-${String(toDate.getDate()).padStart(2, "0")}`;
                  setEndDate(toStr);
                } else {
                  setEndDate("");
                }
              } else {
                setStartDate("");
                setEndDate("");
              }
            }}
          />
        </div>
        {hasCreatePermission && (
          <Button onClick={() => setIsCreateDialogOpen(true)} size="sm">
            <Plus className="h-4 w-4 mr-2" />
            {t("buttons.addVisitReport")}
          </Button>
        )}
      </div>

      {/* Table */}
      <DataTable
        columns={columns}
        data={visitReports}
        isLoading={isLoading}
        emptyMessage={t("empty.table")}
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
        itemName={t("dialogs.itemName")}
        perPageOptions={[10, 20, 50, 100]}
        onResetFilters={() => {
          setSearch("");
          setStatus("");
          setAccountId("");
          setStartDate("");
          setEndDate("");
          setPage(1);
        }}
      />

      {/* Create Dialog */}
      <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
        <DialogContent className="sm:max-w-[600px]">
          <DialogHeader>
            <DialogTitle>{t("dialogs.createTitle")}</DialogTitle>
          </DialogHeader>
          <VisitReportForm
            key="create-visit-report-form"
            open={isCreateDialogOpen}
            onSubmit={async (data) => {
              await handleCreate(data as CreateVisitReportFormData);
            }}
            onCancel={() => setIsCreateDialogOpen(false)}
            isLoading={createVisitReport.isPending}
          />
        </DialogContent>
      </Dialog>

      {/* Edit Dialog */}
      {editingVisitReport && editingVisitReportData?.data && (
        <Dialog open={!!editingVisitReport} onOpenChange={(open) => !open && setEditingVisitReport(null)}>
          <DialogContent className="sm:max-w-[600px]">
            <DialogHeader>
            <DialogTitle>{t("dialogs.editTitle")}</DialogTitle>
            </DialogHeader>
            <VisitReportForm
              visitReport={editingVisitReportData.data}
            onSubmit={(data) => handleUpdate(data)}
              onCancel={() => setEditingVisitReport(null)}
              isLoading={updateVisitReport.isPending}
            />
          </DialogContent>
        </Dialog>
      )}

      {/* Visit Report Detail Modal */}
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
          // Refresh will be handled by query invalidation in hooks
        }}
      />

      {/* Delete Dialog */}
      <DeleteDialog
        open={!!deletingVisitReportId}
        onOpenChange={(open) => {
          if (!open) {
            setDeletingVisitReportId(null);
          }
        }}
        onConfirm={handleDeleteConfirm}
        title={t("dialogs.deleteTitle")}
        description={t("dialogs.deleteDescription")}
        itemName={t("dialogs.itemName")}
        isLoading={deleteVisitReport.isPending}
      />

      {/* Reject Dialog */}
      <Dialog open={!!rejectingVisitReportId} onOpenChange={(open) => {
        if (!open) {
          setRejectingVisitReportId(null);
          setRejectReason("");
        }
      }}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{tDetail("rejectDialog.title")}</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div>
              <Label className="text-sm font-medium mb-2 block">
                {tDetail("rejectDialog.reasonLabel")} *
              </Label>
              <Textarea
                value={rejectReason}
                onChange={(e) => setRejectReason(e.target.value)}
                placeholder={tDetail("rejectDialog.reasonPlaceholder")}
                className="min-h-[100px]"
                rows={4}
              />
            </div>
            <div className="flex justify-end gap-2">
              <Button
                variant="outline"
                onClick={() => {
                  setRejectingVisitReportId(null);
                  setRejectReason("");
                }}
                disabled={rejectVisitReport.isPending}
              >
                {tDetail("rejectDialog.cancel")}
              </Button>
              <Button
                onClick={handleRejectConfirm}
                disabled={rejectVisitReport.isPending || !rejectReason.trim()}
                variant="destructive"
              >
                {rejectVisitReport.isPending ? tDetail("rejectDialog.submitting") : tDetail("rejectDialog.submit")}
              </Button>
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
}

