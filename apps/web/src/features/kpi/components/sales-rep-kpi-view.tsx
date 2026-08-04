"use client";

import React, { useState, useMemo } from "react";
import { format, startOfMonth, endOfMonth } from "date-fns";
import type { DateRange } from "react-day-picker";
import { ArrowLeft, RefreshCw, AlertCircle, Trophy } from "lucide-react";

import { useSalesRepKPI } from "../hooks/useKPIHooks";
import { DateRangePicker } from "@/components/ui/date-range-picker";
import { Button } from "@/components/ui/button";
import CompositeScoreCard from "./CompositeScoreCard";
import TargetGapSummary from "./TargetGapSummary";
import MetricCard from "./MetricCard";
import DiagnosticList from "./DiagnosticList";
import KPIBarChart from "./KPIBarChart";
import { formatCurrency } from "../utils/formatters";
import type { SalesRepKPIParams } from "../types";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

interface Props {
  userId?: string; // Optional, provided for manager drilldown
  userName?: string; // Optional, for manager drilldown header display
  onBack?: () => void; // Optional, callback for back button
  initialStartDate?: string;
  initialEndDate?: string;
}

function createInitialDateRange(initialStartDate?: string, initialEndDate?: string): DateRange {
  if (initialStartDate && initialEndDate) {
    const from = new Date(initialStartDate);
    const to = new Date(initialEndDate);

    if (!Number.isNaN(from.getTime()) && !Number.isNaN(to.getTime())) {
      return { from, to };
    }
  }

  const today = new Date();
  return {
    from: startOfMonth(today),
    to: endOfMonth(today),
  };
}

