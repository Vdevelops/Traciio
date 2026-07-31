"use client";

import React from "react";
import type { TeamSummary } from "../types";

interface Props {
  summary: TeamSummary | null;
}

export default function TeamSummaryCard({ summary }: Props) {
  if (!summary) return null;
  const fmt = (v: number | null) => (v === null ? <span className="text-sm text-gray-500">Belum ada data</span> : <span className="text-lg font-semibold">{v.toLocaleString()}</span>);

  return (
    <div className="p-4 bg-white rounded-md shadow-sm grid grid-cols-3 gap-4">
      <div>
        <div className="text-xs text-gray-400">Total Reps</div>
        <div className="mt-1">{fmt(summary.totalRepsCount)}</div>
      </div>
      <div>
        <div className="text-xs text-gray-400">Total Deals Closed</div>
        <div className="mt-1">{fmt(summary.totalDealsClosed)}</div>
      </div>
      <div>
        <div className="text-xs text-gray-400">Total Revenue</div>
        <div className="mt-1">{fmt(summary.totalRevenue)}</div>
      </div>
      <div>
        <div className="text-xs text-gray-400">Team Conversion Rate</div>
        <div className="mt-1">{summary.teamConversionRate === null ? <span className="text-sm text-gray-500">Belum ada data</span> : <span className="text-lg font-semibold">{summary.teamConversionRate.toFixed(1)}%</span>}</div>
      </div>
      <div>
        <div className="text-xs text-gray-400">Team Visit Compliance</div>
        <div className="mt-1">{summary.teamVisitCompliance === null ? <span className="text-sm text-gray-500">Belum ada data</span> : <span className="text-lg font-semibold">{summary.teamVisitCompliance.toFixed(1)}%</span>}</div>
      </div>
      <div>
        <div className="text-xs text-gray-400">Team Overdue Task Rate</div>
        <div className="mt-1">{summary.teamOverdueTaskRate === null ? <span className="text-sm text-gray-500">Belum ada data</span> : <span className="text-lg font-semibold">{summary.teamOverdueTaskRate.toFixed(1)}%</span>}</div>
      </div>
    </div>
  );
}
