"use client";

import { useEffect, useState } from "react";
import { useTranslations } from "next-intl";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Label } from "@/components/ui/label";
import { DateRangePicker } from "@/components/ui/date-range-picker";
import type { DateRange } from "react-day-picker";

interface ProductAnalyticsFiltersProps {
  readonly period?: "day" | "week" | "month" | "year";
  readonly onPeriodChange?: (period: "day" | "week" | "month" | "year") => void;
  readonly metric?: "quantity" | "revenue" | "growth";
  readonly onMetricChange?: (metric: "quantity" | "revenue" | "growth") => void;
  readonly dateRange?: DateRange;
  readonly onDateRangeChange?: (range: DateRange | undefined) => void;
  readonly groupBy?: "day" | "week" | "month" | "year";
  readonly onGroupByChange?: (groupBy: "day" | "week" | "month" | "year") => void;
  readonly showDateRange?: boolean;
  readonly showGroupBy?: boolean;
}

export function ProductAnalyticsFilters({
  period = "month",
  onPeriodChange,
  metric = "revenue",
  onMetricChange,
  dateRange,
  onDateRangeChange,
  groupBy = "month",
  onGroupByChange,
  showDateRange = false,
  showGroupBy = false,
}: ProductAnalyticsFiltersProps) {
  const t = useTranslations("productAnalytics.filters");
  const [mounted, setMounted] = useState(false);

  // Prevent hydration mismatch
  useEffect(() => {
    setMounted(true);
  }, []);

  if (!mounted) {
    return (
      <div className="flex flex-wrap items-end gap-4">
        <div className="flex-1 min-w-[200px]">
          <Label htmlFor="period">{t("period")}</Label>
          <div className="mt-1 h-9 w-full rounded-md border bg-transparent" />
        </div>
        <div className="flex-1 min-w-[200px]">
          <Label htmlFor="metric">{t("metric")}</Label>
          <div className="mt-1 h-9 w-full rounded-md border bg-transparent" />
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-wrap items-end gap-4">
      {!showDateRange && onPeriodChange && (
        <div className="flex-1 min-w-[200px]">
          <Label htmlFor="period">{t("period")}</Label>
          <Select value={period} onValueChange={onPeriodChange}>
            <SelectTrigger id="period" className="mt-1">
              <SelectValue placeholder={t("period")} />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="day">{t("day")}</SelectItem>
              <SelectItem value="week">{t("week")}</SelectItem>
              <SelectItem value="month">{t("month")}</SelectItem>
              <SelectItem value="year">{t("year")}</SelectItem>
            </SelectContent>
          </Select>
        </div>
      )}

      <div className="flex-1 min-w-[200px]">
        <Label htmlFor="metric">{t("metric")}</Label>
        <Select value={metric} onValueChange={onMetricChange}>
          <SelectTrigger id="metric" className="mt-1">
            <SelectValue placeholder={t("metric")} />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="quantity">{t("quantity")}</SelectItem>
            <SelectItem value="revenue">{t("revenue")}</SelectItem>
            <SelectItem value="growth">{t("growth")}</SelectItem>
          </SelectContent>
        </Select>
      </div>

      {showDateRange && onDateRangeChange && (
        <div className="flex-1">
          <Label>{t("startDate")} - {t("endDate")}</Label>
          <div className="mt-1">
            <DateRangePicker
              dateRange={dateRange}
              onDateChange={onDateRangeChange}
            />
          </div>
        </div>
      )}

      {showGroupBy && onGroupByChange && (
        <div className="flex-1 min-w-[200px]">
          <Label htmlFor="groupBy">{t("groupBy")}</Label>
          <Select value={groupBy} onValueChange={onGroupByChange}>
            <SelectTrigger id="groupBy" className="mt-1">
              <SelectValue placeholder={t("groupBy")} />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="day">{t("day")}</SelectItem>
              <SelectItem value="week">{t("week")}</SelectItem>
              <SelectItem value="month">{t("month")}</SelectItem>
              <SelectItem value="year">{t("year")}</SelectItem>
            </SelectContent>
          </Select>
        </div>
      )}
    </div>
  );
}