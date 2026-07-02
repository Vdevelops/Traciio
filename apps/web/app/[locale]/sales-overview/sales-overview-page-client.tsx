"use client";

import { useEffect, useMemo, useState } from "react";
import { useTranslations } from "next-intl";
import { useQuery } from "@tanstack/react-query";
import {
  Activity,
  AlertTriangle,
  ChartColumn,
  Search,
  TimerReset,
} from "lucide-react";
import { SalesPerformanceList } from "@/features/sales-overview/components/sales-performance-list";
import { SalesOverviewChart } from "@/features/sales-overview/components/SalesOverviewChart";
import { useMonthlySalesOverview } from "@/features/sales-overview/hooks/useMonthlySalesOverview";
import { useSalesPerformanceList } from "@/features/sales-overview/hooks/useSalesPerformanceList";
import { salesOverviewService } from "@/features/sales-overview/services/salesOverviewService";
import type {
  FunnelDiagnosticsData,
  ListProspectOutcomesRequest,
  SalesPerformanceListItem,
} from "@/features/sales-overview/types";
import { useDebounce } from "@/hooks/use-debounce";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Pagination,
  PaginationContent,
  PaginationFirst,
  PaginationItem,
  PaginationLast,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Skeleton } from "@/components/ui/skeleton";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { cn } from "@/lib/utils";
import type { DateRange } from "react-day-picker";
import { startOfYear, endOfYear, format } from "date-fns";

type OutcomeStatusFilter = "all" | "open" | "won" | "lost";
type TrendMode = "monthly" | "mom" | "rolling_30d" | "rolling_90d" | "qoq";

const perPageOptions = [10, 20, 50, 100];

const getInitials = (name?: string) => {
  if (!name) return "SR";
  return (
    name
      .split(" ")
      .filter(Boolean)
      .slice(0, 2)
      .map((part) => part[0]?.toUpperCase())
      .join("") || "SR"
  );
};

const formatOutcomeDate = (date?: string) => {
  if (!date) return "-";
  const parsed = new Date(date);
  if (Number.isNaN(parsed.getTime())) return "-";
  return format(parsed, "d MMM yyyy");
};

const formatDateTimeLabel = (date?: string) => {
  if (!date) return "-";
  const parsed = new Date(date);
  if (Number.isNaN(parsed.getTime())) return "-";
  return format(parsed, "d MMM yyyy, HH:mm");
};

const getPaginationPages = (currentPage: number, totalPages: number) => {
  if (totalPages <= 7) {
    return Array.from({ length: totalPages }, (_, index) => index + 1);
  }

  const pages = new Set([1, totalPages, currentPage]);
  if (currentPage > 1) pages.add(currentPage - 1);
  if (currentPage < totalPages) pages.add(currentPage + 1);
  return Array.from(pages).sort((a, b) => a - b);
};

