"use client";

import { Suspense, useMemo, useState } from "react";
import { useLocale, useTranslations } from "next-intl";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import {
  MapPin,
  Package,
  Building2,
  BarChart3,
  Target,
  Trophy,
  TrendingDown,
  Clock3,
} from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
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
import { Label } from "@/components/ui/label";
import dynamic from "next/dynamic";
import type {
  ProspectOutcomeSummary,
  ProspectReasonBreakdown,
  SalesRepCheckInLocation,
} from "../types";

// Lazy load components with code splitting
const SalesRepCheckInMap = dynamic(
  () =>
    import("./SalesRepCheckInMap").then((mod) => ({
      default: mod.SalesRepCheckInMap,
    })),
  {
    loading: () => <Skeleton className="h-[500px] w-full" />,
    ssr: false,
  },
);

const SalesRepProductSales = dynamic(
  () =>
    import("./SalesRepProductSales").then((mod) => ({
      default: mod.SalesRepProductSales,
    })),
  {
    loading: () => (
      <div className="space-y-4">
        <Skeleton className="h-10 w-full" />
        <Skeleton className="h-[400px] w-full" />
      </div>
    ),
    ssr: false,
  },
);

const SalesRepCustomers = dynamic(
  () =>
    import("./SalesRepCustomers").then((mod) => ({
      default: mod.SalesRepCustomers,
    })),
  {
    loading: () => (
      <div className="space-y-4">
        <Skeleton className="h-10 w-full" />
        <Skeleton className="h-[400px] w-full" />
      </div>
    ),
    ssr: false,
  },
);

interface SalesRepDetailTabsProps {
  readonly userId: string;
  readonly startDate?: string;
  readonly endDate?: string;
  readonly prospectOutcome?: ProspectOutcomeSummary;
  readonly checkInLocationsProps: {
    readonly locations: readonly SalesRepCheckInLocation[];
    readonly isLoading?: boolean;
    readonly totalVisits?: number;
    readonly page?: number;
    readonly perPage?: number;
    readonly onPageChange?: (page: number) => void;
    readonly onPerPageChange?: (perPage: number) => void;
  };
}

