"use client";

import React from "react";
import type { KPIGrade, TrendInfo } from "../types";
import { gradeColor, trendColor, getTrendIcon } from "../utils/formatters";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface Props {
  score: number | null;
  grade: KPIGrade | null;
  trend: TrendInfo | null;
}

export default function CompositeScoreCard({ score, grade, trend }: Props) {
  const isNoData = score === null || score === undefined;
  const colors = gradeColor(grade);

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center gap-2">
          <div className="h-1.5 w-1.5 rounded-full bg-primary" />
          <CardTitle className="text-sm font-medium uppercase tracking-wide text-muted-foreground">
            Weighted Composite Score
          </CardTitle>
        </div>
      </CardHeader>

      <CardContent>
        <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div>
            {isNoData ? (
              <div className="text-base font-medium text-muted-foreground" tabIndex={0}>
                Tidak cukup data untuk evaluasi
              </div>
            ) : (
              <div className="flex items-center gap-4">
                <span
                  className="text-4xl font-semibold tracking-tight text-foreground sm:text-5xl"
                  aria-live="polite"
                >
                  {score.toFixed(1)}
                </span>
                <span
                  className={`inline-flex items-center rounded-md border px-2.5 py-1 text-xs font-semibold ${colors.bg} ${colors.text} ${colors.border}`}
                  role="status"
                  aria-label={`Grade ${grade}`}
                >
                  {grade}
                </span>
              </div>
            )}
          </div>

          {!isNoData && trend && (
            <div className="flex flex-col items-start sm:items-end">
              <span className="text-xs font-medium text-muted-foreground">
                Tren vs Periode Lalu
              </span>
              <div className="mt-1.5 flex items-center gap-1.5">
                {(() => {
                  const Icon = getTrendIcon(trend.direction);
                  const color = trendColor(trend.direction);
                  const deltaText = Math.abs(trend.delta ?? 0).toFixed(1);

                  let label = "";
                  if (trend.direction === "up") {
                    label = `naik ${deltaText} poin dari periode lalu`;
                  } else if (trend.direction === "down") {
                    label = `turun ${deltaText} poin dari periode lalu`;
                  } else {
                    label = "stabil dari periode lalu";
                  }

                  return (
                    <span className={`inline-flex items-center gap-1 text-sm font-medium ${color}`}>
                      <Icon className="h-4 w-4" />
                      <span>{label}</span>
                    </span>
                  );
                })()}
              </div>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
