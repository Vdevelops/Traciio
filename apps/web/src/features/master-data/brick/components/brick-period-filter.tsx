"use client";

import * as React from "react";
import { useTranslations } from "next-intl";
import { DateRangePicker } from "@/components/ui/date-range-picker";
import type { DateRange } from "react-day-picker";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  toYearPeriod,
  useBrickPeriodQueryParams,
} from "../hooks/useBrickPeriodQueryParams";

type UIMode = "year" | "range";

function toDateRange(periodStart: string, periodEnd: string): DateRange {
  return {
    from: new Date(`${periodStart}T00:00:00`),
    to: new Date(`${periodEnd}T00:00:00`),
  };
}

function formatDateOnlyLocal(date: Date): string {
  const d = new Date(date);
  d.setHours(0, 0, 0, 0);
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, "0")}-${String(
    d.getDate()
  ).padStart(2, "0")}`;
}

function generateYearOptions(): number[] {
  const currentYear = new Date().getFullYear();
  const startYear = 2000;
  const endYear = currentYear + 1;
  const years: number[] = [];

  for (let year = endYear; year >= startYear; year--) {
    years.push(year);
  }

  return years;
}

export function BrickPeriodFilter() {
  const t = useTranslations("brickPeriodFilter");
  const { mode, periodStart, periodEnd, setPeriodParams } = useBrickPeriodQueryParams();

  const uiMode: UIMode = mode === "year" ? "year" : "range";
  const selectedYear = React.useMemo(() => {
    const parsed = parseInt(periodStart.slice(0, 4));
    if (Number.isFinite(parsed) && parsed >= 2000) return parsed;
    return new Date().getFullYear();
  }, [periodStart]);
  const years = React.useMemo(() => generateYearOptions(), []);

  const dateRange = React.useMemo(
    () => {
      if (!periodStart) return undefined;
      if (!periodEnd) return { from: new Date(`${periodStart}T00:00:00`), to: undefined };
      return toDateRange(periodStart, periodEnd);
    },
    [periodStart, periodEnd]
  );

  const handleModeChange = (next: UIMode) => {
    if (next === "year") {
      const nextPeriod = toYearPeriod(String(selectedYear));
      setPeriodParams({ mode: "year", ...nextPeriod });
      return;
    }

    // Range (also covers old "month" mode semantics)
    setPeriodParams({ mode: "range" });
  };

  const handleDateRangeChange = (range: DateRange | undefined) => {
    if (!range?.from) {
      // Clear date filter - all time (match Sales Performance)
      setPeriodParams({ mode: "range", periodStart: "", periodEnd: "" });
      return;
    }

    const fromStr = formatDateOnlyLocal(range.from);
    const toStr = range.to ? formatDateOnlyLocal(range.to) : "";
    setPeriodParams({ mode: "range", periodStart: fromStr, periodEnd: toStr });
  };

  const handleYearChange = (value: string) => {
    const nextYear = parseInt(value);
    if (!Number.isFinite(nextYear)) return;
    const next = toYearPeriod(String(nextYear));
    setPeriodParams({ mode: "year", ...next });
  };

  return (
    <div className="flex flex-wrap items-center gap-3">
      <div className="flex items-center bg-muted p-1 rounded-lg">
        <button
          onClick={() => handleModeChange("year")}
          className={`px-3 py-1.5 text-xs font-medium rounded-md transition-all cursor-pointer ${
            uiMode === "year"
              ? "bg-background text-foreground shadow-sm"
              : "text-muted-foreground hover:text-foreground"
          }`}
          type="button"
        >
          {t("year")}
        </button>
        <button
          onClick={() => handleModeChange("range")}
          className={`px-3 py-1.5 text-xs font-medium rounded-md transition-all cursor-pointer ${
            uiMode === "range"
              ? "bg-background text-foreground shadow-sm"
              : "text-muted-foreground hover:text-foreground"
          }`}
          type="button"
        >
          {t("range")}
        </button>
      </div>

      {uiMode === "year" && (
        <Select value={selectedYear.toString()} onValueChange={handleYearChange}>
          <SelectTrigger className="w-[110px] h-9">
            <SelectValue />
          </SelectTrigger>
          <SelectContent
            align="start"
            className="max-h-[min(320px,var(--radix-select-content-available-height))]"
          >
            {years.map((year) => (
              <SelectItem key={year} value={year.toString()}>
                {year}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      )}

      {uiMode === "range" && (
        <div className="w-full sm:w-auto">
          <DateRangePicker
            dateRange={dateRange}
            onDateChange={handleDateRangeChange}
            placeholder={t("selectRange")}
          />
        </div>
      )}
    </div>
  );
}