export function SalesOverviewPageClient() {
  const t = useTranslations("salesOverview");

  // Filter State for Chart
  const [filterMode, setFilterMode] = useState<"year" | "range">("year");
  const [selectedYear, setSelectedYear] = useState<number>(
    new Date().getFullYear(),
  );
  const [selectedMetric, setSelectedMetric] = useState<
    "revenue" | "deals" | "visits" | "tasks"
  >("revenue");
  const [trendMode, setTrendMode] = useState<TrendMode>("monthly");
  const [diagnosticsSalesUserId, setDiagnosticsSalesUserId] = useState("all");
  const [diagnosticsStageId, setDiagnosticsStageId] = useState("all");
  const [dateRange, setDateRange] = useState<DateRange | undefined>(() => {
    const start = startOfYear(new Date());
    const end = endOfYear(new Date());
    return { from: start, to: end };
  });

  // Calculate generic start/end dates based on mode for Chart
  const { startDate: chartStartDate, endDate: chartEndDate } = useMemo(() => {
    if (trendMode === "rolling_30d" || trendMode === "rolling_90d") {
      const daySpan = trendMode === "rolling_30d" ? 29 : 89;
      const end = new Date();
      const start = new Date();
      start.setDate(end.getDate() - daySpan);
      return {
        startDate: format(start, "yyyy-MM-dd"),
        endDate: format(end, "yyyy-MM-dd"),
      };
    }

    if (filterMode === "year") {
      const start = new Date(selectedYear, 0, 1);
      const end = new Date(selectedYear, 11, 31);
      return {
        startDate: format(start, "yyyy-MM-dd"),
        endDate: format(end, "yyyy-MM-dd"),
      };
    } else {
      // Range mode
      if (!dateRange?.from) return { startDate: undefined, endDate: undefined };

      const start = format(dateRange.from, "yyyy-MM-dd");
      let end = undefined;
      if (dateRange.to) {
        end = format(dateRange.to, "yyyy-MM-dd");
      }
      return { startDate: start, endDate: end };
    }
  }, [filterMode, selectedYear, dateRange]);

  // Fetch Chart Data
  const { monthlyData, isLoading: isChartLoading } = useMonthlySalesOverview(
    chartStartDate,
    chartEndDate,
    trendMode,
  );

  const { data: diagnosticsResponse, isLoading: isDiagnosticsLoading } =
    useQuery({
      queryKey: [
        "sales-overview",
        "funnel-diagnostics",
        {
          sales_user_id: diagnosticsSalesUserId,
          stage_id: diagnosticsStageId,
        },
      ],
      queryFn: () =>
        salesOverviewService.getFunnelDiagnostics({
          sales_user_id:
            diagnosticsSalesUserId !== "all"
              ? diagnosticsSalesUserId
              : undefined,
          stage_id:
            diagnosticsStageId !== "all" ? diagnosticsStageId : undefined,
        }),
      staleTime: 30000,
    });
  const diagnostics = diagnosticsResponse?.data;

  // List Data Hook (Lifted state)
  const listProps = useSalesPerformanceList();
  const { setStartDate, setEndDate } = listProps;

  // Sync Chart Date Filter to List
  useEffect(() => {
    if (chartStartDate) setStartDate(chartStartDate);
    if (chartEndDate) setEndDate(chartEndDate);
  }, [chartStartDate, chartEndDate, setEndDate, setStartDate]);

  // Handlers
  const handleYearChange = (year: number) => {
    setSelectedYear(year);
  };

  const handleDateRangeChange = (range: DateRange | undefined) => {
    setDateRange(range);
  };
  return (
    <div className="space-y-8">
      {/* Page Header */}
      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div>
          <h1 className="text-3xl font-medium tracking-tight flex items-center gap-2">
            <ChartColumn className="h-8 w-8 text-primary" aria-hidden="true" />
            {t("title")}
          </h1>
          <p className="text-muted-foreground mt-1 text-sm">
            {t("description")}
          </p>
        </div>
      </div>

      {/* Chart Section */}
      <SalesOverviewChart
        data={monthlyData?.monthly_data ?? []}
        isLoading={isChartLoading}
        filterMode={filterMode}
        onFilterModeChange={setFilterMode}
        selectedYear={selectedYear}
        onYearChange={handleYearChange}
        dateRange={dateRange}
        onDateRangeChange={handleDateRangeChange}
        selectedMetric={selectedMetric}
        onMetricChange={setSelectedMetric}
        trendMode={trendMode}
        onTrendModeChange={setTrendMode}
      />

      {/* List Section */}
      <Tabs defaultValue="performance" className="space-y-4">
        <TabsList>
          <TabsTrigger value="performance" className="cursor-pointer">
            {t("tabs.performance")}
          </TabsTrigger>
          <TabsTrigger value="prospect-outcomes" className="cursor-pointer">
            {t("tabs.prospect_outcomes")}
          </TabsTrigger>
          <TabsTrigger value="diagnostics" className="cursor-pointer">
            {t("tabs.diagnostics")}
          </TabsTrigger>
        </TabsList>
        <TabsContent value="performance" className="mt-0">
          <SalesPerformanceList {...listProps} />
        </TabsContent>
        <TabsContent value="prospect-outcomes" className="mt-0">
          <ProspectOutcomesOverview
            startDate={chartStartDate}
            endDate={chartEndDate}
          />
        </TabsContent>
        <TabsContent value="diagnostics" className="mt-0">
          <FunnelDiagnosticsOverview
            data={diagnostics}
            isLoading={isDiagnosticsLoading}
            selectedSalesUserId={diagnosticsSalesUserId}
            selectedStageId={diagnosticsStageId}
            onSalesUserChange={setDiagnosticsSalesUserId}
            onStageChange={setDiagnosticsStageId}
          />
        </TabsContent>
      </Tabs>
    </div>
  );
}

