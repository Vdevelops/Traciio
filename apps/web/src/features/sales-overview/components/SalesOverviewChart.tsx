"use client";

import * as React from "react";
import { useTranslations } from "next-intl";
import { Bar, BarChart, CartesianGrid, XAxis, YAxis, LabelList } from "recharts";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import {
  ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Skeleton } from "@/components/ui/skeleton";
import { Calendar } from "lucide-react";
import { formatCurrency } from "@/lib/utils";
import { DateRangePicker } from "@/components/ui/date-range-picker";
import type { DateRange } from "react-day-picker";
import type { MonthlySalesData } from "../types";

export interface SalesOverviewChartProps {
  readonly data: MonthlySalesData[];
  readonly isLoading?: boolean;
  readonly filterMode: "year" | "range";
  readonly onFilterModeChange: (mode: "year" | "range") => void;
  readonly selectedYear: number;
  readonly onYearChange: (year: number) => void;
  readonly dateRange: DateRange | undefined;
  readonly onDateRangeChange: (range: DateRange | undefined) => void;
  readonly selectedMetric?: "revenue" | "deals" | "visits" | "tasks";
  readonly onMetricChange?: (metric: "revenue" | "deals" | "visits" | "tasks") => void;
}

const chartConfig = {
  revenue: {
    label: "Revenue",
    color: "#E07B00",
  },
  deals: {
    label: "Deals Closed",
    color: "#3B82F6",
  },
  visits: {
    label: "Visits",
    color: "#10B981",
  },
  tasks: {
    label: "Tasks",
    color: "#8B5CF6",
  },
  target: {
    label: "Target",
    color: "#94A3B8", // Slate-400 for target
  },
} satisfies ChartConfig;

const formatNumber = (value: number): string => {
  return value.toLocaleString("id-ID");
};

// Generate year options dynamically (2000 to current year + 1)
const generateYearOptions = () => {
  const currentYear = new Date().getFullYear();
  const startYear = 2000;
  const endYear = currentYear + 1;
  const years: number[] = [];
  
  for (let i = endYear; i >= startYear; i--) {
    years.push(i);
  }
  return years;
};

