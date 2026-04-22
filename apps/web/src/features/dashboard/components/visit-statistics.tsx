"use client";

import * as React from "react";
import { useTranslations, useLocale } from "next-intl";
import { Area, AreaChart, CartesianGrid, XAxis, YAxis } from "recharts";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import {
  ChartConfig,
  ChartContainer,
  ChartLegend,
  ChartLegendContent,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";
import { useVisitStatistics } from "../hooks/useDashboard";
import { Skeleton } from "@/components/ui/skeleton";


const chartConfig = {
  completed: {
    label: "Completed",
    color: "oklch(0.5234 0.1347 144.1672)", // Primary - Hijau
  },
  approved: {
    label: "Approved",
    color: "oklch(0.55 0.15 240)", // Biru
  },
  pending: {
    label: "Pending",
    color: "oklch(0.55 0.15 300)", // Purple
  },
} satisfies ChartConfig;

import { useHasPermission } from "@/features/auth/providers/permissions-provider";
import { BarChart3, Lock } from "lucide-react";

export function VisitStatistics({ startDate, endDate }: Readonly<{ startDate?: string; endDate?: string }>) {
  const t = useTranslations("visitStatistics");
  const locale = useLocale();
  
  // Permission check
  const hasPermission = useHasPermission("visit-reports.view");

  // Use global date range, default to current month if not provided
  // actually default is handled by backend or parent component usually, but here 
  // we want to ensure we pass what we have.
  // The original code mixed monthParams with props. 
  // If startDate/endDate are provided, they should override.
  const params = { start_date: startDate, end_date: endDate };
  
  const { data, isLoading, isError, error } = useVisitStatistics(params, {
    enabled: hasPermission
  });

  console.log("VisitStatistics Params:", params);
  console.log("VisitStatistics Data:", data);

  const stats = data?.data;

  // Prepare chart data - must be called unconditionally (Rules of Hooks)
  const chartData = React.useMemo(() => {
    if (!stats?.by_date || stats.by_date.length === 0) {
      return [];
    }

    return stats.by_date
      .filter((item) => {
        if (!item?.date) return false;
        const itemDate = new Date(item.date);
        if (Number.isNaN(itemDate.getTime())) return false;
        return true;
      })
      .map((item) => ({
        date: item.date,
        total: item.count ?? 0,
        completed: item.completed ?? 0,
        approved: item.approved ?? 0,
        pending: item.pending ?? 0,
      }))
      .sort((a, b) => {
        const dateA = new Date(a.date).getTime();
        const dateB = new Date(b.date).getTime();
        return dateA - dateB;
      });
  }, [stats]);

  // Show warning if no permission
  if (!hasPermission) {
    return (
      <Card className="border-0 shadow-sm bg-muted/10 h-full">
        <CardHeader className="px-3 sm:px-6 pb-2">
          <CardTitle className="text-sm sm:text-base flex items-center gap-2 text-muted-foreground">
            <Lock className="h-4 w-4" />
            {t("title")}
          </CardTitle>
        </CardHeader>
        <CardContent className="flex flex-col items-center justify-center p-6 text-center h-[250px]">
          <div className="bg-muted/20 p-3 rounded-full mb-3">
            <Lock className="h-6 w-6 text-muted-foreground" />
          </div>
          <p className="text-sm font-medium text-muted-foreground mb-1">
            Access Restricted
          </p>
          <p className="text-xs text-muted-foreground/70 max-w-[200px]">
            You don't have permission to view visit statistics.
          </p>
        </CardContent>
      </Card>
    );
  }

  if (isLoading) {
    return (
      <Card className="border-0 shadow-sm">
        <CardHeader className="px-3 sm:px-6">
          <CardTitle className="text-sm sm:text-base">{t("title")}</CardTitle>
        </CardHeader>
        <CardContent className="px-3 sm:px-6 pb-3 sm:pb-6">
          <Skeleton className="h-48 sm:h-64 w-full" />
        </CardContent>
      </Card>
    );
  }

  if (isError) {
    return (
      <Card className="border-0 shadow-sm">
        <CardHeader className="px-3 sm:px-6">
          <CardTitle className="text-sm sm:text-base">{t("title")}</CardTitle>
        </CardHeader>
        <CardContent className="px-3 sm:px-6 pb-3 sm:pb-6">
          <div className="flex items-center justify-center h-[200px] sm:h-[250px] md:h-[300px] w-full border border-dashed rounded-lg">
            <p className="text-sm text-destructive">
              {error instanceof Error ? error.message : "Failed to load visit statistics"}
            </p>
          </div>
        </CardContent>
      </Card>
    );
  }

  if (!stats) {
    return null;
  }

  return (
    <Card className="border-0 shadow-sm">
      <CardHeader className="relative flex flex-col sm:flex-row items-start sm:items-center gap-2 sm:gap-2 space-y-0 px-3 sm:px-6">
        <div className="grid flex-1 gap-0.5 sm:gap-1 w-full sm:w-auto">
          <div className="flex items-center gap-1.5 sm:gap-2">
            <BarChart3 className="h-4 w-4 sm:h-5 sm:w-5" />
            <CardTitle className="text-sm sm:text-base">{t("title")}</CardTitle>
          </div>
          <CardDescription className="text-xs sm:text-sm">{t("description")}</CardDescription>
        </div>
      </CardHeader>
      <CardContent className="px-3 sm:px-6 pb-3 sm:pb-6">
        <div className="space-y-3 sm:space-y-4">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-2 sm:gap-4">
            <div>
              <div className="text-xs sm:text-sm text-muted-foreground">{t("total")}</div>
              <div className="text-xl sm:text-2xl font-medium">{stats.total}</div>
            </div>
            <div>
              <div className="text-xs sm:text-sm text-muted-foreground">{t("completed")}</div>
              <div className="text-xl sm:text-2xl font-medium">{stats.completed}</div>
            </div>
            <div>
              <div className="text-xs sm:text-sm text-muted-foreground">{t("pending")}</div>
              <div className="text-xl sm:text-2xl font-medium">{stats.pending}</div>
            </div>
            <div>
              <div className="text-xs sm:text-sm text-muted-foreground">{t("approved")}</div>
              <div className="text-xl sm:text-2xl font-medium">{stats.approved}</div>
            </div>
          </div>

          {chartData.length > 0 ? (
            <ChartContainer config={chartConfig} className="aspect-auto h-[200px] sm:h-[250px] md:h-[300px] w-full">
              <AreaChart data={chartData}>
                <defs>
                  {/* Stacked: pending + completed = total (no visual overlap) */}
                  <linearGradient id="fillPending" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="var(--color-pending)" stopOpacity={0.7} />
                    <stop offset="95%" stopColor="var(--color-pending)" stopOpacity={0.1} />
                  </linearGradient>
                  <linearGradient id="fillCompleted" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="var(--color-completed)" stopOpacity={0.8} />
                    <stop offset="95%" stopColor="var(--color-completed)" stopOpacity={0.1} />
                  </linearGradient>
                  <linearGradient id="fillApproved" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="var(--color-approved)" stopOpacity={0.6} />
                    <stop offset="95%" stopColor="var(--color-approved)" stopOpacity={0.1} />
                  </linearGradient>
                </defs>
                <CartesianGrid vertical={false} />
                <XAxis
                  dataKey="date"
                  tickLine={false}
                  axisLine={false}
                  tickMargin={8}
                  minTickGap={32}
                  tickFormatter={(value) => {
                    const date = new Date(value);
                    return date.toLocaleDateString(locale, {
                      month: "short",
                      day: "numeric",
                    });
                  }}
                />
                <YAxis
                  tickLine={false}
                  axisLine={false}
                  tickMargin={8}
                  tickFormatter={(value) => value.toLocaleString(locale)}
                />
                <ChartTooltip
                  cursor={false}
                  content={
                    <ChartTooltipContent
                      labelFormatter={(value) => {
                        const dateValue =
                          typeof value === "string" || typeof value === "number"
                            ? new Date(value)
                            : value instanceof Date
                              ? value
                              : new Date(String(value));
                        return dateValue.toLocaleDateString(locale, {
                          month: "short",
                          day: "numeric",
                          year: "numeric",
                        });
                      }}
                      indicator="dot"
                    />
                  }
                />
                {/* Stack pending + completed so they fill up to the total without overlap */}
                <Area
                  dataKey="pending"
                  type="natural"
                  fill="url(#fillPending)"
                  stroke="var(--color-pending)"
                  strokeWidth={2}
                  stackId="visits"
                />
                <Area
                  dataKey="completed"
                  type="natural"
                  fill="url(#fillCompleted)"
                  stroke="var(--color-completed)"
                  strokeWidth={2}
                  stackId="visits"
                />
                <Area
                  dataKey="approved"
                  type="natural"
                  fill="url(#fillApproved)"
                  stroke="var(--color-approved)"
                  strokeWidth={1.5}
                  strokeDasharray="4 2"
                />
                <ChartLegend content={<ChartLegendContent />} />
              </AreaChart>
            </ChartContainer>
          ) : (
            <div className="flex items-center justify-center h-[200px] sm:h-[250px] md:h-[300px] w-full border border-dashed rounded-lg">
              <p className="text-sm text-muted-foreground">{t("noData")}</p>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}