function ProspectOutcomesOverview({
  startDate,
  endDate,
}: Readonly<{
  startDate?: string;
  endDate?: string;
}>) {
  const t = useTranslations("salesOverview");
  const [page, setPage] = useState(1);
  const [perPage, setPerPage] = useState(10);
  const [search, setSearch] = useState("");
  const [salesUserId, setSalesUserId] = useState("all");
  const [status, setStatus] = useState<OutcomeStatusFilter>("all");
  const debouncedSearch = useDebounce(search, 500);

  const { data: salesOptionsData } = useQuery({
    queryKey: [
      "sales-overview",
      "performance",
      "options",
      { startDate, endDate },
    ],
    queryFn: () =>
      salesOverviewService.listSalesPerformance({
        page: 1,
        per_page: 100,
        start_date: startDate,
        end_date: endDate,
        sort_by: "name",
        order: "asc",
      }),
    staleTime: 30000,
  });

  const queryParams = useMemo<ListProspectOutcomesRequest>(() => {
    const params: ListProspectOutcomesRequest = {
      page,
      per_page: perPage,
      start_date: startDate,
      end_date: endDate,
    };

    const trimmedSearch = debouncedSearch.trim();
    if (trimmedSearch) {
      params.search = trimmedSearch;
    }
    if (salesUserId !== "all") {
      params.sales_user_id = salesUserId;
    }
    if (status !== "all") {
      params.status = status;
    }

    return params;
  }, [debouncedSearch, endDate, page, perPage, salesUserId, startDate, status]);

  const { data, isLoading } = useQuery({
    queryKey: ["sales-overview", "prospect-outcomes", queryParams],
    queryFn: () => salesOverviewService.listProspectOutcomes(queryParams),
    staleTime: 30000,
  });

  const outcomes = Array.isArray(data?.data) ? data.data : [];
  const pagination = data?.meta?.pagination;
  const salesOptions = Array.isArray(salesOptionsData?.data)
    ? salesOptionsData.data
    : [];
  const totalPages = pagination?.total_pages ?? 0;
  const pageNumbers = getPaginationPages(pagination?.page ?? page, totalPages);
  const showingStart =
    pagination && pagination.total > 0
      ? (pagination.page - 1) * pagination.per_page + 1
      : 0;
  const showingEnd = pagination
    ? Math.min(pagination.page * pagination.per_page, pagination.total)
    : 0;

  const handleSalesChange = (value: string) => {
    setSalesUserId(value);
    setPage(1);
  };

  const handleStatusChange = (value: string) => {
    setStatus(value as OutcomeStatusFilter);
    setPage(1);
  };

  const statusBadgeVariant = (value: string) => {
    if (value === "won") return "active";
    if (value === "lost") return "destructive";
    return "outline";
  };

  const formatReason = (reason?: string) => {
    if (!reason) return "-";
    if (reason === "no_reason_provided") {
      return t("prospects.no_reason_provided");
    }
    return reason;
  };

  const formatReasonCategory = (category?: string) => {
    if (!category) return null;
    return t(`prospects.reason_categories.${category}`);
  };

  const formatStatus = (value: string) => {
    if (value === "won" || value === "lost" || value === "open") {
      return t(`prospects.status.${value}`);
    }
    return value;
  };

  return (
    <Card className="gap-0 overflow-hidden py-0">
      <CardHeader className="border-b px-5 py-4">
        <div className="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
          <div>
            <CardTitle className="text-base">
              {t("prospects.overview_title")}
            </CardTitle>
            <p className="mt-1 text-sm text-muted-foreground">
              {t("prospects.overview_description")}
            </p>
          </div>
          <div className="flex flex-col gap-3 sm:flex-row sm:items-center">
            <div className="relative w-full sm:w-[280px]">
              <Search
                className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground"
                aria-hidden="true"
              />
              <Input
                className="h-9 pl-9"
                value={search}
                placeholder={t("prospects.filters.search_placeholder")}
                onChange={(event) => {
                  setSearch(event.target.value);
                  setPage(1);
                }}
              />
            </div>
            <Select value={salesUserId} onValueChange={handleSalesChange}>
              <SelectTrigger className="h-9 w-full sm:w-[220px]">
                <SelectValue placeholder={t("prospects.filters.sales")} />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">
                  {t("prospects.filters.all_sales")}
                </SelectItem>
                {salesOptions.map((sales: SalesPerformanceListItem) => (
                  <SelectItem key={sales.user_id} value={sales.user_id}>
                    {sales.user_name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <Select value={status} onValueChange={handleStatusChange}>
              <SelectTrigger className="h-9 w-full sm:w-[180px]">
                <SelectValue placeholder={t("prospects.filters.status")} />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">
                  {t("prospects.filters.all_status")}
                </SelectItem>
                <SelectItem value="won">{t("prospects.status.won")}</SelectItem>
                <SelectItem value="lost">
                  {t("prospects.status.lost")}
                </SelectItem>
                <SelectItem value="open">
                  {t("prospects.status.open")}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>
      </CardHeader>
      <CardContent className="p-0">
        <div className="overflow-x-auto">
          <Table>
            <TableHeader>
              <TableRow className="bg-muted/40 hover:bg-muted/40">
                <TableHead className="px-5 text-xs uppercase text-muted-foreground">
                  {t("prospects.table.prospect")}
                </TableHead>
                <TableHead className="text-xs uppercase text-muted-foreground">
                  {t("prospects.table.sales")}
                </TableHead>
                <TableHead className="text-xs uppercase text-muted-foreground">
                  {t("prospects.table.account")}
                </TableHead>
                <TableHead className="text-xs uppercase text-muted-foreground">
                  {t("prospects.table.status")}
                </TableHead>
                <TableHead className="text-right text-xs uppercase text-muted-foreground">
                  {t("prospects.table.value")}
                </TableHead>
                <TableHead className="text-xs uppercase text-muted-foreground">
                  {t("prospects.table.reason")}
                </TableHead>
                <TableHead className="text-xs uppercase text-muted-foreground">
                  {t("prospects.table.closed_at")}
                </TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {isLoading ? (
                Array.from({ length: 5 }).map((_, index) => (
                  <TableRow key={index}>
                    <TableCell className="px-5">
                      <Skeleton className="h-4 w-44" />
                    </TableCell>
                    <TableCell>
                      <Skeleton className="h-8 w-36" />
                    </TableCell>
                    <TableCell>
                      <Skeleton className="h-4 w-32" />
                    </TableCell>
                    <TableCell>
                      <Skeleton className="h-6 w-16 rounded-full" />
                    </TableCell>
                    <TableCell>
                      <Skeleton className="ml-auto h-4 w-24" />
                    </TableCell>
                    <TableCell>
                      <Skeleton className="h-4 w-48" />
                    </TableCell>
                    <TableCell>
                      <Skeleton className="h-4 w-24" />
                    </TableCell>
                  </TableRow>
                ))
              ) : outcomes.length === 0 ? (
                <TableRow>
                  <TableCell
                    colSpan={7}
                    className="h-28 px-5 text-center text-sm text-muted-foreground"
                  >
                    {t("prospects.no_recent")}
                  </TableCell>
                </TableRow>
              ) : (
                outcomes.map((outcome) => (
                  <TableRow key={outcome.id}>
                    <TableCell className="px-5 font-medium">
                      <div className="max-w-[260px] truncate">
                        {outcome.title}
                      </div>
                    </TableCell>
                    <TableCell>
                      <div className="flex min-w-[180px] items-center gap-3">
                        <Avatar className="h-8 w-8">
                          <AvatarImage
                            src={outcome.sales_rep_avatar_url || undefined}
                            alt={
                              outcome.sales_rep_name ||
                              t("prospects.filters.sales")
                            }
                          />
                          <AvatarFallback>
                            {getInitials(outcome.sales_rep_name)}
                          </AvatarFallback>
                        </Avatar>
                        <div className="min-w-0">
                          <div className="truncate text-sm font-medium">
                            {outcome.sales_rep_name || "-"}
                          </div>
                          <div className="truncate text-xs text-muted-foreground">
                            {outcome.sales_rep_email || "-"}
                          </div>
                        </div>
                      </div>
                    </TableCell>
                    <TableCell className="text-muted-foreground">
                      {outcome.account_name || "-"}
                    </TableCell>
                    <TableCell>
                      <Badge variant={statusBadgeVariant(outcome.status)}>
                        {formatStatus(outcome.status)}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-right font-medium">
                      {outcome.value_formatted}
                    </TableCell>
                    <TableCell className="max-w-[280px] truncate text-muted-foreground">
                      <div className="space-y-1">
                        <div className="truncate">
                          {formatReason(outcome.reason)}
                        </div>
                        {outcome.reason_category ? (
                          <div className="text-xs uppercase tracking-wide text-muted-foreground/80">
                            {formatReasonCategory(outcome.reason_category)}
                          </div>
                        ) : null}
                      </div>
                    </TableCell>
                    <TableCell className="text-muted-foreground">
                      {formatOutcomeDate(
                        outcome.closed_at ?? outcome.created_at,
                      )}
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </div>
        {pagination && (
          <div className="flex flex-col gap-4 border-t px-5 py-4 lg:flex-row lg:items-center lg:justify-between">
            <div className="flex items-center gap-3">
              <Label
                htmlFor="prospect-outcomes-rows"
                className="text-sm font-medium text-muted-foreground"
              >
                {t("prospects.pagination.rows_per_page")}
              </Label>
              <Select
                value={String(pagination.per_page)}
                onValueChange={(value) => {
                  setPerPage(Number(value));
                  setPage(1);
                }}
              >
                <SelectTrigger
                  id="prospect-outcomes-rows"
                  className="h-8 w-fit"
                >
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {perPageOptions.map((option) => (
                    <SelectItem key={option} value={String(option)}>
                      {option}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="text-sm text-muted-foreground">
              {t("prospects.pagination.showing", {
                start: showingStart,
                end: showingEnd,
                total: pagination.total,
              })}
            </div>

            {totalPages > 1 && (
              <Pagination className="w-fit">
                <PaginationContent>
                  <PaginationItem>
                    <PaginationFirst
                      onClick={() => setPage(1)}
                      className={cn(
                        "h-8 w-8 cursor-pointer",
                        (!pagination.has_prev || isLoading) &&
                          "pointer-events-none opacity-50",
                      )}
                    />
                  </PaginationItem>
                  <PaginationItem>
                    <PaginationPrevious
                      onClick={() => setPage(Math.max(1, pagination.page - 1))}
                      className={cn(
                        "h-8 w-8 cursor-pointer",
                        (!pagination.has_prev || isLoading) &&
                          "pointer-events-none opacity-50",
                      )}
                    />
                  </PaginationItem>
                  {pageNumbers.map((pageNumber) => (
                    <PaginationItem key={pageNumber}>
                      <PaginationLink
                        isActive={pageNumber === pagination.page}
                        onClick={() => setPage(pageNumber)}
                        className={cn(
                          "h-8 w-8 cursor-pointer",
                          isLoading && "pointer-events-none opacity-50",
                        )}
                      >
                        {pageNumber}
                      </PaginationLink>
                    </PaginationItem>
                  ))}
                  <PaginationItem>
                    <PaginationNext
                      onClick={() =>
                        setPage(Math.min(totalPages, pagination.page + 1))
                      }
                      className={cn(
                        "h-8 w-8 cursor-pointer",
                        (!pagination.has_next || isLoading) &&
                          "pointer-events-none opacity-50",
                      )}
                    />
                  </PaginationItem>
                  <PaginationItem>
                    <PaginationLast
                      onClick={() => setPage(totalPages)}
                      className={cn(
                        "h-8 w-8 cursor-pointer",
                        (!pagination.has_next || isLoading) &&
                          "pointer-events-none opacity-50",
                      )}
                    />
                  </PaginationItem>
                </PaginationContent>
              </Pagination>
            )}
          </div>
        )}
      </CardContent>
    </Card>
  );
}

function FunnelDiagnosticsOverview({
  data,
  isLoading,
  selectedSalesUserId,
  selectedStageId,
  onSalesUserChange,
  onStageChange,
}: Readonly<{
  data?: FunnelDiagnosticsData;
  isLoading: boolean;
  selectedSalesUserId: string;
  selectedStageId: string;
  onSalesUserChange: (value: string) => void;
  onStageChange: (value: string) => void;
}>) {
  const t = useTranslations("salesOverview");

  const prioritizedDeals = useMemo(() => {
    if (!data) return [];

    const dealMap = new Map<
      string,
      {
        id: string;
        title: string;
        accountName?: string;
        assignedToName?: string;
        stageName: string;
        value: number;
        valueFormatted: string;
        probability: number;
        daysInStage?: number;
        daysWithoutActivity?: number;
        flags: string[];
      }
    >();

    data.stalled_deals.forEach((deal) => {
      const existing = dealMap.get(deal.id);
      if (existing) {
        existing.daysInStage = Math.max(
          existing.daysInStage ?? 0,
          deal.days_in_stage,
        );
        existing.flags = Array.from(new Set([...existing.flags, "stalled"]));
        return;
      }

      dealMap.set(deal.id, {
        id: deal.id,
        title: deal.title,
        accountName: deal.account_name,
        assignedToName: deal.assigned_to_name,
        stageName: deal.stage_name,
        value: deal.value,
        valueFormatted: deal.value_formatted,
        probability: deal.probability,
        daysInStage: deal.days_in_stage,
        flags: ["stalled"],
      });
    });

    data.no_activity_deals.forEach((deal) => {
      const existing = dealMap.get(deal.id);
      if (existing) {
        existing.daysWithoutActivity = Math.max(
          existing.daysWithoutActivity ?? 0,
          deal.days_without_activity,
        );
        existing.flags = Array.from(
          new Set([...existing.flags, "no_activity"]),
        );
        return;
      }

      dealMap.set(deal.id, {
        id: deal.id,
        title: deal.title,
        accountName: deal.account_name,
        assignedToName: deal.assigned_to_name,
        stageName: deal.stage_name,
        value: deal.value,
        valueFormatted: deal.value_formatted,
        probability: deal.probability,
        daysWithoutActivity: deal.days_without_activity,
        flags: ["no_activity"],
      });
    });

    return Array.from(dealMap.values()).sort((a, b) => {
      if (b.flags.length !== a.flags.length) {
        return b.flags.length - a.flags.length;
      }

      if ((b.daysWithoutActivity ?? 0) !== (a.daysWithoutActivity ?? 0)) {
        return (b.daysWithoutActivity ?? 0) - (a.daysWithoutActivity ?? 0);
      }

      if ((b.daysInStage ?? 0) !== (a.daysInStage ?? 0)) {
        return (b.daysInStage ?? 0) - (a.daysInStage ?? 0);
      }

      return b.value - a.value;
    });
  }, [data]);

  const atRiskValue = useMemo(
    () => prioritizedDeals.reduce((total, deal) => total + deal.value, 0),
    [prioritizedDeals],
  );

  const highestRiskDeal = prioritizedDeals[0];
  const dualRiskCount = prioritizedDeals.filter(
    (deal) => deal.flags.length > 1,
  ).length;
  const slowestTransition = useMemo(() => {
    if (!data?.stage_aging.length) return undefined;
    return [...data.stage_aging].sort((a, b) => {
      if (b.average_days !== a.average_days) {
        return b.average_days - a.average_days;
      }
      return b.transitions - a.transitions;
    })[0];
  }, [data]);

  if (isLoading) {
    return (
      <div className="grid gap-4 lg:grid-cols-3">
        {Array.from({ length: 3 }).map((_, index) => (
          <Card key={index}>
            <CardHeader>
              <Skeleton className="h-5 w-32" />
            </CardHeader>
            <CardContent>
              <Skeleton className="h-40 w-full" />
            </CardContent>
          </Card>
        ))}
      </div>
    );
  }

  if (!data) {
    return (
      <Card>
        <CardContent className="py-12 text-center text-sm text-muted-foreground">
          {t("diagnostics.empty")}
        </CardContent>
      </Card>
    );
  }

  return (
    <div className="space-y-4">
      <Card>
        <CardHeader className="pb-3">
          <CardTitle className="text-base">
            {t("diagnostics.filters.title")}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid gap-3 md:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="diagnostics-sales-filter">
                {t("diagnostics.filters.sales")}
              </Label>
              <Select
                value={selectedSalesUserId}
                onValueChange={onSalesUserChange}
              >
                <SelectTrigger id="diagnostics-sales-filter" className="h-9">
                  <SelectValue
                    placeholder={t("diagnostics.filters.all_sales")}
                  />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">
                    {t("diagnostics.filters.all_sales")}
                  </SelectItem>
                  {data.available_sales_reps.map((salesRep) => (
                    <SelectItem key={salesRep.id} value={salesRep.id}>
                      {salesRep.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <Label htmlFor="diagnostics-stage-filter">
                {t("diagnostics.filters.stage")}
              </Label>
              <Select value={selectedStageId} onValueChange={onStageChange}>
                <SelectTrigger id="diagnostics-stage-filter" className="h-9">
                  <SelectValue
                    placeholder={t("diagnostics.filters.all_stages")}
                  />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">
                    {t("diagnostics.filters.all_stages")}
                  </SelectItem>
                  {data.available_stages.map((stage) => (
                    <SelectItem key={stage.id} value={stage.id}>
                      {stage.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>
        </CardContent>
      </Card>

      <div className="grid gap-4 xl:grid-cols-3">
        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="flex items-center gap-2 text-base">
              <TimerReset className="h-4 w-4 text-amber-600" />
              {t("diagnostics.summary.priority")}
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-2">
            <div className="text-3xl font-semibold">
              {prioritizedDeals.length}
            </div>
            <p className="text-sm text-muted-foreground">
              {t("diagnostics.summary.priority_desc")}
            </p>
            {highestRiskDeal ? (
              <div className="rounded-lg border border-amber-200 bg-amber-50/70 p-3 text-sm text-amber-950 dark:border-amber-900/60 dark:bg-amber-950/20 dark:text-amber-100">
                <div className="font-medium">
                  {t("diagnostics.summary.next_focus")}
                </div>
                <div className="mt-1">
                  {highestRiskDeal.title} · {highestRiskDeal.stageName}
                </div>
                <div className="text-xs text-amber-900/80 dark:text-amber-200/80">
                  {highestRiskDeal.accountName ||
                    highestRiskDeal.assignedToName ||
                    "-"}
                </div>
              </div>
            ) : (
              <div className="text-sm text-muted-foreground">
                {t("diagnostics.summary.no_priority")}
              </div>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="flex items-center gap-2 text-base">
              <Activity className="h-4 w-4 text-red-600" />
              {t("diagnostics.summary.value_at_risk")}
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-2">
            <div className="text-3xl font-semibold">
              {new Intl.NumberFormat("id-ID", {
                style: "currency",
                currency: "IDR",
                maximumFractionDigits: 0,
              }).format(atRiskValue)}
            </div>
            <p className="text-sm text-muted-foreground">
              {t("diagnostics.summary.value_at_risk_desc")}
            </p>
            <div className="flex items-center justify-between rounded-lg border p-3 text-sm">
              <span className="text-muted-foreground">
                {t("diagnostics.summary.dual_risk")}
              </span>
              <span className="font-medium">{dualRiskCount}</span>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="flex items-center gap-2 text-base">
              <AlertTriangle className="h-4 w-4 text-primary" />
              {t("diagnostics.summary.bottleneck")}
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-2">
            {slowestTransition ? (
              <>
                <div className="text-xl font-semibold">
                  {slowestTransition.from_stage_name} →{" "}
                  {slowestTransition.to_stage_name}
                </div>
                <p className="text-sm text-muted-foreground">
                  {t("diagnostics.summary.bottleneck_desc", {
                    days: slowestTransition.average_days.toFixed(1),
                    transitions: slowestTransition.transitions,
                  })}
                </p>
                <div className="flex items-center justify-between rounded-lg border p-3 text-sm">
                  <span className="text-muted-foreground">
                    {t("diagnostics.stage_aging.median_days")}
                  </span>
                  <span className="font-medium">
                    {slowestTransition.median_days.toFixed(1)}
                  </span>
                </div>
              </>
            ) : (
              <>
                <div className="text-3xl font-semibold">
                  {data.summary.stage_aging_transitions}
                </div>
                <p className="text-sm text-muted-foreground">
                  {t("diagnostics.summary.stage_aging_desc")}
                </p>
              </>
            )}
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="text-base">
            {t("diagnostics.action_plan.title")}
          </CardTitle>
          <p className="text-sm text-muted-foreground">
            {t("diagnostics.action_plan.description", {
              generated_at: formatDateTimeLabel(data.generated_at),
            })}
          </p>
        </CardHeader>
        <CardContent className="grid gap-3 md:grid-cols-3">
          <div className="rounded-lg border p-4">
            <div className="text-sm font-medium">
              {t("diagnostics.action_plan.step_1_title")}
            </div>
            <p className="mt-1 text-sm text-muted-foreground">
              {prioritizedDeals.length > 0
                ? t("diagnostics.action_plan.step_1_value", {
                    count: prioritizedDeals.length,
                  })
                : t("diagnostics.action_plan.step_1_empty")}
            </p>
          </div>
          <div className="rounded-lg border p-4">
            <div className="text-sm font-medium">
              {t("diagnostics.action_plan.step_2_title")}
            </div>
            <p className="mt-1 text-sm text-muted-foreground">
              {data.summary.stalled_deals > 0
                ? t("diagnostics.action_plan.step_2_value", {
                    count: data.summary.stalled_deals,
                    days: data.stalled_threshold_days,
                  })
                : t("diagnostics.action_plan.step_2_empty")}
            </p>
          </div>
          <div className="rounded-lg border p-4">
            <div className="text-sm font-medium">
              {t("diagnostics.action_plan.step_3_title")}
            </div>
            <p className="mt-1 text-sm text-muted-foreground">
              {slowestTransition
                ? t("diagnostics.action_plan.step_3_value", {
                    from: slowestTransition.from_stage_name,
                    to: slowestTransition.to_stage_name,
                    days: slowestTransition.average_days.toFixed(1),
                  })
                : t("diagnostics.action_plan.step_3_empty")}
            </p>
          </div>
        </CardContent>
      </Card>

      <div className="grid gap-4 xl:grid-cols-5">
        <Card className="xl:col-span-3">
          <CardHeader>
            <CardTitle className="text-base">
              {t("diagnostics.priority_table.title")}
            </CardTitle>
            <p className="text-sm text-muted-foreground">
              {t("diagnostics.priority_table.description")}
            </p>
          </CardHeader>
          <CardContent>
            {prioritizedDeals.length === 0 ? (
              <p className="text-sm text-muted-foreground">
                {t("diagnostics.priority_table.empty")}
              </p>
            ) : (
              <div className="overflow-x-auto">
                <Table>
                  <TableHeader>
                    <TableRow className="bg-muted/40 hover:bg-muted/40">
                      <TableHead>
                        {t("diagnostics.priority_table.deal")}
                      </TableHead>
                      <TableHead>
                        {t("diagnostics.priority_table.owner")}
                      </TableHead>
                      <TableHead>
                        {t("diagnostics.priority_table.risk")}
                      </TableHead>
                      <TableHead className="text-right">
                        {t("diagnostics.common.value")}
                      </TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {prioritizedDeals.map((deal) => (
                      <TableRow key={deal.id}>
                        <TableCell className="min-w-[260px]">
                          <div className="font-medium">{deal.title}</div>
                          <div className="text-xs text-muted-foreground">
                            {deal.accountName || "-"} · {deal.stageName}
                          </div>
                        </TableCell>
                        <TableCell className="text-sm text-muted-foreground">
                          {deal.assignedToName || "-"}
                        </TableCell>
                        <TableCell>
                          <div className="flex flex-wrap gap-2">
                            {deal.daysInStage ? (
                              <Badge
                                variant="outline"
                                className="border-amber-200 text-amber-700"
                              >
                                {t("diagnostics.stalled.days_in_stage")}{" "}
                                {deal.daysInStage}
                              </Badge>
                            ) : null}
                            {deal.daysWithoutActivity ? (
                              <Badge
                                variant="outline"
                                className="border-red-200 text-red-700"
                              >
                                {t(
                                  "diagnostics.no_activity.days_without_activity",
                                )}{" "}
                                {deal.daysWithoutActivity}
                              </Badge>
                            ) : null}
                            <Badge variant="secondary">
                              {t("diagnostics.priority_table.probability", {
                                value: deal.probability,
                              })}
                            </Badge>
                          </div>
                        </TableCell>
                        <TableCell className="text-right font-medium">
                          {deal.valueFormatted}
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>
            )}
          </CardContent>
        </Card>

        <Card className="xl:col-span-2">
          <CardHeader>
            <CardTitle className="text-base">
              {t("diagnostics.stage_aging.title")}
            </CardTitle>
            <p className="text-sm text-muted-foreground">
              {t("diagnostics.stage_aging.description")}
            </p>
          </CardHeader>
          <CardContent className="space-y-3">
            {data.stage_aging.length === 0 ? (
              <p className="text-sm text-muted-foreground">
                {t("diagnostics.stage_aging.empty")}
              </p>
            ) : (
              data.stage_aging.map((item) => (
                <div key={item.transition_key} className="rounded-lg border p-3">
                  <div className="flex items-start justify-between gap-3">
                    <div>
                      <div className="font-medium">
                        {item.from_stage_name} → {item.to_stage_name}
                      </div>
                      <div className="text-xs text-muted-foreground">
                        {t("diagnostics.stage_aging.transitions")}{" "}
                        {item.transitions}
                      </div>
                    </div>
                    <Badge variant="outline">
                      {item.average_days.toFixed(1)}{" "}
                      {t("diagnostics.stage_aging.day_unit")}
                    </Badge>
                  </div>
                  <div className="mt-3 flex items-center justify-between text-sm">
                    <span className="text-muted-foreground">
                      {t("diagnostics.stage_aging.avg_days")}
                    </span>
                    <span className="font-medium">
                      {item.average_days.toFixed(1)}
                    </span>
                  </div>
                  <div className="mt-1 flex items-center justify-between text-sm">
                    <span className="text-muted-foreground">
                      {t("diagnostics.stage_aging.median_days")}
                    </span>
                    <span className="font-medium">
                      {item.median_days.toFixed(1)}
                    </span>
                  </div>
                </div>
              ))
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
