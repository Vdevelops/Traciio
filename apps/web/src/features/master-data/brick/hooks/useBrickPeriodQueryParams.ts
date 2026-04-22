"use client";

import { useCallback, useMemo } from "react";
import { usePathname, useRouter, useSearchParams } from "next/navigation";
import { useCurrentMonthPeriod } from "./useBrickAnalytics";

export type BrickPeriodMode = "range" | "month" | "year";

function isValidDateOnly(value: string | null): value is string {
  if (!value) return false;
  return /^\d{4}-\d{2}-\d{2}$/.test(value);
}

function clampPeriod(start: string, end: string): { periodStart: string; periodEnd: string } {
  if (start <= end) return { periodStart: start, periodEnd: end };
  return { periodStart: end, periodEnd: start };
}

function lastDayOfMonth(year: number, month1Based: number): number {
  return new Date(year, month1Based, 0).getDate();
}

export function toMonthPeriod(monthValue: string): { periodStart: string; periodEnd: string } {
  // monthValue: YYYY-MM
  const [yearStr, monthStr] = monthValue.split("-");
  const year = Number(yearStr);
  const month = Number(monthStr);
  const lastDay = lastDayOfMonth(year, month);
  const periodStart = `${yearStr}-${monthStr}-01`;
  const periodEnd = `${yearStr}-${monthStr}-${String(lastDay).padStart(2, "0")}`;
  return { periodStart, periodEnd };
}

export function toYearPeriod(yearValue: string): { periodStart: string; periodEnd: string } {
  const year = yearValue.padStart(4, "0");
  return {
    periodStart: `${year}-01-01`,
    periodEnd: `${year}-12-31`,
  };
}

export function useBrickPeriodQueryParams() {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const { periodStart: defaultStart, periodEnd: defaultEnd } = useCurrentMonthPeriod();

  const mode = useMemo<BrickPeriodMode>(() => {
    const raw = searchParams.get("period_mode");
    if (raw === "range" || raw === "month" || raw === "year") return raw;
    return "month";
  }, [searchParams]);

  const { periodStart, periodEnd } = useMemo(() => {
    const rawStart = searchParams.get("period_start");
    const rawEnd = searchParams.get("period_end");

    const start = rawStart === "" ? "" : isValidDateOnly(rawStart) ? rawStart : defaultStart;
    const end = rawEnd === "" ? "" : isValidDateOnly(rawEnd) ? rawEnd : defaultEnd;

    // If either side is explicitly cleared, keep as-is (represents "all time" in range mode)
    if (start === "" || end === "") {
      return { periodStart: start, periodEnd: end };
    }

    return clampPeriod(start, end);
  }, [searchParams, defaultStart, defaultEnd]);

  const setPeriodParams = useCallback(
    (next: {
      mode?: BrickPeriodMode;
      periodStart?: string;
      periodEnd?: string;
    }) => {
      const params = new URLSearchParams(searchParams.toString());

      if (next.mode !== undefined) params.set("period_mode", next.mode);
      if (next.periodStart !== undefined) params.set("period_start", next.periodStart);
      if (next.periodEnd !== undefined) params.set("period_end", next.periodEnd);

      router.replace(`${pathname}?${params.toString()}`, { scroll: false });
    },
    [router, pathname, searchParams]
  );

  return {
    mode,
    periodStart,
    periodEnd,
    setPeriodParams,
  };
}
