"use client";

import React from "react";
import type { KPIGrade } from "../types";

interface Props {
  score: number | null;
  grade: KPIGrade | null;
  trendDelta?: number | null;
}

export default function CompositeScoreCard({ score, grade, trendDelta }: Props) {
  return (
    <section aria-labelledby="kpi-composite-title" className="p-4 bg-white rounded-md shadow-sm" role="region">
      <div className="flex items-center justify-between">
        <div>
          <h2 id="kpi-composite-title" className="text-sm text-gray-500">Composite Score</h2>
          {score === null ? (
            <div className="text-2xl font-semibold" tabIndex={0}>Tidak cukup data untuk evaluasi</div>
          ) : (
            <div className="flex items-baseline gap-3">
              <div className="text-4xl font-bold" aria-live="polite">{score.toFixed(1)}</div>
              <div className="px-2 py-1 rounded-full text-sm bg-gray-100" role="status" aria-label={`Grade ${grade}`}>{grade}</div>
            </div>
          )}
        </div>
        <div>
          {trendDelta === undefined || trendDelta === null ? null : (
            <div className="text-sm text-gray-600" aria-hidden={false}>
              {trendDelta > 0 ? `Naik ${trendDelta.toFixed(1)} poin` : trendDelta < 0 ? `Turun ${Math.abs(trendDelta).toFixed(1)} poin` : `Flat`}
            </div>
          )}
        </div>
      </div>
    </section>
  );
}
