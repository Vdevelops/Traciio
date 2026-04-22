"use client";

import { useState, useMemo } from "react";
import { useTranslations } from "next-intl";
import { SalesPerformanceList } from "@/features/sales-overview/components/sales-performance-list";
import { SalesOverviewChart } from "@/features/sales-overview/components/SalesOverviewChart";
import { useMonthlySalesOverview } from "@/features/sales-overview/hooks/useMonthlySalesOverview";
import { useSalesPerformanceList } from "@/features/sales-overview/hooks/useSalesPerformanceList";
import type { DateRange } from "react-day-picker";
import { startOfYear, endOfYear, format } from "date-fns";

export function SalesOverviewPageClient() {
  const t = useTranslations("salesOverview");

  // Filter State for Chart
  const [filterMode, setFilterMode] = useState<"year" | "range">("year");
  const [selectedYear, setSelectedYear] = useState<number>(new Date().getFullYear());
  const [selectedMetric, setSelectedMetric] = useState<"revenue" | "deals" | "visits" | "tasks">("revenue");
  const [dateRange, setDateRange] = useState<DateRange | undefined>(() => {
    const start = startOfYear(new Date());
    const end = endOfYear(new Date());
    return { from: start, to: end };
  });

  // Calculate generic start/end dates based on mode for Chart
  const { startDate: chartStartDate, endDate: chartEndDate } = useMemo(() => {
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
  const { monthlyData, isLoading: isChartLoading } = useMonthlySalesOverview(chartStartDate, chartEndDate);

  // List Data Hook (Lifted state)
  const listProps = useSalesPerformanceList();

  // Sync Chart Date Filter to List
  useMemo(() => {
    if (chartStartDate) listProps.setStartDate(chartStartDate);
    if (chartEndDate) listProps.setEndDate(chartEndDate);
  }, [chartStartDate, chartEndDate]);

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
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="lucide lucide-chart-column h-8 w-8 text-primary" aria-hidden="true"><path d="M3 3v18h18"/><path d="M18 17V9"/><path d="M13 17V5"/><path d="M8 17v-3"/></svg>
            {t("title")}
          </h1>
          <p className="text-muted-foreground mt-1 text-sm">{t("description")}</p>
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
      />

      {/* List Section */}
      <div className="space-y-4">
        <SalesPerformanceList {...listProps} />
      </div>
    </div>
  );
}
