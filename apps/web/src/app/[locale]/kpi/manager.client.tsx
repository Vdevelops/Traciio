"use client";

import React, { useState } from "react";
import { useSalesManagerKPI } from "@/features/kpi/hooks/useKPIHooks";
import TeamSummaryCard from "@/features/kpi/components/TeamSummaryCard";
import CompositeScoreCard from "@/features/kpi/components/CompositeScoreCard";
import DiagnosticList from "@/features/kpi/components/DiagnosticList";

export default function SalesManagerKPIPageClient() {
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

  const params = { startDate, endDate, includeTeamBreakdown: true };
  const { data, isLoading, error } = useSalesManagerKPI(params as any);

  if (isLoading) return <div>Loading KPI...</div>;
  if (error) return <div>Gagal memuat KPI. Silakan coba lagi.</div>;

  const resp = data as any;

  return (
    <div className="space-y-4">
      <CompositeScoreCard score={resp?.evaluation?.compositeScore ?? null} grade={resp?.evaluation?.grade ?? null} trendDelta={resp?.evaluation?.trend?.delta ?? null} />

      <TeamSummaryCard summary={resp?.teamSummary ?? null} />

      <div className="grid grid-cols-2 gap-4">
        <div>
          <div className="text-sm font-semibold mb-2">Team Breakdown</div>
          {/* simple table */}
          <div className="bg-white rounded-md shadow-sm p-2">
            {resp?.teamBreakdown?.map((t: any) => (
              <div key={t.userId} className="flex justify-between py-2 border-b last:border-b-0">
                <div>{t.name}</div>
                <div className="text-sm">{t.compositeScore?.toFixed(1)} • {t.grade}</div>
              </div>
            ))}
          </div>
        </div>
        <div>
          <div className="text-sm font-semibold mb-2">Brick Breakdown</div>
          <div className="bg-white rounded-md shadow-sm p-2">
            {resp?.brickBreakdown?.map((b: any) => (
              <div key={b.brickId} className="flex justify-between py-2 border-b last:border-b-0">
                <div>{b.name}</div>
                <div className="text-sm">{b.compositeScore?.toFixed(1)}</div>
              </div>
            ))}
          </div>
        </div>
      </div>

      <DiagnosticList items={resp?.diagnostics ?? []} />
    </div>
  );
}
