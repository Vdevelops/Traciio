"use client";

import { useState } from "react";
import { Search, Eye, Calendar, MapPin } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { DataTable, type Column } from "@/components/ui/data-table";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useVisitReportList } from "../hooks/useVisitReportList";
import { VisitReportDetailModal } from "./visit-report-detail-modal";
import { DateRangePicker } from "@/components/ui/date-range-picker";
import { useAccounts } from "../../account-management/hooks/useAccounts";
import { useDeals } from "../../pipeline-management/hooks/useDeals";
import type { Deal } from "../../pipeline-management/types";
import type { VisitReport } from "../types";
import { useTranslations } from "next-intl";

export function VisitReportList() {
  const t = useTranslations("visitReportList");
  const {
    setPage,
    setPerPage,
    search,
    setSearch,
    accountId,
    setAccountId,
    dealId,
    setDealId,
    startDate,
    setStartDate,
    endDate,
    setEndDate,
    visitReports,
    pagination,
    isLoading,
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
          setAccountId("");
          setStartDate("");
          setEndDate("");
          setPage(1);
        }}
      />

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
    </div>
  );
}
