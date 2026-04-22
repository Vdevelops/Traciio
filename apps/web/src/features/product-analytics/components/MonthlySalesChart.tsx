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
import type { MonthlySalesResponse } from "../types/monthly-sales";
import { DateRangePicker } from "@/components/ui/date-range-picker";
import type { DateRange } from "react-day-picker";

export interface MonthlySalesChartProps {
  readonly data: MonthlySalesResponse | undefined;
  readonly isLoading?: boolean;
  readonly filterMode: "year" | "range";
  readonly onFilterModeChange: (mode: "year" | "range") => void;
  readonly selectedYear: number;
  readonly onYearChange: (year: number) => void;
  readonly dateRange: DateRange | undefined;
  readonly onDateChange: (range: DateRange | undefined) => void;
  readonly selectedMetric?: "total_sold" | "total_revenue" | "total_profit" | "sales_count";
  readonly onMetricChange?: (metric: "total_sold" | "total_revenue" | "total_profit" | "sales_count") => void;
}

const chartConfig = {
  total_sold: {
    label: "Total Sold",
    color: "#F39200",
  },
  total_revenue: {
    label: "Revenue",
    color: "#E07B00",
  },
  total_profit: {
    label: "Profit",
    color: "#10B981",
  },
  sales_count: {
    label: "Sales Count",
    color: "#3B82F6",
  },
} satisfies ChartConfig;

const formatNumber = (value: number): string => {
  return value.toLocaleString("id-ID");
};

// Generate year options dynamically (2000 to current year + 1)
// This ensures future years are automatically included
const generateYearOptions = () => {
  const currentYear = new Date().getFullYear();
  const startYear = 2000; // Business started year
  const endYear = currentYear + 1; // Allow next year for planning
  const years: number[] = [];
  
  for (let i = endYear; i >= startYear; i--) {
    years.push(i); // Descending order (newest first)
  }
  return years;
};

export function MonthlySalesChart({
  data,
  isLoading,
  filterMode,
  onFilterModeChange,
  selectedYear,
  onYearChange,
  dateRange,
  onDateChange,
  selectedMetric = "total_sold",
  onMetricChange,
}: MonthlySalesChartProps) {
  const t = useTranslations("productAnalytics.monthlySales");
  const years = React.useMemo(() => generateYearOptions(), []);

  const metricLabels = {
    total_sold: t("totalSold"),
    total_revenue: t("totalRevenue"),
    total_profit: t("totalProfit"),
    sales_count: t("totalSales"),
  };

  const chartData = React.useMemo(() => {
    if (!data?.monthly_sales) return [];

    return data.monthly_sales.map((month) => ({
      month: month.month_name.substring(0, 3), // Jan, Feb, Mar, etc.
      total_sold: month.total_sold,
      total_revenue: month.total_revenue,
      total_profit: month.total_profit,
      sales_count: month.sales_count,
    }));
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

  if (!data) {
    return (
      <Card>
        <CardContent className="text-center py-12 text-muted-foreground">
          {t("noData")}
        </CardContent>
      </Card>
    );
  }

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
                  <SelectItem value="total_sold">{metricLabels.total_sold}</SelectItem>
                  <SelectItem value="total_revenue">{metricLabels.total_revenue}</SelectItem>
                  <SelectItem value="total_profit">{metricLabels.total_profit}</SelectItem>
                  <SelectItem value="sales_count">{metricLabels.sales_count}</SelectItem>
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
                  onDateChange={onDateChange}
                  placeholder="Select range"
                />
              </div>
            )}
          </div>
        </div>

        {/* Summary Stats */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mt-6">
          <div className="space-y-1">
            <p className="text-xs text-muted-foreground">{t("totalSold")}</p>
            <p className="text-lg font-medium">{formatNumber(data.total_sold)}</p>
          </div>
          <div className="space-y-1">
            <p className="text-xs text-muted-foreground">{t("totalRevenue")}</p>
            <p className="text-lg font-medium">{formatCurrency(data.total_revenue)}</p>
          </div>
          <div className="space-y-1">
            <p className="text-xs text-muted-foreground">{t("totalProfit")}</p>
            <p className="text-lg font-medium">{formatCurrency(data.total_profit)}</p>
          </div>
          <div className="space-y-1">
            <p className="text-xs text-muted-foreground">{t("totalSales")}</p>
            <p className="text-lg font-medium">{formatNumber(data.total_sales)}</p>
          </div>
        </div>
      </CardHeader>
      <CardContent>
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
              tickFormatter={(value) => formatNumber(value)}
            />
            <ChartTooltip
              content={
                <ChartTooltipContent
                  formatter={(value, name) => {
                    if (name === "total_revenue") {
                      return formatCurrency(Number(value));
                    }
                    return formatNumber(Number(value));
                  }}
                  labelFormatter={(label) => `${label}`}
                />
              }
            />
            <Bar
              dataKey={selectedMetric}
              fill={`var(--color-${selectedMetric})`}
              radius={[8, 8, 0, 0]}
              maxBarSize={60}
            >
              <LabelList
                dataKey={selectedMetric}
                position="top"
                formatter={(value) => {
                  const numValue = Number(value);
                  if (selectedMetric === "total_revenue" || selectedMetric === "total_profit") {
                    return formatCurrency(numValue);
                  }
                  return formatNumber(numValue);
                }}
                className="fill-foreground text-xs"
              />
            </Bar>
          </BarChart>
        </ChartContainer>
      </CardContent>
    </Card>
  );
}
