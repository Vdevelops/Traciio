"use client";

import React from "react";
import { useState } from "react";
import { useSalesRepKPI } from "@/features/kpi/hooks/useKPIHooks";
import CompositeScoreCard from "@/features/kpi/components/CompositeScoreCard";
import MetricCard from "@/features/kpi/components/MetricCard";
import DiagnosticList from "@/features/kpi/components/DiagnosticList";
import type { SalesRepKPIParams } from "@/features/kpi/types";

export default function SalesRepKPIPageClient() {
  const [startDate, setStartDate] = useState(() => {
    const d = new Date();
    d.setDate(1);
    return d.toISOString().slice(0, 10);
  });
  const [endDate, setEndDate] = useState(() => {
    const d = new Date();
    d.setMonth(d.getMonth() + 1);
    d.setDate(0);
    return d.toISOString().slice(0, 10);
  });

  const params: SalesRepKPIParams = { startDate, endDate, compareWithPrevious: true };
  const { data, isLoading, error } = useSalesRepKPI(params, { enabled: true });

  if (isLoading) return <div>Loading KPI...</div>;
  if (error) return <div>Gagal memuat KPI. Silakan coba lagi.</div>;

  const resp = data as any; // runtime shape guaranteed by API/BE

  return (
    <div className="space-y-4">
      <CompositeScoreCard score={resp?.evaluation?.compositeScore ?? null} grade={resp?.evaluation?.grade ?? null} trendDelta={resp?.evaluation?.trend?.delta ?? null} />

      <div className="grid grid-cols-3 gap-4">
        <MetricCard label="Total Deals Closed" value={resp?.scorecard?.totalDealsClosed ?? null} />
        <MetricCard label="Total Revenue" value={resp?.scorecard?.totalRevenue ?? null} />
        <MetricCard label="Deals Created" value={resp?.scorecard?.dealsCreated ?? null} />

        <MetricCard label="Conversion Rate" value={resp?.scorecard?.conversionRate ?? null} suffix="%" />
        <MetricCard label="Average Deal Value" value={resp?.scorecard?.averageDealValue ?? null} />
        <MetricCard label="Pipeline Movement Score" value={resp?.scorecard?.pipelineMovementScore ?? null} />

        <MetricCard label="Visit Completed" value={resp?.scorecard?.visitCompleted ?? null} />
        <MetricCard label="Visit Planned" value={resp?.scorecard?.visitPlanned ?? null} />
        <MetricCard label="Visit Compliance" value={resp?.scorecard?.visitCompliance ?? null} suffix="%" />

        <MetricCard label="Tasks Completed" value={resp?.scorecard?.tasksCompleted ?? null} />
        <MetricCard label="Overdue Task Rate" value={resp?.scorecard?.overdueTaskRate ?? null} suffix="%" />
        <MetricCard label="Revenue Target Attainment" value={resp?.scorecard?.revenueTargetAttainment ?? null} suffix="%" />
      </div>

      <DiagnosticList items={resp?.diagnostics ?? []} />
    </div>
  );
}