export function SalesRepDetailTabs({
  userId,
  startDate,
  endDate,
  prospectOutcome,
  checkInLocationsProps,
}: SalesRepDetailTabsProps) {
  const t = useTranslations("salesOverview");
  const locale = useLocale();
  const [prospectPage, setProspectPage] = useState(1);
  const [prospectPerPage, setProspectPerPage] = useState(5);

  const formatReason = (reason?: string) => {
    if (!reason) {
      return "-";
    }
    if (reason === "no_reason_provided") {
      return t("prospects.no_reason_provided");
    }
    return reason;
  };

  const formatStatus = (status: string) => {
    if (status === "won") {
      return t("prospects.status.won");
    }
    if (status === "lost") {
      return t("prospects.status.lost");
    }
    if (status === "open") {
      return t("prospects.status.open");
    }
    return status;
  };

  const getStatusClassName = (status: string) => {
    if (status === "won") {
      return "border-emerald-200 bg-emerald-50 text-emerald-700";
    }
    if (status === "lost") {
      return "border-rose-200 bg-rose-50 text-rose-700";
    }
    return "border-amber-200 bg-amber-50 text-amber-700";
  };

  const formatDate = (value?: string) => {
    if (!value) {
      return "-";
    }
    return new Intl.DateTimeFormat(locale, {
      day: "2-digit",
      month: "short",
      year: "numeric",
    }).format(new Date(value));
  };

  const renderReasonBreakdown = (
    title: string,
    reasons: readonly ProspectReasonBreakdown[],
  ) => (
    <Card className="gap-4 py-5">
      <CardHeader className="px-5">
        <CardTitle className="text-sm">{title}</CardTitle>
      </CardHeader>
      <CardContent className="px-5">
        {reasons.length === 0 ? (
          <p className="text-sm text-muted-foreground">
            {t("prospects.no_reasons")}
          </p>
        ) : (
          <div className="space-y-4">
            {reasons.map((reason) => (
              <div key={reason.reason} className="space-y-2">
                <div className="flex items-center justify-between gap-3 text-sm">
                  <span className="min-w-0 truncate font-medium">
                    {formatReason(reason.reason)}
                  </span>
                  <span className="text-muted-foreground">
                    {reason.count} ({reason.percentage.toFixed(1)}%)
                  </span>
                </div>
                <Progress value={Math.min(reason.percentage, 100)} />
              </div>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );

  const prospectMetrics = [
    {
      label: t("prospects.total"),
      value: prospectOutcome?.total_prospects ?? 0,
      icon: Target,
      className: "text-foreground",
    },
    {
      label: t("prospects.won"),
      value: prospectOutcome?.won_prospects ?? 0,
      icon: Trophy,
      className: "text-emerald-600",
    },
    {
      label: t("prospects.lost"),
      value: prospectOutcome?.lost_prospects ?? 0,
      icon: TrendingDown,
      className: "text-rose-600",
    },
    {
      label: t("prospects.open"),
      value: prospectOutcome?.open_prospects ?? 0,
      icon: Clock3,
      className: "text-amber-600",
    },
  ];

  const prospectConversionRate = prospectOutcome?.prospect_conversion_rate ?? 0;
  const recentProspects = useMemo(
    () => prospectOutcome?.recent_prospects ?? [],
    [prospectOutcome?.recent_prospects],
  );
  const prospectTotalPages = Math.max(
    1,
    Math.ceil(recentProspects.length / prospectPerPage),
  );
  const currentProspectPage = Math.min(prospectPage, prospectTotalPages);
  const paginatedProspects = useMemo(() => {
    const start = (currentProspectPage - 1) * prospectPerPage;
    return recentProspects.slice(start, start + prospectPerPage);
  }, [currentProspectPage, prospectPerPage, recentProspects]);
  const prospectStart = recentProspects.length
    ? (currentProspectPage - 1) * prospectPerPage + 1
    : 0;
  const prospectEnd = Math.min(
    currentProspectPage * prospectPerPage,
    recentProspects.length,
  );

  return (
    <Tabs defaultValue="locations" className="w-full">
      <TabsList>
        <TabsTrigger value="locations" className="gap-2 cursor-pointer">
          <MapPin className="h-4 w-4" />
          {t("check_in_locations")}
        </TabsTrigger>
        <TabsTrigger value="products" className="gap-2 cursor-pointer">
          <Package className="h-4 w-4" />
          {t("products_sold")}
        </TabsTrigger>
        <TabsTrigger value="customers" className="gap-2 cursor-pointer">
          <Building2 className="h-4 w-4" />
          {t("customers.title")}
        </TabsTrigger>
        <TabsTrigger value="prospects" className="gap-2 cursor-pointer">
          <BarChart3 className="h-4 w-4" />
          {t("prospects.title")}
        </TabsTrigger>
      </TabsList>

      <TabsContent value="locations" className="mt-6">
        <Suspense fallback={<Skeleton className="h-[500px] w-full" />}>
          <SalesRepCheckInMap {...checkInLocationsProps} />
        </Suspense>
      </TabsContent>

      <TabsContent value="products" className="mt-6">
        <Suspense
          fallback={
            <div className="space-y-4">
              <Skeleton className="h-10 w-full" />
              <Skeleton className="h-[400px] w-full" />
            </div>
          }
        >
          <SalesRepProductSales
            userId={userId}
            startDate={startDate}
            endDate={endDate}
          />
        </Suspense>
      </TabsContent>

      <TabsContent value="customers" className="mt-6">
        <Suspense
          fallback={
            <div className="space-y-4">
              <Skeleton className="h-10 w-full" />
              <Skeleton className="h-[400px] w-full" />
            </div>
          }
        >
          <SalesRepCustomers
            userId={userId}
            startDate={startDate}
            endDate={endDate}
          />
        </Suspense>
      </TabsContent>

      <TabsContent value="prospects" className="mt-6">
        <div className="space-y-6">
          <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
            {prospectMetrics.map((metric) => {
              const Icon = metric.icon;

              return (
                <Card key={metric.label} className="gap-3 py-5">
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 px-5">
                    <CardTitle className="text-xs font-medium uppercase text-muted-foreground">
                      {metric.label}
                    </CardTitle>
                    <Icon className="h-4 w-4 text-muted-foreground" />
                  </CardHeader>
                  <CardContent className="px-5">
                    <div
                      className={`text-3xl font-semibold ${metric.className}`}
                    >
                      {metric.value}
                    </div>
                  </CardContent>
                </Card>
              );
            })}
          </div>

          <Card className="gap-4 py-5">
            <CardHeader className="flex flex-row items-start justify-between gap-4 px-5">
              <div className="space-y-1">
                <CardTitle className="text-sm">
                  {t("prospects.conversion")}
                </CardTitle>
                <p className="text-sm text-muted-foreground">
                  {t("prospects.from_deals")}
                </p>
              </div>
              <Badge variant="active">
                {prospectConversionRate.toFixed(1)}%
              </Badge>
            </CardHeader>
            <CardContent className="px-5">
              <Progress value={Math.min(prospectConversionRate, 100)} />
            </CardContent>
          </Card>

          <div className="grid gap-4 lg:grid-cols-2">
            {renderReasonBreakdown(
              t("prospects.won_reasons"),
              prospectOutcome?.won_reasons ?? [],
            )}
            {renderReasonBreakdown(
              t("prospects.lost_reasons"),
              prospectOutcome?.lost_reasons ?? [],
            )}
          </div>

          <Card className="gap-0 overflow-hidden py-0">
            <CardHeader className="border-b px-5 py-4">
              <CardTitle className="text-sm">{t("prospects.recent")}</CardTitle>
            </CardHeader>
            {recentProspects.length === 0 ? (
              <CardContent className="px-5 py-8 text-center text-sm text-muted-foreground">
                {t("prospects.no_recent")}
              </CardContent>
            ) : (
              <>
                <Table>
                  <TableHeader>
                    <TableRow className="bg-muted/40 hover:bg-muted/40">
                      <TableHead className="px-5 text-xs uppercase text-muted-foreground">
                        {t("prospects.table.prospect")}
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
                    {paginatedProspects.map((prospect) => (
                      <TableRow key={prospect.id}>
                        <TableCell className="px-5 font-medium">
                          {prospect.title}
                        </TableCell>
                        <TableCell className="text-muted-foreground">
                          {prospect.account_name || "-"}
                        </TableCell>
                        <TableCell>
                          <Badge
                            variant="outline"
                            className={getStatusClassName(prospect.status)}
                          >
                            {formatStatus(prospect.status)}
                          </Badge>
                        </TableCell>
                        <TableCell className="text-right font-medium">
                          {prospect.value_formatted}
                        </TableCell>
                        <TableCell className="max-w-[280px] truncate text-muted-foreground">
                          {formatReason(prospect.reason)}
                        </TableCell>
                        <TableCell className="text-muted-foreground">
                          {formatDate(
                            prospect.closed_at ?? prospect.created_at,
                          )}
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>

                <CardContent className="flex flex-col gap-4 border-t bg-muted/20 px-5 py-4 md:flex-row md:items-center md:justify-between">
                  <div className="flex flex-wrap items-center gap-3">
                    <Label
                      htmlFor="prospect-rows-per-page"
                      className="text-sm text-muted-foreground"
                    >
                      {t("prospects.pagination.rows_per_page")}
                    </Label>
                    <Select
                      value={String(prospectPerPage)}
                      onValueChange={(value) => {
                        setProspectPerPage(Number(value));
                        setProspectPage(1);
                      }}
                    >
                      <SelectTrigger
                        id="prospect-rows-per-page"
                        className="h-9 w-20"
                      >
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        {[5, 10, 20].map((option) => (
                          <SelectItem key={option} value={String(option)}>
                            {option}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    <span className="text-sm text-muted-foreground">
                      {t("prospects.pagination.showing", {
                        start: prospectStart,
                        end: prospectEnd,
                        total: recentProspects.length,
                      })}
                    </span>
                  </div>

                  <Pagination className="mx-0 w-auto justify-start md:justify-end">
                    <PaginationContent>
                      <PaginationItem>
                        <PaginationFirst
                          onClick={() => setProspectPage(1)}
                          disabled={currentProspectPage <= 1}
                          aria-disabled={currentProspectPage <= 1}
                        />
                      </PaginationItem>
                      <PaginationItem>
                        <PaginationPrevious
                          onClick={() =>
                            setProspectPage(
                              Math.max(1, currentProspectPage - 1),
                            )
                          }
                          disabled={currentProspectPage <= 1}
                          aria-disabled={currentProspectPage <= 1}
                        />
                      </PaginationItem>
                      {Array.from(
                        { length: prospectTotalPages },
                        (_, index) => index + 1,
                      ).map((pageNumber) => (
                        <PaginationItem key={pageNumber}>
                          <PaginationLink
                            isActive={pageNumber === currentProspectPage}
                            onClick={() => setProspectPage(pageNumber)}
                          >
                            {pageNumber}
                          </PaginationLink>
                        </PaginationItem>
                      ))}
                      <PaginationItem>
                        <PaginationNext
                          onClick={() =>
                            setProspectPage(
                              Math.min(
                                prospectTotalPages,
                                currentProspectPage + 1,
                              ),
                            )
                          }
                          disabled={currentProspectPage >= prospectTotalPages}
                          aria-disabled={
                            currentProspectPage >= prospectTotalPages
                          }
                        />
                      </PaginationItem>
                      <PaginationItem>
                        <PaginationLast
                          onClick={() => setProspectPage(prospectTotalPages)}
                          disabled={currentProspectPage >= prospectTotalPages}
                          aria-disabled={
                            currentProspectPage >= prospectTotalPages
                          }
                        />
                      </PaginationItem>
                    </PaginationContent>
                  </Pagination>
                </CardContent>
              </>
            )}
          </Card>
        </div>
      </TabsContent>
    </Tabs>
  );
}