export default function SalesRepKPIView({
  userId,
  userName,
  onBack,
  initialStartDate,
  initialEndDate,
}: Props) {
  // Date state initialized to current month
  const [dateRange, setDateRange] = useState<DateRange | undefined>(() =>
    createInitialDateRange(initialStartDate, initialEndDate)
  );

  // Format dates for API
  const { startDate, endDate } = useMemo(() => {
    if (!dateRange?.from || !dateRange?.to) {
      const today = new Date();
      return {
        startDate: format(startOfMonth(today), "yyyy-MM-dd"),
        endDate: format(endOfMonth(today), "yyyy-MM-dd"),
      };
    }
    return {
      startDate: format(dateRange.from, "yyyy-MM-dd"),
      endDate: format(dateRange.to, "yyyy-MM-dd"),
    };
  }, [dateRange]);

  const params: SalesRepKPIParams = useMemo(() => {
    const p: SalesRepKPIParams = {
      startDate,
      endDate,
      compareWithPrevious: true,
    };
    if (userId) {
      p.userId = userId;
    }
    return p;
  }, [startDate, endDate, userId]);

  const { data, isLoading, error, refetch, isFetching } = useSalesRepKPI(params);

  const scorecard = data?.scorecard;
  const evaluation = data?.evaluation;
  const diagnostics = data?.diagnostics ?? [];
  const meta = data?.meta;

  const metricChartData = useMemo(() => {
    if (!scorecard) return [];

    const taskHealth =
      scorecard.overdueTaskRate === null || scorecard.overdueTaskRate === undefined
        ? 0
        : Math.max(0, 100 - scorecard.overdueTaskRate);

    return [
      { label: "Konversi", value: scorecard.conversionRate ?? 0 },
      { label: "Target", value: Math.min(scorecard.revenueTargetAttainment ?? 0, 100) },
      { label: "Visit", value: scorecard.visitCompliance ?? 0 },
      { label: "Task", value: taskHealth },
    ];
  }, [scorecard]);

  // Empty state check as per §7 requirements
  const isEmptyState = useMemo(() => {
    if (!data?.scorecard) return false;
    const sc = data.scorecard;
    const noActivity =
      sc.totalDealsClosed === 0 &&
      sc.dealsCreated === 0 &&
      sc.totalRevenue === 0 &&
      sc.visitCompleted === 0 &&
      sc.tasksCompleted === 0;

    return (
      noActivity &&
      (sc.visitPlanned === 0 || sc.visitPlanned === null) &&
      (sc.conversionRate === null || sc.conversionRate === 0)
    );
  }, [data]);

  return (
    <div className="space-y-8">
      {/* Header and Filter Control Area */}
      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div>
          {onBack && (
            <Button
              variant="ghost"
              size="sm"
              onClick={onBack}
              className="mb-2 -ml-2 text-muted-foreground hover:text-foreground"
            >
              <ArrowLeft className="mr-2 h-4 w-4" />
              Kembali ke Tim
            </Button>
          )}
          <h1 className="flex items-center gap-2 text-3xl font-medium tracking-tight">
            <Trophy className="h-8 w-8 text-primary" aria-hidden="true" />
            {userName ? `KPI: ${userName}` : "Evaluasi KPI"}
          </h1>
          <p className="mt-1 text-sm text-muted-foreground">
            {userName 
              ? `Analisis pencapaian target dan efisiensi kerja untuk ${userName}`
              : "Pantau pencapaian KPI, skor komposit, serta diagnosis kinerja pribadi Anda."}
          </p>
        </div>

        <div className="flex flex-wrap items-center gap-2">
          <DateRangePicker dateRange={dateRange} onDateChange={setDateRange} />
          
          <Button
            variant="outline"
            size="icon"
            onClick={() => refetch()}
            disabled={isLoading || isFetching}
            className="h-9 w-9 border-border/60"
          >
            <RefreshCw className={`h-4 w-4 text-muted-foreground ${isFetching ? "animate-spin" : ""}`} />
            <span className="sr-only">Refresh</span>
          </Button>
        </div>
      </div>

      {/* Loading Skeleton */}
      {isLoading ? (
        <div className="space-y-6">
          <div className="h-28 rounded-lg bg-muted" />
          <div className="h-44 rounded-lg bg-muted" />
          <div className="grid gap-4 grid-cols-1 sm:grid-cols-2 lg:grid-cols-3">
            {Array.from({ length: 9 }).map((_, idx) => (
              <div key={idx} className="h-24 rounded-lg bg-muted" />
            ))}
          </div>
        </div>
      ) : error ? (
        /* Error State */
        <div className="rounded-lg border border-destructive/20 bg-destructive/5 p-8 text-center">
          <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-destructive/10">
            <AlertCircle className="h-6 w-6 text-destructive" />
          </div>
          <h3 className="mt-4 text-base font-bold text-foreground">
            Gagal Memuat KPI
          </h3>
          <p className="mx-auto mt-2 max-w-md text-sm text-muted-foreground">
            Terjadi kesalahan koneksi atau Anda tidak memiliki akses untuk melihat data ini.
          </p>
          <Button onClick={() => refetch()} className="mt-6" variant="outline">
            Coba Lagi
          </Button>
        </div>
      ) : isEmptyState ? (
        /* Empty State as per §7 */
        <div className="rounded-lg border bg-muted/30 p-12 text-center">
          <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-muted">
            <AlertCircle className="h-6 w-6 text-muted-foreground" />
          </div>
          <h3 className="mt-4 text-base font-bold text-foreground">
            Tidak Cukup Data untuk Evaluasi
          </h3>
          <p className="mx-auto mt-2 max-w-md text-sm text-muted-foreground">
            Belum ada aktivitas tercatat pada periode ini. Silakan ubah rentang tanggal atau pilih periode lain.
          </p>
        </div>
      ) : (
        /* Main View */
        <div className="space-y-6">
          {/* Blok 1 — Ringkasan */}
          <div className="grid gap-6 lg:grid-cols-3">
            <div className="lg:col-span-1">
              <CompositeScoreCard 
                score={evaluation?.compositeScore ?? null} 
                grade={evaluation?.grade ?? null} 
                trend={evaluation?.trend ?? null} 
              />
            </div>
            <div className="lg:col-span-2">
              <TargetGapSummary 
                revenue={evaluation?.targetGap?.revenue ?? null} 
                deals={evaluation?.targetGap?.deals ?? null} 
              />
            </div>
          </div>

          <KPIBarChart
            title="Komposisi KPI Pribadi"
            description="Normalisasi metrik utama untuk composite score pribadi."
            data={metricChartData}
            valueSuffix="%"
            maxValue={100}
          />

          {/* Blok 2 — Scorecard & Trend Detail */}
          <Card>
            <CardHeader>
              <CardTitle>Rincian Kinerja</CardTitle>
              <CardDescription>Metrik mentah dari aktivitas, deal, target, dan pipeline.</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid gap-4 grid-cols-1 sm:grid-cols-2 lg:grid-cols-3">
              {/* Sales Metrics */}
              <MetricCard 
                label="Total Revenue" 
                value={formatCurrency(scorecard?.totalRevenue ?? null)} 
              />
              <MetricCard 
                label="Average Deal Value" 
                value={formatCurrency(scorecard?.averageDealValue ?? null)} 
              />
              <MetricCard 
                label="Revenue Target Attainment" 
                value={scorecard?.revenueTargetAttainment ?? null} 
                suffix="%" 
              />

              {/* Deal Metrics */}
              <MetricCard 
                label="Total Deals Closed" 
                value={scorecard?.totalDealsClosed ?? null} 
              />
              <MetricCard 
                label="Deals Created" 
                value={scorecard?.dealsCreated ?? null} 
              />
              <MetricCard 
                label="Deal Target Attainment" 
                value={scorecard?.dealTargetAttainment ?? null} 
                suffix="%" 
              />
              
              <MetricCard 
                label="Conversion Rate" 
                value={scorecard?.conversionRate ?? null} 
                suffix="%" 
              />
              <MetricCard 
                label="Pipeline Movement Score" 
                value={scorecard?.pipelineMovementScore ?? null} 
              />
              
              {/* Activity Metrics */}
              <MetricCard 
                label="Visit Compliance" 
                value={scorecard?.visitCompliance ?? null} 
                suffix="%" 
              />
              <MetricCard 
                label="Visit Completed" 
                value={scorecard?.visitCompleted ?? null} 
              />
              <MetricCard 
                label="Visit Planned" 
                value={scorecard?.visitPlanned ?? null} 
              />

              {/* Task Metrics */}
              <MetricCard 
                label="Tasks Completed" 
                value={scorecard?.tasksCompleted ?? null} 
              />
              <MetricCard 
                label="Overdue Task Rate" 
                value={scorecard?.overdueTaskRate ?? null} 
                suffix="%" 
              />
              </div>
            </CardContent>
          </Card>

          {/* Blok 3 — Diagnostic & Actionable Insight */}
          <DiagnosticList items={diagnostics} />

          {/* Metadata information */}
          {meta && (
            <div className="flex items-center justify-between px-2 text-[10px] text-muted-foreground">
              <span>
                Data fallback: {meta.brickInferredCount} dirujuk, {meta.brickMissingCount} tidak terpetakan.
              </span>
              <span>
                Diperbarui pada: {new Date(meta.generatedAt).toLocaleString()}
              </span>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