export function SalesOverviewChart({
  data,
  isLoading,
  filterMode,
  onFilterModeChange,
  selectedYear,
  onYearChange,
  dateRange,
  onDateRangeChange,
  selectedMetric = "revenue",
  onMetricChange,
}: SalesOverviewChartProps) {
  const t = useTranslations("salesOverview.chart");
  const years = React.useMemo(() => generateYearOptions(), []);

  // Use translations for metric labels if available, layout matching Product Analytics
  const metricLabels = {
    revenue: t("metrics.revenue") || "Revenue",
    deals: t("metrics.deals") || "Deals Closed",
    visits: t("metrics.visits") || "Visits Completed",
    tasks: t("metrics.tasks") || "Tasks Completed",
  };

  const chartData = React.useMemo(() => {
    if (!data) return [];

    return data.map((month) => ({
      month: month.month_name.substring(0, 3), 
      revenue: month.total_revenue,
      target: month.target_amount, // Assuming this field exists in MonthlySalesData
      deals: month.total_deals,
      visits: month.total_visits,
      tasks: month.total_tasks,
    }));
  }, [data]);

  // Calculate totals for Summary Stats
  const totals = React.useMemo(() => {
    if (!data) return { revenue: 0, deals: 0, visits: 0, tasks: 0 };
    return data.reduce(
      (acc, curr) => ({
        revenue: acc.revenue + curr.total_revenue,
        deals: acc.deals + curr.total_deals,
        visits: acc.visits + curr.total_visits,
        tasks: acc.tasks + curr.total_tasks,
      }),
      { revenue: 0, deals: 0, visits: 0, tasks: 0 }
    );
  }, [data]);

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <Skeleton className="h-6 w-48" />
          <Skeleton className="h-4 w-64 mt-2" />
        </CardHeader>
        <CardContent>
          <Skeleton className="h-[400px] w-full" />
        </CardContent>
      </Card>
    );
  }

  // Determine if we have valid data to show
  const hasData = data && data.length > 0;

  return (
    <Card>
      <CardHeader>
        <div className="flex items-start justify-between">
          <div className="space-y-1">
            <CardTitle className="flex items-center gap-2">
              <Calendar className="h-5 w-5 text-primary" />
              {t("title")}
            </CardTitle>
            <CardDescription>{t("description")}</CardDescription>
          </div>
          <div className="flex items-center gap-3">
            {/* Filter Mode Toggle */}
            <div className="flex items-center bg-muted p-1 rounded-lg">
              <button
                onClick={() => onFilterModeChange("year")}
                className={`px-3 py-1.5 text-xs font-medium rounded-md transition-all cursor-pointer ${
                  filterMode === "year" 
                    ? "bg-background text-foreground shadow-sm" 
                    : "text-muted-foreground hover:text-foreground"
                }`}
              >
                Yearly
              </button>
              <button
                onClick={() => onFilterModeChange("range")}
                className={`px-3 py-1.5 text-xs font-medium rounded-md transition-all cursor-pointer ${
                  filterMode === "range" 
                    ? "bg-background text-foreground shadow-sm" 
                    : "text-muted-foreground hover:text-foreground"
                }`}
              >
                Range
              </button>
            </div>

            {/* Metric Selector */}
            {onMetricChange && (
              <Select
                value={selectedMetric}
                onValueChange={(value) => onMetricChange(value as typeof selectedMetric)}
              >
                <SelectTrigger className="w-[140px] h-9">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="revenue">{metricLabels.revenue}</SelectItem>
                  <SelectItem value="deals">{metricLabels.deals}</SelectItem>
                  <SelectItem value="visits">{metricLabels.visits}</SelectItem>
                  <SelectItem value="tasks">{metricLabels.tasks}</SelectItem>
                </SelectContent>
              </Select>
            )}

            {/* Year Selector (Yearly Mode) */}
            {filterMode === "year" && (
              <Select
                value={selectedYear.toString()}
                onValueChange={(val) => onYearChange(parseInt(val))}
              >
                <SelectTrigger className="w-[110px] h-9">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {years.map((year) => (
                    <SelectItem key={year} value={year.toString()}>
                      {year}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            )}

            {/* Date Range Picker (Range Mode) */}
            {filterMode === "range" && (
              <div>
                <DateRangePicker
                  dateRange={dateRange}
                  onDateChange={onDateRangeChange}
                  placeholder="Select range"
                />
              </div>
            )}
          </div>
        </div>

        {/* Summary Stats */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mt-6">
          <div className="space-y-1">
            <p className="text-xs text-muted-foreground">{metricLabels.revenue}</p>
            <p className="text-lg font-medium">{formatCurrency(totals.revenue)}</p>
          </div>
          <div className="space-y-1">
            <p className="text-xs text-muted-foreground">{metricLabels.deals}</p>
            <p className="text-lg font-medium">{formatNumber(totals.deals)}</p>
          </div>
          <div className="space-y-1">
            <p className="text-xs text-muted-foreground">{metricLabels.visits}</p>
            <p className="text-lg font-medium">{formatNumber(totals.visits)}</p>
          </div>
          <div className="space-y-1">
            <p className="text-xs text-muted-foreground">{metricLabels.tasks}</p>
            <p className="text-lg font-medium">{formatNumber(totals.tasks)}</p>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        {/* Chart Area */}
        {!hasData ? (
          <div className="text-center py-12 text-muted-foreground">
            {t("noData")}
          </div>
        ) : (
          <ChartContainer config={chartConfig} className="h-[400px] w-full">
            <BarChart data={chartData} margin={{ top: 20, right: 20, bottom: 20, left: 20 }}>
              <CartesianGrid strokeDasharray="3 3" vertical={false} />
              <XAxis
                dataKey="month"
                tickLine={false}
                tickMargin={10}
                axisLine={false}
              />
              <YAxis
                tickLine={false}
                axisLine={false}
                tickFormatter={(value) => {
                  if (selectedMetric === "revenue") {
                      if (value >= 1000000000) return `${(value / 1000000000).toFixed(1)}M`;
                      if (value >= 1000000) return `${(value / 1000000).toFixed(0)}Jt`;
                      return `${(value / 1000).toFixed(0)}k`;
                  }
                  return formatNumber(value);
                }}
              />
              <ChartTooltip
                content={
                  <ChartTooltipContent
                formatter={(value, name) => {
                  if (name === "revenue") {
                    return (
                      <>
                        <span className="font-medium">{metricLabels.revenue}</span>
                        <span className="ml-auto font-mono font-medium tabular-nums text-foreground">
                          {formatCurrency(Number(value))}
                        </span>
                      </>
                    );
                  }
                  if (name === "target") {
                    return (
                      <>
                        <span className="font-medium">Target</span>
                        <span className="ml-auto font-mono font-medium tabular-nums text-foreground">
                          {formatCurrency(Number(value))}
                        </span>
                      </>
                    );
                  }
                  return (
                    <>
                      <span className="font-medium">{chartConfig[name as keyof typeof chartConfig]?.label || name}</span>
                      <span className="ml-auto font-mono font-medium tabular-nums text-foreground">
                        {formatNumber(Number(value))}
                      </span>
                    </>
                  );
                }}
              />
            }
          />
          <Bar
            dataKey={selectedMetric}
            name={selectedMetric}
            fill={`var(--color-${selectedMetric})`}
            radius={[4, 4, 0, 0]}
            maxBarSize={50}
          />
          {selectedMetric === "revenue" && (
            <Bar
              dataKey="target"
              name="target"
              fill="var(--color-target)"
              radius={[4, 4, 0, 0]}
              maxBarSize={50}
            />
          )}
            </BarChart>
          </ChartContainer>
        )}
      </CardContent>
    </Card>
  );
}
