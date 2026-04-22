"use client";

import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import { DataTable, type Column } from "@/components/ui/data-table";
import { Avatar, AvatarImage } from "@/components/ui/avatar";
import type { SalesPerformanceListItem } from "../types";
import { formatEmailToMailto } from "@/lib/utils";

interface SalesPerformanceTableProps {
  readonly data: readonly SalesPerformanceListItem[];
  readonly isLoading?: boolean;
  readonly pagination?: {
    readonly page: number;
    readonly per_page: number;
    readonly total: number;
    readonly total_pages: number;
    readonly has_next: boolean;
    readonly has_prev: boolean;
  };
  readonly onPageChange?: (page: number) => void;
  readonly onPerPageChange?: (perPage: number) => void;
  readonly emptyMessage?: string;
}

type SalesPerformanceWithId = SalesPerformanceListItem & { id: string };

export function SalesPerformanceTable({
  data,
  isLoading,
  pagination,
  onPageChange,
  onPerPageChange,
  emptyMessage,
}: SalesPerformanceTableProps) {
  const t = useTranslations("salesOverview");
  const router = useRouter();

  const handleUserClick = (userId: string) => {
    router.push(`/sales-overview/sales-rep/${userId}`);
  };

  // Map data to add 'id' property from 'user_id' for DataTable key requirement
  const mappedData: SalesPerformanceWithId[] = data.map((item) => ({
    ...item,
    id: item.user_id,
  }));

  const columns: Column<SalesPerformanceWithId>[] = [
    {
      id: "user",
      header: t("table.user"),
      accessor: (row) => (
        <button
          onClick={() => handleUserClick(row.user_id)}
          className="flex items-center gap-3 font-medium text-primary hover:underline cursor-pointer"
        >
          <Avatar className="h-8 w-8">
            <AvatarImage src={row.avatar_url} alt={row.user_name ?? "User"} />
          </Avatar>
          <span>{row.user_name ?? "-"}</span>
        </button>
      ),
      className: "w-[200px]",
    },
    {
      id: "email",
      header: t("table.email"),
      accessor: (row) => (
        <a href={formatEmailToMailto(row.user_email)} className="text-muted-foreground hover:text-primary hover:underline cursor-pointer min-w-0">
          {row.user_email ?? "-"}
        </a>
      ),
    },
    {
      id: "revenue",
      header: t("table.revenue"),
      accessor: (row) => (
        <div className="text-right">
          <span className="font-medium">{row.total_revenue_formatted}</span>
        </div>
      ),
      className: "w-[150px] text-right",
    },
    {
      id: "deals",
      header: t("table.deals"),
      accessor: (row) => (
        <div className="text-right">
          <span className="font-medium">{row.deals_closed}</span>
        </div>
      ),
      className: "w-[120px] text-right",
    },
    {
      id: "visits",
      header: t("table.visits"),
      accessor: (row) => (
        <div className="text-right">
          <span className="font-medium">{row.visits_completed}</span>
        </div>
      ),
      className: "w-[100px] text-right",
    },
    {
      id: "tasks",
      header: t("table.tasks"),
      accessor: (row) => (
        <div className="text-right">
          <span className="font-medium">{row.tasks_completed}</span>
        </div>
      ),
      className: "w-[120px] text-right",
    },
  ];

  return (
    <DataTable
      columns={columns}
      data={mappedData}
      isLoading={isLoading}
      pagination={pagination}
      onPageChange={onPageChange}
      onPerPageChange={onPerPageChange}
      emptyMessage={emptyMessage || t("table.empty")}
      itemName="sales performance"
      perPageOptions={[10, 20, 50, 100]}
    />
  );
}


